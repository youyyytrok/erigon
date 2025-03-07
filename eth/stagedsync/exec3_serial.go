package stagedsync

import (
	"context"
	"errors"
	"fmt"
	"time"

	chaos_monkey "github.com/erigontech/erigon/tests/chaos-monkey"

	"github.com/erigontech/erigon-lib/log/v3"
	state2 "github.com/erigontech/erigon-lib/state"
	"github.com/erigontech/erigon/consensus"
	"github.com/erigontech/erigon/core"
	"github.com/erigontech/erigon/core/rawdb/rawtemporaldb"
	"github.com/erigontech/erigon/core/state"
	"github.com/erigontech/erigon/core/types"
)

type serialExecutor struct {
	txExecutor
	skipPostEvaluation bool
	// outputs
	txCount     uint64
	usedGas     uint64
	blobGasUsed uint64
}

func (se *serialExecutor) wait() error {
	return nil
}

func (se *serialExecutor) status(ctx context.Context, commitThreshold uint64) error {
	return nil
}

func (se *serialExecutor) execute(ctx context.Context, tasks []*state.TxTask) (cont bool, err error) {
	for _, txTask := range tasks {
		if txTask.Error != nil {
			return false, nil
		}

		se.applyWorker.RunTxTaskNoLock(txTask, se.isMining, se.skipPostEvaluation)
		if err := func() error {
			if errors.Is(txTask.Error, context.Canceled) {
				return txTask.Error
			}
			if txTask.Error != nil {
				return fmt.Errorf("%w, txnIdx=%d, %v", consensus.ErrInvalidBlock, txTask.TxIndex, txTask.Error) //same as in stage_exec.go
			}

			se.txCount++
			se.usedGas += txTask.UsedGas
			mxExecGas.Add(float64(txTask.UsedGas))
			mxExecTransactions.Add(1)

			if txTask.Tx != nil {
				se.blobGasUsed += txTask.Tx.GetBlobGas()
			}

			txTask.CreateReceipt(se.applyTx)

			if txTask.Final {
				if !se.isMining && !se.inMemExec && !se.skipPostEvaluation && !se.execStage.CurrentSyncCycle.IsInitialCycle {
					// note this assumes the bloach reciepts is a fixed array shared by
					// all tasks - if that changes this will need to change - robably need to
					// add this to the executor
					se.cfg.notifications.RecentLogs.Add(txTask.BlockReceipts)
				}
				checkReceipts := !se.cfg.vmConfig.StatelessExec && se.cfg.chainConfig.IsByzantium(txTask.BlockNum) && !se.cfg.vmConfig.NoReceipts && !se.isMining
				if txTask.BlockNum > 0 && !se.skipPostEvaluation { //Disable check for genesis. Maybe need somehow improve it in future - to satisfy TestExecutionSpec
					if err := core.BlockPostValidation(se.usedGas, se.blobGasUsed, checkReceipts, txTask.BlockReceipts, txTask.Header, se.isMining, txTask.Txs, se.cfg.chainConfig, se.logger); err != nil {
						return fmt.Errorf("%w, txnIdx=%d, %v", consensus.ErrInvalidBlock, txTask.TxIndex, err) //same as in stage_exec.go
					}
				}

				se.outputBlockNum.SetUint64(txTask.BlockNum)
			}
			if se.cfg.syncCfg.ChaosMonkey {
				chaosErr := chaos_monkey.ThrowRandomConsensusError(se.execStage.CurrentSyncCycle.IsInitialCycle, txTask.TxIndex, se.cfg.badBlockHalt, txTask.Error)
				if chaosErr != nil {
					log.Warn("Monkey in a consensus")
					return chaosErr
				}
			}
			return nil
		}(); err != nil {
			if errors.Is(err, context.Canceled) {
				return false, err
			}
			se.logger.Warn(fmt.Sprintf("[%s] Execution failed", se.execStage.LogPrefix()),
				"block", txTask.BlockNum, "txNum", txTask.TxNum, "hash", txTask.Header.Hash().String(), "err", err, "inMem", se.inMemExec)
			if se.cfg.hd != nil && se.cfg.hd.POSSync() && errors.Is(err, consensus.ErrInvalidBlock) {
				se.cfg.hd.ReportBadHeaderPoS(txTask.Header.Hash(), txTask.Header.ParentHash)
			}
			if se.cfg.badBlockHalt {
				return false, err
			}
			if errors.Is(err, consensus.ErrInvalidBlock) {
				if se.u != nil {
					if err := se.u.UnwindTo(txTask.BlockNum-1, BadBlock(txTask.Header.Hash(), err), se.applyTx); err != nil {
						return false, err
					}
				}
			} else {
				if se.u != nil {
					if err := se.u.UnwindTo(txTask.BlockNum-1, ExecUnwind, se.applyTx); err != nil {
						return false, err
					}
				}
			}
			return false, nil
		}

		if !txTask.Final {
			var receipt *types.Receipt
			if txTask.TxIndex >= 0 {
				receipt = txTask.BlockReceipts[txTask.TxIndex]
			}
			if err := rawtemporaldb.AppendReceipt(se.doms, receipt, se.blobGasUsed); err != nil {
				return false, err
			}
		} else {
			if se.cfg.chainConfig.Bor != nil && txTask.TxIndex >= 1 {
				// get last receipt and store the last log index + 1
				lastReceipt := txTask.BlockReceipts[txTask.TxIndex-1]
				if len(lastReceipt.Logs) > 0 {
					firstIndex := lastReceipt.Logs[len(lastReceipt.Logs)-1].Index + 1
					receipt := types.Receipt{
						CumulativeGasUsed:        lastReceipt.CumulativeGasUsed,
						FirstLogIndexWithinBlock: uint32(firstIndex),
					}
					if err := rawtemporaldb.AppendReceipt(se.doms, &receipt, se.blobGasUsed); err != nil {
						return false, err
					}
				}
			}
		}

		// MA applystate
		if err := se.rs.ApplyState4(ctx, txTask); err != nil {
			return false, err
		}

		se.outputTxNum.Add(1)
	}

	return true, nil
}

func (se *serialExecutor) commit(ctx context.Context, txNum uint64, blockNum uint64, useExternalTx bool) (t2 time.Duration, err error) {
	se.doms.Close()
	if err = se.execStage.Update(se.applyTx, blockNum); err != nil {
		return 0, err
	}

	se.applyTx.CollectMetrics()

	if !useExternalTx {
		tt := time.Now()
		if err = se.applyTx.Commit(); err != nil {
			return 0, err
		}

		t2 = time.Since(tt)
		se.agg.BuildFilesInBackground(se.outputTxNum.Load())

		se.applyTx, err = se.cfg.db.BeginRw(context.Background()) //nolint
		if err != nil {
			return t2, err
		}
	}
	se.doms, err = state2.NewSharedDomains(se.applyTx, se.logger)
	if err != nil {
		return t2, err
	}
	se.doms.SetTxNum(txNum)
	se.rs = state.NewStateV3(se.doms, se.logger)

	se.applyWorker.ResetTx(se.applyTx)
	se.applyWorker.ResetState(se.rs, se.accumulator)

	return t2, nil
}

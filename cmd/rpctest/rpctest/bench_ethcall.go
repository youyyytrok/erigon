// Copyright 2024 The Erigon Authors
// This file is part of Erigon.
//
// Erigon is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Erigon is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Erigon. If not, see <http://www.gnu.org/licenses/>.

package rpctest

import (
	"bufio"
	"fmt"
	"os"
)

// BenchEthCall compares response of Erigon with Geth
// but also can be used for comparing RPCDaemon with Geth or infura
// parameters:
// needCompare - if false - doesn't call Erigon and doesn't compare responses
//
//			    false value - to generate vegeta files, it's faster but we can generate vegeta files for Geth and Erigon
//	                 recordFile stores all eth_call returned with success
//	                 errorFile stores information when erigon and geth doesn't return same data
func BenchEthCall(erigonURL, gethURL string, needCompare, latest bool, blockFrom, blockTo uint64, recordFileName string, errorFileName string) error {
	setRoutes(erigonURL, gethURL)

	var rec *bufio.Writer
	var errs *bufio.Writer
	var resultsCh chan CallResult = nil
	var nTransactions = 0

	if errorFileName != "" {
		f, err := os.Create(errorFileName)
		if err != nil {
			return fmt.Errorf("Cannot create file %s for errorFile: %v\n", errorFileName, err)
		}
		defer f.Close()
		errs = bufio.NewWriter(f)
		defer errs.Flush()
	}

	if recordFileName != "" {
		frec, errRec := os.Create(recordFileName)
		if errRec != nil {
			return fmt.Errorf("Cannot create file %s for errorFile: %v\n", recordFileName, errRec)
		}
		defer frec.Close()
		rec = bufio.NewWriter(frec)
		defer rec.Flush()
	}

	if !needCompare {
		resultsCh = make(chan CallResult, 1000)
		defer close(resultsCh)
		go vegetaWrite(true, []string{"eth_call"}, resultsCh)
	}
	var res CallResult

	reqGen := &RequestGenerator{}

	for bn := blockFrom; bn <= blockTo; bn++ {

		var b EthBlockByNumber
		res = reqGen.Erigon("eth_getBlockByNumber", reqGen.getBlockByNumber(bn, true /* withTxs */), &b)
		if res.Err != nil {
			return fmt.Errorf("Could not retrieve block (Erigon) %d: %v\n", bn, res.Err)
		}

		if b.Error != nil {
			return fmt.Errorf("Error retrieving block (Erigon): %d %s\n", b.Error.Code, b.Error.Message)
		}

		if needCompare {
			var bg EthBlockByNumber
			res = reqGen.Geth("eth_getBlockByNumber", reqGen.getBlockByNumber(bn, true /* withTxs */), &bg)
			if res.Err != nil {
				return fmt.Errorf("Could not retrieve block (geth) %d: %v\n", bn, res.Err)
			}
			if bg.Error != nil {
				return fmt.Errorf("Error retrieving block (geth): %d %s\n", bg.Error.Code, bg.Error.Message)
			}
			if !compareBlocks(&b, &bg) {
				if rec != nil {
					fmt.Fprintf(rec, "Block difference for block=%d\n", bn)
					rec.Flush()
					continue
				} else {
					return fmt.Errorf("Block one or more fields areis different for block %d\n", bn)
				}
			}
		}

		for _, txn := range b.Result.Transactions {

			nTransactions = nTransactions + 1

			var request string
			var insertedOnlyIfSuccess bool
			if latest {
				request = reqGen.ethCallLatest(txn.From, txn.To, &txn.Gas, &txn.GasPrice, &txn.Value, txn.Input)
				insertedOnlyIfSuccess = true
			} else {
				request = reqGen.ethCall(txn.From, txn.To, &txn.Gas, &txn.GasPrice, &txn.Value, txn.Input, bn-1)
				insertedOnlyIfSuccess = false
			}
			errCtx := fmt.Sprintf(" bn=%d hash=%s", bn, txn.Hash)

			if err := requestAndCompare(request, "eth_call", errCtx, reqGen, needCompare, rec, errs, resultsCh,
				insertedOnlyIfSuccess); err != nil {
				return err
			}
		}

		fmt.Println("\nProcessed Transactions: ", nTransactions)
	}
	return nil
}

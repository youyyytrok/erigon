// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/erigontech/erigon/cl/beacon/builder (interfaces: BuilderClient)
//
// Generated by this command:
//
//	mockgen -typed=true -destination=./mock_services/builder_client_mock.go -package=mock_services . BuilderClient
//

// Package mock_services is a generated GoMock package.
package mock_services

import (
	context "context"
	reflect "reflect"

	common "github.com/erigontech/erigon-lib/common"
	builder "github.com/erigontech/erigon/cl/beacon/builder"
	cltypes "github.com/erigontech/erigon/cl/cltypes"
	engine_types "github.com/erigontech/erigon/turbo/engineapi/engine_types"
	gomock "go.uber.org/mock/gomock"
)

// MockBuilderClient is a mock of BuilderClient interface.
type MockBuilderClient struct {
	ctrl     *gomock.Controller
	recorder *MockBuilderClientMockRecorder
	isgomock struct{}
}

// MockBuilderClientMockRecorder is the mock recorder for MockBuilderClient.
type MockBuilderClientMockRecorder struct {
	mock *MockBuilderClient
}

// NewMockBuilderClient creates a new mock instance.
func NewMockBuilderClient(ctrl *gomock.Controller) *MockBuilderClient {
	mock := &MockBuilderClient{ctrl: ctrl}
	mock.recorder = &MockBuilderClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBuilderClient) EXPECT() *MockBuilderClientMockRecorder {
	return m.recorder
}

// GetHeader mocks base method.
func (m *MockBuilderClient) GetHeader(ctx context.Context, slot int64, parentHash common.Hash, pubKey common.Bytes48) (*builder.ExecutionHeader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHeader", ctx, slot, parentHash, pubKey)
	ret0, _ := ret[0].(*builder.ExecutionHeader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHeader indicates an expected call of GetHeader.
func (mr *MockBuilderClientMockRecorder) GetHeader(ctx, slot, parentHash, pubKey any) *MockBuilderClientGetHeaderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHeader", reflect.TypeOf((*MockBuilderClient)(nil).GetHeader), ctx, slot, parentHash, pubKey)
	return &MockBuilderClientGetHeaderCall{Call: call}
}

// MockBuilderClientGetHeaderCall wrap *gomock.Call
type MockBuilderClientGetHeaderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockBuilderClientGetHeaderCall) Return(arg0 *builder.ExecutionHeader, arg1 error) *MockBuilderClientGetHeaderCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockBuilderClientGetHeaderCall) Do(f func(context.Context, int64, common.Hash, common.Bytes48) (*builder.ExecutionHeader, error)) *MockBuilderClientGetHeaderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockBuilderClientGetHeaderCall) DoAndReturn(f func(context.Context, int64, common.Hash, common.Bytes48) (*builder.ExecutionHeader, error)) *MockBuilderClientGetHeaderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetStatus mocks base method.
func (m *MockBuilderClient) GetStatus(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatus", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetStatus indicates an expected call of GetStatus.
func (mr *MockBuilderClientMockRecorder) GetStatus(ctx any) *MockBuilderClientGetStatusCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatus", reflect.TypeOf((*MockBuilderClient)(nil).GetStatus), ctx)
	return &MockBuilderClientGetStatusCall{Call: call}
}

// MockBuilderClientGetStatusCall wrap *gomock.Call
type MockBuilderClientGetStatusCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockBuilderClientGetStatusCall) Return(arg0 error) *MockBuilderClientGetStatusCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockBuilderClientGetStatusCall) Do(f func(context.Context) error) *MockBuilderClientGetStatusCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockBuilderClientGetStatusCall) DoAndReturn(f func(context.Context) error) *MockBuilderClientGetStatusCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// RegisterValidator mocks base method.
func (m *MockBuilderClient) RegisterValidator(ctx context.Context, registers []*cltypes.ValidatorRegistration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterValidator", ctx, registers)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterValidator indicates an expected call of RegisterValidator.
func (mr *MockBuilderClientMockRecorder) RegisterValidator(ctx, registers any) *MockBuilderClientRegisterValidatorCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterValidator", reflect.TypeOf((*MockBuilderClient)(nil).RegisterValidator), ctx, registers)
	return &MockBuilderClientRegisterValidatorCall{Call: call}
}

// MockBuilderClientRegisterValidatorCall wrap *gomock.Call
type MockBuilderClientRegisterValidatorCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockBuilderClientRegisterValidatorCall) Return(arg0 error) *MockBuilderClientRegisterValidatorCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockBuilderClientRegisterValidatorCall) Do(f func(context.Context, []*cltypes.ValidatorRegistration) error) *MockBuilderClientRegisterValidatorCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockBuilderClientRegisterValidatorCall) DoAndReturn(f func(context.Context, []*cltypes.ValidatorRegistration) error) *MockBuilderClientRegisterValidatorCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SubmitBlindedBlocks mocks base method.
func (m *MockBuilderClient) SubmitBlindedBlocks(ctx context.Context, block *cltypes.SignedBlindedBeaconBlock) (*cltypes.Eth1Block, *engine_types.BlobsBundleV1, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitBlindedBlocks", ctx, block)
	ret0, _ := ret[0].(*cltypes.Eth1Block)
	ret1, _ := ret[1].(*engine_types.BlobsBundleV1)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SubmitBlindedBlocks indicates an expected call of SubmitBlindedBlocks.
func (mr *MockBuilderClientMockRecorder) SubmitBlindedBlocks(ctx, block any) *MockBuilderClientSubmitBlindedBlocksCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitBlindedBlocks", reflect.TypeOf((*MockBuilderClient)(nil).SubmitBlindedBlocks), ctx, block)
	return &MockBuilderClientSubmitBlindedBlocksCall{Call: call}
}

// MockBuilderClientSubmitBlindedBlocksCall wrap *gomock.Call
type MockBuilderClientSubmitBlindedBlocksCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockBuilderClientSubmitBlindedBlocksCall) Return(arg0 *cltypes.Eth1Block, arg1 *engine_types.BlobsBundleV1, arg2 error) *MockBuilderClientSubmitBlindedBlocksCall {
	c.Call = c.Call.Return(arg0, arg1, arg2)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockBuilderClientSubmitBlindedBlocksCall) Do(f func(context.Context, *cltypes.SignedBlindedBeaconBlock) (*cltypes.Eth1Block, *engine_types.BlobsBundleV1, error)) *MockBuilderClientSubmitBlindedBlocksCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockBuilderClientSubmitBlindedBlocksCall) DoAndReturn(f func(context.Context, *cltypes.SignedBlindedBeaconBlock) (*cltypes.Eth1Block, *engine_types.BlobsBundleV1, error)) *MockBuilderClientSubmitBlindedBlocksCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

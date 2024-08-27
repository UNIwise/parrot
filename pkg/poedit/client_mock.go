// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package poedit is a generated GoMock package.
package poedit

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// ExportProject mocks base method.
func (m *MockClient) ExportProject(ctx context.Context, req ExportProjectRequest) (*ExportProjectResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExportProject", ctx, req)
	ret0, _ := ret[0].(*ExportProjectResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExportProject indicates an expected call of ExportProject.
func (mr *MockClientMockRecorder) ExportProject(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExportProject", reflect.TypeOf((*MockClient)(nil).ExportProject), ctx, req)
}

// ListProjectLanguages mocks base method.
func (m *MockClient) ListProjectLanguages(ctx context.Context, r ListProjectLanguagesRequest) (*ListProjectLanguagesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProjectLanguages", ctx, r)
	ret0, _ := ret[0].(*ListProjectLanguagesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProjectLanguages indicates an expected call of ListProjectLanguages.
func (mr *MockClientMockRecorder) ListProjectLanguages(ctx, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProjectLanguages", reflect.TypeOf((*MockClient)(nil).ListProjectLanguages), ctx, r)
}

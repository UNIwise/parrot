// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
// Generated by this command:
//
//	mockgen --source=service.go -destination=service_mock.go -package=project
//

// Package project is a generated GoMock package.
package project

import (
	context "context"
	reflect "reflect"

	go_sundheit "github.com/AppsFlyer/go-sundheit"
	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GetAllProjects mocks base method.
func (m *MockService) GetAllProjects(ctx context.Context) ([]Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllProjects", ctx)
	ret0, _ := ret[0].([]Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllProjects indicates an expected call of GetAllProjects.
func (mr *MockServiceMockRecorder) GetAllProjects(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllProjects", reflect.TypeOf((*MockService)(nil).GetAllProjects), ctx)
}

// GetProjectByID mocks base method.
func (m *MockService) GetProjectByID(ctx context.Context, id int) (*Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectByID", ctx, id)
	ret0, _ := ret[0].(*Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectByID indicates an expected call of GetProjectByID.
func (mr *MockServiceMockRecorder) GetProjectByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectByID", reflect.TypeOf((*MockService)(nil).GetProjectByID), ctx, id)
}

// GetProjectVersions mocks base method.
func (m *MockService) GetProjectVersions(ctx context.Context, projectID int) ([]Version, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectVersions", ctx, projectID)
	ret0, _ := ret[0].([]Version)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectVersions indicates an expected call of GetProjectVersions.
func (mr *MockServiceMockRecorder) GetProjectVersions(ctx, projectID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectVersions", reflect.TypeOf((*MockService)(nil).GetProjectVersions), ctx, projectID)
}

// GetTranslation mocks base method.
func (m *MockService) GetTranslation(ctx context.Context, projectID int, languageCode, format string) (*Translation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTranslation", ctx, projectID, languageCode, format)
	ret0, _ := ret[0].(*Translation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTranslation indicates an expected call of GetTranslation.
func (mr *MockServiceMockRecorder) GetTranslation(ctx, projectID, languageCode, format interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTranslation", reflect.TypeOf((*MockService)(nil).GetTranslation), ctx, projectID, languageCode, format)
}

// PurgeProject mocks base method.
func (m *MockService) PurgeProject(ctx context.Context, projectID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PurgeProject", ctx, projectID)
	ret0, _ := ret[0].(error)
	return ret0
}

// PurgeProject indicates an expected call of PurgeProject.
func (mr *MockServiceMockRecorder) PurgeProject(ctx, projectID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PurgeProject", reflect.TypeOf((*MockService)(nil).PurgeProject), ctx, projectID)
}

// PurgeTranslation mocks base method.
func (m *MockService) PurgeTranslation(ctx context.Context, projectID int, languageCode string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PurgeTranslation", ctx, projectID, languageCode)
	ret0, _ := ret[0].(error)
	return ret0
}

// PurgeTranslation indicates an expected call of PurgeTranslation.
func (mr *MockServiceMockRecorder) PurgeTranslation(ctx, projectID, languageCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PurgeTranslation", reflect.TypeOf((*MockService)(nil).PurgeTranslation), ctx, projectID, languageCode)
}

// RegisterChecks mocks base method.
func (m *MockService) RegisterChecks(h go_sundheit.Health) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterChecks", h)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterChecks indicates an expected call of RegisterChecks.
func (mr *MockServiceMockRecorder) RegisterChecks(h interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterChecks", reflect.TypeOf((*MockService)(nil).RegisterChecks), h)
}

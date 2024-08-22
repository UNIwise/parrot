//
// Generated by this command:
//
//	mockgen --source=repository.go -destination=repository_mock.go -package=project
//

// Package project is a generated GoMock package.
package project

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetAllProjects mocks base method.
func (m *MockRepository) GetAllProjects(ctx context.Context) ([]Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllProjects", ctx)
	ret0, _ := ret[0].([]Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllProjects indicates an expected call of GetAllProjects.
func (mr *MockRepositoryMockRecorder) GetAllProjects(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllProjects", reflect.TypeOf((*MockRepository)(nil).GetAllProjects), ctx)
}

// GetProjectByID mocks base method.
func (m *MockRepository) GetProjectByID(ctx context.Context, id int) (*Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectByID", ctx, id)
	ret0, _ := ret[0].(*Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectByID indicates an expected call of GetProjectByID.
func (mr *MockRepositoryMockRecorder) GetProjectByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectByID", reflect.TypeOf((*MockRepository)(nil).GetProjectByID), ctx, id)
}

// GetProjectVersions mocks base method.
func (m *MockRepository) GetProjectVersions(ctx context.Context, projectID int) ([]Version, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectVersions", ctx, projectID)
	ret0, _ := ret[0].([]Version)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectVersions indicates an expected call of GetProjectVersions.
func (mr *MockRepositoryMockRecorder) GetProjectVersions(ctx, projectID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectVersions", reflect.TypeOf((*MockRepository)(nil).GetProjectVersions), ctx, projectID)
}

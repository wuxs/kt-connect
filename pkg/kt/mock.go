// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/kt/types.go

// Package kt is a generated GoMock package.
package kt

import (
	reflect "reflect"

	cluster "github.com/alibaba/kt-connect/pkg/kt/cluster"
	connect "github.com/alibaba/kt-connect/pkg/kt/tunnel"
	exec "github.com/alibaba/kt-connect/pkg/kt/exec"
	gomock "github.com/golang/mock/gomock"
)

// MockCliInterface is a mock of CliInterface interface.
type MockCliInterface struct {
	ctrl     *gomock.Controller
	recorder *MockCliInterfaceMockRecorder
}

// MockCliInterfaceMockRecorder is the mock recorder for MockCliInterface.
type MockCliInterfaceMockRecorder struct {
	mock *MockCliInterface
}

// NewMockCliInterface creates a new mock instance.
func NewMockCliInterface(ctrl *gomock.Controller) *MockCliInterface {
	mock := &MockCliInterface{ctrl: ctrl}
	mock.recorder = &MockCliInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCliInterface) EXPECT() *MockCliInterfaceMockRecorder {
	return m.recorder
}

// Exec mocks base method.
func (m *MockCliInterface) Exec() exec.CliInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec")
	ret0, _ := ret[0].(exec.CliInterface)
	return ret0
}

// Exec indicates an expected call of Exec.
func (mr *MockCliInterfaceMockRecorder) Exec() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockCliInterface)(nil).Exec))
}

// Kubernetes mocks base method.
func (m *MockCliInterface) Kubernetes() (cluster.KubernetesInterface, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Kubernetes")
	ret0, _ := ret[0].(cluster.KubernetesInterface)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Kubernetes indicates an expected call of Kubernetes.
func (mr *MockCliInterfaceMockRecorder) Kubernetes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Kubernetes", reflect.TypeOf((*MockCliInterface)(nil).Kubernetes))
}

// Shadow mocks base method.
func (m *MockCliInterface) Shadow() connect.ShadowInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shadow")
	ret0, _ := ret[0].(connect.ShadowInterface)
	return ret0
}

// Shadow indicates an expected call of Shadow.
func (mr *MockCliInterfaceMockRecorder) Shadow() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shadow", reflect.TypeOf((*MockCliInterface)(nil).Shadow))
}

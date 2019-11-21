// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	statefile "github.com/hashicorp/terraform/states/statefile"

	types "github.com/kyma-incubator/hydroform/types"
)

// Operator is an autogenerated mock type for the Operator type
type Operator struct {
	mock.Mock
}

// Create provides a mock function with given fields: p, cfg
func (_m *Operator) Create(p types.ProviderType, cfg map[string]interface{}) (*types.ClusterInfo, error) {
	ret := _m.Called(p, cfg)

	var r0 *types.ClusterInfo
	if rf, ok := ret.Get(0).(func(types.ProviderType, map[string]interface{}) *types.ClusterInfo); ok {
		r0 = rf(p, cfg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ClusterInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.ProviderType, map[string]interface{}) error); ok {
		r1 = rf(p, cfg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: state, p, cfg
func (_m *Operator) Delete(state *statefile.File, p types.ProviderType, cfg map[string]interface{}) error {
	ret := _m.Called(state, p, cfg)

	var r0 error
	if rf, ok := ret.Get(0).(func(*statefile.File, types.ProviderType, map[string]interface{}) error); ok {
		r0 = rf(state, p, cfg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Status provides a mock function with given fields: state, p, cfg
func (_m *Operator) Status(state *statefile.File, p types.ProviderType, cfg map[string]interface{}) (*types.ClusterStatus, error) {
	ret := _m.Called(state, p, cfg)

	var r0 *types.ClusterStatus
	if rf, ok := ret.Get(0).(func(*statefile.File, types.ProviderType, map[string]interface{}) *types.ClusterStatus); ok {
		r0 = rf(state, p, cfg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ClusterStatus)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*statefile.File, types.ProviderType, map[string]interface{}) error); ok {
		r1 = rf(state, p, cfg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

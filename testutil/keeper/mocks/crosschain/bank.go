// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// CrosschainBankKeeper is an autogenerated mock type for the CrosschainBankKeeper type
type CrosschainBankKeeper struct {
	mock.Mock
}

// BurnCoins provides a mock function with given fields: ctx, name, amt
func (_m *CrosschainBankKeeper) BurnCoins(ctx types.Context, name string, amt types.Coins) error {
	ret := _m.Called(ctx, name, amt)

	if len(ret) == 0 {
		panic("no return value specified for BurnCoins")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, string, types.Coins) error); ok {
		r0 = rf(ctx, name, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MintCoins provides a mock function with given fields: ctx, moduleName, amt
func (_m *CrosschainBankKeeper) MintCoins(ctx types.Context, moduleName string, amt types.Coins) error {
	ret := _m.Called(ctx, moduleName, amt)

	if len(ret) == 0 {
		panic("no return value specified for MintCoins")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, string, types.Coins) error); ok {
		r0 = rf(ctx, moduleName, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewCrosschainBankKeeper creates a new instance of CrosschainBankKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCrosschainBankKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *CrosschainBankKeeper {
	mock := &CrosschainBankKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

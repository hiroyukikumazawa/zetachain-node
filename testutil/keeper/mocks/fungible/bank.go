// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// FungibleBankKeeper is an autogenerated mock type for the FungibleBankKeeper type
type FungibleBankKeeper struct {
	mock.Mock
}

// MintCoins provides a mock function with given fields: ctx, moduleName, amt
func (_m *FungibleBankKeeper) MintCoins(ctx types.Context, moduleName string, amt types.Coins) error {
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

// SendCoinsFromModuleToAccount provides a mock function with given fields: ctx, senderModule, recipientAddr, amt
func (_m *FungibleBankKeeper) SendCoinsFromModuleToAccount(ctx types.Context, senderModule string, recipientAddr types.AccAddress, amt types.Coins) error {
	ret := _m.Called(ctx, senderModule, recipientAddr, amt)

	if len(ret) == 0 {
		panic("no return value specified for SendCoinsFromModuleToAccount")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, string, types.AccAddress, types.Coins) error); ok {
		r0 = rf(ctx, senderModule, recipientAddr, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewFungibleBankKeeper creates a new instance of FungibleBankKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFungibleBankKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *FungibleBankKeeper {
	mock := &FungibleBankKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

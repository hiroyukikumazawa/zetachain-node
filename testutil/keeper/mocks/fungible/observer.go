// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	chains "github.com/zeta-chain/zetacore/pkg/chains"

	types "github.com/cosmos/cosmos-sdk/types"
)

// FungibleObserverKeeper is an autogenerated mock type for the FungibleObserverKeeper type
type FungibleObserverKeeper struct {
	mock.Mock
}

// GetSupportedChains provides a mock function with given fields: ctx
func (_m *FungibleObserverKeeper) GetSupportedChains(ctx types.Context) []chains.Chain {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetSupportedChains")
	}

	var r0 []chains.Chain
	if rf, ok := ret.Get(0).(func(types.Context) []chains.Chain); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]chains.Chain)
		}
	}

	return r0
}

// NewFungibleObserverKeeper creates a new instance of FungibleObserverKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFungibleObserverKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *FungibleObserverKeeper {
	mock := &FungibleObserverKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

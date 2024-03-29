// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import account "github.com/wpe_merge/wpe_merge/account"
import context "context"
import mock "github.com/stretchr/testify/mock"

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// GetAccount provides a mock function with given fields: _a0, _a1
func (_m *Client) GetAccount(_a0 context.Context, _a1 *account.GetAccountRequest) (*account.Account, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *account.Account
	if rf, ok := ret.Get(0).(func(context.Context, *account.GetAccountRequest) *account.Account); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*account.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *account.GetAccountRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccounts provides a mock function with given fields: _a0, _a1
func (_m *Client) GetAccounts(_a0 context.Context, _a1 *account.GetAccountsRequest) (*account.GetAccountsResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *account.GetAccountsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *account.GetAccountsRequest) *account.GetAccountsResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*account.GetAccountsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *account.GetAccountsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

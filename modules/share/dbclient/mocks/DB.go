// Code generated by mockery 2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	sql "database/sql"
)

// DB is an autogenerated mock type for the DB type
type DB struct {
	mock.Mock
}

// BeginTx provides a mock function with given fields: ctx, opts
func (_m *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	ret := _m.Called(ctx, opts)

	var r0 *sql.Tx
	if rf, ok := ret.Get(0).(func(context.Context, *sql.TxOptions) *sql.Tx); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Tx)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *sql.TxOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryRowContext provides a mock function with given fields: ctx, query, args
func (_m *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	var _ca []interface{}
	_ca = append(_ca, ctx, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *sql.Row
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) *sql.Row); ok {
		r0 = rf(ctx, query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Row)
		}
	}

	return r0
}

// SelectContext provides a mock function with given fields: ctx, dest, query, args
func (_m *DB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, ctx, dest, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, string, ...interface{}) error); ok {
		r0 = rf(ctx, dest, query, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
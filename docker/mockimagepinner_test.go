// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package docker_test

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockImagePinner is an autogenerated mock type for the ImagePinner type
type mockImagePinner struct {
	mock.Mock
}

// Pin provides a mock function with given fields: ctx, image
func (_m *mockImagePinner) Pin(ctx context.Context, image string) (string, error) {
	ret := _m.Called(ctx, image)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, image)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, image)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Unpin provides a mock function with given fields: ctx, image, hash
func (_m *mockImagePinner) Unpin(ctx context.Context, image string, hash string) (string, error) {
	ret := _m.Called(ctx, image, hash)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, image, hash)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, image, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
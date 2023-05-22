package mediator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMediator_Publish(t *testing.T) {
	tests := []struct {
		name     string
		handlers []any
		args     []any
		wantErr  error
	}{
		{
			name: "successful publish with 1 argument handler",
			handlers: []any{
				func(first int) {},
			},
			args:    []any{1},
			wantErr: nil,
		},
		{
			name: "successful publish with 2 arguments handler",
			handlers: []any{
				func(first, second int) {},
			},
			args:    []any{1, 2},
			wantErr: nil,
		},
		{
			name: "mismatch signature (type)",
			handlers: []any{
				func(first int) {},
			},
			args:    []any{1.1},
			wantErr: ErrHandlerNotFound,
		},
		{
			name: "mismatch signature (count)",
			handlers: []any{
				func(first int) {},
			},
			args:    []any{1, 2},
			wantErr: ErrHandlerNotFound,
		},
		{
			name: "too many arguments",
			handlers: []any{
				func(first int) {},
			},
			args:    []any{1, 2, 3},
			wantErr: ErrTooManyArguments,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()

			for _, handler := range tt.handlers {
				err := m.Register(handler)
				require.NoError(t, err)
			}

			err := m.Publish(tt.args...)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestMediator_Register(t *testing.T) {
	tests := []struct {
		name    string
		handler any
		wantErr error
	}{
		{
			name:    "register handler with 1 argument",
			handler: func(first int) {},
			wantErr: nil,
		},
		{
			name:    "register handler with 2 arguments",
			handler: func(first, second int) {},
			wantErr: nil,
		},
		{
			name:    "too many arguments",
			handler: func(first, second, third int) {},
			wantErr: ErrTooManyArguments,
		},
		{
			name:    "no func handler",
			handler: struct{}{},
			wantErr: ErrPassedNotFunc,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			gotErr := m.Register(tt.handler)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}

	t.Run("try register duplicate signature handler", func(t *testing.T) {
		m := New()
		err := m.Register(func(a int) {})
		require.NoError(t, err)
		err = m.Register(func(a int) {})
		assert.Equal(t, ErrHandlerAlreadyRegistered, err)
	})
}

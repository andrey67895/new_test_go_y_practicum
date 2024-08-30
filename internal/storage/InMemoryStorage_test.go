package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemStorage_GetCounter(t *testing.T) {
	type args struct {
		id    string
		value int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{
				id:    "TEST",
				value: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := InMemStorage{}
			ctx := context.Background()
			mem.RetrySaveCounter(ctx, tt.args.id, tt.args.value)
			got, err := mem.GetCounter(ctx, tt.args.id)
			assert.True(t, err == nil)
			assert.Equal(t, tt.args.value, got)
		})
	}
}

func TestInMemStorage_GetData(t *testing.T) {
	type args struct {
		id    string
		value int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{
				id:    "TEST",
				value: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mem := InMemStorage{}
			mem.RetrySaveCounter(ctx, tt.args.id, tt.args.value)
			got, err := mem.GetData(ctx)
			assert.False(t, err != nil)
			assert.True(t, len(got) > 0)
		})
	}
}

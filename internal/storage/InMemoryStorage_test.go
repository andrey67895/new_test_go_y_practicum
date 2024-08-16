package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
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

//func TestInMemStorage_GetGauge(t *testing.T) {
//	type args struct {
//		in0 context.Context
//		id  string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    float64
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mem := InMemStorage{}
//			got, err := mem.GetGauge(tt.args.in0, tt.args.id)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetGauge() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("GetGauge() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestInMemStorage_Ping(t *testing.T) {
//	tests := []struct {
//		name    string
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mem := InMemStorage{}
//			if err := mem.Ping(); (err != nil) != tt.wantErr {
//				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestInMemStorage_RetrySaveCounter(t *testing.T) {
//	type args struct {
//		in0   context.Context
//		id    string
//		value int64
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mem := InMemStorage{}
//			if err := mem.RetrySaveCounter(tt.args.in0, tt.args.id, tt.args.value); (err != nil) != tt.wantErr {
//				t.Errorf("RetrySaveCounter() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestInMemStorage_RetrySaveGauge(t *testing.T) {
//	type args struct {
//		in0   context.Context
//		id    string
//		delta float64
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mem := InMemStorage{}
//			if err := mem.RetrySaveGauge(tt.args.in0, tt.args.id, tt.args.delta); (err != nil) != tt.wantErr {
//				t.Errorf("RetrySaveGauge() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

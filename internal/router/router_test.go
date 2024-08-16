package router

import (
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRoutersForServer(t *testing.T) {
	type args struct {
		iStorage storage.IStorageData
	}
	tests := []struct {
		name string
		args args
		want *chi.Mux
	}{
		{
			name: "positive test #1",
			args: args{
				iStorage: storage.InMemStorage{},
			},
		},
		{
			name: "positive test #2",
			args: args{
				iStorage: storage.DBStorage{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				GetRoutersForServer(tt.args.iStorage)
			})
		})
	}
}

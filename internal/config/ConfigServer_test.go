package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitServerConfigForKey(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive test #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				os.Setenv("KEY", "KEY")
				os.Setenv("DATABASE_DSN", "/")
				os.Setenv("STORE_INTERVAL", "10")
				os.Setenv("FILE_STORAGE_PATH", "/")
				os.Setenv("RESTORE", "true")
				InitServerConfig()
			})
		})
	}
}

func Test_getBool(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive test #1",
			args: args{
				env: "true",
			},
			want: true,
		},
		{
			name: "positive test #2",
			args: args{
				env: "false",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getBool(tt.args.env), "getBool(%v)", tt.args.env)
		})
	}
}

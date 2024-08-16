package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getValueInEnv(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "positive test #1",
			args: args{
				env: "10",
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getValueInEnv(tt.args.env))
		})
	}
}

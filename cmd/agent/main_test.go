package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getMemByGopsutil(t *testing.T) {
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
				getMemByGopsutil()
			})
		})
	}
}

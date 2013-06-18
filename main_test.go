package main

import (
	"bytes"
	"github.com/bmizerany/assert"
	"testing"
)

func TestServeZero(t *testing.T) {
	i, o := bytes.NewBuffer([]byte{0, 0}), new(bytes.Buffer)
	assert.Equal(t, ErrInvalidRequest, serve(i, o))
}

func TestServeMoreThanZero(t *testing.T) {
	i, o := bytes.NewBuffer([]byte{1, 0}), new(bytes.Buffer)
	serve(i, o)
	assert.Equal(t, 8, o.Len())

	i, o = bytes.NewBuffer([]byte{2, 0}), new(bytes.Buffer)
	serve(i, o)
	assert.Equal(t, 16, o.Len())

	i, o = bytes.NewBuffer([]byte{255, 0}), new(bytes.Buffer)
	serve(i, o)
	assert.Equal(t, 255*8, o.Len())
}

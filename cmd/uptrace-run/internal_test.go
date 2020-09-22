package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTailWriter(t *testing.T) {
	w := NewTailWriter(make([]byte, 16))

	_, _ = w.Write([]byte("hello"))
	assert.Equal(t, "hello", w.Text())

	_, _ = w.Write([]byte(" world"))
	assert.Equal(t, "hello world", w.Text())

	_, _ = w.Write([]byte(" foo bar"))
	assert.Equal(t, "lo world foo bar", w.Text())

	_, _ = w.Write([]byte("----------------------------"))
	assert.Equal(t, "----------------", w.Text())
}

package generics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOverwrite(t *testing.T) {
	// bool
	assert.Equal(t, false, Overwrite(false, false))
	assert.Equal(t, true, Overwrite(false, true))
	assert.Equal(t, true, Overwrite(true, false))
	assert.Equal(t, true, Overwrite(true, true))

	// int
	assert.Equal(t, 0, Overwrite(0, 0))
	assert.Equal(t, 1, Overwrite(0, 1))
	assert.Equal(t, 1, Overwrite(1, 0))
	assert.Equal(t, 1, Overwrite(1, 1))

	// int64
	assert.Equal(t, int64(0), Overwrite(int64(0), int64(0)))
	assert.Equal(t, int64(1), Overwrite(int64(0), int64(1)))
	assert.Equal(t, int64(1), Overwrite(int64(1), int64(0)))
	assert.Equal(t, int64(1), Overwrite(int64(1), int64(1)))

	// float
	assert.Equal(t, 0.0, Overwrite(0.0, 0.0))
	assert.Equal(t, 1.1, Overwrite(0.0, 1.1))
	assert.Equal(t, 1.1, Overwrite(1.1, 0.0))
	assert.Equal(t, 1.1, Overwrite(1.1, 1.1))

	// string
	assert.Equal(t, "", Overwrite("", ""))
	assert.Equal(t, "hello", Overwrite("", "hello"))
	assert.Equal(t, "hello", Overwrite("hello", ""))
	assert.Equal(t, "hello", Overwrite("hello", "hello"))
}

func TestIsTruthy(t *testing.T) {
	ptr := ""
	assert.Equal(t, false, isTruthy(&ptr))
}

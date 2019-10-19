package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRotation0(t *testing.T) {
	vec := []int8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	result := []int8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	RotateSquareVector(vec, 0)

	assert.Equal(t, result, vec)
}

func TestRotation1(t *testing.T) {
	vec := []int8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	// 13 09 05 01
	// 14 10 06 02
	// 15 11 07 03
	// 16 12 08 04
	result := []int8{13, 9, 5, 1, 14, 10, 6, 2, 15, 11, 7, 3, 16, 12, 8, 4}
	RotateSquareVector(vec, 1)

	assert.Equal(t, result, vec)
}

func TestRotation2(t *testing.T) {
	vec := []int8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	// 16 15 14 13
	// 12 11 10 09
	// 08 07 06 05
	// 04 03 02 01
	result := []int8{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	RotateSquareVector(vec, 2)

	assert.Equal(t, result, vec)
}

func TestRotation3(t *testing.T) {
	vec := []int8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	// 04 08 12 16
	// 03 07 11 15
	// 02 06 10 14
	// 01 05 09 13
	result := []int8{4, 8, 12, 16, 3, 7, 11, 15, 2, 6, 10, 14, 1, 5, 9, 13}
	RotateSquareVector(vec, 3)

	assert.Equal(t, result, vec)
}

func TestSymmetry1(t *testing.T) {
	vec := []int8{1, 2, 3, 4}

	// 3 4
	// 1 2
	result := []int8{3, 4, 1, 2}
	PerformSymmetryVector1(result)

	assert.Equal(t, result, vec)
}

func TestSymmetry2(t *testing.T) {
	vec := []int8{1, 2, 3, 4}

	// 2 1
	// 4 3
	result := []int8{2, 1, 4, 3}
	PerformSymmetryVector2(result)

	assert.Equal(t, result, vec)
}

func TestSymmetry3(t *testing.T) {
	vec := []int8{1, 2, 3, 4, 5, 6, 7, 8, 9}

	// 1 4 7
	// 2 5 8
	// 3 6 9
	result := []int8{1, 4, 7, 2, 5, 8, 3, 6, 9}
	PerformSymmetryVector3(result)

	assert.Equal(t, result, vec)
}

func TestSymmetry4(t *testing.T) {
	vec := []int8{1, 2, 3, 4, 5, 6, 7, 8, 9}

	// 9 6 3
	// 8 5 2
	// 7 4 1
	result := []int8{9, 6, 3, 8, 5, 2, 7, 4, 1}
	PerformSymmetryVector4(result)

	assert.Equal(t, result, vec)
}

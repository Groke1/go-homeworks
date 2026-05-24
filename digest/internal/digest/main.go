package digest

import (
	"math/cmplx"
	"strings"
	"unsafe"
)

// GetCharByIndex returns the i-th character from the given string.
func GetCharByIndex(str string, idx int) rune {
	runesCount := 0
	for _, r := range str {
		if runesCount == idx {
			return r
		}
		runesCount++
	}
	panic("index out of range")
}

// GetStringBySliceOfIndexes returns a string formed by concatenating specific characters from the input string based
// on the provided indexes.
func GetStringBySliceOfIndexes(str string, indexes []int) string {
	var newString strings.Builder
	newString.Grow(len(indexes))
	runes := []rune(str)
	for _, idx := range indexes {
		newString.WriteRune(runes[idx])
	}
	return newString.String()
}

// ShiftPointer shifts the given pointer by the specified number of bytes using unsafe.Add.
func ShiftPointer(pointer **int, shift int) {
	*pointer = (*int)(unsafe.Add(unsafe.Pointer(*pointer), shift))
}

// IsComplexEqual compares two complex numbers and determines if they are equal.
func IsComplexEqual(a, b complex128) bool {
	return cmplx.Abs(a-b) < 1e-6 || a == b
}

// GetRootsOfQuadraticEquation returns two roots of a quadratic equation ax^2 + bx + c = 0.
func GetRootsOfQuadraticEquation(a, b, c float64) (complex128, complex128) {
	sqrtD := cmplx.Sqrt(complex(b*b-4*a*c, 0))

	x1 := (-complex(b, 0) + sqrtD) / complex(2*a, 0)
	x2 := (-complex(b, 0) - sqrtD) / complex(2*a, 0)
	return x1, x2
}

func Heap(arr []int, n, i int) {
	for {
		largest := i
		left := 2*i + 1
		right := left + 1

		if left < n && arr[left] > arr[largest] {
			largest = left
		}

		if right < n && arr[right] > arr[largest] {
			largest = right
		}

		if largest == i {
			break
		}

		arr[i], arr[largest] = arr[largest], arr[i]

		i = largest
	}
}

func HeapSort(arr []int) {
	n := len(arr)

	for i := n/2 - 1; i >= 0; i-- {
		Heap(arr, n, i)
	}

	for i := n - 1; i > 0; i-- {
		arr[0], arr[i] = arr[i], arr[0]
		Heap(arr, i, 0)
	}
}

// Sort sorts in-place the given slice of integers in ascending order.
func Sort(source []int) {
	HeapSort(source)
}

// ReverseSliceOne in-place reverses the order of elements in the given slice.
func ReverseSliceOne(s []int) {
	for i := 0; i < len(s)/2; i++ {
		s[i], s[len(s)-i-1] = s[len(s)-i-1], s[i]
	}
}

// ReverseSliceTwo returns a new slice of integers with elements in reverse order compared to the input slice.
// The original slice remains unmodified.
func ReverseSliceTwo(s []int) []int {
	newSlice := make([]int, len(s))
	for i := 0; i < len(s); i++ {
		newSlice[i] = s[len(s)-i-1]
	}
	return newSlice
}

// SwapPointers swaps the values of two pointers.
func SwapPointers(a, b *int) {
	*a, *b = *b, *a
}

// IsSliceEqual compares two slices of integers and returns true if they contain the same elements in the same order.
func IsSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for ind, val := range a {
		if val != b[ind] {
			return false
		}
	}
	return true
}

// DeleteByIndex deletes the element at the specified index from the slice and returns a new slice.
// The original slice remains unmodified.
func DeleteByIndex(s []int, idx int) []int {
	newSlice := make([]int, len(s)-1)
	copy(newSlice, s[:idx])
	copy(newSlice[idx:], s[idx+1:])
	return newSlice
}

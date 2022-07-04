package arraystringslice

import "testing"

func TestArrayBasis(t *testing.T) {
	t.Run("ArrayBasis", func(t *testing.T) {
		ArrayBasis()
	})
}

func TestStringBasis(t *testing.T) {
	t.Run("StringBasis", func(t *testing.T) {
		StringBasis()
	})
}

func TestSliceBasis(t *testing.T) {
	t.Run("SliceBasis", func(t *testing.T) {
		SliceBasis()
	})
}

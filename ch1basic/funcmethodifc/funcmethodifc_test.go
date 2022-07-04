package funcmethodifc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncBasis(t *testing.T) {
	t.Run("FuncBasis", func(t *testing.T) {
		FuncBasis()

		assert.Equal(t, 3, Add(1, 2))
		assert.Equal(t, 3, Add2(1, 2))
		assert.Equal(t, 1, Inc())
	})
}

func TestMethodBasis(t *testing.T) {
	t.Run("MethodBasis", func(t *testing.T) {
		MethodBasis()
	})
}

func TestIfcBasis(t *testing.T) {
	t.Run("IfcBasis", func(t *testing.T) {

	})
}

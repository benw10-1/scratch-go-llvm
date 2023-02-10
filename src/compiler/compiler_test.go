package compiler

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestCompiler(t *testing.T) {
	a := assert.New(t)

	path := `C:\Users\benja\Desktop\scratch-compiler\src\compiler\testdata\test.sb3`

	compiler, err := NewCompiler(true, path, "")
	a.NoError(err)

	err = compiler.Compile()
	a.NoError(err)
}

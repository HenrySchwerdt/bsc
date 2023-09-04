package compiler_test

import (
	"bsc/src/compiler"
	"strings"
	"testing"
)

func TestStackFrame(t *testing.T) {

	t.Run("Should get right offset to local variables", func(t *testing.T) {
		// given
		sm := compiler.NewStackManager()
		sm.StartNewStackFrame(&strings.Builder{})

		// when
		sm.Push(8)
		sm.AddVariable("x", 8)
		sm.Push(8)
		sm.AddVariable("y", 8)

		offset, _ := sm.GetVariableOffset("x")
		if offset != -8 {
			t.Fatalf("Expected -8 but got %d", offset)
		}
		offset, _ = sm.GetVariableOffset("y")
		if offset != -16 {
			t.Fatalf("Expected -16 but got %d", offset)
		}
		sm.StartNewStackFrame(&strings.Builder{})
		offset, _ = sm.GetVariableOffset("x")
		if offset != 24 {
			t.Fatalf("Expected 24 but got %d", offset)
		}
		offset, _ = sm.GetVariableOffset("y")
		if offset != 16 {
			t.Fatalf("Expected 16 but got %d", offset)
		}
		sm.CloseCurrentStackFrame(&strings.Builder{})
		offset, _ = sm.GetVariableOffset("x")
		if offset != -8 {
			t.Fatalf("Expected -8 but got %d", offset)
		}
		offset, _ = sm.GetVariableOffset("y")
		if offset != -16 {
			t.Fatalf("Expected -16 but got %d", offset)
		}

	})

}

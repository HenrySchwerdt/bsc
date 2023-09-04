package compiler

import "strings"

const WordSize = 8 // size of a word in bytes on x86_64 architecture

type StackVariable struct {
	Offset int
	Size   int
}

type StackFrame struct {
	Variables       map[string]*StackVariable
	StackFrameStart int
}

func NewStackFrame(start int) *StackFrame {
	return &StackFrame{
		Variables:       make(map[string]*StackVariable),
		StackFrameStart: start,
	}
}

type StackManager struct {
	Frames    []*StackFrame
	StackSize int
}

func NewStackManager() *StackManager {
	return &StackManager{
		Frames:    []*StackFrame{},
		StackSize: 0,
	}
}

func (sm *StackManager) CurrentFrame() *StackFrame {
	if len(sm.Frames) == 0 {
		return nil
	}
	return sm.Frames[len(sm.Frames)-1]
}

func (sm *StackManager) AddVariable(name string, size int) {
	frame := sm.CurrentFrame()
	if frame == nil {
		panic("No active stack frame to add a variable to.")
	}

	offset := sm.StackSize // Calculate offset based on current stack size

	frame.Variables[name] = &StackVariable{Offset: offset, Size: size}
}

func (sm *StackManager) AddParam(name string, paramCount int, nth int, size int) {
	frame := sm.CurrentFrame()
	if frame == nil {
		panic("No active stack frame to add a variable to.")
	}
	offset := (frame.StackFrameStart) - (paramCount-nth)*8
	frame.Variables[name] = &StackVariable{Offset: offset, Size: size}
}

func (sm *StackManager) GetVariableOffset(name string) (int, bool) {
	for i := len(sm.Frames) - 1; i >= 0; i-- {
		frame := sm.Frames[i]

		if v, exists := frame.Variables[name]; exists {
			// 24 8
			sum := sm.CurrentFrame().StackFrameStart - v.Offset
			if sum > 0 {
				return sum + WordSize, true
			}
			return sum, true
		}
	}
	return 0, false
}

func (sm *StackManager) StartNewStackFrame(out *strings.Builder) {
	// push rbp
	// 24
	out.WriteString("    push rbp\n")
	out.WriteString("    mov rbp, rsp\n")
	sm.Push(WordSize)
	sm.Frames = append(sm.Frames, NewStackFrame(sm.StackSize))

}

func (sm *StackManager) CloseCurrentStackFrame(out *strings.Builder) {
	if len(sm.Frames) == 0 {
		return
	}
	sm.Frames = sm.Frames[:len(sm.Frames)-1]
	sm.Pop(WordSize)
	out.WriteString("    mov rsp, rbp\n")
	out.WriteString("    pop rbp\n")
}

func (sm *StackManager) Push(size int) {
	sm.StackSize += size
}

func (sm *StackManager) Pop(size int) {
	sm.StackSize -= size
}

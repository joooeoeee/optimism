package exec

import (
	"errors"
	"fmt"

	"github.com/ethereum-optimism/optimism/cannon/mipsevm32"
	"github.com/ethereum-optimism/optimism/cannon/mipsevm32/program"
)

type StackTracker interface {
	PushStack(caller uint32, target uint32)
	PopStack()
}

type TraceableStackTracker interface {
	StackTracker
	Traceback()
}

type NoopStackTracker struct{}

func (n *NoopStackTracker) PushStack(caller uint32, target uint32) {}

func (n *NoopStackTracker) PopStack() {}

func (n *NoopStackTracker) Traceback() {}

type StackTrackerImpl struct {
	state mipsevm32.FPVMState

	stack  []uint32
	caller []uint32
	meta   *program.Metadata
}

func NewStackTracker(state mipsevm32.FPVMState, meta *program.Metadata) (*StackTrackerImpl, error) {
	if meta == nil {
		return nil, errors.New("metadata is nil")
	}
	return NewStackTrackerUnsafe(state, meta), nil
}

// NewStackTrackerUnsafe creates a new TraceableStackTracker without verifying meta is not nil
func NewStackTrackerUnsafe(state mipsevm32.FPVMState, meta *program.Metadata) *StackTrackerImpl {
	return &StackTrackerImpl{state: state, meta: meta}
}

func (s *StackTrackerImpl) PushStack(caller uint32, target uint32) {
	s.caller = append(s.caller, caller)
	s.stack = append(s.stack, target)
}

func (s *StackTrackerImpl) PopStack() {
	if len(s.stack) != 0 {
		fn := s.meta.LookupSymbol(s.state.GetPC())
		topFn := s.meta.LookupSymbol(s.stack[len(s.stack)-1])
		if fn != topFn {
			// most likely the function was inlined. Snap back to the last return.
			i := len(s.stack) - 1
			for ; i >= 0; i-- {
				if s.meta.LookupSymbol(s.stack[i]) == fn {
					s.stack = s.stack[:i]
					s.caller = s.caller[:i]
					break
				}
			}
		} else {
			s.stack = s.stack[:len(s.stack)-1]
			s.caller = s.caller[:len(s.caller)-1]
		}
	} else {
		fmt.Printf("ERROR: stack underflow at pc=%x. step=%d\n", s.state.GetPC(), s.state.GetStep())
	}
}

func (s *StackTrackerImpl) Traceback() {
	fmt.Printf("traceback at pc=%x. step=%d\n", s.state.GetPC(), s.state.GetStep())
	for i := len(s.stack) - 1; i >= 0; i-- {
		jumpAddr := s.stack[i]
		idx := len(s.stack) - i - 1
		fmt.Printf("\t%d %x in %s caller=%08x\n", idx, jumpAddr, s.meta.LookupSymbol(jumpAddr), s.caller[i])
	}
}

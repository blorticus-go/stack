// Package stack implements a LIFO stack safe for concurrent operation.
package stack

import "fmt"

// Stack represents a LIFO stack of arbitrary, untyped values.
type Stack struct {
	manipulator                       *stackManipulator
	channelOfOperationsForManipulator chan<- *stackManipulationMessage
}

// NewStack returns an empty stack.
func NewStack() *Stack {
	return NewStackWithInitialSizeHint(100)
}

// NewStackWithInitialSizeHint returns an empty stack using a backing store with the specified
// number of elements.
func NewStackWithInitialSizeHint(initialElementStorageSize uint) *Stack {
	m := newStackManipulator(initialElementStorageSize)
	go m.Start()

	return &Stack{
		manipulator:                       m,
		channelOfOperationsForManipulator: m.requestChannel(),
	}
}

// NewBoundedDiscardingStack returns an unbounded, discarding stack which can contain
// no more than the specified number of elements.  When the stack contains that number
// of elements, a Push() will succeed, but the element at the bottom of the stack will
// be discarded and all stack elements will move down one slot.
func NewBoundedDiscardingStack(maximumNumberOfAllowedElements uint) *Stack {
	initialSizeHint := uint(100)
	if maximumNumberOfAllowedElements < 100 {
		initialSizeHint = maximumNumberOfAllowedElements
	}

	m := newStackManipulator(initialSizeHint).whichDiscardsAtSize(maximumNumberOfAllowedElements)
	go m.Start()

	return &Stack{
		manipulator:                       m,
		channelOfOperationsForManipulator: m.requestChannel(),
	}
}

// WithAMaximumDepthOf sets the maximum number of elements allowed in the stack.  This
// will panic if an attempt is made to set a maximum depth on a discarding stack (which
// already has a maximum).  If an attempt is made to Push() to a stack that has the
// maximum number of elements, the pushed element will be discarded and Push() will
// indicate that the stack was full before the Push().  This method will panic if an
// attempt is made to set a maximum depth of zero.
func (stack *Stack) WithAMaximumDepthOf(maximumNumberOfAllowedElements uint) *Stack {
	if stack.manipulator.discardsFIFOAfterMaxSize {
		panic("You may not set a maximum stack depth with a discarding stack")
	}

	responseChannel := make(chan *stackManipulationResponse)
	stack.channelOfOperationsForManipulator <- &stackManipulationMessage{
		operation:       setMaximumDepth,
		depth:           maximumNumberOfAllowedElements,
		responseChannel: responseChannel,
	}

	response := <-responseChannel

	if response.operationError != nil {
		panic(response.operationError.Error())
	}

	return stack
}

// SetMaximumDepthTo is the same as WithAMaximumDepthOf().  There are two versions so that
// the chosen method can improve readability.  Usually, WithAMaximumDepthOf() is used as
// a chained method with the constructor, as in:
//		s := stack.NewStack().WithMaximumDepthOf(100)
// whereas SetMaximumDepthTo() is used to later change the maximum stack depth.  If the
// provided new maximum is smaller than the previous maximum, all elements between the
// top of the stack and the new (smaller) maximum will be silently discarded.
func (stack *Stack) SetMaximumDepthTo(maximumNumberOfAllowedElements uint) *Stack {
	return stack.WithAMaximumDepthOf(maximumNumberOfAllowedElements)
}

// Push pushes a value to the top of the stack.  If this is a standard stack that
// has no maximum depth, it will succeed and return false, meaning the stack was
// not full before the Push (because the stack cannot be full).  If this is a standard
// stack that has a maximum depth, Push() will discard the pushed value and return
// true if the stack was full before the Push() attempt.  Otherwise, it will push
// the value and return false.  If this is a discarding stack and the stack has the
// maximum number of elements, the new value will be added after silently discarding
// the item at the bottom of the stack.  In this case, Push() will return true.  If
// the discarding stack isn't full, the value will be added and false will be returned.
func (stack *Stack) Push(value interface{}) (cannotPushBecauseStackIsFull bool) {
	responseChannel := make(chan *stackManipulationResponse)
	stack.channelOfOperationsForManipulator <- &stackManipulationMessage{
		operation:       push,
		valueToPush:     value,
		responseChannel: responseChannel,
	}

	response := <-responseChannel

	return response.stackIsEmptyOrFullBeforeOperation
}

// Pop removes the value from the top of the stack and returns it.  If the stack was
// empty before the operation, Pop will return an undefined value and true.  If it
// was not empty before the operation, it will return the popped value and false.
func (stack *Stack) Pop() (value interface{}, stackWasEmptyBeforePop bool) {
	responseChannel := make(chan *stackManipulationResponse)
	stack.channelOfOperationsForManipulator <- &stackManipulationMessage{
		operation:       pop,
		responseChannel: responseChannel,
	}

	response := <-responseChannel

	return response.poppedValueOrCurrentDepth, response.stackIsEmptyOrFullBeforeOperation
}

// PopUint is a convenience function that will typecast the returned value as a uint.
// Naturally, if the element isn't really a uint, a runtime error will be raised.
func (stack *Stack) PopUint() (uint, bool) {
	v, b := stack.Pop()
	return v.(uint), b
}

// PopInt is a convenience function that will typecast the returned value as an int.
// Naturally, if the element isn't really an iint, a runtime error will be raised.
func (stack *Stack) PopInt() (int, bool) {
	v, b := stack.Pop()
	return v.(int), b
}

// PopByte is a convenience function that will typecast the returned value as a byte.
// Naturally, if the element isn't really a byte, a runtime error will be raised.
func (stack *Stack) PopByte() (byte, bool) {
	v, b := stack.Pop()
	return v.(byte), b
}

// PopString is a convenience function that will typecast the returned value as a string.
// Naturally, if the element isn't really a string, a runtime error will be raised.
func (stack *Stack) PopString() (string, bool) {
	v, b := stack.Pop()
	return v.(string), b
}

// Depth returns the number of values currently on the stack.
func (stack *Stack) Depth() uint {
	responseChannel := make(chan *stackManipulationResponse)
	stack.channelOfOperationsForManipulator <- &stackManipulationMessage{
		operation:       getDepth,
		responseChannel: responseChannel,
	}

	response := <-responseChannel

	return response.poppedValueOrCurrentDepth.(uint)
}

// IsEmpty returns true if the stack is empty (i.e., the depth is 0), or false otherwise.
func (stack *Stack) IsEmpty() bool {
	responseChannel := make(chan *stackManipulationResponse)
	stack.channelOfOperationsForManipulator <- &stackManipulationMessage{
		operation:       getDepth,
		responseChannel: responseChannel,
	}

	response := <-responseChannel

	return response.poppedValueOrCurrentDepth.(uint) == 0
}

// ResetToEmpty silently discards all elements on the stack and sets the stack depth to 0.
func (stack *Stack) ResetToEmpty() {
	responseChannel := make(chan *stackManipulationResponse)
	stack.channelOfOperationsForManipulator <- &stackManipulationMessage{
		operation:       resetToEmpty,
		responseChannel: responseChannel,
	}

	<-responseChannel
}

type stackOperation int

const (
	push stackOperation = iota
	pop
	resetToEmpty
	setMaximumDepth
	getDepth
)

type stackManipulationResponse struct {
	poppedValueOrCurrentDepth         interface{}
	stackIsEmptyOrFullBeforeOperation bool
	operationError                    error
}

type stackManipulationMessage struct {
	operation       stackOperation
	valueToPush     interface{}
	depth           uint
	responseChannel chan<- *stackManipulationResponse
}

type stackManipulator struct {
	channelOfRequestedOperations chan *stackManipulationMessage
	stackBackingSlice            []interface{}
	currentStackDepth            uint
	maximumStackDepth            uint
	indexInSliceOfHead           int
	discardsFIFOAfterMaxSize     bool
}

func newStackManipulator(initialSizeHint uint) *stackManipulator {
	return &stackManipulator{
		channelOfRequestedOperations: make(chan *stackManipulationMessage),
		stackBackingSlice:            make([]interface{}, initialSizeHint),
		currentStackDepth:            0,
		maximumStackDepth:            0,
		indexInSliceOfHead:           -1,
		discardsFIFOAfterMaxSize:     false,
	}
}

func (manipulator *stackManipulator) whichDiscardsAtSize(maximumDepth uint) *stackManipulator {
	manipulator.discardsFIFOAfterMaxSize = true
	manipulator.maximumStackDepth = maximumDepth
	return manipulator
}

func (manipulator *stackManipulator) requestChannel() chan<- *stackManipulationMessage {
	return manipulator.channelOfRequestedOperations
}

func (manipulator *stackManipulator) Start() {
	for {
		nextRequest := <-manipulator.channelOfRequestedOperations

		switch nextRequest.operation {
		case push:
			wasStackAlreadyFull := manipulator.push(nextRequest.valueToPush)
			nextRequest.responseChannel <- &stackManipulationResponse{nil, wasStackAlreadyFull, nil}

		case pop:
			topOfStackValue, wasStackAlreadyEmpty := manipulator.pop()
			nextRequest.responseChannel <- &stackManipulationResponse{topOfStackValue, wasStackAlreadyEmpty, nil}

		case resetToEmpty:
			manipulator.resetToEmpty()
			nextRequest.responseChannel <- &stackManipulationResponse{nil, false, nil}

		case setMaximumDepth:
			err := manipulator.setMaximumDepth(nextRequest.depth)
			nextRequest.responseChannel <- &stackManipulationResponse{nil, false, err}

		case getDepth:
			depth := manipulator.getCurrentDepth()
			nextRequest.responseChannel <- &stackManipulationResponse{depth, false, nil}
		}
	}
}

func (manipulator *stackManipulator) push(value interface{}) (stackWasAlreadyFull bool) {
	if manipulator.discardsFIFOAfterMaxSize {
		return manipulator.pushWithDiscarding(value)
	}

	return manipulator.pushWithoutDiscarding(value)

}

func (manipulator *stackManipulator) pushWithDiscarding(value interface{}) (stackWasAlreadyFull bool) {
	manipulator.indexInSliceOfHead++

	if manipulator.indexInSliceOfHead == int(manipulator.maximumStackDepth) {
		manipulator.stackBackingSlice[0] = value
		manipulator.indexInSliceOfHead = 0
	} else {
		if manipulator.indexInSliceOfHead == len(manipulator.stackBackingSlice) {
			manipulator.stackBackingSlice = append(manipulator.stackBackingSlice, value)
			manipulator.currentStackDepth++
		} else {
			manipulator.stackBackingSlice[manipulator.indexInSliceOfHead] = value
			if manipulator.currentStackDepth < manipulator.maximumStackDepth {
				manipulator.currentStackDepth++
			}
		}
	}

	return manipulator.currentStackDepth >= manipulator.maximumStackDepth
}

func (manipulator *stackManipulator) pushWithoutDiscarding(value interface{}) (stackWasAlreadyFull bool) {
	if manipulator.maximumStackDepth > 0 && manipulator.currentStackDepth == manipulator.maximumStackDepth {
		return true
	}

	manipulator.indexInSliceOfHead++

	if manipulator.currentStackDepth >= uint(len(manipulator.stackBackingSlice)) {
		manipulator.stackBackingSlice = append(manipulator.stackBackingSlice, value)
	} else {
		manipulator.stackBackingSlice[manipulator.indexInSliceOfHead] = value
	}

	manipulator.currentStackDepth++

	return false
}

func (manipulator *stackManipulator) pop() (value interface{}, stackWasAlreadyEmpty bool) {
	if manipulator.currentStackDepth == 0 {
		return nil, true
	}

	value = manipulator.stackBackingSlice[manipulator.indexInSliceOfHead]
	manipulator.indexInSliceOfHead--
	manipulator.currentStackDepth--

	if manipulator.indexInSliceOfHead < 0 {
		manipulator.indexInSliceOfHead = int(manipulator.maximumStackDepth) - 1
	}

	return value, false
}

func (manipulator *stackManipulator) resetToEmpty() {
	manipulator.indexInSliceOfHead = -1
	manipulator.currentStackDepth = 0
}

func (manipulator *stackManipulator) setMaximumDepth(newMaximumDepth uint) error {
	if newMaximumDepth < 1 {
		return fmt.Errorf("stack size must be at least 1")
	}

	if newMaximumDepth < manipulator.maximumStackDepth {
		if manipulator.indexInSliceOfHead >= int(newMaximumDepth) {
			manipulator.indexInSliceOfHead = int(newMaximumDepth) - 1
		}

		if manipulator.currentStackDepth >= newMaximumDepth {
			manipulator.currentStackDepth = newMaximumDepth
		}
	}

	manipulator.maximumStackDepth = newMaximumDepth

	return nil
}

func (manipulator *stackManipulator) getCurrentDepth() uint {
	return manipulator.currentStackDepth
}

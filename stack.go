package stack

import "fmt"

type Stack struct {
	manipulator                       *stackManipulator
	channelOfOperationsForManipulator chan<- *stackManipulationMessage
}

func NewStack() *Stack {
	m := newStackManipulator(100)
	return &Stack{
		manipulator:                       m,
		channelOfOperationsForManipulator: m.requestChannel(),
	}
}

func (stack *Stack) NewStackWithInitialSizeHint(initialElementStorageSize uint) {

}

func (stack *Stack) WithAMaximumDepthOf(maximumNumberOfAllowedElements uint) *Stack {
	return stack
}

func (stack *Stack) SetMaximumDepthTo(maximumNumberOfAllowedElements uint) *Stack {
	return stack
}

func (stack *Stack) Push(value interface{}) (cannotPushBecauseStackIsFull bool) {
	return true
}

func (stack *Stack) Pop() (value interface{}, stackWasEmptyBeforePop bool) {
	return nil, true
}

func (stack *Stack) Depth() uint {
	return 0
}

func (stack *Stack) IsEmpty() bool {
	return true
}

func (stack *Stack) ResetToEmpty() {

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
}

func newStackManipulator(initialSizeHint uint) *stackManipulator {
	return &stackManipulator{
		channelOfRequestedOperations: make(chan *stackManipulationMessage),
		stackBackingSlice:            make([]interface{}, initialSizeHint),
		currentStackDepth:            0,
		maximumStackDepth:            0,
		indexInSliceOfHead:           -1,
	}
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
	if manipulator.maximumStackDepth > 0 && manipulator.currentStackDepth == manipulator.maximumStackDepth {
		return true
	}

	manipulator.indexInSliceOfHead++
	manipulator.stackBackingSlice[manipulator.indexInSliceOfHead] = value
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

	if newMaximumDepth < manipulator.maximumStackDepth && newMaximumDepth > manipulator.currentStackDepth {
		manipulator.currentStackDepth = newMaximumDepth
		manipulator.indexInSliceOfHead = int(newMaximumDepth) - 1
	}

	manipulator.maximumStackDepth = newMaximumDepth

	return nil
}

func (manipulator *stackManipulator) getCurrentDepth() uint {
	return manipulator.currentStackDepth
}
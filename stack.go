package stack

import "fmt"

type Stack struct {
	manipulator                       *stackManipulator
	channelOfOperationsForManipulator chan<- *stackManipulationMessage
}

func NewStack() *Stack {
	return NewStackWithInitialSizeHint(100)
}

func NewStackWithInitialSizeHint(initialElementStorageSize uint) *Stack {
	m := newStackManipulator(initialElementStorageSize)
	go m.Start()

	return &Stack{
		manipulator:                       m,
		channelOfOperationsForManipulator: m.requestChannel(),
	}
}

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

func (stack *Stack) SetMaximumDepthTo(maximumNumberOfAllowedElements uint) *Stack {
	return stack.WithAMaximumDepthOf(maximumNumberOfAllowedElements)
}

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

func (stack *Stack) Pop() (value interface{}, stackWasEmptyBeforePop bool) {
	responseChannel := make(chan *stackManipulationResponse)
	stack.channelOfOperationsForManipulator <- &stackManipulationMessage{
		operation:       pop,
		responseChannel: responseChannel,
	}

	response := <-responseChannel

	return response.poppedValueOrCurrentDepth, response.stackIsEmptyOrFullBeforeOperation
}

func (stack *Stack) Depth() uint {
	responseChannel := make(chan *stackManipulationResponse)
	stack.channelOfOperationsForManipulator <- &stackManipulationMessage{
		operation:       getDepth,
		responseChannel: responseChannel,
	}

	response := <-responseChannel

	return response.poppedValueOrCurrentDepth.(uint)
}

func (stack *Stack) IsEmpty() bool {
	responseChannel := make(chan *stackManipulationResponse)
	stack.channelOfOperationsForManipulator <- &stackManipulationMessage{
		operation:       getDepth,
		responseChannel: responseChannel,
	}

	response := <-responseChannel

	return response.poppedValueOrCurrentDepth.(uint) == 0
}

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

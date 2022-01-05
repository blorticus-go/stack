package stack_test

import (
	"fmt"
	"testing"

	"github.com/blorticus-go/stack"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
)

func TestNonConcurrentNonCircularStackWithoutMax(t *testing.T) {
	g := NewGomegaWithT(t)

	s := stack.NewStackWithInitialSizeHint(4)

	for _, testCase := range []*stackOperationTestCase{
		{testname: "New Clear Stack Initial Check", operation: "check", expectedStackDepthAfterOperation: 0},
		{testname: "New Clear Stack Pop", operation: "pop", expectStackToHaveBeenEmpty: true, expectedStackDepthAfterOperation: 0},

		{testname: "Push first value", operation: "push", valueToPush: "first", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 1},
		{testname: "Push second value", operation: "push", valueToPush: "second", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 2},
		{testname: "Push third value", operation: "push", valueToPush: "third", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 3},

		{testname: "Pop with first three values", operation: "pop", expectedPopValue: "third", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 2},

		{testname: "Push fourth value", operation: "push", valueToPush: "fourth", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 3},

		{testname: "First pop after fourth push", operation: "pop", expectedPopValue: "fourth", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 2},
		{testname: "Second pop after fourth push", operation: "pop", expectedPopValue: "second", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 1},

		{testname: "Push fifth value", operation: "push", valueToPush: "fifth", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 2},
		{testname: "Push sixth value", operation: "push", valueToPush: "sixth", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 3},
		{testname: "Push seventh value", operation: "push", valueToPush: "seventh", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 4},
		{testname: "Push eighth", operation: "push", valueToPush: "eighth", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 5},
		{testname: "Push ninth value", operation: "push", valueToPush: "ninth", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 6},

		{testname: "First pop after ninth push", operation: "pop", expectedPopValue: "ninth", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 5},
		{testname: "Second pop after ninth push", operation: "pop", expectedPopValue: "eighth", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 4},
		{testname: "Third pop after ninth push", operation: "pop", expectedPopValue: "seventh", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 3},
		{testname: "Fourth pop after ninth push", operation: "pop", expectedPopValue: "sixth", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 2},
		{testname: "Fifth pop after ninth push", operation: "pop", expectedPopValue: "fifth", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 1},
		{testname: "Sixth pop after ninth push", operation: "pop", expectedPopValue: "first", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 0},
		{testname: "Seventh pop after ninth push", operation: "pop", expectStackToHaveBeenEmpty: true, expectedStackDepthAfterOperation: 0},
		{testname: "Eighth pop after ninth push", operation: "pop", expectStackToHaveBeenEmpty: true, expectedStackDepthAfterOperation: 0},
	} {
		testCase.evaluateAgainstStack(s, g)
	}
}

func TestNonConcurrentNonCircularStackWithMax(t *testing.T) {
	g := NewGomegaWithT(t)

	s := stack.NewStack().WithAMaximumDepthOf(4)

	for _, testCase := range []*stackOperationTestCase{
		{testname: "New Clear Stack Initial Check", operation: "check", expectedStackDepthAfterOperation: 0},
		{testname: "New Clear Stack Pop", operation: "pop", expectStackToHaveBeenEmpty: true, expectedStackDepthAfterOperation: 0},

		{testname: "Push first value", operation: "push", valueToPush: "first", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 1},
		{testname: "Push second value", operation: "push", valueToPush: "second", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 2},
		{testname: "Push third value", operation: "push", valueToPush: "third", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 3},

		{testname: "Pop with first three values", operation: "pop", expectedPopValue: "third", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 2},

		{testname: "Push fourth value", operation: "push", valueToPush: "fourth", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 3},

		{testname: "First pop after fourth push", operation: "pop", expectedPopValue: "fourth", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 2},
		{testname: "Second pop after fourth push", operation: "pop", expectedPopValue: "second", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 1},

		{testname: "Push fifth value", operation: "push", valueToPush: "fifth", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 2},
		{testname: "Push sixth value", operation: "push", valueToPush: "sixth", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 3},
		{testname: "Push seventh value", operation: "push", valueToPush: "seventh", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 4},
		{testname: "Push eighth", operation: "push", valueToPush: "eighth", expectStackToHaveBeenFull: true, expectedStackDepthAfterOperation: 4},
		{testname: "Push ninth value", operation: "push", valueToPush: "ninth", expectStackToHaveBeenFull: true, expectedStackDepthAfterOperation: 4},

		{testname: "First pop after ninth push", operation: "pop", expectedPopValue: "seventh", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 3},
		{testname: "Second pop after ninth push", operation: "pop", expectedPopValue: "sixth", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 2},
		{testname: "Third pop after ninth push", operation: "pop", expectedPopValue: "fifth", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 1},
		{testname: "Fourth pop after ninth push", operation: "pop", expectedPopValue: "first", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 0},
		{testname: "Seventh pop after ninth push", operation: "pop", expectStackToHaveBeenEmpty: true, expectedStackDepthAfterOperation: 0},
		{testname: "Eighth pop after ninth push", operation: "pop", expectStackToHaveBeenEmpty: true, expectedStackDepthAfterOperation: 0},
	} {
		testCase.evaluateAgainstStack(s, g)
	}
}

func TestNonConcurrentRoundedDiscardingStack(t *testing.T) {
	g := NewGomegaWithT(t)

	s := stack.NewBoundedDiscardingStack(4)

	for _, testCase := range []*stackOperationTestCase{
		{testname: "New Clear Stack Initial Check", operation: "check", expectedStackDepthAfterOperation: 0},
		{testname: "New Clear Stack Pop", operation: "pop", expectStackToHaveBeenEmpty: true, expectedStackDepthAfterOperation: 0},

		{testname: "Push first value", operation: "push", valueToPush: "first", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 1},
		{testname: "Push second value", operation: "push", valueToPush: "second", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 2},
		{testname: "Push third value", operation: "push", valueToPush: "third", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 3},

		{testname: "Pop with first three values", operation: "pop", expectedPopValue: "third", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 2},

		{testname: "Push fourth value", operation: "push", valueToPush: "fourth", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 3},

		{testname: "First pop after fourth push", operation: "pop", expectedPopValue: "fourth", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 2},
		{testname: "Second pop after fourth push", operation: "pop", expectedPopValue: "second", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 1},

		{testname: "Push fifth value", operation: "push", valueToPush: "fifth", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 2},
		{testname: "Push sixth value", operation: "push", valueToPush: "sixth", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 3},
		{testname: "Push seventh value", operation: "push", valueToPush: "seventh", expectStackToHaveBeenFull: true, expectedStackDepthAfterOperation: 4},
		{testname: "Push eighth", operation: "push", valueToPush: "eighth", expectStackToHaveBeenFull: true, expectedStackDepthAfterOperation: 4},
		{testname: "Push ninth value", operation: "push", valueToPush: "ninth", expectStackToHaveBeenFull: true, expectedStackDepthAfterOperation: 4},

		{testname: "First pop after ninth push", operation: "pop", expectedPopValue: "ninth", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 3},
		{testname: "Second pop after ninth push", operation: "pop", expectedPopValue: "eighth", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 2},
		{testname: "Third pop after ninth push", operation: "pop", expectedPopValue: "seventh", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 1},
		{testname: "Fourth pop after ninth push", operation: "pop", expectedPopValue: "sixth", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 0},
		{testname: "Seventh pop after ninth push", operation: "pop", expectStackToHaveBeenEmpty: true, expectedStackDepthAfterOperation: 0},
		{testname: "Eighth pop after ninth push", operation: "pop", expectStackToHaveBeenEmpty: true, expectedStackDepthAfterOperation: 0},
	} {
		testCase.evaluateAgainstStack(s, g)
	}

	s = stack.NewBoundedDiscardingStack(150)

	// the code arbitrarily uses 100 as the append cutoff for discarding stacks
	for i := 0; i < 149; i++ {
		if stackWasAlreadyFull := s.Push("v"); stackWasAlreadyFull {
			t.Fatalf("[Discarding Stack with 150 Depth] On push %d received stackWasAlreadyFull, expected stack to not be full", i)
		}
	}

	if stackWasAlreadyFull := s.Push("v"); !stackWasAlreadyFull {
		t.Fatalf("[Discarding Stack with 150 Depth] On push 150 received !stackWasAlreadyFull, expected stack to be full")
	}

	for i := 0; i < 150; i++ {
		if _, stackWasAlreadyEmpty := s.Pop(); stackWasAlreadyEmpty {
			t.Fatalf("[Discarding Stack with 150 Depth] On pop %d received stackWasAlreadyEmpty, expected stack to not be already empty", i)
		}
	}

	if _, stackWasAlreadyEmpty := s.Pop(); !stackWasAlreadyEmpty {
		t.Fatalf("[Discarding Stack with 150 Depth] On pop 151 received !stackWasAlreadyEmpty, expected stack to be already empty")
	}

}

func TestReset(t *testing.T) {
	g := NewGomegaWithT(t)
	s := stack.NewStack()

	for _, testCase := range []*stackOperationTestCase{
		{testname: "New Clear Stack Initial Check", operation: "check", expectedStackDepthAfterOperation: 0},
		{testname: "Reset of Empty Stack", operation: "reset", expectedStackDepthAfterOperation: 0},
		{testname: "Push First Value", operation: "push", valueToPush: "first", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 1},
		{testname: "Push second value", operation: "push", valueToPush: "second", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 2},
		{testname: "Push third value", operation: "push", valueToPush: "third", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 3},
		{testname: "Reset of Stack With Three Values", operation: "reset", expectedStackDepthAfterOperation: 0},
		{testname: "Push Third Value", operation: "push", valueToPush: "third", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 1},
		{testname: "Push Fourth value", operation: "push", valueToPush: "fourth", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 2},
		{testname: "Push Fifth value", operation: "push", valueToPush: "fifth", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 3},
	} {
		testCase.evaluateAgainstStack(s, g)
	}
}

func testForPanic(functionThatShouldPanic func()) (functionDidPanic bool) {
	c := make(chan bool)
	go func(c chan bool, f func()) {
		defer func() {
			err := recover()
			if err != nil {
				c <- true
			}
		}()

		functionThatShouldPanic()

		c <- false
	}(c, functionThatShouldPanic)

	return <-c
}

func TestMaximumResizeSmaller(t *testing.T) {
	g := NewGomegaWithT(t)

	s := stack.NewStack().WithAMaximumDepthOf(6)

	for _, testCase := range []*stackOperationTestCase{
		{testname: "Push First Value", operation: "push", valueToPush: "first", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 1},
		{testname: "Push Second value", operation: "push", valueToPush: "second", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 2},
		{testname: "Push Third value", operation: "push", valueToPush: "third", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 3},
		{testname: "Push Fourth value", operation: "push", valueToPush: "fourth", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 4},
		{testname: "Push Fifth value", operation: "push", valueToPush: "fifth", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 5},
		{testname: "Push Sixth value", operation: "push", valueToPush: "sixth", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 6},
	} {
		testCase.evaluateAgainstStack(s, g)
	}

	s.SetMaximumDepthTo(3)

	for _, testCase := range []*stackOperationTestCase{
		{testname: "First Pop After Resize", operation: "pop", expectedPopValue: "third", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 2},
		{testname: "Second Pop After Resize", operation: "pop", expectedPopValue: "second", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 1},
		{testname: "Third Pop After Resize", operation: "pop", expectedPopValue: "first", expectStackToHaveBeenEmpty: false, expectedStackDepthAfterOperation: 0},
		{testname: "Fourth Pop After Resize", operation: "pop", expectStackToHaveBeenEmpty: true, expectedStackDepthAfterOperation: 0},
	} {
		testCase.evaluateAgainstStack(s, g)
	}

	for _, testCase := range []*stackOperationTestCase{
		{testname: "Push First Value after resize", operation: "push", valueToPush: "first", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 1},
		{testname: "Push Second value after resize", operation: "push", valueToPush: "second", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 2},
		{testname: "Push Third value after resize", operation: "push", valueToPush: "third", expectStackToHaveBeenFull: false, expectedStackDepthAfterOperation: 3},
		{testname: "Push Fourth value after resize", operation: "push", valueToPush: "fourth", expectStackToHaveBeenFull: true, expectedStackDepthAfterOperation: 3},
		{testname: "Push Fifth value after resize", operation: "push", valueToPush: "fifth", expectStackToHaveBeenFull: true, expectedStackDepthAfterOperation: 3},
		{testname: "Push Sixth value after resize", operation: "push", valueToPush: "sixth", expectStackToHaveBeenFull: true, expectedStackDepthAfterOperation: 3},
	} {
		testCase.evaluateAgainstStack(s, g)
	}

}

func TestPanicConditions(t *testing.T) {
	f := func() { stack.NewStack().WithAMaximumDepthOf(0) }
	if functionDidPanic := testForPanic(f); !functionDidPanic {
		t.Errorf("On NewStack WithMaximumDepthOf 0 expected panic, did not panic")
	}

	f = func() { stack.NewBoundedDiscardingStack(5).WithAMaximumDepthOf(10) }
	if functionDidPanic := testForPanic(f); !functionDidPanic {
		t.Errorf("On attempt to set MaximumDepth on DiscardingStack expected panic, did not panic")
	}

}

type typedPopTestCase struct {
	testname                   string
	expectedValue              interface{}
	expectStackToHaveBeenEmpty bool
}

type conversionAttemptMessage struct {
	msgType    string // "success", "does not match", "stack is empty", "panic"
	panicError error
}

func (testCase *typedPopTestCase) attemptTypeConvertedPop(s *stack.Stack, conversionAttemptMessageChannel chan *conversionAttemptMessage) {
	defer func(c chan *conversionAttemptMessage) {
		if err := recover(); err != nil {
			c <- &conversionAttemptMessage{
				msgType:    "panic",
				panicError: err.(error),
			}
		}
	}(conversionAttemptMessageChannel)

	switch testCase.expectedValue.(type) {
	case int:
		v, stackIsEmpty := s.PopInt()
		if stackIsEmpty {
			conversionAttemptMessageChannel <- &conversionAttemptMessage{msgType: "stack is empty"}
			return
		}

		if v != testCase.expectedValue.(int) {
			conversionAttemptMessageChannel <- &conversionAttemptMessage{msgType: "does not match"}
			return
		}

		conversionAttemptMessageChannel <- &conversionAttemptMessage{msgType: "success"}
		return

	case uint:
		v, stackIsEmpty := s.PopUint()
		if stackIsEmpty {
			conversionAttemptMessageChannel <- &conversionAttemptMessage{msgType: "stack is empty"}
			return
		}

		if v != testCase.expectedValue.(uint) {
			conversionAttemptMessageChannel <- &conversionAttemptMessage{msgType: "does not match"}
			return
		}

		conversionAttemptMessageChannel <- &conversionAttemptMessage{msgType: "success"}
		return

	case string:
		v, stackIsEmpty := s.PopString()
		if stackIsEmpty {
			conversionAttemptMessageChannel <- &conversionAttemptMessage{msgType: "stack is empty"}
			return
		}

		if v != testCase.expectedValue.(string) {
			conversionAttemptMessageChannel <- &conversionAttemptMessage{msgType: "does not match"}
			return
		}

		conversionAttemptMessageChannel <- &conversionAttemptMessage{msgType: "success"}
		return

	case byte:
		v, stackIsEmpty := s.PopByte()
		if stackIsEmpty {
			conversionAttemptMessageChannel <- &conversionAttemptMessage{msgType: "stack is empty"}
			return
		}

		if v != testCase.expectedValue.(byte) {
			conversionAttemptMessageChannel <- &conversionAttemptMessage{msgType: "does not match"}
			return
		}

		conversionAttemptMessageChannel <- &conversionAttemptMessage{msgType: "success"}
		return
	}
}

func (testCase *typedPopTestCase) evaulateAgainstStack(s *stack.Stack) error {
	conversionAttemptMessageChannel := make(chan *conversionAttemptMessage)
	go testCase.attemptTypeConvertedPop(s, conversionAttemptMessageChannel)

	msg := <-conversionAttemptMessageChannel

	switch msg.msgType {
	case "success":
		return nil

	case "does not match":
		return fmt.Errorf("popped value does not match expected value")

	case "stack is empty":
		if !testCase.expectStackToHaveBeenEmpty {
			return fmt.Errorf("stack was empty before pop but that was not expected")
		}

		return nil

	case "panic":
		return fmt.Errorf("a panic occurred: %s", msg.panicError.Error())
	}

	return fmt.Errorf("an unexpected condition was returned")
}

func TestTypedPopping(t *testing.T) {
	s := stack.NewStack()

	for i, v := range []interface{}{
		"string",
		uint(10),
		int(-10),
		byte(100),
	} {
		if stackWasAlreadyFull := s.Push(v); stackWasAlreadyFull {
			t.Errorf("On push %d expected stack to not be full but was", i+1)
		}
	}

	for _, testCase := range []*typedPopTestCase{
		{"pop of byte 100", byte(100), false},
		{"pop of int -10", int(-10), false},
		{"pop of uint 10", uint(10), false},
		{"pop of string 'string'", "string", false},
		{testname: "pop int when stack is empty", expectedValue: int(0), expectStackToHaveBeenEmpty: true},
		{testname: "pop uint when stack is empty", expectedValue: uint(0), expectStackToHaveBeenEmpty: true},
		{testname: "pop string when stack is empty", expectedValue: "", expectStackToHaveBeenEmpty: true},
		{testname: "pop byte when stack is empty", expectedValue: byte(0), expectStackToHaveBeenEmpty: true},
	} {
		if err := testCase.evaulateAgainstStack(s); err != nil {
			t.Errorf("[%s] %s", testCase.testname, err.Error())
		}
	}
}

type stackOperationTestCase struct {
	testname                         string
	operation                        string // "push", "pop", "check", "reset"
	valueToPush                      string
	expectedPopValue                 string
	expectStackToHaveBeenEmpty       bool
	expectStackToHaveBeenFull        bool
	expectedStackDepthAfterOperation uint
}

func (testCase *stackOperationTestCase) evaluateAgainstStack(s *stack.Stack, g *gomega.WithT) {
	switch testCase.operation {
	case "push":
		stackWasAlreadyFull := s.Push(testCase.valueToPush)
		if testCase.expectStackToHaveBeenFull {
			g.Expect(stackWasAlreadyFull).To(BeTrue(), fmt.Sprintf("[%s] stack should have been full on push", testCase.testname))
		} else {
			g.Expect(stackWasAlreadyFull).To(BeFalse(), fmt.Sprintf("[%s] stack should have been full on push", testCase.testname))
		}
	case "pop":
		poppedValue, stackWasAlreadyEmpty := s.Pop()

		if testCase.expectStackToHaveBeenEmpty {
			g.Expect(stackWasAlreadyEmpty).To(BeTrue(), fmt.Sprintf("[%s] stack should be empty on pop", testCase.testname))
		} else {
			g.Expect(stackWasAlreadyEmpty).To(BeFalse(), fmt.Sprintf("[%s] stack should be empty on pop", testCase.testname))
			g.Expect(poppedValue.(string)).To(Equal(testCase.expectedPopValue), fmt.Sprintf("[%s] popped value should be", testCase.expectedPopValue))
		}

	case "reset":
		s.ResetToEmpty()
	}

	if s.Depth() != testCase.expectedStackDepthAfterOperation {
		g.Expect(s.Depth()).To(Equal(testCase.expectedStackDepthAfterOperation), fmt.Sprintf("[%s] after %s stack depth should be %d", testCase.testname, testCase.operation, testCase.expectedStackDepthAfterOperation))
	}

	if testCase.expectedStackDepthAfterOperation == 0 {
		g.Expect(s.IsEmpty()).To(BeTrue(), fmt.Sprintf("[%s] stack IsEmpty should be true", testCase.testname))
	} else {
		g.Expect(s.IsEmpty()).To(BeFalse(), fmt.Sprintf("[%s] stack IsEmpty should be false", testCase.testname))
	}
}

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
			g.Expect(stackWasAlreadyFull).To(BeTrue(), fmt.Sprintf("[%s] stack have be full on push", testCase.testname))
		} else {
			g.Expect(stackWasAlreadyFull).To(BeFalse(), fmt.Sprintf("[%s] stack have be full on push", testCase.testname))
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
		g.Expect(s.Depth()).To(Equal(testCase.expectedStackDepthAfterOperation), fmt.Sprintf("[%s] after %s stack depth should be %d", testCase.expectedPopValue, testCase.operation, testCase.expectedStackDepthAfterOperation))
	}

	if testCase.expectedStackDepthAfterOperation == 0 {
		g.Expect(s.IsEmpty()).To(BeTrue(), fmt.Sprintf("[%s] stack IsEmpty should be true", testCase.testname))
	} else {
		g.Expect(s.IsEmpty()).To(BeFalse(), fmt.Sprintf("[%s] stack IsEmpty should be false", testCase.testname))
	}
}

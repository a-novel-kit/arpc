package arpc_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a-novel-kit/arpc"
)

type FooErr struct {
	message string
}

func (e *FooErr) Error() string {
	return e.message
}

func NewFooErr(message string) error {
	return &FooErr{message}
}

func TestHandleError(t *testing.T) {
	var (
		ErrA = errors.New("A")
		ErrB = errors.New("B")
	)

	testCases := []struct {
		name string

		handler arpc.ErrorHandler
		err     error

		expectMessage string
		expectCode    codes.Code
	}{
		{
			name: "DefaultOnly",

			handler: arpc.HandleError(codes.Internal),
			err:     ErrA,

			expectMessage: "A",
			expectCode:    codes.Internal,
		},
		{
			name: "Is",

			handler: arpc.HandleError(codes.Internal).Is(ErrA, codes.InvalidArgument),
			err:     ErrA,

			expectMessage: "A",
			expectCode:    codes.InvalidArgument,
		},
		{
			name: "IsNot",

			handler: arpc.HandleError(codes.Internal).Is(ErrA, codes.InvalidArgument),
			err:     ErrB,

			expectMessage: "B",
			expectCode:    codes.Internal,
		},
		{
			name: "IsW",

			handler: arpc.HandleError(codes.Internal).IsW(ErrA, codes.InvalidArgument, ErrB),
			err:     ErrA,

			expectMessage: "B\nA",
			expectCode:    codes.InvalidArgument,
		},
		{
			name: "IsWNot",

			handler: arpc.HandleError(codes.Internal).IsW(ErrA, codes.InvalidArgument, ErrB),
			err:     ErrB,

			expectMessage: "B",
			expectCode:    codes.Internal,
		},
		{
			name: "IsWF",

			handler: arpc.HandleError(codes.Internal).IsWF(ErrA, codes.InvalidArgument, "Hello %s", "World"),
			err:     ErrA,

			expectMessage: "Hello World",
			expectCode:    codes.InvalidArgument,
		},
		{
			name: "IsWFNot",

			handler: arpc.HandleError(codes.Internal).IsWF(ErrA, codes.InvalidArgument, "Hello %s", "World"),
			err:     ErrB,

			expectMessage: "B",
			expectCode:    codes.Internal,
		},
		{
			name: "As",

			handler: arpc.HandleError(codes.Internal).As(lo.ToPtr(&FooErr{}), codes.InvalidArgument),
			err:     NewFooErr("A"),

			expectMessage: "A",
			expectCode:    codes.InvalidArgument,
		},
		{
			name: "AsNot",

			handler: arpc.HandleError(codes.Internal).As(lo.ToPtr(&FooErr{}), codes.InvalidArgument),
			err:     ErrA,

			expectMessage: "A",
			expectCode:    codes.Internal,
		},
		{
			name: "AsW",

			handler: arpc.HandleError(codes.Internal).AsW(lo.ToPtr(&FooErr{}), codes.InvalidArgument, ErrB),
			err:     NewFooErr("A"),

			expectMessage: "B\nA",
			expectCode:    codes.InvalidArgument,
		},
		{
			name: "AsWNot",

			handler: arpc.HandleError(codes.Internal).AsW(lo.ToPtr(&FooErr{}), codes.InvalidArgument, ErrB),
			err:     ErrA,

			expectMessage: "A",
			expectCode:    codes.Internal,
		},
		{
			name: "AsWF",

			handler: arpc.HandleError(codes.Internal).AsWF(lo.ToPtr(&FooErr{}), codes.InvalidArgument, "Hello %s", "World"),
			err:     NewFooErr("A"),

			expectMessage: "Hello World",
			expectCode:    codes.InvalidArgument,
		},
		{
			name: "AsWFNot",

			handler: arpc.HandleError(codes.Internal).AsWF(lo.ToPtr(&FooErr{}), codes.InvalidArgument, "Hello %s", "World"),
			err:     ErrA,

			expectMessage: "A",
			expectCode:    codes.Internal,
		},
		{
			name: "Test",

			handler: arpc.HandleError(codes.Internal).Test(func(err error) (error, bool) {
				if err.Error() == "A" {
					return status.Errorf(codes.InvalidArgument, err.Error()), true
				}

				return nil, false
			}),
			err: ErrA,

			expectMessage: "A",
			expectCode:    codes.InvalidArgument,
		},
		{
			name: "TestNot",

			handler: arpc.HandleError(codes.Internal).Test(func(err error) (error, bool) {
				if err.Error() == "A" {
					return status.Errorf(codes.InvalidArgument, err.Error()), true
				}

				return nil, false
			}),
			err: ErrB,

			expectMessage: "B",
			expectCode:    codes.Internal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.handler.Handle(testCase.err)
			require.Error(t, err)

			st, ok := status.FromError(err)
			require.True(t, ok)

			require.Equal(t, testCase.expectMessage, st.Message())
			require.Equal(t, testCase.expectCode, st.Code())
		})
	}
}

func TestHandleErrorConcurrency(t *testing.T) {
	var (
		ErrA = errors.New("A")
		ErrB = errors.New("B")
	)

	handler := arpc.HandleError(codes.Internal).
		Is(ErrA, codes.InvalidArgument).
		Is(ErrB, codes.NotFound)

	errs := []error{
		ErrA, ErrB,
		ErrA, ErrB,
		ErrA, ErrB,
		ErrA, ErrB,
		ErrA, ErrB,
		ErrA, ErrB,
		ErrA, ErrB,
		ErrA, ErrB,
		ErrA, ErrB,
		ErrA, ErrB,
	}

	collectCodes := sync.Map{}

	wg := new(sync.WaitGroup)

	callHandler := func(err error, i int) {
		handled := handler.Handle(err)
		collectCodes.Store(i, status.Code(handled))
		wg.Done()
	}

	for i, err := range errs {
		wg.Add(1)
		go callHandler(err, i)
	}

	wg.Wait()

	for i, err := range errs {
		code, ok := collectCodes.Load(i)
		require.True(t, ok)

		if errors.Is(err, ErrA) {
			require.Equal(t, codes.InvalidArgument, code)
		} else if errors.Is(err, ErrB) {
			require.Equal(t, codes.NotFound, code)
		} else {
			require.Fail(t, "Unknown error")
		}
	}
}

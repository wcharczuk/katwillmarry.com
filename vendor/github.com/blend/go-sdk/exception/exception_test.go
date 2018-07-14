package exception

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"strings"

	"github.com/blend/go-sdk/assert"
)

func TestNewOfString(t *testing.T) {
	a := assert.New(t)
	ex := New("this is a test")
	a.Equal("this is a test", fmt.Sprintf("%v", ex))
	a.NotNil(ex.Stack())
	a.Nil(ex.Inner())
}

func TestNewOfError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("This is an error")
	wrappedErr := New(err)
	a.NotNil(wrappedErr)
	typedWrapped := As(wrappedErr)
	a.NotNil(typedWrapped)
	a.Equal("This is an error", fmt.Sprintf("%v", typedWrapped))
}

func TestNewOfException(t *testing.T) {
	a := assert.New(t)
	ex := New(Class("This is an exception"))
	wrappedEx := New(ex)
	a.NotNil(wrappedEx)
	typedWrappedEx := As(wrappedEx)
	a.Equal("This is an exception", fmt.Sprintf("%v", typedWrappedEx))
	a.Equal(ex, typedWrappedEx)
}

func TestNewOfNil(t *testing.T) {
	a := assert.New(t)

	shouldBeNil := New(nil)
	a.Nil(shouldBeNil)
	a.Equal(nil, shouldBeNil)
}

func TestNewOfTypedNil(t *testing.T) {
	a := assert.New(t)

	var nilError error
	a.Nil(nilError)
	a.Equal(nil, nilError)

	shouldBeNil := New(nilError)
	a.Nil(shouldBeNil)
	a.True(shouldBeNil == nil)
}

func TestNewOfReturnedNil(t *testing.T) {
	a := assert.New(t)

	returnsNil := func() error {
		return nil
	}

	shouldBeNil := New(returnsNil())
	a.Nil(shouldBeNil)
	a.True(shouldBeNil == nil)

	returnsTypedNil := func() error {
		return New(nil)
	}

	shouldAlsoBeNil := returnsTypedNil()
	a.Nil(shouldAlsoBeNil)
	a.True(shouldAlsoBeNil == nil)
}

func TestError(t *testing.T) {
	a := assert.New(t)

	ex := New(Class("this is a test"))
	message := ex.Error()
	a.NotEmpty(message)
}

func TestCallers(t *testing.T) {
	a := assert.New(t)

	callStack := func() StackTrace { return callers(defaultStartDepth) }()

	a.NotNil(callStack)
	callstackStr := callStack.String()
	a.True(strings.Contains(callstackStr, "TestCallers"), callstackStr)
}

func TestExceptionFormatters(t *testing.T) {
	assert := assert.New(t)

	// test the "%v" formatter with just the exception class.
	class := &Ex{class: Class("this is a test")}
	assert.Equal("this is a test", fmt.Sprintf("%v", class))

	classAndMessage := &Ex{class: Class("foo"), message: "bar"}
	assert.Equal("foo\nmessage: bar", fmt.Sprintf("%v", classAndMessage))
}

func TestMarshalJSON(t *testing.T) {
	type ReadableStackTrace struct {
		Class   string   `json:"Class"`
		Message string   `json:"Message"`
		Stack   []string `json:"Stack"`
	}

	a := assert.New(t)
	ex := New("new test error")
	a.NotNil(ex)
	stackTrace := ex.Stack()
	typed, isTyped := stackTrace.(*StackPointers)
	a.True(isTyped)
	a.NotNil(typed)
	stackDepth := len(*typed)

	jsonErr, err := json.Marshal(ex)
	a.Nil(err)
	a.NotNil(jsonErr)

	ex2 := &ReadableStackTrace{}
	err = json.Unmarshal(jsonErr, ex2)
	a.Nil(err)
	a.Len(ex2.Stack, stackDepth)
}

func TestNest(t *testing.T) {
	a := assert.New(t)

	ex1 := New("this is an error")
	ex2 := New("this is another error")
	err := Nest(ex1, ex2)

	a.NotNil(err)
	a.NotNil(err.Inner())
	a.NotEmpty(err.Error())

	a.True(Is(ex1, Class("this is an error")))
	a.True(Is(ex1.Inner(), Class("this is another error")))
}

func TestNestNil(t *testing.T) {
	a := assert.New(t)

	var ex1 error
	var ex2 error
	var ex3 error

	err := Nest(ex1, ex2, ex3)
	a.Nil(err)
	a.Equal(nil, err)
}

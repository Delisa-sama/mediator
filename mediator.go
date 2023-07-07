package mediator

import (
	"errors"
	"reflect"
)

const ArgumentsVectorSize = 2

type ArgumentsVector [ArgumentsVectorSize]reflect.Type

type Mediator struct {
	handlers map[ArgumentsVector]reflect.Value
}

func New() *Mediator {
	return &Mediator{handlers: make(map[ArgumentsVector]reflect.Value)}
}

var (
	ErrPassedNotFunc            = errors.New("passed not a function")
	ErrTooManyArguments         = errors.New("too many arguments")
	ErrHandlerAlreadyRegistered = errors.New("handler signature already registered")
	ErrHandlerNotFound          = errors.New("handler not found")
)

func (m *Mediator) Register(handler any) error {
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		return ErrPassedNotFunc
	}

	argsCount := handlerType.NumIn()
	if argsCount > ArgumentsVectorSize {
		return ErrTooManyArguments
	}

	var argsVector ArgumentsVector
	for i := 0; i < argsCount; i++ {
		argsVector[i] = handlerType.In(i)
	}

	_, registered := m.handlers[argsVector]
	if registered {
		return ErrHandlerAlreadyRegistered
	}

	m.handlers[argsVector] = reflect.ValueOf(handler)
	return nil
}

func (m *Mediator) Publish(args ...any) error {
	if len(args) > ArgumentsVectorSize {
		return ErrTooManyArguments
	}

	var argsVector ArgumentsVector
	callArgs := make([]reflect.Value, 0, len(args))
	for i, arg := range args {
		argsVector[i] = reflect.TypeOf(arg)
		callArgs = append(callArgs, reflect.ValueOf(arg))
	}

	handler, exists := m.handlers[argsVector]
	if !exists {
		return ErrHandlerNotFound
	}

	result := handler.Call(callArgs)
	if len(result) == 0 {
		return nil
	}

	if result[0].IsNil() {
		return nil
	}

	return result[0].Interface().(error)
}

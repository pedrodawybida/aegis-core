// Package nsep implements the Nexo Secure Execution Protocol (NSEP) sandbox engine.
package nsep

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dop251/goja"
)

var (
	ErrExecutionTimeout = errors.New("nsep: execution timeout exceeded")
	ErrMaxCallsExceeded = errors.New("nsep: maximum tool calls limit exceeded")
	ErrNilCallHandler   = errors.New("nsep: call handler cannot be nil")
)

// CallHandler defines the signature for handling intercepting nsep.call() operations in Go.
type CallHandler func(opID string, params map[string]interface{}) (interface{}, error)

// SearchHandler defines the signature for handling nsep.search() discovery queries in Go.
type SearchHandler func(query string) interface{}

// Sandbox manages the execution environment for JS scripts in a secure Goja VM.
type Sandbox struct {
	Timeout       time.Duration
	MaxCalls      int
	Handler       CallHandler
	SearchHandler SearchHandler
}

// NewSandbox creates a new NSEP execution sandbox instance.
func NewSandbox(timeout time.Duration, maxCalls int, handler CallHandler) *Sandbox {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	if maxCalls <= 0 {
		maxCalls = 50
	}
	return &Sandbox{
		Timeout:  timeout,
		MaxCalls: maxCalls,
		Handler:  handler,
	}
}

// Execute runs the provided JavaScript script within an isolated Goja VM runtime.
func (s *Sandbox) Execute(script string) (interface{}, error) {
	if s.Handler == nil {
		return nil, ErrNilCallHandler
	}

	vm := goja.New()

	var mu sync.Mutex
	callCount := 0

	// Expose nsep.call() binding to JavaScript
	nsepObj := vm.NewObject()
	err := nsepObj.Set("call", func(call goja.FunctionCall) goja.Value {
		mu.Lock()
		callCount++
		currentCalls := callCount
		mu.Unlock()

		if currentCalls > s.MaxCalls {
			panic(vm.ToValue(ErrMaxCallsExceeded.Error()))
		}

		if len(call.Arguments) < 1 {
			panic(vm.ToValue("nsep.call requires at least operationId argument"))
		}

		opID := call.Arguments[0].String()
		params := make(map[string]interface{})

		if len(call.Arguments) >= 2 && !goja.IsUndefined(call.Arguments[1]) && !goja.IsNull(call.Arguments[1]) {
			if obj, ok := call.Arguments[1].Export().(map[string]interface{}); ok {
				params = obj
			}
		}

		// Execute Go handler (will be wired to compliance pipeline)
		result, err := s.Handler(opID, params)
		if err != nil {
			panic(vm.ToValue(fmt.Sprintf("nsep.call error: %v", err)))
		}

		return vm.ToValue(result)
	})

	if err != nil {
		return nil, fmt.Errorf("nsep: failed to bind nsep.call: %w", err)
	}

	// Expose nsep.search() binding to JavaScript
	err = nsepObj.Set("search", func(call goja.FunctionCall) goja.Value {
		query := ""
		if len(call.Arguments) >= 1 && !goja.IsUndefined(call.Arguments[0]) {
			query = call.Arguments[0].String()
		}

		if s.SearchHandler != nil {
			return vm.ToValue(s.SearchHandler(query))
		}
		return vm.ToValue([]interface{}{})
	})

	if err != nil {
		return nil, fmt.Errorf("nsep: failed to bind nsep.search: %w", err)
	}

	if err := vm.Set("nsep", nsepObj); err != nil {
		return nil, fmt.Errorf("nsep: failed to set nsep global object: %w", err)
	}

	// Execution timeout protection
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	timer := time.AfterFunc(s.Timeout, func() {
		vm.Interrupt(ErrExecutionTimeout)
	})
	defer timer.Stop()

	// Wrap execution in a channel to capture context cancellation
	type execResult struct {
		val goja.Value
		err error
	}

	ch := make(chan execResult, 1)
	go func() {
		val, err := vm.RunString(script)
		ch <- execResult{val: val, err: err}
	}()

	select {
	case <-ctx.Done():
		vm.Interrupt(ErrExecutionTimeout)
		return nil, ErrExecutionTimeout
	case res := <-ch:
		if res.err != nil {
			if errors.Is(res.err, ErrExecutionTimeout) || ctx.Err() != nil {
				return nil, ErrExecutionTimeout
			}
			return nil, fmt.Errorf("nsep execution error: %w", res.err)
		}
		if res.val == nil || goja.IsUndefined(res.val) || goja.IsNull(res.val) {
			return nil, nil
		}
		return res.val.Export(), nil
	}
}

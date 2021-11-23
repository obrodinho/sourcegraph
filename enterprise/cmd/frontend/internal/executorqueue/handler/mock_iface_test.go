// Code generated by go-mockgen 1.1.2; DO NOT EDIT.

package handler

import (
	"context"
	"sync"
	"time"

	executor "github.com/sourcegraph/sourcegraph/cmd/frontend/services/executors/store"
	basestore "github.com/sourcegraph/sourcegraph/internal/database/basestore"
	types "github.com/sourcegraph/sourcegraph/internal/types"
)

// MockExecutorStore is a mock implementation of the ExecutorStore interface
// (from the package github.com/sourcegraph/sourcegraph/internal/database)
// used for unit testing.
type MockExecutorStore struct {
	// DeleteInactiveHeartbeatsFunc is an instance of a mock function object
	// controlling the behavior of the method DeleteInactiveHeartbeats.
	DeleteInactiveHeartbeatsFunc *ExecutorStoreDeleteInactiveHeartbeatsFunc
	// DoneFunc is an instance of a mock function object controlling the
	// behavior of the method Done.
	DoneFunc *ExecutorStoreDoneFunc
	// GetByIDFunc is an instance of a mock function object controlling the
	// behavior of the method GetByID.
	GetByIDFunc *ExecutorStoreGetByIDFunc
	// HandleFunc is an instance of a mock function object controlling the
	// behavior of the method Handle.
	HandleFunc *ExecutorStoreHandleFunc
	// ListFunc is an instance of a mock function object controlling the
	// behavior of the method List.
	ListFunc *ExecutorStoreListFunc
	// TransactFunc is an instance of a mock function object controlling the
	// behavior of the method Transact.
	TransactFunc *ExecutorStoreTransactFunc
	// UpsertHeartbeatFunc is an instance of a mock function object
	// controlling the behavior of the method UpsertHeartbeat.
	UpsertHeartbeatFunc *ExecutorStoreUpsertHeartbeatFunc
	// WithFunc is an instance of a mock function object controlling the
	// behavior of the method With.
	WithFunc *ExecutorStoreWithFunc
}

// NewMockExecutorStore creates a new mock of the ExecutorStore interface.
// All methods return zero values for all results, unless overwritten.
func NewMockExecutorStore() *MockExecutorStore {
	return &MockExecutorStore{
		DeleteInactiveHeartbeatsFunc: &ExecutorStoreDeleteInactiveHeartbeatsFunc{
			defaultHook: func(context.Context, time.Duration) error {
				return nil
			},
		},
		DoneFunc: &ExecutorStoreDoneFunc{
			defaultHook: func(error) error {
				return nil
			},
		},
		GetByIDFunc: &ExecutorStoreGetByIDFunc{
			defaultHook: func(context.Context, int) (types.Executor, bool, error) {
				return types.Executor{}, false, nil
			},
		},
		HandleFunc: &ExecutorStoreHandleFunc{
			defaultHook: func() *basestore.TransactableHandle {
				return nil
			},
		},
		ListFunc: &ExecutorStoreListFunc{
			defaultHook: func(context.Context, executor.ExecutorStoreListOptions) ([]types.Executor, int, error) {
				return nil, 0, nil
			},
		},
		TransactFunc: &ExecutorStoreTransactFunc{
			defaultHook: func(context.Context) (executor.ExecutorStore, error) {
				return nil, nil
			},
		},
		UpsertHeartbeatFunc: &ExecutorStoreUpsertHeartbeatFunc{
			defaultHook: func(context.Context, types.Executor) error {
				return nil
			},
		},
		WithFunc: &ExecutorStoreWithFunc{
			defaultHook: func(basestore.ShareableStore) executor.ExecutorStore {
				return nil
			},
		},
	}
}

// NewMockExecutorStoreFrom creates a new mock of the MockExecutorStore
// interface. All methods delegate to the given implementation, unless
// overwritten.
func NewMockExecutorStoreFrom(i executor.ExecutorStore) *MockExecutorStore {
	return &MockExecutorStore{
		DeleteInactiveHeartbeatsFunc: &ExecutorStoreDeleteInactiveHeartbeatsFunc{
			defaultHook: i.DeleteInactiveHeartbeats,
		},
		GetByIDFunc: &ExecutorStoreGetByIDFunc{
			defaultHook: i.GetByID,
		},
		ListFunc: &ExecutorStoreListFunc{
			defaultHook: i.List,
		},
		UpsertHeartbeatFunc: &ExecutorStoreUpsertHeartbeatFunc{
			defaultHook: i.UpsertHeartbeat,
		},
	}
}

// ExecutorStoreDeleteInactiveHeartbeatsFunc describes the behavior when the
// DeleteInactiveHeartbeats method of the parent MockExecutorStore instance
// is invoked.
type ExecutorStoreDeleteInactiveHeartbeatsFunc struct {
	defaultHook func(context.Context, time.Duration) error
	hooks       []func(context.Context, time.Duration) error
	history     []ExecutorStoreDeleteInactiveHeartbeatsFuncCall
	mutex       sync.Mutex
}

// DeleteInactiveHeartbeats delegates to the next hook function in the queue
// and stores the parameter and result values of this invocation.
func (m *MockExecutorStore) DeleteInactiveHeartbeats(v0 context.Context, v1 time.Duration) error {
	r0 := m.DeleteInactiveHeartbeatsFunc.nextHook()(v0, v1)
	m.DeleteInactiveHeartbeatsFunc.appendCall(ExecutorStoreDeleteInactiveHeartbeatsFuncCall{v0, v1, r0})
	return r0
}

// SetDefaultHook sets function that is called when the
// DeleteInactiveHeartbeats method of the parent MockExecutorStore instance
// is invoked and the hook queue is empty.
func (f *ExecutorStoreDeleteInactiveHeartbeatsFunc) SetDefaultHook(hook func(context.Context, time.Duration) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// DeleteInactiveHeartbeats method of the parent MockExecutorStore instance
// invokes the hook at the front of the queue and discards it. After the
// queue is empty, the default hook function is invoked for any future
// action.
func (f *ExecutorStoreDeleteInactiveHeartbeatsFunc) PushHook(hook func(context.Context, time.Duration) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ExecutorStoreDeleteInactiveHeartbeatsFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context, time.Duration) error {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ExecutorStoreDeleteInactiveHeartbeatsFunc) PushReturn(r0 error) {
	f.PushHook(func(context.Context, time.Duration) error {
		return r0
	})
}

func (f *ExecutorStoreDeleteInactiveHeartbeatsFunc) nextHook() func(context.Context, time.Duration) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ExecutorStoreDeleteInactiveHeartbeatsFunc) appendCall(r0 ExecutorStoreDeleteInactiveHeartbeatsFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of
// ExecutorStoreDeleteInactiveHeartbeatsFuncCall objects describing the
// invocations of this function.
func (f *ExecutorStoreDeleteInactiveHeartbeatsFunc) History() []ExecutorStoreDeleteInactiveHeartbeatsFuncCall {
	f.mutex.Lock()
	history := make([]ExecutorStoreDeleteInactiveHeartbeatsFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ExecutorStoreDeleteInactiveHeartbeatsFuncCall is an object that describes
// an invocation of method DeleteInactiveHeartbeats on an instance of
// MockExecutorStore.
type ExecutorStoreDeleteInactiveHeartbeatsFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 time.Duration
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ExecutorStoreDeleteInactiveHeartbeatsFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ExecutorStoreDeleteInactiveHeartbeatsFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// ExecutorStoreDoneFunc describes the behavior when the Done method of the
// parent MockExecutorStore instance is invoked.
type ExecutorStoreDoneFunc struct {
	defaultHook func(error) error
	hooks       []func(error) error
	history     []ExecutorStoreDoneFuncCall
	mutex       sync.Mutex
}

// Done delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockExecutorStore) Done(v0 error) error {
	r0 := m.DoneFunc.nextHook()(v0)
	m.DoneFunc.appendCall(ExecutorStoreDoneFuncCall{v0, r0})
	return r0
}

// SetDefaultHook sets function that is called when the Done method of the
// parent MockExecutorStore instance is invoked and the hook queue is empty.
func (f *ExecutorStoreDoneFunc) SetDefaultHook(hook func(error) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Done method of the parent MockExecutorStore instance invokes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ExecutorStoreDoneFunc) PushHook(hook func(error) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ExecutorStoreDoneFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(error) error {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ExecutorStoreDoneFunc) PushReturn(r0 error) {
	f.PushHook(func(error) error {
		return r0
	})
}

func (f *ExecutorStoreDoneFunc) nextHook() func(error) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ExecutorStoreDoneFunc) appendCall(r0 ExecutorStoreDoneFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ExecutorStoreDoneFuncCall objects
// describing the invocations of this function.
func (f *ExecutorStoreDoneFunc) History() []ExecutorStoreDoneFuncCall {
	f.mutex.Lock()
	history := make([]ExecutorStoreDoneFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ExecutorStoreDoneFuncCall is an object that describes an invocation of
// method Done on an instance of MockExecutorStore.
type ExecutorStoreDoneFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 error
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ExecutorStoreDoneFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ExecutorStoreDoneFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// ExecutorStoreGetByIDFunc describes the behavior when the GetByID method
// of the parent MockExecutorStore instance is invoked.
type ExecutorStoreGetByIDFunc struct {
	defaultHook func(context.Context, int) (types.Executor, bool, error)
	hooks       []func(context.Context, int) (types.Executor, bool, error)
	history     []ExecutorStoreGetByIDFuncCall
	mutex       sync.Mutex
}

// GetByID delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockExecutorStore) GetByID(v0 context.Context, v1 int) (types.Executor, bool, error) {
	r0, r1, r2 := m.GetByIDFunc.nextHook()(v0, v1)
	m.GetByIDFunc.appendCall(ExecutorStoreGetByIDFuncCall{v0, v1, r0, r1, r2})
	return r0, r1, r2
}

// SetDefaultHook sets function that is called when the GetByID method of
// the parent MockExecutorStore instance is invoked and the hook queue is
// empty.
func (f *ExecutorStoreGetByIDFunc) SetDefaultHook(hook func(context.Context, int) (types.Executor, bool, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// GetByID method of the parent MockExecutorStore instance invokes the hook
// at the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ExecutorStoreGetByIDFunc) PushHook(hook func(context.Context, int) (types.Executor, bool, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ExecutorStoreGetByIDFunc) SetDefaultReturn(r0 types.Executor, r1 bool, r2 error) {
	f.SetDefaultHook(func(context.Context, int) (types.Executor, bool, error) {
		return r0, r1, r2
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ExecutorStoreGetByIDFunc) PushReturn(r0 types.Executor, r1 bool, r2 error) {
	f.PushHook(func(context.Context, int) (types.Executor, bool, error) {
		return r0, r1, r2
	})
}

func (f *ExecutorStoreGetByIDFunc) nextHook() func(context.Context, int) (types.Executor, bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ExecutorStoreGetByIDFunc) appendCall(r0 ExecutorStoreGetByIDFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ExecutorStoreGetByIDFuncCall objects
// describing the invocations of this function.
func (f *ExecutorStoreGetByIDFunc) History() []ExecutorStoreGetByIDFuncCall {
	f.mutex.Lock()
	history := make([]ExecutorStoreGetByIDFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ExecutorStoreGetByIDFuncCall is an object that describes an invocation of
// method GetByID on an instance of MockExecutorStore.
type ExecutorStoreGetByIDFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 int
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 types.Executor
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 bool
	// Result2 is the value of the 3rd result returned from this method
	// invocation.
	Result2 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ExecutorStoreGetByIDFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ExecutorStoreGetByIDFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1, c.Result2}
}

// ExecutorStoreHandleFunc describes the behavior when the Handle method of
// the parent MockExecutorStore instance is invoked.
type ExecutorStoreHandleFunc struct {
	defaultHook func() *basestore.TransactableHandle
	hooks       []func() *basestore.TransactableHandle
	history     []ExecutorStoreHandleFuncCall
	mutex       sync.Mutex
}

// Handle delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockExecutorStore) Handle() *basestore.TransactableHandle {
	r0 := m.HandleFunc.nextHook()()
	m.HandleFunc.appendCall(ExecutorStoreHandleFuncCall{r0})
	return r0
}

// SetDefaultHook sets function that is called when the Handle method of the
// parent MockExecutorStore instance is invoked and the hook queue is empty.
func (f *ExecutorStoreHandleFunc) SetDefaultHook(hook func() *basestore.TransactableHandle) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Handle method of the parent MockExecutorStore instance invokes the hook
// at the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ExecutorStoreHandleFunc) PushHook(hook func() *basestore.TransactableHandle) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ExecutorStoreHandleFunc) SetDefaultReturn(r0 *basestore.TransactableHandle) {
	f.SetDefaultHook(func() *basestore.TransactableHandle {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ExecutorStoreHandleFunc) PushReturn(r0 *basestore.TransactableHandle) {
	f.PushHook(func() *basestore.TransactableHandle {
		return r0
	})
}

func (f *ExecutorStoreHandleFunc) nextHook() func() *basestore.TransactableHandle {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ExecutorStoreHandleFunc) appendCall(r0 ExecutorStoreHandleFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ExecutorStoreHandleFuncCall objects
// describing the invocations of this function.
func (f *ExecutorStoreHandleFunc) History() []ExecutorStoreHandleFuncCall {
	f.mutex.Lock()
	history := make([]ExecutorStoreHandleFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ExecutorStoreHandleFuncCall is an object that describes an invocation of
// method Handle on an instance of MockExecutorStore.
type ExecutorStoreHandleFuncCall struct {
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 *basestore.TransactableHandle
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ExecutorStoreHandleFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ExecutorStoreHandleFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// ExecutorStoreListFunc describes the behavior when the List method of the
// parent MockExecutorStore instance is invoked.
type ExecutorStoreListFunc struct {
	defaultHook func(context.Context, executor.ExecutorStoreListOptions) ([]types.Executor, int, error)
	hooks       []func(context.Context, executor.ExecutorStoreListOptions) ([]types.Executor, int, error)
	history     []ExecutorStoreListFuncCall
	mutex       sync.Mutex
}

// List delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockExecutorStore) List(v0 context.Context, v1 executor.ExecutorStoreListOptions) ([]types.Executor, int, error) {
	r0, r1, r2 := m.ListFunc.nextHook()(v0, v1)
	m.ListFunc.appendCall(ExecutorStoreListFuncCall{v0, v1, r0, r1, r2})
	return r0, r1, r2
}

// SetDefaultHook sets function that is called when the List method of the
// parent MockExecutorStore instance is invoked and the hook queue is empty.
func (f *ExecutorStoreListFunc) SetDefaultHook(hook func(context.Context, executor.ExecutorStoreListOptions) ([]types.Executor, int, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// List method of the parent MockExecutorStore instance invokes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ExecutorStoreListFunc) PushHook(hook func(context.Context, executor.ExecutorStoreListOptions) ([]types.Executor, int, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ExecutorStoreListFunc) SetDefaultReturn(r0 []types.Executor, r1 int, r2 error) {
	f.SetDefaultHook(func(context.Context, executor.ExecutorStoreListOptions) ([]types.Executor, int, error) {
		return r0, r1, r2
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ExecutorStoreListFunc) PushReturn(r0 []types.Executor, r1 int, r2 error) {
	f.PushHook(func(context.Context, executor.ExecutorStoreListOptions) ([]types.Executor, int, error) {
		return r0, r1, r2
	})
}

func (f *ExecutorStoreListFunc) nextHook() func(context.Context, executor.ExecutorStoreListOptions) ([]types.Executor, int, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ExecutorStoreListFunc) appendCall(r0 ExecutorStoreListFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ExecutorStoreListFuncCall objects
// describing the invocations of this function.
func (f *ExecutorStoreListFunc) History() []ExecutorStoreListFuncCall {
	f.mutex.Lock()
	history := make([]ExecutorStoreListFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ExecutorStoreListFuncCall is an object that describes an invocation of
// method List on an instance of MockExecutorStore.
type ExecutorStoreListFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 executor.ExecutorStoreListOptions
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []types.Executor
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 int
	// Result2 is the value of the 3rd result returned from this method
	// invocation.
	Result2 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ExecutorStoreListFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ExecutorStoreListFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1, c.Result2}
}

// ExecutorStoreTransactFunc describes the behavior when the Transact method
// of the parent MockExecutorStore instance is invoked.
type ExecutorStoreTransactFunc struct {
	defaultHook func(context.Context) (executor.ExecutorStore, error)
	hooks       []func(context.Context) (executor.ExecutorStore, error)
	history     []ExecutorStoreTransactFuncCall
	mutex       sync.Mutex
}

// Transact delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockExecutorStore) Transact(v0 context.Context) (executor.ExecutorStore, error) {
	r0, r1 := m.TransactFunc.nextHook()(v0)
	m.TransactFunc.appendCall(ExecutorStoreTransactFuncCall{v0, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the Transact method of
// the parent MockExecutorStore instance is invoked and the hook queue is
// empty.
func (f *ExecutorStoreTransactFunc) SetDefaultHook(hook func(context.Context) (executor.ExecutorStore, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Transact method of the parent MockExecutorStore instance invokes the hook
// at the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ExecutorStoreTransactFunc) PushHook(hook func(context.Context) (executor.ExecutorStore, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ExecutorStoreTransactFunc) SetDefaultReturn(r0 executor.ExecutorStore, r1 error) {
	f.SetDefaultHook(func(context.Context) (executor.ExecutorStore, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ExecutorStoreTransactFunc) PushReturn(r0 executor.ExecutorStore, r1 error) {
	f.PushHook(func(context.Context) (executor.ExecutorStore, error) {
		return r0, r1
	})
}

func (f *ExecutorStoreTransactFunc) nextHook() func(context.Context) (executor.ExecutorStore, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ExecutorStoreTransactFunc) appendCall(r0 ExecutorStoreTransactFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ExecutorStoreTransactFuncCall objects
// describing the invocations of this function.
func (f *ExecutorStoreTransactFunc) History() []ExecutorStoreTransactFuncCall {
	f.mutex.Lock()
	history := make([]ExecutorStoreTransactFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ExecutorStoreTransactFuncCall is an object that describes an invocation
// of method Transact on an instance of MockExecutorStore.
type ExecutorStoreTransactFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 executor.ExecutorStore
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ExecutorStoreTransactFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ExecutorStoreTransactFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// ExecutorStoreUpsertHeartbeatFunc describes the behavior when the
// UpsertHeartbeat method of the parent MockExecutorStore instance is
// invoked.
type ExecutorStoreUpsertHeartbeatFunc struct {
	defaultHook func(context.Context, types.Executor) error
	hooks       []func(context.Context, types.Executor) error
	history     []ExecutorStoreUpsertHeartbeatFuncCall
	mutex       sync.Mutex
}

// UpsertHeartbeat delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockExecutorStore) UpsertHeartbeat(v0 context.Context, v1 types.Executor) error {
	r0 := m.UpsertHeartbeatFunc.nextHook()(v0, v1)
	m.UpsertHeartbeatFunc.appendCall(ExecutorStoreUpsertHeartbeatFuncCall{v0, v1, r0})
	return r0
}

// SetDefaultHook sets function that is called when the UpsertHeartbeat
// method of the parent MockExecutorStore instance is invoked and the hook
// queue is empty.
func (f *ExecutorStoreUpsertHeartbeatFunc) SetDefaultHook(hook func(context.Context, types.Executor) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// UpsertHeartbeat method of the parent MockExecutorStore instance invokes
// the hook at the front of the queue and discards it. After the queue is
// empty, the default hook function is invoked for any future action.
func (f *ExecutorStoreUpsertHeartbeatFunc) PushHook(hook func(context.Context, types.Executor) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ExecutorStoreUpsertHeartbeatFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context, types.Executor) error {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ExecutorStoreUpsertHeartbeatFunc) PushReturn(r0 error) {
	f.PushHook(func(context.Context, types.Executor) error {
		return r0
	})
}

func (f *ExecutorStoreUpsertHeartbeatFunc) nextHook() func(context.Context, types.Executor) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ExecutorStoreUpsertHeartbeatFunc) appendCall(r0 ExecutorStoreUpsertHeartbeatFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ExecutorStoreUpsertHeartbeatFuncCall
// objects describing the invocations of this function.
func (f *ExecutorStoreUpsertHeartbeatFunc) History() []ExecutorStoreUpsertHeartbeatFuncCall {
	f.mutex.Lock()
	history := make([]ExecutorStoreUpsertHeartbeatFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ExecutorStoreUpsertHeartbeatFuncCall is an object that describes an
// invocation of method UpsertHeartbeat on an instance of MockExecutorStore.
type ExecutorStoreUpsertHeartbeatFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 types.Executor
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ExecutorStoreUpsertHeartbeatFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ExecutorStoreUpsertHeartbeatFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// ExecutorStoreWithFunc describes the behavior when the With method of the
// parent MockExecutorStore instance is invoked.
type ExecutorStoreWithFunc struct {
	defaultHook func(basestore.ShareableStore) executor.ExecutorStore
	hooks       []func(basestore.ShareableStore) executor.ExecutorStore
	history     []ExecutorStoreWithFuncCall
	mutex       sync.Mutex
}

// With delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockExecutorStore) With(v0 basestore.ShareableStore) executor.ExecutorStore {
	r0 := m.WithFunc.nextHook()(v0)
	m.WithFunc.appendCall(ExecutorStoreWithFuncCall{v0, r0})
	return r0
}

// SetDefaultHook sets function that is called when the With method of the
// parent MockExecutorStore instance is invoked and the hook queue is empty.
func (f *ExecutorStoreWithFunc) SetDefaultHook(hook func(basestore.ShareableStore) executor.ExecutorStore) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// With method of the parent MockExecutorStore instance invokes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ExecutorStoreWithFunc) PushHook(hook func(basestore.ShareableStore) executor.ExecutorStore) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ExecutorStoreWithFunc) SetDefaultReturn(r0 executor.ExecutorStore) {
	f.SetDefaultHook(func(basestore.ShareableStore) executor.ExecutorStore {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ExecutorStoreWithFunc) PushReturn(r0 executor.ExecutorStore) {
	f.PushHook(func(basestore.ShareableStore) executor.ExecutorStore {
		return r0
	})
}

func (f *ExecutorStoreWithFunc) nextHook() func(basestore.ShareableStore) executor.ExecutorStore {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ExecutorStoreWithFunc) appendCall(r0 ExecutorStoreWithFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ExecutorStoreWithFuncCall objects
// describing the invocations of this function.
func (f *ExecutorStoreWithFunc) History() []ExecutorStoreWithFuncCall {
	f.mutex.Lock()
	history := make([]ExecutorStoreWithFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ExecutorStoreWithFuncCall is an object that describes an invocation of
// method With on an instance of MockExecutorStore.
type ExecutorStoreWithFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 basestore.ShareableStore
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 executor.ExecutorStore
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ExecutorStoreWithFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ExecutorStoreWithFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

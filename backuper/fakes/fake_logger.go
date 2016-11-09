// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/pivotal-cf/pcf-backup-and-restore/backuper"
)

type FakeLogger struct {
	DebugStub        func(tag, msg string, args ...interface{})
	debugMutex       sync.RWMutex
	debugArgsForCall []struct {
		tag  string
		msg  string
		args []interface{}
	}
	InfoStub        func(tag, msg string, args ...interface{})
	infoMutex       sync.RWMutex
	infoArgsForCall []struct {
		tag  string
		msg  string
		args []interface{}
	}
	WarnStub        func(tag, msg string, args ...interface{})
	warnMutex       sync.RWMutex
	warnArgsForCall []struct {
		tag  string
		msg  string
		args []interface{}
	}
	ErrorStub        func(tag, msg string, args ...interface{})
	errorMutex       sync.RWMutex
	errorArgsForCall []struct {
		tag  string
		msg  string
		args []interface{}
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeLogger) Debug(tag string, msg string, args ...interface{}) {
	fake.debugMutex.Lock()
	fake.debugArgsForCall = append(fake.debugArgsForCall, struct {
		tag  string
		msg  string
		args []interface{}
	}{tag, msg, args})
	fake.recordInvocation("Debug", []interface{}{tag, msg, args})
	fake.debugMutex.Unlock()
	if fake.DebugStub != nil {
		fake.DebugStub(tag, msg, args...)
	}
}

func (fake *FakeLogger) DebugCallCount() int {
	fake.debugMutex.RLock()
	defer fake.debugMutex.RUnlock()
	return len(fake.debugArgsForCall)
}

func (fake *FakeLogger) DebugArgsForCall(i int) (string, string, []interface{}) {
	fake.debugMutex.RLock()
	defer fake.debugMutex.RUnlock()
	return fake.debugArgsForCall[i].tag, fake.debugArgsForCall[i].msg, fake.debugArgsForCall[i].args
}

func (fake *FakeLogger) Info(tag string, msg string, args ...interface{}) {
	fake.infoMutex.Lock()
	fake.infoArgsForCall = append(fake.infoArgsForCall, struct {
		tag  string
		msg  string
		args []interface{}
	}{tag, msg, args})
	fake.recordInvocation("Info", []interface{}{tag, msg, args})
	fake.infoMutex.Unlock()
	if fake.InfoStub != nil {
		fake.InfoStub(tag, msg, args...)
	}
}

func (fake *FakeLogger) InfoCallCount() int {
	fake.infoMutex.RLock()
	defer fake.infoMutex.RUnlock()
	return len(fake.infoArgsForCall)
}

func (fake *FakeLogger) InfoArgsForCall(i int) (string, string, []interface{}) {
	fake.infoMutex.RLock()
	defer fake.infoMutex.RUnlock()
	return fake.infoArgsForCall[i].tag, fake.infoArgsForCall[i].msg, fake.infoArgsForCall[i].args
}

func (fake *FakeLogger) Warn(tag string, msg string, args ...interface{}) {
	fake.warnMutex.Lock()
	fake.warnArgsForCall = append(fake.warnArgsForCall, struct {
		tag  string
		msg  string
		args []interface{}
	}{tag, msg, args})
	fake.recordInvocation("Warn", []interface{}{tag, msg, args})
	fake.warnMutex.Unlock()
	if fake.WarnStub != nil {
		fake.WarnStub(tag, msg, args...)
	}
}

func (fake *FakeLogger) WarnCallCount() int {
	fake.warnMutex.RLock()
	defer fake.warnMutex.RUnlock()
	return len(fake.warnArgsForCall)
}

func (fake *FakeLogger) WarnArgsForCall(i int) (string, string, []interface{}) {
	fake.warnMutex.RLock()
	defer fake.warnMutex.RUnlock()
	return fake.warnArgsForCall[i].tag, fake.warnArgsForCall[i].msg, fake.warnArgsForCall[i].args
}

func (fake *FakeLogger) Error(tag string, msg string, args ...interface{}) {
	fake.errorMutex.Lock()
	fake.errorArgsForCall = append(fake.errorArgsForCall, struct {
		tag  string
		msg  string
		args []interface{}
	}{tag, msg, args})
	fake.recordInvocation("Error", []interface{}{tag, msg, args})
	fake.errorMutex.Unlock()
	if fake.ErrorStub != nil {
		fake.ErrorStub(tag, msg, args...)
	}
}

func (fake *FakeLogger) ErrorCallCount() int {
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	return len(fake.errorArgsForCall)
}

func (fake *FakeLogger) ErrorArgsForCall(i int) (string, string, []interface{}) {
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	return fake.errorArgsForCall[i].tag, fake.errorArgsForCall[i].msg, fake.errorArgsForCall[i].args
}

func (fake *FakeLogger) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.debugMutex.RLock()
	defer fake.debugMutex.RUnlock()
	fake.infoMutex.RLock()
	defer fake.infoMutex.RUnlock()
	fake.warnMutex.RLock()
	defer fake.warnMutex.RUnlock()
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeLogger) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ backuper.Logger = new(FakeLogger)

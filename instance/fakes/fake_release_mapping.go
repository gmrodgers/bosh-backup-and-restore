// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/instance"
)

type FakeReleaseMapping struct {
	FindReleaseNameStub        func(instanceGroupName, jobName string) (string, error)
	findReleaseNameMutex       sync.RWMutex
	findReleaseNameArgsForCall []struct {
		instanceGroupName string
		jobName           string
	}
	findReleaseNameReturns struct {
		result1 string
		result2 error
	}
	findReleaseNameReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	IsJobBackupOneRestoreAllStub        func(instanceGroupName, jobName string) (bool, error)
	isJobBackupOneRestoreAllMutex       sync.RWMutex
	isJobBackupOneRestoreAllArgsForCall []struct {
		instanceGroupName string
		jobName           string
	}
	isJobBackupOneRestoreAllReturns struct {
		result1 bool
		result2 error
	}
	isJobBackupOneRestoreAllReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeReleaseMapping) FindReleaseName(instanceGroupName string, jobName string) (string, error) {
	fake.findReleaseNameMutex.Lock()
	ret, specificReturn := fake.findReleaseNameReturnsOnCall[len(fake.findReleaseNameArgsForCall)]
	fake.findReleaseNameArgsForCall = append(fake.findReleaseNameArgsForCall, struct {
		instanceGroupName string
		jobName           string
	}{instanceGroupName, jobName})
	fake.recordInvocation("FindReleaseName", []interface{}{instanceGroupName, jobName})
	fake.findReleaseNameMutex.Unlock()
	if fake.FindReleaseNameStub != nil {
		return fake.FindReleaseNameStub(instanceGroupName, jobName)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.findReleaseNameReturns.result1, fake.findReleaseNameReturns.result2
}

func (fake *FakeReleaseMapping) FindReleaseNameCallCount() int {
	fake.findReleaseNameMutex.RLock()
	defer fake.findReleaseNameMutex.RUnlock()
	return len(fake.findReleaseNameArgsForCall)
}

func (fake *FakeReleaseMapping) FindReleaseNameArgsForCall(i int) (string, string) {
	fake.findReleaseNameMutex.RLock()
	defer fake.findReleaseNameMutex.RUnlock()
	return fake.findReleaseNameArgsForCall[i].instanceGroupName, fake.findReleaseNameArgsForCall[i].jobName
}

func (fake *FakeReleaseMapping) FindReleaseNameReturns(result1 string, result2 error) {
	fake.FindReleaseNameStub = nil
	fake.findReleaseNameReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeReleaseMapping) FindReleaseNameReturnsOnCall(i int, result1 string, result2 error) {
	fake.FindReleaseNameStub = nil
	if fake.findReleaseNameReturnsOnCall == nil {
		fake.findReleaseNameReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.findReleaseNameReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeReleaseMapping) IsJobBackupOneRestoreAll(instanceGroupName string, jobName string) (bool, error) {
	fake.isJobBackupOneRestoreAllMutex.Lock()
	ret, specificReturn := fake.isJobBackupOneRestoreAllReturnsOnCall[len(fake.isJobBackupOneRestoreAllArgsForCall)]
	fake.isJobBackupOneRestoreAllArgsForCall = append(fake.isJobBackupOneRestoreAllArgsForCall, struct {
		instanceGroupName string
		jobName           string
	}{instanceGroupName, jobName})
	fake.recordInvocation("IsJobBackupOneRestoreAll", []interface{}{instanceGroupName, jobName})
	fake.isJobBackupOneRestoreAllMutex.Unlock()
	if fake.IsJobBackupOneRestoreAllStub != nil {
		return fake.IsJobBackupOneRestoreAllStub(instanceGroupName, jobName)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.isJobBackupOneRestoreAllReturns.result1, fake.isJobBackupOneRestoreAllReturns.result2
}

func (fake *FakeReleaseMapping) IsJobBackupOneRestoreAllCallCount() int {
	fake.isJobBackupOneRestoreAllMutex.RLock()
	defer fake.isJobBackupOneRestoreAllMutex.RUnlock()
	return len(fake.isJobBackupOneRestoreAllArgsForCall)
}

func (fake *FakeReleaseMapping) IsJobBackupOneRestoreAllArgsForCall(i int) (string, string) {
	fake.isJobBackupOneRestoreAllMutex.RLock()
	defer fake.isJobBackupOneRestoreAllMutex.RUnlock()
	return fake.isJobBackupOneRestoreAllArgsForCall[i].instanceGroupName, fake.isJobBackupOneRestoreAllArgsForCall[i].jobName
}

func (fake *FakeReleaseMapping) IsJobBackupOneRestoreAllReturns(result1 bool, result2 error) {
	fake.IsJobBackupOneRestoreAllStub = nil
	fake.isJobBackupOneRestoreAllReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeReleaseMapping) IsJobBackupOneRestoreAllReturnsOnCall(i int, result1 bool, result2 error) {
	fake.IsJobBackupOneRestoreAllStub = nil
	if fake.isJobBackupOneRestoreAllReturnsOnCall == nil {
		fake.isJobBackupOneRestoreAllReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.isJobBackupOneRestoreAllReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeReleaseMapping) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.findReleaseNameMutex.RLock()
	defer fake.findReleaseNameMutex.RUnlock()
	fake.isJobBackupOneRestoreAllMutex.RLock()
	defer fake.isJobBackupOneRestoreAllMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeReleaseMapping) recordInvocation(key string, args []interface{}) {
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

var _ instance.ReleaseMapping = new(FakeReleaseMapping)

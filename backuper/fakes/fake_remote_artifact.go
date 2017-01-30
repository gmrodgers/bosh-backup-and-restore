// This file was generated by counterfeiter
package fakes

import (
	"io"
	"sync"

	"github.com/pivotal-cf/pcf-backup-and-restore/backuper"
)

type FakeRemoteArtifact struct {
	NameStub        func() string
	nameMutex       sync.RWMutex
	nameArgsForCall []struct{}
	nameReturns     struct {
		result1 string
	}
	IndexStub        func() string
	indexMutex       sync.RWMutex
	indexArgsForCall []struct{}
	indexReturns     struct {
		result1 string
	}
	IDStub        func() string
	iDMutex       sync.RWMutex
	iDArgsForCall []struct{}
	iDReturns     struct {
		result1 string
	}
	BackupSizeStub        func() (string, error)
	backupSizeMutex       sync.RWMutex
	backupSizeArgsForCall []struct{}
	backupSizeReturns     struct {
		result1 string
		result2 error
	}
	BackupChecksumStub        func() (backuper.BackupChecksum, error)
	backupChecksumMutex       sync.RWMutex
	backupChecksumArgsForCall []struct{}
	backupChecksumReturns     struct {
		result1 backuper.BackupChecksum
		result2 error
	}
	StreamBackupFromRemoteStub        func(io.Writer) error
	streamBackupFromRemoteMutex       sync.RWMutex
	streamBackupFromRemoteArgsForCall []struct {
		arg1 io.Writer
	}
	streamBackupFromRemoteReturns struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeRemoteArtifact) Name() string {
	fake.nameMutex.Lock()
	fake.nameArgsForCall = append(fake.nameArgsForCall, struct{}{})
	fake.recordInvocation("Name", []interface{}{})
	fake.nameMutex.Unlock()
	if fake.NameStub != nil {
		return fake.NameStub()
	}
	return fake.nameReturns.result1
}

func (fake *FakeRemoteArtifact) NameCallCount() int {
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	return len(fake.nameArgsForCall)
}

func (fake *FakeRemoteArtifact) NameReturns(result1 string) {
	fake.NameStub = nil
	fake.nameReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeRemoteArtifact) Index() string {
	fake.indexMutex.Lock()
	fake.indexArgsForCall = append(fake.indexArgsForCall, struct{}{})
	fake.recordInvocation("Index", []interface{}{})
	fake.indexMutex.Unlock()
	if fake.IndexStub != nil {
		return fake.IndexStub()
	}
	return fake.indexReturns.result1
}

func (fake *FakeRemoteArtifact) IndexCallCount() int {
	fake.indexMutex.RLock()
	defer fake.indexMutex.RUnlock()
	return len(fake.indexArgsForCall)
}

func (fake *FakeRemoteArtifact) IndexReturns(result1 string) {
	fake.IndexStub = nil
	fake.indexReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeRemoteArtifact) ID() string {
	fake.iDMutex.Lock()
	fake.iDArgsForCall = append(fake.iDArgsForCall, struct{}{})
	fake.recordInvocation("ID", []interface{}{})
	fake.iDMutex.Unlock()
	if fake.IDStub != nil {
		return fake.IDStub()
	}
	return fake.iDReturns.result1
}

func (fake *FakeRemoteArtifact) IDCallCount() int {
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	return len(fake.iDArgsForCall)
}

func (fake *FakeRemoteArtifact) IDReturns(result1 string) {
	fake.IDStub = nil
	fake.iDReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeRemoteArtifact) BackupSize() (string, error) {
	fake.backupSizeMutex.Lock()
	fake.backupSizeArgsForCall = append(fake.backupSizeArgsForCall, struct{}{})
	fake.recordInvocation("BackupSize", []interface{}{})
	fake.backupSizeMutex.Unlock()
	if fake.BackupSizeStub != nil {
		return fake.BackupSizeStub()
	}
	return fake.backupSizeReturns.result1, fake.backupSizeReturns.result2
}

func (fake *FakeRemoteArtifact) BackupSizeCallCount() int {
	fake.backupSizeMutex.RLock()
	defer fake.backupSizeMutex.RUnlock()
	return len(fake.backupSizeArgsForCall)
}

func (fake *FakeRemoteArtifact) BackupSizeReturns(result1 string, result2 error) {
	fake.BackupSizeStub = nil
	fake.backupSizeReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteArtifact) BackupChecksum() (backuper.BackupChecksum, error) {
	fake.backupChecksumMutex.Lock()
	fake.backupChecksumArgsForCall = append(fake.backupChecksumArgsForCall, struct{}{})
	fake.recordInvocation("BackupChecksum", []interface{}{})
	fake.backupChecksumMutex.Unlock()
	if fake.BackupChecksumStub != nil {
		return fake.BackupChecksumStub()
	}
	return fake.backupChecksumReturns.result1, fake.backupChecksumReturns.result2
}

func (fake *FakeRemoteArtifact) BackupChecksumCallCount() int {
	fake.backupChecksumMutex.RLock()
	defer fake.backupChecksumMutex.RUnlock()
	return len(fake.backupChecksumArgsForCall)
}

func (fake *FakeRemoteArtifact) BackupChecksumReturns(result1 backuper.BackupChecksum, result2 error) {
	fake.BackupChecksumStub = nil
	fake.backupChecksumReturns = struct {
		result1 backuper.BackupChecksum
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteArtifact) StreamBackupFromRemote(arg1 io.Writer) error {
	fake.streamBackupFromRemoteMutex.Lock()
	fake.streamBackupFromRemoteArgsForCall = append(fake.streamBackupFromRemoteArgsForCall, struct {
		arg1 io.Writer
	}{arg1})
	fake.recordInvocation("StreamBackupFromRemote", []interface{}{arg1})
	fake.streamBackupFromRemoteMutex.Unlock()
	if fake.StreamBackupFromRemoteStub != nil {
		return fake.StreamBackupFromRemoteStub(arg1)
	}
	return fake.streamBackupFromRemoteReturns.result1
}

func (fake *FakeRemoteArtifact) StreamBackupFromRemoteCallCount() int {
	fake.streamBackupFromRemoteMutex.RLock()
	defer fake.streamBackupFromRemoteMutex.RUnlock()
	return len(fake.streamBackupFromRemoteArgsForCall)
}

func (fake *FakeRemoteArtifact) StreamBackupFromRemoteArgsForCall(i int) io.Writer {
	fake.streamBackupFromRemoteMutex.RLock()
	defer fake.streamBackupFromRemoteMutex.RUnlock()
	return fake.streamBackupFromRemoteArgsForCall[i].arg1
}

func (fake *FakeRemoteArtifact) StreamBackupFromRemoteReturns(result1 error) {
	fake.StreamBackupFromRemoteStub = nil
	fake.streamBackupFromRemoteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRemoteArtifact) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	fake.indexMutex.RLock()
	defer fake.indexMutex.RUnlock()
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	fake.backupSizeMutex.RLock()
	defer fake.backupSizeMutex.RUnlock()
	fake.backupChecksumMutex.RLock()
	defer fake.backupChecksumMutex.RUnlock()
	fake.streamBackupFromRemoteMutex.RLock()
	defer fake.streamBackupFromRemoteMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeRemoteArtifact) recordInvocation(key string, args []interface{}) {
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

var _ backuper.RemoteArtifact = new(FakeRemoteArtifact)
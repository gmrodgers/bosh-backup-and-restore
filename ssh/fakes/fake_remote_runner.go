// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"io"
	"sync"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/ssh"
)

type FakeRemoteRunner struct {
	ArchiveAndDownloadStub        func(string, io.Writer) error
	archiveAndDownloadMutex       sync.RWMutex
	archiveAndDownloadArgsForCall []struct {
		arg1 string
		arg2 io.Writer
	}
	archiveAndDownloadReturns struct {
		result1 error
	}
	archiveAndDownloadReturnsOnCall map[int]struct {
		result1 error
	}
	ChecksumDirectoryStub        func(string) (map[string]string, error)
	checksumDirectoryMutex       sync.RWMutex
	checksumDirectoryArgsForCall []struct {
		arg1 string
	}
	checksumDirectoryReturns struct {
		result1 map[string]string
		result2 error
	}
	checksumDirectoryReturnsOnCall map[int]struct {
		result1 map[string]string
		result2 error
	}
	ConnectedUsernameStub        func() string
	connectedUsernameMutex       sync.RWMutex
	connectedUsernameArgsForCall []struct {
	}
	connectedUsernameReturns struct {
		result1 string
	}
	connectedUsernameReturnsOnCall map[int]struct {
		result1 string
	}
	CreateDirectoryStub        func(string) error
	createDirectoryMutex       sync.RWMutex
	createDirectoryArgsForCall []struct {
		arg1 string
	}
	createDirectoryReturns struct {
		result1 error
	}
	createDirectoryReturnsOnCall map[int]struct {
		result1 error
	}
	DirectoryExistsStub        func(string) (bool, error)
	directoryExistsMutex       sync.RWMutex
	directoryExistsArgsForCall []struct {
		arg1 string
	}
	directoryExistsReturns struct {
		result1 bool
		result2 error
	}
	directoryExistsReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	ExtractAndUploadStub        func(io.Reader, string) error
	extractAndUploadMutex       sync.RWMutex
	extractAndUploadArgsForCall []struct {
		arg1 io.Reader
		arg2 string
	}
	extractAndUploadReturns struct {
		result1 error
	}
	extractAndUploadReturnsOnCall map[int]struct {
		result1 error
	}
	FindFilesStub        func(string) ([]string, error)
	findFilesMutex       sync.RWMutex
	findFilesArgsForCall []struct {
		arg1 string
	}
	findFilesReturns struct {
		result1 []string
		result2 error
	}
	findFilesReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	IsWindowsStub        func() (bool, error)
	isWindowsMutex       sync.RWMutex
	isWindowsArgsForCall []struct {
	}
	isWindowsReturns struct {
		result1 bool
		result2 error
	}
	isWindowsReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	RemoveDirectoryStub        func(string) error
	removeDirectoryMutex       sync.RWMutex
	removeDirectoryArgsForCall []struct {
		arg1 string
	}
	removeDirectoryReturns struct {
		result1 error
	}
	removeDirectoryReturnsOnCall map[int]struct {
		result1 error
	}
	RunScriptStub        func(string, string) (string, error)
	runScriptMutex       sync.RWMutex
	runScriptArgsForCall []struct {
		arg1 string
		arg2 string
	}
	runScriptReturns struct {
		result1 string
		result2 error
	}
	runScriptReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	RunScriptWithEnvStub        func(string, map[string]string, string) (string, error)
	runScriptWithEnvMutex       sync.RWMutex
	runScriptWithEnvArgsForCall []struct {
		arg1 string
		arg2 map[string]string
		arg3 string
	}
	runScriptWithEnvReturns struct {
		result1 string
		result2 error
	}
	runScriptWithEnvReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	SizeOfStub        func(string) (string, error)
	sizeOfMutex       sync.RWMutex
	sizeOfArgsForCall []struct {
		arg1 string
	}
	sizeOfReturns struct {
		result1 string
		result2 error
	}
	sizeOfReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeRemoteRunner) ArchiveAndDownload(arg1 string, arg2 io.Writer) error {
	fake.archiveAndDownloadMutex.Lock()
	ret, specificReturn := fake.archiveAndDownloadReturnsOnCall[len(fake.archiveAndDownloadArgsForCall)]
	fake.archiveAndDownloadArgsForCall = append(fake.archiveAndDownloadArgsForCall, struct {
		arg1 string
		arg2 io.Writer
	}{arg1, arg2})
	fake.recordInvocation("ArchiveAndDownload", []interface{}{arg1, arg2})
	fake.archiveAndDownloadMutex.Unlock()
	if fake.ArchiveAndDownloadStub != nil {
		return fake.ArchiveAndDownloadStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.archiveAndDownloadReturns
	return fakeReturns.result1
}

func (fake *FakeRemoteRunner) ArchiveAndDownloadCallCount() int {
	fake.archiveAndDownloadMutex.RLock()
	defer fake.archiveAndDownloadMutex.RUnlock()
	return len(fake.archiveAndDownloadArgsForCall)
}

func (fake *FakeRemoteRunner) ArchiveAndDownloadCalls(stub func(string, io.Writer) error) {
	fake.archiveAndDownloadMutex.Lock()
	defer fake.archiveAndDownloadMutex.Unlock()
	fake.ArchiveAndDownloadStub = stub
}

func (fake *FakeRemoteRunner) ArchiveAndDownloadArgsForCall(i int) (string, io.Writer) {
	fake.archiveAndDownloadMutex.RLock()
	defer fake.archiveAndDownloadMutex.RUnlock()
	argsForCall := fake.archiveAndDownloadArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeRemoteRunner) ArchiveAndDownloadReturns(result1 error) {
	fake.archiveAndDownloadMutex.Lock()
	defer fake.archiveAndDownloadMutex.Unlock()
	fake.ArchiveAndDownloadStub = nil
	fake.archiveAndDownloadReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRemoteRunner) ArchiveAndDownloadReturnsOnCall(i int, result1 error) {
	fake.archiveAndDownloadMutex.Lock()
	defer fake.archiveAndDownloadMutex.Unlock()
	fake.ArchiveAndDownloadStub = nil
	if fake.archiveAndDownloadReturnsOnCall == nil {
		fake.archiveAndDownloadReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.archiveAndDownloadReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRemoteRunner) ChecksumDirectory(arg1 string) (map[string]string, error) {
	fake.checksumDirectoryMutex.Lock()
	ret, specificReturn := fake.checksumDirectoryReturnsOnCall[len(fake.checksumDirectoryArgsForCall)]
	fake.checksumDirectoryArgsForCall = append(fake.checksumDirectoryArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("ChecksumDirectory", []interface{}{arg1})
	fake.checksumDirectoryMutex.Unlock()
	if fake.ChecksumDirectoryStub != nil {
		return fake.ChecksumDirectoryStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.checksumDirectoryReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRemoteRunner) ChecksumDirectoryCallCount() int {
	fake.checksumDirectoryMutex.RLock()
	defer fake.checksumDirectoryMutex.RUnlock()
	return len(fake.checksumDirectoryArgsForCall)
}

func (fake *FakeRemoteRunner) ChecksumDirectoryCalls(stub func(string) (map[string]string, error)) {
	fake.checksumDirectoryMutex.Lock()
	defer fake.checksumDirectoryMutex.Unlock()
	fake.ChecksumDirectoryStub = stub
}

func (fake *FakeRemoteRunner) ChecksumDirectoryArgsForCall(i int) string {
	fake.checksumDirectoryMutex.RLock()
	defer fake.checksumDirectoryMutex.RUnlock()
	argsForCall := fake.checksumDirectoryArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRemoteRunner) ChecksumDirectoryReturns(result1 map[string]string, result2 error) {
	fake.checksumDirectoryMutex.Lock()
	defer fake.checksumDirectoryMutex.Unlock()
	fake.ChecksumDirectoryStub = nil
	fake.checksumDirectoryReturns = struct {
		result1 map[string]string
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteRunner) ChecksumDirectoryReturnsOnCall(i int, result1 map[string]string, result2 error) {
	fake.checksumDirectoryMutex.Lock()
	defer fake.checksumDirectoryMutex.Unlock()
	fake.ChecksumDirectoryStub = nil
	if fake.checksumDirectoryReturnsOnCall == nil {
		fake.checksumDirectoryReturnsOnCall = make(map[int]struct {
			result1 map[string]string
			result2 error
		})
	}
	fake.checksumDirectoryReturnsOnCall[i] = struct {
		result1 map[string]string
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteRunner) ConnectedUsername() string {
	fake.connectedUsernameMutex.Lock()
	ret, specificReturn := fake.connectedUsernameReturnsOnCall[len(fake.connectedUsernameArgsForCall)]
	fake.connectedUsernameArgsForCall = append(fake.connectedUsernameArgsForCall, struct {
	}{})
	fake.recordInvocation("ConnectedUsername", []interface{}{})
	fake.connectedUsernameMutex.Unlock()
	if fake.ConnectedUsernameStub != nil {
		return fake.ConnectedUsernameStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.connectedUsernameReturns
	return fakeReturns.result1
}

func (fake *FakeRemoteRunner) ConnectedUsernameCallCount() int {
	fake.connectedUsernameMutex.RLock()
	defer fake.connectedUsernameMutex.RUnlock()
	return len(fake.connectedUsernameArgsForCall)
}

func (fake *FakeRemoteRunner) ConnectedUsernameCalls(stub func() string) {
	fake.connectedUsernameMutex.Lock()
	defer fake.connectedUsernameMutex.Unlock()
	fake.ConnectedUsernameStub = stub
}

func (fake *FakeRemoteRunner) ConnectedUsernameReturns(result1 string) {
	fake.connectedUsernameMutex.Lock()
	defer fake.connectedUsernameMutex.Unlock()
	fake.ConnectedUsernameStub = nil
	fake.connectedUsernameReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeRemoteRunner) ConnectedUsernameReturnsOnCall(i int, result1 string) {
	fake.connectedUsernameMutex.Lock()
	defer fake.connectedUsernameMutex.Unlock()
	fake.ConnectedUsernameStub = nil
	if fake.connectedUsernameReturnsOnCall == nil {
		fake.connectedUsernameReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.connectedUsernameReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeRemoteRunner) CreateDirectory(arg1 string) error {
	fake.createDirectoryMutex.Lock()
	ret, specificReturn := fake.createDirectoryReturnsOnCall[len(fake.createDirectoryArgsForCall)]
	fake.createDirectoryArgsForCall = append(fake.createDirectoryArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("CreateDirectory", []interface{}{arg1})
	fake.createDirectoryMutex.Unlock()
	if fake.CreateDirectoryStub != nil {
		return fake.CreateDirectoryStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.createDirectoryReturns
	return fakeReturns.result1
}

func (fake *FakeRemoteRunner) CreateDirectoryCallCount() int {
	fake.createDirectoryMutex.RLock()
	defer fake.createDirectoryMutex.RUnlock()
	return len(fake.createDirectoryArgsForCall)
}

func (fake *FakeRemoteRunner) CreateDirectoryCalls(stub func(string) error) {
	fake.createDirectoryMutex.Lock()
	defer fake.createDirectoryMutex.Unlock()
	fake.CreateDirectoryStub = stub
}

func (fake *FakeRemoteRunner) CreateDirectoryArgsForCall(i int) string {
	fake.createDirectoryMutex.RLock()
	defer fake.createDirectoryMutex.RUnlock()
	argsForCall := fake.createDirectoryArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRemoteRunner) CreateDirectoryReturns(result1 error) {
	fake.createDirectoryMutex.Lock()
	defer fake.createDirectoryMutex.Unlock()
	fake.CreateDirectoryStub = nil
	fake.createDirectoryReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRemoteRunner) CreateDirectoryReturnsOnCall(i int, result1 error) {
	fake.createDirectoryMutex.Lock()
	defer fake.createDirectoryMutex.Unlock()
	fake.CreateDirectoryStub = nil
	if fake.createDirectoryReturnsOnCall == nil {
		fake.createDirectoryReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createDirectoryReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRemoteRunner) DirectoryExists(arg1 string) (bool, error) {
	fake.directoryExistsMutex.Lock()
	ret, specificReturn := fake.directoryExistsReturnsOnCall[len(fake.directoryExistsArgsForCall)]
	fake.directoryExistsArgsForCall = append(fake.directoryExistsArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("DirectoryExists", []interface{}{arg1})
	fake.directoryExistsMutex.Unlock()
	if fake.DirectoryExistsStub != nil {
		return fake.DirectoryExistsStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.directoryExistsReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRemoteRunner) DirectoryExistsCallCount() int {
	fake.directoryExistsMutex.RLock()
	defer fake.directoryExistsMutex.RUnlock()
	return len(fake.directoryExistsArgsForCall)
}

func (fake *FakeRemoteRunner) DirectoryExistsCalls(stub func(string) (bool, error)) {
	fake.directoryExistsMutex.Lock()
	defer fake.directoryExistsMutex.Unlock()
	fake.DirectoryExistsStub = stub
}

func (fake *FakeRemoteRunner) DirectoryExistsArgsForCall(i int) string {
	fake.directoryExistsMutex.RLock()
	defer fake.directoryExistsMutex.RUnlock()
	argsForCall := fake.directoryExistsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRemoteRunner) DirectoryExistsReturns(result1 bool, result2 error) {
	fake.directoryExistsMutex.Lock()
	defer fake.directoryExistsMutex.Unlock()
	fake.DirectoryExistsStub = nil
	fake.directoryExistsReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteRunner) DirectoryExistsReturnsOnCall(i int, result1 bool, result2 error) {
	fake.directoryExistsMutex.Lock()
	defer fake.directoryExistsMutex.Unlock()
	fake.DirectoryExistsStub = nil
	if fake.directoryExistsReturnsOnCall == nil {
		fake.directoryExistsReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.directoryExistsReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteRunner) ExtractAndUpload(arg1 io.Reader, arg2 string) error {
	fake.extractAndUploadMutex.Lock()
	ret, specificReturn := fake.extractAndUploadReturnsOnCall[len(fake.extractAndUploadArgsForCall)]
	fake.extractAndUploadArgsForCall = append(fake.extractAndUploadArgsForCall, struct {
		arg1 io.Reader
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("ExtractAndUpload", []interface{}{arg1, arg2})
	fake.extractAndUploadMutex.Unlock()
	if fake.ExtractAndUploadStub != nil {
		return fake.ExtractAndUploadStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.extractAndUploadReturns
	return fakeReturns.result1
}

func (fake *FakeRemoteRunner) ExtractAndUploadCallCount() int {
	fake.extractAndUploadMutex.RLock()
	defer fake.extractAndUploadMutex.RUnlock()
	return len(fake.extractAndUploadArgsForCall)
}

func (fake *FakeRemoteRunner) ExtractAndUploadCalls(stub func(io.Reader, string) error) {
	fake.extractAndUploadMutex.Lock()
	defer fake.extractAndUploadMutex.Unlock()
	fake.ExtractAndUploadStub = stub
}

func (fake *FakeRemoteRunner) ExtractAndUploadArgsForCall(i int) (io.Reader, string) {
	fake.extractAndUploadMutex.RLock()
	defer fake.extractAndUploadMutex.RUnlock()
	argsForCall := fake.extractAndUploadArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeRemoteRunner) ExtractAndUploadReturns(result1 error) {
	fake.extractAndUploadMutex.Lock()
	defer fake.extractAndUploadMutex.Unlock()
	fake.ExtractAndUploadStub = nil
	fake.extractAndUploadReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRemoteRunner) ExtractAndUploadReturnsOnCall(i int, result1 error) {
	fake.extractAndUploadMutex.Lock()
	defer fake.extractAndUploadMutex.Unlock()
	fake.ExtractAndUploadStub = nil
	if fake.extractAndUploadReturnsOnCall == nil {
		fake.extractAndUploadReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.extractAndUploadReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRemoteRunner) FindFiles(arg1 string) ([]string, error) {
	fake.findFilesMutex.Lock()
	ret, specificReturn := fake.findFilesReturnsOnCall[len(fake.findFilesArgsForCall)]
	fake.findFilesArgsForCall = append(fake.findFilesArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("FindFiles", []interface{}{arg1})
	fake.findFilesMutex.Unlock()
	if fake.FindFilesStub != nil {
		return fake.FindFilesStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.findFilesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRemoteRunner) FindFilesCallCount() int {
	fake.findFilesMutex.RLock()
	defer fake.findFilesMutex.RUnlock()
	return len(fake.findFilesArgsForCall)
}

func (fake *FakeRemoteRunner) FindFilesCalls(stub func(string) ([]string, error)) {
	fake.findFilesMutex.Lock()
	defer fake.findFilesMutex.Unlock()
	fake.FindFilesStub = stub
}

func (fake *FakeRemoteRunner) FindFilesArgsForCall(i int) string {
	fake.findFilesMutex.RLock()
	defer fake.findFilesMutex.RUnlock()
	argsForCall := fake.findFilesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRemoteRunner) FindFilesReturns(result1 []string, result2 error) {
	fake.findFilesMutex.Lock()
	defer fake.findFilesMutex.Unlock()
	fake.FindFilesStub = nil
	fake.findFilesReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteRunner) FindFilesReturnsOnCall(i int, result1 []string, result2 error) {
	fake.findFilesMutex.Lock()
	defer fake.findFilesMutex.Unlock()
	fake.FindFilesStub = nil
	if fake.findFilesReturnsOnCall == nil {
		fake.findFilesReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.findFilesReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteRunner) IsWindows() (bool, error) {
	fake.isWindowsMutex.Lock()
	ret, specificReturn := fake.isWindowsReturnsOnCall[len(fake.isWindowsArgsForCall)]
	fake.isWindowsArgsForCall = append(fake.isWindowsArgsForCall, struct {
	}{})
	fake.recordInvocation("IsWindows", []interface{}{})
	fake.isWindowsMutex.Unlock()
	if fake.IsWindowsStub != nil {
		return fake.IsWindowsStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.isWindowsReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRemoteRunner) IsWindowsCallCount() int {
	fake.isWindowsMutex.RLock()
	defer fake.isWindowsMutex.RUnlock()
	return len(fake.isWindowsArgsForCall)
}

func (fake *FakeRemoteRunner) IsWindowsCalls(stub func() (bool, error)) {
	fake.isWindowsMutex.Lock()
	defer fake.isWindowsMutex.Unlock()
	fake.IsWindowsStub = stub
}

func (fake *FakeRemoteRunner) IsWindowsReturns(result1 bool, result2 error) {
	fake.isWindowsMutex.Lock()
	defer fake.isWindowsMutex.Unlock()
	fake.IsWindowsStub = nil
	fake.isWindowsReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteRunner) IsWindowsReturnsOnCall(i int, result1 bool, result2 error) {
	fake.isWindowsMutex.Lock()
	defer fake.isWindowsMutex.Unlock()
	fake.IsWindowsStub = nil
	if fake.isWindowsReturnsOnCall == nil {
		fake.isWindowsReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.isWindowsReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteRunner) RemoveDirectory(arg1 string) error {
	fake.removeDirectoryMutex.Lock()
	ret, specificReturn := fake.removeDirectoryReturnsOnCall[len(fake.removeDirectoryArgsForCall)]
	fake.removeDirectoryArgsForCall = append(fake.removeDirectoryArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("RemoveDirectory", []interface{}{arg1})
	fake.removeDirectoryMutex.Unlock()
	if fake.RemoveDirectoryStub != nil {
		return fake.RemoveDirectoryStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.removeDirectoryReturns
	return fakeReturns.result1
}

func (fake *FakeRemoteRunner) RemoveDirectoryCallCount() int {
	fake.removeDirectoryMutex.RLock()
	defer fake.removeDirectoryMutex.RUnlock()
	return len(fake.removeDirectoryArgsForCall)
}

func (fake *FakeRemoteRunner) RemoveDirectoryCalls(stub func(string) error) {
	fake.removeDirectoryMutex.Lock()
	defer fake.removeDirectoryMutex.Unlock()
	fake.RemoveDirectoryStub = stub
}

func (fake *FakeRemoteRunner) RemoveDirectoryArgsForCall(i int) string {
	fake.removeDirectoryMutex.RLock()
	defer fake.removeDirectoryMutex.RUnlock()
	argsForCall := fake.removeDirectoryArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRemoteRunner) RemoveDirectoryReturns(result1 error) {
	fake.removeDirectoryMutex.Lock()
	defer fake.removeDirectoryMutex.Unlock()
	fake.RemoveDirectoryStub = nil
	fake.removeDirectoryReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRemoteRunner) RemoveDirectoryReturnsOnCall(i int, result1 error) {
	fake.removeDirectoryMutex.Lock()
	defer fake.removeDirectoryMutex.Unlock()
	fake.RemoveDirectoryStub = nil
	if fake.removeDirectoryReturnsOnCall == nil {
		fake.removeDirectoryReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.removeDirectoryReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRemoteRunner) RunScript(arg1 string, arg2 string) (string, error) {
	fake.runScriptMutex.Lock()
	ret, specificReturn := fake.runScriptReturnsOnCall[len(fake.runScriptArgsForCall)]
	fake.runScriptArgsForCall = append(fake.runScriptArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("RunScript", []interface{}{arg1, arg2})
	fake.runScriptMutex.Unlock()
	if fake.RunScriptStub != nil {
		return fake.RunScriptStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.runScriptReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRemoteRunner) RunScriptCallCount() int {
	fake.runScriptMutex.RLock()
	defer fake.runScriptMutex.RUnlock()
	return len(fake.runScriptArgsForCall)
}

func (fake *FakeRemoteRunner) RunScriptCalls(stub func(string, string) (string, error)) {
	fake.runScriptMutex.Lock()
	defer fake.runScriptMutex.Unlock()
	fake.RunScriptStub = stub
}

func (fake *FakeRemoteRunner) RunScriptArgsForCall(i int) (string, string) {
	fake.runScriptMutex.RLock()
	defer fake.runScriptMutex.RUnlock()
	argsForCall := fake.runScriptArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeRemoteRunner) RunScriptReturns(result1 string, result2 error) {
	fake.runScriptMutex.Lock()
	defer fake.runScriptMutex.Unlock()
	fake.RunScriptStub = nil
	fake.runScriptReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteRunner) RunScriptReturnsOnCall(i int, result1 string, result2 error) {
	fake.runScriptMutex.Lock()
	defer fake.runScriptMutex.Unlock()
	fake.RunScriptStub = nil
	if fake.runScriptReturnsOnCall == nil {
		fake.runScriptReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.runScriptReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteRunner) RunScriptWithEnv(arg1 string, arg2 map[string]string, arg3 string) (string, error) {
	fake.runScriptWithEnvMutex.Lock()
	ret, specificReturn := fake.runScriptWithEnvReturnsOnCall[len(fake.runScriptWithEnvArgsForCall)]
	fake.runScriptWithEnvArgsForCall = append(fake.runScriptWithEnvArgsForCall, struct {
		arg1 string
		arg2 map[string]string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("RunScriptWithEnv", []interface{}{arg1, arg2, arg3})
	fake.runScriptWithEnvMutex.Unlock()
	if fake.RunScriptWithEnvStub != nil {
		return fake.RunScriptWithEnvStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.runScriptWithEnvReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRemoteRunner) RunScriptWithEnvCallCount() int {
	fake.runScriptWithEnvMutex.RLock()
	defer fake.runScriptWithEnvMutex.RUnlock()
	return len(fake.runScriptWithEnvArgsForCall)
}

func (fake *FakeRemoteRunner) RunScriptWithEnvCalls(stub func(string, map[string]string, string) (string, error)) {
	fake.runScriptWithEnvMutex.Lock()
	defer fake.runScriptWithEnvMutex.Unlock()
	fake.RunScriptWithEnvStub = stub
}

func (fake *FakeRemoteRunner) RunScriptWithEnvArgsForCall(i int) (string, map[string]string, string) {
	fake.runScriptWithEnvMutex.RLock()
	defer fake.runScriptWithEnvMutex.RUnlock()
	argsForCall := fake.runScriptWithEnvArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeRemoteRunner) RunScriptWithEnvReturns(result1 string, result2 error) {
	fake.runScriptWithEnvMutex.Lock()
	defer fake.runScriptWithEnvMutex.Unlock()
	fake.RunScriptWithEnvStub = nil
	fake.runScriptWithEnvReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteRunner) RunScriptWithEnvReturnsOnCall(i int, result1 string, result2 error) {
	fake.runScriptWithEnvMutex.Lock()
	defer fake.runScriptWithEnvMutex.Unlock()
	fake.RunScriptWithEnvStub = nil
	if fake.runScriptWithEnvReturnsOnCall == nil {
		fake.runScriptWithEnvReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.runScriptWithEnvReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteRunner) SizeOf(arg1 string) (string, error) {
	fake.sizeOfMutex.Lock()
	ret, specificReturn := fake.sizeOfReturnsOnCall[len(fake.sizeOfArgsForCall)]
	fake.sizeOfArgsForCall = append(fake.sizeOfArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("SizeOf", []interface{}{arg1})
	fake.sizeOfMutex.Unlock()
	if fake.SizeOfStub != nil {
		return fake.SizeOfStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.sizeOfReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeRemoteRunner) SizeOfCallCount() int {
	fake.sizeOfMutex.RLock()
	defer fake.sizeOfMutex.RUnlock()
	return len(fake.sizeOfArgsForCall)
}

func (fake *FakeRemoteRunner) SizeOfCalls(stub func(string) (string, error)) {
	fake.sizeOfMutex.Lock()
	defer fake.sizeOfMutex.Unlock()
	fake.SizeOfStub = stub
}

func (fake *FakeRemoteRunner) SizeOfArgsForCall(i int) string {
	fake.sizeOfMutex.RLock()
	defer fake.sizeOfMutex.RUnlock()
	argsForCall := fake.sizeOfArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRemoteRunner) SizeOfReturns(result1 string, result2 error) {
	fake.sizeOfMutex.Lock()
	defer fake.sizeOfMutex.Unlock()
	fake.SizeOfStub = nil
	fake.sizeOfReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteRunner) SizeOfReturnsOnCall(i int, result1 string, result2 error) {
	fake.sizeOfMutex.Lock()
	defer fake.sizeOfMutex.Unlock()
	fake.SizeOfStub = nil
	if fake.sizeOfReturnsOnCall == nil {
		fake.sizeOfReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.sizeOfReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeRemoteRunner) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.archiveAndDownloadMutex.RLock()
	defer fake.archiveAndDownloadMutex.RUnlock()
	fake.checksumDirectoryMutex.RLock()
	defer fake.checksumDirectoryMutex.RUnlock()
	fake.connectedUsernameMutex.RLock()
	defer fake.connectedUsernameMutex.RUnlock()
	fake.createDirectoryMutex.RLock()
	defer fake.createDirectoryMutex.RUnlock()
	fake.directoryExistsMutex.RLock()
	defer fake.directoryExistsMutex.RUnlock()
	fake.extractAndUploadMutex.RLock()
	defer fake.extractAndUploadMutex.RUnlock()
	fake.findFilesMutex.RLock()
	defer fake.findFilesMutex.RUnlock()
	fake.isWindowsMutex.RLock()
	defer fake.isWindowsMutex.RUnlock()
	fake.removeDirectoryMutex.RLock()
	defer fake.removeDirectoryMutex.RUnlock()
	fake.runScriptMutex.RLock()
	defer fake.runScriptMutex.RUnlock()
	fake.runScriptWithEnvMutex.RLock()
	defer fake.runScriptWithEnvMutex.RUnlock()
	fake.sizeOfMutex.RLock()
	defer fake.sizeOfMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeRemoteRunner) recordInvocation(key string, args []interface{}) {
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

var _ ssh.RemoteRunner = new(FakeRemoteRunner)

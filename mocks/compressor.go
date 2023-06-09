// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	os "os"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	compression "github.com/konstellation-io/kli/internal/krt/fileutils/compression"
)

// MockCompressor is a mock of Compressor interface.
type MockCompressor struct {
	ctrl     *gomock.Controller
	recorder *MockCompressorMockRecorder
}

// MockCompressorMockRecorder is the mock recorder for MockCompressor.
type MockCompressorMockRecorder struct {
	mock *MockCompressor
}

// NewMockCompressor creates a new mock instance.
func NewMockCompressor(ctrl *gomock.Controller) *MockCompressor {
	mock := &MockCompressor{ctrl: ctrl}
	mock.recorder = &MockCompressorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCompressor) EXPECT() *MockCompressorMockRecorder {
	return m.recorder
}

// Extract mocks base method.
func (m *MockCompressor) Extract(krtPath, dst string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Extract", krtPath, dst)
	ret0, _ := ret[0].(error)
	return ret0
}

// Extract indicates an expected call of Extract.
func (mr *MockCompressorMockRecorder) Extract(krtPath, dst interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Extract", reflect.TypeOf((*MockCompressor)(nil).Extract), krtPath, dst)
}

// NewCompressedFile mocks base method.
func (m *MockCompressor) NewCompressedFile(filePath string) (compression.CompressedFile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewCompressedFile", filePath)
	ret0, _ := ret[0].(compression.CompressedFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewCompressedFile indicates an expected call of NewCompressedFile.
func (mr *MockCompressorMockRecorder) NewCompressedFile(filePath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCompressedFile", reflect.TypeOf((*MockCompressor)(nil).NewCompressedFile), filePath)
}

// MockCompressedFile is a mock of CompressedFile interface.
type MockCompressedFile struct {
	ctrl     *gomock.Controller
	recorder *MockCompressedFileMockRecorder
}

// MockCompressedFileMockRecorder is the mock recorder for MockCompressedFile.
type MockCompressedFileMockRecorder struct {
	mock *MockCompressedFile
}

// NewMockCompressedFile creates a new mock instance.
func NewMockCompressedFile(ctrl *gomock.Controller) *MockCompressedFile {
	mock := &MockCompressedFile{ctrl: ctrl}
	mock.recorder = &MockCompressedFileMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCompressedFile) EXPECT() *MockCompressedFileMockRecorder {
	return m.recorder
}

// AddFile mocks base method.
func (m *MockCompressedFile) AddFile(file os.FileInfo, path, relativePath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFile", file, path, relativePath)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddFile indicates an expected call of AddFile.
func (mr *MockCompressedFileMockRecorder) AddFile(file, path, relativePath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFile", reflect.TypeOf((*MockCompressedFile)(nil).AddFile), file, path, relativePath)
}

// Compress mocks base method.
func (m *MockCompressedFile) Compress() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Compress")
	ret0, _ := ret[0].(error)
	return ret0
}

// Compress indicates an expected call of Compress.
func (mr *MockCompressedFileMockRecorder) Compress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Compress", reflect.TypeOf((*MockCompressedFile)(nil).Compress))
}

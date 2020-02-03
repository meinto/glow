package mock_git

import "github.com/meinto/glow/git"

type MockServiceInterface interface {
	EXPECT() *MockServiceMockRecorder
	git.Service
}

type MockNativeServiceInterface interface {
	EXPECT() *MockNativeServiceMockRecorder
	git.NativeService
}

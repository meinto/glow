package mock_git

type MockServiceInterface interface {
	EXPECT() *MockServiceMockRecorder
}

type MockNativeServiceInterface interface {
	EXPECT() *MockNativeServiceMockRecorder
}

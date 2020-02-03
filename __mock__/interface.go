package mock_glow

import glow "github.com/meinto/glow"

type MockBranchInterface interface {
	EXPECT() *MockBranchMockRecorder
	glow.Branch
}

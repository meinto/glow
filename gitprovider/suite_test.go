package gitprovider_test

import (
	"net/http"
	"testing"

	. "github.com/meinto/glow/gitprovider"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGlow(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Git Provider Suite")
}

type MockHTTPClient struct {
	HttpClient
	RequestIntercaptionCallback func(req *http.Request)
	RequestCounter              int
}

func (c *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.RequestIntercaptionCallback(req)
	c.RequestCounter++
	return &http.Response{Body: http.NoBody}, nil
}

func (c *MockHTTPClient) SetRequestIntercaptionCallback(cb func(req *http.Request)) {
	c.RequestIntercaptionCallback = cb
}

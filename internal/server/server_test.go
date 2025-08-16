package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	var api LibraryApi

	testRouter := gin.New()
	gin.SetMode(gin.ReleaseMode)

	testRouter.GET("/health", api.healthHandler)

	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	req := resty.New().R()

	req.Method = http.MethodGet
	req.URL = httpTest.URL + "/health"

	resp, err := req.Send()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

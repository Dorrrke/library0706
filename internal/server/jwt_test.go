package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthMiddleware(t *testing.T) {
	token, err := createTestToken()
	assert.NoError(t, err)

	api := LibraryApi{}
	testRouter := gin.New()
	gin.SetMode(gin.ReleaseMode)

	testRouter.GET("/health", api.JWTAuthMiddleware(), api.healthHandler)

	httpTest := httptest.NewServer(testRouter)
	defer httpTest.Close()

	type want struct {
		statuCode int
		body      string
	}

	type test struct {
		name    string
		request string
		method  string
		token   string
		want    want
	}

	tests := []test{
		{
			name:    "Case 1 - valid token",
			request: "/health",
			method:  http.MethodGet,
			token:   token,
			want: want{
				statuCode: http.StatusOK,
				body:      `{"status":"ok"}`,
			},
		},
		{
			name:    "Case 2 - invalid token",
			request: "/health",
			method:  http.MethodGet,
			token:   "invalid_token",
			want: want{
				statuCode: http.StatusUnauthorized,
				body:      `{"error":"token is malformed: token contains an invalid number of segments"}`,
			},
		},
		{
			name:    "Case 3 - empty token",
			request: "/health",
			method:  http.MethodGet,
			token:   "",
			want: want{
				statuCode: http.StatusUnauthorized,
				body:      `{"error":"unauthorized"}`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()

			req.Method = tc.method
			req.URL = httpTest.URL + tc.request
			req.Header.Set("Authorization", "Bearer "+tc.token)

			resp, err := req.Send()
			assert.NoError(t, err)

			assert.Equal(t, tc.want.statuCode, resp.StatusCode())
			assert.Equal(t, tc.want.body, resp.String())
		})
	}
}

func createTestToken() (string, error) {
	userID := "test_user_id"

	token, _, err := createToken(userID, accessTokenDuration)
	if err != nil {
		return "", err
	}
	return token, nil
}

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dorrrke/library0706/internal/domain/models"
	"github.com/Dorrrke/library0706/internal/server/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetBook(t *testing.T) {
	api := LibraryApi{}
	testRouter := gin.New()
	gin.SetMode(gin.ReleaseMode)

	testRouter.POST("/books/get/:bookID", api.getBook)
	testSrv := httptest.NewServer(testRouter)
	defer testSrv.Close()

	type want struct {
		resultMsg string
		status    int
	}

	type test struct {
		name      string
		request   string
		method    string
		bookID    string
		bookModel models.Book
		want      want
	}

	tests := []test{
		{
			name:    "Case 1 - get book",
			request: "/books/get/%s",
			method:  "POST",
			bookID:  "12356789-1234-1234-1234-123456789012",
			bookModel: models.Book{
				BookID:      "12356789-1234-1234-1234-123456789012",
				Author:      "Джордж Оруэлл",
				Lable:       "1984",
				Description: `Роман-антиутопия о тоталитарном обществе, где правительство контролирует каждую мысль гражданина.`,
				Genre:       "Антиутопия",
				WritedAt:    "1949",
				Count:       5,
			},
			want: want{
				resultMsg: `{"book_id":"12356789-1234-1234-1234-123456789012","author":"Джордж Оруэлл","lable":"1984","description":"Роман-антиутопия о тоталитарном обществе, где правительство контролирует каждую мысль гражданина.","genre":"Антиутопия","writed_at":"1949","count":5}`,
				status:    200,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewRepository(t)
			mockRepo.On("GetBook", tc.bookID).Return(tc.bookModel, nil)

			api.db = mockRepo

			req := resty.New().R()

			req.Method = tc.method
			req.URL = testSrv.URL + fmt.Sprintf(tc.request, tc.bookID)

			resp, err := req.Send()
			assert.NoError(t, err)

			assert.Equal(t, tc.want.status, resp.StatusCode())
			assert.Equal(t, tc.want.resultMsg, resp.String())
		})
	}

}

func TestNewBook(t *testing.T) {
	api := LibraryApi{}
	testRouter := gin.New()
	gin.SetMode(gin.ReleaseMode)

	testRouter.POST("/create", api.newBook)
	testSrv := httptest.NewServer(testRouter)
	defer testSrv.Close()

	type want struct {
		book   models.Book
		status int
	}

	type test struct {
		name    string
		request string
		method  string
		book    models.Book
		want    want
	}

	tests := []test{
		{
			name:    "Case 1 - new book",
			request: "/create",
			method:  http.MethodPost,
			book: models.Book{
				BookID:      "12356789-1234-1234-1234-123456789012",
				Author:      "Джордж Оруэлл",
				Lable:       "1984",
				Description: `Роман-антиутопия о тоталитарном обществе, где правительство контролирует каждую мысль гражданина.`,
				Genre:       "Антиутопия",
				WritedAt:    "1949",
				Count:       5,
			},
			want: want{
				book: models.Book{
					BookID:      "12356789-1234-1234-1234-123456789012",
					Author:      "Джордж Оруэлл",
					Lable:       "1984",
					Description: `Роман-антиутопия о тоталитарном обществе, где правительство контролирует каждую мысль гражданина.`,
					Genre:       "Антиутопия",
					WritedAt:    "1949",
					Count:       5,
				},
				status: http.StatusCreated,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewRepository(t)
			mockRepo.On("SaveBook", mock.Anything).Return(nil)

			api.db = mockRepo

			req := resty.New().R()

			req.Method = tc.method
			req.URL = testSrv.URL + tc.request

			bodyJson, err := json.Marshal(tc.book)
			assert.NoError(t, err)

			req.Body = bodyJson

			resp, err := req.Send()
			assert.NoError(t, err)

			nBody := resp.String()
			var book models.Book
			err = json.Unmarshal([]byte(nBody), &book)
			assert.NoError(t, err)

			assert.Equal(t, tc.want.status, resp.StatusCode())

			assert.NotEmpty(t, book.BookID)
			assert.Equal(t, tc.want.book.Author, book.Author)
			assert.Equal(t, tc.want.book.Lable, book.Lable)
			assert.Equal(t, tc.want.book.Description, book.Description)
			assert.Equal(t, tc.want.book.Genre, book.Genre)
			assert.Equal(t, tc.want.book.WritedAt, book.WritedAt)
			assert.Equal(t, tc.want.book.Count, book.Count)
		})
	}
}

func BenchmarkNewBook(b *testing.B) {
	api := LibraryApi{}
	testRouter := gin.New()
	gin.SetMode(gin.ReleaseMode)

	testRouter.POST("/create", api.newBook)
	testSrv := httptest.NewServer(testRouter)
	defer testSrv.Close()

	mockRepo := mocks.NewRepository(b)
	mockRepo.On("SaveBook", mock.Anything).Return(nil)

	api.db = mockRepo

	req := resty.New().R()

	req.Method = http.MethodPost
	req.URL = testSrv.URL + "/create"

	book := models.Book{
		BookID:      "12356789-1234-1234-1234-123456789012",
		Author:      "Джордж Оруэлл",
		Lable:       "1984",
		Description: `Роман-антиутопия о тоталитарном обществе, где правительство контролирует каждую мысль гражданина.`,
		Genre:       "Антиутопия",
		WritedAt:    "1949",
		Count:       5,
	}

	bodyJson, err := json.Marshal(book)
	assert.NoError(b, err)

	req.Body = bodyJson

	for range b.N {
		_, err = req.Send()
	}
}

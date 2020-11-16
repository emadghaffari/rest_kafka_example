package twitter

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/emadghaffari/res_errors/errors"
	twitterModel "github.com/emadghaffari/rest_kafka_example/model/twitter"
	service "github.com/emadghaffari/rest_kafka_example/services/twitter"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type twitterServiceMock struct {
	SearchmockFuc func() error
	StoremockFuc  func() (*twitter.Tweet, error)
}

// Search func
func (tw *twitterServiceMock) Search(text string) error {
	text = strings.TrimSpace(text)
	if text == "" {
		return errors.HandlerBadRequest("invalid text")
	}
	return tw.SearchmockFuc()
}

// Store func
func (tw *twitterServiceMock) Store(request twitterModel.StoreRequest) (*twitter.Tweet, error) {
	request.Text = strings.TrimSpace(request.Text)
	if request.Text == "" {
		return nil, errors.HandlerBadRequest("invalid text")
	}

	return tw.StoremockFuc()
}

func TestSuccessSearch(t *testing.T) {
	mock := &twitterServiceMock{}
	mock.SearchmockFuc = func() error {
		return nil
	}
	service.TwitterService = mock

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.POST("/search", Search)

	c.Request, _ = http.NewRequest(http.MethodPost, "/search", bytes.NewBuffer([]byte("{}")))

	fmt.Println(c.Request)
	r.ServeHTTP(w, c.Request)

	request := twitterModel.SearchRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(w.Body)
	assert.Equal(t, w.Code, http.StatusOK)
}

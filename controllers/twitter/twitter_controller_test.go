package twitter

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	twitterModel "github.com/emadghaffari/rest_kafka_example/model/twitter"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSearch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{gin.Param{Key: "request", Value: "golang"}}
	request := twitterModel.SearchRequest{}
	c.ShouldBindJSON(&request)
	fmt.Println(request)
	Search(c)

	if w.Code != 200 {
		b, _ := ioutil.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}
}
func TestStore(t *testing.T) {

}

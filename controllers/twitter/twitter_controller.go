package twitter

import (
	"fmt"
	"net/http"

	"github.com/emadghaffari/res_errors/errors"
	twitterModel "github.com/emadghaffari/rest_kafka_example/model/twitter"
	"github.com/emadghaffari/rest_kafka_example/services/twitter"
	"github.com/gin-gonic/gin"
)

// Search func
func Search(c *gin.Context) {
	request := twitterModel.SearchRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		resErr := errors.HandlerBadRequest("Invalid JSON Body.")
		c.JSON(resErr.Status(), resErr.Message())
		return
	}
	err := twitter.Search(request.Request)
	if err != nil {
		resErr := errors.HandlerInternalServerError("internal search Error", err)
		c.JSON(resErr.Status(), resErr.Message())
		return
	}
	c.JSON(http.StatusOK, "Success")
}

// Store new Tweet
func Store(c *gin.Context) {
	request := twitterModel.StoreRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		resErr := errors.HandlerBadRequest("Invalid JSON Body.")
		c.JSON(resErr.Status(), resErr.Message())
		return
	}

	result, err := twitter.Store(request)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, result)

}

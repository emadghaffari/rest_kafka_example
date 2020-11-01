package twitter

import (
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/emadghaffari/res_errors/errors"
	twitterModel "github.com/emadghaffari/rest_kafka_example/model/twitter"
)

// Search func
func Search(text string) (error) {
	text = strings.TrimSpace(text)
	if text == ""{
		return errors.HandlerBadRequest("invalid text")
	}
	return twitterModel.Search(text)
}

// Store func
func Store(request twitterModel.StoreRequest) (*twitter.Tweet ,error) {
	request.Text = strings.TrimSpace(request.Text)
	if request.Text == ""{
		return nil, errors.HandlerBadRequest("invalid text")
	}
	result,err := twitterModel.Store(request)
	if err != nil {
		return nil,err
	}
	return result,nil 
}
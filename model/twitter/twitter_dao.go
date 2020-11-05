package twitter

import (
	"encoding/json"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/emadghaffari/rest_kafka_example/databases/kafka"
)

// Search func
func Search(text string) error {
	client, err := Account.GetClient()
	if err != nil {
		return err
	}

	result, err := Account.Search(client, text)
	if err != nil {
		return err
	}
	var items []string
	for _, item := range result.Statuses {
		res, _ := json.Marshal(item.User)
		items = append(items, string(res))
	}
	kafka.Producer(items)
	return nil
}

// Store func
func Store(request StoreRequest) (*twitter.Tweet, error) {
	client, err := Account.GetClient()
	if err != nil {
		return nil, err
	}
	result, _, err := Account.NewTweet(client, request.Text, nil)
	if err != nil {
		return nil, err
	}
	kafka.Producer([]string{result.Text})
	return result, nil
}

package twitter

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/emadghaffari/res_errors/errors"
	"github.com/spf13/viper"
)

var (
	// Account var
	Account credentialsInterface = &credentials{}
)

type credentialsInterface interface {
	Init()
	GetClient() (*twitter.Client, error)
	Search(client *twitter.Client, query string) (*twitter.Search, error)
	NewTweet(client *twitter.Client, text string, params *twitter.StatusUpdateParams) (*twitter.Tweet, *http.Response, error)
}

// credentials stores all of our access/consumer tokens
// and secret keys needed for authentication against
// the twitter REST API.
type credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// SearchRequest struct
type SearchRequest struct {
	Request string `'json:"request"`
	Type    string `'json:"type"`
}

// StoreRequest struct
type StoreRequest struct {
	Text  string `'json:"text"`
	Title string `'json:"title"`
}

// Init func
func (creds *credentials) Init() {
	creds.ConsumerKey = viper.GetString("twitter.consumerKey")
	creds.ConsumerSecret = viper.GetString("twitter.consumerSecret")
	creds.AccessToken = viper.GetString("twitter.accessToken")
	creds.AccessTokenSecret = viper.GetString("twitter.accessTokenSecret")
}

// GetClient is a helper function that will return a twitter client
// that we can subsequently use to send tweets, or to stream new tweets
// this will take in a pointer to a Credential struct which will contain
// everything needed to authenticate and return a pointer to a twitter Client
// or an error
func (creds credentials) GetClient() (*twitter.Client, error) {
	// Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	// Pass in your Access Token and your Access Token Secret
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	// we can retrieve the user and verify if the credentials
	// we have used successfully allow us to log in!
	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("twitter: 215 Bad Authentication data")
		return nil, err
	}

	return client, nil
}

// Search func
func (creds credentials) Search(client *twitter.Client, query string) (*twitter.Search, error) {
	search, hresp, err := client.Search.Tweets(
		&twitter.SearchTweetParams{
			Query: query,
		},
	)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("internal search error")
		return nil, errors.HandlerBadRequest("Error in Search Tweet Params")
	}
	if hresp.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Error in http request for search twitter")
		return nil, errors.HandlerInternalServerError(fmt.Sprintf("Error in http request for search: %v", hresp), nil)
	}

	return search, nil
}

// NewTweet func
func (creds credentials) NewTweet(client *twitter.Client, text string, params *twitter.StatusUpdateParams) (*twitter.Tweet, *http.Response, error) {
	tweet, resp, err := client.Statuses.Update(text, params)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	return tweet, resp, nil
}

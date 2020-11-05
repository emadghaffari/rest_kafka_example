package elasticsearch

import (
	"context"
	"fmt"
	"time"

	"github.com/emadghaffari/res_errors/errors"
	"github.com/emadghaffari/res_errors/logger"
	"github.com/olivere/elastic"
	"github.com/spf13/viper"
)

var (
	// Client var
	Client elasticsearchInterface = &elasticsearchClient{}
)

type elasticsearchInterface interface {
	Index(string, string, string, interface{}) (*elastic.IndexResponse, errors.ResError)
	Get(string, string, string) (*elastic.GetResult, errors.ResError)
	SetClient(*elastic.Client)
	Search(string, elastic.Query) (*elastic.SearchResult, errors.ResError)
	Delete(string, string, string) (*elastic.DeleteResponse, errors.ResError)
}

type elasticsearchClient struct {
	client *elastic.Client
}

// Init func
// init the elasticsearch with default setup
func Init() {
	logger := logger.GetLogger()
	client, err := elastic.NewClient(
		elastic.SetURL(viper.GetString("elastic.url")),
		elastic.SetBasicAuth(viper.GetString("elastic.auth.username"), viper.GetString("elastic.auth.password")),
		elastic.SetHealthcheckInterval(50*time.Second),
		elastic.SetErrorLog(logger),
		elastic.SetInfoLog(logger),
	)
	if err != nil {
		panic(err)
	}
	Client.SetClient(client)
}

// SetClient func
// Set New Client
func (c *elasticsearchClient) SetClient(client *elastic.Client) {
	c.client = client
}

func (c *elasticsearchClient) Index(id string, index string, docType string, doc interface{}) (*elastic.IndexResponse, errors.ResError) {
	ctx := context.Background()
	elk, err := c.client.Index().
		Index(index).
		Id(id).
		BodyJson(doc).
		Type(docType).
		Do(ctx)
	if err != nil {
		logger.Error("error in index esClient", err)
		return nil, errors.HandlerInternalServerError("internal ELK error in Index", err)
	}
	return elk, nil
}
func (c *elasticsearchClient) Get(index string, docType string, id string) (*elastic.GetResult, errors.ResError) {
	ctx := context.Background()
	result, err := c.client.Get().
		Index(index).
		Type(docType).
		Id(id).
		Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error in index esClient %s", id), err)
		return nil, errors.HandlerInternalServerError("internal ELK error in Index", err)
	}
	return result, nil
}
func (c *elasticsearchClient) Search(index string, query elastic.Query) (*elastic.SearchResult, errors.ResError) {
	ctx := context.Background()
	result, err := c.client.Search(index).Query(query).RestTotalHitsAsInt(true).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error in search esClient %s", index), err)
		return nil, errors.HandlerInternalServerError("internal ELK error in Index", err)
	}
	return result, nil
}
func (c *elasticsearchClient) Delete(index string, docType string, id string) (*elastic.DeleteResponse, errors.ResError) {
	ctx := context.Background()
	result, err := c.client.Delete().
		Index(index).
		Type(docType).
		Id(id).
		Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error in index esClient %s", id), err)
		return nil, errors.HandlerInternalServerError("internal ELK error in Index", err)
	}
	return result, nil
}

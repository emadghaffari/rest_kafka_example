package elasticsearch

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/emadghaffari/res_errors/errors"
	"github.com/emadghaffari/rest_kafka_example/databases/elasticsearch"
)

const (
	indexES = "twitter"
	docType = "_doc"
)

// Save method
// store new item
func Save(id string, i interface{}) errors.ResError {
	_, err := elasticsearch.Client.Index(id, indexES, docType, i)
	if err != nil {
		return err
	}
	return nil
}

// Get func
// get item with id
func Get(i string) (interface{}, errors.ResError) {
	result, err := elasticsearch.Client.Get(indexES, docType, i)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, errors.HandlerNotFoundError(fmt.Sprintf("item not found %s", i))
		}
		return nil, err
	}
	if !result.Found {
		return nil, errors.HandlerNotFoundError(fmt.Sprintf("item not found %s", i))
	}

	bytes, marshalErr := result.Source.MarshalJSON()
	if marshalErr != nil {
		return nil, errors.HandlerInternalServerError(fmt.Sprintf("error in MarshalJSON from DB %s", i), err)
	}
	var unmarshaled interface{}
	if err := json.Unmarshal(bytes, &unmarshaled); err != nil {
		return nil, errors.HandlerInternalServerError(fmt.Sprintf("error in unmarshal data %s", i), err)
	}
	return unmarshaled, nil
}

// Delete func
// get item with id
func Delete(i string) errors.ResError {
	result, err := elasticsearch.Client.Delete(indexES, docType, i)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return errors.HandlerNotFoundError(fmt.Sprintf("item not found %s", i))
		}
		return err
	}
	if result.Shards.Successful > 0 {
		return errors.HandlerNotFoundError(fmt.Sprintf("item not found %s", i))
	}

	return nil
}

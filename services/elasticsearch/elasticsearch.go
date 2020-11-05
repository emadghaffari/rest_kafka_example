package elasticsearch

import (
	"github.com/emadghaffari/res_errors/errors"
	"github.com/emadghaffari/rest_kafka_example/model/elasticsearch"
)

// Save method
// store new item
func Save(id string, i interface{}) errors.ResError {
	err := elasticsearch.Save(id, i)
	if err != nil {
		return err
	}
	return nil
}

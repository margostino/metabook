package db

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/margostino/metabook/common"
	"log"
)

const INDEX = "books"

var clientConnection *elasticsearch.Client

func Connect() {
	//config := elasticsearch.Config{
	//	Addresses: []string{
	//		"https://localhost:9200",
	//	},
	//}
	//client, err := elasticsearch.NewClient(config)
	client, err := elasticsearch.NewDefaultClient()
	common.Check(err)

	response, err := client.Info()
	common.Check(err)

	defer response.Body.Close()
	log.Println(response)

	clientConnection = client
}

func Index(document map[string]string) {
	data, err := json.Marshal(document)
	common.Check(err)
	request := esapi.IndexRequest{
		Index:   INDEX,
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	response, err := request.Do(context.Background(), clientConnection)
	common.Check(err)
	defer response.Body.Close()

	if response.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", response.Status(), 1)
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			log.Printf("[%s] %s; version=%d", response.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

}

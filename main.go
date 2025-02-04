package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"strconv"
	"strings"
	"time"
)

// I write this implementation after reading official docs of using elasticsearch with go,
// You can read that here: https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/getting-started-go.html

type ProductDocument struct {
	Id          int        `json:"id"`
	Seller_id   string     `json:"seller_id"`
	Name        string     `json:"name"`
	Category    string     `json:"category"`
	Quantity    int        `json:"quantity"`
	Price       float64    `json:"price"`
	Weight      int        `json:"weight"`
	Size        string     `json:"size"`
	Status      string     `json:"status"`
	Description string     `json:"description"`
	Created_at  *time.Time `json:"created_at"`
	Updated_at  *time.Time `json:"updated_at"`
}

func main() {
	cfg := elasticsearch.Config{
		// no replica so only 1
		Addresses: []string{
			"http://localhost:9200",
		},
		Username: "elastic",
		Password: "elastic123",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalln("Error creating elasticsearch client: ", err)
	}

	_, err = es.Info()
	if err != nil {
		log.Fatalln("Error getting response/ping from elasticsearch: ", err)
	}

	now := time.Now()

	doc := ProductDocument{
		Id:          1,
		Seller_id:   "ea838b0c-a235-48c0-a84f-90a8aca284e2",
		Name:        "Martabak Manis",
		Category:    "Makanan",
		Quantity:    10,
		Price:       10.000,
		Weight:      10,
		Size:        "XXL",
		Status:      "Ready",
		Description: "Martabak manis adalah kudapan sejenis panekuk yang biasa dijajakan di pinggir jalan di seluruh Indonesia, Malaysia, Brunei Darussalam, Filipina dan Singapura.[5]",
		Created_at:  &now,
		Updated_at:  &now,
	}

	// with elasticsearch you can with search fuzzy(typo) query, autocomplete, suggest, filtered search,

	// When insert some document you would need index, if the index haven't created then the index will be created automatically by elasticsearch
	// Usually people set Product id or User id in PostgreSQL as document id in Elasticsearch
	// 1. by default if no id is specified when insert a document, the document id will be random
	// 2. insert document and set the document id as product id
	// 3. get one document by document id (in this case its 1)
	// 4. search all document in "products" index, (search usually dont use document id because you can just the previous implementation for that)
	// 5 .search document with wildcard (not recommended, due to slow performance) its like %elastic% in relational database so all data that contains "elastic" in that field will be retrieve
	// 6. search document with match value data (recommended), "match" mean will retrieve data that looks like same to the query data for example "dika","DIKA","dIka","dikA"
	// 7. search document with term value data (recommended), "term" mean will retrieve data that exactly looks like that for example "dika" only get retrieve by "dika" if query "dikA" wont retrieve "dika"
	// 8. search document with multimatch (recommended), "multimatch" mean it will search that query in many field
	// ?. search document with bool ?
	// 9. update docoument by document id
	// 10. delete document by document id

	insertDocument(es, doc)

	indexName := "products"

	insertDocumentWithDocumentID(es, indexName, doc)

	getDocumentFromIndexByDocumentID(es, indexName, doc.Id)

	searchAllDocumentFromIndex(es, indexName)

	fieldName := "name"
	matchQuery := "Martabak Manis"

	searchDocumentWithWildcardQuery(es, indexName, fieldName, matchQuery)

	searchDocumentWithMatchQuery(es, indexName, fieldName, matchQuery)

	searchDocumentWithTermQuery(es, indexName, fieldName, matchQuery)

	fieldName2 := "category"
	fieldName3 := "description"

	searchDocumentWithMultiMatch(es, indexName, fieldName, fieldName2, fieldName3, matchQuery)

	updatedFields := map[string]interface{}{
		"status": "Not Ready",
		"size":   "L",
	}

	updateDocumentByDocumentID(es, indexName, doc.Id, updatedFields)

	getDocumentFromIndexByDocumentID(es, indexName, doc.Id)

	deleteDocumentByDocumentID(es, indexName, doc.Id)

	getDocumentFromIndexByDocumentID(es, indexName, doc.Id)

}

func insertDocument(es *elasticsearch.Client, doc ProductDocument) {
	data, _ := json.Marshal(doc)
	indexResponse, err := es.Index("my_products", bytes.NewReader(data))
	if err != nil {
		log.Fatalln("Error sending a post request to index a document to elasticsearch: ", err.Error())
	}

	defer indexResponse.Body.Close()

	if indexResponse.IsError() {
		log.Fatalln("Error inserting a document in elasticsearch (elasticsearch response): ", indexResponse.Status())
	}

	log.Println("Success to insert a document in elasticsearch")
}

func insertDocumentWithDocumentID(es *elasticsearch.Client, indexName string, doc ProductDocument) {
	data, _ := json.Marshal(doc)

	indexResponse, err := es.Index(
		indexName,
		bytes.NewReader(data),
		es.Index.WithDocumentID(strconv.Itoa(doc.Id)),
	)

	if err != nil {
		log.Fatalln("Error sending a post request to index a document to elasticsearch: ", err.Error())
	}

	defer indexResponse.Body.Close()

	if indexResponse.IsError() {
		log.Fatalln("Error indexing a document in elasticsearch (elasticsearch response): ", indexResponse.Status())
	}

	log.Println("Success to index document in elasticsearch")
}

func getDocumentFromIndexByDocumentID(es *elasticsearch.Client, indexName string, productID int) {
	res, err := es.Get(indexName, strconv.Itoa(productID))
	if err != nil {
		log.Fatalln("Error sending a get request to get a document from index to elasticsearch: ", err.Error())
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Fatalln("Error get a document in elasticsearch (elasticsearch response): ", res.Status())
	}

	var searchResult map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&searchResult)
	if err != nil {
		log.Fatalln("Error decode the response body: ", err)
	}

	source := searchResult["_source"].(map[string]interface{})

	log.Println("Document: ", source)

	log.Println("Success to get getDocumentFromIndexByDocumentID")
}

func searchAllDocumentFromIndex(es *elasticsearch.Client, indexName string) {
	query := `{ "query":{ 
				"match_all": {} 
               }
              }`

	res, err := es.Search(
		es.Search.WithIndex(indexName),
		es.Search.WithBody(strings.NewReader(query)),
	)

	if err != nil {
		log.Fatalln("Error sending a post/get request to search a document to elasticsearch: ", err.Error())
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Fatalln("Error search a document in elasticsearch (elasticsearch response): ", res.Status())
	}

	var searchResult map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&searchResult)
	if err != nil {
		log.Fatalln("Error decode document: ", err)
	}

	// access all the hits in response object from elasticsearch result
	hits := searchResult["hits"].(map[string]interface{})["hits"].([]interface{})

	// loop through each document in the hits
	for _, hit := range hits {
		document := hit.(map[string]interface{})
		source := document["_source"].(map[string]interface{})
		fmt.Printf("Document: %v\n", source)
	}

	log.Println("Success to searchAllDocumentFromIndex")
}

func searchDocumentWithWildcardQuery(es *elasticsearch.Client, indexName string, fieldName string, matchQuery string) {
	query := fmt.Sprintf(`{
		"query": {
			"wildcard": {
				"%s": "%s"
			}
		}
	}`, fieldName, matchQuery)

	res, err := es.Search(
		es.Search.WithIndex(indexName),
		es.Search.WithBody(strings.NewReader(query)),
	)

	if err != nil {
		log.Fatalln("Error sending a post request to search a document to elasticsearch: ", err)
	}

	if res.IsError() {
		log.Fatalln("Error search a document in elasticsearch (elasticsearch response): ", res.Status())
	}

	var searchResult map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&searchResult)
	if err != nil {
		log.Fatalln("Error decode document: ", err)
	}

	hits := searchResult["hits"].(map[string]interface{})["hits"].([]interface{})

	for _, hit := range hits {
		document := hit.(map[string]interface{})
		source := document["_source"].(map[string]interface{})
		fmt.Printf("Document: %v\n", source)
	}

	log.Println("Success to searchDocumentWithWildcardQuery")
}

func searchDocumentWithMatchQuery(es *elasticsearch.Client, indexName string, fieldName string, matchQuery string) {
	query := fmt.Sprintf(`{
		"query": {
			"match": {
				"%s": "%s"
			}
		}
	}`, fieldName, matchQuery)

	res, err := es.Search(
		es.Search.WithIndex(indexName),
		es.Search.WithBody(strings.NewReader(query)),
	)

	if err != nil {
		log.Fatalln("Error sending a post request to search a document to elasticsearch: ", err)
	}

	if res.IsError() {
		log.Fatalln("Error search a document in elasticsearch (elasticsearch response): ", res.Status())
	}

	var searchResult map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&searchResult)
	if err != nil {
		log.Fatalln("Error decode document: ", err)
	}

	hits := searchResult["hits"].(map[string]interface{})["hits"].([]interface{})

	for _, hit := range hits {
		document := hit.(map[string]interface{})
		source := document["_source"].(map[string]interface{})
		fmt.Printf("Document: %v\n", source)
	}

	log.Println("Success to searchDocumentWithMatchQuery")
}

func searchDocumentWithTermQuery(es *elasticsearch.Client, indexName string, fieldName string, matchQuery string) {
	query := fmt.Sprintf(`{
		"query": {
			"term": {
				"%s.keyword": "%s"
			}
		}
	}`, fieldName, matchQuery)

	res, err := es.Search(
		es.Search.WithIndex(indexName),
		es.Search.WithBody(strings.NewReader(query)),
	)

	if err != nil {
		log.Fatalln("Error sending a post request to search a document to elasticsearch: ", err)
	}

	if res.IsError() {
		log.Fatalln("Error search a document in elasticsearch (elasticsearch response): ", res.Status())
	}

	var searchResult map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&searchResult)
	if err != nil {
		log.Fatalln("Error decode document: ", err)
	}

	hits := searchResult["hits"].(map[string]interface{})["hits"].([]interface{})

	for _, hit := range hits {
		document := hit.(map[string]interface{})
		source := document["_source"].(map[string]interface{})
		fmt.Printf("Document: %v\n", source)
	}

	log.Println("Success to searchDocumentWithTermQuery")
}

func searchDocumentWithMultiMatch(es *elasticsearch.Client, indexName string, fieldName string, fieldName2 string, fieldName3 string, matchQuery string) {
	query := fmt.Sprintf(`{
		"query": {
			"multi_match": {
				"query": "%s",
				"fields": ["%s", "%s", "%s"]
			}
		}
	}`, matchQuery, fieldName, fieldName2, fieldName3)

	res, err := es.Search(
		es.Search.WithIndex(indexName),
		es.Search.WithBody(strings.NewReader(query)),
	)

	if err != nil {
		log.Fatalln("Error sending a post request to search a document to elasticsearch: ", err)
	}

	if res.IsError() {
		log.Fatalln("Error search a document in elasticsearch (elasticsearch response): ", res.Status())
	}

	var searchResult map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&searchResult)
	if err != nil {
		log.Fatalln("Error decode document: ", err)
	}

	hits := searchResult["hits"].(map[string]interface{})["hits"].([]interface{})

	for _, hit := range hits {
		document := hit.(map[string]interface{})
		source := document["_source"].(map[string]interface{})
		fmt.Printf("Document: %v\n", source)
	}

	log.Println("Success to searchDocumentWithMultiMatch")
}

func updateDocumentByDocumentID(es *elasticsearch.Client, indexName string, documentID int, updatedFields map[string]interface{}) {
	data, err := json.Marshal(map[string]interface{}{
		"doc": updatedFields,
	})

	if err != nil {
		log.Fatalln("Error marshalling update data: ", err)
	}

	res, err := es.Update(indexName,
		strconv.Itoa(documentID),
		bytes.NewReader(data),
	)

	if err != nil {
		log.Fatalln("Error sending a put/patch request to elasticsearch: ", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Fatalln("Error updating document in elasticsearch (elasticsearch response): ", res.Status())
	}

	log.Println("Success to update document")
}

func deleteDocumentByDocumentID(es *elasticsearch.Client, indexName string, documentID int) {
	res, err := es.Delete(indexName,
		strconv.Itoa(documentID),
	)

	if err != nil {
		log.Fatalln("Error sending a delete request to elasticsearch: ", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Fatalln("Error delete document in elasticsearch (elasticsearch response): ", res.Status())
	}

	log.Println("Success to delete document")
}

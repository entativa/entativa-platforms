package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Client wraps Elasticsearch client
type Client struct {
	es *elasticsearch.Client
}

// NewClient creates a new Elasticsearch client
func NewClient(addresses []string, username, password string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}

	if username != "" && password != "" {
		cfg.Username = username
		cfg.Password = password
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating elasticsearch client: %w", err)
	}

	// Test connection
	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("error getting elasticsearch info: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error response from elasticsearch: %s", res.String())
	}

	log.Println("Connected to Elasticsearch successfully")

	client := &Client{es: es}

	// Initialize indices
	if err := client.InitializeIndices(); err != nil {
		log.Printf("Warning: Failed to initialize indices: %v", err)
	}

	return client, nil
}

// GetClient returns the underlying Elasticsearch client
func (c *Client) GetClient() *elasticsearch.Client {
	return c.es
}

// IndexExists checks if an index exists
func (c *Client) IndexExists(ctx context.Context, indexName string) (bool, error) {
	res, err := c.es.Indices.Exists(
		[]string{indexName},
		c.es.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	return res.StatusCode == 200, nil
}

// CreateIndex creates an index with settings and mappings
func (c *Client) CreateIndex(ctx context.Context, indexName string, body map[string]interface{}) error {
	exists, err := c.IndexExists(ctx, indexName)
	if err != nil {
		return err
	}

	if exists {
		log.Printf("Index %s already exists, skipping creation", indexName)
		return nil
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshaling index body: %w", err)
	}

	res, err := c.es.Indices.Create(
		indexName,
		c.es.Indices.Create.WithContext(ctx),
		c.es.Indices.Create.WithBody(bytes.NewReader(bodyJSON)),
	)
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response creating index: %s", res.String())
	}

	log.Printf("Created index: %s", indexName)
	return nil
}

// DeleteIndex deletes an index
func (c *Client) DeleteIndex(ctx context.Context, indexName string) error {
	res, err := c.es.Indices.Delete(
		[]string{indexName},
		c.es.Indices.Delete.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting index: %s", res.String())
	}

	log.Printf("Deleted index: %s", indexName)
	return nil
}

// IndexDocument indexes a document
func (c *Client) IndexDocument(ctx context.Context, indexName, documentID string, document interface{}) error {
	docJSON, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("error marshaling document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: documentID,
		Body:       bytes.NewReader(docJSON),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("error indexing document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response indexing document: %s", res.String())
	}

	return nil
}

// UpdateDocument updates a document
func (c *Client) UpdateDocument(ctx context.Context, indexName, documentID string, doc interface{}) error {
	updateDoc := map[string]interface{}{
		"doc": doc,
	}

	docJSON, err := json.Marshal(updateDoc)
	if err != nil {
		return fmt.Errorf("error marshaling update document: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: documentID,
		Body:       bytes.NewReader(docJSON),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response updating document: %s", res.String())
	}

	return nil
}

// DeleteDocument deletes a document
func (c *Client) DeleteDocument(ctx context.Context, indexName, documentID string) error {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: documentID,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("error deleting document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error response deleting document: %s", res.String())
	}

	return nil
}

// BulkIndex performs bulk indexing
func (c *Client) BulkIndex(ctx context.Context, operations []map[string]interface{}) (int, int, error) {
	var buf bytes.Buffer
	for _, op := range operations {
		opJSON, err := json.Marshal(op)
		if err != nil {
			return 0, 0, fmt.Errorf("error marshaling operation: %w", err)
		}
		buf.Write(opJSON)
		buf.WriteByte('\n')
	}

	res, err := c.es.Bulk(
		bytes.NewReader(buf.Bytes()),
		c.es.Bulk.WithContext(ctx),
		c.es.Bulk.WithRefresh("true"),
	)
	if err != nil {
		return 0, 0, fmt.Errorf("error performing bulk operation: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, 0, fmt.Errorf("error response from bulk operation: %s", res.String())
	}

	// Parse response
	var bulkResp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&bulkResp); err != nil {
		return 0, 0, fmt.Errorf("error parsing bulk response: %w", err)
	}

	items := bulkResp["items"].([]interface{})
	successCount := 0
	failedCount := 0

	for _, item := range items {
		itemMap := item.(map[string]interface{})
		for _, action := range itemMap {
			actionMap := action.(map[string]interface{})
			if status, ok := actionMap["status"].(float64); ok {
				if status >= 200 && status < 300 {
					successCount++
				} else {
					failedCount++
				}
			}
		}
	}

	return successCount, failedCount, nil
}

// Search performs a search query
func (c *Client) Search(ctx context.Context, indices []string, query map[string]interface{}) (*esapi.Response, error) {
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("error marshaling query: %w", err)
	}

	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex(indices...),
		c.es.Search.WithBody(bytes.NewReader(queryJSON)),
		c.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("error performing search: %w", err)
	}

	if res.IsError() {
		defer res.Body.Close()
		return nil, fmt.Errorf("error response from search: %s", res.String())
	}

	return res, nil
}

// MultiSearch performs multiple searches in one request
func (c *Client) MultiSearch(ctx context.Context, queries []map[string]interface{}) (*esapi.Response, error) {
	var buf bytes.Buffer
	for _, query := range queries {
		// Empty header for each search
		buf.WriteString("{}\n")
		
		// Query body
		queryJSON, err := json.Marshal(query)
		if err != nil {
			return nil, fmt.Errorf("error marshaling query: %w", err)
		}
		buf.Write(queryJSON)
		buf.WriteByte('\n')
	}

	res, err := c.es.Msearch(
		bytes.NewReader(buf.Bytes()),
		c.es.Msearch.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("error performing multi-search: %w", err)
	}

	if res.IsError() {
		defer res.Body.Close()
		return nil, fmt.Errorf("error response from multi-search: %s", res.String())
	}

	return res, nil
}

// Suggest performs a suggest query (for autocomplete)
func (c *Client) Suggest(ctx context.Context, index string, suggesterName string, text string, field string) ([]string, error) {
	query := map[string]interface{}{
		"suggest": map[string]interface{}{
			suggesterName: map[string]interface{}{
				"prefix": text,
				"completion": map[string]interface{}{
					"field": field,
					"skip_duplicates": true,
					"size": 10,
				},
			},
		},
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("error marshaling suggest query: %w", err)
	}

	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex(index),
		c.es.Search.WithBody(bytes.NewReader(queryJSON)),
	)
	if err != nil {
		return nil, fmt.Errorf("error performing suggest: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error response from suggest: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing suggest response: %w", err)
	}

	// Extract suggestions
	var suggestions []string
	if suggest, ok := result["suggest"].(map[string]interface{}); ok {
		if suggester, ok := suggest[suggesterName].([]interface{}); ok {
			for _, item := range suggester {
				itemMap := item.(map[string]interface{})
				if options, ok := itemMap["options"].([]interface{}); ok {
					for _, opt := range options {
						optMap := opt.(map[string]interface{})
						if text, ok := optMap["text"].(string); ok {
							suggestions = append(suggestions, text)
						}
					}
				}
			}
		}
	}

	return suggestions, nil
}

// RefreshIndex refreshes an index
func (c *Client) RefreshIndex(ctx context.Context, indices ...string) error {
	res, err := c.es.Indices.Refresh(
		c.es.Indices.Refresh.WithContext(ctx),
		c.es.Indices.Refresh.WithIndex(indices...),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error refreshing index: %s", res.String())
	}

	return nil
}

// GetIndexStats gets statistics for an index
func (c *Client) GetIndexStats(ctx context.Context, indexName string) (map[string]interface{}, error) {
	res, err := c.es.Indices.Stats(
		c.es.Indices.Stats.WithContext(ctx),
		c.es.Indices.Stats.WithIndex(indexName),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error getting index stats: %s", res.String())
	}

	var stats map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("error parsing index stats: %w", err)
	}

	return stats, nil
}

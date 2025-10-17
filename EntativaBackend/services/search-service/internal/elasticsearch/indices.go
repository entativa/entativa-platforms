package elasticsearch

import (
	"context"
	"log"
)

// Index names
const (
	IndexUsers     = "users"
	IndexPosts     = "posts"
	IndexTakes     = "takes"
	IndexHashtags  = "hashtags"
	IndexLocations = "locations"
)

// InitializeIndices creates all required indices with proper mappings
func (c *Client) InitializeIndices() error {
	ctx := context.Background()

	indices := map[string]map[string]interface{}{
		IndexUsers:     getUsersIndexMapping(),
		IndexPosts:     getPostsIndexMapping(),
		IndexTakes:     getTakesIndexMapping(),
		IndexHashtags:  getHashtagsIndexMapping(),
		IndexLocations: getLocationsIndexMapping(),
	}

	for indexName, mapping := range indices {
		if err := c.CreateIndex(ctx, indexName, mapping); err != nil {
			log.Printf("Error creating index %s: %v", indexName, err)
			continue
		}
	}

	return nil
}

// getUsersIndexMapping returns the mapping for users index
func getUsersIndexMapping() map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   3,
			"number_of_replicas": 1,
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"username_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "keyword",
						"filter": []string{
							"lowercase",
						},
					},
					"text_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"asciifolding",
						},
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "keyword",
				},
				"username": map[string]interface{}{
					"type":     "text",
					"analyzer": "username_analyzer",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
						"suggest": map[string]interface{}{
							"type": "completion",
						},
					},
				},
				"display_name": map[string]interface{}{
					"type":     "text",
					"analyzer": "text_analyzer",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"bio": map[string]interface{}{
					"type":     "text",
					"analyzer": "text_analyzer",
				},
				"location": map[string]interface{}{
					"type": "text",
				},
				"verified": map[string]interface{}{
					"type": "boolean",
				},
				"avatar_url": map[string]interface{}{
					"type":  "keyword",
					"index": false,
				},
				"follower_count": map[string]interface{}{
					"type": "integer",
				},
				"following_count": map[string]interface{}{
					"type": "integer",
				},
				"post_count": map[string]interface{}{
					"type": "integer",
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
			},
		},
	}
}

// getPostsIndexMapping returns the mapping for posts index
func getPostsIndexMapping() map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   5,
			"number_of_replicas": 1,
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"text_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"asciifolding",
							"english_stop",
						},
					},
				},
				"filter": map[string]interface{}{
					"english_stop": map[string]interface{}{
						"type":      "stop",
						"stopwords": "_english_",
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "keyword",
				},
				"user_id": map[string]interface{}{
					"type": "keyword",
				},
				"username": map[string]interface{}{
					"type": "keyword",
				},
				"caption": map[string]interface{}{
					"type":     "text",
					"analyzer": "text_analyzer",
				},
				"content": map[string]interface{}{
					"type":     "text",
					"analyzer": "text_analyzer",
				},
				"media_ids": map[string]interface{}{
					"type": "keyword",
				},
				"media_type": map[string]interface{}{
					"type": "keyword",
				},
				"has_media": map[string]interface{}{
					"type": "boolean",
				},
				"hashtags": map[string]interface{}{
					"type": "keyword",
				},
				"location": map[string]interface{}{
					"type": "text",
				},
				"likes_count": map[string]interface{}{
					"type": "integer",
				},
				"comments_count": map[string]interface{}{
					"type": "integer",
				},
				"shares_count": map[string]interface{}{
					"type": "integer",
				},
				"views_count": map[string]interface{}{
					"type": "integer",
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
			},
		},
	}
}

// getTakesIndexMapping returns the mapping for takes index
func getTakesIndexMapping() map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   5,
			"number_of_replicas": 1,
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"text_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"asciifolding",
						},
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "keyword",
				},
				"user_id": map[string]interface{}{
					"type": "keyword",
				},
				"username": map[string]interface{}{
					"type": "keyword",
				},
				"caption": map[string]interface{}{
					"type":     "text",
					"analyzer": "text_analyzer",
				},
				"media_id": map[string]interface{}{
					"type": "keyword",
				},
				"thumbnail_url": map[string]interface{}{
					"type":  "keyword",
					"index": false,
				},
				"duration": map[string]interface{}{
					"type": "float",
				},
				"hashtags": map[string]interface{}{
					"type": "keyword",
				},
				"audio_track_id": map[string]interface{}{
					"type": "keyword",
				},
				"filter_used": map[string]interface{}{
					"type": "keyword",
				},
				"trend_id": map[string]interface{}{
					"type": "keyword",
				},
				"template_id": map[string]interface{}{
					"type": "keyword",
				},
				"views_count": map[string]interface{}{
					"type": "integer",
				},
				"likes_count": map[string]interface{}{
					"type": "integer",
				},
				"comments_count": map[string]interface{}{
					"type": "integer",
				},
				"remix_count": map[string]interface{}{
					"type": "integer",
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
			},
		},
	}
}

// getHashtagsIndexMapping returns the mapping for hashtags index
func getHashtagsIndexMapping() map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   3,
			"number_of_replicas": 1,
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"hashtag_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "keyword",
						"filter": []string{
							"lowercase",
						},
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"tag": map[string]interface{}{
					"type":     "text",
					"analyzer": "hashtag_analyzer",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
						"suggest": map[string]interface{}{
							"type": "completion",
						},
					},
				},
				"display_tag": map[string]interface{}{
					"type": "keyword",
				},
				"usage_count": map[string]interface{}{
					"type": "long",
				},
				"post_count": map[string]interface{}{
					"type": "long",
				},
				"take_count": map[string]interface{}{
					"type": "long",
				},
				"growth_rate": map[string]interface{}{
					"type": "float",
				},
				"is_trending": map[string]interface{}{
					"type": "boolean",
				},
				"category": map[string]interface{}{
					"type": "keyword",
				},
				"first_used": map[string]interface{}{
					"type": "date",
				},
				"last_used": map[string]interface{}{
					"type": "date",
				},
			},
		},
	}
}

// getLocationsIndexMapping returns the mapping for locations index
func getLocationsIndexMapping() map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   3,
			"number_of_replicas": 1,
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "keyword",
				},
				"name": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
						"suggest": map[string]interface{}{
							"type": "completion",
						},
					},
				},
				"city": map[string]interface{}{
					"type": "text",
				},
				"country": map[string]interface{}{
					"type": "keyword",
				},
				"coordinates": map[string]interface{}{
					"type": "geo_point",
				},
				"latitude": map[string]interface{}{
					"type": "float",
				},
				"longitude": map[string]interface{}{
					"type": "float",
				},
				"post_count": map[string]interface{}{
					"type": "long",
				},
				"take_count": map[string]interface{}{
					"type": "long",
				},
				"checkin_count": map[string]interface{}{
					"type": "long",
				},
				"category": map[string]interface{}{
					"type": "keyword",
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
			},
		},
	}
}

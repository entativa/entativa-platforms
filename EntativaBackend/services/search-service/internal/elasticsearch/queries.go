package elasticsearch

import (
	"socialink/search-service/internal/model"
)

// BuildSearchQuery builds an Elasticsearch query for search
func BuildSearchQuery(req *model.SearchRequest) map[string]interface{} {
	query := map[string]interface{}{
		"from": req.Offset,
		"size": req.Limit,
		"query": buildMultiMatchQuery(req.Query),
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"*": map[string]interface{}{},
			},
			"pre_tags":  []string{"<mark>"},
			"post_tags": []string{"</mark>"},
		},
	}

	// Add filters
	if len(getFilters(req)) > 0 {
		query["query"] = map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   buildMultiMatchQuery(req.Query),
				"filter": getFilters(req),
			},
		}
	}

	// Add sorting
	query["sort"] = []map[string]interface{}{
		{
			"_score": map[string]interface{}{
				"order": "desc",
			},
		},
		{
			"created_at": map[string]interface{}{
				"order": "desc",
			},
		},
	}

	return query
}

// buildMultiMatchQuery builds a multi-match query
func buildMultiMatchQuery(query string) map[string]interface{} {
	return map[string]interface{}{
		"multi_match": map[string]interface{}{
			"query": query,
			"fields": []string{
				"username^3",        // Username gets highest boost
				"display_name^2",    // Display name gets medium boost
				"caption^2",
				"content^2",
				"bio",
				"tag^3",            // Hashtag gets high boost
				"name^2",           // Location name gets medium boost
			},
			"type":                "best_fields",
			"tie_breaker":         0.3,
			"minimum_should_match": "75%",
			"fuzziness":           "AUTO",
		},
	}
}

// getFilters builds filters based on search filters
func getFilters(req *model.SearchRequest) []map[string]interface{} {
	var filters []map[string]interface{}

	// User filters
	if req.Filters.Verified != nil {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"verified": *req.Filters.Verified,
			},
		})
	}

	if req.Filters.MinFollowers > 0 {
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{
				"follower_count": map[string]interface{}{
					"gte": req.Filters.MinFollowers,
				},
			},
		})
	}

	// Post/Take filters
	if req.Filters.HasMedia != nil {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"has_media": *req.Filters.HasMedia,
			},
		})
	}

	if req.Filters.MediaType != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"media_type": req.Filters.MediaType,
			},
		})
	}

	if !req.Filters.DateFrom.IsZero() || !req.Filters.DateTo.IsZero() {
		rangeFilter := map[string]interface{}{}
		if !req.Filters.DateFrom.IsZero() {
			rangeFilter["gte"] = req.Filters.DateFrom.Format("2006-01-02")
		}
		if !req.Filters.DateTo.IsZero() {
			rangeFilter["lte"] = req.Filters.DateTo.Format("2006-01-02")
		}
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{
				"created_at": rangeFilter,
			},
		})
	}

	if req.Filters.MinLikes > 0 {
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{
				"likes_count": map[string]interface{}{
					"gte": req.Filters.MinLikes,
				},
			},
		})
	}

	if req.Filters.MinViews > 0 {
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{
				"views_count": map[string]interface{}{
					"gte": req.Filters.MinViews,
				},
			},
		})
	}

	// Hashtag filters
	if req.Filters.TrendingOnly {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"is_trending": true,
			},
		})
	}

	// Location filters (geo-distance)
	if req.Filters.Latitude != 0 && req.Filters.Longitude != 0 && req.Filters.Radius > 0 {
		filters = append(filters, map[string]interface{}{
			"geo_distance": map[string]interface{}{
				"distance": map[string]interface{}{
					"coordinates": map[string]interface{}{
						"lat": req.Filters.Latitude,
						"lon": req.Filters.Longitude,
					},
				},
				"coordinates": map[string]interface{}{
					"lat": req.Filters.Latitude,
					"lon": req.Filters.Longitude,
				},
			},
		})
	}

	return filters
}

// BuildAutocompleteQuery builds a completion suggester query
func BuildAutocompleteQuery(text string, field string) map[string]interface{} {
	return map[string]interface{}{
		"suggest": map[string]interface{}{
			"autocomplete": map[string]interface{}{
				"prefix": text,
				"completion": map[string]interface{}{
					"field":           field,
					"skip_duplicates": true,
					"size":            10,
					"fuzzy": map[string]interface{}{
						"fuzziness": "AUTO",
					},
				},
			},
		},
	}
}

// BuildTrendingHashtagsQuery builds a query for trending hashtags
func BuildTrendingHashtagsQuery(limit int) map[string]interface{} {
	return map[string]interface{}{
		"size": limit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"is_trending": true,
						},
					},
				},
			},
		},
		"sort": []map[string]interface{}{
			{
				"growth_rate": map[string]interface{}{
					"order": "desc",
				},
			},
			{
				"usage_count": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}
}

// BuildRelatedHashtagsQuery builds a query for related hashtags
func BuildRelatedHashtagsQuery(tag string, limit int) map[string]interface{} {
	return map[string]interface{}{
		"size": limit,
		"query": map[string]interface{}{
			"more_like_this": map[string]interface{}{
				"fields": []string{"tag"},
				"like": []map[string]interface{}{
					{
						"_index": IndexHashtags,
						"_id":    tag,
					},
				},
				"min_term_freq":     1,
				"max_query_terms":   12,
				"min_doc_freq":      1,
				"minimum_should_match": "30%",
			},
		},
		"sort": []map[string]interface{}{
			{
				"usage_count": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}
}

// BuildAggregationQuery builds an aggregation query
func BuildAggregationQuery(field string, size int) map[string]interface{} {
	return map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"top_items": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": field,
					"size":  size,
					"order": map[string]interface{}{
						"_count": "desc",
					},
				},
			},
		},
	}
}

// BuildRangeQuery builds a range query for numeric fields
func BuildRangeQuery(field string, gte, lte interface{}) map[string]interface{} {
	rangeMap := make(map[string]interface{})
	if gte != nil {
		rangeMap["gte"] = gte
	}
	if lte != nil {
		rangeMap["lte"] = lte
	}

	return map[string]interface{}{
		"range": map[string]interface{}{
			field: rangeMap,
		},
	}
}

// BuildTermQuery builds a term query for exact matching
func BuildTermQuery(field string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"term": map[string]interface{}{
			field: value,
		},
	}
}

// BuildMatchPhraseQuery builds a match phrase query
func BuildMatchPhraseQuery(field, query string) map[string]interface{} {
	return map[string]interface{}{
		"match_phrase": map[string]interface{}{
			field: map[string]interface{}{
				"query": query,
				"slop":  2,
			},
		},
	}
}

// BuildFuzzyQuery builds a fuzzy query for typo tolerance
func BuildFuzzyQuery(field, query string, fuzziness string) map[string]interface{} {
	return map[string]interface{}{
		"fuzzy": map[string]interface{}{
			field: map[string]interface{}{
				"value":      query,
				"fuzziness":  fuzziness,
				"prefix_length": 1,
			},
		},
	}
}

// BuildGeoDistanceQuery builds a geo-distance query
func BuildGeoDistanceQuery(lat, lon, distance float64, unit string) map[string]interface{} {
	return map[string]interface{}{
		"geo_distance": map[string]interface{}{
			"distance": fmt.Sprintf("%f%s", distance, unit),
			"coordinates": map[string]interface{}{
				"lat": lat,
				"lon": lon,
			},
		},
	}
}

// BuildBoolQuery builds a bool query combining multiple queries
func BuildBoolQuery(must, should, mustNot, filter []map[string]interface{}) map[string]interface{} {
	boolQuery := map[string]interface{}{}

	if len(must) > 0 {
		boolQuery["must"] = must
	}
	if len(should) > 0 {
		boolQuery["should"] = should
	}
	if len(mustNot) > 0 {
		boolQuery["must_not"] = mustNot
	}
	if len(filter) > 0 {
		boolQuery["filter"] = filter
	}

	return map[string]interface{}{
		"bool": boolQuery,
	}
}

// BuildFunctionScoreQuery builds a function score query for custom ranking
func BuildFunctionScoreQuery(query map[string]interface{}, functions []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"function_score": map[string]interface{}{
			"query":      query,
			"functions":  functions,
			"score_mode": "sum",
			"boost_mode": "multiply",
		},
	}
}

// BuildDecayFunction builds a decay function for scoring
func BuildDecayFunction(field string, origin interface{}, scale, offset, decay string) map[string]interface{} {
	return map[string]interface{}{
		"gauss": map[string]interface{}{
			field: map[string]interface{}{
				"origin": origin,
				"scale":  scale,
				"offset": offset,
				"decay":  decay,
			},
		},
	}
}

// BuildScriptScore builds a script score function
func BuildScriptScore(source string, params map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"script_score": map[string]interface{}{
			"script": map[string]interface{}{
				"source": source,
				"params": params,
			},
		},
	}
}

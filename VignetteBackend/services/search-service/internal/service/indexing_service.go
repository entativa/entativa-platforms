package service

import (
	"context"
	"fmt"
	"log"

	"vignette/search-service/internal/elasticsearch"
	"vignette/search-service/internal/model"

	"github.com/go-redis/redis/v8"
)

type IndexingService struct {
	es    *elasticsearch.Client
	redis *redis.Client
}

func NewIndexingService(es *elasticsearch.Client, redis *redis.Client) *IndexingService {
	return &IndexingService{
		es:    es,
		redis: redis,
	}
}

// IndexDocument indexes a single document
func (s *IndexingService) IndexDocument(ctx context.Context, req *model.IndexingRequest) (*model.IndexingResponse, error) {
	// Get index name
	indexName := s.getIndexName(req.DocumentType)
	
	switch req.Action {
	case "index":
		return s.indexDoc(ctx, indexName, req.DocumentID, req.Data)
	case "update":
		return s.updateDoc(ctx, indexName, req.DocumentID, req.Data)
	case "delete":
		return s.deleteDoc(ctx, indexName, req.DocumentID)
	default:
		return nil, fmt.Errorf("invalid action: %s", req.Action)
	}
}

// BulkIndex indexes multiple documents in a single request
func (s *IndexingService) BulkIndex(ctx context.Context, req *model.BulkIndexingRequest) (*model.BulkIndexingResponse, error) {
	if len(req.Documents) == 0 {
		return &model.BulkIndexingResponse{
			Success:     true,
			TotalDocs:   0,
			IndexedDocs: 0,
			FailedDocs:  0,
		}, nil
	}

	// Build bulk operations
	var operations []map[string]interface{}
	for _, doc := range req.Documents {
		indexName := s.getIndexName(doc.DocumentType)
		
		// Action metadata
		actionMeta := map[string]interface{}{
			"_index": indexName,
			"_id":    doc.DocumentID,
		}

		switch doc.Action {
		case "index":
			operations = append(operations, map[string]interface{}{
				"index": actionMeta,
			})
			operations = append(operations, doc.Data.(map[string]interface{}))
			
		case "update":
			operations = append(operations, map[string]interface{}{
				"update": actionMeta,
			})
			operations = append(operations, map[string]interface{}{
				"doc": doc.Data,
			})
			
		case "delete":
			operations = append(operations, map[string]interface{}{
				"delete": actionMeta,
			})
		}
	}

	// Execute bulk operation
	successCount, failedCount, err := s.es.BulkIndex(ctx, operations)
	if err != nil {
		return nil, fmt.Errorf("bulk index error: %w", err)
	}

	response := &model.BulkIndexingResponse{
		Success:     failedCount == 0,
		TotalDocs:   len(req.Documents),
		IndexedDocs: successCount,
		FailedDocs:  failedCount,
	}

	// Invalidate cache
	s.invalidateCache(ctx)

	log.Printf("Bulk indexed %d documents (%d successful, %d failed)", 
		len(req.Documents), successCount, failedCount)

	return response, nil
}

// indexDoc indexes a document
func (s *IndexingService) indexDoc(ctx context.Context, indexName, documentID string, data interface{}) (*model.IndexingResponse, error) {
	if err := s.es.IndexDocument(ctx, indexName, documentID, data); err != nil {
		return &model.IndexingResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	// Invalidate cache
	s.invalidateCache(ctx)

	return &model.IndexingResponse{
		Success:    true,
		Message:    "Document indexed successfully",
		DocumentID: documentID,
	}, nil
}

// updateDoc updates a document
func (s *IndexingService) updateDoc(ctx context.Context, indexName, documentID string, data interface{}) (*model.IndexingResponse, error) {
	if err := s.es.UpdateDocument(ctx, indexName, documentID, data); err != nil {
		return &model.IndexingResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	// Invalidate cache
	s.invalidateCache(ctx)

	return &model.IndexingResponse{
		Success:    true,
		Message:    "Document updated successfully",
		DocumentID: documentID,
	}, nil
}

// deleteDoc deletes a document
func (s *IndexingService) deleteDoc(ctx context.Context, indexName, documentID string) (*model.IndexingResponse, error) {
	if err := s.es.DeleteDocument(ctx, indexName, documentID); err != nil {
		return &model.IndexingResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	// Invalidate cache
	s.invalidateCache(ctx)

	return &model.IndexingResponse{
		Success:    true,
		Message:    "Document deleted successfully",
		DocumentID: documentID,
	}, nil
}

// getIndexName returns the index name for a document type
func (s *IndexingService) getIndexName(docType model.SearchType) string {
	switch docType {
	case model.SearchTypeUser:
		return elasticsearch.IndexUsers
	case model.SearchTypePost:
		return elasticsearch.IndexPosts
	case model.SearchTypeTake:
		return elasticsearch.IndexTakes
	case model.SearchTypeHashtag:
		return elasticsearch.IndexHashtags
	case model.SearchTypeLocation:
		return elasticsearch.IndexLocations
	default:
		return ""
	}
}

// invalidateCache invalidates search cache
func (s *IndexingService) invalidateCache(ctx context.Context) {
	// Delete all search cache keys
	pattern := "search:*"
	
	iter := s.redis.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		if err := s.redis.Del(ctx, iter.Val()).Err(); err != nil {
			log.Printf("Error deleting cache key %s: %v", iter.Val(), err)
		}
	}

	if err := iter.Err(); err != nil {
		log.Printf("Error iterating cache keys: %v", err)
	}
}

// ReindexAll reindexes all documents (admin operation)
func (s *IndexingService) ReindexAll(ctx context.Context) error {
	log.Println("Starting full reindex...")

	// This would typically fetch all documents from the main database
	// and index them. For now, we'll just refresh all indices.

	indices := []string{
		elasticsearch.IndexUsers,
		elasticsearch.IndexPosts,
		elasticsearch.IndexTakes,
		elasticsearch.IndexHashtags,
		elasticsearch.IndexLocations,
	}

	for _, index := range indices {
		if err := s.es.RefreshIndex(ctx, index); err != nil {
			log.Printf("Error refreshing index %s: %v", index, err)
			continue
		}
		log.Printf("Refreshed index: %s", index)
	}

	// Invalidate all caches
	s.invalidateCache(ctx)

	log.Println("Full reindex completed")
	return nil
}

// GetIndexStats returns statistics for all indices
func (s *IndexingService) GetIndexStats(ctx context.Context) (map[string]interface{}, error) {
	indices := []string{
		elasticsearch.IndexUsers,
		elasticsearch.IndexPosts,
		elasticsearch.IndexTakes,
		elasticsearch.IndexHashtags,
		elasticsearch.IndexLocations,
	}

	stats := make(map[string]interface{})

	for _, index := range indices {
		indexStats, err := s.es.GetIndexStats(ctx, index)
		if err != nil {
			log.Printf("Error getting stats for index %s: %v", index, err)
			stats[index] = map[string]interface{}{
				"error": err.Error(),
			}
			continue
		}

		// Extract relevant stats
		if indices, ok := indexStats["indices"].(map[string]interface{}); ok {
			if indexData, ok := indices[index].(map[string]interface{}); ok {
				if primaries, ok := indexData["primaries"].(map[string]interface{}); ok {
					if docs, ok := primaries["docs"].(map[string]interface{}); ok {
						stats[index] = map[string]interface{}{
							"document_count": docs["count"],
							"deleted_count":  docs["deleted"],
						}
					}
				}
			}
		}
	}

	return stats, nil
}

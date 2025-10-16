package handler

import (
	"net/http"

	"vignette/search-service/internal/model"
	"vignette/search-service/internal/service"

	"github.com/gin-gonic/gin"
)

type IndexingHandler struct {
	indexingService *service.IndexingService
}

func NewIndexingHandler(indexingService *service.IndexingService) *IndexingHandler {
	return &IndexingHandler{
		indexingService: indexingService,
	}
}

// IndexDocument indexes a single document
// @Summary Index document
// @Description Index a single document in Elasticsearch
// @Tags indexing
// @Accept json
// @Produce json
// @Param request body model.IndexingRequest true "Indexing request"
// @Success 200 {object} model.IndexingResponse
// @Router /index/document [post]
func (h *IndexingHandler) IndexDocument(c *gin.Context) {
	var req model.IndexingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.indexingService.IndexDocument(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// BulkIndex indexes multiple documents
// @Summary Bulk index documents
// @Description Index multiple documents in a single request
// @Tags indexing
// @Accept json
// @Produce json
// @Param request body model.BulkIndexingRequest true "Bulk indexing request"
// @Success 200 {object} model.BulkIndexingResponse
// @Router /index/bulk [post]
func (h *IndexingHandler) BulkIndex(c *gin.Context) {
	var req model.BulkIndexingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.indexingService.BulkIndex(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateDocument updates a document
// @Summary Update document
// @Description Update a document in Elasticsearch
// @Tags indexing
// @Accept json
// @Produce json
// @Param request body model.IndexingRequest true "Update request"
// @Success 200 {object} model.IndexingResponse
// @Router /index/document [put]
func (h *IndexingHandler) UpdateDocument(c *gin.Context) {
	var req model.IndexingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Action = "update"

	result, err := h.indexingService.IndexDocument(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteDocument deletes a document
// @Summary Delete document
// @Description Delete a document from Elasticsearch
// @Tags indexing
// @Accept json
// @Produce json
// @Param document_type query string true "Document type"
// @Param document_id query string true "Document ID"
// @Success 200 {object} model.IndexingResponse
// @Router /index/document [delete]
func (h *IndexingHandler) DeleteDocument(c *gin.Context) {
	documentType := model.SearchType(c.Query("document_type"))
	documentID := c.Query("document_id")

	if documentType == "" || documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document_type and document_id are required"})
		return
	}

	req := &model.IndexingRequest{
		Action:       "delete",
		DocumentType: documentType,
		DocumentID:   documentID,
	}

	result, err := h.indexingService.IndexDocument(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ReindexAll reindexes all documents (admin operation)
// @Summary Reindex all documents
// @Description Reindex all documents across all indices
// @Tags indexing
// @Success 200 {object} map[string]interface{}
// @Router /index/reindex [post]
func (h *IndexingHandler) ReindexAll(c *gin.Context) {
	// This should be protected by admin auth in production
	err := h.indexingService.ReindexAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Reindex completed successfully",
	})
}

// GetIndexStats gets statistics for all indices
// @Summary Get index statistics
// @Description Get document counts and stats for all indices
// @Tags indexing
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /index/stats [get]
func (h *IndexingHandler) GetIndexStats(c *gin.Context) {
	stats, err := h.indexingService.GetIndexStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

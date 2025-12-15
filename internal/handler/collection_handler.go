package handler

import (
	"net/http"
	"strconv"

	"Alpha_Strike_Helper/internal/domain"
	"Alpha_Strike_Helper/internal/service"

	"github.com/gin-gonic/gin"
)

type CollectionHandler struct {
	service *service.CollectionService
}

func NewCollectionHandler(service *service.CollectionService) *CollectionHandler {
	return &CollectionHandler{service: service}
}

// List получает список коллекций пользователя
func (h *CollectionHandler) List(c *gin.Context) {
	// В реальном приложении нужно получить userID из контекста (JWT токен)
	userID := uint(1) // Placeholder

	collections, err := h.service.ListCollections(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": collections})
}

// Get получает коллекцию по ID
func (h *CollectionHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	collection, err := h.service.GetCollection(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": collection})
}

// Create создаёт новую коллекцию
func (h *CollectionHandler) Create(c *gin.Context) {
	var collection domain.Collection
	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// В реальном приложении нужно получить userID из контекста
	collection.UserID = 1 // Placeholder

	if err := h.service.CreateCollection(&collection); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": collection})
}

// Update обновляет коллекцию
func (h *CollectionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var collection domain.Collection
	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection.ID = uint(id)
	if err := h.service.UpdateCollection(&collection); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": collection})
}

// Delete удаляет коллекцию
func (h *CollectionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteCollection(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddCard добавляет карточку в коллекцию
func (h *CollectionHandler) AddCard(c *gin.Context) {
	collectionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	cardID, err := strconv.ParseUint(c.Param("cardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}

	if err := h.service.AddCardToCollection(uint(collectionID), uint(cardID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card added to collection"})
}

// RemoveCard удаляет карточку из коллекции
func (h *CollectionHandler) RemoveCard(c *gin.Context) {
	collectionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	cardID, err := strconv.ParseUint(c.Param("cardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}

	if err := h.service.RemoveCardFromCollection(uint(collectionID), uint(cardID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

package handler

import (
	"Alpha_Strike_Helper/internal/domain"
	"Alpha_Strike_Helper/internal/service"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	service *service.CardService
}

func NewCardHandler(service *service.CardService) *CardHandler {
	return &CardHandler{service: service}
}

// List получает список карточек с пагинацией и фильтрами
// List получает список карточек с пагинацией и фильтрами
func (h *CardHandler) List(c *gin.Context) {
	// Параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 3604 // Разумное значение по умолчанию

	// Поддержка параметров pagesize и page_size
	if ps := c.Query("pagesize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	} else if ps := c.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	if page < 1 {
		page = 1
	}

	// Фильтры
	filters := make(map[string]interface{})
	if role := c.Query("role"); role != "" {
		filters["role"] = role
	}
	if size := c.Query("size"); size != "" {
		filters["size"] = size
	}
	if faction := c.Query("faction"); faction != "" {
		filters["faction"] = faction
	}
	if cardType := c.Query("type"); cardType != "" {
		filters["type"] = cardType
	}
	if techBase := c.Query("techbase"); techBase != "" {
		filters["tech_base"] = techBase
	}
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if pvMin := c.Query("pvmin"); pvMin != "" {
		filters["pv_min"] = pvMin
	}
	if pvMax := c.Query("pvmax"); pvMax != "" {
		filters["pv_max"] = pvMax
	}

	// ✅ КЛЮЧЕВАЯ СТРОКА: cards, total, err
	cards, total, err := h.service.ListCards(filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if cards == nil {
		cards = []domain.Card{}
	}

	// ✅ РАСЧЕТ total_pages (исправление основной проблемы)
	totalPages := int64(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        cards,
		"page":        page,
		"page_size":   pageSize,
		"total":       total,
		"total_pages": totalPages, // ✅ ИСПРАВЛЕНО: было "1"
	})
}

// Get получает карточку по ID
func (h *CardHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	card, err := h.service.GetCard(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": card})
}

// Create создаёт новую карточку (админ только)
func (h *CardHandler) Create(c *gin.Context) {
	var card domain.Card
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateCard(&card); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": card})
}

// Search ищет карточки по названию и номеру модели
func (h *CardHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	cards, err := h.service.SearchCards(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if cards == nil {
		cards = []domain.Card{}
	}

	c.JSON(http.StatusOK, gin.H{"data": cards})
}

// Update обновляет карточку (админ только)
func (h *CardHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var card domain.Card
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card.ID = uint(id)

	if err := h.service.UpdateCard(&card); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": card})
}

// Delete удаляет карточку (админ только)
func (h *CardHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteCard(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

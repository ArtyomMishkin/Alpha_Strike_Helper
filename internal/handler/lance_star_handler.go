package handler

import (
	"Alpha_Strike_Helper/internal/domain"
	"Alpha_Strike_Helper/internal/repository"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type LanceHandler struct {
	lanceRepo repository.LanceRepository
}

type StarHandler struct {
	starRepo repository.StarRepository
}

func NewLanceHandler(lanceRepo repository.LanceRepository) *LanceHandler {
	return &LanceHandler{lanceRepo: lanceRepo}
}

func NewStarHandler(starRepo repository.StarRepository) *StarHandler {
	return &StarHandler{starRepo: starRepo}
}

// ========== LANCE HANDLERS ==========

// CreateLance POST /api/v1/lances
func (h *LanceHandler) CreateLance(c *gin.Context) {
	userID := uint(1)
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Faction     string `json:"faction"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lance := &domain.Lance{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Faction:     req.Faction,
	}

	log.Printf("Creating lance: %+v", lance)
	if err := h.lanceRepo.Create(lance); err != nil {
		log.Printf("Failed to create lance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, lance)
}

// GetLances GET /api/v1/lances
func (h *LanceHandler) GetLances(c *gin.Context) {
	userID := uint(1)
	lances, err := h.lanceRepo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get lances"})
		return
	}

	if lances == nil {
		lances = []*domain.Lance{}
	}

	c.JSON(http.StatusOK, lances)
}

// GetLance GET /api/v1/lances/:id
// ✅ ИСПРАВЛЕНО: Теперь загружает Members.Card
func (h *LanceHandler) GetLance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lance id"})
		return
	}

	lance, err := h.lanceRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lance not found"})
		return
	}

	// Убедимся, что Members не nil (пустой массив вместо null в JSON)
	if lance.Members == nil {
		lance.Members = []domain.LanceMember{}
	}

	c.JSON(http.StatusOK, gin.H{
		"lance": lance,
	})
}

// UpdateLance PUT /api/v1/lances/:id
func (h *LanceHandler) UpdateLance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lance id"})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Faction     string `json:"faction"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lance, err := h.lanceRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lance not found"})
		return
	}

	lance.Name = req.Name
	lance.Description = req.Description
	if req.Faction != "" {
		lance.Faction = req.Faction
	}

	if err := h.lanceRepo.Update(lance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update lance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lance": lance})
}

// DeleteLance DELETE /api/v1/lances/:id
func (h *LanceHandler) DeleteLance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lance id"})
		return
	}

	if err := h.lanceRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete lance"})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddCardToLance POST /api/v1/lances/:id/cards/:cardId
func (h *LanceHandler) AddCardToLance(c *gin.Context) {
	lanceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lance id"})
		return
	}

	cardID, err := strconv.ParseUint(c.Param("cardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card id"})
		return
	}

	var req struct {
		Position int `json:"position" binding:"required,min=1,max=4"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lance, err := h.lanceRepo.GetByID(uint(lanceID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lance not found"})
		return
	}

	if len(lance.Members) >= 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Lance не может содержать более 4 мехов, текущее количество: %d", len(lance.Members))})
		return
	}

	member := &domain.LanceMember{
		LanceID:  uint(lanceID),
		CardID:   uint(cardID),
		Position: req.Position,
	}

	if err := h.lanceRepo.AddMember(member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add card to lance"})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// RemoveCardFromLance DELETE /api/v1/lances/:id/cards/:cardId
func (h *LanceHandler) RemoveCardFromLance(c *gin.Context) {
	lanceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lance id"})
		return
	}

	cardID, err := strconv.ParseUint(c.Param("cardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card id"})
		return
	}

	if err := h.lanceRepo.RemoveMember(uint(lanceID), uint(cardID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove card from lance"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ExportLance GET /api/v1/lances/:id/export
// ExportLance GET /api/v1/lances/:id/export
func (h *LanceHandler) ExportLance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lance id"})
		return
	}

	lance, err := h.lanceRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lance not found"})
		return
	}

	if lance.Members == nil {
		lance.Members = []domain.LanceMember{}
	}

	// ✅ Считаем статистику
	totalPV := 0
	totalArmor := 0
	for _, member := range lance.Members {
		if member.Card != nil {
			totalPV += member.Card.PointValue
			totalArmor += member.Card.Armor
		}
	}

	// Подготовка данных для экспорта
	exportData := gin.H{
		"lance_name":   lance.Name,
		"lance_id":     lance.ID,
		"faction":      lance.Faction,
		"members":      lance.Members,
		"total_pv":     totalPV,
		"total_armor":  totalArmor,
		"member_count": len(lance.Members),
		"timestamp":    time.Now().Format("2006-01-02 15:04:05"),
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, exportData)
}

// ========== STAR HANDLERS ==========

// CreateStar POST /api/v1/stars
func (h *StarHandler) CreateStar(c *gin.Context) {
	userID := uint(1)
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Faction     string `json:"faction"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	star := &domain.Star{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Faction:     req.Faction,
	}

	if err := h.starRepo.Create(star); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, star)
}

// GetStars GET /api/v1/stars
func (h *StarHandler) GetStars(c *gin.Context) {
	userID := uint(1)
	stars, err := h.starRepo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stars"})
		return
	}

	if stars == nil {
		stars = []*domain.Star{}
	}

	c.JSON(http.StatusOK, stars)
}

// GetStar GET /api/v1/stars/:id
func (h *StarHandler) GetStar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid star id"})
		return
	}

	star, err := h.starRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "star not found"})
		return
	}

	if star.Members == nil {
		star.Members = []domain.StarMember{}
	}

	c.JSON(http.StatusOK, gin.H{
		"star": star,
	})
}

// UpdateStar PUT /api/v1/stars/:id
func (h *StarHandler) UpdateStar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid star id"})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Faction     string `json:"faction"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	star, err := h.starRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "star not found"})
		return
	}

	star.Name = req.Name
	star.Description = req.Description
	if req.Faction != "" {
		star.Faction = req.Faction
	}

	if err := h.starRepo.Update(star); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update star"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"star": star})
}

// DeleteStar DELETE /api/v1/stars/:id
func (h *StarHandler) DeleteStar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid star id"})
		return
	}

	if err := h.starRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete star"})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddCardToStar POST /api/v1/stars/:id/cards/:cardId
func (h *StarHandler) AddCardToStar(c *gin.Context) {
	starID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid star id"})
		return
	}

	cardID, err := strconv.ParseUint(c.Param("cardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card id"})
		return
	}

	var req struct {
		Position int `json:"position" binding:"required,min=1,max=5"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	star, err := h.starRepo.GetByID(uint(starID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "star not found"})
		return
	}

	if len(star.Members) >= 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Star не может содержать более 5 мехов, текущее количество: %d", len(star.Members))})
		return
	}

	member := &domain.StarMember{
		StarID:   uint(starID),
		CardID:   uint(cardID),
		Position: req.Position,
	}

	if err := h.starRepo.AddMember(member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add card to star"})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// RemoveCardFromStar DELETE /api/v1/stars/:id/cards/:cardId
func (h *StarHandler) RemoveCardFromStar(c *gin.Context) {
	starID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid star id"})
		return
	}

	cardID, err := strconv.ParseUint(c.Param("cardId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card id"})
		return
	}

	if err := h.starRepo.RemoveMember(uint(starID), uint(cardID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove card from star"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ExportStar GET /api/v1/stars/:id/export
func (h *StarHandler) ExportStar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid star id"})
		return
	}

	star, err := h.starRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "star not found"})
		return
	}

	if star.Members == nil {
		star.Members = []domain.StarMember{}
	}

	// ✅ Считаем статистику
	totalPV := 0
	totalArmor := 0
	for _, member := range star.Members {
		if member.Card != nil {
			totalPV += member.Card.PointValue
			totalArmor += member.Card.Armor
		}
	}

	// Подготовка данных для экспорта
	exportData := gin.H{
		"star_name":    star.Name,
		"star_id":      star.ID,
		"faction":      star.Faction,
		"members":      star.Members,
		"total_pv":     totalPV,
		"total_armor":  totalArmor,
		"member_count": len(star.Members),
		"timestamp":    time.Now().Format("2006-01-02 15:04:05"),
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, exportData)
}

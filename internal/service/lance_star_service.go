package service

import (
	"Alpha_Strike_Helper/internal/domain"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ValidateTechBaseForLance проверяет, что все мехи Inner Sphere или Mixed
func ValidateTechBaseForLance(lance *domain.Lance) error {
	if len(lance.Members) == 0 {
		return nil
	}

	for _, member := range lance.Members {
		if member.Card == nil {
			continue
		}

		unitType := strings.TrimSpace(member.Card.UnitType)
		// Lance in your UI is designed for BattleMechs only.
		if unitType != "" && unitType != "BATTLEMECH" {
			return fmt.Errorf("Lance может содержать только BATTLEMECH (мехи), а не %s", unitType)
		}

		techBase := strings.TrimSpace(member.Card.TechBase)
		// Some catalogs may not provide tech-base.
		// Treat empty as "unknown" and skip strict validation.
		if techBase == "" {
			continue
		}
		// Lance может содержать: Inner Sphere, Mixed, Primitive
		if techBase != "Inner Sphere" && techBase != "Mixed" && techBase != "Primitive" {
			return fmt.Errorf("Lance может содержать только Inner Sphere/Mixed/Primitive, а не %s", techBase)
		}
	}

	return nil
}

// ValidateTechBaseForStar проверяет, что все мехи Clan или Mixed
func ValidateTechBaseForStar(star *domain.Star) error {
	if len(star.Members) == 0 {
		return nil
	}

	for _, member := range star.Members {
		if member.Card == nil {
			continue
		}

		unitType := strings.TrimSpace(member.Card.UnitType)
		// Star in your UI is designed for Clan BattleMechs only.
		if unitType != "" && unitType != "BATTLEMECH" {
			return fmt.Errorf("Star может содержать только BATTLEMECH (мехи), а не %s", unitType)
		}

		techBase := strings.TrimSpace(member.Card.TechBase)
		// Some catalogs may not provide tech-base.
		// Treat empty as "unknown" and skip strict validation.
		if techBase == "" {
			continue
		}
		// Star может содержать: Clan, Mixed
		if techBase != "Clan" && techBase != "Mixed" {
			return fmt.Errorf("Star может содержать только Clan/Mixed, а не %s", techBase)
		}
	}

	return nil
}

// CalculateLanceStats вычисляет статистику Lance
func CalculateLanceStats(lance *domain.Lance) *domain.LanceStats {
	stats := &domain.LanceStats{
		LanceID:   lance.ID,
		LanceName: lance.Name,
		TechBases: []string{},
		Roles:     []string{},
	}

	if len(lance.Members) == 0 {
		return stats
	}

	stats.UnitCount = len(lance.Members)
	techBaseMap := make(map[string]bool)
	roleMap := make(map[string]bool)
	var tmms []float64

	for _, member := range lance.Members {
		if member.Card == nil {
			continue
		}

		card := member.Card

		// Tonnage (string, нужен парсинг)
		if tonnageStr := strings.TrimSpace(card.Tonnage); tonnageStr != "" {
			if tonnage, err := strconv.Atoi(tonnageStr); err == nil && tonnage > 0 {
				stats.TotalTonnage += tonnage
			}
		}

		// PointValue
		if card.PointValue > 0 {
			stats.TotalPV += card.PointValue
		}

		// TMM
		if card.TMM > 0 {
			tmms = append(tmms, float64(card.TMM))
		}

		// Damage (short range) - парсим DamageShort
		if damageStr := strings.TrimSpace(card.DamageShort); damageStr != "" {
			if damage, err := strconv.Atoi(damageStr); err == nil && damage > 0 {
				stats.TotalDamage += damage
			}
		}

		// TechBase и Role
		if card.TechBase != "" {
			techBaseMap[card.TechBase] = true
		}
		if card.Role != "" {
			roleMap[card.Role] = true
		}
	}

	// Усреднённые значения
	if stats.UnitCount > 0 {
		stats.AvgTonnage = float64(stats.TotalTonnage) / float64(stats.UnitCount)
		stats.AvgPV = float64(stats.TotalPV) / float64(stats.UnitCount)
		stats.TotalDamage = stats.TotalDamage / stats.UnitCount // Среднее damage
	}

	// Average TMM
	if len(tmms) > 0 {
		sum := 0.0
		for _, tmm := range tmms {
			sum += tmm
		}
		stats.AvgTMM = math.Round((sum/float64(len(tmms)))*100) / 100
	}

	// Собрать уникальные tech bases и roles
	for tb := range techBaseMap {
		stats.TechBases = append(stats.TechBases, tb)
	}
	for r := range roleMap {
		stats.Roles = append(stats.Roles, r)
	}

	return stats
}

// CalculateStarStats вычисляет статистику Star
func CalculateStarStats(star *domain.Star) *domain.StarStats {
	stats := &domain.StarStats{
		StarID:    star.ID,
		StarName:  star.Name,
		TechBases: []string{},
		Roles:     []string{},
	}

	if len(star.Members) == 0 {
		return stats
	}

	stats.UnitCount = len(star.Members)
	techBaseMap := make(map[string]bool)
	roleMap := make(map[string]bool)
	var tmms []float64

	for _, member := range star.Members {
		if member.Card == nil {
			continue
		}

		card := member.Card

		// Tonnage (string, нужен парсинг)
		if tonnageStr := strings.TrimSpace(card.Tonnage); tonnageStr != "" {
			if tonnage, err := strconv.Atoi(tonnageStr); err == nil && tonnage > 0 {
				stats.TotalTonnage += tonnage
			}
		}

		// PointValue
		if card.PointValue > 0 {
			stats.TotalPV += card.PointValue
		}

		// TMM
		if card.TMM > 0 {
			tmms = append(tmms, float64(card.TMM))
		}

		// Damage (short range) - парсим DamageShort
		if damageStr := strings.TrimSpace(card.DamageShort); damageStr != "" {
			if damage, err := strconv.Atoi(damageStr); err == nil && damage > 0 {
				stats.TotalDamage += damage
			}
		}

		// TechBase и Role
		if card.TechBase != "" {
			techBaseMap[card.TechBase] = true
		}
		if card.Role != "" {
			roleMap[card.Role] = true
		}
	}

	// Усреднённые значения
	if stats.UnitCount > 0 {
		stats.AvgTonnage = float64(stats.TotalTonnage) / float64(stats.UnitCount)
		stats.AvgPV = float64(stats.TotalPV) / float64(stats.UnitCount)
		stats.TotalDamage = stats.TotalDamage / stats.UnitCount // Среднее damage
	}

	// Average TMM
	if len(tmms) > 0 {
		sum := 0.0
		for _, tmm := range tmms {
			sum += tmm
		}
		stats.AvgTMM = math.Round((sum/float64(len(tmms)))*100) / 100
	}

	// Собрать уникальные tech bases и roles
	for tb := range techBaseMap {
		stats.TechBases = append(stats.TechBases, tb)
	}
	for r := range roleMap {
		stats.Roles = append(stats.Roles, r)
	}

	return stats
}

// ValidateLanceComposition проверяет корректность состава Lance
func ValidateLanceComposition(lance *domain.Lance) error {
	// Минимум 3 меха, максимум 4
	if len(lance.Members) < 3 {
		return fmt.Errorf("Lance должен содержать минимум 3 меха, текущее количество: %d", len(lance.Members))
	}
	if len(lance.Members) > 4 {
		return fmt.Errorf("Lance не может содержать более 4 мехов, текущее количество: %d", len(lance.Members))
	}

	// Проверка позиций (1-4)
	positionMap := make(map[int]bool)
	for _, member := range lance.Members {
		if member.Position < 1 || member.Position > 4 {
			return fmt.Errorf("позиция меха должна быть от 1 до 4, получена: %d", member.Position)
		}
		if positionMap[member.Position] {
			return fmt.Errorf("позиция %d занята несколько раз", member.Position)
		}
		positionMap[member.Position] = true
	}

	// Проверка Tech Base
	if err := ValidateTechBaseForLance(lance); err != nil {
		return err
	}

	return nil
}

// ValidateStarComposition проверяет корректность состава Star
func ValidateStarComposition(star *domain.Star) error {
	// Минимум 3 меха, максимум 5
	if len(star.Members) < 3 {
		return fmt.Errorf("Star должна содержать минимум 3 меха, текущее количество: %d", len(star.Members))
	}
	if len(star.Members) > 5 {
		return fmt.Errorf("Star не может содержать более 5 мехов, текущее количество: %d", len(star.Members))
	}

	// Проверка позиций (1-5)
	positionMap := make(map[int]bool)
	for _, member := range star.Members {
		if member.Position < 1 || member.Position > 5 {
			return fmt.Errorf("позиция меха должна быть от 1 до 5, получена: %d", member.Position)
		}
		if positionMap[member.Position] {
			return fmt.Errorf("позиция %d занята несколько раз", member.Position)
		}
		positionMap[member.Position] = true
	}

	// Проверка Tech Base
	if err := ValidateTechBaseForStar(star); err != nil {
		return err
	}

	return nil
}

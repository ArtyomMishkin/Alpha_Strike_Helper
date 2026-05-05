package masterunitlist

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"Alpha_Strike_Helper/internal/domain"
	"gorm.io/datatypes"
)

var moveFirstNumberRe = regexp.MustCompile(`\d+`)

func MapUnitToCard(
	u Unit,
	availableFactions []string,
	availableFactionGroups []string,
	availableEras []string,
	factionEra map[string][]string,
) domain.Card {
	return domain.Card{
		ModelNumber:            strconv.Itoa(u.ID),
		Name:                   u.Name,
		UnitType:               u.Type.Name,
		Type:                   classifyWeight(u.Tonnage),
		Size:                   strconv.Itoa(u.BFSize),
		Move:                   u.BFMove,
		TMM:                    resolveTMM(u.BFTMM, u.BFMove),
		PointValue:             u.BFPointValue,
		Armor:                  u.BFArmor,
		Structure:              u.BFStructure,
		DamageShort:            formatDamage(u.BFDamageShort, u.BFDamageShortMin),
		DamageMedium:           formatDamage(u.BFDamageMedium, u.BFDamageMediumMin),
		DamageLong:             formatDamage(u.BFDamageLong, u.BFDamageLongMin),
		Overheat:               u.BFOverheat,
		Abilities:              abilitiesToJSON(u.BFAbilities),
		Tonnage:                formatTonnage(u.Tonnage),
		TechBase:               u.Technology.Name,
		Role:                   u.Role.Name,
		Source:                 u.TRO,
		Faction:                firstOrEmpty(availableFactions),
		FactionGroup:           firstOrEmpty(availableFactionGroups),
		Era:                    firstOrEmpty(availableEras),
		AvailableFactions:      datatypes.JSONSlice[string](availableFactions),
		AvailableFactionGroups: datatypes.JSONSlice[string](availableFactionGroups),
		AvailableEras:          datatypes.JSONSlice[string](availableEras),
		FactionEraAvailability: mustJSON(factionEra),
		ImageURL:               u.ImageURL,
		CardURL:                "/Unit/Card/" + strconv.Itoa(u.ID) + "?skill=4",
	}
}

func resolveTMM(sourceTMM int, move string) int {
	if sourceTMM > 0 {
		return sourceTMM
	}
	return estimateTMMFromMove(move)
}

func estimateTMMFromMove(move string) int {
	move = strings.TrimSpace(move)
	if move == "" {
		return 0
	}
	first := moveFirstNumberRe.FindString(move)
	if first == "" {
		return 0
	}
	n, err := strconv.Atoi(first)
	if err != nil {
		return 0
	}
	switch {
	case n <= 4:
		return 0
	case n <= 8:
		return 1
	case n <= 12:
		return 2
	case n <= 18:
		return 3
	case n <= 34:
		return 4
	default:
		return 5
	}
}

func abilitiesToJSON(s string) datatypes.JSONSlice[string] {
	s = strings.TrimSpace(s)
	if s == "" {
		return datatypes.JSONSlice[string]{}
	}
	parts := strings.Fields(s)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	return datatypes.JSONSlice[string](out)
}

func formatDamage(value int, minimal bool) string {
	if minimal {
		return strconv.Itoa(value) + "*"
	}
	return strconv.Itoa(value)
}

func classifyWeight(tonnage float64) string {
	switch {
	case tonnage <= 35:
		return "Light"
	case tonnage <= 55:
		return "Medium"
	case tonnage <= 75:
		return "Heavy"
	default:
		return "Assault"
	}
}

func formatTonnage(tonnage float64) string {
	return fmt.Sprintf("%g", tonnage)
}

func firstOrEmpty(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func sortUnique(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

func mustJSON(v any) datatypes.JSON {
	b, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON([]byte("{}"))
	}
	return datatypes.JSON(b)
}

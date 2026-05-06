package masterunitlist

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"Alpha_Strike_Helper/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ImportOptions struct {
	UnitTypeID         int
	IncludeFactionEras bool
	SkipEraBootstrap   bool
	Replace            bool
	BatchSize          int
	EraIDs             []int
	FactionIDs         []int
	OnProgress         func(ImportProgress)
}

type ImportProgress struct {
	Stage              string
	UnitTypeID         int
	EraID              int
	EraName            string
	FactionID          int
	FactionName        string
	UnitsFetched       int
	PairsDone          int
	PairsTotal         int
	DiscoveredFactions int
	SelectedEras       int
	Error              string
}

type ImportStats struct {
	UnitsTotal        int
	Upserted          int
	Factions          int
	Eras              int
	PairCalls         int
	DuplicatesBefore  int64
	DuplicatesRemoved int64
}

type Importer struct {
	client *Client
	db     *gorm.DB
}

func NewImporter(client *Client, db *gorm.DB) *Importer {
	return &Importer{client: client, db: db}
}

func (imp *Importer) Import(ctx context.Context, opts ImportOptions) (ImportStats, error) {
	var stats ImportStats
	emitProgress := func(p ImportProgress) {
		if opts.OnProgress == nil {
			return
		}
		p.UnitTypeID = opts.UnitTypeID
		opts.OnProgress(p)
	}

	if imp.client == nil || imp.db == nil {
		return stats, fmt.Errorf("client or db is nil")
	}
	if opts.UnitTypeID == 0 {
		opts.UnitTypeID = 18
	}
	if opts.BatchSize <= 0 {
		opts.BatchSize = 300
	}

	eras := filterEras(DefaultEras(), opts.EraIDs)
	stats.Eras = len(eras)

	factions, err := imp.discoverFactions(ctx)
	if err != nil {
		return stats, err
	}
	factions = filterFactionsByIDs(factions, opts.FactionIDs)
	stats.Factions = len(factions)
	log.Printf("discovered factions: %d", len(factions))
	totalPairs := len(eras) * len(factions)
	emitProgress(ImportProgress{
		Stage:              "discovered",
		PairsDone:          0,
		PairsTotal:         totalPairs,
		DiscoveredFactions: len(factions),
		SelectedEras:       len(eras),
	})

	unitByID := make(map[int]Unit, 8192)
	unitEras := make(map[int][]string, 8192)
	unitFactions := make(map[int][]string, 8192)
	unitFactionGroups := make(map[int][]string, 8192)
	unitFactionEra := make(map[int]map[string][]string, 8192)

	factionGroups, err := imp.client.FactionGroups(ctx, factions)
	if err != nil {
		return stats, fmt.Errorf("load faction groups: %w", err)
	}
	for i := range factions {
		if g, ok := factionGroups[factions[i].Label]; ok {
			factions[i].Group = g
		}
	}

	// Process per era -> per faction with explicit "ok" logs.
	for _, era := range eras {
		activeFactions := 0
		log.Printf("era start: %s (%d)", era.Name, era.ID)
		emitProgress(ImportProgress{
			Stage:              "era_start",
			EraID:              era.ID,
			EraName:            era.Name,
			PairsDone:          stats.PairCalls,
			PairsTotal:         totalPairs,
			DiscoveredFactions: len(factions),
			SelectedEras:       len(eras),
		})
		for _, f := range factions {
			stats.PairCalls++
			units, err := imp.client.QuickList(ctx, map[string]string{
				"Types":         strconv.Itoa(opts.UnitTypeID),
				"Factions":      strconv.Itoa(f.Value),
				"AvailableEras": eraIDString(era.ID),
			})
			if err != nil {
				emitProgress(ImportProgress{
					Stage:              "error",
					EraID:              era.ID,
					EraName:            era.Name,
					FactionID:          f.Value,
					FactionName:        f.Label,
					PairsDone:          stats.PairCalls,
					PairsTotal:         totalPairs,
					DiscoveredFactions: len(factions),
					SelectedEras:       len(eras),
					Error:              err.Error(),
				})
				return stats, fmt.Errorf("quicklist faction=%d era=%d failed: %w", f.Value, era.ID, err)
			}
			emitProgress(ImportProgress{
				Stage:              "pair_done",
				EraID:              era.ID,
				EraName:            era.Name,
				FactionID:          f.Value,
				FactionName:        f.Label,
				UnitsFetched:       len(units),
				PairsDone:          stats.PairCalls,
				PairsTotal:         totalPairs,
				DiscoveredFactions: len(factions),
				SelectedEras:       len(eras),
			})
			if len(units) == 0 {
				continue
			}
			activeFactions++
			log.Printf("ok era=%s faction=%s(%d) units=%d", era.Name, f.Label, f.Value, len(units))
			for _, u := range units {
				unitByID[u.ID] = u
				unitEras[u.ID] = append(unitEras[u.ID], era.Name)
				unitFactions[u.ID] = append(unitFactions[u.ID], f.Label)
				if f.Group != "" {
					unitFactionGroups[u.ID] = append(unitFactionGroups[u.ID], f.Group)
				}
				if _, ok := unitFactionEra[u.ID]; !ok {
					unitFactionEra[u.ID] = make(map[string][]string)
				}
				unitFactionEra[u.ID][f.Label] = append(unitFactionEra[u.ID][f.Label], era.Name)
			}
		}
		log.Printf("era done: %s active_factions=%d", era.Name, activeFactions)
		emitProgress(ImportProgress{
			Stage:              "era_done",
			EraID:              era.ID,
			EraName:            era.Name,
			PairsDone:          stats.PairCalls,
			PairsTotal:         totalPairs,
			DiscoveredFactions: len(factions),
			SelectedEras:       len(eras),
		})
	}

	cards := make([]domain.Card, 0, len(unitByID))
	for id, u := range unitByID {
		factionsList := sortUnique(unitFactions[id])
		factionGroupsList := sortUnique(unitFactionGroups[id])
		erasList := sortUnique(unitEras[id])
		factionEra := make(map[string][]string, len(unitFactionEra[id]))
		for factionName, eraList := range unitFactionEra[id] {
			factionEra[factionName] = sortUnique(eraList)
		}
		cards = append(cards, MapUnitToCard(u, factionsList, factionGroupsList, erasList, factionEra))
	}
	stats.UnitsTotal = len(cards)

	dupBefore, err := imp.countDuplicateModelNumbers()
	if err != nil {
		return stats, fmt.Errorf("count duplicates before import: %w", err)
	}
	stats.DuplicatesBefore = dupBefore
	if dupBefore > 0 {
		removed, err := imp.removeDuplicateModelNumbers()
		if err != nil {
			return stats, fmt.Errorf("remove duplicates before import: %w", err)
		}
		stats.DuplicatesRemoved += removed
		log.Printf("duplicates removed before import: %d", removed)
	}

	if opts.Replace {
		if err := imp.db.Exec("DELETE FROM cards").Error; err != nil {
			return stats, fmt.Errorf("delete cards: %w", err)
		}
	}

	if len(cards) > 0 {
		if err := imp.db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "model_number"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name", "unit_type", "type", "size", "move", "tmm", "point_value", "armor", "structure",
				"damage_short", "damage_medium", "damage_long", "overheat", "abilities", "tonnage",
				"tech_base", "role", "source", "faction", "faction_group", "era", "available_factions",
				"available_faction_groups", "available_eras",
				"faction_era_availability", "image_url", "card_url",
			}),
		}).CreateInBatches(cards, opts.BatchSize).Error; err != nil {
			return stats, fmt.Errorf("upsert cards: %w", err)
		}
	}
	stats.Upserted = len(cards)
	emitProgress(ImportProgress{
		Stage:              "db_upsert_done",
		UnitsFetched:       len(cards),
		PairsDone:          stats.PairCalls,
		PairsTotal:         totalPairs,
		DiscoveredFactions: len(factions),
		SelectedEras:       len(eras),
	})

	dupAfter, err := imp.countDuplicateModelNumbers()
	if err != nil {
		return stats, fmt.Errorf("count duplicates after import: %w", err)
	}
	if dupAfter > 0 {
		removed, err := imp.removeDuplicateModelNumbers()
		if err != nil {
			return stats, fmt.Errorf("remove duplicates after import: %w", err)
		}
		stats.DuplicatesRemoved += removed
		log.Printf("duplicates removed after import: %d", removed)
	}

	emitProgress(ImportProgress{
		Stage:              "completed",
		UnitsFetched:       len(cards),
		PairsDone:          stats.PairCalls,
		PairsTotal:         totalPairs,
		DiscoveredFactions: len(factions),
		SelectedEras:       len(eras),
	})
	return stats, nil
}

func (imp *Importer) discoverFactions(ctx context.Context) ([]LabelValue, error) {
	terms := []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
		"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}
	seen := map[int]LabelValue{}
	for _, term := range terms {
		items, err := imp.client.FactionAutocomplete(ctx, term)
		if err != nil {
			return nil, fmt.Errorf("autocomplete term=%q: %w", term, err)
		}
		for _, item := range items {
			if strings.TrimSpace(item.Label) == "" || item.Value == 0 {
				continue
			}
			seen[item.Value] = item
		}
	}

	out := make([]LabelValue, 0, len(seen))
	for _, v := range seen {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Label < out[j].Label
	})
	return out, nil
}

func filterEras(all []Era, selectedIDs []int) []Era {
	if len(selectedIDs) == 0 {
		return all
	}
	selected := make(map[int]struct{}, len(selectedIDs))
	for _, id := range selectedIDs {
		selected[id] = struct{}{}
	}
	out := make([]Era, 0, len(all))
	for _, era := range all {
		if _, ok := selected[era.ID]; ok {
			out = append(out, era)
		}
	}
	return out
}

func filterFactionsByIDs(all []LabelValue, selectedIDs []int) []LabelValue {
	if len(selectedIDs) == 0 {
		return all
	}
	selected := make(map[int]struct{}, len(selectedIDs))
	for _, id := range selectedIDs {
		selected[id] = struct{}{}
	}
	out := make([]LabelValue, 0, len(selectedIDs))
	for _, f := range all {
		if _, ok := selected[f.Value]; ok {
			out = append(out, f)
		}
	}
	return out
}

func (imp *Importer) countDuplicateModelNumbers() (int64, error) {
	var n int64
	err := imp.db.Raw(`
		SELECT COALESCE(COUNT(*), 0)
		FROM (
			SELECT model_number
			FROM cards
			WHERE COALESCE(model_number, '') <> ''
			GROUP BY model_number
			HAVING COUNT(*) > 1
		) d
	`).Scan(&n).Error
	return n, err
}

func (imp *Importer) removeDuplicateModelNumbers() (int64, error) {
	res := imp.db.Exec(`
		DELETE FROM cards a
		USING cards b
		WHERE a.id < b.id
		  AND COALESCE(a.model_number, '') <> ''
		  AND a.model_number = b.model_number
	`)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

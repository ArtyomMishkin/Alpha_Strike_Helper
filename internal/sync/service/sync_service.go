package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"Alpha_Strike_Helper/internal/domain"
	"Alpha_Strike_Helper/internal/sync/masterunitlist"
	"gorm.io/gorm"
)

type RunRequest struct {
	BaseURL            string `json:"base_url"`
	UnitTypeIDs        []int  `json:"unit_type_ids"`
	EraIDs             []int  `json:"era_ids"`
	FactionIDs         []int  `json:"faction_ids"`
	ReplaceFirst       bool   `json:"replace_first"`
	IncludeFactionEras bool   `json:"include_faction_eras"`
	SkipEraBootstrap   bool   `json:"skip_era_bootstrap"`
	HTTPTimeoutSeconds int    `json:"http_timeout_seconds"`
	BatchSize          int    `json:"batch_size"`
}

type PerTypeRunResult struct {
	UnitTypeID int                         `json:"unit_type_id"`
	Stats      *masterunitlist.ImportStats `json:"stats,omitempty"`
	Error      string                      `json:"error,omitempty"`
}

type RunStatus struct {
	RunID      string             `json:"run_id,omitempty"`
	Status     string             `json:"status"`
	StartedAt  *time.Time         `json:"started_at,omitempty"`
	FinishedAt *time.Time         `json:"finished_at,omitempty"`
	Request    *RunRequest        `json:"request,omitempty"`
	Progress   *RunProgress       `json:"progress,omitempty"`
	Results    []PerTypeRunResult `json:"per_type_results,omitempty"`
	Error      string             `json:"error,omitempty"`
}

type RunProgress struct {
	Stage                string `json:"stage"`
	UnitTypeID           int    `json:"unit_type_id"`
	EraID                int    `json:"era_id,omitempty"`
	EraName              string `json:"era_name,omitempty"`
	FactionID            int    `json:"faction_id,omitempty"`
	FactionName          string `json:"faction_name,omitempty"`
	UnitsFetched         int    `json:"units_fetched,omitempty"`
	PairsDone            int    `json:"pairs_done,omitempty"`
	PairsTotal           int    `json:"pairs_total,omitempty"`
	DiscoveredFactions   int    `json:"discovered_factions,omitempty"`
	SelectedEras         int    `json:"selected_eras,omitempty"`
	CurrentTypeIndex     int    `json:"current_type_index,omitempty"`
	TotalTypes           int    `json:"total_types,omitempty"`
	CurrentTypeCompleted bool   `json:"current_type_completed,omitempty"`
	Error                string `json:"error,omitempty"`
	UpdatedAt            string `json:"updated_at"`
}

type SyncService struct {
	db *gorm.DB

	mu      sync.RWMutex
	current RunStatus
}

func NewSyncService(db *gorm.DB) *SyncService {
	return &SyncService{
		db:      db,
		current: RunStatus{Status: "idle"},
	}
}

func (s *SyncService) CurrentStatus() RunStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneStatus(s.current)
}

func (s *SyncService) CardsCount(ctx context.Context) (int64, error) {
	if s.db == nil {
		return 0, fmt.Errorf("database is not initialized")
	}
	var n int64
	if err := s.db.WithContext(ctx).Model(&domain.Card{}).Count(&n).Error; err != nil {
		return 0, err
	}
	return n, nil
}

func (s *SyncService) StartRun(ctx context.Context, req RunRequest) (RunStatus, error) {
	s.mu.Lock()
	if s.current.Status == "running" {
		cur := cloneStatus(s.current)
		s.mu.Unlock()
		return cur, fmt.Errorf("sync already running")
	}

	normalized := normalizeRequest(req)
	now := time.Now()
	runID := fmt.Sprintf("run-%d", now.UnixNano())
	s.current = RunStatus{
		RunID:     runID,
		Status:    "running",
		StartedAt: &now,
		Request:   &normalized,
		Results:   []PerTypeRunResult{},
	}
	initial := cloneStatus(s.current)
	s.mu.Unlock()

	go s.run(context.Background(), runID, normalized)
	return initial, nil
}

func (s *SyncService) run(ctx context.Context, runID string, req RunRequest) {
	timeout := time.Duration(req.HTTPTimeoutSeconds) * time.Second
	client := masterunitlist.NewClient(req.BaseURL, timeout)
	importer := masterunitlist.NewImporter(client, s.db)

	results := make([]PerTypeRunResult, 0, len(req.UnitTypeIDs))
	var topErr error

	for idx, unitTypeID := range req.UnitTypeIDs {
		stats, err := importer.Import(ctx, masterunitlist.ImportOptions{
			UnitTypeID:         unitTypeID,
			IncludeFactionEras: req.IncludeFactionEras,
			SkipEraBootstrap:   req.SkipEraBootstrap,
			Replace:            req.ReplaceFirst && idx == 0,
			BatchSize:          req.BatchSize,
			EraIDs:             req.EraIDs,
			FactionIDs:         req.FactionIDs,
			OnProgress: func(p masterunitlist.ImportProgress) {
				s.mu.Lock()
				defer s.mu.Unlock()
				if s.current.RunID != runID {
					return
				}
				s.current.Progress = &RunProgress{
					Stage:                p.Stage,
					UnitTypeID:           p.UnitTypeID,
					EraID:                p.EraID,
					EraName:              p.EraName,
					FactionID:            p.FactionID,
					FactionName:          p.FactionName,
					UnitsFetched:         p.UnitsFetched,
					PairsDone:            p.PairsDone,
					PairsTotal:           p.PairsTotal,
					DiscoveredFactions:   p.DiscoveredFactions,
					SelectedEras:         p.SelectedEras,
					CurrentTypeIndex:     idx + 1,
					TotalTypes:           len(req.UnitTypeIDs),
					CurrentTypeCompleted: false,
					Error:                p.Error,
					UpdatedAt:            time.Now().UTC().Format(time.RFC3339),
				}
			},
		})
		if err != nil {
			results = append(results, PerTypeRunResult{
				UnitTypeID: unitTypeID,
				Error:      err.Error(),
			})
			s.mu.Lock()
			if s.current.RunID == runID {
				s.current.Progress = &RunProgress{
					Stage:                "failed",
					UnitTypeID:           unitTypeID,
					CurrentTypeIndex:     idx + 1,
					TotalTypes:           len(req.UnitTypeIDs),
					CurrentTypeCompleted: false,
					Error:                err.Error(),
					UpdatedAt:            time.Now().UTC().Format(time.RFC3339),
				}
			}
			s.mu.Unlock()
			if topErr == nil {
				topErr = err
			}
			continue
		}
		statsCopy := stats
		results = append(results, PerTypeRunResult{
			UnitTypeID: unitTypeID,
			Stats:      &statsCopy,
		})
		s.mu.Lock()
		if s.current.RunID == runID {
			s.current.Progress = &RunProgress{
				Stage:                "unit_type_completed",
				UnitTypeID:           unitTypeID,
				UnitsFetched:         statsCopy.UnitsTotal,
				PairsDone:            statsCopy.PairCalls,
				PairsTotal:           statsCopy.Factions * statsCopy.Eras,
				DiscoveredFactions:   statsCopy.Factions,
				SelectedEras:         statsCopy.Eras,
				CurrentTypeIndex:     idx + 1,
				TotalTypes:           len(req.UnitTypeIDs),
				CurrentTypeCompleted: true,
				UpdatedAt:            time.Now().UTC().Format(time.RFC3339),
			}
		}
		s.mu.Unlock()
	}

	finished := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	// Ignore stale goroutine write if a newer run somehow started.
	if s.current.RunID != runID {
		return
	}
	s.current.FinishedAt = &finished
	s.current.Results = results
	if topErr != nil {
		s.current.Status = "failed"
		s.current.Error = topErr.Error()
		return
	}
	s.current.Status = "completed"
	s.current.Error = ""
}

func normalizeRequest(req RunRequest) RunRequest {
	if req.BaseURL == "" {
		req.BaseURL = "https://masterunitlist.azurewebsites.net"
	}
	if len(req.UnitTypeIDs) == 0 {
		req.UnitTypeIDs = []int{18, 19, 17, 21}
	}
	if req.HTTPTimeoutSeconds <= 0 {
		req.HTTPTimeoutSeconds = 180
	}
	if req.BatchSize <= 0 {
		req.BatchSize = 300
	}
	if !req.IncludeFactionEras {
		req.IncludeFactionEras = true
	}
	return req
}

func cloneStatus(in RunStatus) RunStatus {
	out := in
	if in.Request != nil {
		req := *in.Request
		out.Request = &req
	}
	if in.StartedAt != nil {
		t := *in.StartedAt
		out.StartedAt = &t
	}
	if in.FinishedAt != nil {
		t := *in.FinishedAt
		out.FinishedAt = &t
	}
	if in.Progress != nil {
		p := *in.Progress
		out.Progress = &p
	}
	if in.Results != nil {
		out.Results = append([]PerTypeRunResult(nil), in.Results...)
	}
	return out
}

func BuildProgressMessage(status RunStatus) string {
	if status.Progress == nil {
		switch status.Status {
		case "idle":
			return "Импорт не запущен"
		case "completed":
			return "Импорт завершен"
		case "failed":
			return "Импорт завершился с ошибкой"
		default:
			return "Прогресс пока недоступен"
		}
	}
	p := status.Progress
	typePart := ""
	if p.TotalTypes > 0 && p.CurrentTypeIndex > 0 {
		typePart = fmt.Sprintf("Type %d [%d/%d]", p.UnitTypeID, p.CurrentTypeIndex, p.TotalTypes)
	} else if p.UnitTypeID > 0 {
		typePart = fmt.Sprintf("Type %d", p.UnitTypeID)
	}
	pairPart := ""
	if p.PairsTotal > 0 {
		pairPart = fmt.Sprintf("Pairs %d/%d", p.PairsDone, p.PairsTotal)
	}
	eraPart := ""
	if p.EraName != "" {
		eraPart = fmt.Sprintf("Era: %s", p.EraName)
	}
	factionPart := ""
	if p.FactionName != "" {
		factionPart = fmt.Sprintf("Faction: %s", p.FactionName)
	}
	unitsPart := ""
	if p.UnitsFetched > 0 {
		unitsPart = fmt.Sprintf("Units: %d", p.UnitsFetched)
	}
	stagePart := fmt.Sprintf("Stage: %s", p.Stage)
	parts := make([]string, 0, 6)
	for _, v := range []string{typePart, pairPart, eraPart, factionPart, unitsPart, stagePart} {
		if v != "" {
			parts = append(parts, v)
		}
	}
	if p.Error != "" {
		parts = append(parts, "Error: "+p.Error)
	}
	if len(parts) == 0 {
		return "Прогресс обновляется"
	}
	return joinWithSeparator(parts, " | ")
}

func joinWithSeparator(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += sep + parts[i]
	}
	return out
}

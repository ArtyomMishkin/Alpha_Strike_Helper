package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"Alpha_Strike_Helper/internal/sync/masterunitlist"
	"Alpha_Strike_Helper/pkg/config"
	"Alpha_Strike_Helper/pkg/database"
)

func main() {
	var (
		baseURL            = flag.String("base-url", "https://masterunitlist.azurewebsites.net", "Master Unit List base URL")
		unitTypeIDsCSV     = flag.String("unit-type-ids", "18,19,17,21", "comma-separated MUL type IDs")
		eraIDsCSV          = flag.String("era-ids", "13,247,14", "comma-separated era IDs")
		factionIDsCSV      = flag.String("faction-ids", "", "comma-separated faction IDs")
		includeFactionEras = flag.Bool("include-faction-eras", true, "build detailed faction->eras availability map")
		skipEraBootstrap   = flag.Bool("skip-era-bootstrap", true, "skip broad era-only requests")
		replace            = flag.Bool("replace", false, "replace cards table before each import")
		httpTimeout        = flag.Duration("http-timeout", 120*time.Second, "HTTP timeout per request (MUL can be slow)")
		batchSize          = flag.Int("batch-size", 300, "database upsert batch size")
		interval           = flag.Duration("interval", 7*24*time.Hour, "sync interval")
		runNow             = flag.Bool("run-now", true, "run one sync immediately on startup")
	)
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg := config.Load()
	db, err := database.NewPostgresDB(database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}

	client := masterunitlist.NewClient(*baseURL, *httpTimeout)
	importer := masterunitlist.NewImporter(client, db)
	unitTypeIDs := parseCSVInts(*unitTypeIDsCSV)
	eraIDs := parseCSVInts(*eraIDsCSV)
	factionIDs := parseCSVInts(*factionIDsCSV)

	if len(unitTypeIDs) == 0 {
		log.Fatalf("no unit type ids provided")
	}

	runSync := func(runCtx context.Context) {
		started := time.Now()
		log.Printf("weekly sync started: unit_type_ids=%v era_ids=%v", unitTypeIDs, eraIDs)
		for _, unitTypeID := range unitTypeIDs {
			stats, err := importer.Import(runCtx, masterunitlist.ImportOptions{
				UnitTypeID:         unitTypeID,
				IncludeFactionEras: *includeFactionEras,
				SkipEraBootstrap:   *skipEraBootstrap,
				Replace:            *replace,
				BatchSize:          *batchSize,
				EraIDs:             eraIDs,
				FactionIDs:         factionIDs,
			})
			if err != nil {
				log.Printf("weekly sync failed for unit_type_id=%d: %v", unitTypeID, err)
				continue
			}
			log.Printf("weekly sync complete for unit_type_id=%d: units=%d upserted=%d factions=%d eras=%d pair_calls=%d duplicates_before=%d duplicates_removed=%d",
				unitTypeID, stats.UnitsTotal, stats.Upserted, stats.Factions, stats.Eras, stats.PairCalls, stats.DuplicatesBefore, stats.DuplicatesRemoved)
		}
		log.Printf("weekly sync run finished in %s", time.Since(started).Round(time.Second))
	}

	if *runNow {
		runSync(ctx)
	}

	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	log.Printf("weekly sync scheduler is running: interval=%s", interval.String())
	for {
		select {
		case <-ctx.Done():
			log.Printf("weekly sync scheduler stopped")
			return
		case <-ticker.C:
			runSync(ctx)
		}
	}
}

func parseCSVInts(v string) []int {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			continue
		}
		out = append(out, n)
	}
	return out
}

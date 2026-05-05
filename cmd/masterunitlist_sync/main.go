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
		unitTypeID         = flag.Int("unit-type-id", 18, "MUL type id (18=BattleMech)")
		replace            = flag.Bool("replace", true, "replace existing cards table before import")
		includeFactionEras = flag.Bool("include-faction-eras", true, "build detailed faction->eras availability map")
		skipEraBootstrap   = flag.Bool("skip-era-bootstrap", false, "skip broad era-only requests (useful when an era endpoint is unstable)")
		eraIDsCSV          = flag.String("era-ids", "", "comma-separated era IDs, e.g. 13,247,14")
		factionIDsCSV      = flag.String("faction-ids", "", "comma-separated faction IDs, e.g. 24,29,27")
		httpTimeout        = flag.Duration("http-timeout", 30*time.Second, "HTTP timeout per request")
		batchSize          = flag.Int("batch-size", 300, "database upsert batch size")
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
	stats, err := importer.Import(ctx, masterunitlist.ImportOptions{
		UnitTypeID:         *unitTypeID,
		IncludeFactionEras: *includeFactionEras,
		SkipEraBootstrap:   *skipEraBootstrap,
		Replace:            *replace,
		BatchSize:          *batchSize,
		EraIDs:             parseCSVInts(*eraIDsCSV),
		FactionIDs:         parseCSVInts(*factionIDsCSV),
	})
	if err != nil {
		log.Fatalf("import failed: %v", err)
	}

	log.Printf("import complete: units=%d upserted=%d factions=%d eras=%d pair_calls=%d duplicates_before=%d duplicates_removed=%d",
		stats.UnitsTotal, stats.Upserted, stats.Factions, stats.Eras, stats.PairCalls, stats.DuplicatesBefore, stats.DuplicatesRemoved)
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

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	syncservice "Alpha_Strike_Helper/internal/sync/service"
	"Alpha_Strike_Helper/pkg/config"
	"Alpha_Strike_Helper/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
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
		log.Fatalf("sync-service db connect failed: %v", err)
	}
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("sync-service auto-migrate failed: %v", err)
	}

	svc := syncservice.NewSyncService(db)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "sync-service",
			"time":    time.Now().UTC(),
		})
	})

	router.GET("/sync/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, svc.CurrentStatus())
	})

	router.GET("/sync/progress", func(c *gin.Context) {
		status := svc.CurrentStatus()
		message := syncservice.BuildProgressMessage(status)
		if status.Progress == nil {
			c.JSON(http.StatusOK, gin.H{
				"run_id":  status.RunID,
				"status":  status.Status,
				"message": message,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"run_id":          status.RunID,
			"status":          status.Status,
			"started_at":      status.StartedAt,
			"finished_at":     status.FinishedAt,
			"message":         message,
			"progress":        status.Progress,
			"per_type_result": status.Results,
			"error":           status.Error,
		})
	})

	router.POST("/sync/run", func(c *gin.Context) {
		var req syncservice.RunRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
			return
		}
		status, err := svc.StartRun(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error(), "status": status})
			return
		}
		c.JSON(http.StatusAccepted, status)
	})

	runOnStart := strings.EqualFold(strings.TrimSpace(os.Getenv("RUN_ON_START")), "true")
	if runOnStart {
		count, err := svc.CardsCount(context.Background())
		if err != nil {
			log.Printf("sync-service bootstrap check failed: %v", err)
		} else if count == 0 {
			req := syncservice.RunRequest{
				ReplaceFirst:       true,
				IncludeFactionEras: true,
				HTTPTimeoutSeconds: 180,
				UnitTypeIDs:        []int{18, 19, 17, 21},
			}
			if raw := strings.TrimSpace(os.Getenv("RUN_ON_START_UNIT_TYPE_IDS")); raw != "" {
				req.UnitTypeIDs = parseCSVInts(raw)
			}
			if v := strings.TrimSpace(os.Getenv("RUN_ON_START_HTTP_TIMEOUT_SECONDS")); v != "" {
				if n, convErr := strconv.Atoi(v); convErr == nil && n > 0 {
					req.HTTPTimeoutSeconds = n
				}
			}
			if status, startErr := svc.StartRun(context.Background(), req); startErr != nil {
				log.Printf("sync-service bootstrap run start failed: %v", startErr)
			} else {
				log.Printf("sync-service bootstrap run started: run_id=%s unit_type_ids=%v", status.RunID, req.UnitTypeIDs)
			}
		} else {
			log.Printf("sync-service bootstrap skipped: cards already present (%d)", count)
		}
	}

	port := "8081"
	if v := cfg.ServerPort; v != "" {
		// Keep sync-service separate from main app port by default.
		if v == "8080" {
			port = "8081"
		} else {
			port = v
		}
	}
	log.Printf("sync-service listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("sync-service start failed: %v", err)
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

package cron

import (
	"context"
	"net/http"
	"time"

	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/WahyuSiddarta/be_saham_go/utime"
	"github.com/getsentry/sentry-go"
	robfigcron "github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

const STOCK_DATASOURCE = "https://api.datasectors.com/api/stocks/v2/"

type Runner struct {
	logger     *zerolog.Logger
	httpClient *http.Client
}

func NewRunner(logger *zerolog.Logger) *Runner {
	return &Runner{
		logger: logger,
		httpClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        20,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

func (r *Runner) Start(ctx context.Context) {
	r.logger.Info().Msg("Cron runner started")

	scheduler := robfigcron.New()

	r.UpsertStockInformation(ctx)
	if _, err := scheduler.AddFunc("0 14 * * *", func() {
		r.UpsertStockInformation(ctx)
	}); err != nil {
		r.captureException(err, map[string]string{
			"module": "cron",
			"job":    "upsertStockInformation",
			"action": "register_schedule",
		}, nil)
		r.logger.Fatal().Err(err).Msg("Failed to register cron schedule")
		return
	}
	r.logger.Info().Str("job", "upsertStockInformation").Str("schedule", "0 14 * * *").Msg("Cron job registered")

	scheduler.Start()

	<-ctx.Done()
	stopCtx := scheduler.Stop()
	<-stopCtx.Done()
	r.logger.Info().Msg("Cron runner stopped")
}

func (r *Runner) UpsertStockInformation(ctx context.Context) {
	startTime := utime.Utime.Now().ToTime()
	r.logger.Info().Str("job", "upsertStockInformation").Msg("Cron job execution started")

	const processStockInterval = 20 * time.Second

	stockRepo := models.NewStockRepository()
	targetStocks, err := stockRepo.GetStockApiKey()
	if err != nil {
		r.captureException(err, map[string]string{
			"module": "cron",
			"job":    "upsertStockInformation",
			"action": "get_stock_api_key",
		}, nil)
		r.logger.Error().Err(err).Str("job", "upsertStockInformation").Msg("Failed to get stock API key")
		return
	}

	if len(targetStocks) == 0 {
		r.logger.Warn().Str("job", "upsertStockInformation").Msg("No target stocks configured")
		return
	}

	for index, stock := range targetStocks {
		select {
		case <-ctx.Done():
			r.logger.Info().Str("job", "upsertStockInformation").Msg("Cron job canceled")
			return
		default:
		}

		processStartTime := utime.Utime.Now().ToTime()
		r.processStock(ctx, stockRepo, stock)

		if index < len(targetStocks)-1 {
			remainingWait := processStockInterval - utime.Utime.Now().ToTime().Sub(processStartTime)
			if remainingWait > 0 {
				select {
				case <-ctx.Done():
					r.logger.Info().Str("job", "upsertStockInformation").Msg("Cron job canceled while waiting for rate limit")
					return
				case <-time.After(remainingWait):
				}
			}
		}
	}

	r.logger.Info().Str("job", "upsertStockInformation").Dur("duration", utime.Utime.Now().ToTime().Sub(startTime)).Msg("Cron job execution completed")
}

func (r *Runner) captureException(err error, tags map[string]string, extra map[string]interface{}) {
	if err == nil {
		return
	}

	hub := sentry.CurrentHub()
	if hub == nil {
		return
	}

	hub.WithScope(func(scope *sentry.Scope) {
		for key, value := range tags {
			scope.SetTag(key, value)
		}

		for key, value := range extra {
			scope.SetExtra(key, value)
		}

		hub.CaptureException(err)
	})
}

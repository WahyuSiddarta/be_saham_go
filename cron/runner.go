package cron

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/bytedance/sonic"
	robfigcron "github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

const STOCK_DATASOURCE = "https://api.datasectors.com/api/stocks/v2/"

type Runner struct {
	logger     *zerolog.Logger
	httpClient *http.Client
}

func NewRunner(logger *zerolog.Logger, _ time.Duration) *Runner {
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
	if _, err := scheduler.AddFunc("* * * * *", func() {
		r.UpsertStockInformation(ctx)
	}); err != nil {
		r.logger.Fatal().Err(err).Msg("Failed to register cron schedule")
		return
	}
	r.logger.Info().Str("job", "upsertStockInformation").Str("schedule", "* * * * *").Msg("Cron job registered")

	scheduler.Start()

	<-ctx.Done()
	stopCtx := scheduler.Stop()
	<-stopCtx.Done()
	r.logger.Info().Msg("Cron runner stopped")
}

func (r *Runner) UpsertStockInformation(ctx context.Context) {
	startTime := time.Now()
	r.logger.Info().Str("job", "upsertStockInformation").Msg("Cron job execution started")

	stockRepo := models.NewStockRepository()
	targetStocks, err := stockRepo.GetStockApiKey()
	if err != nil {
		r.logger.Error().Err(err).Str("job", "upsertStockInformation").Msg("Failed to get stock API key")
		return
	}

	if len(targetStocks) == 0 {
		r.logger.Warn().Str("job", "upsertStockInformation").Msg("No target stocks configured")
		return
	}

	for _, stock := range targetStocks {
		select {
		case <-ctx.Done():
			r.logger.Info().Str("job", "upsertStockInformation").Msg("Cron job canceled")
			return
		default:
		}
		r.processStock(ctx, stockRepo, stock)
	}

	r.logger.Info().Str("job", "upsertStockInformation").Dur("duration", time.Since(startTime)).Msg("Cron job execution completed")
}

// fetchStockData fetches earnings and equities data concurrently for a single stock.
func (r *Runner) fetchStockData(ctx context.Context, stock models.StockInformation) (earningsResp *helper.ExternalResponse, equitiesResp *helper.ExternalResponse, earningsErr error, equitiesErr error) {
	var wg sync.WaitGroup
	wg.Add(2)
	r.logger.Info().
		Str("job", "fetchStockData").
		Str("IMPORTANT", stock.ApiKey).
		Msg("Fetching stock data")
	go func() {
		defer wg.Done()
		ctxEarning, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		earningsResp, earningsErr = helper.DoExternalJSONRequest(
			ctxEarning,
			r.httpClient,
			http.MethodGet,
			STOCK_DATASOURCE+"earnings",
			helper.ExternalJSONRequestOptions{
				Headers: map[string]string{"X-API-Key": stock.ApiKey},
				Query:   map[string]string{"symbol": stock.Ticker, "market": "id-id"},
			},
		)
	}()

	go func() {
		defer wg.Done()
		ctxEquities, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		equitiesResp, equitiesErr = helper.DoExternalJSONRequest(
			ctxEquities,
			r.httpClient,
			http.MethodGet,
			STOCK_DATASOURCE+"equities",
			helper.ExternalJSONRequestOptions{
				Headers: map[string]string{"X-API-Key": stock.ApiKey},
				Query:   map[string]string{"symbol": stock.Ticker, "market": "id-id"},
			},
		)
	}()

	wg.Wait()
	return
}

func (r *Runner) processStock(ctx context.Context, stockRepo models.StockRepository, stock models.StockInformation) {
	earningsResp, equitiesResp, earningsErr, equitiesErr := r.fetchStockData(ctx, stock)

	var overviewRecord *models.StockOverviewMetricsRecord
	isOverviewRecordFromEarnings := false

	if earningsErr != nil {
		r.logger.Warn().
			Err(earningsErr).
			Str("job", "upsertStockInformation").
			Str("ticker", stock.Ticker).
			Msg("Failed to fetch earnings data; continue with equities data")
	} else {
		var earnings EarningsResponse
		if err := sonic.Unmarshal(earningsResp.Body, &earnings); err != nil {
			r.logger.Error().
				Err(err).
				Str("job", "upsertStockInformation").
				Str("ticker", stock.Ticker).
				Msg("Failed to decode earnings data")
		} else {
			quarterlyRecords, err := earnings.ToQuarterlyHistoryRecords()
			if err != nil {
				r.logger.Error().
					Err(err).
					Str("job", "upsertStockInformation").
					Str("ticker", stock.Ticker).
					Msg("Failed to parse quarterly history records")
			} else if err := stockRepo.UpsertStockEarningQuarterlyHistory(quarterlyRecords); err != nil {
				r.logger.Error().
					Err(err).
					Str("job", "upsertStockInformation").
					Str("ticker", stock.Ticker).
					Msg("Failed to upsert quarterly history records")
			} else {
				r.logger.Info().
					Str("job", "upsertStockInformation").
					Str("ticker", stock.Ticker).
					Str("symbol", earnings.Symbol).
					Int("quarterlyRecordCount", len(quarterlyRecords)).
					Msg("Quarterly history upserted")
			}

			overviewRecord, err = earnings.ToOverviewMetricsRecord()
			if err != nil {
				r.logger.Error().
					Err(err).
					Str("job", "upsertStockInformation").
					Str("ticker", stock.Ticker).
					Msg("Failed to parse overview metrics from earnings data")
			} else {
				isOverviewRecordFromEarnings = true
			}
		}
	}

	if overviewRecord == nil {
		overviewRecord = &models.StockOverviewMetricsRecord{Symbol: stock.Ticker}
	}

	if equitiesErr != nil {
		r.logger.Error().
			Err(equitiesErr).
			Str("job", "upsertStockInformation").
			Str("ticker", stock.Ticker).
			Msg("Failed to fetch overview data")
	} else {
		var equities EquitiesResponse
		if err := sonic.Unmarshal(equitiesResp.Body, &equities); err != nil {
			r.logger.Error().
				Err(err).
				Str("job", "upsertStockInformation").
				Str("ticker", stock.Ticker).
				Msg("Failed to decode overview data")
		} else {
			if err := equities.MergeIntoOverviewMetricsRecord(overviewRecord); err != nil {
				r.logger.Error().
					Err(err).
					Str("job", "upsertStockInformation").
					Str("ticker", stock.Ticker).
					Msg("Failed to merge overview metrics from equities data")
			}
		}
	}

	if err := stockRepo.UpsertStockOverviewMetrics(overviewRecord); err != nil {
		r.logger.Error().
			Err(err).
			Str("job", "upsertStockInformation").
			Str("ticker", stock.Ticker).
			Msg("Failed to upsert overview metrics")
		return
	}

	r.logger.Info().
		Str("job", "upsertStockInformation").
		Str("ticker", stock.Ticker).
		Str("symbol", overviewRecord.Symbol).
		Bool("overviewFromEarnings", isOverviewRecordFromEarnings).
		Msg("Overview metrics upserted")
}

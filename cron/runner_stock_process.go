package cron

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/bytedance/sonic"
)

// fetchStockData fetches earnings and equities data concurrently for a single stock.
func (r *Runner) fetchStockData(ctx context.Context, stock models.StockInformation) (earningsResp *helper.ExternalResponse, equitiesResp *helper.ExternalResponse, earningsErr error, equitiesErr error) {
	var wg sync.WaitGroup
	wg.Add(2)
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
		r.captureException(earningsErr, map[string]string{
			"module": "cron",
			"job":    "upsertStockInformation",
			"action": "fetch_earnings",
		}, map[string]interface{}{
			"ticker": stock.Ticker,
		})
		r.logger.Warn().
			Err(earningsErr).
			Str("job", "upsertStockInformation").
			Str("ticker", stock.Ticker).
			Msg("Failed to fetch earnings data; continue with equities data")
	} else {
		var earnings EarningsResponse
		if err := sonic.Unmarshal(earningsResp.Body, &earnings); err != nil {
			r.captureException(err, map[string]string{
				"module": "cron",
				"job":    "upsertStockInformation",
				"action": "decode_earnings",
			}, map[string]interface{}{
				"ticker": stock.Ticker,
			})
			r.logger.Error().
				Err(err).
				Str("job", "upsertStockInformation").
				Str("ticker", stock.Ticker).
				Msg("Failed to decode earnings data")
		} else {
			quarterlyRecords, err := earnings.ToQuarterlyHistoryRecords()
			if err != nil {
				r.captureException(err, map[string]string{
					"module": "cron",
					"job":    "upsertStockInformation",
					"action": "parse_quarterly_history",
				}, map[string]interface{}{
					"ticker": stock.Ticker,
				})
				r.logger.Error().
					Err(err).
					Str("job", "upsertStockInformation").
					Str("ticker", stock.Ticker).
					Msg("Failed to parse quarterly history records")
			} else if err := stockRepo.UpsertStockEarningQuarterlyHistory(quarterlyRecords); err != nil {
				r.captureException(err, map[string]string{
					"module": "cron",
					"job":    "upsertStockInformation",
					"action": "upsert_quarterly_history",
				}, map[string]interface{}{
					"ticker":               stock.Ticker,
					"quarterlyRecordCount": len(quarterlyRecords),
				})
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
				r.captureException(err, map[string]string{
					"module": "cron",
					"job":    "upsertStockInformation",
					"action": "parse_overview_from_earnings",
				}, map[string]interface{}{
					"ticker": stock.Ticker,
				})
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
		r.captureException(equitiesErr, map[string]string{
			"module": "cron",
			"job":    "upsertStockInformation",
			"action": "fetch_equities",
		}, map[string]interface{}{
			"ticker": stock.Ticker,
		})
		r.logger.Error().
			Err(equitiesErr).
			Str("job", "upsertStockInformation").
			Str("ticker", stock.Ticker).
			Msg("Failed to fetch overview data")
	} else {
		var equities EquitiesResponse
		if err := sonic.Unmarshal(equitiesResp.Body, &equities); err != nil {
			r.captureException(err, map[string]string{
				"module": "cron",
				"job":    "upsertStockInformation",
				"action": "decode_equities",
			}, map[string]interface{}{
				"ticker": stock.Ticker,
			})
			r.logger.Error().
				Err(err).
				Str("job", "upsertStockInformation").
				Str("ticker", stock.Ticker).
				Msg("Failed to decode overview data")
		} else {
			if err := equities.MergeIntoOverviewMetricsRecord(overviewRecord); err != nil {
				r.captureException(err, map[string]string{
					"module": "cron",
					"job":    "upsertStockInformation",
					"action": "merge_overview_from_equities",
				}, map[string]interface{}{
					"ticker": stock.Ticker,
				})
				r.logger.Error().
					Err(err).
					Str("job", "upsertStockInformation").
					Str("ticker", stock.Ticker).
					Msg("Failed to merge overview metrics from equities data")
			}
		}
	}

	if err := stockRepo.UpsertStockOverviewMetrics(overviewRecord); err != nil {
		r.captureException(err, map[string]string{
			"module": "cron",
			"job":    "upsertStockInformation",
			"action": "upsert_overview_metrics",
		}, map[string]interface{}{
			"ticker": stock.Ticker,
			"symbol": overviewRecord.Symbol,
		})
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

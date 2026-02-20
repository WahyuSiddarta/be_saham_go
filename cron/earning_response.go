package cron

import (
	"fmt"
	"sort"

	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/WahyuSiddarta/be_saham_go/models"
)

type EarningsResponse struct {
	Success bool         `json:"success"`
	Symbol  string       `json:"symbol"`
	SecID   string       `json:"secId"`
	Data    EarningsData `json:"data"`
}

type EarningsData struct {
	InstrumentID           string             `json:"InstrumentId"`
	Market                 string             `json:"Market"`
	Currency               string             `json:"Currency"`
	MarketCap              *float64           `json:"MarketCap"`
	LastActualFiscalPeriod string             `json:"LastActualFiscalPeriod"`
	ExpectedReportDate     *string            `json:"ExpectedReportDate"`
	TimeLastUpdated        *string            `json:"TimeLastUpdated"`
	LastActual             EarningsLastActual `json:"LastActual"`
	History                EarningsHistory    `json:"History"`
}

type EarningsLastActual struct {
	EpsActual     *float64 `json:"EpsActual"`
	RevenueActual *float64 `json:"RevenueActual"`
}

type EarningsHistory struct {
	Quarterly map[string]EarningsQuarterlyHistory `json:"quarterly"`
}

type EarningsQuarterlyHistory struct {
	EpsActual                    *float64 `json:"EpsActual"`
	EpsSurprise                  *float64 `json:"EpsSurprise"`
	EpsSurprisePercent           *float64 `json:"EpsSurprisePercent"`
	RevenueActual                *float64 `json:"RevenueActual"`
	RevenueSurprise              *float64 `json:"RevenueSurprise"`
	RevenueSurprisePercent       *float64 `json:"RevenueSurprisePercent"`
	ForecastSource               *string  `json:"ForecastSource"`
	EpsForecast                  *float64 `json:"EpsForecast"`
	RevenueForecast              *float64 `json:"RevenueForecast"`
	EarningReleaseDate           *string  `json:"EarningReleaseDate"`
	EPSGAAPConsensusMedian       *float64 `json:"EPSGAAPConsensusMedian"`
	EPSNormalizedConsensusMedian *float64 `json:"EPSNormalizedConsensusMedian"`
	CiqFiscalPeriodType          *string  `json:"CiqFiscalPeriodType"`
	CalendarPeriodType           *string  `json:"CalendarPeriodType"`
	CalendarPeriodStartDate      *string  `json:"CalendarPeriodStartDate"`
	CalendarPeriodEndDate        *string  `json:"CalendarPeriodEndDate"`
	PrimaryEPS                   *string  `json:"PrimaryEPS"`
}

func (e EarningsResponse) ToQuarterlyHistoryRecords() ([]models.StockEarningQuarterlyHistoryRecord, error) {
	if len(e.Data.History.Quarterly) == 0 {
		return []models.StockEarningQuarterlyHistoryRecord{}, nil
	}

	keys := make([]string, 0, len(e.Data.History.Quarterly))
	for period := range e.Data.History.Quarterly {
		keys = append(keys, period)
	}
	sort.Strings(keys)

	records := make([]models.StockEarningQuarterlyHistoryRecord, 0, len(keys))
	for _, period := range keys {
		quarter := e.Data.History.Quarterly[period]

		releaseDate, err := helper.ParseRFC3339Pointer(quarter.EarningReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("invalid earning release date for period %s: %w", period, err)
		}
		startDate, err := helper.ParseRFC3339Pointer(quarter.CalendarPeriodStartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid calendar period start date for period %s: %w", period, err)
		}
		endDate, err := helper.ParseRFC3339Pointer(quarter.CalendarPeriodEndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid calendar period end date for period %s: %w", period, err)
		}

		record := models.StockEarningQuarterlyHistoryRecord{
			Symbol:                       e.Symbol,
			SecID:                        helper.StringPointerOrNil(e.SecID),
			InstrumentID:                 helper.StringPointerOrNil(e.Data.InstrumentID),
			PeriodCode:                   period,
			EpsActual:                    quarter.EpsActual,
			EpsSurprise:                  quarter.EpsSurprise,
			EpsSurprisePercent:           quarter.EpsSurprisePercent,
			RevenueActual:                quarter.RevenueActual,
			RevenueSurprise:              quarter.RevenueSurprise,
			RevenueSurprisePercent:       quarter.RevenueSurprisePercent,
			ForecastSource:               quarter.ForecastSource,
			EpsForecast:                  quarter.EpsForecast,
			RevenueForecast:              quarter.RevenueForecast,
			EarningReleaseDate:           releaseDate,
			EPSGAAPConsensusMedian:       quarter.EPSGAAPConsensusMedian,
			EPSNormalizedConsensusMedian: quarter.EPSNormalizedConsensusMedian,
			CiqFiscalPeriodType:          quarter.CiqFiscalPeriodType,
			CalendarPeriodType:           quarter.CalendarPeriodType,
			CalendarPeriodStartDate:      startDate,
			CalendarPeriodEndDate:        endDate,
			PrimaryEPS:                   quarter.PrimaryEPS,
		}

		records = append(records, record)
	}

	return records, nil
}

func (e EarningsResponse) ToOverviewMetricsRecord() (*models.StockOverviewMetricsRecord, error) {
	nextExpectedReportDate, err := helper.ParseRFC3339Pointer(e.Data.ExpectedReportDate)
	if err != nil {
		return nil, fmt.Errorf("invalid expected report date: %w", err)
	}

	sourceTimeLastUpdated, err := helper.ParseRFC3339Pointer(e.Data.TimeLastUpdated)
	if err != nil {
		return nil, fmt.Errorf("invalid source time last updated: %w", err)
	}

	record := &models.StockOverviewMetricsRecord{
		Symbol:                   e.Symbol,
		SecID:                    helper.StringPointerOrNil(e.SecID),
		InstrumentID:             helper.StringPointerOrNil(e.Data.InstrumentID),
		Market:                   helper.StringPointerOrNil(e.Data.Market),
		Currency:                 helper.StringPointerOrNil(e.Data.Currency),
		MarketCap:                e.Data.MarketCap,
		LastActualPeriodCode:     helper.StringPointerOrNil(e.Data.LastActualFiscalPeriod),
		LastActualQuarterEPS:     e.Data.LastActual.EpsActual,
		LastActualQuarterRevenue: e.Data.LastActual.RevenueActual,
		NextExpectedReportDate:   nextExpectedReportDate,
		SourceTimeLastUpdated:    sourceTimeLastUpdated,
	}

	return record, nil
}

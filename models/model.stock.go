package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type StockInformation struct {
	Ticker string `json:"ticker" db:"ticker"`
	ApiKey string `json:"api_key" db:"api_key"`
}

type StockEarningQuarterlyHistoryRecord struct {
	ID                           int64      `json:"id" db:"id"`
	Symbol                       string     `json:"symbol" db:"symbol"`
	PeriodCode                   string     `json:"period_code" db:"period_code"`
	EpsActual                    *float64   `json:"eps_actual" db:"eps_actual"`
	EpsSurprise                  *float64   `json:"eps_surprise" db:"eps_surprise"`
	EpsSurprisePercent           *float64   `json:"eps_surprise_percent" db:"eps_surprise_percent"`
	RevenueActual                *float64   `json:"revenue_actual" db:"revenue_actual"`
	RevenueSurprise              *float64   `json:"revenue_surprise" db:"revenue_surprise"`
	RevenueSurprisePercent       *float64   `json:"revenue_surprise_percent" db:"revenue_surprise_percent"`
	ForecastSource               *string    `json:"forecast_source" db:"forecast_source"`
	EpsForecast                  *float64   `json:"eps_forecast" db:"eps_forecast"`
	RevenueForecast              *float64   `json:"revenue_forecast" db:"revenue_forecast"`
	EarningReleaseDate           *time.Time `json:"earning_release_date" db:"earning_release_date"`
	EPSGAAPConsensusMedian       *float64   `json:"eps_gaap_consensus_median" db:"eps_gaap_consensus_median"`
	EPSNormalizedConsensusMedian *float64   `json:"eps_normalized_consensus_median" db:"eps_normalized_consensus_median"`
	CiqFiscalPeriodType          *string    `json:"ciq_fiscal_period_type" db:"ciq_fiscal_period_type"`
	CalendarPeriodType           *string    `json:"calendar_period_type" db:"calendar_period_type"`
	CalendarPeriodStartDate      *time.Time `json:"calendar_period_start_date" db:"calendar_period_start_date"`
	CalendarPeriodEndDate        *time.Time `json:"calendar_period_end_date" db:"calendar_period_end_date"`
	PrimaryEPS                   *string    `json:"primary_eps" db:"primary_eps"`
	CreatedAt                    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                    time.Time  `json:"updated_at" db:"updated_at"`
}

type StockOverviewMetricsRecord struct {
	Symbol                        string
	Beta                          *float64
	Eps                           *float64
	BookValuePerShare             *float64
	LatestRevenuePerShare         *float64
	Profitability                 *string
	StockGrowth                   *float64
	LatestRevenue                 *float64
	LatestIncome                  *float64
	LatestNetProfitMargin         *float64
	Assets                        *float64
	Liabilities                   *float64
	DebtToEquityRatio             *float64
	CurrentRatio                  *float64
	QuickRatio                    *float64
	LeverageRatio                 *float64
	DebtAssetRatio                *float64
	InterestCoverage              *float64
	ReturnOnAssets                *float64
	ReturnOnEquity                *float64
	ReturnOnCapital               *float64
	RoaTTM                        *float64
	GrossMargin                   *float64
	OperatingMargin               *float64
	PretaxMargin                  *float64
	NetProfitMargin               *float64
	NetMarginPercent              *float64
	AverageGrossMargin5Y          *float64
	AveragePretaxMargin5Y         *float64
	AverageNetProfitMargin5Y      *float64
	ReturnOnAssets5YAvg           *float64
	ReturnOnEquity5YAvg           *float64
	ReturnOnCapital5YAvg          *float64
	RevenueQQLastYearGrowthRate   *float64
	NetIncomeQQLastYearGrowthRate *float64
	RevenueYTDYTD                 *float64
	NetIncomeYTDYTDGrowthRate     *float64
	Revenue3YAvg                  *float64
	Revenue5YAvgGrowthRate        *float64
	NetIncome5YAvgGrowthRate      *float64
	DilutedEPS3YGrowth            *float64
	PE5YHighRatio                 *float64
	PE5YLowRatio                  *float64
	TrailingAnnualDividendYield   *float64
	Dividend5YAvgGrowthRate       *float64
	OperatingCashFlow             *float64
	IncomeEmployee                *float64
	RevenueEmployee               *float64
	AssetTurnover                 *float64
	InventoryTurnover             *float64
	ReceivableTurnover            *float64
	PriceCashFlowRatio            *float64
	PEGrowthRatio                 *float64
	PayoutRatio                   *float64
	BookValueShareRatio           *float64
	Current                       *float64
	MarketCap                     *float64
	EnterpriseValue               *float64
	SharesOutstanding             *int64
	AverageDividendYield5Y        *float64
	LastSplitFactor               *string
	LastSplitDate                 *time.Time
	DeclarationDate               *time.Time
	DividendDate                  *time.Time
	ExDividendDate                *time.Time
	ExDividendAmount              *float64
	PriceToBookRatio              *float64
	PriceToSalesRatio             *float64
	ForwardPriceToEPS             *float64
	ForwardDividendYield          *float64
	DividendYield                 *float64
	LastActualPeriodCode          *string
	LastActualQuarterEPS          *float64
	LastActualQuarterRevenue      *float64
	NextExpectedReportDate        *time.Time
	SourceTimeLastUpdated         *time.Time
}

// StockRepository defines operations for stocks.
type StockRepository interface {
	GetStockApiKey() ([]StockInformation, error)
	UpsertStockEarningQuarterlyHistory(records []StockEarningQuarterlyHistoryRecord) error
	UpsertStockOverviewMetrics(record *StockOverviewMetricsRecord) error
}

type stockRepository struct{}

func NewStockRepository() StockRepository {
	return &stockRepository{}
}

func (r *stockRepository) getDB() (*sqlx.DB, error) {
	db := GetDB().PostgreDBManager.RW
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return db, nil
}

func (r *stockRepository) GetStockApiKey() ([]StockInformation, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	const query = `SELECT a.ticker, a.api_key FROM stock a WHERE a.last_update < NOW() - INTERVAL '14 days' AND a.api_key IS NOT NULL`

	var stocks []StockInformation
	err = db.Select(&stocks, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching stock API keys: %w", err)
	}

	return stocks, nil
}

func (r *stockRepository) UpsertStockEarningQuarterlyHistory(records []StockEarningQuarterlyHistoryRecord) error {
	if len(records) == 0 {
		return nil
	}

	db, err := r.getDB()
	if err != nil {
		return err
	}

	const updateTimestampQuery = `UPDATE stock SET last_update = CURRENT_TIMESTAMP WHERE ticker = $1`
	const query = `
		INSERT INTO stock_earning_quarterly_history (
			symbol,
			period_code,
			eps_actual,
			eps_surprise,
			eps_surprise_percent,
			revenue_actual,
			revenue_surprise,
			revenue_surprise_percent,
			forecast_source,
			eps_forecast,
			revenue_forecast,
			earning_release_date,
			eps_gaap_consensus_median,
			eps_normalized_consensus_median,
			ciq_fiscal_period_type,
			calendar_period_type,
			calendar_period_start_date,
			calendar_period_end_date,
			primary_eps
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13,
			$14,
			$15,
			$16,
			$17,
			$18,
			$19
		)
		ON CONFLICT (symbol, period_code) 
		DO UPDATE SET 
			eps_actual = EXCLUDED.eps_actual,
			eps_surprise = EXCLUDED.eps_surprise,
			eps_surprise_percent = EXCLUDED.eps_surprise_percent,
			revenue_actual = EXCLUDED.revenue_actual,
			revenue_surprise = EXCLUDED.revenue_surprise,
			revenue_surprise_percent = EXCLUDED.revenue_surprise_percent,
			forecast_source = EXCLUDED.forecast_source,
			eps_forecast = EXCLUDED.eps_forecast,
			revenue_forecast = EXCLUDED.revenue_forecast,
			earning_release_date = EXCLUDED.earning_release_date,
			eps_gaap_consensus_median = EXCLUDED.eps_gaap_consensus_median,
			eps_normalized_consensus_median = EXCLUDED.eps_normalized_consensus_median,
			ciq_fiscal_period_type = EXCLUDED.ciq_fiscal_period_type,
			calendar_period_type = EXCLUDED.calendar_period_type,
			calendar_period_start_date = EXCLUDED.calendar_period_start_date,
			calendar_period_end_date = EXCLUDED.calendar_period_end_date,
			primary_eps = EXCLUDED.primary_eps,
			updated_at = CURRENT_TIMESTAMP
	`

	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("error starting transaction for stock earning quarterly upsert: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex(query)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error preparing stock earning quarterly upsert statement: %w", err)
	}
	defer stmt.Close()

	stmtUpdate, err := tx.Preparex(updateTimestampQuery)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error preparing stock last update timestamp statement: %w", err)
	}
	defer stmtUpdate.Close()

	for _, record := range records {
		if _, err := stmtUpdate.Exec(record.Symbol); err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating last update timestamp for symbol %s: %w", record.Symbol, err)
		}

		if _, err := stmt.Exec(
			record.Symbol,
			record.PeriodCode,
			record.EpsActual,
			record.EpsSurprise,
			record.EpsSurprisePercent,
			record.RevenueActual,
			record.RevenueSurprise,
			record.RevenueSurprisePercent,
			record.ForecastSource,
			record.EpsForecast,
			record.RevenueForecast,
			record.EarningReleaseDate,
			record.EPSGAAPConsensusMedian,
			record.EPSNormalizedConsensusMedian,
			record.CiqFiscalPeriodType,
			record.CalendarPeriodType,
			record.CalendarPeriodStartDate,
			record.CalendarPeriodEndDate,
			record.PrimaryEPS,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("error upserting stock earning quarterly history for symbol %s period %s: %w", record.Symbol, record.PeriodCode, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing stock earning quarterly upsert transaction: %w", err)
	}

	return nil
}

func (r *stockRepository) UpsertStockOverviewMetrics(record *StockOverviewMetricsRecord) error {
	if record == nil {
		return nil
	}

	db, err := r.getDB()
	if err != nil {
		return err
	}

	const updateTimestampQuery = `UPDATE stock SET last_update = CURRENT_TIMESTAMP WHERE ticker = $1`
	const query = `
		INSERT INTO stock_overview_metrics (
			symbol,
			beta,
			eps,
			book_value_per_share,
			latest_revenue_per_share,
			profitability,
			stock_growth,
			latest_revenue,
			latest_income,
			latest_net_profit_margin,
			assets,
			liabilities,
			debt_to_equity_ratio,
			current_ratio,
			quick_ratio,
			leverage_ratio,
			debt_asset_ratio,
			interest_coverage,
			return_on_assets,
			return_on_equity,
			return_on_capital,
			roa_ttm,
			gross_margin,
			operating_margin,
			pretax_margin,
			net_profit_margin,
			net_margin_percent,
			average_gross_margin_5y,
			average_pretax_margin_5y,
			average_net_profit_margin_5y,
			return_on_assets_5y_avg,
			return_on_equity_5y_avg,
			return_on_capital_5y_avg,
			revenue_qq_last_year_growth_rate,
			net_income_qq_last_year_growth_rate,
			revenue_ytdytd,
			net_income_ytdytd_growth_rate,
			revenue_3y_avg,
			revenue_5y_avg_growth_rate,
			net_income_5y_avg_growth_rate,
			diluted_eps_3y_growth,
			pe_5y_high_ratio,
			pe_5y_low_ratio,
			trailing_annual_dividend_yield,
			dividend_5y_avg_growth_rate,
			operating_cash_flow,
			income_employee,
			revenue_employee,
			asset_turnover,
			inventory_turnover,
			receivable_turnover,
			price_cash_flow_ratio,
			peg_growth_ratio,
			payout_ratio,
			book_value_share_ratio,
			current,
			market_cap,
			enterprise_value,
			shares_outstanding,
			average_dividend_yield_5y,
			last_split_factor,
			last_split_date,
			declaration_date,
			dividend_date,
			ex_dividend_date,
			ex_dividend_amount,
			price_to_book_ratio,
			price_to_sales_ratio,
			forward_price_to_eps,
			forward_dividend_yield,
			dividend_yield,
			last_actual_period_code,
			last_actual_quarter_eps,
			last_actual_quarter_revenue,
			next_expected_report_date,
			source_time_last_updated
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
			$31, $32, $33, $34, $35, $36, $37, $38, $39, $40,
			$41, $42, $43, $44, $45, $46, $47, $48, $49, $50,
			$51, $52, $53, $54, $55, $56, $57, $58, $59, $60,
			$61, $62, $63, $64, $65, $66, $67, $68, $69, $70,
			$71, $72, $73, $74, $75, $76
		)
		ON CONFLICT (symbol) 
		DO UPDATE SET 
			beta = EXCLUDED.beta,
			eps = EXCLUDED.eps,
			book_value_per_share = EXCLUDED.book_value_per_share,
			latest_revenue_per_share = EXCLUDED.latest_revenue_per_share,
			profitability = EXCLUDED.profitability,
			stock_growth = EXCLUDED.stock_growth,
			latest_revenue = EXCLUDED.latest_revenue,
			latest_income = EXCLUDED.latest_income,
			latest_net_profit_margin = EXCLUDED.latest_net_profit_margin,
			assets = EXCLUDED.assets,
			liabilities = EXCLUDED.liabilities,
			debt_to_equity_ratio = EXCLUDED.debt_to_equity_ratio,
			current_ratio = EXCLUDED.current_ratio,
			quick_ratio = EXCLUDED.quick_ratio,
			leverage_ratio = EXCLUDED.leverage_ratio,
			debt_asset_ratio = EXCLUDED.debt_asset_ratio,
			interest_coverage = EXCLUDED.interest_coverage,
			return_on_assets = EXCLUDED.return_on_assets,
			return_on_equity = EXCLUDED.return_on_equity,
			return_on_capital = EXCLUDED.return_on_capital,
			roa_ttm = EXCLUDED.roa_ttm,
			gross_margin = EXCLUDED.gross_margin,
			operating_margin = EXCLUDED.operating_margin,
			pretax_margin = EXCLUDED.pretax_margin,
			net_profit_margin = EXCLUDED.net_profit_margin,
			net_margin_percent = EXCLUDED.net_margin_percent,
			average_gross_margin_5y = EXCLUDED.average_gross_margin_5y,
			average_pretax_margin_5y = EXCLUDED.average_pretax_margin_5y,
			average_net_profit_margin_5y = EXCLUDED.average_net_profit_margin_5y,
			return_on_assets_5y_avg = EXCLUDED.return_on_assets_5y_avg,
			return_on_equity_5y_avg = EXCLUDED.return_on_equity_5y_avg,
			return_on_capital_5y_avg = EXCLUDED.return_on_capital_5y_avg,
			revenue_qq_last_year_growth_rate = EXCLUDED.revenue_qq_last_year_growth_rate,
			net_income_qq_last_year_growth_rate = EXCLUDED.net_income_qq_last_year_growth_rate,
			revenue_ytdytd = EXCLUDED.revenue_ytdytd,
			net_income_ytdytd_growth_rate = EXCLUDED.net_income_ytdytd_growth_rate,
			revenue_3y_avg = EXCLUDED.revenue_3y_avg,
			revenue_5y_avg_growth_rate = EXCLUDED.revenue_5y_avg_growth_rate,
			net_income_5y_avg_growth_rate = EXCLUDED.net_income_5y_avg_growth_rate,
			diluted_eps_3y_growth = EXCLUDED.diluted_eps_3y_growth,
			pe_5y_high_ratio = EXCLUDED.pe_5y_high_ratio,
			pe_5y_low_ratio = EXCLUDED.pe_5y_low_ratio,
			trailing_annual_dividend_yield = EXCLUDED.trailing_annual_dividend_yield,
			dividend_5y_avg_growth_rate = EXCLUDED.dividend_5y_avg_growth_rate,
			operating_cash_flow = EXCLUDED.operating_cash_flow,
			income_employee = EXCLUDED.income_employee,
			revenue_employee = EXCLUDED.revenue_employee,
			asset_turnover = EXCLUDED.asset_turnover,
			inventory_turnover = EXCLUDED.inventory_turnover,
			receivable_turnover = EXCLUDED.receivable_turnover,
			price_cash_flow_ratio = EXCLUDED.price_cash_flow_ratio,
			peg_growth_ratio = EXCLUDED.peg_growth_ratio,
			payout_ratio = EXCLUDED.payout_ratio,
			book_value_share_ratio = EXCLUDED.book_value_share_ratio,
			current = EXCLUDED.current,
			market_cap = EXCLUDED.market_cap,
			enterprise_value = EXCLUDED.enterprise_value,
			shares_outstanding = EXCLUDED.shares_outstanding,
			average_dividend_yield_5y = EXCLUDED.average_dividend_yield_5y,
			last_split_factor = EXCLUDED.last_split_factor,
			last_split_date = EXCLUDED.last_split_date,
			declaration_date = EXCLUDED.declaration_date,
			dividend_date = EXCLUDED.dividend_date,
			ex_dividend_date = EXCLUDED.ex_dividend_date,
			ex_dividend_amount = EXCLUDED.ex_dividend_amount,
			price_to_book_ratio = EXCLUDED.price_to_book_ratio,
			price_to_sales_ratio = EXCLUDED.price_to_sales_ratio,
			forward_price_to_eps = EXCLUDED.forward_price_to_eps,
			forward_dividend_yield = EXCLUDED.forward_dividend_yield,
			dividend_yield = EXCLUDED.dividend_yield,
			last_actual_period_code = EXCLUDED.last_actual_period_code,
			last_actual_quarter_eps = EXCLUDED.last_actual_quarter_eps,
			last_actual_quarter_revenue = EXCLUDED.last_actual_quarter_revenue,
			next_expected_report_date = EXCLUDED.next_expected_report_date,
			source_time_last_updated = EXCLUDED.source_time_last_updated
	`

	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("error starting transaction for stock overview metrics upsert: %w", err)
	}

	stmt, err := tx.Preparex(query)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error preparing stock earning quarterly upsert statement: %w", err)
	}
	defer stmt.Close()

	stmtUpdate, err := tx.Preparex(updateTimestampQuery)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error preparing stock last update timestamp statement: %w", err)
	}
	defer stmtUpdate.Close()

	if _, err := stmt.Exec(
		&record.Symbol,
		&record.Beta,
		&record.Eps,
		&record.BookValuePerShare,
		&record.LatestRevenuePerShare,
		&record.Profitability,
		&record.StockGrowth,
		&record.LatestRevenue,
		&record.LatestIncome,
		&record.LatestNetProfitMargin,
		&record.Assets,
		&record.Liabilities,
		&record.DebtToEquityRatio,
		&record.CurrentRatio,
		&record.QuickRatio,
		&record.LeverageRatio,
		&record.DebtAssetRatio,
		&record.InterestCoverage,
		&record.ReturnOnAssets,
		&record.ReturnOnEquity,
		&record.ReturnOnCapital,
		&record.RoaTTM,
		&record.GrossMargin,
		&record.OperatingMargin,
		&record.PretaxMargin,
		&record.NetProfitMargin,
		&record.NetMarginPercent,
		&record.AverageGrossMargin5Y,
		&record.AveragePretaxMargin5Y,
		&record.AverageNetProfitMargin5Y,
		&record.ReturnOnAssets5YAvg,
		&record.ReturnOnEquity5YAvg,
		&record.ReturnOnCapital5YAvg,
		&record.RevenueQQLastYearGrowthRate,
		&record.NetIncomeQQLastYearGrowthRate,
		&record.RevenueYTDYTD,
		&record.NetIncomeYTDYTDGrowthRate,
		&record.Revenue3YAvg,
		&record.Revenue5YAvgGrowthRate,
		&record.NetIncome5YAvgGrowthRate,
		&record.DilutedEPS3YGrowth,
		&record.PE5YHighRatio,
		&record.PE5YLowRatio,
		&record.TrailingAnnualDividendYield,
		&record.Dividend5YAvgGrowthRate,
		&record.OperatingCashFlow,
		&record.IncomeEmployee,
		&record.RevenueEmployee,
		&record.AssetTurnover,
		&record.InventoryTurnover,
		&record.ReceivableTurnover,
		&record.PriceCashFlowRatio,
		&record.PEGrowthRatio,
		&record.PayoutRatio,
		&record.BookValueShareRatio,
		&record.Current,
		&record.MarketCap,
		&record.EnterpriseValue,
		&record.SharesOutstanding,
		&record.AverageDividendYield5Y,
		&record.LastSplitFactor,
		&record.LastSplitDate,
		&record.DeclarationDate,
		&record.DividendDate,
		&record.ExDividendDate,
		&record.ExDividendAmount,
		&record.PriceToBookRatio,
		&record.PriceToSalesRatio,
		&record.ForwardPriceToEPS,
		&record.ForwardDividendYield,
		&record.DividendYield,
		&record.LastActualPeriodCode,
		&record.LastActualQuarterEPS,
		&record.LastActualQuarterRevenue,
		&record.NextExpectedReportDate,
		&record.SourceTimeLastUpdated,
	); err != nil {
		tx.Rollback()
		Logger.Error().Err(err).Msgf("Error upserting stock overview metrics : %+v", record)
		return fmt.Errorf("error upserting stock overview metrics for symbol %s: %w", record.Symbol, err)
	}

	if _, err := stmtUpdate.Exec(record.Symbol); err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating last update timestamp for symbol %s: %w", record.Symbol, err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("error committing transaction for symbol %s: %w", record.Symbol, err)
	}

	tx.Commit()
	return nil
}

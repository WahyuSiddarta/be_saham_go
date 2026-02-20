-- Quarterly-only DDL for datasource/earning.json and datasource/equities.json
-- Fresh schema (ignores previous DDL drafts)

-- public.stock definition

-- Drop table

-- DROP TABLE public.stock;

CREATE TABLE public.stock (
	ticker bpchar(4) NOT NULL,
	name varchar NULL,
	description varchar NULL,
	icon varchar NULL,
	establish_year int4 NULL,
	street varchar NULL,
	city varchar NULL,
	state varchar NULL,
	zip varchar NULL,
	phone varchar NULL,
	industry public.stock_industry NULL,
	sector public.stock_sector NULL,
	last_update timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
  api_key varchar NULL,
	CONSTRAINT stock_pk PRIMARY KEY (ticker)
);

-- ============================================================================
-- EARNING JSON (Quarterly)
-- Source path:
-- - data.History.quarterly
-- ============================================================================


CREATE TABLE stock_earning_quarterly_history (
    id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR(32) NOT NULL,
    sec_id VARCHAR(64),
    instrument_id VARCHAR(64),
    period_code CHAR(6) NOT NULL CHECK (period_code ~ '^[0-9]{6}$'),

    eps_actual NUMERIC(20, 6),
    eps_surprise NUMERIC(20, 6),
    eps_surprise_percent NUMERIC(20, 6),

    revenue_actual NUMERIC(20, 2),
    revenue_surprise NUMERIC(20, 2),
    revenue_surprise_percent NUMERIC(20, 6),

    forecast_source VARCHAR(64),
    eps_forecast NUMERIC(20, 6),
    revenue_forecast NUMERIC(20, 2),
    earning_release_date TIMESTAMP WITH TIME ZONE,
    eps_gaap_consensus_median NUMERIC(20, 6),
    eps_normalized_consensus_median NUMERIC(20, 6),

    ciq_fiscal_period_type VARCHAR(32),
    calendar_period_type VARCHAR(32),
    calendar_period_start_date TIMESTAMP WITH TIME ZONE,
    calendar_period_end_date TIMESTAMP WITH TIME ZONE,
    primary_eps VARCHAR(32),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE stock_earning_quarterly_history
    ADD CONSTRAINT uq_stock_earning_qh_symbol_period UNIQUE (symbol, period_code);

CREATE INDEX idx_stock_earning_qh_symbol ON stock_earning_quarterly_history(symbol);
CREATE INDEX idx_stock_earning_qh_period_code ON stock_earning_quarterly_history(period_code);
CREATE INDEX idx_stock_earning_qh_release_date ON stock_earning_quarterly_history(earning_release_date);

-- ============================================================================
-- OVERVIEW METRICS (EPS, asset, liabilities, revenue, etc.)
-- Source paths:
-- - equities.json: data.analysis.keyMetrics, data.analysis.companyMetrics,
--   data.analysis.shareStatistics, data.marketCap / data.enterpriseValue,
--   data.beta,
--   latest annual statement values from data.analysis.annualStatements
-- - earning.json: data.LastActual, data.ExpectedReportDate
-- ============================================================================

CREATE TABLE stock_overview_metrics (
    id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR(32) NOT NULL,
    sec_id VARCHAR(64),
    instrument_id VARCHAR(64),
    market VARCHAR(32),
    currency VARCHAR(16),
    beta NUMERIC(30, 6),

    eps NUMERIC(30, 6),
    book_value_per_share NUMERIC(30, 6),
    latest_revenue_per_share NUMERIC(30, 6),
    profitability VARCHAR(255),
    stock_growth NUMERIC(30, 6),
    latest_revenue NUMERIC(30, 2),
    latest_income NUMERIC(30, 2),
    latest_net_profit_margin NUMERIC(30, 6),
    assets NUMERIC(30, 2),
    liabilities NUMERIC(30, 2),
    debt_to_equity_ratio NUMERIC(30, 6),
    current_ratio NUMERIC(30, 6),
    quick_ratio NUMERIC(30, 6),
    leverage_ratio NUMERIC(30, 6),
    debt_asset_ratio NUMERIC(30, 6),
    interest_coverage NUMERIC(30, 6),
    return_on_assets NUMERIC(30, 6),
    return_on_equity NUMERIC(30, 6),
    return_on_capital NUMERIC(30, 6),
    roa_ttm NUMERIC(30, 6),

    gross_margin NUMERIC(30, 6),
    operating_margin NUMERIC(30, 6),
    pretax_margin NUMERIC(30, 6),
    net_profit_margin NUMERIC(30, 6),
    net_margin_percent NUMERIC(30, 6),

    average_gross_margin_5y NUMERIC(30, 6),
    average_pretax_margin_5y NUMERIC(30, 6),
    average_net_profit_margin_5y NUMERIC(30, 6),
    return_on_assets_5y_avg NUMERIC(30, 6),
    return_on_equity_5y_avg NUMERIC(30, 6),
    return_on_capital_5y_avg NUMERIC(30, 6),

    revenue_qq_last_year_growth_rate NUMERIC(30, 6),
    net_income_qq_last_year_growth_rate NUMERIC(30, 6),
    revenue_ytdytd NUMERIC(30, 6),
    net_income_ytdytd_growth_rate NUMERIC(30, 6),
    revenue_3y_avg NUMERIC(30, 6),
    revenue_5y_avg_growth_rate NUMERIC(30, 6),
    net_income_5y_avg_growth_rate NUMERIC(30, 6),
    diluted_eps_3y_growth NUMERIC(30, 6),
    pe_5y_high_ratio NUMERIC(30, 6),
    pe_5y_low_ratio NUMERIC(30, 6),
    trailing_annual_dividend_yield NUMERIC(30, 6),
    dividend_5y_avg_growth_rate NUMERIC(30, 6),

    operating_cash_flow NUMERIC(30, 2),
    income_employee NUMERIC(30, 6),
    revenue_employee NUMERIC(30, 6),
    asset_turnover NUMERIC(30, 6),
    inventory_turnover NUMERIC(30, 6),
    receivable_turnover NUMERIC(30, 6),
    price_cash_flow_ratio NUMERIC(30, 6),
    peg_growth_ratio NUMERIC(30, 6),
    payout_ratio NUMERIC(30, 6),
    book_value_share_ratio NUMERIC(30, 6),
    current NUMERIC(30, 6),

    market_cap NUMERIC(30, 2),
    enterprise_value NUMERIC(30, 2),
    shares_outstanding BIGINT,
    average_dividend_yield_5y NUMERIC(30, 6),
    last_split_factor VARCHAR(64),
    last_split_date TIMESTAMP WITH TIME ZONE,
    declaration_date TIMESTAMP WITH TIME ZONE,
    dividend_date TIMESTAMP WITH TIME ZONE,
    ex_dividend_date TIMESTAMP WITH TIME ZONE,
    ex_dividend_amount NUMERIC(30, 6),
    price_to_book_ratio NUMERIC(30, 6),
    price_to_sales_ratio NUMERIC(30, 6),
    forward_price_to_eps NUMERIC(30, 6),
    forward_dividend_yield NUMERIC(30, 6),
    dividend_yield NUMERIC(30, 6),

    last_actual_period_code CHAR(6) CHECK (last_actual_period_code ~ '^[0-9]{6}$'),
    last_actual_quarter_eps NUMERIC(30, 6),
    last_actual_quarter_revenue NUMERIC(30, 2),
    next_expected_report_date TIMESTAMP WITH TIME ZONE,

    source_time_last_updated TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE stock_overview_metrics
    ADD CONSTRAINT uq_stock_overview_metrics_symbol UNIQUE (symbol);

CREATE INDEX idx_stock_overview_symbol ON stock_overview_metrics(symbol);
CREATE INDEX idx_stock_overview_source_time ON stock_overview_metrics(source_time_last_updated);

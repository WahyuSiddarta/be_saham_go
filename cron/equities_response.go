package cron

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/WahyuSiddarta/be_saham_go/models"
)

type EquitiesResponse struct {
	Success bool         `json:"success"`
	Symbol  string       `json:"symbol"`
	SecID   string       `json:"secId"`
	Data    EquitiesData `json:"data"`
}

type EquitiesData struct {
	InstrumentID    string           `json:"instrumentId"`
	Market          string           `json:"market"`
	Currency        string           `json:"currency"`
	MarketCap       *float64         `json:"marketCap"`
	EnterpriseValue *float64         `json:"enterpriseValue"`
	Beta            *float64         `json:"beta"`
	TimeLastUpdated *string          `json:"timeLastUpdated"`
	Analysis        EquitiesAnalysis `json:"analysis"`
}

type EquitiesAnalysis struct {
	AnnualStatements map[string]EquitiesAnnualStatement `json:"annualStatements"`
	KeyMetrics       EquitiesKeyMetrics                 `json:"keyMetrics"`
	CompanyMetrics   EquitiesCompanyMetrics             `json:"companyMetrics"`
	ShareStatistics  EquitiesShareStatistics            `json:"shareStatistics"`
}

type EquitiesAnnualStatement struct {
	Assets      *float64 `json:"assets"`
	Liabilities *float64 `json:"liabilities"`
}

type EquitiesKeyMetrics struct {
	Eps                   *float64 `json:"eps"`
	BookValuePerShare     *float64 `json:"bookValuePerShare"`
	LatestRevenuePerShare *float64 `json:"latestRevenuePerShare"`
	Profitability         *string  `json:"profitability"`
	StockGrowth           *float64 `json:"stockGrowth"`
	LatestRevenue         *float64 `json:"latestRevenue"`
	LatestIncome          *float64 `json:"latestIncome"`
	LatestNetProfitMargin *float64 `json:"latestNetProfitMargin"`
	CurrentRatio          *float64 `json:"currentRatio"`
	DebtToEquityRatio     *float64 `json:"debtToEquityRatio"`
	ForwardPriceToEPS     *float64 `json:"forwardPriceToEPS"`
	ForwardDividendYield  *float64 `json:"forwardDividendYield"`
	PayoutRatio           *float64 `json:"payoutRatio"`
	PriceToBookRatio      *float64 `json:"priceToBookRatio"`
	ReturnOnAssets        *float64 `json:"returnOnAssets"`
	ReturnOnCapital       *float64 `json:"returnOnCapital"`
	ReturnOnEquity        *float64 `json:"returnOnEquity"`
}

type EquitiesCompanyMetrics struct {
	PE5YearHighRatio                *float64 `json:"pE5YearHighRatio"`
	PE5YearLowRatio                 *float64 `json:"pE5YearLowRatio"`
	RevenueYTDYTD                   *float64 `json:"revenueYTDYTD"`
	RevenueQQLastYearGrowthRate     *float64 `json:"revenueQQLastYearGrowthRate"`
	NetIncomeYTDYTDGrowthRate       *float64 `json:"netIncomeYTDYTDGrowthRate"`
	NetIncomeQQLastYearGrowthRate   *float64 `json:"netIncomeQQLastYearGrowthRate"`
	Revenue5YearAverageGrowthRate   *float64 `json:"revenue5YearAverageGrowthRate"`
	NetIncome5YearAverageGrowthRate *float64 `json:"netIncome5YearAverageGrowthRate"`
	Dividend5YearAverageGrowthRate  *float64 `json:"dividend5YearAverageGrowthRate"`
	ForwardDividendYield            *float64 `json:"forwardDividendYield"`
	DividendYield                   *float64 `json:"dividendYield"`
	CurrentRatio                    *float64 `json:"currentRatio"`
	DebtAssetRatio                  *float64 `json:"debtAssetRatio"`
	LeverageRatio                   *float64 `json:"leverageRatio"`
	InterestCoverage                *float64 `json:"interestCoverage"`
	PriceCashFlowRatio              *float64 `json:"priceCashFlowRatio"`
	Revenue3YearAverage             *float64 `json:"revenue3YearAverage"`
	TrailingAnnualDividendYield     *float64 `json:"trailingAnnualDividendYield"`
	PriceBookRatio                  *float64 `json:"priceBookRatio"`
	PriceSalesRatio                 *float64 `json:"priceSalesRatio"`
	BookValueShareRatio             *float64 `json:"bookValueShareRatio"`
	OperatingCashFlow               *float64 `json:"operatingCashFlow"`
	PayoutRatio                     *float64 `json:"payoutRatio"`
	QuickRatio                      *float64 `json:"quickRatio"`
	Current                         *float64 `json:"current"`
	DebtEquityRatio                 *float64 `json:"debtEquityRatio"`
	DilutedEPS3YearGrowth           *float64 `json:"dilutedEPS3YearGrowth"`
	PEGrowthRatio                   *float64 `json:"pEGrowthRatio"`
	GrossMargin                     *float64 `json:"grossMargin"`
	PreTaxMargin                    *float64 `json:"preTaxMargin"`
	NetProfitMargin                 *float64 `json:"netProfitMargin"`
	AverageGrossMargin5Year         *float64 `json:"averageGrossMargin5Year"`
	AveragePreTaxMargin5Year        *float64 `json:"averagePreTaxMargin5Year"`
	AverageNetProfitMargin5Year     *float64 `json:"averageNetProfitMargin5Year"`
	OperatingMargin                 *float64 `json:"operatingMargin"`
	NetMarginPercent                *float64 `json:"netMarginPercent"`
	ReturnOnEquityCurrent           *float64 `json:"returnOnEquityCurrent"`
	ReturnOnEquity5YearAverage      *float64 `json:"returnOnEquity5YearAverage"`
	ReturnOnAssetCurrent            *float64 `json:"returnOnAssetCurrent"`
	ReturnOnAsset5YearAverage       *float64 `json:"returnOnAsset5YearAverage"`
	ReturnOnCapitalCurrent          *float64 `json:"returnOnCapitalCurrent"`
	ReturnOnCapital5YearAverage     *float64 `json:"returnOnCapital5YearAverage"`
	IncomeEmployee                  *float64 `json:"incomeEmployee"`
	RevenueEmployee                 *float64 `json:"revenueEmployee"`
	AssetTurnover                   *float64 `json:"assetTurnover"`
	InventoryTurnover               *float64 `json:"inventoryTurnover"`
	ReceivableTurnover              *float64 `json:"receivableTurnover"`
	RoaTTM                          *float64 `json:"roaTTM"`
}

type EquitiesShareStatistics struct {
	AverageDividendYield5Year *float64 `json:"averageDividendYield5Year"`
	LastSplitFactor           *string  `json:"lastSplitFactor"`
	LastSplitDate             *string  `json:"lastSplitDate"`
	DeclarationDate           *string  `json:"declarationDate"`
	DividendDate              *string  `json:"dividendDate"`
	ExDividendDate            *string  `json:"exDividendDate"`
	ExDividendAmount          *float64 `json:"exDividendAmount"`
	SharesOutstanding         *int64   `json:"sharesOutstanding"`
	EnterpriseValue           *float64 `json:"enterpriseValue"`
	DividendYield             *float64 `json:"dividendYield"`
}

func (e EquitiesResponse) MergeIntoOverviewMetricsRecord(record *models.StockOverviewMetricsRecord) error {
	if record == nil {
		return fmt.Errorf("overview metrics record is nil")
	}

	if record.Symbol == "" {
		record.Symbol = e.Symbol
	}

	setFloat(&record.Beta, e.Data.Beta)
	setFloat(&record.Eps, e.Data.Analysis.KeyMetrics.Eps)
	setFloat(&record.BookValuePerShare, e.Data.Analysis.KeyMetrics.BookValuePerShare)
	setFloat(&record.LatestRevenuePerShare, e.Data.Analysis.KeyMetrics.LatestRevenuePerShare)
	setString(&record.Profitability, e.Data.Analysis.KeyMetrics.Profitability)
	setFloat(&record.StockGrowth, e.Data.Analysis.KeyMetrics.StockGrowth)
	setFloat(&record.LatestRevenue, e.Data.Analysis.KeyMetrics.LatestRevenue)
	setFloat(&record.LatestIncome, e.Data.Analysis.KeyMetrics.LatestIncome)
	setFloat(&record.LatestNetProfitMargin, e.Data.Analysis.KeyMetrics.LatestNetProfitMargin)
	setFloat(&record.CurrentRatio, e.Data.Analysis.KeyMetrics.CurrentRatio)
	setFloat(&record.DebtToEquityRatio, e.Data.Analysis.KeyMetrics.DebtToEquityRatio)
	setFloat(&record.ForwardPriceToEPS, e.Data.Analysis.KeyMetrics.ForwardPriceToEPS)
	setFloat(&record.ForwardDividendYield, e.Data.Analysis.KeyMetrics.ForwardDividendYield)
	setFloat(&record.PayoutRatio, e.Data.Analysis.KeyMetrics.PayoutRatio)
	setFloat(&record.PriceToBookRatio, e.Data.Analysis.KeyMetrics.PriceToBookRatio)
	setFloat(&record.ReturnOnAssets, e.Data.Analysis.KeyMetrics.ReturnOnAssets)
	setFloat(&record.ReturnOnCapital, e.Data.Analysis.KeyMetrics.ReturnOnCapital)
	setFloat(&record.ReturnOnEquity, e.Data.Analysis.KeyMetrics.ReturnOnEquity)

	setFloat(&record.PE5YHighRatio, e.Data.Analysis.CompanyMetrics.PE5YearHighRatio)
	setFloat(&record.PE5YLowRatio, e.Data.Analysis.CompanyMetrics.PE5YearLowRatio)
	setFloat(&record.RevenueYTDYTD, e.Data.Analysis.CompanyMetrics.RevenueYTDYTD)
	setFloat(&record.RevenueQQLastYearGrowthRate, e.Data.Analysis.CompanyMetrics.RevenueQQLastYearGrowthRate)
	setFloat(&record.NetIncomeYTDYTDGrowthRate, e.Data.Analysis.CompanyMetrics.NetIncomeYTDYTDGrowthRate)
	setFloat(&record.NetIncomeQQLastYearGrowthRate, e.Data.Analysis.CompanyMetrics.NetIncomeQQLastYearGrowthRate)
	setFloat(&record.Revenue5YAvgGrowthRate, e.Data.Analysis.CompanyMetrics.Revenue5YearAverageGrowthRate)
	setFloat(&record.NetIncome5YAvgGrowthRate, e.Data.Analysis.CompanyMetrics.NetIncome5YearAverageGrowthRate)
	setFloat(&record.Dividend5YAvgGrowthRate, e.Data.Analysis.CompanyMetrics.Dividend5YearAverageGrowthRate)
	setFloat(&record.DividendYield, e.Data.Analysis.CompanyMetrics.DividendYield)
	setFloat(&record.CurrentRatio, e.Data.Analysis.CompanyMetrics.CurrentRatio)
	setFloat(&record.DebtAssetRatio, e.Data.Analysis.CompanyMetrics.DebtAssetRatio)
	setFloat(&record.LeverageRatio, e.Data.Analysis.CompanyMetrics.LeverageRatio)
	setFloat(&record.InterestCoverage, e.Data.Analysis.CompanyMetrics.InterestCoverage)
	setFloat(&record.PriceCashFlowRatio, e.Data.Analysis.CompanyMetrics.PriceCashFlowRatio)
	setFloat(&record.Revenue3YAvg, e.Data.Analysis.CompanyMetrics.Revenue3YearAverage)
	setFloat(&record.TrailingAnnualDividendYield, e.Data.Analysis.CompanyMetrics.TrailingAnnualDividendYield)
	setFloat(&record.PriceToBookRatio, e.Data.Analysis.CompanyMetrics.PriceBookRatio)
	setFloat(&record.PriceToSalesRatio, e.Data.Analysis.CompanyMetrics.PriceSalesRatio)
	setFloat(&record.BookValueShareRatio, e.Data.Analysis.CompanyMetrics.BookValueShareRatio)
	setFloat(&record.OperatingCashFlow, e.Data.Analysis.CompanyMetrics.OperatingCashFlow)
	setFloat(&record.PayoutRatio, e.Data.Analysis.CompanyMetrics.PayoutRatio)
	setFloat(&record.QuickRatio, e.Data.Analysis.CompanyMetrics.QuickRatio)
	setFloat(&record.Current, e.Data.Analysis.CompanyMetrics.Current)
	setFloat(&record.DebtToEquityRatio, e.Data.Analysis.CompanyMetrics.DebtEquityRatio)
	setFloat(&record.DilutedEPS3YGrowth, e.Data.Analysis.CompanyMetrics.DilutedEPS3YearGrowth)
	setFloat(&record.PEGrowthRatio, e.Data.Analysis.CompanyMetrics.PEGrowthRatio)
	setFloat(&record.GrossMargin, e.Data.Analysis.CompanyMetrics.GrossMargin)
	setFloat(&record.PretaxMargin, e.Data.Analysis.CompanyMetrics.PreTaxMargin)
	setFloat(&record.NetProfitMargin, e.Data.Analysis.CompanyMetrics.NetProfitMargin)
	setFloat(&record.AverageGrossMargin5Y, e.Data.Analysis.CompanyMetrics.AverageGrossMargin5Year)
	setFloat(&record.AveragePretaxMargin5Y, e.Data.Analysis.CompanyMetrics.AveragePreTaxMargin5Year)
	setFloat(&record.AverageNetProfitMargin5Y, e.Data.Analysis.CompanyMetrics.AverageNetProfitMargin5Year)
	setFloat(&record.OperatingMargin, e.Data.Analysis.CompanyMetrics.OperatingMargin)
	setFloat(&record.NetMarginPercent, e.Data.Analysis.CompanyMetrics.NetMarginPercent)
	setFloat(&record.ReturnOnEquity, e.Data.Analysis.CompanyMetrics.ReturnOnEquityCurrent)
	setFloat(&record.ReturnOnEquity5YAvg, e.Data.Analysis.CompanyMetrics.ReturnOnEquity5YearAverage)
	setFloat(&record.ReturnOnAssets, e.Data.Analysis.CompanyMetrics.ReturnOnAssetCurrent)
	setFloat(&record.ReturnOnAssets5YAvg, e.Data.Analysis.CompanyMetrics.ReturnOnAsset5YearAverage)
	setFloat(&record.ReturnOnCapital, e.Data.Analysis.CompanyMetrics.ReturnOnCapitalCurrent)
	setFloat(&record.ReturnOnCapital5YAvg, e.Data.Analysis.CompanyMetrics.ReturnOnCapital5YearAverage)
	setFloat(&record.IncomeEmployee, e.Data.Analysis.CompanyMetrics.IncomeEmployee)
	setFloat(&record.RevenueEmployee, e.Data.Analysis.CompanyMetrics.RevenueEmployee)
	setFloat(&record.AssetTurnover, e.Data.Analysis.CompanyMetrics.AssetTurnover)
	setFloat(&record.InventoryTurnover, e.Data.Analysis.CompanyMetrics.InventoryTurnover)
	setFloat(&record.ReceivableTurnover, e.Data.Analysis.CompanyMetrics.ReceivableTurnover)
	setFloat(&record.RoaTTM, e.Data.Analysis.CompanyMetrics.RoaTTM)

	setFloat(&record.AverageDividendYield5Y, e.Data.Analysis.ShareStatistics.AverageDividendYield5Year)
	setString(&record.LastSplitFactor, e.Data.Analysis.ShareStatistics.LastSplitFactor)
	setFloat(&record.ExDividendAmount, e.Data.Analysis.ShareStatistics.ExDividendAmount)
	setInt64(&record.SharesOutstanding, e.Data.Analysis.ShareStatistics.SharesOutstanding)
	setFloat(&record.EnterpriseValue, e.Data.Analysis.ShareStatistics.EnterpriseValue)
	setFloat(&record.DividendYield, e.Data.Analysis.ShareStatistics.DividendYield)

	setFloat(&record.MarketCap, e.Data.MarketCap)
	setFloat(&record.EnterpriseValue, e.Data.EnterpriseValue)

	if latestAnnual := latestAnnualStatement(e.Data.Analysis.AnnualStatements); latestAnnual != nil {
		setFloat(&record.Assets, latestAnnual.Assets)
		setFloat(&record.Liabilities, latestAnnual.Liabilities)
	}

	lastSplitDate, err := helper.ParseRFC3339Pointer(e.Data.Analysis.ShareStatistics.LastSplitDate)
	if err != nil {
		return fmt.Errorf("invalid last split date: %w", err)
	}
	declarationDate, err := helper.ParseRFC3339Pointer(e.Data.Analysis.ShareStatistics.DeclarationDate)
	if err != nil {
		return fmt.Errorf("invalid declaration date: %w", err)
	}
	dividendDate, err := helper.ParseRFC3339Pointer(e.Data.Analysis.ShareStatistics.DividendDate)
	if err != nil {
		return fmt.Errorf("invalid dividend date: %w", err)
	}
	exDividendDate, err := helper.ParseRFC3339Pointer(e.Data.Analysis.ShareStatistics.ExDividendDate)
	if err != nil {
		return fmt.Errorf("invalid ex-dividend date: %w", err)
	}
	sourceTimeLastUpdated, err := helper.ParseRFC3339Pointer(e.Data.TimeLastUpdated)
	if err != nil {
		return fmt.Errorf("invalid source time last updated: %w", err)
	}

	setTime(&record.LastSplitDate, lastSplitDate)
	setTime(&record.DeclarationDate, declarationDate)
	setTime(&record.DividendDate, dividendDate)
	setTime(&record.ExDividendDate, exDividendDate)
	setTime(&record.SourceTimeLastUpdated, sourceTimeLastUpdated)

	return nil
}

func latestAnnualStatement(statements map[string]EquitiesAnnualStatement) *EquitiesAnnualStatement {
	if len(statements) == 0 {
		return nil
	}

	years := make([]int, 0, len(statements))
	for year := range statements {
		parsedYear, err := strconv.Atoi(year)
		if err != nil {
			continue
		}
		years = append(years, parsedYear)
	}
	if len(years) == 0 {
		return nil
	}

	sort.Ints(years)
	latestKey := strconv.Itoa(years[len(years)-1])
	latest, ok := statements[latestKey]
	if !ok {
		return nil
	}

	return &latest
}

func setFloat(dest **float64, src *float64) {
	if src == nil {
		return
	}
	if math.IsNaN(*src) || math.IsInf(*src, 0) {
		return
	}
	*dest = src
}

func setString(dest **string, src *string) {
	if src == nil || *src == "" {
		return
	}
	*dest = src
}

// func setStringIfNil(dest **string, src *string) {
// 	if *dest != nil {
// 		return
// 	}
// 	if src == nil || *src == "" {
// 		return
// 	}
// 	*dest = src
// }

func setTime(dest **time.Time, src *time.Time) {
	if src == nil {
		return
	}
	*dest = src
}

func setInt64(dest **int64, src *int64) {
	if src == nil {
		return
	}
	*dest = src
}

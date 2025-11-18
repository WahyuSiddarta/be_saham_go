package validator

import (
	"time"

	"github.com/WahyuSiddarta/be_saham_go/middleware"
	"github.com/WahyuSiddarta/be_saham_go/models"
)

// LoginRequest represents login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
}

// RegisterRequest represents registration request payload
type RegisterRequest struct {
	Email           string             `json:"email" validate:"required,email"`
	Password        string             `json:"password" validate:"required,min=6"`
	ConfirmPassword string             `json:"confirmPassword" validate:"required,eqfield=Password"`
	Status          *models.UserStatus `json:"status,omitempty" validate:"omitempty,user_status"`
	UserLevel       *models.UserLevel  `json:"user_level,omitempty" validate:"omitempty,user_level"`
}

// UpdateUserLevelRequest represents request to update user level
type UpdateUserLevelRequest struct {
	UserLevel   models.UserLevel    `json:"user_level" validate:"required,user_level"`
	PaymentData *PaymentDataRequest `json:"payment_data,omitempty"`
}

// PaymentDataRequest represents payment data in requests
type PaymentDataRequest struct {
	OriginalPrice  float64  `json:"original_price" validate:"required,gt=0"`
	PaidPrice      float64  `json:"paid_price" validate:"required,gte=0"`
	DiscountAmount *float64 `json:"discount_amount,omitempty" validate:"omitempty,gte=0"`
	DiscountReason *string  `json:"discount_reason,omitempty" validate:"omitempty,max=255"`
	PaymentMethod  *string  `json:"payment_method,omitempty" validate:"omitempty,max=50"`
	PaymentDate    *string  `json:"payment_date,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Notes          *string  `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

// UpdateUserStatusRequest represents request to update user status
type UpdateUserStatusRequest struct {
	Status models.UserStatus `json:"status" validate:"required,user_status"`
}

// GetUsersQuery represents query parameters for getting users
type GetUsersQuery struct {
	Page        int                `query:"page" validate:"omitempty,min=1"`
	Limit       int                `query:"limit" validate:"omitempty,min=1,max=100"`
	Status      *models.UserStatus `query:"status" validate:"omitempty,user_status"`
	UserLevel   *models.UserLevel  `query:"user_level" validate:"omitempty,user_level"`
	EmailFilter *string            `query:"email_filter" validate:"omitempty,max=100"`
}

// LoginData contains login response data
type LoginData struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

// RegisterData contains registration response data
type RegisterData struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

// ProfileData contains user profile data
type ProfileData struct {
	User *middleware.AuthUser `json:"user"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag,omitempty"`
}

// Portfolio Cash Requests

// CreatePortfolioCashRequest represents request to create cash portfolio
type CreatePortfolioCashRequest struct {
	Account             string   `json:"account" validate:"required"`
	Bank                string   `json:"bank" validate:"required"`
	Amount              float64  `json:"amount" validate:"required,min=0"`
	YieldRate           *float64 `json:"yield_rate" validate:"omitempty,min=0,max=100"`
	YieldPeriod         string   `json:"yield_period"`
	YieldFrequencyType  string   `json:"yield_frequency_type" validate:"required,oneof=daily monthly yearly"`
	YieldFrequencyValue int      `json:"yield_frequency_value" validate:"required,min=1"`
	YieldPaymentType    string   `json:"yield_payment_type" validate:"required"`
	HasMaturity         bool     `json:"has_maturity"`
	MaturityDate        *string  `json:"maturity_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Note                *string  `json:"note"`
	Category            string   `json:"category" validate:"required,oneof=liquid time_deposit money_market other"`
}

// ParsedMaturityDate returns the parsed maturity date if provided.
func (r *CreatePortfolioCashRequest) ParsedMaturityDate() (*time.Time, error) {
	return parseDateString(r.MaturityDate)
}

// UpdatePortfolioCashRequest represents request to update cash portfolio
type UpdatePortfolioCashRequest struct {
	Account             string   `json:"account"`
	Bank                string   `json:"bank"`
	Amount              *float64 `json:"amount" validate:"omitempty,min=0"`
	YieldRate           *float64 `json:"yield_rate" validate:"omitempty,min=0,max=100"`
	YieldPeriod         string   `json:"yield_period"`
	YieldFrequencyType  string   `json:"yield_frequency_type" validate:"omitempty,oneof=daily monthly yearly"`
	YieldFrequencyValue *int     `json:"yield_frequency_value" validate:"omitempty,min=1"`
	YieldPaymentType    string   `json:"yield_payment_type"`
	HasMaturity         *bool    `json:"has_maturity"`
	MaturityDate        *string  `json:"maturity_date" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Note                *string  `json:"note"`
	Status              string   `json:"status" validate:"omitempty,oneof=active maturity"`
	Category            string   `json:"category" validate:"omitempty,oneof=liquid time_deposit money_market other"`
}

// ParsedMaturityDate returns the parsed maturity date if provided.
func (r *UpdatePortfolioCashRequest) ParsedMaturityDate() (*time.Time, error) {
	return parseDateString(r.MaturityDate)
}

// MoveAssetRequest represents request to move asset between portfolios
type MoveAssetRequest struct {
	SourcePortfolioID int `json:"source_portfolio_id" validate:"required,gt=0"`
	TargetPortfolioID int `json:"target_portfolio_id" validate:"required,gt=0"`
}

// RealizeCashPortfolioRequest represents request to realize a cash portfolio
type RealizeCashPortfolioRequest struct {
	PortfolioCashID int     `json:"portfolio_cash_id" validate:"required,gt=0"`
	FinalSaldo      float64 `json:"final_saldo" validate:"required"`
	Amount          float64 `json:"amount" validate:"required"`
	RealizedAt      string  `json:"realized_at" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

// PnL Requests

// CreatePnlRealizedCashRequest represents request to create PnL entry
type CreatePnlRealizedCashRequest struct {
	PortfolioCashID int     `json:"portfolio_cash_id" validate:"required,gt=0"`
	Amount          float64 `json:"amount" validate:"required"`
	Note            *string `json:"note"`
	RealizedAt      *string `json:"realized_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

// UpdatePnlRealizedCashRequest represents request to update PnL entry
type UpdatePnlRealizedCashRequest struct {
	PortfolioCashID int     `json:"portfolio_cash_id" validate:"required,gt=0"`
	Amount          float64 `json:"amount" validate:"required"`
	Note            *string `json:"note"`
	RealizedAt      *string `json:"realized_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

// PnlListQuery represents query parameters for listing PnL
type PnlListQuery struct {
	Limit  int `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset int `query:"offset" validate:"omitempty,min=0"`
}

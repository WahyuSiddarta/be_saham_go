-- PostgreSQL DDL Schema for Saham Application
-- Run this script to recreate the database structure

-- ============================================================================
-- ENUMS / CUSTOM TYPES
-- ============================================================================

CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended', 'banned');
CREATE TYPE user_level AS ENUM ('free', 'premium', 'premium+', 'admin');
CREATE TYPE coupon_frequency AS ENUM ('monthly', 'quarterly', 'semi-annual', 'annual');
CREATE TYPE portfolio_status AS ENUM ('active', 'inactive', 'closed');
CREATE TYPE coupon_status AS ENUM ('pending', 'paid', 'missed', 'cancelled');
CREATE TYPE payment_status AS ENUM ('pending', 'completed', 'failed', 'refunded');
CREATE TYPE yield_period AS ENUM ('daily', 'weekly', 'monthly', 'quarterly', 'semi-annual', 'annual');
CREATE TYPE yield_frequency_type AS ENUM ('daily', 'weekly', 'monthly', 'quarterly', 'semi-annual', 'annual');
CREATE TYPE yield_payment_type AS ENUM ('automatic', 'manual', 'reinvest');

-- ============================================================================
-- USERS TABLE
-- ============================================================================

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    status user_status NOT NULL DEFAULT 'active',
    user_level user_level NOT NULL DEFAULT 'free',
    premium_expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_user_level ON users(user_level);

-- ============================================================================
-- PAYMENT RECORDS TABLE
-- ============================================================================

CREATE TABLE payment_records (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subscription_type user_level NOT NULL,
    original_price NUMERIC(12, 2) NOT NULL,
    paid_price NUMERIC(12, 2) NOT NULL,
    discount_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
    discount_reason TEXT,
    payment_method VARCHAR(50) NOT NULL,
    payment_status payment_status NOT NULL,
    payment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    notes TEXT,
    processed_by_admin_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payment_records_user_id ON payment_records(user_id);
CREATE INDEX idx_payment_records_payment_status ON payment_records(payment_status);
CREATE INDEX idx_payment_records_payment_date ON payment_records(payment_date);

-- ============================================================================
-- BOND TRACKER TABLE
-- ============================================================================

CREATE TABLE bond_tracker (
    bond_id VARCHAR(100) PRIMARY KEY,
    market_price NUMERIC(12, 4) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_bond_tracker_created_at ON bond_tracker(created_at);
CREATE INDEX idx_bond_tracker_updated_at ON bond_tracker(updated_at);
CREATE INDEX idx_bond_tracker_deleted_at ON bond_tracker(deleted_at) WHERE deleted_at IS NULL;

-- ============================================================================
-- PORTFOLIO BONDS TABLE
-- ============================================================================

CREATE TABLE portfolio_bonds (
    id SERIAL PRIMARY KEY,
    bond_id VARCHAR(100) NOT NULL REFERENCES bond_tracker(bond_id),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    purchase_price NUMERIC(12, 4) NOT NULL,
    coupon_rate NUMERIC(6, 4) NOT NULL,
    coupon_frequency coupon_frequency NOT NULL,
    next_coupon_date DATE,
    maturity_date DATE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    status portfolio_status NOT NULL DEFAULT 'active',
    note TEXT,
    market_price NUMERIC(12, 4),
    market_price_override NUMERIC(12, 4),
    market_price_override_date TIMESTAMP WITH TIME ZONE,
    secondary_market BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_portfolio_bonds_user_id ON portfolio_bonds(user_id);
CREATE INDEX idx_portfolio_bonds_bond_id ON portfolio_bonds(bond_id);
CREATE INDEX idx_portfolio_bonds_status ON portfolio_bonds(status);
CREATE INDEX idx_portfolio_bonds_deleted_at ON portfolio_bonds(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_portfolio_bonds_user_id_status ON portfolio_bonds(user_id, status) WHERE deleted_at IS NULL;

-- ============================================================================
-- PORTFOLIO BOND COUPONS TABLE
-- ============================================================================

CREATE TABLE portfolio_bond_coupons (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    portfolio_bond_id INTEGER NOT NULL REFERENCES portfolio_bonds(id) ON DELETE CASCADE,
    coupon_number INTEGER NOT NULL,
    payment_date DATE NOT NULL,
    amount NUMERIC(12, 4) NOT NULL,
    status coupon_status NOT NULL DEFAULT 'pending',
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_portfolio_bond_coupons_user_id ON portfolio_bond_coupons(user_id);
CREATE INDEX idx_portfolio_bond_coupons_portfolio_bond_id ON portfolio_bond_coupons(portfolio_bond_id);
CREATE INDEX idx_portfolio_bond_coupons_status ON portfolio_bond_coupons(status);
CREATE INDEX idx_portfolio_bond_coupons_deleted_at ON portfolio_bond_coupons(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_portfolio_bond_coupons_payment_date ON portfolio_bond_coupons(payment_date);

-- ============================================================================
-- PORTFOLIO CASH TABLE
-- ============================================================================

CREATE TABLE portfolio_cash (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account VARCHAR(100) NOT NULL,
    bank VARCHAR(100) NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    yield_rate NUMERIC(6, 4),
    yield_period yield_period NOT NULL DEFAULT 'annual',
    yield_frequency_type yield_frequency_type NOT NULL DEFAULT 'annual',
    yield_frequency_value INTEGER NOT NULL DEFAULT 1,
    yield_payment_type yield_payment_type NOT NULL DEFAULT 'manual',
    has_maturity BOOLEAN NOT NULL DEFAULT FALSE,
    maturity_date DATE,
    note TEXT,
    status portfolio_status NOT NULL DEFAULT 'active',
    category VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_portfolio_cash_user_id ON portfolio_cash(user_id);
CREATE INDEX idx_portfolio_cash_status ON portfolio_cash(status);
CREATE INDEX idx_portfolio_cash_category ON portfolio_cash(category);
CREATE INDEX idx_portfolio_cash_deleted_at ON portfolio_cash(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_portfolio_cash_user_id_status ON portfolio_cash(user_id, status) WHERE deleted_at IS NULL;

-- ============================================================================
-- PORTFOLIO PNL REALIZED CASH TABLE
-- ============================================================================

CREATE TABLE portfolio_pnl_realized_cash (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    portfolio_cash_id INTEGER NOT NULL REFERENCES portfolio_cash(id) ON DELETE CASCADE,
    amount NUMERIC(15, 2) NOT NULL,
    realized_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_portfolio_pnl_realized_cash_user_id ON portfolio_pnl_realized_cash(user_id);
CREATE INDEX idx_portfolio_pnl_realized_cash_portfolio_cash_id ON portfolio_pnl_realized_cash(portfolio_cash_id);
CREATE INDEX idx_portfolio_pnl_realized_cash_realized_at ON portfolio_pnl_realized_cash(realized_at);
CREATE INDEX idx_portfolio_pnl_realized_cash_deleted_at ON portfolio_pnl_realized_cash(deleted_at) WHERE deleted_at IS NULL;

-- ============================================================================
-- TRIGGERS FOR UPDATED_AT
-- ============================================================================

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to all tables with updated_at
CREATE TRIGGER trigger_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trigger_payment_records_updated_at BEFORE UPDATE ON payment_records
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trigger_bond_tracker_updated_at BEFORE UPDATE ON bond_tracker
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trigger_portfolio_bonds_updated_at BEFORE UPDATE ON portfolio_bonds
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trigger_portfolio_bond_coupons_updated_at BEFORE UPDATE ON portfolio_bond_coupons
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trigger_portfolio_cash_updated_at BEFORE UPDATE ON portfolio_cash
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER trigger_portfolio_pnl_realized_cash_updated_at BEFORE UPDATE ON portfolio_pnl_realized_cash
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

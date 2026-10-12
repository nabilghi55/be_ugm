-- Enable UUID Extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Branches
CREATE TABLE IF NOT EXISTS branches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    address TEXT,
    phone VARCHAR(30),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id) ON DELETE SET NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role VARCHAR(30) NOT NULL CHECK (role IN ('SUPER_ADMIN', 'ADMIN', 'APPRAISER', 'CASHIER', 'MANAGER')),
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Refresh Tokens
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    replaced_by_token VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Password Reset Tokens
CREATE TABLE IF NOT EXISTS password_resets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Customers
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_code VARCHAR(30) UNIQUE NOT NULL,
    nik VARCHAR(20) UNIQUE NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    birth_place VARCHAR(50),
    birth_date DATE,
    gender VARCHAR(10) CHECK (gender IN ('MALE', 'FEMALE')),
    address TEXT,
    province VARCHAR(50),
    city VARCHAR(50),
    district VARCHAR(50),
    village VARCHAR(50),
    phone VARCHAR(30) NOT NULL,
    whatsapp VARCHAR(30),
    email VARCHAR(100),
    occupation VARCHAR(50),
    mother_maiden_name VARCHAR(100),
    emergency_contact_name VARCHAR(100),
    emergency_contact_phone VARCHAR(30),
    emergency_contact_relationship VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Customer Documents
CREATE TABLE IF NOT EXISTS customer_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    document_type VARCHAR(30) NOT NULL CHECK (document_type IN ('KTP', 'SELFIE', 'SELFIE_WITH_KTP', 'OTHER')),
    object_key VARCHAR(255) NOT NULL,
    ocr_data JSONB,
    ocr_confidence NUMERIC(5,2),
    verification_status VARCHAR(20) DEFAULT 'UNVERIFIED' CHECK (verification_status IN ('UNVERIFIED', 'VERIFIED', 'REJECTED')),
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Asset Categories & LTV Rules
CREATE TABLE IF NOT EXISTS asset_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    ltv_percent NUMERIC(5,2) NOT NULL DEFAULT 70.00,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Category Inspection Checklist Config
CREATE TABLE IF NOT EXISTS category_checklists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID NOT NULL REFERENCES asset_categories(id) ON DELETE CASCADE,
    item_name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Assets
CREATE TABLE IF NOT EXISTS assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_code VARCHAR(30) UNIQUE NOT NULL,
    category_id UUID NOT NULL REFERENCES asset_categories(id),
    subcategory VARCHAR(50),
    brand VARCHAR(50) NOT NULL,
    model VARCHAR(100) NOT NULL,
    variant VARCHAR(100),
    serial_number VARCHAR(100),
    imei VARCHAR(50),
    year INT,
    color VARCHAR(30),
    specification TEXT,
    accessories TEXT,
    physical_condition VARCHAR(20) CHECK (physical_condition IN ('EXCELLENT', 'GOOD', 'FAIR', 'POOR', 'BROKEN')),
    functional_condition VARCHAR(20) CHECK (functional_condition IN ('EXCELLENT', 'GOOD', 'FAIR', 'POOR', 'BROKEN')),
    purchase_price BIGINT DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Asset Photos
CREATE TABLE IF NOT EXISTS asset_photos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    object_key VARCHAR(255) NOT NULL,
    image_type VARCHAR(30) NOT NULL CHECK (image_type IN ('FRONT', 'BACK', 'LEFT', 'RIGHT', 'SERIAL_NUMBER', 'DAMAGE', 'ACCESSORY', 'PACKAGING', 'OTHER')),
    is_primary BOOLEAN DEFAULT FALSE,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Asset Inspection Checklist Results
CREATE TABLE IF NOT EXISTS asset_inspections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    item_name VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('GOOD', 'FAIR', 'POOR', 'BROKEN')),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Market Comparables
CREATE TABLE IF NOT EXISTS market_comparables (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source VARCHAR(50) NOT NULL,
    external_id VARCHAR(100),
    title VARCHAR(255) NOT NULL,
    url TEXT,
    price BIGINT NOT NULL,
    location VARCHAR(100),
    condition VARCHAR(50),
    seller_type VARCHAR(50),
    brand VARCHAR(50) NOT NULL,
    model VARCHAR(100) NOT NULL,
    variant VARCHAR(100),
    specification TEXT,
    fetched_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    similarity_score NUMERIC(5,2) DEFAULT 0.00
);

-- Appraisals
CREATE TABLE IF NOT EXISTS appraisals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_id UUID NOT NULL REFERENCES assets(id),
    market_price BIGINT NOT NULL DEFAULT 0,
    market_price_min BIGINT NOT NULL DEFAULT 0,
    market_price_max BIGINT NOT NULL DEFAULT 0,
    market_confidence INT DEFAULT 0,
    historical_price BIGINT DEFAULT 0,
    condition_adjustment_percent NUMERIC(5,2) DEFAULT 0.00,
    adjusted_market_price BIGINT NOT NULL DEFAULT 0,
    appraiser_price BIGINT NOT NULL DEFAULT 0,
    appraiser_notes TEXT,
    loan_to_value NUMERIC(5,2) NOT NULL DEFAULT 70.00,
    maximum_loan BIGINT NOT NULL DEFAULT 0,
    recommended_loan BIGINT NOT NULL DEFAULT 0,
    appraiser_id UUID REFERENCES users(id),
    appraised_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Approval Rules
CREATE TABLE IF NOT EXISTS approval_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    min_amount BIGINT NOT NULL,
    max_amount BIGINT NOT NULL,
    required_role VARCHAR(30) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Pawn Transactions
CREATE TABLE IF NOT EXISTS pawn_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pawn_number VARCHAR(30) UNIQUE NOT NULL,
    customer_id UUID NOT NULL REFERENCES customers(id),
    branch_id UUID NOT NULL REFERENCES branches(id),
    appraisal_id UUID NOT NULL REFERENCES appraisals(id),
    principal_amount BIGINT NOT NULL,
    interest_rate NUMERIC(5,2) NOT NULL DEFAULT 2.00,
    administration_fee BIGINT NOT NULL DEFAULT 50000,
    storage_fee BIGINT NOT NULL DEFAULT 0,
    penalty_rate NUMERIC(5,2) NOT NULL DEFAULT 0.50,
    loan_date DATE NOT NULL,
    due_date DATE NOT NULL,
    grace_period INT DEFAULT 7,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    requested_by UUID REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP WITH TIME ZONE,
    approval_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Disbursements
CREATE TABLE IF NOT EXISTS disbursements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pawn_transaction_id UUID NOT NULL REFERENCES pawn_transactions(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL,
    method VARCHAR(30) NOT NULL CHECK (method IN ('CASH', 'BANK_TRANSFER', 'EWALLET')),
    bank_name VARCHAR(50),
    account_number VARCHAR(50),
    reference_number VARCHAR(100),
    processed_by UUID REFERENCES users(id),
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Payments
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pawn_transaction_id UUID NOT NULL REFERENCES pawn_transactions(id),
    payment_number VARCHAR(30) UNIQUE NOT NULL,
    payment_type VARCHAR(30) NOT NULL CHECK (payment_type IN ('INTEREST_PAYMENT', 'PARTIAL_PAYMENT', 'EXTENSION_PAYMENT', 'FULL_PAYOFF')),
    principal_paid BIGINT NOT NULL DEFAULT 0,
    interest_paid BIGINT NOT NULL DEFAULT 0,
    admin_paid BIGINT NOT NULL DEFAULT 0,
    storage_paid BIGINT NOT NULL DEFAULT 0,
    penalty_paid BIGINT NOT NULL DEFAULT 0,
    total_amount BIGINT NOT NULL,
    payment_method VARCHAR(30) NOT NULL CHECK (payment_method IN ('CASH', 'BANK_TRANSFER', 'EWALLET')),
    reference_number VARCHAR(100),
    processed_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Pawn Extensions
CREATE TABLE IF NOT EXISTS pawn_extensions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pawn_transaction_id UUID NOT NULL REFERENCES pawn_transactions(id),
    payment_id UUID NOT NULL REFERENCES payments(id),
    old_due_date DATE NOT NULL,
    new_due_date DATE NOT NULL,
    days_extended INT NOT NULL DEFAULT 30,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Asset Releases
CREATE TABLE IF NOT EXISTS asset_releases (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pawn_transaction_id UUID UNIQUE NOT NULL REFERENCES pawn_transactions(id),
    receiver_name VARCHAR(100) NOT NULL,
    receiver_id_number VARCHAR(50),
    proof_object_key VARCHAR(255),
    processed_by UUID NOT NULL REFERENCES users(id),
    released_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    notes TEXT
);

-- Inventories
CREATE TABLE IF NOT EXISTS inventories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id),
    warehouse VARCHAR(50) NOT NULL DEFAULT 'Gudang Utama',
    rack VARCHAR(50) NOT NULL DEFAULT 'Rak A',
    box VARCHAR(50) NOT NULL DEFAULT 'Box 01',
    asset_id UUID UNIQUE NOT NULL REFERENCES assets(id),
    pawn_transaction_id UUID REFERENCES pawn_transactions(id),
    status VARCHAR(30) NOT NULL DEFAULT 'STORED',
    qr_code VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Audit Logs
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(50) NOT NULL,
    entity VARCHAR(50) NOT NULL,
    entity_id VARCHAR(100),
    old_value JSONB,
    new_value JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Seed Initial Data (Kota Padang Branches)
INSERT INTO branches (id, code, name, address, phone) VALUES
('00000000-0000-0000-0000-000000000001', 'BR-PDG-001', 'Cabang Nanggalo', 'Jl. Siteba No. 12, Nanggalo, Kota Padang', '0751-441122'),
('00000000-0000-0000-0000-000000000002', 'BR-PDG-002', 'Cabang Pasar Raya', 'Jl. M. Yamin No. 45, Pasar Raya, Kota Padang', '0751-442233'),
('00000000-0000-0000-0000-000000000003', 'BR-PDG-003', 'Cabang Azizi', 'Jl. Dr. Sutomo No. 88, Azizi, Kota Padang', '0751-443344')
ON CONFLICT (code) DO NOTHING;

-- Seed Super Admin (password: Admin123!)
INSERT INTO users (id, branch_id, username, email, password_hash, full_name, role) VALUES
('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 'admin', 'admin@samplegadai.com', '$2a$10$tZ922R4.4y8Gq/Wd3q3eXe5T1e00e23805f159154f248e3532f1.', 'Super Administrator', 'SUPER_ADMIN')
ON CONFLICT (username) DO NOTHING;

-- Seed Asset Categories
INSERT INTO asset_categories (id, name, ltv_percent) VALUES
('c0000000-0000-0000-0000-000000000001', 'Smartphone', 70.00),
('c0000000-0000-0000-0000-000000000002', 'Laptop', 65.00),
('c0000000-0000-0000-0000-000000000003', 'Kamera', 60.00),
('c0000000-0000-0000-0000-000000000004', 'Emas', 85.00),
('c0000000-0000-0000-0000-000000000005', 'Kendaraan', 70.00)
ON CONFLICT (name) DO NOTHING;

-- Seed Smartphone Checklist Config
INSERT INTO category_checklists (category_id, item_name, description) VALUES
('c0000000-0000-0000-0000-000000000001', 'Screen', 'Fisik kaca dan keretakan'),
('c0000000-0000-0000-0000-000000000001', 'Touchscreen', 'Respon sentuhan seluruh layar'),
('c0000000-0000-0000-0000-000000000001', 'Camera Front', 'Fungsi kamera depan'),
('c0000000-0000-0000-0000-000000000001', 'Camera Back', 'Fungsi kamera belakang & zoom'),
('c0000000-0000-0000-0000-000000000001', 'Speaker', 'Kejelasan suara media'),
('c0000000-0000-0000-0000-000000000001', 'Microphone', 'Input suara rekaman & telepon'),
('c0000000-0000-0000-0000-000000000001', 'Battery', 'Kesehatan baterai & pemakaian'),
('c0000000-0000-0000-0000-000000000001', 'Charging', 'Pengisian daya port Lightning/Type-C'),
('c0000000-0000-0000-0000-000000000001', 'WiFi', 'Koneksi jaringan nirkabel'),
('c0000000-0000-0000-0000-000000000001', 'Bluetooth', 'Koneksi perangkat sekitar'),
('c0000000-0000-0000-0000-000000000001', 'SIM', 'Pembacaan kartu SIM & jaringan seluler'),
('c0000000-0000-0000-0000-000000000001', 'Fingerprint', 'Sensor sidik jari'),
('c0000000-0000-0000-0000-000000000001', 'Face ID', 'Pengenalan wajah'),
('c0000000-0000-0000-0000-000000000001', 'Body', 'Mulus/lecet/dent casing'),
('c0000000-0000-0000-0000-000000000001', 'Buttons', 'Tombol Power & Volume'),
('c0000000-0000-0000-0000-000000000001', 'IMEI', 'Kesesuaian IMEI fisik vs sistem');

-- Seed Approval Rules
INSERT INTO approval_rules (min_amount, max_amount, required_role) VALUES
(0, 5000000, 'APPRAISER'),
(5000001, 20000000, 'MANAGER'),
(20000001, 999999999999, 'SUPER_ADMIN');

-- Create Indexes
CREATE INDEX IF NOT EXISTS idx_customers_nik ON customers(nik);
CREATE INDEX IF NOT EXISTS idx_assets_code ON assets(asset_code);
CREATE INDEX IF NOT EXISTS idx_pawn_status ON pawn_transactions(status);
CREATE INDEX IF NOT EXISTS idx_pawn_due_date ON pawn_transactions(due_date);
CREATE INDEX IF NOT EXISTS idx_audit_created ON audit_logs(created_at);

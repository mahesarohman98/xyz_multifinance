CREATE TABLE IF NOT EXISTS Customers (
    customer_id VARCHAR(100) PRIMARY KEY,
    nik VARCHAR(100) NOT NULL UNIQUE,
    full_name VARCHAR(100) NOT NULL,
    legal_name VARCHAR(100) NOT NULL,
    place_of_birth VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    wages VARCHAR(100) NOT NULL,
    ktp_photo_url VARCHAR(100) NOT NULL,
    photo_url VARCHAR(100) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS TenorLimits (
    customer_id VARCHAR(100),
    month INTEGER,
    amount DECIMAL NOT NULL,
    used_amount DECIMAL NOT NULL,
    PRIMARY KEY(customer_id, month),
    FOREIGN KEY (customer_id) REFERENCES Customers(customer_id)
);

-- source example ecommerce, dealer, web
CREATE TABLE IF NOT EXISTS Sources (
    source_id VARCHAR(100) PRIMARY KEY,
    category VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS Transactions (
    transaction_id VARCHAR(100) PRIMARY KEY,
    contract_number VARCHAR(100) UNIQUE,
    customer_id VARCHAR(100),
    external_id VARCHAR(100),
    source_id VARCHAR(100),
    otr DECIMAL NOT NULL,
    admin_fee DECIMAL NOT NULL,
    total_borowed DECIMAL NOT NULL,
    installment_amount DECIMAL NOT NULL,
    amount_of_interest DECIMAL NOT NULL,
    asset_name VARCHAR(100) NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES Customers(customer_id),
    FOREIGN KEY (source_id) REFERENCES Sources(source_id)
);

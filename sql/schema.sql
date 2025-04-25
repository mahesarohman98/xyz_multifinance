CREATE TABLE IF NOT EXISTS Customers (
    customer_id VARCHAR(100) PRIMARY KEY,
    nik VARCHAR(100) UNIQUE,
    full_name VARCHAR(100),
    legal_name VARCHAR(100),
    place_of_birth VARCHAR(100),
    date_of_birth VARCHAR(100),
    wages VARCHAR(100),
    ktp_photo_url VARCHAR(100),
    photo_url VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS Tenors (
    customer_id VARCHAR(100) PRIMARY KEY,
    month INTEGER PRIMARY KEY,
    amount DECIMAL,
    FOREIGN KEY (customer_id) REFERENCES Customers(customer_id)
);

-- source example ecommerce, dealer, web
CREATE TABLE IF NOT EXISTS Source (
    source_id VARCHAR(100) PRIMARY KEY
    source_name VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS Transactions (
    contract_number VARCHAR(100) PRIMARY KEY,
    customer_id VARCHAR(100),
    external_id VARCHAR(100),
    source_id VARCHAR(100),
    admin_fee DECIMAL,
    installment_amount DECIMAL,
    amount_of_interest DECIMAL,
    asset_name VARCHAR(100),
    FOREIGN KEY (customer_id) REFERENCES Customers(customer_id)
    FOREIGN KEY (source_id) REFERENCES Customers(source_id)
);

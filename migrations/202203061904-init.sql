--
-- Initialize database with basic tables
--
CREATE TABLE IF NOT EXISTS migrations (
    id VARCHAR PRIMARY KEY
);

SELECT EXISTS (
    SELECT id FROM migrations WHERE id = :'MIGRATION_ID'
) as migrated \gset

\if :migrated
    \echo 'migration' :MIGRATION_ID 'already exists, skipping'
\else
    \echo 'migration' :MIGRATION_ID 'does not exist'

    CREATE TABLE users (
        id VARCHAR(16) PRIMARY KEY,
        email VARCHAR NOT NULL,
        phone VARCHAR NOT NULL,
        password VARCHAR NOT NULL,
        name VARCHAR NOT NULL,
        blocked_until TIMESTAMP,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP
    );

    CREATE INDEX idx_user_indices_deleted_at ON users(deleted_at);

    CREATE TABLE banks (
        id VARCHAR PRIMARY KEY,
        name VARCHAR UNIQUE NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP
    );

    CREATE INDEX idx_bank_indices_deleted_at ON banks(deleted_at);

    CREATE TABLE accounts (
        id VARCHAR(16) PRIMARY KEY,
        is_bank BOOLEAN NOT NULL,
        user_id VARCHAR(16),
        bank_id VARCHAR NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP,
        CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
        CONSTRAINT fk_bank FOREIGN KEY (bank_id) REFERENCES banks(id)
    );

    CREATE INDEX idx_account_indices_deleted_at ON accounts(deleted_at);

    CREATE TABLE products (
        id VARCHAR(16) PRIMARY KEY,
        title VARCHAR NOT NULL,
        description VARCHAR NOT NULL,
        is_on_market BOOLEAN NOT NULL DEFAULT FALSE,
        owner_id VARCHAR(16) NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP,
        CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users(id)
    );

    CREATE INDEX idx_product_indices_deleted_at ON products(deleted_at);

    CREATE TABLE product_images (
        id VARCHAR(16) PRIMARY KEY,
        filename VARCHAR NOT NULL,
        path VARCHAR NOT NULL,
        product_id VARCHAR(16) NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP,
        CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id)
    );

    CREATE INDEX idx_product_image_indices_deleted_at ON product_images(deleted_at);

    CREATE TABLE offers (
        id VARCHAR(16) PRIMARY KEY,
        delete_on_sell BOOLEAN NOT NULL DEFAULT TRUE, 
        product_id VARCHAR(16) NOT NULL,
        user_id VARCHAR(16) NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP,
        CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id),
        CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
    );

    CREATE INDEX idx_offer_indices_deleted_at ON offers(deleted_at);

    CREATE TYPE transaction_status AS ENUM (
        'offering',
        'processing',
        'error',
        'ok'
    );

    CREATE TYPE transaction_type AS ENUM (
        'deposit',
        'buy',
        'fee',
        'withdrawal'
    );

    CREATE TYPE transaction_currency AS ENUM (
        'rub',
        'eur',
        'usd'
    );

    CREATE TABLE transactions (
        id SERIAL PRIMARY KEY,
        date TIMESTAMP,
        status transaction_status NOT NULL,
        type transaction_type NOT NULL,
        currency transaction_currency NOT NULL,
        amount DECIMAL(12, 2) NOT NULL,
        error VARCHAR,
        account_from_id VARCHAR(16) NOT NULL,
        account_to_id VARCHAR(16) NOT NULL,
        offer_id VARCHAR(16),
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP,
        CONSTRAINT fk_accunt_from FOREIGN KEY (account_from_id) REFERENCES accounts(id),
        CONSTRAINT fk_accunt_to FOREIGN KEY (account_to_id) REFERENCES accounts(id),
        CONSTRAINT fk_offer FOREIGN KEY (offer_id) REFERENCES offers(id)
    );

    CREATE INDEX idx_transaction_indices_deleted_at ON transactions(deleted_at);

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;
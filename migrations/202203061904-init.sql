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

    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    CREATE EXTENSION IF NOT EXISTS "pg_cron";

    CREATE TABLE guests (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        expires_at TIMESTAMP NOT NULL
    );

    -- FIXME: change to 0 0 * * * and add VACUUM
    SELECT cron.schedule('*/5 * * * *', $$DELETE FROM guests WHERE expires_at < NOW()$$);

    CREATE TABLE users (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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

    CREATE TABLE token_creators (
        ID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        user_id UUID,
        guest_id UUID,
        CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
        CONSTRAINT fk_guest FOREIGN KEY (guest_id) REFERENCES guests(id) ON DELETE SET NULL,
        CONSTRAINT chk_user_or_guest CHECK (
            user_id IS NULL OR guest_id IS NULL OR user_id = guest_id),
        UNIQUE (user_id, guest_id)
    );

    -- FIXME: change to 0 0 * * * and add VACUUM
    SELECT cron.schedule('*/5 * * * *', $$DELETE FROM token_creators WHERE user_id IS NULL AND guest_id IS NULL$$);

    CREATE OR REPLACE FUNCTION random_between(low INT ,high INT) 
    RETURNS INT AS
    $$
    BEGIN
        RETURN floor(random()* (high-low + 1) + low);
    END;
    $$ language 'plpgsql' STRICT;

    CREATE TYPE token_action AS ENUM (
      'APPROVE_USER_EMAIL',
      'APPROVE_USER_PHONE'
    );

    CREATE TABLE tokens (
        id DECIMAL(6, 0) PRIMARY KEY DEFAULT random_between(100000, 999999)::DECIMAL(6, 0),
        activated_at TIMESTAMP,
        expires_at TIMESTAMP NOT NULL,
        action token_action NOT NULL,
        data JSONB,
        creator_id UUID NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP,
        CONSTRAINT fk_creator FOREIGN KEY (creator_id) REFERENCES token_creators(id) ON DELETE CASCADE
    );

    CREATE INDEX idx_token_indices_deleted_at ON tokens(deleted_at);

    -- FIXME: change to 0 0 * * * and add VACUUM
    SELECT cron.schedule('*/5 * * * *', $$DELETE FROM tokens WHERE expires_at < NOW() AND activated_at IS NULL$$);

    CREATE TABLE banks (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR UNIQUE NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP
    );

    CREATE INDEX idx_bank_indices_deleted_at ON banks(deleted_at);

    CREATE TABLE accounts (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        is_bank BOOLEAN NOT NULL,
        user_id UUID,
        bank_id UUID NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP,
        CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
        CONSTRAINT fk_bank FOREIGN KEY (bank_id) REFERENCES banks(id)
    );

    CREATE INDEX idx_account_indices_deleted_at ON accounts(deleted_at);

    CREATE TABLE products (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        title VARCHAR NOT NULL,
        description VARCHAR NOT NULL,
        is_on_market BOOLEAN NOT NULL DEFAULT FALSE,
        creator_id UUID NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP,
        CONSTRAINT fk_creator FOREIGN KEY (creator_id) REFERENCES users(id)
    );

    CREATE INDEX idx_product_indices_deleted_at ON products(deleted_at);

    CREATE TABLE product_images (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        filename VARCHAR NOT NULL,
        path VARCHAR NOT NULL,
        product_id UUID NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP,
        CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id)
    );

    CREATE INDEX idx_product_image_indices_deleted_at ON product_images(deleted_at);

    CREATE TYPE offer_state AS ENUM (
        'CREATED',
        'CANCELLED',
        'TRANSFERRING_MONEY',
        'TRANSFER_MONEY_FAILED',
        'TRANSFERRING_PRODUCT',
        'TRANSFER_PRODUCT_FAILED',
        'SUCCEEDED',
        'RETURNING_MONEY',
        'RETURN_MONEY_FAILED',
        'MONEY_RETURNED'
    );

    CREATE TABLE offers (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        state offer_state NOT NULL DEFAULT 'CREATED',
        fail_message VARCHAR,
        delete_on_sell BOOLEAN NOT NULL DEFAULT TRUE, 
        product_id UUID NOT NULL,
        consumer_id UUID NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP,
        CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id),
        CONSTRAINT fk_consumer FOREIGN KEY (consumer_id) REFERENCES users(id)
    );

    CREATE INDEX idx_offer_indices_deleted_at ON offers(deleted_at);

    CREATE TYPE transaction_state AS ENUM (
        'CREATED',
        'CANCELLED',
        'PROCESSING',
        'ERROR',
        'SUCCEEDED',
        'FAILED'
    );

    CREATE TYPE transaction_type AS ENUM (
        'DEPOSIT',
        'BUY',
        'FEE',
        'WITHDRAWAL'
    );

    CREATE TYPE transaction_currency AS ENUM (
        'RUB',
        'EUR',
        'USD'
    );

    CREATE TABLE transactions (
        id SERIAL PRIMARY KEY,
        date TIMESTAMP,
        state transaction_state NOT NULL DEFAULT 'CREATED',
        type transaction_type NOT NULL,
        currency transaction_currency NOT NULL,
        amount DECIMAL(12, 2) NOT NULL,
        error VARCHAR,
        account_from_id UUID NOT NULL,
        account_to_id UUID NOT NULL,
        offer_id UUID,
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
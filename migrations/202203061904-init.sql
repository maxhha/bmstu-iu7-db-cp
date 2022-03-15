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

    CREATE EXTENSION IF NOT EXISTS "pgcrypto";
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    CREATE EXTENSION IF NOT EXISTS "pg_cron";

    -- shorkey used for prettier urls
    CREATE DOMAIN SHORTKEY as varchar(11);

    CREATE OR REPLACE FUNCTION shortkey_generate()
    RETURNS TRIGGER AS $$
    DECLARE
        gkey TEXT;
        key SHORTKEY;
        qry TEXT;
        found TEXT;
        user_id BOOLEAN;
    BEGIN
        -- generate the first part of a query as a string with safely
        -- escaped table name, using || to concat the parts
        qry := 'SELECT id FROM ' || quote_ident(TG_TABLE_NAME) || ' WHERE id=';

        LOOP
            -- deal with user-supplied keys, they don't have to be valid base64
            -- only the right length for the type
            IF NEW.id IS NOT NULL AND length(NEW.id) > 0 THEN
                key := NEW.id;
                user_id := TRUE;

                IF length(key) <> 11 THEN
                    RAISE 'User defined key value % has invalid length. Expected 11, got %.', key, length(key);
                END IF;
            ELSE
                -- 8 bytes gives a collision p = .5 after 5.1 x 10^9 values
                gkey := encode(gen_random_bytes(8), 'base64');
                gkey := replace(gkey, '/', '_');  -- url safe replacement
                gkey := replace(gkey, '+', '-');  -- url safe replacement
                key := rtrim(gkey, '=');          -- cut off padding
                user_id := FALSE;
            END IF;

            -- Concat the generated key (safely quoted) with the generated query
            -- and run it.
            -- SELECT id FROM "test" WHERE id='blahblah' INTO found
            -- Now "found" will be the duplicated id or NULL.
            EXECUTE qry || quote_literal(key) INTO found;

            -- Check to see if found is NULL.
            -- If we checked to see if found = NULL it would always be FALSE
            -- because (NULL = NULL) is always FALSE.
            IF found IS NULL THEN
                -- If we didn't find a collision then leave the LOOP.
                EXIT;
            END IF;

            IF user_id THEN
                -- User supplied ID but it violates the PK unique constraint
                RAISE 'ID % already exists in table %', key, TG_TABLE_NAME;
            END IF;

            -- We haven't EXITed yet, so return to the top of the LOOP
            -- and try again.
        END LOOP;

        -- NEW and OLD are available in TRIGGER PROCEDURES.
        -- NEW is the mutated row that will actually be INSERTed.
        -- We're replacing id, regardless of what it was before
        -- with our key variable.
        NEW.id = key;

        -- The RECORD returned here is what will actually be INSERTed,
        -- or what the next trigger will get if there is one.
        RETURN NEW;
    END
    $$ language 'plpgsql';

    -- guests are users that dont have anything but tokens
    CREATE TABLE guests (
        id SHORTKEY PRIMARY KEY,
        expires_at TIMESTAMP NOT NULL
    );

    -- generate shortkey for each insert
    CREATE TRIGGER trgr_guests_genid 
        BEFORE INSERT ON guests FOR EACH ROW 
        EXECUTE PROCEDURE shortkey_generate();

    -- remove expired guests record
    -- FIXME: change to 0 0 * * * and add VACUUM
    SELECT cron.schedule('*/5 * * * *', $$DELETE FROM guests WHERE expires_at < NOW()$$);

    CREATE TABLE users (
        id SHORTKEY PRIMARY KEY,
        email VARCHAR NOT NULL,
        phone VARCHAR NOT NULL,
        password VARCHAR NOT NULL,
        name VARCHAR NOT NULL,
        blocked_until TIMESTAMP,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP
    );

    -- generate shortkey for each insert
    CREATE TRIGGER trgr_users_genid 
        BEFORE INSERT ON users FOR EACH ROW 
        EXECUTE PROCEDURE shortkey_generate();

    CREATE INDEX idx_user_indices_deleted_at ON users(deleted_at);

    -- first, each guest have its own token creator
    -- when guest become a user he inherits token creator
    CREATE TABLE token_creators (
        ID SHORTKEY PRIMARY KEY,
        user_id SHORTKEY,
        guest_id SHORTKEY,
        CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
        CONSTRAINT fk_guest FOREIGN KEY (guest_id) REFERENCES guests(id) ON DELETE SET NULL,
        CONSTRAINT chk_user_or_guest CHECK (
            user_id IS NULL OR guest_id IS NULL OR user_id = guest_id),
        UNIQUE (user_id, guest_id)
    );

    -- generate shortkey for each insert
    CREATE TRIGGER trgr_token_creators_genid 
        BEFORE INSERT ON token_creators FOR EACH ROW 
        EXECUTE PROCEDURE shortkey_generate();

    -- remove token creators without any foreign key
    -- FIXME: change to 0 0 * * * and add VACUUM
    SELECT cron.schedule('*/5 * * * *', $$DELETE FROM token_creators WHERE user_id IS NULL AND guest_id IS NULL$$);

    CREATE OR REPLACE FUNCTION insert_token_creator_for_guest()
    RETURNS TRIGGER AS $$
    BEGIN
        INSERT INTO token_creators (guest_id)
        SELECT NEW.id;

        RETURN NEW;
    END
    $$ LANGUAGE 'plpgsql' STRICT;

    CREATE TRIGGER trger_guests_token_creator
        AFTER INSERT ON guests FOR EACH ROW 
        EXECUTE PROCEDURE insert_token_creator_for_guest();

    CREATE TYPE token_action AS ENUM (
      'APPROVE_USER_EMAIL',
      'APPROVE_USER_PHONE'
    );

    -- tokens are sent to clients by other means
    -- they must be activated to perform some actions 
    CREATE TABLE tokens (
        id DECIMAL(6, 0),
        activated_at TIMESTAMP,
        expires_at TIMESTAMP NOT NULL,
        action token_action NOT NULL,
        data JSONB,
        creator_id SHORTKEY NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP,
        CONSTRAINT fk_creator FOREIGN KEY (creator_id) REFERENCES token_creators(id) ON DELETE CASCADE,
        PRIMARY KEY (id, creator_id)
    );

    CREATE INDEX idx_token_indices_deleted_at ON tokens(deleted_at);

    -- delete expired and not activated tokens
    -- FIXME: change to 0 0 * * * and add VACUUM
    SELECT cron.schedule('*/5 * * * *', $$DELETE FROM tokens WHERE expires_at < NOW() AND activated_at IS NULL$$);

    CREATE OR REPLACE FUNCTION token_generate()
    RETURNS TRIGGER AS $$
    DECLARE
        key DECIMAL(6, 0);
        qry TEXT;
        found DECIMAL(6, 0);
    BEGIN
        -- tokens must be generated by this function
        -- id = 0 is also means empty
        IF NEW.id IS NOT NULL AND NEW.id <> 0 THEN
            RAISE 'Tokens id must be generated by database';
        END IF;

        IF NEW.creator_id IS NULL OR length(NEW.creator_id) <> 11 THEN
            RAISE 'creator_id must be provided';
        END IF;

        -- query to check if this token already exists
        qry := 'SELECT id, creator_id FROM tokens WHERE creator_id = ' || quote_literal(NEW.creator_id) || ' AND id =';

        LOOP
            key := floor(random() * 999999 + 1);

            EXECUTE qry || key::TEXT INTO found;

            IF found IS NULL THEN
                EXIT;
            END IF;
        END LOOP;

        NEW.id = key;

        RETURN NEW;
    END
    $$ LANGUAGE 'plpgsql';

    CREATE TRIGGER trgr_tokens_genid
        BEFORE INSERT ON tokens FOR EACH ROW 
        EXECUTE PROCEDURE token_generate();

    -- CREATE TABLE banks (
    --     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    --     name VARCHAR UNIQUE NOT NULL,
    --     created_at TIMESTAMP NOT NULL,
    --     updated_at TIMESTAMP NOT NULL,
    --     deleted_at TIMESTAMP
    -- );

    -- CREATE INDEX idx_bank_indices_deleted_at ON banks(deleted_at);

    -- CREATE TABLE accounts (
    --     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    --     is_bank BOOLEAN NOT NULL,
    --     user_id UUID,
    --     bank_id UUID NOT NULL,
    --     created_at TIMESTAMP NOT NULL,
    --     updated_at TIMESTAMP NOT NULL,
    --     deleted_at TIMESTAMP,
    --     CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    --     CONSTRAINT fk_bank FOREIGN KEY (bank_id) REFERENCES banks(id)
    -- );

    -- CREATE INDEX idx_account_indices_deleted_at ON accounts(deleted_at);

    -- CREATE TABLE products (
    --     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    --     title VARCHAR NOT NULL,
    --     description VARCHAR NOT NULL,
    --     is_on_market BOOLEAN NOT NULL DEFAULT FALSE,
    --     creator_id UUID NOT NULL,
    --     created_at TIMESTAMP NOT NULL,
    --     updated_at TIMESTAMP NOT NULL,
    --     deleted_at TIMESTAMP,
    --     CONSTRAINT fk_creator FOREIGN KEY (creator_id) REFERENCES users(id)
    -- );

    -- CREATE INDEX idx_product_indices_deleted_at ON products(deleted_at);

    -- CREATE TABLE product_images (
    --     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    --     filename VARCHAR NOT NULL,
    --     path VARCHAR NOT NULL,
    --     product_id UUID NOT NULL,
    --     created_at TIMESTAMP NOT NULL,
    --     updated_at TIMESTAMP NOT NULL,
    --     deleted_at TIMESTAMP,
    --     CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id)
    -- );

    -- CREATE INDEX idx_product_image_indices_deleted_at ON product_images(deleted_at);

    -- CREATE TYPE offer_state AS ENUM (
    --     'CREATED',
    --     'CANCELLED',
    --     'TRANSFERRING_MONEY',
    --     'TRANSFER_MONEY_FAILED',
    --     'TRANSFERRING_PRODUCT',
    --     'TRANSFER_PRODUCT_FAILED',
    --     'SUCCEEDED',
    --     'RETURNING_MONEY',
    --     'RETURN_MONEY_FAILED',
    --     'MONEY_RETURNED'
    -- );

    -- CREATE TABLE offers (
    --     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    --     state offer_state NOT NULL DEFAULT 'CREATED',
    --     fail_message VARCHAR,
    --     delete_on_sell BOOLEAN NOT NULL DEFAULT TRUE, 
    --     product_id UUID NOT NULL,
    --     consumer_id UUID NOT NULL,
    --     created_at TIMESTAMP NOT NULL,
    --     updated_at TIMESTAMP NOT NULL,
    --     deleted_at TIMESTAMP,
    --     CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id),
    --     CONSTRAINT fk_consumer FOREIGN KEY (consumer_id) REFERENCES users(id)
    -- );

    -- CREATE INDEX idx_offer_indices_deleted_at ON offers(deleted_at);

    -- CREATE TYPE transaction_state AS ENUM (
    --     'CREATED',
    --     'CANCELLED',
    --     'PROCESSING',
    --     'ERROR',
    --     'SUCCEEDED',
    --     'FAILED'
    -- );

    -- CREATE TYPE transaction_type AS ENUM (
    --     'DEPOSIT',
    --     'BUY',
    --     'FEE',
    --     'WITHDRAWAL'
    -- );

    -- CREATE TYPE transaction_currency AS ENUM (
    --     'RUB',
    --     'EUR',
    --     'USD'
    -- );

    -- CREATE TABLE transactions (
    --     id SERIAL PRIMARY KEY,
    --     date TIMESTAMP,
    --     state transaction_state NOT NULL DEFAULT 'CREATED',
    --     type transaction_type NOT NULL,
    --     currency transaction_currency NOT NULL,
    --     amount DECIMAL(12, 2) NOT NULL,
    --     error VARCHAR,
    --     account_from_id UUID NOT NULL,
    --     account_to_id UUID NOT NULL,
    --     offer_id UUID,
    --     created_at TIMESTAMP NOT NULL,
    --     updated_at TIMESTAMP NOT NULL,
    --     deleted_at TIMESTAMP,
    --     CONSTRAINT fk_accunt_from FOREIGN KEY (account_from_id) REFERENCES accounts(id),
    --     CONSTRAINT fk_accunt_to FOREIGN KEY (account_to_id) REFERENCES accounts(id),
    --     CONSTRAINT fk_offer FOREIGN KEY (offer_id) REFERENCES offers(id)
    -- );

    -- CREATE INDEX idx_transaction_indices_deleted_at ON transactions(deleted_at);

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;
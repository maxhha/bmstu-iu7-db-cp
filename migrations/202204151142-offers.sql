--
-- Add offers table and fail_reason column in auctions
--
SELECT EXISTS (
    SELECT id FROM migrations WHERE id = :'MIGRATION_ID'
) as migrated \gset

\if :migrated
    \echo 'migration' :MIGRATION_ID 'already exists, skipping'
\else
    \echo 'migration' :MIGRATION_ID 'does not exist'

    CREATE TYPE offer_state AS ENUM (
        'CREATED',
        'CANCELLED',
        'PROCESSING',
        'SUCCEEDED',
        'FAILED'
    );

    CREATE TABLE offers (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        state offer_state NOT NULL DEFAULT 'CREATED',
        auction_id UUID NOT NULL,
        user_id SHORTKEY NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
        CONSTRAINT fk_auction FOREIGN KEY (auction_id) REFERENCES auctions(id),
        CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
    );

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;
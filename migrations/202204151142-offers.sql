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

    \set SERVER_USER `echo $SERVER_USER`
    \set VIEWER_USER `echo $VIEWER_USER`
    \set DEALER_USER `echo $DEALER_USER`

    \set exit_error false

    SELECT (:'SERVER_USER' = '') as is_not_empty \gset
    \if :is_not_empty
        \warn 'SERVER_USER is empty'
        \set exit_error true
    \endif

    SELECT (:'VIEWER_USER' = '') as is_not_empty \gset
    \if :is_not_empty
        \warn 'VIEWER_USER is empty'
        \set exit_error true
    \endif

    SELECT (:'DEALER_USER' = '') as is_not_empty \gset
    \if :is_not_empty
        \warn 'DEALER_USER is empty'
        \set exit_error true
    \endif

    \if :exit_error 
        DO $$
        BEGIN
        RAISE EXCEPTION 'all required environment variables must not be empty';
        END;
        $$;
    \endif

    CREATE TYPE offer_state AS ENUM (
        'CREATED',
        'CANCELLED',
        'ACCEPTED',
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

    GRANT SELECT, INSERT, UPDATE 
        ON offers
        TO :SERVER_USER, :DEALER_USER;

    GRANT SELECT, UPDATE 
        ON offers
        TO :DEALER_USER;

    GRANT SELECT
        ON offers
        TO :VIEWER_USER;

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;
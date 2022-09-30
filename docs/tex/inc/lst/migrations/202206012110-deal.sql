--
-- Adds deals
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

    CREATE TYPE deal_state AS ENUM (
        'TRANSFERRING_MONEY',
        'TRANSFER_MONEY_FAILED',
        'TRANSFERRING_PRODUCT',
        'TRANSFER_PRODUCT_FAILED',
        'SUCCEEDED',
        'RETURNING_MONEY',
        'RETURN_MONEY_FAILED',
        'MONEY_RETURNED'
    );

    CREATE TABLE deal_states (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        state deal_state NOT NULL DEFAULT 'TRANSFERRING_MONEY',
        creator_id SHORTKEY,
        offer_id UUID NOT NULL,
        comment TEXT,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        CONSTRAINT fk_creator FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE SET NULL,
        CONSTRAINT fk_offer FOREIGN KEY (offer_id) REFERENCES offers(id) ON DELETE CASCADE
    );

    GRANT SELECT, INSERT, UPDATE 
        ON deal_states
        TO :SERVER_USER, :DEALER_USER;

    GRANT SELECT
        ON deal_states
        TO :VIEWER_USER;

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;

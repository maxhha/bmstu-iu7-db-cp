--
-- Add transactions table
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

    \if :exit_error 
        DO $$
        BEGIN
        RAISE EXCEPTION 'all required environment variables must not be empty';
        END;
        $$;
    \endif

    CREATE TYPE transaction_state AS ENUM (
        'CREATED',
        'CANCELLED',
        'NEW',
        'PROCESSING',
        'ERROR',
        'SUCCEEDED',
        'FAILED'
    );

    CREATE TYPE transaction_type AS ENUM (
        'DEPOSIT',
        'CURRENCY_CONVERTION',
        'RETURN',
        'FEE_RETURN',
        'BUY',
        'FEE_BUY',
        'WITHDRAWAL',
        'FEE_WITHDRAWAL'
    );

    CREATE TABLE transactions (
        id SERIAL PRIMARY KEY,
        date TIMESTAMP,
        state transaction_state NOT NULL DEFAULT 'CREATED',
        type transaction_type NOT NULL,
        currency currency NOT NULL,
        amount DECIMAL(12, 2) NOT NULL CHECK (amount > 0),
        error VARCHAR,
        account_from_id UUID,
        account_to_id UUID,
        offer_id UUID,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted_at TIMESTAMP,
        CONSTRAINT fk_accunt_from FOREIGN KEY (account_from_id) REFERENCES accounts(id),
        CONSTRAINT fk_accunt_to FOREIGN KEY (account_to_id) REFERENCES accounts(id),
        CONSTRAINT fk_offer FOREIGN KEY (offer_id) REFERENCES offers(id),
        CONSTRAINT chk_account_from CHECK (
            CASE WHEN type IN ('DEPOSIT') 
                THEN account_from_id IS NULL 
            WHEN type IN (
                'CURRENCY_CONVERTION',
                'RETURN',
                'FEE_RETURN',
                'BUY',
                'FEE_BUY',
                'WITHDRAWAL',
                'FEE_WITHDRAWAL'
            )
                THEN account_from_id IS NOT NULL
            ELSE false
            END
        ),
        CONSTRAINT chk_account_to CHECK (
            CASE WHEN type IN ('WITHDRAWAL') 
                THEN account_to_id IS NULL 
            WHEN type IN (
                'DEPOSIT',
                'CURRENCY_CONVERTION',
                'RETURN',
                'FEE_RETURN',
                'BUY',
                'FEE_BUY',
                'FEE_WITHDRAWAL'
            )
                THEN account_to_id IS NOT NULL
            ELSE false
            END
        ),
        CONSTRAINT chk_offer CHECK (
            CASE WHEN type IN ('DEPOSIT', 'CURRENCY_CONVERTION') 
                THEN offer_id IS NULL 
            WHEN type IN (
                'RETURN',
                'FEE_RETURN',
                'BUY',
                'FEE_BUY', 
                'WITHDRAWAL',
                'FEE_WITHDRAWAL'
            )
                THEN offer_id IS NOT NULL
            ELSE false
            END
        ),
        CONSTRAINT chk_date CHECK (
            state <> 'SUCCEEDED' OR date IS NOT NULL
        )
    );

    CREATE INDEX idx_transaction_indices_deleted_at ON transactions(deleted_at);

    GRANT SELECT, INSERT, UPDATE 
        ON transactions
        TO :SERVER_USER;
    
    GRANT USAGE, SELECT
        ON SEQUENCE transactions_id_seq
        TO :SERVER_USER;

    GRANT SELECT
        ON transactions
        TO :VIEWER_USER;

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;
--
-- Create transfers and transfer_algs tables
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

    CREATE TABLE transfer_algs (
        id SERIAL PRIMARY KEY,
        name VARCHAR NOT NULL,
        type VARCHAR NOT NULL,
        params JSONB,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMP
    );

    CREATE INDEX idx_transfer_alg_indices_deleted_at ON transfer_algs(deleted_at);

    CREATE TABLE transfers (
        id SERIAL PRIMARY KEY,
        currency_from currency NOT NULL,
        currency_to currency NOT NULL,
        account_from_id UUID NOT NULL,
        account_to_id UUID NOT NULL,
        alg_id INT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMP,
        CONSTRAINT fk_accunt_from FOREIGN KEY (account_from_id) REFERENCES nominal_accounts(id) ON DELETE CASCADE,
        CONSTRAINT fk_accunt_to FOREIGN KEY (account_to_id) REFERENCES nominal_accounts(id) ON DELETE CASCADE,
        CONSTRAINT fk_transfer_alg FOREIGN KEY (alg_id) REFERENCES transfer_algs(id)
    );

    CREATE INDEX idx_transfer_indices_deleted_at ON transfers(deleted_at);

    GRANT SELECT, INSERT, UPDATE 
        ON transfers, transfer_algs
        TO :SERVER_USER;

    GRANT USAGE, SELECT
        ON SEQUENCE transfers_id_seq, transfer_algs_id_seq
        TO :SERVER_USER;

    GRANT SELECT
        ON transfers, transfer_algs
        TO :VIEWER_USER;

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;
--
-- Grant permissions on transactions table for bankgate
--
SELECT EXISTS (
    SELECT id FROM migrations WHERE id = :'MIGRATION_ID'
) as migrated \gset

\if :migrated
    \echo 'migration' :MIGRATION_ID 'already exists, skipping'
\else
    \echo 'migration' :MIGRATION_ID 'does not exist'

    \set BANKGATE_USER `echo $BANKGATE_USER`

    \set exit_error false

    SELECT (:'BANKGATE_USER' = '') as is_not_empty \gset
    \if :is_not_empty
        \warn 'BANKGATE_USER is empty'
        \set exit_error true
    \endif

    \if :exit_error 
        DO $$
        BEGIN
        RAISE EXCEPTION 'all required environment variables must not be empty';
        END;
        $$;
    \endif

    GRANT SELECT, INSERT, UPDATE 
        ON transactions
        TO :BANKGATE_USER;

    GRANT USAGE, SELECT
        ON SEQUENCE transactions_id_seq
        TO :BANKGATE_USER;

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;
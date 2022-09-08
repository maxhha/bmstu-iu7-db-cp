--
-- Create database roles for server, bankgate and viewer
--
SELECT EXISTS (
    SELECT id FROM migrations WHERE id = :'MIGRATION_ID'
) as migrated \gset

\if :migrated
    \echo 'migration' :MIGRATION_ID 'already exists, skipping'
\else
    \echo 'migration' :MIGRATION_ID 'does not exist'

    \set DEALER_USER `echo $DEALER_USER`

    \set exit_error false

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

    GRANT SELECT, INSERT, UPDATE 
        ON ALL TABLES IN SCHEMA public
        TO :DEALER_USER;

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;
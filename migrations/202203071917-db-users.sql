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

    \set SERVER_USER `echo $SERVER_USER`
    \set SERVER_PASSWORD `echo $SERVER_PASSWORD`
    \set BANKGATE_USER `echo $BANKGATE_USER`
    \set BANKGATE_PASSWORD `echo $BANKGATE_PASSWORD`
    \set VIEWER_USER `echo $VIEWER_USER`
    \set VIEWER_PASSWORD `echo $VIEWER_PASSWORD`

    \set exit_error false

    SELECT (:'SERVER_USER' = '') as is_not_empty \gset
    \if :is_not_empty
        \warn 'SERVER_USER is empty'
        \set exit_error true
    \endif

    SELECT (:'SERVER_PASSWORD' = '') as is_not_empty \gset
    \if :is_not_empty
        \warn 'SERVER_PASSWORD is empty'
        \set exit_error true
    \endif

    SELECT (:'BANKGATE_USER' = '') as is_not_empty \gset
    \if :is_not_empty
        \warn 'BANKGATE_USER is empty'
        \set exit_error true
    \endif

    SELECT (:'BANKGATE_PASSWORD' = '') as is_not_empty \gset
    \if :is_not_empty
        \warn 'BANKGATE_PASSWORD is empty'
        \set exit_error true
    \endif

    SELECT (:'VIEWER_USER' = '') as is_not_empty \gset
    \if :is_not_empty
        \warn 'VIEWER_USER is empty'
        \set exit_error true
    \endif

    SELECT (:'VIEWER_PASSWORD' = '') as is_not_empty \gset
    \if :is_not_empty
        \warn 'VIEWER_PASSWORD is empty'
        \set exit_error true
    \endif

    \if :exit_error 
        DO $$
        BEGIN
        RAISE EXCEPTION 'all required environment variables must not be empty';
        END;
        $$;
    \endif

    CREATE ROLE :SERVER_USER
        WITH LOGIN PASSWORD :'SERVER_PASSWORD';
    GRANT SELECT, INSERT, UPDATE 
        ON ALL TABLES IN SCHEMA public
        TO :SERVER_USER;

    CREATE ROLE :VIEWER_USER 
        WITH LOGIN PASSWORD :'VIEWER_PASSWORD';
    GRANT SELECT
        ON ALL TABLES IN SCHEMA public
        TO :VIEWER_USER;

    -- CREATE ROLE :BANKGATE_USER
    --     WITH LOGIN PASSWORD :'BANKGATE_PASSWORD'
    --     IN ROLE :VIEWER_USER;
    -- GRANT INSERT, UPDATE
    --     ON TABLE transactions, banks
    --     TO :BANKGATE_USER;

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;
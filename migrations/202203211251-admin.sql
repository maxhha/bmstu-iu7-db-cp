--
-- Add admin user in project
--
SELECT EXISTS (
    SELECT id FROM migrations WHERE id = :'MIGRATION_ID'
) as migrated \gset

\if :migrated
    \echo 'migration' :MIGRATION_ID 'already exists, skipping'
\else
    \echo 'migration' :MIGRATION_ID 'does not exist'

    \set PROJECT_ADMIN_EMAIL `echo $PROJECT_ADMIN_EMAIL`
    \set PROJECT_ADMIN_PASSWORD `echo $PROJECT_ADMIN_PASSWORD`

    \set exit_error false

    SELECT (:'PROJECT_ADMIN_EMAIL' = '') as is_not_empty \gset
    \if :is_not_empty
        \warn 'PROJECT_ADMIN_EMAIL is empty'
        \set exit_error true
    \endif

    SELECT (:'PROJECT_ADMIN_PASSWORD' = '') as is_not_empty \gset
    \if :is_not_empty
        \warn 'PROJECT_ADMIN_PASSWORD is empty'
        \set exit_error true
    \endif

    \if :exit_error 
        DO $$
        BEGIN
        RAISE EXCEPTION 'all required environment variables must not be empty';
        END;
        $$;
    \endif

    INSERT INTO users (deleted_at, blocked_until)
    VALUES (NULL, NULL) RETURNING id AS admin_id \gset

    INSERT INTO user_forms (user_id, state, email, password, name)
    VALUES (:'admin_id', 'APPROVED', :'PROJECT_ADMIN_EMAIL', :'PROJECT_ADMIN_PASSWORD', 'admin');

    INSERT INTO roles (type, user_id, issuer_id)
    VALUES ('MANAGER', :'admin_id', :'admin_id'),
        ('ADMIN', :'admin_id', :'admin_id');

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;
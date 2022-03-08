--
-- Add is_admin column in user table
--
SELECT EXISTS (
    SELECT id FROM migrations WHERE id = :'MIGRATION_ID'
) as migrated \gset

\if :migrated
    \echo 'migration' :MIGRATION_ID 'already exists, skipping'
\else
    \echo 'migration' :MIGRATION_ID 'does not exist'

    ALTER TABLE users ADD is_admin BOOLEAN NOT NULL DEFAULT false;

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;
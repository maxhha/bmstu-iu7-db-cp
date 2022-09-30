--
-- Add declain_reason to products table
--
SELECT EXISTS (
    SELECT id FROM migrations WHERE id = :'MIGRATION_ID'
) as migrated \gset

\if :migrated
    \echo 'migration' :MIGRATION_ID 'already exists, skipping'
\else
    \echo 'migration' :MIGRATION_ID 'does not exist'

    ALTER TABLE products ADD COLUMN declain_reason TEXT;

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;
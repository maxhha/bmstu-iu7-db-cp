--
-- Add token action MODERATE_PRODUCT
--
SELECT EXISTS (
    SELECT id FROM migrations WHERE id = :'MIGRATION_ID'
) as migrated \gset

\if :migrated
    \echo 'migration' :MIGRATION_ID 'already exists, skipping'
\else
    \echo 'migration' :MIGRATION_ID 'does not exist'

    ALTER TYPE token_action ADD VALUE 'MODERATE_PRODUCT';

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;
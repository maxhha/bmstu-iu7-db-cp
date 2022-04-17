--
-- Creates fake bank
--
SELECT EXISTS (
    SELECT id FROM migrations WHERE id = :'MIGRATION_ID'
) as migrated \gset

\if :migrated
    \echo 'migration' :MIGRATION_ID 'already exists, skipping'
\else
    \echo 'migration' :MIGRATION_ID 'does not exist'

    INSERT INTO banks (name)
    VALUES ('fake') RETURNING id as bank_id \gset

    INSERT INTO accounts (type, bank_id)
    VALUES ('BANK', :'bank_id');

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;
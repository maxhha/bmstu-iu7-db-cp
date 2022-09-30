--
-- Add seller account column to auctions
--
SELECT EXISTS (
    SELECT id FROM migrations WHERE id = :'MIGRATION_ID'
) as migrated \gset

\if :migrated
    \echo 'migration' :MIGRATION_ID 'already exists, skipping'
\else
    \echo 'migration' :MIGRATION_ID 'does not exist'

    ALTER TABLE auctions ADD COLUMN seller_account_id UUID;

    ALTER TABLE auctions ADD CONSTRAINT fk_seller_account 
        FOREIGN KEY (seller_account_id) REFERENCES accounts(id);

    ALTER TABLE auctions ADD CONSTRAINT chk_seller_account CHECK (
        state = 'CREATED' OR seller_account_id IS NOT NULL
    );

    INSERT INTO migrations(id) VALUES (:'MIGRATION_ID');
\endif

COMMIT;
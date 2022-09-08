      OR (
        transactions.type IN ('BUY') 
        AND transactions.state IN ('CREATED')
      )
    ) 
  WHERE "transactions"."deleted_at" IS NULL 
  UNION ALL 
  SELECT currency, account_to_id as account_id, amount 
  FROM "transactions" 
  JOIN (
    SELECT * 
    FROM "accounts" 
    WHERE id = 'a5c4d021-67b8-4bf5-8a22-8869523c9ab9' 
    AND "accounts"."deleted_at" IS NULL
  ) a ON account_to_id = a.id 
  AND transactions.state IN (
    'SUCCEEDED', 'PROCESSING', 'ERROR'
  ) 
  WHERE "transactions"."deleted_at" IS NULL
) trs 
GROUP BY trs.currency, trs.account_id

-- Создание предложения
INSERT INTO "offers" (
  "state", "auction_id", "user_id", 
  "created_at", "updated_at"
)
VALUES (
  'CREATED', 'a75d4700-7058-4749-b1e4-8aef840103fc', 
  'zTUZBNT-9b0', '2022-06-18 09:56:16.263', 
  '2022-06-18 09:56:16.263'
)
RETURNING "id";
-- Создание транзакция предложения
INSERT INTO "transactions" (
  "date", "state", "type", "currency", 
  "amount", "error", "account_from_id", 
  "account_to_id", "offer_id", "created_at", 
  "updated_at", "deleted_at"
) 
VALUES (
  NULL, 'CREATED', 'BUY', 'RUB', '10', 
  NULL, 'a5c4d021-67b8-4bf5-8a22-8869523c9ab9', 
  '5d0c6a55-23cd-4e70-b342-5707578a92cb', 
  '31e167d9-4341-4b47-8922-e8b5281fbf21', 
  '2022-06-18 09:56:16.268', '2022-06-18 09:56:16.268', 
  NULL
)
RETURNING "id";
COMMIT;

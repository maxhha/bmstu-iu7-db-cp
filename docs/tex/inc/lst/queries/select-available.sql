SELECT trs.currency, trs.account_id, SUM(trs.amount) as amount 
FROM (
  SELECT currency, account_from_id as account_id, -amount as amount 
  FROM "transactions" 
  JOIN (
      SELECT * 
      FROM "accounts" 
      WHERE id = 'test-account-id' 
      AND "accounts"."deleted_at" IS NULL
  ) a ON account_from_id = a.id 
    AND (
      transactions.state IN (
        'SUCCEEDED', 'PROCESSING', 'ERROR'
      )
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
    WHERE id = 'test-account-id' 
    AND "accounts"."deleted_at" IS NULL
  ) a ON account_to_id = a.id 
  AND transactions.state IN (
    'SUCCEEDED', 'PROCESSING', 'ERROR'
  ) 
  WHERE "transactions"."deleted_at" IS NULL
) trs 
GROUP BY trs.currency, trs.account_id

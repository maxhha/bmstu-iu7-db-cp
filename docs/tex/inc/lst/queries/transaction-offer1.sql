-- Текущий пользователь
SELECT * FROM "users" 
WHERE id = 'zTUZBNT-9b0' AND "users"."deleted_at" IS NULL;
-- Аукцион, в котором планируется сделать ставку
SELECT * FROM "auctions" 
WHERE id = 'a75d4700-7058-4749-b1e4-8aef840103fc';
-- Счет, с которого будет списана сумма
SELECT * FROM "accounts" 
WHERE id = 'a5c4d021-67b8-4bf5-8a22-8869523c9ab9' 
AND "accounts"."deleted_at" IS NULL;

BEGIN TRANSACTION;
-- Блокировка аукциона 
SELECT * FROM "auctions" 
WHERE id = 'a75d4700-7058-4749-b1e4-8aef840103fc' 
FOR SHARE OF "auctions";
-- Блокировка счета
SELECT * FROM "accounts" 
WHERE id = 'a5c4d021-67b8-4bf5-8a22-8869523c9ab9' 
AND "accounts"."deleted_at" IS NULL 
FOR UPDATE OF "accounts";
-- Отмена прошлого предложения, если оно есть
SELECT * FROM "offers" 
WHERE auction_id IN ('a75d4700-7058-4749-b1e4-8aef840103fc')
AND state IN ('CREATED') 
AND user_id IN ('zTUZBNT-9b0');
-- Здесь его не оказалось, поэтому отмены нет.
-- Получение баланса
SELECT trs.currency, trs.account_id, SUM(trs.amount) as amount 
FROM (
  SELECT currency, account_from_id as account_id, -amount as amount 
  FROM "transactions" 
  JOIN (
      SELECT * 
      FROM "accounts" 
      WHERE id = 'a5c4d021-67b8-4bf5-8a22-8869523c9ab9' 
      AND "accounts"."deleted_at" IS NULL
  ) a ON account_from_id = a.id 
    AND (
      transactions.state IN (
        'SUCCEEDED', 'PROCESSING', 'ERROR'
      )
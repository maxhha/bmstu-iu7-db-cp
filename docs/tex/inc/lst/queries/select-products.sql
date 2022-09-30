SELECT *
FROM "products" 
JOIN (
  SELECT *, 
    ROW_NUMBER() OVER(
      PARTITION BY ofd.product_id 
      ORDER BY ofd.from_date DESC
    ) as owner_n 
  FROM (
    SELECT id as product_id, creator_id as owner_id, created_at as from_date 
    FROM "products" 
    WHERE "products"."deleted_at" IS NULL 
    UNION ALL 
    SELECT products.id as product_id, auctions.buyer_id as owner_id, auctions.finished_at as from_date 
    FROM "products" 
    JOIN auctions ON auctions.product_id = products.id 
    AND auctions.state = 'SUCCEEDED' 
    WHERE "products"."deleted_at" IS NULL
  ) as ofd
) ofd ON products.id = ofd.product_id 
AND ofd.owner_n = 1 
AND ofd.owner_id IN ('test-owner-id') 
WHERE "products"."deleted_at" IS NULL 
ORDER BY created_at desc

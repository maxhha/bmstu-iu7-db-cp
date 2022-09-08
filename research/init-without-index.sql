CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS user_forms;
DROP TYPE IF EXISTS user_form_state;
DROP TYPE IF EXISTS currency;

DROP TRIGGER IF EXISTS trgr_users_genid ON users;
DROP TABLE IF EXISTS users;
DROP DOMAIN IF EXISTS SHORTKEY;

-- shorkey used for prettier urls
CREATE DOMAIN SHORTKEY as varchar(11);

CREATE OR REPLACE FUNCTION shortkey_generate()
RETURNS TRIGGER AS $$
DECLARE
    gkey TEXT;
    key SHORTKEY;
    qry TEXT;
    found TEXT;
    user_id BOOLEAN;
BEGIN
    -- generate the first part of a query as a string with safely
    -- escaped table name, using || to concat the parts
    qry := 'SELECT id FROM ' || quote_ident(TG_TABLE_NAME) || ' WHERE id=';

    LOOP
        -- deal with user-supplied keys, they don't have to be valid base64
        -- only the right length for the type
        IF NEW.id IS NOT NULL THEN
            key := NEW.id;
            user_id := TRUE;

            IF length(key) <> 11 THEN
                RAISE 'User defined key value % has invalid length. Expected 11, got %.', key, length(key);
            END IF;
        ELSE
            -- 8 bytes gives a collision p = .5 after 5.1 x 10^9 values
            gkey := encode(gen_random_bytes(8), 'base64');
            gkey := replace(gkey, '/', '_');  -- url safe replacement
            gkey := replace(gkey, '+', '-');  -- url safe replacement
            key := rtrim(gkey, '=');          -- cut off padding
            user_id := FALSE;
        END IF;

        -- Concat the generated key (safely quoted) with the generated query
        -- and run it.
        -- SELECT id FROM "test" WHERE id='blahblah' INTO found
        -- Now "found" will be the duplicated id or NULL.
        EXECUTE qry || quote_literal(key) INTO found;

        -- Check to see if found is NULL.
        -- If we checked to see if found = NULL it would always be FALSE
        -- because (NULL = NULL) is always FALSE.
        IF found IS NULL THEN
            -- If we didn't find a collision then leave the LOOP.
            EXIT;
        END IF;

        IF user_id THEN
            -- User supplied ID but it violates the PK unique constraint
            RAISE 'ID % already exists in table %', key, TG_TABLE_NAME;
        END IF;

        -- We haven't EXITed yet, so return to the top of the LOOP
        -- and try again.
    END LOOP;

    -- NEW and OLD are available in TRIGGER PROCEDURES.
    -- NEW is the mutated row that will actually be INSERTed.
    -- We're replacing id, regardless of what it was before
    -- with our key variable.
    NEW.id = key;

    -- The RECORD returned here is what will actually be INSERTed,
    -- or what the next trigger will get if there is one.
    RETURN NEW;
END
$$ language 'plpgsql';

CREATE TABLE users (
    id SHORTKEY PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- generate shortkey for each insert
CREATE TRIGGER trgr_users_genid 
    BEFORE INSERT ON users FOR EACH ROW
    EXECUTE PROCEDURE shortkey_generate();

CREATE TYPE user_form_state AS ENUM (
    'CREATED',
    'MODERATING',
    'APPROVED',
    'DECLAINED'
);

CREATE TYPE currency AS ENUM (
    'RUB',
    'EUR',
    'USD'
);

-- represents user data
CREATE TABLE user_forms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id SHORTKEY NOT NULL,
    state user_form_state NOT NULL DEFAULT 'CREATED',
    name VARCHAR,
    password VARCHAR,
    phone VARCHAR,
    email VARCHAR,
    currency currency,
    declain_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

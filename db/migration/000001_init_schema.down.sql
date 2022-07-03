ALTER TABLE IF EXISTS "transfers" DROP CONSTRAINT IF EXISTS "transfer_accounts1key";

ALTER TABLE IF EXISTS "transfers" DROP CONSTRAINT IF EXISTS "transfer_accounts2key";

ALTER TABLE IF EXISTS "entries" DROP CONSTRAINT IF EXISTS "entries_accounts_key";

DROP TABLE IF EXISTS "entries";

DROP TABLE IF EXISTS "transfers";

DROP TABLE IF EXISTS "accounts";
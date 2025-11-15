DROP TRIGGER IF EXISTS trg_set_transactions_updated_at ON transactions;
DROP FUNCTION IF EXISTS set_transactions_updated_at();

ALTER TABLE transactions
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS updated_at;
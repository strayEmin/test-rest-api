ALTER TABLE transactions
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE OR REPLACE FUNCTION set_transaction_updated_at()
RETURNS TRIGGER AS $$
BEGIN
     NEW.updated_at = NOW();
     RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_set_transactions_updated_at
BEFORE UPDATE ON transactions
FOR EACH ROW
EXECUTE FUNCTION set_transaction_updated_at();

BEGIN;

ALTER TABLE models ADD COLUMN IF NOT EXISTS updated timestamp NOT NULL DEFAULT current_timestamp;

CREATE OR REPLACE FUNCTION models_updated() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated = current_timestamp;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS models_updated_trigger ON models;
CREATE TRIGGER models_updated_trigger
BEFORE INSERT OR UPDATE ON models
FOR EACH ROW EXECUTE FUNCTION models_updated();

END;

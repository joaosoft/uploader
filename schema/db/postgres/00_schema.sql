-- migrate up
CREATE OR REPLACE FUNCTION function_updated_at()
  RETURNS TRIGGER AS $$
  BEGIN
   NEW.updated_at = now();
   RETURN NEW;
  END;
  $$ LANGUAGE 'plpgsql';


CREATE TABLE upload (
  id_upload         TEXT NOT NULL,
  name              TEXT NOT NULL,
  file              BYTEA NOT NULL,
  created_at        TIMESTAMP DEFAULT NOW(),
  updated_at        TIMESTAMP DEFAULT NOW(),
  CONSTRAINT upload_id_upload_pkey PRIMARY KEY (id_upload)
);


CREATE TRIGGER trigger_upload_updated_at BEFORE UPDATE
  ON upload FOR EACH ROW EXECUTE PROCEDURE function_updated_at();


-- migrate down
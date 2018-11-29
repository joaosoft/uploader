-- migrate up
CREATE SCHEMA IF NOT EXISTS uploader;


CREATE OR REPLACE FUNCTION uploader.function_updated_at()
  RETURNS TRIGGER AS $$
  BEGIN
   NEW.updated_at = now();
   RETURN NEW;
  END;
  $$ LANGUAGE 'plpgsql';


CREATE TABLE uploader.upload (
  id_upload         TEXT NOT NULL,
  name              TEXT NOT NULL,
  file              BYTEA NOT NULL,
  created_at        TIMESTAMP DEFAULT NOW(),
  updated_at        TIMESTAMP DEFAULT NOW(),
  CONSTRAINT upload_id_upload_pkey PRIMARY KEY (id_upload)
);


CREATE TRIGGER trigger_upload_updated_at BEFORE UPDATE
  ON uploader.upload FOR EACH ROW EXECUTE PROCEDURE uploader.function_updated_at();


-- migrate down
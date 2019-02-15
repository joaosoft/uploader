-- migrate up

CREATE SCHEMA IF NOT EXISTS uploader;

-- upload configuration
CREATE TABLE uploader.section (
	id_section        SERIAL PRIMARY KEY,
	"name" 		      TEXT NOT NULL,
	"path"            TEXT NOT NULL,
	active            BOOLEAN DEFAULT TRUE NOT NULL
);

CREATE TABLE uploader.image_size (
	id_image_size     SERIAL PRIMARY KEY,
	"name" 		      TEXT NOT NULL,
	"path"            TEXT NOT NULL,
	width             SMALLINT NOT NULL,
	height            SMALLINT NOT NULL,
	active            BOOLEAN DEFAULT TRUE NOT NULL
);

CREATE TABLE uploader.section_image_size (
	fk_section        INTEGER REFERENCES uploader.section (id_section) NOT NULL,
	fk_image_size 	  INTEGER REFERENCES uploader.image_size (id_image_size) NOT NULL,
	active            BOOLEAN DEFAULT TRUE NOT NULL,
	CONSTRAINT        section_image_size_pkey PRIMARY KEY (fk_section, fk_image_size)
);






-- table to insert images
CREATE OR REPLACE FUNCTION uploader.function_updated_at()
  RETURNS TRIGGER AS $$
  BEGIN
   NEW.updated_at = now();
   RETURN NEW;
  END;
  $$ LANGUAGE 'plpgsql';


CREATE TABLE uploader.upload (
  id_upload         TEXT NOT NULL,
  file              BYTEA NOT NULL,
  created_at        TIMESTAMP DEFAULT NOW(),
  updated_at        TIMESTAMP DEFAULT NOW(),
  CONSTRAINT upload_id_upload_pkey PRIMARY KEY (id_upload)
);


CREATE TRIGGER trigger_upload_updated_at BEFORE UPDATE
  ON uploader.upload FOR EACH ROW EXECUTE PROCEDURE uploader.function_updated_at();


-- migrate down
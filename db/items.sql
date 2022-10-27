-- Table: rfq.items

-- DROP TABLE IF EXISTS rfq.items;

CREATE TABLE IF NOT EXISTS rfq.items
(
    "HeaderID" integer NOT NULL,
    "Nama" text COLLATE pg_catalog."default" NOT NULL,
    "Harga" integer NOT NULL,
    "Qty" integer NOT NULL
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS rfq.items
    OWNER to postgres;
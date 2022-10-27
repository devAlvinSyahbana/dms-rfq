-- Table: rfq.header

-- DROP TABLE IF EXISTS rfq.header;

CREATE TABLE IF NOT EXISTS rfq.header
(
    "CompanyName" text COLLATE pg_catalog."default",
    "CompanyAddress" text COLLATE pg_catalog."default",
    "CompanyWebsite" text COLLATE pg_catalog."default",
    "QuotationDate" text COLLATE pg_catalog."default",
    "QuotationNo" text COLLATE pg_catalog."default",
    "QuotationExpires" text COLLATE pg_catalog."default",
    "MadeForName" text COLLATE pg_catalog."default",
    "MadeForAddress" text COLLATE pg_catalog."default",
    "MadeForPhone" text COLLATE pg_catalog."default",
    "SentToName" text COLLATE pg_catalog."default",
    "SentToAddress" text COLLATE pg_catalog."default",
    "SentToPhone" text COLLATE pg_catalog."default",
    "Disc" text COLLATE pg_catalog."default",
    "Tax" text COLLATE pg_catalog."default",
    "Interest" text COLLATE pg_catalog."default",
    "SNK" text[] COLLATE pg_catalog."default",
    "ID" integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    CONSTRAINT id PRIMARY KEY ("ID")
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS rfq.header
    OWNER to postgres;
-- Table: rfq.user

-- DROP TABLE IF EXISTS rfq."user";

CREATE TABLE IF NOT EXISTS rfq."user"
(
    email text COLLATE pg_catalog."default",
    password text COLLATE pg_catalog."default"
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS rfq."user"
    OWNER to postgres;
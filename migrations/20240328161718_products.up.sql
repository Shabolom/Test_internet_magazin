CREATE TABLE IF NOT EXISTS public.products
(
    id integer NOT NULL,
    name text COLLATE pg_catalog."default",
    CONSTRAINT products_pkey PRIMARY KEY (id)
    );

CREATE UNIQUE INDEX idx_products_id
    ON public.products (id);

ALTER TABLE IF EXISTS public.products
    OWNER to postgres;
CREATE TABLE IF NOT EXISTS public.orders
(
    id integer NOT NULL,
    CONSTRAINT orders_pkey PRIMARY KEY (id)
    );

CREATE UNIQUE INDEX idx_orders_id
    ON public.orders (id);

ALTER TABLE IF EXISTS public.orders
    OWNER to postgres;
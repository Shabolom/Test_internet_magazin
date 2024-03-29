CREATE TABLE IF NOT EXISTS public.orders_products
(
    id  integer NOT NULL,
    order_id integer NOT NULL REFERENCES public.orders(id),
    product_id integer NOT NULL REFERENCES public.products(id),
    palette_id integer NOT NULL REFERENCES public.palettes(id),
    product_count integer NOT NULL,
    CONSTRAINT product_order_unique UNIQUE (product_id, order_id, palette_id),
    CONSTRAINT orders_products_pk PRIMARY KEY (id)
    );

CREATE INDEX idx_orders_products_order_id
    ON public.orders_products (order_id);

ALTER TABLE IF EXISTS public.orders_products
    OWNER to postgres;
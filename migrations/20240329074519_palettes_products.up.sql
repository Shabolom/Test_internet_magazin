CREATE TABLE IF NOT EXISTS public.palettes_products
(
    id  integer NOT NULL,
    palette_id integer NOT NULL REFERENCES public.palettes(id),
    product_id integer NOT NULL REFERENCES public.products(id),
    product_count integer NOT NULL ,
    palette_status bool,
    CONSTRAINT product_palette_unique UNIQUE (palette_id, product_id),
    CONSTRAINT products_palettes_pk PRIMARY KEY (id)
    );

CREATE INDEX idx_palettes_products_palette_id
    ON public.palettes_products (palette_id);

ALTER TABLE IF EXISTS public.palettes_products
    OWNER to postgres;
CREATE TABLE IF NOT EXISTS public.palettes
(
    id integer NOT NULL,
    name text COLLATE pg_catalog."default",
    CONSTRAINT palettes_pk PRIMARY KEY (id)
    );

CREATE UNIQUE INDEX idx_palettes
    ON public.palettes (id);

ALTER TABLE IF EXISTS public.palettes
    OWNER to postgres;
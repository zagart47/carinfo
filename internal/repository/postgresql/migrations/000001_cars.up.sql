CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE IF NOT EXISTS public.cars
(
    id               SERIAL      NOT NULL PRIMARY KEY,
    mark             VARCHAR,
    model            VARCHAR,
    year             INT,
    reg_num          VARCHAR,
    owner_name       VARCHAR,
    owner_surname    VARCHAR,
    owner_patronymic VARCHAR,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON public.cars
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp()



CREATE TABLE securebin (
    id bigserial not null primary key,
    img bytea not null,
    encrypted_password varchar not null,
    lifetime timestamp without time zone DEFAULT (CURRENT_TIMESTAMP + '1 day'::interval) NOT NULL
);

CREATE TABLE public.requests
(
	id bigserial NOT NULL PRIMARY KEY,
	data_id bigint NOT NULL,
    remote_addr cidr NOT NULL,
    "time" timestamp without time zone,
    CONSTRAINT req_id FOREIGN KEY (data_id)
        REFERENCES public.securebin (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
);

CREATE FUNCTION public.del_expired()
    RETURNS trigger
    LANGUAGE 'plpgsql'
     NOT LEAKPROOF
AS $BODY$BEGIN
DELETE FROM securebin WHERE lifetime < CURRENT_TIMESTAMP;
RETURN NULL;
END$BODY$;

CREATE TRIGGER del_exp
    AFTER INSERT
    ON public.securebin
    FOR EACH STATEMENT
    EXECUTE PROCEDURE public.del_expired();
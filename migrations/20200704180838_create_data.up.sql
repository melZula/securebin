CREATE TABLE securebin (
    id bigserial not null primary key,
    img bytea not null,
    encrypted_password varchar not null
);
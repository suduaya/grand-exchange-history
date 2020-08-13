CREATE TABLE ge.items(
    id integer,
    name text,
    PRIMARY KEY (id,name)
);

ALTER TABLE ge.items
    OWNER to postgres;
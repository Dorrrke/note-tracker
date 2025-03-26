CREATE TABLE IF NOT EXISTS users(
    uid varchar(36) NOT NULL PRIMARY KEY,
    name text NOT NULL,
    login text NOT NULL,
    password text NOT NULL,
    UNIQUE (login)
);

CREATE TABLE IF NOT EXISTS tasks(
    tid varchar(36) NOT NULL PRIMARY KEY,
    title text NOT NULL,
    description text NOT NULL,
    status text NOT NULL,
    created_at timestamp NOT NULL default NOW(),
    updated_at timestamp NOT NULL default NOW(),
    done_at timestamp,
    UNIQUE (title)
);


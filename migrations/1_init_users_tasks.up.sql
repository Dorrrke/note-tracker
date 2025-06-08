CREATE TABLE IF NOT EXISTS users (
    uid VARCHAR(36) NOT NULL PRIMARY KEY,
    name text NOT NULL,
    login text UNIQUE NOT NULL,
    email text UNIQUE NOT NULL,
    password text NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks(
    tid VARCHAR(36) NOT NULL PRIMARY KEY,
    title text UNIQUE NOT NULL,
    description text NOT NULL,
    status text NOT NULL,
    created_at timestamp NOT NULL default NOW(),
    updated_at timestamp NOT NULL default NOW(),
    done_at timestamp
);

CREATE TABLE IF NOT EXISTS users(
    uid varchar(36) NOT NULL PRIMARY KEY,
    name text NOT NULL,
    age int NOT NULL,
    email text NOT NULL,
    pass text NOT NULL,
    UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS books(
    bid varchar(36) NOT NULL PRIMARY KEY,
    author text NOT NULL,
    lable text NOT NULL,
    description text NOT NULL,
    genre text NOT NULL,
    writed_at varchar(4) NOT NULL,
    count int NOT NULL,
    UNIQUE (lable)
);
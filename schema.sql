Create table if not exists users
(
    id       Bigserial primary key,
    login    Text UNIQUE not null,
    password text        not null
);

INSERT INTO users (login, password)
VALUES ('admin', 'pass'),
       ('qwe', 'pass');

INSERT INTO users (login, password)
VALUES (?, ?);

Create table if not exists rooms
(
    id       Bigserial primary key,
    status   bool,
    timestart text NOT NULL ,
    timestop text NOT NULL ,
    filename  text NOT NULL ,
    removed   bool DEFAULT FALSE
);

INSERT INTO mitings (status)
VALUES (?);

SELECT id, timeinhour, timeinminutes, timeouthour, timeoutminutes, filename FROM mitings;
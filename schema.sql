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

CREATE TABLE If Not Exists rooms
(
    id       BIGSERIAL PRIMARY KEY,
    name     text NOT NULL,
    status   BOOLEAN DEFAULT FALSE,
    filename  TEXT NOT NULL,
    removed  BOOLEAN DEFAULT FALSE
);

INSERT INTO mitings (status)
VALUES (?);

SELECT id, timeinhour, timeinminutes, timeouthour, timeoutminutes, filename FROM mitings;

Update rooms  Set status = true where  timestart < ? < timestop;
update products set removed = true where id = $1`, id

SELECT id, status, timestart, timestop, filename FROM rooms where removed = false and status = false;

CREATE TABLE If Not Exists rooms_history
(
    id       BIGSERIAL PRIMARY KEY,
    room_id  INTEGER REFERENCES rooms(id),
    user_login TEXT NOT NULL,
    name_meeting TEXT NOT NULL,
    start_time BIGINT NOT NULL ,
    end_time BIGINT NOT NULL ,
    result TEXT DEFAULT ' '
);

INSERT INTO rooms_history(room_id, user_id, name_meeting, start_time, end_time)
VALUES (?, ?, ?, ?, ?);
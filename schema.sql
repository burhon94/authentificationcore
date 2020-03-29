Create table if not exists users
(
    id          BIGSERIAL primary key,
    login       TEXT UNIQUE not null,
    password    TEXT        not null,
    namesurname TEXT        not null,
    avatar      TEXT,
    roles       TEXT[]
);

INSERT INTO users (login, password, namesurname, avatar, roles)
VALUES ('vasya', '$2y$12$XgMXfNNdB/Zb8I0Du36lwuDHPH.LxK5MVlpy/fDiFoM7NnS.1bPOC', 'Vasya', 'https://i.pravatar.cc/200',
        '{user, admin}'),
       ('petya', '$2y$12$1BP55i1Y9mpveKj4MTiwKOqcp391Eam2hXkgW8cxrSlE2sw6PAJFK', 'Petya', 'https://i.pravatar.cc/200',
        '{user}');

INSERT INTO users (login, password, namesurname)
VALUES (?, ?, ?);

SELECT login, password from users WHERE login = 'vasya';

SELECT login, password from users WHERE id = ?;

SELECT login, namesurname, avatar, roles FROM users WHERE id = ?;

SELECT login, password, namesurname, avatar FROM users WHERE id = ?;

UPDATE users SET namesurname = ? WHERE id = ?;

SELECT password from users WHERE id = ?;

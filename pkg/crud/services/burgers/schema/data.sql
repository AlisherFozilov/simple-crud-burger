INSERT INTO burgers (name, price)
VALUES ('Big Mac', 16000),
       ('Chicken Mac', 12000);

UPDATE burgers SET removed = true WHERE id = ?;

INSERT INTO burgers (name, price) VALUES (?, ?);
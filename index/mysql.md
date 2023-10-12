``` sql
CREATE TABLE cards (
   mark text,
   number text
);

CREATE INDEX idx_cards ON cards(mark(16), number(16));


-- echo {1..13} | xargs -n 1 bash -c "mysql --host=localhost --user=user --password=password --database=db -e \"INSERT INTO cards VALUES ('💚', '\${0}');\""
-- echo {1..13} | xargs -n 1 bash -c "mysql --host=localhost --user=user --password=password --database=db -e \"INSERT INTO cards VALUES ('🍀', '\${0}');\""
-- echo {1..13} | xargs -n 1 bash -c "mysql --host=localhost --user=user --password=password --database=db -e \"INSERT INTO cards VALUES ('🔷', '\${0}');\""
-- echo {1..13} | xargs -n 1 bash -c "mysql --host=localhost --user=user --password=password --database=db -e \"INSERT INTO cards VALUES ('♠️', '\${0}');\""



echo {1..13} | xargs -n 1 bash -c "mysql --host=localhost --user=user --password=password --database=db -e \"INSERT INTO cards VALUES ('A', '\${0}');\""
echo {1..13} | xargs -n 1 bash -c "mysql --host=localhost --user=user --password=password --database=db -e \"INSERT INTO cards VALUES ('B', '\${0}');\""
echo {1..13} | xargs -n 1 bash -c "mysql --host=localhost --user=user --password=password --database=db -e \"INSERT INTO cards VALUES ('C', '\${0}');\""
echo {1..13} | xargs -n 1 bash -c "mysql --host=localhost --user=user --password=password --database=db -e \"INSERT INTO cards VALUES ('D', '\${0}');\""



-- ANALYZE cards;
ANALYZE TABLE cards;


EXPLAIN SELECT * FROM cards WHERE mark = 'A';


ALTER DATABASE db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE cards CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

```

DROP TABLE IF EXISTS tbl_users_to_items cascade;
DROP TABLE IF EXISTS tbl_item cascade;
DROP TABLE IF EXISTS tbl_users cascade;

CREATE TABLE tbl_users
(
    ID       SERIAL PRIMARY KEY,
    name     VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);


CREATE TABLE tbl_item
(
    ID          SERIAL PRIMARY KEY,
    image       VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    item_name   VARCHAR(255) NOT NULL,
    cost        FLOAT        NOT NULL,
    item_type   INT          NOT NULL
);

CREATE TABLE tbl_users_to_items
(
    ID         SERIAL PRIMARY KEY,
    user_id    INT REFERENCES tbl_users (ID) ON DELETE SET NULL,
    item_id    INT REFERENCES tbl_item (ID) ON DELETE SET NULL,
    item_count INT NOT NULL
);



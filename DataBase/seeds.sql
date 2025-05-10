INSERT INTO tbl_users (name, password)
VALUES ('Alice', 'password123'),
       ('Bob', 'securepass'),
       ('Charlie', 'mypassword');
INSERT INTO tbl_item (image, description, item_name, cost)
VALUES ('image1.jpg', 'A beautiful item 1', 'Item 1', 19.99),
       ('image2.jpg', 'A useful item 2', 'Item 2', 29.99),
       ('image3.jpg', 'An amazing item 3', 'Item 3', 39.99);
INSERT INTO tbl_users_to_items (user_id, item_id)
VALUES (1, 1), -- Alice owns Item 1
       (1, 2), -- Alice owns Item 2
       (2, 3), -- Bob owns Item 3
       (3, 1); -- Charlie owns Item 1

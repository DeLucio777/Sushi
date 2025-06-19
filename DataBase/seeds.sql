INSERT INTO tbl_users (name, password)
VALUES ('Alice', 'password123'),
       ('Bob', 'securepass'),
       ('Charlie', 'mypassword');
INSERT INTO tbl_item (image, description, item_name, cost, item_type)
VALUES ('image1.jpg', 'A beautiful item 1', 'Item 1', 19.99,1),
       ('image2.jpg', 'A useful item 2', 'Item 2', 29.99,2),
       ('image3.jpg', 'An amazing item 3', 'Item 3', 39.99,3);
INSERT INTO tbl_users_to_items (user_id, item_id)
VALUES (1, 1), 
       (1, 2), 
       (2, 3), 
       (3, 1);


-- INSERT INTO tbl_item (image, description, item_name, cost, item_type)
-- VALUES ('photo0jpg.jpg', 'Суши из форели, водорослей, сливочного соуса, огурцов и соуса терияки', 'Суши с терияки', 19.99,1),
--        ('1eefcf4a097f5e3f325acf638484be---jpg_800x500_whitepadding15_fb108_convert.jpg', 'Суши из красной рыбы, водорослей, сливочного соуса и огурцов', 'Суши Филадельфия', 14.99,1),
--        ('images.jpg', 'Суши из красной рыбы, риса и водорослей', 'Суши Нигири', 24.99,1),
--        ('kanada-sushi-roll.jpg', 'Суши из красной рыбы, риса, водорослей, огурца и сливочного соуса', 'Канада суши', 9.99,1),
--        ('roll1.jpg', 'Ролл с крабом, авокадо и огурцом', 'Ролл Калифорния', 29.99,2),
--        ('drink1.jpg', 'Освежающий зеленый чай, 0,5 л', 'Зеленый чай', 1.99,3),
--        ('464731543027662.jpg', 'Освежающий газированный напиток. 0,5 л', 'Кока-кола', 2.99,3);

INSERT INTO tbl_users (name, password) VALUES
                                           ('user1', 'password1'),
                                           ('user2', 'password2'),
                                           ('user3', 'password3');
INSERT INTO tbl_tresh (user_id, items_count) VALUES
                                                 (1, 5),
                                                 (2, 10),
                                                 (3, 3);
INSERT INTO tbl_item (image, description, item_name, cost, FK_tresh_ID) VALUES
                                                                            ('image1.jpg', 'Description for item 1', 'Item 1', 19.99, 1),
                                                                            ('image2.jpg', 'Description for item 2', 'Item 2', 29.99, 1),
                                                                            ('image3.jpg', 'Description for item 3', 'Item 3', 39.99, 2),
                                                                            ('image4.jpg', 'Description for item 4', 'Item 4', 49.99, 3);

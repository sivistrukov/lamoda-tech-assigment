INSERT INTO warehouses (name, is_available)
VALUES ('Main', true),
       ('Second', true);

INSERT INTO products (code, name, size)
VALUES ('A001', 'Product A', '?'),
       ('B001', 'Product B', '?'),
       ('C001', 'Product C', '?');

INSERT INTO warehouse_products (warehouse_id, product_code, quantity, status)
VALUES (1, 'A001', 100, 'available'),
       (1, 'B001', 50, 'available'),
       (1, 'C001', 100, 'reservation'),
       (1, 'A001', 100, 'reservation'),
       (2, 'C001', 50, 'available'),
       (2, 'A001', 50, 'available'),
       (2, 'A001', 50, 'reservation');

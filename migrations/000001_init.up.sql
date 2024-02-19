CREATE TABLE warehouses
(
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    is_available BOOLEAN      NOT NULL
);

CREATE TABLE products
(
    code VARCHAR(30) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    size VARCHAR(255) DEFAULT ''
);

CREATE TABLE warehouse_products
(
    id           serial PRIMARY KEY,
    warehouse_id INTEGER REFERENCES warehouses (id)     NOT NULL,
    product_code VARCHAR(30) REFERENCES products (Code) NOT NULL,
    quantity     INTEGER CHECK (quantity >= 0)          NOT NULL,
    status       VARCHAR(30)                            NOT NULL,
    CONSTRAINT unique_product_in_warehouse UNIQUE (warehouse_id, product_code, status)
);

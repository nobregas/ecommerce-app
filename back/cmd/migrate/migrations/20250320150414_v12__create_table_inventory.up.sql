CREATE TABLE inventory (
    `product_id` INT UNSIGNED,
    `stock_quantity` INT UNSIGNED NOT NULL DEFAULT 0,
    `version` INT UNSIGNED NOT NULL DEFAULT 1,

    PRIMARY KEY (`product_id`),
    FOREIGN KEY (`product_id`) REFERENCES products(`id`) ON DELETE CASCADE
);
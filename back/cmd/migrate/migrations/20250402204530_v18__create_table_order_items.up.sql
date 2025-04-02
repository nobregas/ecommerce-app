CREATE TABLE IF NOT EXISTS order_items (
    `orderId` INT UNSIGNED NOT NULL,
    `productId` INT UNSIGNED NOT NULL,
    `quantity` INT UNSIGNED NOT NULL,
    `price` DECIMAL(10,2) UNSIGNED NOT NULL,
    PRIMARY KEY (`orderId`, `productId`),
    FOREIGN KEY (`orderId`) REFERENCES order_history(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`productId`) REFERENCES products(`id`) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS product_categories (
  `productId` INT UNSIGNED NOT NULL,
  `categoryId` INT UNSIGNED NOT NULL,

  PRIMARY KEY (`productId`, `categoryId`),
  FOREIGN KEY (`productId`) REFERENCES products(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`categoryId`) REFERENCES categories(`id`) ON DELETE CASCADE
);
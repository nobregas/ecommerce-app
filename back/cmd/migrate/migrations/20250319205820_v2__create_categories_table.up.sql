CREATE TABLE IF NOT EXISTS categories (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `imageUrl` VARCHAR(512),
  `parentCategoryId` INT UNSIGNED,

  PRIMARY KEY (`id`),
  FOREIGN KEY (`parentCategoryId`) REFERENCES categories(`id`) ON DELETE SET NULL
);
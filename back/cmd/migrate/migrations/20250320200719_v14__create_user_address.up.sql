CREATE TABLE IF NOT EXISTS user_address (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `userId` INT UNSIGNED NOT NULL,
  `street` VARCHAR(255) NOT NULL,
  `city` VARCHAR(100) NOT NULL,
  `state` VARCHAR(50) NOT NULL,
  `postalCode` VARCHAR(20) NOT NULL,
  `country` VARCHAR(50) NOT NULL DEFAULT 'Brasil',
  `isDefault` BOOLEAN NOT NULL DEFAULT false,
  `createdAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  PRIMARY KEY (`id`),
  FOREIGN KEY (`userId`) REFERENCES users(`id`) ON DELETE CASCADE,
  INDEX `idx_address_user` (`userId`),
  INDEX `idx_address_city` (`city`),
  UNIQUE KEY `unique_default_address` (`userId`, `isDefault`)
)
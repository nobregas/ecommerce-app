ALTER TABLE cart_items
    ADD COLUMN productTitle VARCHAR(255) NOT NULL AFTER productImage;
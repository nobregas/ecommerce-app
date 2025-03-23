package product

import (
	"database/sql"
	"fmt"
	types2 "github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"strings"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]*types2.Product, error) {
	// find products
	rows, err := s.db.Query(
		`SELECT p.*, i.stock_quantity, i.version
		FROM products p
		INNER JOIN inventory i ON p.id = i.product_id
		`)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	// get products and their inventory
	productIds := make([]int, 0)
	products := make([]*types2.Product, 0)
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		productIds = append(productIds, p.ID)
		products = append(products, p)
	}

	// find images
	if len(productIds) > 0 {
		images, err := s.GetImagesForProducts(productIds)
		if err != nil {
			return nil, fmt.Errorf("failed to get images: %w", err)
		}

		for _, p := range products {
			p.Images = images[p.ID]
		}
	}

	// find categories
	for _, p := range products {
		cat, err := s.GetProductCategories(p.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get categories: %w", err)
		}
		p.Categories = cat
	}

	return products, nil
}

func (s *Store) GetProductsByCategory(categoryID int) ([]*types2.Product, error) {
	rows, err := s.db.Query(`
        SELECT 
            p.id, 
            p.title, 
            p.description, 
            p.basePrice, 
            p.createdAt, 
            p.updatedAt,
            i.stock_quantity, 
            i.version
        FROM products p
        INNER JOIN inventory i ON p.id = i.product_id
        INNER JOIN product_categories pc ON p.id = pc.productId
        WHERE pc.categoryId = ?
    `, categoryID)

	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	var products []*types2.Product
	var productIDs []int

	for rows.Next() {
		p := new(types2.Product)

		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&p.BasePrice,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.Inventory.StockQuantity,
			&p.Inventory.Version,
		)
		if err != nil {
			return nil, err
		}
		p.Inventory.ProductID = p.ID

		productIDs = append(productIDs, p.ID)
		products = append(products, p)
	}

	if len(productIDs) > 0 {
		imagesMap, err := s.GetImagesForProducts(productIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get images: %w", err)
		}

		for _, p := range products {
			p.Images = imagesMap[p.ID]

			categories, err := s.GetProductCategories(p.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get categories for product %d: %w", p.ID, err)
			}
			p.Categories = categories
		}
	}
	return products, nil
}

func (s *Store) GetImagesForProducts(productIDs []int) (map[int][]types2.ProductImage, error) {
	query := `
        SELECT id, productId, imageUrl, sortOrder 
        FROM product_images 
        WHERE productId IN (?` + strings.Repeat(",?", len(productIDs)-1) + `)`

	args := make([]interface{}, len(productIDs))
	for i, id := range productIDs {
		args[i] = id
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query images: %w", err)
	}
	defer rows.Close()

	imagesMap := make(map[int][]types2.ProductImage)

	for rows.Next() {
		var img types2.ProductImage
		err := rows.Scan(&img.ID, &img.ProductID, &img.ImageUrl, &img.SortOrder)
		if err != nil {
			return nil, fmt.Errorf("failed to scan image: %w", err)
		}
		imagesMap[img.ProductID] = append(imagesMap[img.ProductID], img)
	}

	return imagesMap, nil
}

func (s *Store) GetProductByID(productID int) (*types2.Product, error) {
	// find products
	row := s.db.QueryRow(`SELECT * FROM products WHERE id = ?`, productID)
	product, err := scanRowIntoProduct(row)
	if err != nil {
		return nil, err
	}

	// find inventory
	inv, err := s.GetInventory(productID)
	if err != nil {
		return nil, err
	}
	product.Inventory = *inv

	// find images
	rows, err := s.db.Query(`SELECT * FROM product_images WHERE productId = ?`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var img types2.ProductImage
		err := rows.Scan(&img.ID, &img.ProductID, &img.ImageUrl, &img.SortOrder)
		if err != nil {
			return nil, err
		}
		product.Images = append(product.Images, img)
	}

	categories, err := s.GetProductCategories(productID)
	if err != nil {
		return nil, err
	}
	product.Categories = categories

	return product, nil
}

func (s *Store) GetProductCategories(productID int) ([]types2.Category, error) {
	rows, err := s.db.Query(`
        SELECT c.id, c.name, c.imageUrl, c.parentCategoryId 
        FROM product_categories pc
        JOIN categories c ON pc.categoryId = c.id
        WHERE pc.productId = ?
    `, productID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []types2.Category
	for rows.Next() {
		var c types2.Category
		err := rows.Scan(&c.ID, &c.Name, &c.ImageUrl, &c.ParentCategoryId)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (s *Store) CreateProduct(product types2.CreateProductPayload) error {
	// create transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// create product
	res, err := tx.Exec(
		`INSERT INTO products (title, description, basePrice) VALUES (?, ?, ?)`,
		product.Title, product.Description, product.BasePrice)
	if err != nil {
		return err
	}

	// get product id
	productID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// create inventory
	_, err = tx.Exec(
		`INSERT INTO inventory (product_id, stock_quantity) VALUES (?, ?)`,
		productID, product.StockQuantity,
	)
	if err != nil {
		return err
	}

	// add categories
	for _, categoryID := range product.CategoryIDs {
		_, err := tx.Exec(
			`INSERT INTO product_categories (productId, categoryId) VALUES (?, ?)`,
			productID, categoryID,
		)
		if err != nil {
			return fmt.Errorf("failed to add category: %w", err)
		}
	}

	return tx.Commit()
}

func (s *Store) CreateProductWithImages(payload types2.CreateProductWithImagesPayload) (*types2.Product, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// create product
	res, err := tx.Exec(
		`INSERT INTO products (title, description, basePrice) VALUES (?, ?, ?)`,
		payload.Title, payload.Description, payload.BasePrice,
	)
	if err != nil {
		return nil, err
	}

	// get product id
	productID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// create inventory
	_, err = tx.Exec(
		`INSERT INTO inventory (product_id, stock_quantity) VALUES (?, ?)`,
		productID, payload.StockQuantity,
	)
	if err != nil {
		return nil, err
	}

	// create product images
	for _, img := range payload.Images {
		_, err := tx.Exec(
			`INSERT INTO product_images (productId, imageUrl, sortOrder) VALUES (?, ?, ?)`,
			productID, img.ImageUrl, img.SortOrder,
		)
		if err != nil {
			return nil, err
		}
	}

	// add categories
	for _, categoryID := range payload.CategoryIDs {
		_, err := tx.Exec(
			`INSERT INTO product_categories (productId, categoryId) VALUES (?, ?)`,
			productID, categoryID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to add category: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// find created product
	return s.GetProductByID(int(productID))
}

func (s *Store) UpdateStock(productID int, quantityChange int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentStock, version int
	err = tx.QueryRow(
		`SELECT stock_quantity, version FROM inventory WHERE product_id = ? FOR UPDATE`,
		productID,
	).Scan(&currentStock, &version)
	if err != nil {
		return err
	}

	newStock := currentStock + quantityChange
	if newStock < 0 {
		return fmt.Errorf("insufficient stock")
	}

	_, err = tx.Exec(
		`UPDATE inventory SET stock_quantity = ?, version = version + 1 WHERE product_id = ? AND version = ?`,
		newStock, productID, version,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) GetInventory(productID int) (*types2.Inventory, error) {
	const query = `
        SELECT product_id, stock_quantity, version 
        FROM inventory 
        WHERE product_id = ?
    `

	var inventory types2.Inventory

	err := s.db.QueryRow(query, productID).Scan(
		&inventory.ProductID,
		&inventory.StockQuantity,
		&inventory.Version,
	)
	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("inventory not found for product ID %d", productID)
	case err != nil:
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	return &inventory, nil
}

func (s *Store) UpdateProduct(productID int, payload types2.UpdateProductPayload) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// update product data
	if err := updateProductDetails(tx, productID, payload); err != nil {
		return err
	}

	// update images
	if payload.Images != nil {
		if err := updateProductImages(tx, productID, payload.Images); err != nil {
			return err
		}
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Store) DeleteProduct(productID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// deletes product
	result, err := tx.Exec("DELETE FROM products WHERE id = ?", productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	// verify if product exists
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // product not found
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func updateProductImages(tx *sql.Tx, productID int, images []types2.ImageUpdatePayload) error {
	// get current images
	currentImages, err := getCurrentImages(tx, productID)
	if err != nil {
		return err
	}

	toDelete := make([]int, 0)
	toUpdate := make([]types2.ImageUpdatePayload, 0)
	toCreate := make([]types2.ImageUpdatePayload, 0)

	// aux map for received images (consider imageUrl + sortOrder as key)
	receivedImageMap := make(map[string]types2.ImageUpdatePayload)
	for _, img := range images {
		key := fmt.Sprintf("%s|%d", img.ImageUrl, img.SortOrder)
		receivedImageMap[key] = img
	}

	// process current images
	for currentID, currentImg := range currentImages {
		key := fmt.Sprintf("%s|%d", currentImg.ImageUrl, currentImg.SortOrder)
		if receivedImg, exists := receivedImageMap[key]; exists {
			// update image id with the correspondent one
			receivedImg.ID = &currentID
			toUpdate = append(toUpdate, receivedImg)
			delete(receivedImageMap, key) // remove from map
		} else {
			// mark delete if wasnt find
			toDelete = append(toDelete, currentID)
		}
	}

	// process new images
	for _, img := range receivedImageMap {
		if img.ID == nil {
			toCreate = append(toCreate, img)
		}
	}

	// exec operations
	if err := deleteImages(tx, toDelete); err != nil {
		return err
	}
	if err := updateImages(tx, toUpdate); err != nil {
		return err
	}
	if err := createImages(tx, productID, toCreate); err != nil {
		return err
	}

	return nil
}

func getCurrentImages(tx *sql.Tx, productID int) (map[int]types2.ProductImage, error) {
	images := make(map[int]types2.ProductImage)

	rows, err := tx.Query("SELECT id, imageUrl, sortOrder FROM product_images WHERE productId = ?", productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var img types2.ProductImage
		err := rows.Scan(&img.ID, &img.ImageUrl, &img.SortOrder)
		if err != nil {
			return nil, err
		}
		images[img.ID] = img
	}

	return images, nil
}

func deleteImages(tx *sql.Tx, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	query := "DELETE FROM product_images WHERE id IN (?" + strings.Repeat(",?", len(ids)-1) + ")"
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	_, err := tx.Exec(query, args...)
	return err
}

func updateImages(tx *sql.Tx, images []types2.ImageUpdatePayload) error {
	for _, img := range images {
		_, err := tx.Exec(
			"UPDATE product_images SET imageUrl = ?, sortOrder = ? WHERE id = ?",
			img.ImageUrl, img.SortOrder, *img.ID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func createImages(tx *sql.Tx, productID int, images []types2.ImageUpdatePayload) error {
	for _, img := range images {
		_, err := tx.Exec(
			"INSERT INTO product_images (productId, imageUrl, sortOrder) VALUES (?, ?, ?)",
			productID, img.ImageUrl, img.SortOrder,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateProductDetails(tx *sql.Tx, productID int, payload types2.UpdateProductPayload) error {
	query := "UPDATE products SET"
	args := make([]interface{}, 0)
	updates := make([]string, 0)

	if payload.Title != nil {
		updates = append(updates, " title = ?")
		args = append(args, *payload.Title)
	}
	if payload.Description != nil {
		updates = append(updates, " description = ?")
		args = append(args, *payload.Description)
	}
	if payload.BasePrice != nil {
		updates = append(updates, " basePrice = ?")
		args = append(args, *payload.BasePrice)
	}

	if len(updates) == 0 {
		return nil
	}

	query += strings.Join(updates, ",") + " WHERE id = ?"
	args = append(args, productID)

	_, err := tx.Exec(query, args...)
	return err
}

func scanRowsIntoProduct(rows *sql.Rows) (*types2.Product, error) {
	product := new(types2.Product)

	err := rows.Scan(
		&product.ID,
		&product.Title,
		&product.Description,
		&product.BasePrice,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.Inventory.StockQuantity,
		&product.Inventory.Version,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan product: %w", err)
	}
	product.Inventory.ProductID = product.ID

	return product, nil
}

func scanRowIntoProduct(row *sql.Row) (*types2.Product, error) {
	product := new(types2.Product)
	err := row.Scan(
		&product.ID,
		&product.Title,
		&product.Description,
		&product.BasePrice,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	return product, err
}

package product

import (
	"database/sql"
	"fmt"
	types "github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"strings"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]*types.Product, error) {
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
	products := make([]*types.Product, 0)
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

func (s *Store) GetProductsByCategory(categoryID int) ([]*types.Product, error) {
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

	var products []*types.Product
	var productIDs []int

	for rows.Next() {
		p := new(types.Product)

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

func (s *Store) GetImagesForProducts(productIDs []int) (map[int][]types.ProductImage, error) {
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

	imagesMap := make(map[int][]types.ProductImage)

	for rows.Next() {
		var img types.ProductImage
		err := rows.Scan(&img.ID, &img.ProductID, &img.ImageUrl, &img.SortOrder)
		if err != nil {
			return nil, fmt.Errorf("failed to scan image: %w", err)
		}
		imagesMap[img.ProductID] = append(imagesMap[img.ProductID], img)
	}

	return imagesMap, nil
}

func (s *Store) GetProductByID(productID int) (*types.Product, error) {
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
		var img types.ProductImage
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

func (s *Store) GetProductCategories(productID int) ([]types.Category, error) {
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

	var categories []types.Category
	for rows.Next() {
		var c types.Category
		err := rows.Scan(&c.ID, &c.Name, &c.ImageUrl, &c.ParentCategoryId)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (s *Store) CreateProduct(product types.CreateProductPayload) error {
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

func (s *Store) CreateProductWithImages(payload types.CreateProductWithImagesPayload) (*types.Product, error) {
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

func (s *Store) GetInventory(productID int) (*types.Inventory, error) {
	const query = `
        SELECT product_id, stock_quantity, version 
        FROM inventory 
        WHERE product_id = ?
    `

	var inventory types.Inventory

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

func (s *Store) UpdateProduct(productID int, payload types.UpdateProductPayload) error {
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

func (s *Store) GetProductDetails(userID int, productID int) (*types.ProductDetails, error) {
	query := `
        SELECT 
            p.id,
            p.title,
            p.description,
            p.basePrice,
            COALESCE(MAX(CASE WHEN NOW() BETWEEN d.startDate AND d.endDate THEN d.discountPercent END), 0) AS discount,
            COALESCE(AVG(r.rating), 0) AS avg_rating,
            EXISTS(SELECT 1 FROM user_favorites uf WHERE uf.userId = ? AND uf.productId = p.id) AS is_favorite
        FROM products p
        LEFT JOIN product_discounts d ON p.id = d.productId
        LEFT JOIN product_ratings r ON p.id = r.productId
        WHERE p.id = ?
        GROUP BY p.id
    `
	var detail types.ProductDetails
	var discount float64

	err := s.db.QueryRow(query, userID, productID).Scan(
		&detail.ID,
		&detail.Title,
		&detail.Description,
		&detail.BasePrice,
		&discount,
		&detail.AverageRating,
		&detail.IsFavorite,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get product details: %w", err)
	}

	detail.DiscountPercentage = discount
	detail.Price = detail.BasePrice * (1 - discount/100)

	images, err := s.GetImagesForProducts([]int{productID})
	if err != nil {
		return nil, fmt.Errorf("failed to get images: %w", err)
	}
	if imgList, exists := images[productID]; exists {
		detail.Images = imgList
	}

	return &detail, nil
}

func (s *Store) GetSimpleProductDetails(userID int) (*[]*types.SimpleProductObject, error) {
	query := `
        SELECT 
            p.id,
            p.title,
            p.basePrice,
            COALESCE(MAX(CASE WHEN NOW() BETWEEN d.startDate AND d.endDate THEN d.discountPercent END), 0) AS discount,
            COALESCE(AVG(r.rating), 0) AS avg_rating,
            EXISTS(SELECT 1 FROM user_favorites uf WHERE uf.userId = ? AND uf.productId = p.id) AS is_favorite,
            (SELECT imageUrl FROM product_images WHERE productId = p.id ORDER BY sortOrder LIMIT 1) AS main_image
        FROM products p
        LEFT JOIN product_discounts d ON p.id = d.productId
        LEFT JOIN product_ratings r ON p.id = r.productId
        GROUP BY p.id
    `
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query simple products: %w", err)
	}
	defer rows.Close()

	var products []*types.SimpleProductObject
	for rows.Next() {
		var sp types.SimpleProductObject
		var discount float64
		var imageUrl string

		err := rows.Scan(
			&sp.ID,
			&sp.Title,
			&sp.BasePrice,
			&discount,
			&sp.AverageRating,
			&sp.IsFavorite,
			&imageUrl,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan simple product: %w", err)
		}

		// Calcular preço final
		sp.Price = sp.BasePrice * (1 - discount/100)

		// Construir objeto de imagem
		sp.Image = types.ProductImage{
			ImageUrl:  imageUrl,
			SortOrder: 0,
		}

		products = append(products, &sp)
	}

	return &products, nil
}

func updateProductImages(tx *sql.Tx, productID int, images []types.ImageUpdatePayload) error {
	// get current images
	currentImages, err := getCurrentImages(tx, productID)
	if err != nil {
		return err
	}

	toDelete := make([]int, 0)
	toUpdate := make([]types.ImageUpdatePayload, 0)
	toCreate := make([]types.ImageUpdatePayload, 0)

	// aux map for received images (consider imageUrl + sortOrder as key)
	receivedImageMap := make(map[string]types.ImageUpdatePayload)
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

func getCurrentImages(tx *sql.Tx, productID int) (map[int]types.ProductImage, error) {
	images := make(map[int]types.ProductImage)

	rows, err := tx.Query("SELECT id, imageUrl, sortOrder FROM product_images WHERE productId = ?", productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var img types.ProductImage
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

func updateImages(tx *sql.Tx, images []types.ImageUpdatePayload) error {
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

func createImages(tx *sql.Tx, productID int, images []types.ImageUpdatePayload) error {
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

func updateProductDetails(tx *sql.Tx, productID int, payload types.UpdateProductPayload) error {
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

func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)

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

func scanRowIntoProduct(row *sql.Row) (*types.Product, error) {
	product := new(types.Product)
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

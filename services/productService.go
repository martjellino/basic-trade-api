package services

import (
	"basic-trade-api/models/product"
	"database/sql"
	"errors"
	"time"
)

func CreateProductService(db *sql.DB, productRequest product.ProductRequest, adminId int) (*product.ProductResponse, error) {
	var productResponse product.ProductResponse
	query := `INSERT INTO products (name, image_url, admin_id) VALUES ($1, $2, $3) RETURNING id, uuid, name, image_url, admin_id, created_at, updated_at`
	err := db.QueryRow(query, productRequest.Name, productRequest.ImageURL, adminId).Scan(&productResponse.ID, &productResponse.UUID, &productResponse.Name, &productResponse.ImageURL, &productResponse.AdminID, &productResponse.CreatedAt, &productResponse.UpdatedAt)
	if err != nil {
		return nil, err
	}
	productResponse.ImageFileHeader = productRequest.ImageFile

	// Assign the uploaded file URL to the ProductResponse struct
	productResponse.ImageURL = productRequest.ImageURL

	return &productResponse, nil
}

func GetAllProductService(db *sql.DB, pageSize int, offset int, name string) ([]product.ProductResponse, int, error) {
	var products []product.ProductResponse
	var total int

	// Count total number of products
	query := `SELECT COUNT(*) FROM products`
	err := db.QueryRow(query).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	baseQuery := `SELECT id, uuid, name, image_url, admin_id, created_at, updated_at FROM products`
	// Add WHERE clause if variantName is provided
	if name != "" {
		name = "%" + name + "%"
		baseQuery += ` WHERE name ILIKE $1`
		baseQuery += ` LIMIT $2 OFFSET $3`

		// Execute the query
		rows, err := db.Query(baseQuery, name, pageSize, offset)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()

		// Process the query results
		for rows.Next() {
			var productResponse product.ProductResponse
			err := rows.Scan(&productResponse.ID, &productResponse.UUID, &productResponse.Name, &productResponse.ImageURL, &productResponse.AdminID, &productResponse.CreatedAt, &productResponse.UpdatedAt)
			if err != nil {
				return nil, 0, err
			}
			products = append(products, productResponse)
		}
		if err = rows.Err(); err != nil {
			return nil, 0, err
		}
	} else {
		baseQuery += ` LIMIT $1 OFFSET $2`

		// Execute the query
		rows, err := db.Query(baseQuery, pageSize, offset)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()

		// Process the query results
		for rows.Next() {
			var productResponse product.ProductResponse
			err := rows.Scan(&productResponse.ID, &productResponse.UUID, &productResponse.Name, &productResponse.ImageURL, &productResponse.AdminID, &productResponse.CreatedAt, &productResponse.UpdatedAt)
			if err != nil {
				return nil, 0, err
			}
			products = append(products, productResponse)
		}
		if err = rows.Err(); err != nil {
			return nil, 0, err
		}
	}

	return products, total, nil
}

func GetProductByIDService(db *sql.DB, productUUID string) (*product.ProductResponse, error) {
	var product product.ProductResponse

	query := `SELECT id, uuid, name, image_url, admin_id, created_at, updated_at FROM products WHERE UUID = $1`
	err := db.QueryRow(query, productUUID).Scan(&product.ID, &product.UUID, &product.Name, &product.ImageURL, &product.AdminID, &product.CreatedAt, &product.UpdatedAt)
	if err == sql.ErrNoRows {
		// If no product is found with the given UUID, return a custom error
		return nil, errors.New("product not found")
	} else if err != nil {
		return nil, err
	}

	return &product, nil
}

func UpdateProductService(db *sql.DB, productRequest product.ProductRequest, productUUID string, adminId int) (*product.ProductResponse, error) {
	var product product.ProductResponse
	query := `SELECT id, uuid, name, image_url, admin_id, created_at, updated_at FROM products WHERE UUID = $1 AND admin_id = $2`
	err := db.QueryRow(query, productUUID, adminId).Scan(&product.ID, &product.UUID, &product.Name, &product.ImageURL, &product.AdminID, &product.CreatedAt, &product.UpdatedAt)
	if err == sql.ErrNoRows {
		// If no product is found with the given UUID, return a custom error
		return nil, errors.New("product not found")
	} else if err != nil {
		return nil, err
	}

	product.Name = productRequest.Name
	// Update product.ImageURL with the new uploaded file URL
	if productRequest.ImageURL != "" {
		product.ImageURL = productRequest.ImageURL
	}

	product.ImageFileHeader = productRequest.ImageFile
	product.UpdatedAt = time.Now()

	query = `UPDATE products SET name = $1, image_url = $2, updated_at = $3 WHERE UUID = $4 AND admin_id = $5`
	_, err = db.Exec(query, product.Name, product.ImageURL, product.UpdatedAt, productUUID, adminId)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func DeleteProductService(db *sql.DB, productUUID string, adminId int) (*product.ProductResponse, error) {
	var product product.ProductResponse
	query := `SELECT id, uuid, name, image_url, admin_id, created_at, updated_at FROM products WHERE uuid = $1 AND admin_id = $2`
	rows, err := db.Query(query, productUUID, adminId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		// No product found with the given UUID and adminId
		return nil, errors.New("product not found")
	}
	err = rows.Scan(&product.ID, &product.UUID, &product.Name, &product.ImageURL, &product.AdminID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}

	query = `DELETE FROM products WHERE uuid = $1 AND admin_id = $2`
	_, err = db.Exec(query, productUUID, adminId)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

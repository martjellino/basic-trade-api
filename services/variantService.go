package services

import (
	"basic-trade-api/models/variant"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func CreateVariantService(db *sql.DB, variantReq variant.VariantRequest, adminId int) (*variant.VariantResponse, error) {
	var variantResponse variant.VariantResponse
	query := `INSERT INTO variants (variant_name, quantity, product_id) VALUES ($1, $2, $3) RETURNING id, uuid, variant_name, quantity, product_id, created_at, updated_at`
	err := db.QueryRow(query, variantReq.VariantName, variantReq.Quantity, variantReq.ProductID).Scan(&variantResponse.ID, &variantResponse.UUID, &variantResponse.VariantName, &variantResponse.Quantity, &variantResponse.ProductID, &variantResponse.CreatedAt, &variantResponse.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &variantResponse, nil
}

func GetAllVariantService(db *sql.DB, pageSize int, offset int, variantName string) ([]variant.VariantResponse, int, error) {
	var variants []variant.VariantResponse
	var total int

	// Construct the base query
	baseQuery := `SELECT COUNT(*) FROM variants`
	err := db.QueryRow(baseQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	baseQuery = `SELECT id, uuid, variant_name, quantity, product_id, created_at, updated_at FROM variants`

	// Add WHERE clause if variantName is provided
	if variantName != "" {
		variantName = "%" + variantName + "%"
		baseQuery += ` WHERE variant_name ILIKE $1`
		baseQuery += ` LIMIT $2 OFFSET $3`

		// Execute the query
		rows, err := db.Query(baseQuery, variantName, pageSize, offset)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()

		// Process the query results
		for rows.Next() {
			var variantResponse variant.VariantResponse
			err := rows.Scan(&variantResponse.ID, &variantResponse.UUID, &variantResponse.VariantName, &variantResponse.Quantity, &variantResponse.ProductID, &variantResponse.CreatedAt, &variantResponse.UpdatedAt)
			if err != nil {
				return nil, 0, err
			}
			variants = append(variants, variantResponse)
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
			var variantResponse variant.VariantResponse
			err := rows.Scan(&variantResponse.ID, &variantResponse.UUID, &variantResponse.VariantName, &variantResponse.Quantity, &variantResponse.ProductID, &variantResponse.CreatedAt, &variantResponse.UpdatedAt)
			if err != nil {
				return nil, 0, err
			}
			variants = append(variants, variantResponse)
		}
		if err = rows.Err(); err != nil {
			return nil, 0, err
		}
	}

	return variants, total, nil
}


func GetVariantByIDService(db *sql.DB, variantUUID string) (*variant.VariantResponse, error) {
	var variantResponse variant.VariantResponse

	query := `SELECT id, uuid, variant_name, quantity, product_id, created_at, updated_at FROM variants WHERE uuid = $1`
	err := db.QueryRow(query, variantUUID).Scan(&variantResponse.ID, &variantResponse.UUID, &variantResponse.VariantName, &variantResponse.Quantity, &variantResponse.ProductID, &variantResponse.CreatedAt, &variantResponse.UpdatedAt)
	if err == sql.ErrNoRows {
		// If no product is found with the given UUID, return a custom error
		return nil, errors.New("variant not found")
	} else if err != nil {
		return nil, err
	}
	return &variantResponse, nil
}

func UpdateVariantService(db *sql.DB, variantRequest variant.VariantRequest, variantUUID string, adminId int) (*variant.VariantResponse, error) {
    var variantResponse variant.VariantResponse
    query := `SELECT id, uuid, variant_name, quantity, product_id, created_at, updated_at FROM variants WHERE uuid = $1`
    err := db.QueryRow(query, variantUUID).Scan(&variantResponse.ID, &variantResponse.UUID, &variantResponse.VariantName, &variantResponse.Quantity, &variantResponse.ProductID, &variantResponse.CreatedAt, &variantResponse.UpdatedAt)
    if err == sql.ErrNoRows {
        // If no product is found with the given UUID, return a custom error
        return nil, errors.New("variant not found")
    } else if err != nil {
        return nil, err
    }

    // Update variant details with new values
    variantResponse.VariantName = variantRequest.VariantName
    variantResponse.Quantity = variantRequest.Quantity
    variantResponse.ProductID = variantRequest.ProductID
    variantResponse.UpdatedAt = time.Now()

    query = `UPDATE variants SET variant_name = $1, quantity = $2, product_id = $3, updated_at = $4 WHERE uuid = $5`
    fmt.Printf("Query: %s\nParams: %v\n", query, []interface{}{variantResponse.VariantName, variantResponse.Quantity, variantResponse.ProductID, variantResponse.UpdatedAt, variantUUID})
    _, err = db.Exec(query, variantResponse.VariantName, variantResponse.Quantity, variantResponse.ProductID, variantResponse.UpdatedAt, variantUUID)
    if err != nil {
        return nil, err
    }

    return &variantResponse, nil
}

func DeleteVariantService(db *sql.DB, variantUUID string, adminId int) error {
    query := `SELECT id FROM variants WHERE uuid = $1`
    var variantID int
    err := db.QueryRow(query, variantUUID).Scan(&variantID)
    if err != nil {
        if err == sql.ErrNoRows {
            // Variant not found
            return errors.New("variant not found")
        }
        // Other error occurred
        return err
    }

    // Delete the variant
    deleteQuery := `DELETE FROM variants WHERE uuid = $1`
    _, err = db.Exec(deleteQuery, variantUUID)
    if err != nil {
        return err
    }

    return nil
}
package services

import (
	"basic-trade-api/helpers"
	"basic-trade-api/models/admin"
	"database/sql"
	"errors"
	"fmt"
)

func AdminRegisterService(db *sql.DB, adminRequest admin.AdminRegisterRequest) (*admin.AdminResponse, error) {
	// Check if the email already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM admins WHERE email = $1", adminRequest.Email).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("email already exists")
	}

	// Validate the admin request data
	err = admin.Validate.Struct(adminRequest)
	if err != nil {
		validationErrors := helpers.GeneralValidator(err)
		return nil, fmt.Errorf(fmt.Sprintf("Validation errors: %v", validationErrors))
	}

	// Hash the password before storing it in the database
	hashedPassword, err := helpers.HashPassword(adminRequest.Password)
	if err != nil {
		return nil, err
	}

	// Prepare the SQL query to insert a new admin
	query := `
		INSERT INTO admins (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, uuid, created_at, updated_at
	`

	// Execute the query and get the new admin's ID, UUID, created_at, and updated_at
	var newAdmin admin.AdminResponse
	err = db.QueryRow(query, adminRequest.Name, adminRequest.Email, hashedPassword).Scan(&newAdmin.ID, &newAdmin.UUID, &newAdmin.CreatedAt, &newAdmin.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Set the remaining fields of the AdminResponse struct
	newAdmin.Name = adminRequest.Name
	newAdmin.Email = adminRequest.Email
	newAdmin.Password = "" // Don't return the password

	// Return the new admin
	return &newAdmin, nil
}

func AdminLoginService(db *sql.DB, adminRequest admin.AdminLoginRequest) (*admin.AdminResponse, error) {
	// Validate the admin request data
	err := admin.Validate.Struct(adminRequest)
	if err != nil {
		validationErrors := helpers.GeneralValidator(err)
		return nil, fmt.Errorf(fmt.Sprintf("Validation errors: %v", validationErrors))
	}

	query := `
		SELECT id, name, email, password 
		FROM admins 
		WHERE email = $1
	`

	// Execute the query and get the new admin's ID, UUID, created_at, and updated_at
	var adminResponse admin.AdminResponse
	err = db.QueryRow(query, adminRequest.Email).Scan(&adminResponse.ID, &adminResponse.Name, &adminResponse.Email, &adminResponse.Password)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	// Compare password
	comparePass := helpers.ComparePassword([]byte(adminResponse.Password), []byte(adminRequest.Password))
	if !comparePass {
		return nil, errors.New("invalid password")
	}

	return &adminResponse, nil
}

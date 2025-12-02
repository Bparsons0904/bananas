package generators

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

// IDMap holds generated IDs for foreign key references
type IDMap struct {
	CategoryIDs  []uuid.UUID
	SupplierIDs  []uuid.UUID
	WarehouseIDs []uuid.UUID
	ProductIDs   []uuid.UUID
	CustomerIDs  []uuid.UUID
}

// GenerateCategories creates hierarchical category data
func GenerateCategories(count int) ([][]interface{}, *IDMap) {
	idMap := &IDMap{
		CategoryIDs: make([]uuid.UUID, 0, count),
	}

	rows := make([][]interface{}, 0, count)
	now := time.Now()

	// Create root categories (20% of total)
	rootCount := count / 5
	if rootCount == 0 {
		rootCount = 1
	}

	// Root categories
	for i := 0; i < rootCount; i++ {
		id := uuid.New()
		idMap.CategoryIDs = append(idMap.CategoryIDs, id)

		rows = append(rows, []interface{}{
			id,                          // id
			gofakeit.ProductCategory(),  // name
			gofakeit.Sentence(10),       // description
			nil,                         // parent_id (NULL for root)
			now,                         // created_at
			now,                         // updated_at
			nil,                         // deleted_at
		})
	}

	// Child categories
	for i := rootCount; i < count; i++ {
		id := uuid.New()
		idMap.CategoryIDs = append(idMap.CategoryIDs, id)

		// Randomly assign to a parent category
		parentID := idMap.CategoryIDs[gofakeit.Number(0, len(idMap.CategoryIDs)-1)]

		rows = append(rows, []interface{}{
			id,                          // id
			gofakeit.ProductCategory(),  // name
			gofakeit.Sentence(10),       // description
			parentID,                    // parent_id
			now,                         // created_at
			now,                         // updated_at
			nil,                         // deleted_at
		})
	}

	return rows, idMap
}

// GenerateSuppliers creates supplier data
func GenerateSuppliers(count int, idMap *IDMap) [][]interface{} {
	rows := make([][]interface{}, 0, count)
	now := time.Now()

	for i := 0; i < count; i++ {
		id := uuid.New()
		idMap.SupplierIDs = append(idMap.SupplierIDs, id)

		contactName := gofakeit.Name()
		email := gofakeit.Email()
		phone := gofakeit.Phone()
		address := gofakeit.Address().Address
		city := gofakeit.City()
		state := gofakeit.StateAbr()
		postalCode := gofakeit.Zip()
		country := gofakeit.Country()

		rows = append(rows, []interface{}{
			id,                    // id
			gofakeit.Company(),    // name
			&contactName,          // contact_name
			&email,                // email
			&phone,                // phone
			&address,              // address
			&city,                 // city
			&state,                // state
			&postalCode,           // postal_code
			&country,              // country
			now,                   // created_at
			now,                   // updated_at
			nil,                   // deleted_at
		})
	}

	return rows
}

// GenerateWarehouses creates warehouse data
func GenerateWarehouses(count int, idMap *IDMap) [][]interface{} {
	rows := make([][]interface{}, 0, count)
	now := time.Now()

	for i := 0; i < count; i++ {
		id := uuid.New()
		idMap.WarehouseIDs = append(idMap.WarehouseIDs, id)

		address := gofakeit.Address().Address
		city := gofakeit.City()
		state := gofakeit.StateAbr()
		postalCode := gofakeit.Zip()
		country := gofakeit.Country()

		rows = append(rows, []interface{}{
			id,                           // id
			fmt.Sprintf("Warehouse %s", city), // name
			fmt.Sprintf("WH-%04d", i+1),  // code
			&address,                     // address
			&city,                        // city
			&state,                       // state
			&postalCode,                  // postal_code
			&country,                     // country
			now,                          // created_at
			now,                          // updated_at
			nil,                          // deleted_at
		})
	}

	return rows
}

// CategoryColumns returns the column names for categories table
func CategoryColumns() []string {
	return []string{"id", "name", "description", "parent_id", "created_at", "updated_at", "deleted_at"}
}

// SupplierColumns returns the column names for suppliers table
func SupplierColumns() []string {
	return []string{"id", "name", "contact_name", "email", "phone", "address", "city", "state", "postal_code", "country",
		"created_at", "updated_at", "deleted_at"}
}

// WarehouseColumns returns the column names for warehouses table
func WarehouseColumns() []string {
	return []string{"id", "name", "code", "address", "city", "state", "postal_code", "country",
		"created_at", "updated_at", "deleted_at"}
}

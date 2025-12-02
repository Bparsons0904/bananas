package generators

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

// GenerateCustomers creates customer data
func GenerateCustomers(count int, idMap *IDMap) [][]interface{} {
	rows := make([][]interface{}, 0, count)
	now := time.Now()

	for i := 0; i < count; i++ {
		id := uuid.New()
		idMap.CustomerIDs = append(idMap.CustomerIDs, id)

		person := gofakeit.Person()
		phone := person.Contact.Phone

		rows = append(rows, []interface{}{
			id,                   // id
			person.FirstName,     // first_name
			person.LastName,      // last_name
			person.Contact.Email, // email
			&phone,               // phone
			now,                  // created_at
			now,                  // updated_at
			nil,                  // deleted_at
		})
	}

	return rows
}

// GenerateCustomerAddresses creates customer addresses
func GenerateCustomerAddresses(idMap *IDMap, addressesPerCustomer int) [][]interface{} {
	rows := make([][]interface{}, 0, len(idMap.CustomerIDs)*addressesPerCustomer)
	now := time.Now()

	addressTypes := []string{"billing", "shipping", "both"}

	for _, customerID := range idMap.CustomerIDs {
		// Generate 1-3 addresses per customer
		numAddresses := gofakeit.Number(1, addressesPerCustomer+1)

		for j := 0; j < numAddresses; j++ {
			addr := gofakeit.Address()
			addressType := addressTypes[gofakeit.Number(0, len(addressTypes)-1)]
			isDefault := j == 0 // First address is default
			addressLine2 := ""
			if gofakeit.Bool() {
				addressLine2 = fmt.Sprintf("Apt %d", gofakeit.Number(1, 999))
			}

			rows = append(rows, []interface{}{
				uuid.New(),       // id
				customerID,       // customer_id
				addressType,      // address_type
				addr.Address,     // address_line1
				&addressLine2,    // address_line2
				addr.City,        // city
				addr.State,       // state
				addr.Zip,         // postal_code
				addr.Country,     // country
				isDefault,        // is_default
				now,              // created_at
				now,              // updated_at
			})
		}
	}

	return rows
}

// Customer-related column functions
func CustomerColumns() []string {
	return []string{"id", "first_name", "last_name", "email", "phone",
		"created_at", "updated_at", "deleted_at"}
}

func CustomerAddressColumns() []string {
	return []string{"id", "customer_id", "address_type", "address_line1", "address_line2", "city", "state",
		"postal_code", "country", "is_default", "created_at", "updated_at"}
}

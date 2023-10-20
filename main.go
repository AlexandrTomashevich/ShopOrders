package main

import (
	"ShopOrders/db"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"sort"
	"strings"
)

type Product struct {
	ShelfName         string
	ProductName       string
	ProductID         int
	OrderID           int
	ProductCount      int
	AdditionalShelves string
}

func main() {

	args := os.Args[1:]
	orderNumbers := strings.Join(args, ",")
	fmt.Println("Страница сборки заказов", orderNumbers)
	fmt.Println("=+=+=+=")

	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT shelf_name, product_name, product_id, order_id, product_count, additional_shelves FROM GetSummaryWithShelves(ARRAY[%s]);`, orderNumbers)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product
		var additionalShelves sql.NullString
		err := rows.Scan(&product.ShelfName, &product.ProductName, &product.ProductID, &product.OrderID, &product.ProductCount, &additionalShelves)
		if err != nil {
			log.Fatal(err)
		}

		product.AdditionalShelves = additionalShelves.String // Заменяем NULL на пустую строку, если не NULL
		products = append(products, product)
	}

	sort.SliceStable(products, func(i, j int) bool {
		if products[i].ShelfName != products[j].ShelfName {
			return products[i].ShelfName < products[j].ShelfName
		}
		return products[i].ProductID < products[j].ProductID
	})

	var prevShelfName string

	for _, product := range products {
		shelfName := product.ShelfName
		productName := product.ProductName
		productID := product.ProductID
		orderID := product.OrderID
		productCount := product.ProductCount
		additionalShelves := product.AdditionalShelves

		if shelfName != prevShelfName {
			fmt.Printf("\n%s\n", shelfName)
			prevShelfName = shelfName
		}

		fmt.Printf("%s (id=%d)\norder %d, %d шт\n", productName, productID, orderID, productCount)

		if productID != 1 && productID != 2 && productID != 4 && productID != 6 {
			fmt.Printf("доп стеллаж: %s", additionalShelves)
		}

		fmt.Println("\n")
	}
}

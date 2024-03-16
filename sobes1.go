package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "lew"
	password = "qwert"
	dbname   = "sobes"
)

func main() {
	db, err := sql.Open("postgres", getConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	orderNumbers := readOrderNumbersFromArgs()
	startTime := time.Now()

	printOrdersByShelves(db, orderNumbers)
	elapsedTime := time.Since(startTime)

	fmt.Printf("Время выполнения: %v\n", elapsedTime)

}

func getConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func readOrderNumbersFromArgs() []int {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		log.Fatal("Please provide order numbers as command-line arguments")
	}

	var orderNumbers []int
	for _, arg := range args {
		orderNumber := 0
		fmt.Sscanf(arg, "%d", &orderNumber)
		orderNumbers = append(orderNumbers, orderNumber)
	}

	return orderNumbers
}

func printOrdersByShelves(db *sql.DB, orderNumbers []int) {
	shelfItems := make(map[string]map[int]string)
	shelfAdds := make(map[string]map[int][]string)

	for _, orderNumber := range orderNumbers {
		query := `
			SELECT Shelves.Name AS Main_Shelf, Products.Name, OrderItems.Quantity, Products.ID AS Product_ID, ProductShelfRelations.Add_Shelf_ID
			FROM OrderItems, Products, ProductShelfRelations, Shelves
			WHERE OrderItems.Product_ID = Products.ID AND OrderItems.Product_ID = ProductShelfRelations.Product_ID AND ProductShelfRelations.Main_Shelf_ID = Shelves.ID AND OrderItems.Order_Number = $1
		`

		rows, err := db.Query(query, orderNumber)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var mainShelf, productName, addShelf sql.NullString
			var quantity, productID int

			err := rows.Scan(&mainShelf, &productName, &quantity, &productID, &addShelf)
			if err != nil {
				log.Fatal(err)
			}

			item := fmt.Sprintf("%s (id=%d)\nзаказ %d, %d шт", productName.String, productID, orderNumber, quantity)

			if _, ok := shelfItems[mainShelf.String]; !ok {
				shelfItems[mainShelf.String] = make(map[int]string)
				shelfAdds[mainShelf.String] = make(map[int][]string)
			}

			if _, ok := shelfItems[mainShelf.String][productID]; !ok {
				shelfItems[mainShelf.String][productID] = item
			}

			if addShelf.Valid {
				shelfAdds[mainShelf.String][productID] = append(shelfAdds[mainShelf.String][productID], addShelf.String)
			}
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
	}

	for shelf, items := range shelfItems {
		fmt.Printf("=== %s\n", shelf)
		for productID, item := range items {
			adds := strings.Join(shelfAdds[shelf][productID], ", ")
			if adds != "" {
				item += fmt.Sprintf("\nдоп стеллаж: %s", adds)
			}
			fmt.Printf("%s\n\n", item)
		}
	}
}

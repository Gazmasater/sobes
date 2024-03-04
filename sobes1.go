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
	shelfItems := make(map[string][]string)

	for _, orderNumber := range orderNumbers {
		query := `
			SELECT Shelves_Main.Name AS Main_Shelf, Products.Name, OrderItems.Quantity, Products.ID AS Product_ID, ARRAY_AGG(DISTINCT Shelves_Add.Name) AS Add_Shelves
			FROM OrderItems
			INNER JOIN Products ON OrderItems.Product_ID = Products.ID
			INNER JOIN ProductShelfRelations ON OrderItems.Product_ID = ProductShelfRelations.Product_ID
			INNER JOIN Shelves Shelves_Main ON ProductShelfRelations.Main_Shelf_ID = Shelves_Main.ID
			LEFT JOIN Shelves Shelves_Add ON ProductShelfRelations.Add_Shelf_ID = Shelves_Add.ID
			WHERE OrderItems.Order_Number = $1
			GROUP BY Main_Shelf, Products.Name, OrderItems.Quantity, Products.ID
		`

		rows, err := db.Query(query, orderNumber)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var mainShelf, productName sql.NullString
			var quantity, productID int
			var addShelvesArray []byte

			err := rows.Scan(&mainShelf, &productName, &quantity, &productID, &addShelvesArray)
			if err != nil {
				log.Fatal(err)
			}

			addShelves := parseAddShelvesArray(addShelvesArray)

			item := fmt.Sprintf("%s (id=%d)\nзаказ %d, %d шт", productName.String, productID, orderNumber, quantity)

			if len(addShelves) > 0 {
				item += fmt.Sprintf("\nдоп стеллаж:%s", strings.Join(addShelves, ", "))
			}

			shelfItems[mainShelf.String] = append(shelfItems[mainShelf.String], item)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
	}

	for shelf, items := range shelfItems {
		fmt.Printf("=== %s\n", shelf)
		for _, item := range items {
			fmt.Printf("%s\n\n", item)
		}
	}
}

func parseAddShelvesArray(addShelvesArray []byte) []string {
	var addShelves []string

	// Check if the array is not empty and not "NULL"
	if len(addShelvesArray) > 2 {
		// Remove curly braces from the array representation
		addShelvesArray = addShelvesArray[1 : len(addShelvesArray)-1]

		// Check if the array is not "NULL"
		if string(addShelvesArray) != "NULL" {
			// Split the array into individual elements
			addShelfElements := strings.Split(string(addShelvesArray), ",")

			// Trim spaces from each element and add to the slice
			for _, element := range addShelfElements {
				addShelves = append(addShelves, strings.TrimSpace(element))
			}
		}
	}

	return addShelves
}

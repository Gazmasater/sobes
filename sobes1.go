package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/lib/pq"
)

func printItemsByShelves(rows *sql.Rows, first bool) error {
	itemChan := make(chan []interface{})
	errChan := make(chan error)
	done := make(chan bool)

	go func() {
		for {
			select {
			case item, ok := <-itemChan:
				if !ok {
					itemChan = nil
					continue
				}
				mainShelf := item[0].(string)
				orderID := item[1].(int)
				productName := item[2].(string)
				totalItems := item[3].(int)
				additionalShelves := item[4].(pq.StringArray)

				// Вывод информации о стеллаже только при первом вызове функции
				if first {
					fmt.Printf("Стеллаж: %s\n", mainShelf)
					first = false
				}

				// Вывод информации о товаре на стеллаже
				fmt.Printf("Заказ %d\n", orderID)
				fmt.Printf("Товар: %s\n", productName)
				fmt.Printf("Общее количество товара: %d\n", totalItems)
				fmt.Println("Дополнительные стеллажи:")
				for _, shelf := range additionalShelves {
					fmt.Printf("- %s\n", shelf)
				}
				fmt.Println()

			case err := <-errChan:
				if err != nil {
					log.Fatal(err)
				}
				done <- true
				return
			}
		}
	}()

	go func() {
		defer close(itemChan)
		defer close(errChan)

		for rows.Next() {
			var mainShelf, productName string
			var orderID, totalItems int
			var additionalShelves pq.StringArray

			err := rows.Scan(&mainShelf, &orderID, &productName, &totalItems, &additionalShelves)
			if err != nil {
				errChan <- err
				return
			}

			item := []interface{}{mainShelf, orderID, productName, totalItems, additionalShelves}
			itemChan <- item
		}

		if err := rows.Err(); err != nil {
			errChan <- err
		}

		done <- true
	}()

	<-done
	return nil
}

func main() {
	// Подключение к базе данных
	db, err := sql.Open("postgres", "host=localhost port=5432 user=lew password=qwert dbname=sobes sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Добавление индекса на столбец order_id в таблице OrderItems
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_order_id ON OrderItems (order_id)")
	if err != nil {
		log.Fatal(err)
	}

	// Проверка наличия аргументов командной строки
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run . <order_id1> <order_id2> ...")
		os.Exit(1)
	}

	startTime := time.Now()

	// Обработка каждого переданного номера заказа для стеллажей
	mainShelves := []string{"Стеллаж А", "Стеллаж Б", "Стеллаж Ж"}
	for _, shelf := range mainShelves {
		first := true // Флаг для определения, был ли уже выведен основной стеллаж
		for _, arg := range os.Args[1:] {
			orderID, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("Неверный номер заказа: %s\n", arg)
				continue
			}

			query := `
				SELECT 
					ms.name AS main_shelf,
					oi.order_id,
					p.name AS product_name,
					(SELECT COUNT(*) FROM OrderItems oi2 WHERE oi2.product_id = p.id AND oi2.order_id = oi.order_id) AS total_items,
					(SELECT ARRAY(
							SELECT DISTINCT asl.name
							FROM AdditionalShelves asl
							WHERE asl.id IN (
								SELECT psr2.add_shelf_id
								FROM ProductShelfRelations psr2
								WHERE psr2.product_id = p.id
							)
						)) AS additional_shelves
				FROM 
					OrderItems oi,
					Products p,
					MainShelves ms,
					ProductShelfRelations psr
				WHERE 
					oi.product_id = p.id
					AND p.id = psr.product_id
					AND psr.main_shelf_id = ms.id
					AND oi.order_id = $1
					AND ms.name = $2
				GROUP BY 
					ms.name, oi.order_id, p.name, p.id;
			`

			rows, err := db.Query(query, orderID, shelf)
			if err != nil {
				fmt.Printf("Ошибка при выполнении запроса для заказа %d и стеллажа %s: %v\n", orderID, shelf, err)
				continue
			}
			defer rows.Close()

			err = printItemsByShelves(rows, first)
			if err != nil {
				fmt.Printf("Ошибка при обработке результатов запроса для заказа %d и стеллажа %s: %v\n", orderID, shelf, err)
			}
			first = false // Устанавливаем флаг в false после первого вызова функции
		}
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Время выполнения программы: %s\n", elapsed)
}

package middlewares

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"stocks-api/models"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type response struct {
	Success bool   `josn:"success"`
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error handling .env file")
	}

	db, err := sql.Open("postgres", os.Getenv(("POSTGRES_URL")))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to postgres database")

	return db

}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Error finding id, \n %v", err)
	}

	stock, err := findStock(int64(id))

	if err != nil {
		log.Fatalf("error finding stock\n %v", err)
	}

	json.NewEncoder(w).Encode(stock)

}

func GetAllStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := findAllStocks()
	if err != nil {
		log.Fatalf("error finding alll stocks\n %v", err)
	}

	json.NewEncoder(w).Encode(stocks)
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("unable to decode the request body. \n %v", err)
	}

	insertId := insertStock(stock)

	res := response{
		Success: true,
		ID:      insertId,
		Message: "stock inserted successfully",
	}

	json.NewEncoder(w).Encode(res)

}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to convert the string \n %v", err)
	}

	var stock models.Stock

	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("unable to decode the request\n %v", err)
	}

	updatedRow, _ := updateStock(int64(id), stock)

	res := response{
		Success: true,
		ID:      int64(id),
		Message: fmt.Sprintf("stock updated successfully \n %v", updatedRow),
	}

	json.NewEncoder(w).Encode(res)

}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to convert the string \n %v", err)
	}

	deletedRows := deleteOneStock(int64(id))
	res := response{
		Success: true,
		ID:      int64(id),
		Message: fmt.Sprintf("stock deleted successfully \n %v", deletedRows),
	}

	json.NewEncoder(w).Encode(res)

}

func insertStock(stock models.Stock) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStr := `INSERT INTO stocks(name, price, company) VALUES ($1, $2, $3) RETURNING stockid`

	var id int64

	err := db.QueryRow(sqlStr, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("unable to insert the stock \n %v", err)
	}
	fmt.Println("inserted a single stock with id ", id)
	return id
}

func findStock(id int64) (models.Stock, error) {
	db := CreateConnection()
	defer db.Close()
	sqlStr := `SELECT * FROM stocks WHERE stockid = $1`
	var stock models.Stock
	err := db.QueryRow(sqlStr, id).Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		return stock, errors.New("no rows found")
	case nil:
		return stock, nil
	default:
		fmt.Println("cannot scan row")
	}

	return stock, err

}

func findAllStocks() ([]models.Stock, error) {
	db := CreateConnection()
	defer db.Close()
	sqlStr := `SELECT * FROM stocks ORDER BY stockid DESC ;`

	var stocks []models.Stock
	rows, err := db.Query(sqlStr)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		err = rows.Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	// Check for any error encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	switch err {
	case sql.ErrNoRows:
		return stocks, errors.New("no rows found")
	case nil:
		return stocks, nil
	default:
		fmt.Println("cannot scan row")
	}

	return stocks, err

}

func updateStock(id int64, stock models.Stock) (models.Stock, error) {
	db := CreateConnection()
	defer db.Close()
	sqlStr := `UPDATE stocks SET name = $1, price = $2, company = $3 WHERE stockid = $4`
	var updatedStock models.Stock
	err := db.QueryRow(sqlStr, stock.Name, stock.Price, stock.Company, id).Scan(&updatedStock.Name, &updatedStock.Price, &updatedStock.Company, &updatedStock.StockId)

	switch err {
	case sql.ErrNoRows:
		return updatedStock, errors.New("no rows found")
	case nil:
		return updatedStock, nil
	default:
		fmt.Println("cannot scan row")
	}

	return updatedStock, err

}

func deleteOneStock(id int64) int64 {
	db := CreateConnection()
	defer db.Close()
	sqlStr := `DELETE FROM stocks WHERE stockid = $1`
	err := db.QueryRow(sqlStr, id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

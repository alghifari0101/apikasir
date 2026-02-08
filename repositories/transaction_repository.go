package repositories

import (
	"apikasir/models"
	"database/sql"
	"fmt"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM product WHERE id = $1 FOR UPDATE", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE product SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount, created_at) VALUES ($1, NOW()) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		err = tx.QueryRow("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal).Scan(&details[i].ID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) GetDailyReport() (*models.DailyReport, error) {
	var report models.DailyReport

	// 1. Total Revenue & Total Transaksi
	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(id) 
		FROM transactions 
		WHERE created_at::date = CURRENT_DATE
	`).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// 2. Produk Terlaris
	err = repo.db.QueryRow(`
		SELECT p.name, SUM(td.quantity) as qty
		FROM transaction_details td
		JOIN product p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE t.created_at::date = CURRENT_DATE
		GROUP BY p.name
		ORDER BY qty DESC
		LIMIT 1
	`).Scan(&report.ProdukTerlaris.Nama, &report.ProdukTerlaris.QtyTerjual)

	if err == sql.ErrNoRows {
		report.ProdukTerlaris.Nama = "-"
		report.ProdukTerlaris.QtyTerjual = 0
	} else if err != nil {
		return nil, err
	}

	return &report, nil
}

func (repo *TransactionRepository) GetReportByRange(startDate, endDate string) (*models.ReportResponse, error) {
	var response models.ReportResponse

	// Query untuk mengambil LIST transaksi
	query := `
        SELECT id, total_amount 
        FROM transactions 
        WHERE created_at::date BETWEEN $1 AND $2
        ORDER BY created_at DESC
    `
	rows, err := repo.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totalSales := 0
	transactions := []models.Transaction{}

	for rows.Next() {
		var t models.Transaction
		// Scan ke field yang sesuai
		if err := rows.Scan(&t.ID, &t.TotalAmount); err != nil {
			return nil, err
		}
		totalSales += t.TotalAmount
		transactions = append(transactions, t)
	}

	response.Transactions = transactions
	response.TotalAmount = totalSales

	return &response, nil
}

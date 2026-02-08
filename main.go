package main

import (
	"apikasir/database"
	"apikasir/handlers"
	"apikasir/repositories"
	"apikasir/services"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

// ubah Config
type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	viper.AutomaticEnv()

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize Layers
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Setup routes
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	// Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout) // POST

	// Report
	reportHandler := handlers.NewReportHandler(transactionService)
	http.HandleFunc("/api/report/hari-ini", reportHandler.GetDailyReport)
	http.HandleFunc("/api/report", reportHandler.GetReportByRange)

	// Start Server
	port := config.Port
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

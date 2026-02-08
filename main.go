package main

import (
	"api-pos/database"
	handlers "api-pos/handler"
	"api-pos/repositories"
	"api-pos/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBconn string `mapstructure:"DB_CONN"`
}

func main() {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "-"))

	if _, err := os.Stat(".env"); err == nil {

		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBconn: viper.GetString("DB_CONN"),
	}

	db, err := database.InitDB(config.DBconn)
	if err != nil {
		log.Fatal("Failed connect DB", err)
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	http.HandleFunc("/", productHandler.RootHandler)
	http.HandleFunc("/product", productHandler.ProductHandler)
	http.HandleFunc("/product/", productHandler.ProductByIDHandler)
	http.HandleFunc("/category", categoryHandler.CategoryHandler)
	http.HandleFunc("/category/", categoryHandler.CategoryByIDHandler)
	http.HandleFunc("/checkout/", transactionHandler.HandleCheckout)
	http.HandleFunc("/report", transactionHandler.GetReport)
	http.HandleFunc("/report/today", transactionHandler.GetReport)

	fmt.Printf("Server running on port %s\n", config.Port)

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}

}

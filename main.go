package main

import (
    "encoding/json"
    "net/http"
)

type Product struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Category    string  `json:"category"`
}

var beautyProducts = []Product{
    {
        ID:          1,
        Name:        "Facial Moisturizer",
        Description: "Hydrating daily moisturizer for all skin types",
        Price:       29.99,
        Category:    "Skincare",
    },
    {
        ID:          2,
        Name:        "Volumizing Mascara",
        Description: "Long-lasting mascara for dramatic lashes",
        Price:       19.99,
        Category:    "Makeup",
    },
    {
        ID:          3,
        Name:        "Hair Serum",
        Description: "Anti-frizz hair treatment with argan oil",
        Price:       24.99,
        Category:    "Hair Care",
    },
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(beautyProducts)
}

func main() {
    http.HandleFunc("/products", productsHandler)
    http.ListenAndServe(":8080", nil)
}
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "time"
)

type Product struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Category    string  `json:"category"`
}

type Address struct {
    CEP         string `json:"cep"`
    Street      string `json:"logradouro"`
    Complement  string `json:"complemento"`
    District    string `json:"bairro"`
    City        string `json:"localidade"`
    State       string `json:"uf"`
    Source      string `json:"source"`
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

func fetchViaCEP(ctx context.Context, cep string, resultChan chan<- Address) {
    url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
    
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    var address Address
    if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
        return
    }
    
    address.Source = "ViaCEP"
    resultChan <- address
}

func fetchBrasilAPI(ctx context.Context, cep string, resultChan chan<- Address) {
    url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
    
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    var address Address
    if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
        return
    }
    
    address.Source = "BrasilAPI"
    resultChan <- address
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(beautyProducts)
}

func searchCEP(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    cep := strings.ReplaceAll(r.URL.Query().Get("cep"), "-", "")
    if len(cep) != 8 {
        http.Error(w, "CEP inválido", http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    resultChan := make(chan Address, 2)

    // Inicia as duas requisições concorrentemente
    go fetchViaCEP(ctx, cep, resultChan)
    go fetchBrasilAPI(ctx, cep, resultChan)

    // Aguarda o primeiro resultado ou timeout
    select {
    case result := <-resultChan:
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(result)
    case <-ctx.Done():
        http.Error(w, "Timeout ao buscar CEP", http.StatusGatewayTimeout)
    }
}

func main() {
    http.HandleFunc("/products", productsHandler)
    http.HandleFunc("/cep", searchCEP)
    fmt.Println("Server running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
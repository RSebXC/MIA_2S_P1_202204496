package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "fmt"
    analyzer"MIA_202204496/analyzer"
)

type Request struct {
    InputText string `json:"inputText"`
}

type Response struct {
    Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    // Configurar los encabezados CORS
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // Manejar la solicitud OPTIONS
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

    var req Request
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Can't read body", http.StatusBadRequest)
        return
    }
    err = json.Unmarshal(body, &req)
    if err != nil {
        http.Error(w, "Can't unmarshal JSON", http.StatusBadRequest)
        return
    }
    fmt.Println(req.InputText)
    analyzer.Analyzer(req.InputText)
    response := Response{Message: "Received: " + req.InputText}
    json.NewEncoder(w).Encode(response)
}

func main() {
    http.HandleFunc("/api/message", handler)
    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

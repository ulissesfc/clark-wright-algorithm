package main

import (
	"fmt"
	"net/http"

	"github.com/ulissesfc/clark-wright-algorithm.git/internal/handler"
)

func main() {

	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/", fs)

	http.HandleFunc("/solve", handler.SolveHandler)

	fmt.Println("rodando na porta 8080")
	fmt.Println("Acesse http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

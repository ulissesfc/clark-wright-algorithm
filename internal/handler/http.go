package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ulissesfc/clark-wright-algorithm.git/internal/clarkewright"
	"github.com/ulissesfc/clark-wright-algorithm.git/internal/graph"
)

func SolveHandler(w http.ResponseWriter, r *http.Request) {

	var input graph.InputGraph
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	solver, err := clarkewright.NewSolver(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	solution, err := solver.Solve()
	if err != nil {
		http.Error(w, "Erro ao resolver o grafo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(solution)
}

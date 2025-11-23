package graph

type Node struct {
	ID     string `json:"id"`
	Demand int    `json:"demand"`
}

type InputGraph struct {
	Nodes           []Node                        `json:"nodes"`
	DistanceMatrix  map[string]map[string]float64 `json:"distance_matrix"` // De -> Para -> Distância
	VehicleCapacity int                           `json:"vehicle_capacity"`
}

type Route struct {
	Sequence      []string `json:"sequence"` // Ex: ["Depot", "A", "B", "Depot"]
	TotalDistance float64  `json:"total_distance"`
	TotalLoad     int      `json:"total_load"`
}

func GetDistance(matrix map[string]map[string]float64, from, to string) float64 {
	if val, ok := matrix[from][to]; ok {
		return val
	}
	// Assume simetria se não encontrar (A->B é igual B->A)
	if val, ok := matrix[to][from]; ok {
		return val
	}
	return 0.0
}

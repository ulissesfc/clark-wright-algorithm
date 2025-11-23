package clarkewright

import (
	"errors"
	"sort"

	"github.com/ulissesfc/clark-wright-algorithm.git/internal/graph"
)

type Saving struct {
	i, j  string
	Score float64
}

// Solver mantém o estado da execução do algoritmo
type Solver struct {
	input     graph.InputGraph
	depot     graph.Node
	customers []graph.Node
	savings   []Saving
	routed    map[string]bool // Set para verificação de visitas
}

// NewSolver cria a instância
func NewSolver(input graph.InputGraph) (*Solver, error) {
	if len(input.Nodes) < 2 {
		return nil, errors.New("o grafo deve conter pelo menos 1 depósito e 1 cliente")
	}

	// Assume que o primeiro nó é o depósito
	return &Solver{
		input:     input,
		depot:     input.Nodes[0],
		customers: input.Nodes[1:],
		routed:    make(map[string]bool),
	}, nil
}

// Solve é o método público que orquestra as etapas
func (s *Solver) Solve() ([]graph.Route, error) {
	// 1. Calcula a matriz de economias (s_ij)
	s.calculateSavings()

	// 2. Ordena economias da maior para a menor
	s.sortSavings()

	// 3. Constrói as rotas sequencialmente
	routes := s.buildRoutes()

	return routes, nil
}

// --- Métodos Internos ---

func (s *Solver) calculateSavings() {
	for i := 0; i < len(s.customers); i++ {
		for j := i + 1; j < len(s.customers); j++ {
			custI := s.customers[i]
			custJ := s.customers[j]

			d0i := graph.GetDistance(s.input.DistanceMatrix, s.depot.ID, custI.ID)
			d0j := graph.GetDistance(s.input.DistanceMatrix, s.depot.ID, custJ.ID)
			dij := graph.GetDistance(s.input.DistanceMatrix, custI.ID, custJ.ID)

			score := d0i + d0j - dij
			s.savings = append(s.savings, Saving{i: custI.ID, j: custJ.ID, Score: score})
		}
	}
}

func (s *Solver) sortSavings() {
	sort.Slice(s.savings, func(i, j int) bool {
		return s.savings[i].Score > s.savings[j].Score
	})
}

func (s *Solver) buildRoutes() []graph.Route {
	var routes []graph.Route
	nodesRemaining := len(s.customers)

	for nodesRemaining > 0 {
		// Cria uma nova rota
		routeNodes, load := s.buildSingleRoute()

		// Atualiza contadores
		nodesRemaining -= len(routeNodes)

		// Finaliza a rota (adiciona depósito nas pontas e calcula custo)
		routes = append(routes, s.finalizeRoute(routeNodes, load))
	}
	return routes
}

// buildSingleRoute tenta encher UM veículo até a capacidade máxima ou acabar as economias
func (s *Solver) buildSingleRoute() (nodes []string, load int) {
	// Começa uma rota vazia (apenas lógica, sem depósito ainda)
	// Na lógica sequencial pura, geralmente pega o melhor saving disponível para começar a rota.
	// Se não houver saving válido, pega um nó isolado.

	// Procura o primeiro par de saving válido para INICIAR a rota
	startPairFound := false

	for k, sav := range s.savings {
		if sav.Score == -1 || s.routed[sav.i] || s.routed[sav.j] {
			continue
		}

		demandI := s.getDemand(sav.i)
		demandJ := s.getDemand(sav.j)

		if demandI+demandJ <= s.input.VehicleCapacity {
			nodes = []string{sav.i, sav.j}
			load = demandI + demandJ
			s.routed[sav.i] = true
			s.routed[sav.j] = true
			s.savings[k].Score = -1 // Marca como usado
			startPairFound = true
			break
		}
	}

	// Se não achou par para começar, pega o primeiro cliente não visitado isolado
	if !startPairFound {
		for _, c := range s.customers {
			if !s.routed[c.ID] {
				if s.getDemand(c.ID) <= s.input.VehicleCapacity {
					nodes = []string{c.ID}
					load = s.getDemand(c.ID)
					s.routed[c.ID] = true
					return // Rota de um nó só
				}
			}
		}
		return // Não sobrou ninguém (segurança)
	}

	// Agora tenta estender essa rota pelas pontas (esquerda/direita)
	// Enquanto couber carga e houver savings
	for {
		added := false
		left := nodes[0]
		right := nodes[len(nodes)-1]

		// Busca o melhor saving que conecte em Left ou Right
		for k, sav := range s.savings {
			if sav.Score == -1 {
				continue
			}

			// Verifica candidatos
			var candidate string

			// Logica de conexão nas extremidades
			if sav.i == right && !s.routed[sav.j] {
				candidate = sav.j
			}
			if sav.j == right && !s.routed[sav.i] {
				candidate = sav.i
			}
			if sav.i == left && !s.routed[sav.j] {
				candidate = sav.j
			}
			if sav.j == left && !s.routed[sav.i] {
				candidate = sav.i
			}

			if candidate != "" {
				dem := s.getDemand(candidate)
				if load+dem <= s.input.VehicleCapacity {
					// Insere na ponta correta
					if candidate == sav.i && sav.j == left || candidate == sav.j && sav.i == left {
						nodes = append([]string{candidate}, nodes...) // Prepend
					} else {
						nodes = append(nodes, candidate) // Append
					}

					s.routed[candidate] = true
					load += dem
					s.savings[k].Score = -1
					added = true
					break // Volta para o loop externo para reavaliar as pontas novas
				}
			}
		}

		if !added {
			break
		} // Se rodou a lista toda e não adicionou ninguém, a rota fecha.
	}

	return
}

func (s *Solver) finalizeRoute(clientNodes []string, load int) graph.Route {
	// Adiciona depósito no início e fim
	seq := append([]string{s.depot.ID}, clientNodes...)
	seq = append(seq, s.depot.ID)

	dist := 0.0
	for i := 0; i < len(seq)-1; i++ {
		dist += graph.GetDistance(s.input.DistanceMatrix, seq[i], seq[i+1])
	}

	return graph.Route{
		Sequence:      seq,
		TotalDistance: dist,
		TotalLoad:     load,
	}
}

func (s *Solver) getDemand(id string) int {
	for _, n := range s.input.Nodes {
		if n.ID == id {
			return n.Demand
		}
	}
	return 0
}

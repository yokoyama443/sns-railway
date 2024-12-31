package main

import (
	"container/heap"
	"database/sql"
	"math"
)

const (
	EARTH_RADIUS  = 6371
	INF           = math.MaxFloat64
	MAX_DISTANCE  = 10.0
	TRANSFER_COST = 10.0 // 乗換ペナルティ
)

// Station represents a railway station
type Station struct {
	StationCD   int
	StationGCD  int
	StationName string
	LineCD      int
	Lon         float64
	Lat         float64
}

type RouteRequest struct {
	StartCD string `json:"startCD"`
	EndCD   string `json:"endCD"`
}

type StationPath struct {
	StationName string  `json:"stationName"`
	LineCD      int     `json:"lineCD"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Transfer    bool    `json:"transfer"`
}

type RouteResponse struct {
	Path          []StationPath `json:"path"`
	TotalDistance float64       `json:"totalDistance"`
}

type Item struct {
	stationCD int
	priority  float64
	index     int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].priority < pq[j].priority }
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func getStations(db *sql.DB) (map[int]*Station, error) {
	stations := make(map[int]*Station)

	rows, err := db.Query(`
		SELECT station_cd, station_g_cd, station_name, line_cd, lon, lat
		FROM m_station
		WHERE e_status IS NULL OR e_status = 0
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := &Station{}
		err = rows.Scan(
			&s.StationCD,
			&s.StationGCD,
			&s.StationName,
			&s.LineCD,
			&s.Lon,
			&s.Lat,
		)
		if err != nil {
			return nil, err
		}
		stations[s.StationCD] = s
	}

	return stations, nil
}

func buildGraph(stations map[int]*Station) map[int]map[int]float64 {
	graph := make(map[int]map[int]float64)

	for id := range stations {
		graph[id] = make(map[int]float64)
	}

	for id1, s1 := range stations {
		for id2, s2 := range stations {
			if id1 == id2 {
				continue
			}

			if s1.LineCD == s2.LineCD {
				dist := calcDistance(s1.Lat, s1.Lon, s2.Lat, s2.Lon)
				if dist < MAX_DISTANCE {
					graph[id1][id2] = dist
				}
			}

			if s1.StationGCD == s2.StationGCD && s1.StationGCD != 0 {
				graph[id1][id2] = TRANSFER_COST
			}
		}
	}

	return graph
}

// ダイクストラ法で最短経路を求める
func dijkstra(graph map[int]map[int]float64, start int) (map[int]float64, map[int]int) {
	dist := make(map[int]float64)
	prev := make(map[int]int)

	for id := range graph {
		dist[id] = INF
	}
	dist[start] = 0

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item{stationCD: start, priority: 0})

	for pq.Len() > 0 {
		u := heap.Pop(&pq).(*Item)

		if u.priority > dist[u.stationCD] {
			continue
		}

		for v, weight := range graph[u.stationCD] {
			alt := dist[u.stationCD] + weight
			if alt < dist[v] {
				dist[v] = alt
				prev[v] = u.stationCD
				heap.Push(&pq, &Item{stationCD: v, priority: alt})
			}
		}
	}

	return dist, prev
}

// 最短経路の復元
func reconstructPath(prev map[int]int, start, end int) []int {
	path := make([]int, 0)
	current := end

	for current != start {
		path = append([]int{current}, path...)
		var ok bool
		current, ok = prev[current]
		if !ok {
			return nil
		}
	}
	path = append([]int{start}, path...)

	return path
}

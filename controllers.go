package main

import (
	"encoding/json"
	"net/http"
)

func handleRoute(w http.ResponseWriter, r *http.Request) {
	var req RouteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stations, err := getStations(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	start := stations[atoi(req.StartCD)]
	end := stations[atoi(req.EndCD)]

	if start == nil || end == nil {
		http.Error(w, "指定された駅が見つかりません。", http.StatusNotFound)
		return
	}

	graph := buildGraph(stations)
	dist, prev := dijkstra(graph, start.StationCD)
	path := reconstructPath(prev, start.StationCD, end.StationCD)

	if path == nil {
		http.Error(w, "経路が見つかりませんでした。", http.StatusNotFound)
		return
	}

	resp := buildRouteResponse(path, dist[end.StationCD], stations)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func buildRouteResponse(path []int, totalDistance float64, stations map[int]*Station) RouteResponse {
	resp := RouteResponse{
		TotalDistance: totalDistance,
	}

	prevLine := -1
	for i, stationCD := range path {
		station := stations[stationCD]
		pathItem := StationPath{
			StationName: station.StationName,
			LineCD:      station.LineCD,
			Lat:         station.Lat,
			Lon:         station.Lon,
			Transfer:    false,
		}

		if i > 0 && prevLine != station.LineCD {
			if stations[path[i-1]].StationGCD == station.StationGCD {
				pathItem.Transfer = true
			}
		}
		prevLine = station.LineCD

		resp.Path = append(resp.Path, pathItem)
	}

	return resp
}

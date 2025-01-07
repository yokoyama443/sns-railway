package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RouteRequest struct {
	StartCD string `json:"startCD"`
	EndCD   string `json:"endCD"`
	Allow78 bool   `json:"allow78"`
}

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

	allow78flag := req.Allow78
	fmt.Println("allow78flag: ", allow78flag)

	graph, err := buildGraph(stations, db, allow78flag)
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

	dist, prev := dijkstra(graph, start.StationCD)
	path := reconstructPath(prev, start.StationCD, end.StationCD)
	if path == nil {
		http.Error(w, "経路が見つかりませんでした。", http.StatusNotFound)
		return
	}

	resp := buildRouteResponse(path, dist[end.StationCD], stations)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
				resp.TotalDistance -= TRANSFER_COST
			}
		}
		prevLine = station.LineCD

		resp.Path = append(resp.Path, pathItem)
	}

	return resp
}

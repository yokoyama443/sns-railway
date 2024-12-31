package main

import (
	"fmt"
	"math"
)

func atoi(s string) int {
	var n int
	if _, err := fmt.Sscanf(s, "%d", &n); err != nil {
		return 0
	}
	return n
}

// Haversineの公式を使って2点間の距離を計算
func calcDistance(lat1, lon1, lat2, lon2 float64) float64 {
	lat1 = lat1 * math.Pi / 180
	lon1 = lon1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180
	lon2 = lon2 * math.Pi / 180

	h := math.Pow(math.Sin((lat2-lat1)/2), 2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Pow(math.Sin((lon2-lon1)/2), 2)

	return 2 * EARTH_RADIUS * math.Asin(math.Sqrt(h))
}

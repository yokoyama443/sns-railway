let map;

// 初期化処理
window.onload = () => {
    map = L.map("mapid").setView([35.681236, 139.767125], 11);
    L.tileLayer('https://cyberjapandata.gsi.go.jp/xyz/std/{z}/{x}/{y}.png', {
        attribution: "<a href='https://maps.gsi.go.jp/development/ichiran.html' target='_blank'>地理院タイル</a>",
        minZoom: 5,
        maxZoom: 18,
    }).addTo(map);

    document.getElementById('routeForm').addEventListener('submit', handleFormSubmit);
};


// lineCD->色
function getLineColor(lineCD) {
    const r = (lineCD * 95413) % 256;
    const g = (lineCD * 37237) % 256;
    const b = (lineCD * 49279) % 256;
    return `rgb(${r}, ${g}, ${b})`;
}

// フォーム送信ハンドラ
async function handleFormSubmit(e) {
    e.preventDefault();
    
    const startStation = document.getElementById('startStation').value;
    const endStation = document.getElementById('endStation').value;
    const allow78 = document.querySelector('input[name="allow78"]:checked').value === 'true';

    try {
        const response = await fetch('/api/route', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                startCD: startStation,
                endCD: endStation,
                allow78: !allow78,
            })
        });

        if (!response.ok) {
            throw new Error('経路の取得に失敗しました');
        }

        const data = await response.json();
        drawRoute(data);
    } catch (error) {
        alert(error.message);
    }
}


// 経路を地図上に表示
function drawRoute(routeData) {
    clearExistingRoute();
    drawRouteInfo(routeData);
    drawStations(routeData);
    adjustMapBounds(routeData);
}

// 既存の経路をクリア
function clearExistingRoute() {
    map.eachLayer((layer) => {
        if (layer instanceof L.Polyline || layer instanceof L.Marker) {
            map.removeLayer(layer);
        }
    });
}

// 経路情報の表示
function drawRouteInfo(routeData) {
    const routeInfo = document.getElementById('routeInfo');
    routeInfo.innerHTML = `
        <h3>経路案内</h3>
        <p>総距離: ${routeData.totalDistance.toFixed(2)} km</p>
        <div id="stationList"></div>
    `;
}

// 駅と線路を表示
function drawStations(routeData) {
    const stationList = document.getElementById('stationList');
    console.log(routeData);
    let currentLine = null;
    let currentPoints = [];
    
    routeData.path.forEach((station, index) => {
        addStationMarker(station);
        addStationToList(station, index, stationList);
        currentPoints = updateRouteLine(station, currentLine, currentPoints);
        currentLine = station.lineCD;
    });

    if (currentPoints.length > 0) {
        drawRouteLine(currentPoints, currentLine);
    }
}

// 駅マーカーの追加
function addStationMarker(station) {
    L.marker([station.lat, station.lon])
        .addTo(map)
        .bindPopup(station.stationName === '東高円寺' ? '東高円寺 (制作者の最寄り)' : station.stationName);
}

// 駅リストへの追加
function addStationToList(station, index, stationList) {
    const stationDiv = document.createElement('div');
    stationDiv.className = 'station';
    stationDiv.innerHTML = `${index + 1}. ${station.stationName} (路線: ${station.lineCD})`;
    if (station.transfer) {
        stationDiv.innerHTML += '<div class="transfer">※ 乗り換え</div>';
    }
    stationList.appendChild(stationDiv);
}

// 路線の更新と描画
function updateRouteLine(station, currentLine, currentPoints) {
    if (currentLine !== station.lineCD && currentPoints.length > 0) {
        drawRouteLine(currentPoints, currentLine);
        return [[station.lat, station.lon]];
    }
    currentPoints.push([station.lat, station.lon]);
    return currentPoints;
}

function drawRouteLine(points, lineCD) {
    // 線路がメトロセブン・エイトライナーのときは二重線
    if (lineCD === 9999) {
        L.polyline(points, {
            color: 'white',
            weight: 2,
            opacity: 1
        }).addTo(map);
        L.polyline(points, {
            color: getLineColor(lineCD),
            weight: 6,
            opacity: 0.8
        }).addTo(map);
    } else {
        L.polyline(points, {
            color: getLineColor(lineCD),
            weight: 4,
            opacity: 0.8
        }).addTo(map);
    }
}



// 地図の表示範囲調整
function adjustMapBounds(routeData) {
    const points = routeData.path.map(station => [station.lat, station.lon]);
    if (points.length > 0) {
        map.fitBounds(points);
    }
}

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Response struct {
	Name          string `json:"name"`
	Keterangan    string `json:"keterangan"`
	Status        Status `json:"status"`
	StatusBencana string `json:"status_bencana"`
}

type Status struct {
	Water int64 `json:"water"`
	Wind  int64 `json:"wind"`
}

var (
	numberSumsMutex sync.RWMutex
	angkaWater      int64
	angkaWind       int64
)

func StatusWaterWind(w http.ResponseWriter, r *http.Request) {
	numberSumsMutex.RLock()
	defer numberSumsMutex.RUnlock()

	var stBencana string
	if angkaWater < 5 && angkaWind < 6 {
		stBencana = "Aman"
	} else if angkaWater < 8 && angkaWind < 16 {
		stBencana = "Siaga"
	} else if angkaWater > 8 && angkaWind > 15 {
		stBencana = "Bahaya"
	} else {
		stBencana = "Bahaya"
	}

	var data2 = []Response{
		{
			"Cek Status Bencana Setiap 15 Detik",
			"Hanya Tekan Send request/Refresh Untuk Update status",
			Status{angkaWater, angkaWind},
			stBencana,
		},
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var result, err = json.Marshal(data2)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
	}
	http.Error(w, "", http.StatusBadRequest)

	// fmt.Println(angkaWater)
	// fmt.Println(angkaWind)
}

func main() {
	fmt.Println("web berjalan di server http://localhost:8080/")
	go runDataLoop()
	http.HandleFunc("/cekstatus", StatusWaterWind)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func runDataLoop() {
	for {
		numberSumsMutex.Lock()
		angkaWater = int64(rand.Intn(100))
		angkaWind = int64(rand.Intn(100))
		numberSumsMutex.Unlock()
		time.Sleep(15 * time.Second)
	}
}

package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sort"
    "sync"
    "time"
)

type SortRequest struct {
    ToSort [][]int `json:"to_sort"`
}

type SortResponse struct {
    SortedArrays [][]int `json:"sorted_arrays"`
    TimeNS      int64   `json:"time_us"`
}

func main() {
    http.HandleFunc("/process-single", processSingle)
    http.HandleFunc("/process-concurrent", processConcurrent)

    port := 8000
    fmt.Printf("Server is listening on port %d...\n", port)
    err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
    if err != nil {
        fmt.Println("Error starting server:", err)
    }
}

func processSingle(w http.ResponseWriter, r *http.Request) {
    processRequest(w, r, false)
}

func processConcurrent(w http.ResponseWriter, r *http.Request) {
    processRequest(w, r, true)
}

func processRequest(w http.ResponseWriter, r *http.Request, concurrent bool) {
    var req SortRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
        return
    }

    startTime := time.Now()

    var sortedArrays [][]int
    if concurrent {
        sortedArrays = sortConcurrently(req.ToSort)
    } else {
        sortedArrays = sortSequentially(req.ToSort)
    }

    timeTaken := time.Since(startTime).Microseconds()

    response := SortResponse{
        SortedArrays: sortedArrays,
        TimeNS:       timeTaken,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func sortSequentially(arrays [][]int) [][]int {
    var sortedArrays [][]int
    for _, arr := range arrays {
        sortedArr := make([]int, len(arr))
        copy(sortedArr, arr)
        sort.Ints(sortedArr)
        sortedArrays = append(sortedArrays, sortedArr)
    }
    return sortedArrays
}

func sortConcurrently(arrays [][]int) [][]int {
    var wg sync.WaitGroup
    var mu sync.Mutex
    var sortedArrays [][]int

    for _, arr := range arrays {
        wg.Add(1)
        go func(a []int) {
            defer wg.Done()

            sortedArr := make([]int, len(a))
            copy(sortedArr, a)
            sort.Ints(sortedArr)

            mu.Lock()
            sortedArrays = append(sortedArrays, sortedArr)
            mu.Unlock()
        }(arr)
    }

    wg.Wait()
    return sortedArrays
}

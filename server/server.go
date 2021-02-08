package server

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func GenerateSampleItems() map[string][]Item {
	return map[string][]Item{
		"Fruit": {
			{ID: 1, Name: "Apple", ExpiryDate: time.Now().Add(24 * time.Hour), Value: 1.00, Category: "Fruit"},
			{ID: 2, Name: "Banana", ExpiryDate: time.Now().Add(24 * time.Hour), Value: 2.50, Category: "Fruit"},
			{ID: 3, Name: "Blueberry", ExpiryDate: time.Now().Add(24 * time.Hour), Value: 2.00, Category: "Fruit"},
			{ID: 4, Name: "Orange", ExpiryDate: time.Now().Add(24 * time.Hour), Value: 1.50, Category: "Fruit"},
		},
		"Electronics": {
			{ID: 5, Name: "iPad", ExpiryDate: time.Now().Add(24 * time.Hour), Value: 300, Category: "Electronics"},
			{ID: 6, Name: "Macbook", ExpiryDate: time.Now().Add(24 * time.Hour), Value: 1500, Category: "Electronics"},
			{ID: 7, Name: "Headphones", ExpiryDate: time.Now().Add(24 * time.Hour), Value: 200, Category: "Electronics"},
			{ID: 8, Name: "Microphone", ExpiryDate: time.Now().Add(24 * time.Hour), Value: 40, Category: "Electronics"},
			{ID: 9, Name: "Monitor", ExpiryDate: time.Now().Add(24 * time.Hour), Value: 200, Category: "Electronics"},
		},
		"Clothes": {
			{ID: 10, Name: "Shirt", ExpiryDate: time.Now().Add(24 * time.Hour), Value: 50, Category: "Clothes"},
			{ID: 11, Name: "Pants", ExpiryDate: time.Now().Add(24 * time.Hour), Value: 40, Category: "Clothes"},
			{ID: 12, Name: "Belt", ExpiryDate: time.Now().Add(24 * time.Hour), Value: 20, Category: "Clothes"},
			{ID: 13, Name: "Shoes", ExpiryDate: time.Now().Add(24 * time.Hour), Value: 50, Category: "Clothes"},
			{ID: 14, Name: "Hat", ExpiryDate: time.Now().Add(24 * time.Hour), Value: 15, Category: "Clothes"},
			{ID: 15, Name: "Glasses", ExpiryDate: time.Now().Add(24 * time.Hour), Value: 25, Category: "Clothes"},
		},
	}
}

type APIConfig struct {
	items    map[string][]Item
	allItems []Item
}

func NewAPIConfig(items map[string][]Item) *APIConfig {
	var allItems []Item
	for _, value := range items {
		allItems = append(allItems, value...)
	}

	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].ID < allItems[j].ID
	})
	return &APIConfig{items: items, allItems: allItems}
}

func (a *APIConfig) Init() http.Handler {
	s := http.NewServeMux()
	rand.New(rand.NewSource(time.Now().UnixNano()))
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		count := r.URL.Query().Get("count")
		time.Sleep(time.Duration(time.Duration(rand.Intn(5)) * time.Second))
		w.Write([]byte(fmt.Sprintf("Finished request %s", count)))
	})

	s.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		currentPage, _ := strconv.Atoi(r.URL.Query().Get("page"))
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
		category := r.URL.Query().Get("category")

		var results []Item
		if category != "" {
			results = a.items[category]
		} else {
			results = a.allItems
		}

		totalItems := len(results)
		totalPages := math.Ceil(float64(totalItems) / float64(pageSize))
		indexStart := (currentPage - 1) * pageSize
		indexFinish := pageSize * currentPage

		//Establishing Bounds for Indexes and pages
		if indexStart > totalItems || totalItems == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Page not found."))
			return
		}

		if indexFinish > totalItems {
			indexFinish = totalItems
		}

		results = results[indexStart:indexFinish]
		itemsResponse := ItemResponse{Items: results, PageNumber: currentPage, PageSize: pageSize, TotalPages: int(totalPages), TotalItems: totalItems}
		response, _ := json.Marshal(itemsResponse)
		w.Write(response)
	})

	return s
}

// Item - JSON schema
type Item struct {
	ID         int       `json:"id"`
	ExpiryDate time.Time `json:"expiry_date"`
	Value      float64   `json:"value"`
	Name       string    `json:"name"`
	Category   string    `json:"category"`
}

// ItemResponse - JSON schema
type ItemResponse struct {
	Items      []Item `json:"items"`
	PageNumber int    `json:"page_number"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
	TotalItems int    `json:"total_items"`
}

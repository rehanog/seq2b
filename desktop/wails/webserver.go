// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Simple web server that exposes the App methods via HTTP
func main() {
	app := NewApp()
	
	// Load the test directory
	testDir := "../../testdata/library_test_0/pages"
	if absDir, err := filepath.Abs(testDir); err == nil {
		app.LoadDirectory(absDir)
	}

	// API endpoints
	http.HandleFunc("/api/getPage", func(w http.ResponseWriter, r *http.Request) {
		pageName := r.URL.Query().Get("name")
		page, err := app.GetPage(pageName)
		
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(page)
	})

	http.HandleFunc("/api/getPageList", func(w http.ResponseWriter, r *http.Request) {
		pages := app.GetPageList()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(pages)
	})

	http.HandleFunc("/api/updateBlock", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		var req struct {
			PageName string `json:"pageName"`
			BlockID  string `json:"blockId"`
			Content  string `json:"content"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		err := app.UpdateBlock(req.PageName, req.BlockID, req.Content)
		
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Serve frontend files
	fs := http.FileServer(http.Dir("./frontend/dist"))
	http.Handle("/", fs)

	fmt.Println("Web server running at http://localhost:8080")
	fmt.Println("API endpoints:")
	fmt.Println("  GET  /api/getPage?name=PageName")
	fmt.Println("  GET  /api/getPageList")
	fmt.Println("  POST /api/updateBlock")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
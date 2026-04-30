package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

var globalPlayer *Player

func StartServer(p *Player) {
	globalPlayer = p

	// Serve API
	http.HandleFunc("/api/player", playerHandler)

	// Serve Static Files
	public, _ := fs.Sub(staticFiles, "static")
	http.Handle("/", http.FileServer(http.FS(public)))

	fmt.Println("🌐 Dashboard live at: http://localhost:8080")
	
	// Run server in background
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("❌ Web Server Error: %v\n", err)
		}
	}()
}

func playerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if globalPlayer == nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(globalPlayer)
}

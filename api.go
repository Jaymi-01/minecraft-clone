package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

var globalPlayer *Player

func StartServer(p *Player) {
	globalPlayer = p
	
	// Start Game Cycles
	p.StartRegeneration()
	p.StartRaids()
	p.StartGateSpawning()
	p.StartSubordinateAutonomy()

	// Register Handlers
	http.HandleFunc("/api/player", playerHandler)
	public, _ := fs.Sub(staticFiles, "static")
	http.Handle("/", http.FileServer(http.FS(public)))

	// Find an available port starting from 8080
	port := 8080
	var listener net.Listener
	var err error

	for {
		addr := fmt.Sprintf(":%d", port)
		listener, err = net.Listen("tcp", addr)
		if err == nil {
			break
		}
		port++
		if port > 8100 { // Safety break
			fmt.Printf("❌ Failed to find an open port after 20 attempts: %v\n", err)
			return
		}
	}

	fmt.Printf("🌐 Dashboard live at: http://localhost:%d\n", port)
	
	// Run server in background
	go func() {
		if err := http.Serve(listener, nil); err != nil {
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

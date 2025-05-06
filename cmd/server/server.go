package main

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"

	_ "github.com/matetirpak/chessbot-playground-server/docs"
	"github.com/matetirpak/chessbot-playground-server/pkg/server"
	"github.com/matetirpak/chessbot-playground-server/web"
)

// @title           Chessbot Playground API
// @version         1.0
// @description     This is a RESTful API for a chess server featuring parallel games for bot experimentation.

// @contact.name   Mate Tirpak
// @contact.url    https://github.com/matetirpak
// @contact.email  mate.tirpak@gmail.com

// @license.name  MIT
// @license.url   https://github.com/matetirpak/chessbot-playground-server/blob/main/LICENSE

// @BasePath  /chessserver/v1
func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	srvApi := initApiServer(":8080")
	srvFrontend := initHttpServer(":8081")

	go func() {
		log.Println("Serving API on http://localhost:8080")
		if err := srvApi.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("API server error: %v", err)
		}
	}()

	go func() {
		log.Println("Serving frontend on http://localhost:8081")
		if err := srvFrontend.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Frontend server error: %v", err)
		}
	}()

	// Wait for termination signal (Ctrl+C)
	sig := <-signalChan
	log.Printf("%s signal caught", sig)
	log.Println("Shutdown signal received, shutting down servers...")

	closeHttp(srvFrontend)
	closeHttp(srvApi)

	log.Println("Servers exited.")
}

func initApiServer(port string) *http.Server {
	log.Printf("Server started")

	router := server.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8081",
		},
		AllowedMethods:   []string{"DELETE", "GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	srv := &http.Server{
		Addr:    port,
		Handler: handler,
	}

	return srv
}

func initHttpServer(port string) *http.Server {
	subFS, err := fs.Sub(web.EmbeddedWebFiles, "webui/src")
	if err != nil {
		log.Printf("Error when trying to access frontend subdirectory: %v", err)
		panic(err)
	}

	fs := http.FileServer(http.FS(subFS))

	handler := func(w http.ResponseWriter, r *http.Request) {
		// Add necessary headers for SharedArrayBuffer support
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")

		// Serve the requested file
		fs.ServeHTTP(w, r)
	}

	srv := &http.Server{
		Addr:    port,
		Handler: http.HandlerFunc(handler),
	}

	return srv
}

func closeHttp(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown failed: %+v", err)
	}
}

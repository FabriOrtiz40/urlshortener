package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Mapa de paths cortos a URLs reales
	pathsToUrls := map[string]string{
		"/dogs": "https://en.wikipedia.org/wiki/Dog",
		"/cats": "https://en.wikipedia.org/wiki/Cat",
	}

	// Handler de fallback si no coincide ningún path
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hola! Esta página no existe.")
	})

	// Crea el handler principal usando tu MapHandler
	handler := MapHandler(pathsToUrls, fallback)

	// Inicia el servidor en el puerto 8080
	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}

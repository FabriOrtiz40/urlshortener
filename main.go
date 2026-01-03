package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hola! Esta página no existe.")
	})

	yamlFile, err := os.ReadFile("paths.yaml")
	if err != nil {
		log.Fatalf("Error leyendo archivo YAML: %v", err)
	}

	yamlHandler, err := YAMLHandler(yamlFile, fallback)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", yamlHandler)
}

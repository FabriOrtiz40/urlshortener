package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	bolt "go.etcd.io/bbolt"
)

func main() {

	yamlPath := flag.String("yaml", "paths.yaml", "Ruta al archivo YAML con los paths")
	jsonPath := flag.String("json", "", "Ruta al archivo JSON con los paths")
	dbPath := flag.String("db", "", "Ruta al archivo BoltDB")

	flag.Parse()

	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hola! Esta página no existe.")
	})

	if *dbPath != "" {
		db, err := bolt.Open(*dbPath, 0600, nil)
		if err != nil {
			log.Fatalf("No se pudo abrir la base de datos: %v", err)
		}
		defer db.Close()

		dbHandler := DBHandler(db, fallback)

		fmt.Println("Servidor corriendo en http://localhost:8080 (DB mode)")
		http.ListenAndServe(":8080", dbHandler)
		return
	}

	if *jsonPath != "" {
		jsonData, err := os.ReadFile(*jsonPath)
		if err != nil {
			log.Fatalf("Error leyendo archivo JSON: %v", err)
		}

		jsonHandler, err := JSONHandler(jsonData, fallback)
		if err != nil {
			log.Fatalf("Error creando JSONHandler: %v", err)
		}

		fmt.Println("Servidor corriendo en http://localhost:8080 (JSON mode)")
		http.ListenAndServe(":8080", jsonHandler)
		return
	}

	yamlFile, err := os.ReadFile(*yamlPath)
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

package main

import (
	"log"

	bolt "go.etcd.io/bbolt"
)

func main() {
	db, err := bolt.Open("paths.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("paths"))
		if err != nil {
			return err
		}
		// Insertar rutas
		b.Put([]byte("/yt"), []byte("https://www.youtube.com"))
		b.Put([]byte("/r"), []byte("https://www.reddit.com"))
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Base de datos creada con éxito.")
}

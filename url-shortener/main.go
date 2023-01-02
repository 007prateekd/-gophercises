package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"url-shortener/handlers"

	"github.com/boltdb/bolt"
)

var (
	port       string
	yamlPath   string
	jsonPath   string
	bucketName = []byte("mapping")
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.StringVar(&port, "port", "8081", "Port on which the server runs")
	flag.StringVar(&yamlPath, "yaml", "", "Path of YAML file")
	flag.StringVar(&jsonPath, "json", "", "Path of JSON file")
}

func updateDb(db *bolt.DB, keys, vals []string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		for i := range keys {
			err = bucket.Put([]byte(keys[i]), []byte(vals[i]))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func createMap(db *bolt.DB, keys []string) (map[string]string, error) {
	pathMap := make(map[string]string)
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found", bucketName)
		}
		for i := range keys {
			val := string(bucket.Get([]byte(keys[i])))
			pathMap[keys[i]] = val
		}
		return nil
	})
	return pathMap, err
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})
	return mux
}

func main() {
	flag.Parse()
	db, err := bolt.Open("my.db", 0644, &bolt.Options{Timeout: 1 * time.Second})
	checkError(err)
	defer db.Close()
	// Store key-value pairs in DB
	keys := []string{"/github", "/youtube"}
	vals := []string{"https://github.com", "https://youtube.com/"}
	err = updateDb(db, keys, vals)
	checkError(err)
	// Retrieve data from DB
	pathMap, err := createMap(db, keys)
	checkError(err)
	// Handlers and ServeMuxes
	mux := defaultMux()
	mapHandler := handlers.MapHandler(pathMap, mux)
	handler := mapHandler
	if yamlPath != "" {
		data, err := ioutil.ReadFile(yamlPath)
		checkError(err)
		handler, err = handlers.YAMLHandler([]byte(data), mapHandler)
		checkError(err)
	} else if jsonPath != "" {
		data, err := ioutil.ReadFile(jsonPath)
		checkError(err)
		handler, err = handlers.JSONHandler([]byte(data), mapHandler)
		checkError(err)
	}
	fmt.Printf("Starting the server on :%s\n", port)
	http.ListenAndServe(":"+port, handler)
}

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/valyala/fasthttp"
)

type Filtro struct {
	C  []Campos `json:"C"`
	E  []Evals  `json:"E"`
	Id int32    `json:"Id"`
}
type Campos struct {
	T int      `json:"T"`
	N string   `json:"N"`
	V []string `json:"V"`
}
type Evals struct {
	T int    `json:"T"`
	N string `json:"N"`
}
type MyHandler struct {
	Dbs   *bolt.DB `json:"Dbs"`
	Count int64    `json:"Count"`
	Total int64    `json:"Total"`
}

func main() {

	total := 1

	filtro := Filtro{}
	filtro.C = []Campos{Campos{T: 1, N: "Procesador", V: []string{"X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71", "X15", "IntelC3", "Amd71"}}, Campos{T: 1, N: "Pantalla", V: []string{"4", "4.5", "5.5", "4.5", "5.5", "4.5", "5.5", "4.5", "5.5", "4.5", "5.5"}}, Campos{T: 1, N: "Memoria", V: []string{"2GB", "4GB", "8GB", "16GB", "32GB", "64GB", "128GB"}}, Campos{T: 1, N: "Marca", V: []string{"Samsung", "Motorola", "Nokia", "Samsung", "Motorola", "Nokia", "Samsung", "Motorola", "Nokia", "Samsung", "Motorola", "Nokia", "Samsung", "Motorola", "Nokia"}}}
	filtro.E = []Evals{Evals{T: 1, N: "Buena"}, Evals{T: 1, N: "Nelson"}, Evals{T: 1, N: "Hola"}, Evals{T: 1, N: "Mundo"}}

	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		fmt.Errorf("could not open db, %v", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("DB"))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		for i := 1; i <= total; i++ {
			_, err = root.CreateBucketIfNotExists([]byte(strconv.Itoa(i)))
			if err != nil {
				fmt.Errorf("could not create weight bucket: %v", err)
			} else {
				filtro.Id = int32(i)
				filtroBytes, err := json.Marshal(filtro)
				if err == nil {
					fmt.Println("BUCKET CREATED", i)
					err = tx.Bucket([]byte("DB")).Put([]byte("1000-"+strconv.Itoa(i)), filtroBytes)
					if err != nil {
						fmt.Errorf("could not set config: %v", err)
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Errorf("could not set up buckets, %v", err)
	} else {
		fmt.Println("DB Setup Done")
	}

	/*
		conf := Config{Height: 186.0, Birthday: time.Now()}
		err = setConfig(db, conf)
		if err != nil {
			log.Fatal(err)
		}
		err = addWeight(db, "85.0", time.Now())
		if err != nil {
			log.Fatal(err)
		}
		err = addEntry(db, 100, "apple", time.Now())
		if err != nil {
			log.Fatal(err)
		}
		err = addEntry(db, 100, "orange", time.Now().AddDate(0, 0, -2))
		if err != nil {
			log.Fatal(err)
		}
		err = db.View(func(tx *bolt.Tx) error {
			conf := tx.Bucket([]byte("DB")).Get([]byte("CONFIG"))
			fmt.Printf("Config: %s\n", conf)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("DB")).Bucket([]byte("WEIGHT"))
			b.ForEach(func(k, v []byte) error {
				fmt.Println(string(k), string(v))
				return nil
			})
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		err = db.View(func(tx *bolt.Tx) error {
			c := tx.Bucket([]byte("DB")).Bucket([]byte("ENTRIES")).Cursor()
			min := []byte(time.Now().AddDate(0, 0, -7).Format(time.RFC3339))
			max := []byte(time.Now().AddDate(0, 0, 0).Format(time.RFC3339))
			for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
				fmt.Println(string(k), string(v))
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	*/

	h := &MyHandler{Dbs: db, Total: int64(total), Count: 1}
	fasthttp.ListenAndServe(":80", h.HandleFastHTTP)
}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")

	switch string(ctx.Path()) {
	case "/get":

		fmt.Println("bytes")
		err := h.Dbs.View(func(tx *bolt.Tx) error {
			bytes := tx.Bucket([]byte("DB")).Get([]byte("1000-1"))
			fmt.Fprintf(ctx, string(bytes))
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}

	default:
		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}

}

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ServiceWeaver/weaver"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	if err := weaver.Run(context.Background(), run); err != nil {
		log.Fatal(err)
	}
}

type JsonReturn struct {
	Id      int    `json:"id"`
	Nama    string `json:"nama,omitempty"`
	Kondisi bool   `json:"kondisi"`
	Email   string `json:"email,omitempty"`
}

type app struct {
	weaver.Implements[weaver.Main]
	dengar weaver.Listener `weaver:"hello"`
}

func run(ctx context.Context, app *app) error {
	fmt.Printf("Listener alamat %s:", app.dengar)

	http.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hai2")
	})

	// Router
	r := mux.NewRouter()

	http.HandleFunc("/read", read)
	r.HandleFunc("/insert", insert).Methods("POST")
	r.HandleFunc("/update", update).Methods("UPDATE")
	r.HandleFunc("/delete/{id:[0-9]+}", delete).Methods("DELETE")
	http.Handle("/", r)

	return http.Serve(app.dengar, nil)
}

func insert(w http.ResponseWriter, r *http.Request) {
	// KONEKSI DB
	db, err := sql.Open("mysql", "root:admin@tcp(localhost:3306)/golang")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	// END KONEKSI DB
	c := context.Background()

	if r.Method != http.MethodPost {
		http.Error(w, "Method Tidak Diperbolehkan", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Gagal Memproses Form", http.StatusBadRequest)
		return
	}

	// GET VALUE FORM
	nama := r.FormValue("nama")
	tanggal := r.FormValue("tanggal")
	email := r.FormValue("email")
	kondisi := r.FormValue("kondisi")

	// PREPARE
	stmt, errPre := db.PrepareContext(c, "INSERT INTO test (nama, tanggal, email, kondisi) VALUES(?, ?, ?, ?)")
	if errPre != nil {
		http.Error(w, "Prepare SQL Gagal", http.StatusInternalServerError)
		return
	}

	// EXEC DB
	_, errExec := stmt.ExecContext(c, nama, tanggal, email, kondisi)
	if errExec != nil {
		http.Error(w, "SALAH SQL", http.StatusInternalServerError)
		return
	}

	defer stmt.Close()
	fmt.Fprint(w, "Berhasil Ditambah")
}

func read(w http.ResponseWriter, r *http.Request) {
	// KONEKSI DB
	db, err := sql.Open("mysql", "root:admin@tcp(localhost:3306)/golang")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	c := context.Background()

	// Get Table
	data, err := db.QueryContext(c, "SELECT id, nama, email, kondisi FROM test")
	if err != nil {
		panic(err.Error())
	}

	// Inisiasi Array
	var kembaliJSON []JsonReturn

	// Looping Data
	for data.Next() {
		var id int
		var kondisi bool
		var nama, email sql.NullString
		err := data.Scan(&id, &nama, &email, &kondisi)
		if err != nil {
			panic(err.Error())
		}

		kembaliJSON = append(kembaliJSON,
			JsonReturn{
				Id:      id,
				Nama:    nama.String,
				Kondisi: kondisi,
				Email:   email.String,
			})
	}

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	errEncode := encoder.Encode(kembaliJSON)
	if errEncode != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Tidak Diperbolehkan", http.StatusBadRequest)
		return
	}

	// KONEKSI DB
	db, err := sql.Open("mysql", "root:admin@tcp(localhost:3306)/golang")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Get Param
	// vars := mux.Vars(r)
	// id := vars["id"]

	// Context
	c := context.Background()

	// Parse Data
	if errParse := r.ParseForm(); errParse != nil {
		http.Error(w, "Kesalahan Data", http.StatusBadRequest)
		return
	}

	// Get Data Form
	id := r.FormValue("idEdit")
	nama := r.FormValue("namaEdit")
	tgl := r.FormValue("tglEdit")
	email := r.FormValue("emailEdit")
	timestamp := r.FormValue("timeEdit")
	kondisi := r.FormValue("kondisi")

	// Prepare Data
	stmt, errPre := db.PrepareContext(c, "UPDATE test SET nama=?, tanggal=?, email=?, timestamp=?, kondisi=? WHERE id=?")
	if errPre != nil {
		http.Error(w, "Kesalahan Prepare", http.StatusInternalServerError)
		return
	}

	// Exec DB
	_, errExec := stmt.ExecContext(c, nama, tgl, email, timestamp, kondisi, id)

	if errExec != nil {
		http.Error(w, "Kesalahan Exec DB", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Berhasil Update", id)
	defer stmt.Close()
}

func delete(w http.ResponseWriter, r *http.Request) {

	// Cek Method
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Tidak Boleh", http.StatusBadRequest)
		return
	}

	// KONEKSI DB
	db, err := sql.Open("mysql", "root:admin@tcp(localhost:3306)/golang")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Context
	c := context.Background()

	// Prepare Context
	stmt, errPre := db.PrepareContext(c, "DELETE FROM test WHERE id=?")
	if errPre != nil {
		http.Error(w, "Prepare Salah", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	// Exec DB
	_, errExec := stmt.ExecContext(c, id)
	if errExec != nil {
		http.Error(w, "Exec DB", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Berhasil Hapus", id)
}

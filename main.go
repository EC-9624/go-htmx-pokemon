package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

//go:embed views/*
var views embed.FS

var t = template.Must(template.ParseFS(views, "views/*"))
func main(){
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/pichu")
		if err != nil {
			http.Error(w, "Unable to grab the pokemon data", http.StatusInternalServerError)
			return
		}
		data := Pokemon{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			http.Error(w, "Unable to parse the Pokemon data", http.StatusInternalServerError)
			return
		}
	
		if err := t.ExecuteTemplate(w, "index.html", data); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	})

	router.HandleFunc("POST /poke",func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w,"unable to parse form", http.StatusInternalServerError)
			return
		}

		resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + strings.ToLower(r.FormValue("pokemon")))
		if err != nil {
			http.Error(w, "Unable to fetch new pokemon", http.StatusNotFound)
			fmt.Println(err)
			return
		}
		data := Pokemon{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			http.Error(w, "Unable to parse the Pokemon data", http.StatusUnprocessableEntity)
			return
		}
		
		if err := t.ExecuteTemplate(w, "response.html", data); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server Listening on port 3000...")
	http.ListenAndServe(":3000",router);
}

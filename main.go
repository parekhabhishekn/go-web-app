package main

import (
	"fmt"
	"net/http"
	"html/template"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"encoding/json"
)

type SearchResult struct {
	Title string;
	Author string;
	Year string;
	ISBN string;
	ID string;
}

func main() {

	templates := template.Must(template.ParseFiles("templates/home.html"));
	db,_ := sql.Open("mysql", "USER:PASSWORD@/DATABASE");

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		templates.ExecuteTemplate(w, "home.html", nil);
	});

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request){
		term := r.FormValue("search-term");
		rows, err := db.Query("SELECT title,author,year,isbn FROM books WHERE title like ? ","%"+term+"%");
		if err != nil {
			panic(err.Error());
		}
		var sr []SearchResult;
		for rows.Next() {
			result := SearchResult{};
			err = rows.Scan(&result.Title, &result.Author, &result.Year, &result.ISBN);
			if err != nil {
				panic(err.Error());
			}
			sr = append(sr,result);
		}
		rows.Close();
		encoder := json.NewEncoder(w);
		encoder.Encode(sr);

	});

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request){
		title := r.FormValue("title");
		author := r.FormValue("author");
		year := r.FormValue("year");
		isbn  := r.FormValue("isbn");
		var result = SearchResult{title,author,year,isbn,""};
		_,err := db.Exec("INSERT INTO books (title, author, year, isbn) values ( ?, ?, ?, ?)", 
			result.Title, result.Author, result.Year, result.ISBN);
		if err != nil {
			return;
		}
		rows, err := db.Query("SELECT title,author,year,isbn FROM books");
		if err != nil {
			panic(err.Error());
		}
		var sr []SearchResult;
		for rows.Next() {
			result := SearchResult{};
			err = rows.Scan(&result.Title, &result.Author, &result.Year, &result.ISBN);
			if err != nil {
				panic(err.Error());
			}
			sr = append(sr,result);
		}
		rows.Close();
		encoder := json.NewEncoder(w);
		encoder.Encode(sr);

	});

	
	fmt.Println(http.ListenAndServe(":8080",nil));
}

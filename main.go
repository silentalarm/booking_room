package main
import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request){
		http.ServeFile(w, r, "templates/about.html")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		http.ServeFile(w, r, "templates/index.html")
		//r.ServeFiles("/public/*filepath", http.Dir("public"))

	})
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}
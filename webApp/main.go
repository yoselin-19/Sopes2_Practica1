package main

//Importaciones
import (
	"net/http"
	"fmt"

	//Para hacer el api rest
	"github.com/gorilla/mux"
)

//=======================================================================

//Funcion Principal
func main() {
	router := mux.NewRouter().StrictSlash(true)
	// Para los archivos staticos (css,js)
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
   
	//Rutas para cliente -Si ya tiene en la ruta .html ignora si send a un procedimiento y redirige a la pagina.html-
	router.HandleFunc("/Principal.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w,r, "./public/Principal.html")
	})

	router.HandleFunc("/RAM.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w,r, "./public/RAM.html")
	})

	router.HandleFunc("/Arbol.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w,r, "./public/Arbol.html")
	})

	//Servidor levantado
	fmt.Println("Servidor levantado en el puerto: 3000")
	http.ListenAndServe(":3000", router)
}

//=======================================================================
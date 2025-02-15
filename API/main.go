package main

//Importaciones
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"

	//Para lectura de los archivos
	"strings"

	"./librerias"

	//Para usar json
	"encoding/json"

	"github.com/tidwall/gjson"

	//Para conversiones

	//Para hacer el api rest
	"github.com/gorilla/mux"
)

//=======================================================================

//Funcion Principal
func main() {
	router := mux.NewRouter().StrictSlash(true)
	// Para los archivos staticos (css,js)
	router.PathPrefix("../webApp/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	//Rutas de API-REST
	router.HandleFunc("/PROCESS", lista_procesos)
	router.HandleFunc("/RAM", memoria_proceso)
	router.HandleFunc("/kill/{id}", kill_proceso)
	router.HandleFunc("/Arbol", arbol_procesos)

	//Rutas para cliente -Si ya tiene en la ruta .html ignora si send a un procedimiento y redirige a la pagina.html-
	router.HandleFunc("/Principal.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/Principal.html")
	})

	router.HandleFunc("/RAM.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/RAM.html")
	})

	router.HandleFunc("/Arbol.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/Arbol.html")
	})

	//Servidor levantado
	fmt.Println("Servidor levantado en el puerto: 3000")
	http.ListenAndServe(":3000", router)
}

func memoria_proceso(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("/proc/mem_grupo14")
	if err != nil {
		panic(err)
	}

	memoria_total := gjson.Get(string(data), "memoria_total_mb")
	memoria_utilizada := gjson.Get(string(data), "memoria_utilizada_porcentaje")
	memoria_consumida := gjson.Get(string(data), "memoria_consumida_mb")
	//Conversiones y calculos
	MemTotal_, _ := strconv.Atoi(memoria_total.String())

	MemConsumida, _ := strconv.Atoi(memoria_consumida.String())

	numero_porcentaje, _ := strconv.Atoi(memoria_utilizada.String())
	PorcentajeConsumo := float32(numero_porcentaje)

	info_ram := RAM{
		Total_Ram_Servidor:     MemTotal_,
		Total_Ram_Consumida:    MemConsumida,
		Porcentaje_Consumo_Ram: PorcentajeConsumo,
	}

	JSON_Data, _ := json.Marshal(info_ram)
	w.Write(JSON_Data)

}

func lista_procesos(w http.ResponseWriter, r *http.Request) {
	var arr_process []PROCESS
	data, err := ioutil.ReadFile("/proc/cpu_grupo14")
	if err != nil {
		panic(err)
	}
	procesos := gjson.Get(string(data), "cpu")
	arr_process = readProcesos(procesos, "0", arr_process)

	//Agregando informacion general
	info_general := Info_general{
		Procesos_en_ejecucion: librerias.NumeroRun,
		Procesos_suspendidos:  librerias.NumeroSleep,
		Procesos_detenidos:    librerias.NumeroStop,
		Procesos_zombie:       librerias.NumeroZombie,
		Total_procesos:        len(arr_process),
		List_Procesos:         arr_process,
	}

	JSON_Data, _ := json.Marshal(info_general)
	w.Write(JSON_Data)
}

func readProcesos(procesos gjson.Result, padre string, arr_process []PROCESS) []PROCESS {

	for _, proceso := range procesos.Array() {

		Pid_ := gjson.Get(proceso.String(), "pid")

		Nombre_ := gjson.Get(proceso.String(), "nombre")

		Estado_ := gjson.Get(proceso.String(), "estado")

		Usuario_ := gjson.Get(proceso.String(), "usuario")

		hijos := gjson.Get(proceso.String(), "hijos")

		if strings.Contains(hijos.String(), "{") {
			arr_process = readProcesos(hijos, Pid_.String(), arr_process)
		}

		info_process := PROCESS{
			PID:           Pid_.String(),
			Nombre:        Nombre_.String(),
			Usuario:       librerias.GetNombreUsuario(Usuario_.String()),
			Estado:        librerias.GetStatus(Estado_.String()),
			PorcentajeRAM: librerias.GetPorcentajeRAM(Pid_.String()),
			Proceso_padre: padre,
		}
		arr_process = append(arr_process, info_process)

	}
	return arr_process
}

func armarProcesos(procesos gjson.Result, padre string, arr_process []librerias.Arbol, raiz librerias.Arbol) []librerias.Arbol {

	for _, proceso := range procesos.Array() {

		Pid_ := gjson.Get(proceso.String(), "pid")
		PidNum, _ := strconv.Atoi(Pid_.String())
		Nombre_ := gjson.Get(proceso.String(), "nombre")

		hijos := gjson.Get(proceso.String(), "hijos")

		if strings.Contains(hijos.String(), "{") {
			arr_process = armarProcesos(hijos, Pid_.String(), arr_process, raiz)
		}

		numero_padre, _ := strconv.Atoi(padre)

		raiz = librerias.Arbol{
			Pid:    PidNum,
			Nombre: Nombre_.String(),
			Ppid:   numero_padre,
			Hijos:  nil,
		}
		arr_process = append(arr_process, raiz)

	}
	return arr_process
}

func kill_proceso(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["id"]
	librerias.MatarProceso(key)
	http.Redirect(w, r, "/public/Principal.html", http.StatusFound)
}

func arbol_procesos(w http.ResponseWriter, r *http.Request) {

	//Obteniendo lista de directorios
	data, err := ioutil.ReadFile("/proc/cpu_grupo14")
	if err != nil {
		panic(err)
	}
	//Variables para crear el arreglo de Arbol de procesos
	var raiz librerias.Arbol
	var arreglo []librerias.Arbol

	//Recorriendo cada directorio
	procesos := gjson.Get(string(data), "cpu")

	arreglo = armarProcesos(procesos, "0", arreglo, raiz)

	// Sort by age, keeping original order or equal elements.
	sort.SliceStable(arreglo, func(i, j int) bool {
		return arreglo[i].Ppid < arreglo[j].Ppid
	})

	//Construir texto de arbol
	var nuevoB librerias.Arbol
	for _, item := range arreglo {
		librerias.Insertar(&nuevoB, item)
	}

	TextoArbol := librerias.GetTextoArbol(nuevoB)
	info_tree := Tree{Arbol: TextoArbol}

	JSON_Data, _ := json.Marshal(info_tree)
	w.Write(JSON_Data)

}

//=======================================================================

//Estructuras a utilizar
type RAM struct {
	Total_Ram_Servidor     int
	Total_Ram_Consumida    int
	Porcentaje_Consumo_Ram float32
}

type PROCESS struct {
	PID           string `json:"pid"`
	Nombre        string `json:"nombre"`
	Usuario       string
	Estado        string `json:"estado"`
	PorcentajeRAM string
	Proceso_padre string
	hijos         []PROCESS `json:"hijos"`
}

type Info_general struct {
	Procesos_en_ejecucion int
	Procesos_suspendidos  int
	Procesos_detenidos    int
	Procesos_zombie       int
	Total_procesos        int
	List_Procesos         []PROCESS
}

type Tree struct {
	Arbol string
}

type CPU struct {
	procesos []PROCESS
}

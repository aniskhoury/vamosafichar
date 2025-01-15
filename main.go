package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var plantilla *template.Template

const servidor = "127.0.0.1"
const puerto = "3306"
const usuarioBBDD = "vamosafichar"
const contrasenaBBDD = "vamosafichar"
const basededatos = "vamosafichar"

type Usuario struct {
	idusuario     int
	identificador string
	contrasena    string
	fechaalta     string
}

var autentificacion = "" + usuarioBBDD + ":" + contrasenaBBDD + "@tcp(" + servidor + ":" + puerto + ")/" + basededatos

func init() {
	plantilla = template.Must(template.ParseGlob("plantillas/*.html"))
}
func index(w http.ResponseWriter, r *http.Request) {
	plantilla.ExecuteTemplate(w, "index.html", nil)
}
func autentificacionFunc(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	varmap := map[string]interface{}{
		"identificador":    r.FormValue("identificador"),
		"estaIdentificado": false,
	}

	var consultaSQL = "SELECT * FROM usuarios where identificador = ? and contrasena = ?"
	var estaAutentificado bool = false
	db, err := sql.Open("mysql", autentificacion)
	if err != nil {
		panic(err.Error())
	}
	filas, err := db.Query(consultaSQL, r.FormValue("identificador"), r.FormValue("contrasena"))
	if err != nil {
		panic(err.Error())
	}
	var usuarioDB Usuario
	for filas.Next() {

		err = filas.Scan(&usuarioDB.idusuario, &usuarioDB.identificador, &usuarioDB.contrasena, &usuarioDB.fechaalta)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("ID Usuario:", usuarioDB.idusuario)
		estaAutentificado = true
	}
	filas.Close()
	if estaAutentificado {
		varmap["estaIdentificado"] = true
	}

	plantilla.ExecuteTemplate(w, "autentificacion.html", varmap)

}
func insertarUsuario(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	varmap := map[string]interface{}{
		"identificador": r.FormValue("identificador"),
		"error":         false,
		"mensajeError":  "",
	}
	var consulta = "INSERT INTO `usuarios` (`identificador`,`contrasena`) VALUES (?,?)"
	db, err := sql.Open("mysql", autentificacion)
	if err != nil {
		varmap["error"] = true
		varmap["mensajeError"] = err.Error()
	}
	insertar, err := db.ExecContext(context.Background(), consulta, r.FormValue("identificador"), r.FormValue("contrasena"))
	if err != nil {
		varmap["error"] = true
		varmap["mensajeError"] = err.Error()
	} else {
		id, err := insertar.LastInsertId()
		if err != nil {
			varmap["error"] = true
			varmap["mensajeError"] = err.Error()
		} else {
			fmt.Printf("id %d usuario insertado correctamente", id)
		}
	}

	plantilla.ExecuteTemplate(w, "usuarioInsertado.html", varmap)
}
func main() {

	servidorDeFicheros := http.FileServer(http.Dir("recursos"))
	http.Handle("/recursos/", http.StripPrefix("/recursos", servidorDeFicheros))
	http.HandleFunc("/", index)
	http.HandleFunc("/autentificacion", autentificacionFunc)
	http.HandleFunc("/insertarUsuario", insertarUsuario)
	http.ListenAndServe(":9999", nil)

}

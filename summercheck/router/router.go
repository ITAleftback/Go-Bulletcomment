package router

import (
	"net/http"
	"summercheck/controllers"
)

func Registerrouter(){

	http.HandleFunc("/addroom",controllers.Addroom)
}

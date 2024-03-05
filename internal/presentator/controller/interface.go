package controller

import "net/http"

type ControllerType interface {
	Health(http.ResponseWriter, *http.Request)
	Leader(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	Set(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

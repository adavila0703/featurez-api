package settings

import (
	"net/http"
)

func Routes(route string, mux *http.ServeMux, subroute string) {
	route = route + subroute

	mux.Handle(route+"/UpdateSettings", UpdateSettingsHandler)
	mux.Handle(route+"/GetSettings", GetSettingsHandler)
}

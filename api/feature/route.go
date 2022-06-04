package feature

import (
	"net/http"
)

func Routes(route string, mux *http.ServeMux, subroute string) {
	route = route + "/feature"

	mux.Handle(route+"/CreateFeature", CreateFeatureHandler)
	mux.Handle(route+"/GetFeatureList", GetFeatureListHandler)
	mux.Handle(route+"/DeleteFeature", DeleteFeatureHandler)
	mux.Handle(route+"/UpdateFeature", UpdateFeatureHandler)
}

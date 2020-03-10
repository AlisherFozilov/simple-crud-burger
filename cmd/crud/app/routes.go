package app

import "net/http"

func (receiver *server) InitRoutes() {
	mux := receiver.router.(*exactMux)

	mux.GET("/", receiver.handleBurgersList())
	mux.POST("/burgers/save", receiver.handleBurgersSave())
	mux.POST("/burgers/remove", receiver.handleBurgersRemove())
	mux.GET("/favicon.ico", receiver.handleFavicon())
	mux.notFoundHandler = http.HandlerFunc(receiver.handleNotFound())
}

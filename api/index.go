package handler

import (
	"gin_jwt/app"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	engine := app.SetupRouter()
	engine.ServeHTTP(w, r)
}

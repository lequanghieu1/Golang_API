package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/conglt10/web-golang/middlewares"
	"github.com/conglt10/web-golang/routes"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	router := httprouter.New()

	router.POST("/login", routes.Login)
	router.POST("/register", routes.Register)

	//__________________________Manager____________________________
	router.GET("/managers", middlewares.CheckJwt(routes.GetAllUsers))
	router.POST("/managers", middlewares.CheckJwt(routes.CreateManager))
	router.PUT("/managers/:id", middlewares.CheckJwt(routes.EditManager))
	router.DELETE("/managers/:id", middlewares.CheckJwt(routes.DeleteManager))
	//__________________________Event-Code____________________________
	router.GET("/event-code", middlewares.CheckJwt(routes.GetAllEvents))
	router.POST("/event-code", middlewares.CheckJwt(routes.CreateEventCode))
	router.PUT("/event-code/:id", middlewares.CheckJwt(routes.EditEventCode))
	router.DELETE("/event-code/:id", middlewares.CheckJwt(routes.DeleteEventCode))
	//__________________________Model-Device____________________________
	router.GET("/model-device", middlewares.CheckJwt(routes.GetAllModels))
	router.POST("/model-device", middlewares.CheckJwt(routes.CreateModelDevices))
	router.PUT("/model-device/:id", middlewares.CheckJwt(routes.EditModelDevice))
	router.DELETE("/model-device/:id", middlewares.CheckJwt(routes.DeleteModelDevice))
	//__________________________Page-Schema____________________________
	router.GET("/page-schema", middlewares.CheckJwt(routes.GetAllPages))
	router.POST("/page-schema", middlewares.CheckJwt(routes.CreatePageSchema))
	router.PUT("/page-schema/:id", middlewares.CheckJwt(routes.EditPageSchema))
	router.DELETE("/page-schema/:id", middlewares.CheckJwt(routes.DeletePageSchema))
	//__________________________Roles-Access____________________________
	router.GET("/roles-access", middlewares.CheckJwt(routes.GetAllRoles))
	router.POST("/roles-access", middlewares.CheckJwt(routes.CreateRolesAccess))
	router.PUT("/roles-access/:id", middlewares.CheckJwt(routes.EditRolesAccess))
	router.DELETE("/roles-access/:id", middlewares.CheckJwt(routes.DeleteRolesAccess))

	fmt.Println("Listening to port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

package main

import (
	"log"

	"comprarmas.com.mx/internal/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	log.Println("Chargeing ...")
	app := fiber.New()
	app.Server().MaxConnsPerIP = 5

	app.Use(cors.New())

	app.Static("/", "./public")
	app.Post("/login/", controller.PostLogin)
	app.Post("/register/", controller.PostRegisterUsr)
	app.Post("/control/lp/", controller.GetListProducts)
	app.Get("/cancel/:id", controller.PostCancelProduct)
	app.Post("/regproduct", controller.PostRegProduct)
	app.Post("/product/search/", controller.PostProductSearch)
	app.Get("/product/admin/", controller.PostProductAdmin)

	log.Fatal(app.Listen(":3000"))
}

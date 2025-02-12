package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors" // Importar el middleware de CORS
	analyzer "backend/analyzer" // Importar el paquete analyzer
)

func main() {
	// Crear una nueva instancia de Fiber
	app := fiber.New()

	// Configurar el middleware de CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000", // Permitir solicitudes desde el frontend
		AllowHeaders: "Origin, Content-Type, Accept", // Headers permitidos
		AllowMethods: "GET, POST, PUT, DELETE", // Métodos HTTP permitidos
	}))

	// Ruta de prueba
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("¡Hola, mundo!")
	})

	// Ruta para ejecutar comandos
	app.Post("/execute", func(c *fiber.Ctx) error {
		// Obtener el cuerpo de la solicitud
		var requestBody struct {
			Command string `json:"command"`
		}

		// Parsear el cuerpo de la solicitud
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No se pudo procesar la solicitud",
			})
		}

		// Mostrar el comando recibido
		fmt.Println("Comando recibido:", requestBody.Command)

		// analizar la entrada
		cmd, err := analyzer.Analyzer(requestBody.Command)
		if err != nil {
			// Si hay un error al analizar el comando, imprime el error y continúa con el siguiente comando
			fmt.Println("Error:", err)
		}
		// Comentado: Aquí podrías imprimir el comando analizado
		fmt.Printf("Parsed Command: %+v\n", cmd)

		// Responder con un mensaje de éxito
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("Comando ejecutado: %s", requestBody.Command),
		})
	})

	// Iniciar el servidor en el puerto 8000
	log.Fatal(app.Listen(":8000"))
}
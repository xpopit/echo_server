package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	urlBigIP := os.Getenv("URL_BIG_IP")

	// Create a new Fiber instance
	app := fiber.New()

	// Add middleware to log requests
	app.Use(logger.New())

	app.Get("/del", func(c *fiber.Ctx) error {
		queryParams := c.Request().URI().QueryString()

		// ถ้า URL_BIG_IP ไม่มีการตั้งค่า ให้ใช้ IP:PORT ของต้นทางที่เรียกมา
		if urlBigIP == "" {
			urlBigIP = c.IP() + ":" + c.Port()
		}

		targetUrl := fmt.Sprintf("http://%s/del?%s", urlBigIP, queryParams)

		// เพิ่มการบันทึก targetUrl
		log.Printf("Calling target URL: %s", targetUrl)

		// Make the GET request to the target URL
		resp, err := http.Get(targetUrl)
		if err != nil {
			log.Printf("Error fetching from target URL: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error fetching from target URL: %v", err))
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error reading response body: %v", err))
		}

		// Set the status code and return the response body
		return c.Status(resp.StatusCode).Send(body)
	})

	log.Printf("Server running at http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}

package main

import (
	"employee_crud/api"
	"employee_crud/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	db, err := database.SetUpDatabase()
	if err != nil {
		log.Fatal("Failed to connect db", err)
	}

	server := api.NewServer(db)

	router.POST("/employees", server.CreateEmployee)
	router.GET("/employees", server.GetEmployees)
	router.GET("/employees/:id", server.GetEmployeeByID)
	router.PATCH("/employees/:id", server.UpdateEmployee)
	router.DELETE("/employees/:id", server.DeleteEmployeeByID)

	router.Run(":8080")

}

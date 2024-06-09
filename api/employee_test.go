package api

import (
	"bytes"
	"employee_crud/database"
	"employee_crud/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func TestCRUDEmployee(t *testing.T) {
	router := gin.Default()

	db, err := database.SetUpDatabase()
	if err != nil {
		log.Fatal("Failed to connect db", err)
	}

	server := NewServer(db)

	router.POST("/employees", server.CreateEmployee)
	router.GET("/employees/:id", server.GetEmployeeByID)
	router.PATCH("/employees/:id", server.UpdateEmployee)
	router.DELETE("/employees/:id", server.UpdateEmployee)

	// Test create a new employee
	w := performRequest(router, "POST", "/employees", `{"name": "Rakesh", "position": "Developer", "salary": 60000.00}`)

	var response model.Employee

	err = json.Unmarshal(w.Body.Bytes(), &response)

	require.NoError(t, err)
	assert.Equal(t, "Rakesh", response.Name)
	assert.Equal(t, "Developer", response.Position)
	assert.Equal(t, 60000.00, response.Salary)

	// Test update employee
	url := fmt.Sprintf("/employees/%d", response.ID)

	performRequest(router, "PATCH", url, `{"position": "Architect"}`)

	// Test Get employee by id
	w = performRequest(router, "GET", url, ``)

	var response2 model.Employee

	err = json.Unmarshal(w.Body.Bytes(), &response2)

	require.NoError(t, err)
	assert.Equal(t, "Rakesh", response2.Name)
	assert.Equal(t, "Architect", response2.Position)
	assert.Equal(t, 60000.00, response2.Salary)

	// Delete employee
	performRequest(router, "DELETE", url, ``)
}

func performRequest(r http.Handler, method, path string, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

package api

import (
	"database/sql"
	"employee_crud/model"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) CreateEmployee(c *gin.Context) {
	var employee model.Employee

	if err := c.ShouldBindBodyWithJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO employee (name, position,salary) VALUES ($1, $2, $3) RETURNING id"
	err := server.store.QueryRow(query, employee.Name, employee.Position, employee.Salary).Scan(&employee.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, employee)
}

func (server *Server) GetEmployees(c *gin.Context) {
	defaultPage := 1
	defaultPageSize := 10

	page, err := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(defaultPage)))
	if err != nil || page < 0 {
		page = defaultPage
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(defaultPageSize)))
	if err != nil || page < 0 {
		pageSize = defaultPageSize
	}
	offset := (page - 1) * defaultPageSize

	query := "SELECT id, name,position, salary FROM employee  LIMIT $1 OFFSET $2"
	rows, err := server.store.Query(query, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()
	var employees []model.Employee

	for rows.Next() {

		var employee model.Employee

		err = rows.Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Salary)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		employees = append(employees, employee)
	}

	c.JSON(http.StatusOK, gin.H{"employees": employees, "page": page, "pageSize": pageSize})

}

func (server *Server) GetEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return

	}
	var employee model.Employee

	query := "SELECT id, name,position, salary FROM employee WHERE id=$1"
	err = server.store.QueryRow(query, id).Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Salary)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No employees found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employee)

}

type EmployeeUpdate struct {
	Name     *string  `json:"name"`
	Position *string  `json:"position"`
	Salary   *float64 `json:"salary"`
}

func (server *Server) UpdateEmployee(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return

	}

	var employeeUpdate EmployeeUpdate

	if err := c.ShouldBindBodyWithJSON(&employeeUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v := reflect.ValueOf(employeeUpdate)
	typeOfUpdate := v.Type()

	query := "UPDATE employee SET "
	params := []interface{}{}
	counter := 1

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			if counter > 1 {
				query += ", "
			}
			query += typeOfUpdate.Field(i).Tag.Get("json") + " =$" + strconv.Itoa(counter)
			params = append(params, field.Elem().Interface())
			counter++
		}
	}

	query += " WHERE id = $" + strconv.Itoa(counter)
	params = append(params, id)

	fmt.Println(query, params)
	_, err = server.store.Exec(query, params...)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Employee updated successfully"})

}

func (server *Server) DeleteEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return

	}

	query := "DELETE FROM employee WHERE id=$1"
	_, err = server.store.Exec(query, id)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No employees found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Employee deleted")

}

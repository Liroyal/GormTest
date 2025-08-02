package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/yourname/employee-api/models"
	"github.com/yourname/employee-api/config"
	"github.com/yourname/employee-api/middleware"
	"github.com/yourname/employee-api/utils"
)

// CreateEmployeeHandler handles the creation of a new employee
func CreateEmployeeHandler(c *gin.Context) {
	logger := utils.GetLogger()
	requestID := middleware.GetRequestID(c)
	
	// Log request start
	logger.WithFields(logrus.Fields{
		"request_id": requestID,
		"operation":  "create_employee",
	}).Info("Processing create employee request")

	var employee models.Employee

	// Bind JSON to the employee struct and validate
	if err := c.ShouldBindJSON(&employee); err != nil {
		// Log validation error
		utils.LogValidationError(c, "employee_data", employee, err, logrus.Fields{
			"operation": "create_employee",
		})
		
		// Return standardized error response
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Invalid request format",
		})
		return
	}

	// Validate required fields
	if employee.FirstName == "" {
		err := errors.New("first_name is required")
		utils.LogValidationError(c, "first_name", employee.FirstName, err)
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "first_name is required",
		})
		return
	}

	if employee.LastName == "" {
		err := errors.New("last_name is required")
		utils.LogValidationError(c, "last_name", employee.LastName, err)
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "last_name is required",
		})
		return
	}

	// Get database connection
	db := config.GetDB()

	// Insert into database
	if err := db.Create(&employee).Error; err != nil {
		// Log database error with context
		utils.LogDBError(c, "create_employee", err, logrus.Fields{
			"employee_first_name": employee.FirstName,
			"employee_last_name":  employee.LastName,
		})
		
		// Return standardized error response
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Error: "Failed to create employee",
		})
		return
	}

	// Log successful creation
	logger.WithFields(logrus.Fields{
		"request_id":          requestID,
		"operation":           "create_employee",
		"employee_id":         employee.ID,
		"employee_first_name": employee.FirstName,
		"employee_last_name":  employee.LastName,
	}).Info("Employee created successfully")

	// Return created employee
	c.JSON(http.StatusCreated, employee)
}

// GetEmployeeHandler handles retrieving an employee by ID
func GetEmployeeHandler(c *gin.Context) {
	logger := utils.GetLogger()
	requestID := middleware.GetRequestID(c)
	
	// Get employee ID from URL parameter
	employeeID := c.Param("id")
	if employeeID == "" {
		err := errors.New("employee ID is required")
		utils.LogValidationError(c, "employee_id", employeeID, err)
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Employee ID is required",
		})
		return
	}

	// Log request start
	logger.WithFields(logrus.Fields{
		"request_id":  requestID,
		"operation":   "get_employee",
		"employee_id": employeeID,
	}).Info("Processing get employee request")

	var employee models.Employee
	db := config.GetDB()

	// Find employee by ID
	if err := db.First(&employee, employeeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Log warning for not found
			utils.LogDBError(c, "get_employee", err, logrus.Fields{
				"employee_id": employeeID,
			})
			c.JSON(http.StatusNotFound, middleware.ErrorResponse{
				Error: "Employee not found",
			})
		} else {
			// Log database error
			utils.LogDBError(c, "get_employee", err, logrus.Fields{
				"employee_id": employeeID,
			})
			c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
				Error: "Failed to retrieve employee",
			})
		}
		return
	}

	// Log successful retrieval
	logger.WithFields(logrus.Fields{
		"request_id":          requestID,
		"operation":           "get_employee",
		"employee_id":         employee.ID,
		"employee_first_name": employee.FirstName,
		"employee_last_name":  employee.LastName,
	}).Info("Employee retrieved successfully")

	c.JSON(http.StatusOK, employee)
}

// UpdateEmployeeHandler handles updating an employee
func UpdateEmployeeHandler(c *gin.Context) {
	logger := utils.GetLogger()
	requestID := middleware.GetRequestID(c)
	
	// Get employee ID from URL parameter
	employeeID := c.Param("id")
	if employeeID == "" {
		err := errors.New("employee ID is required")
		utils.LogValidationError(c, "employee_id", employeeID, err)
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Employee ID is required",
		})
		return
	}

	// Log request start
	logger.WithFields(logrus.Fields{
		"request_id":  requestID,
		"operation":   "update_employee",
		"employee_id": employeeID,
	}).Info("Processing update employee request")

	var updateData models.Employee

	// Bind JSON to the employee struct
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.LogValidationError(c, "employee_data", updateData, err, logrus.Fields{
			"operation":   "update_employee",
			"employee_id": employeeID,
		})
		c.JSON(http.StatusBadRequest, middleware.ErrorResponse{
			Error: "Invalid request format",
		})
		return
	}

	db := config.GetDB()

	// Check if employee exists first
	var existingEmployee models.Employee
	if err := db.First(&existingEmployee, employeeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.LogDBError(c, "update_employee", err, logrus.Fields{
				"employee_id": employeeID,
			})
			c.JSON(http.StatusNotFound, middleware.ErrorResponse{
				Error: "Employee not found",
			})
		} else {
			utils.LogDBError(c, "update_employee", err, logrus.Fields{
				"employee_id": employeeID,
			})
			c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
				Error: "Failed to retrieve employee",
			})
		}
		return
	}

	// Update employee
	if err := db.Model(&existingEmployee).Updates(updateData).Error; err != nil {
		utils.LogDBError(c, "update_employee", err, logrus.Fields{
			"employee_id":         employeeID,
			"employee_first_name": updateData.FirstName,
			"employee_last_name":  updateData.LastName,
		})
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Error: "Failed to update employee",
		})
		return
	}

	// Fetch updated employee
	var updatedEmployee models.Employee
	if err := db.First(&updatedEmployee, employeeID).Error; err != nil {
		utils.LogDBError(c, "update_employee_fetch", err, logrus.Fields{
			"employee_id": employeeID,
		})
		c.JSON(http.StatusInternalServerError, middleware.ErrorResponse{
			Error: "Failed to retrieve updated employee",
		})
		return
	}

	// Log successful update
	logger.WithFields(logrus.Fields{
		"request_id":          requestID,
		"operation":           "update_employee",
		"employee_id":         updatedEmployee.ID,
		"employee_first_name": updatedEmployee.FirstName,
		"employee_last_name":  updatedEmployee.LastName,
	}).Info("Employee updated successfully")

	c.JSON(http.StatusOK, updatedEmployee)
}

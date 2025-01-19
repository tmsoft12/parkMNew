package carcontrol

import (
	"fmt"
	"math"
	"os"
	"park/database"
	modelscar "park/models/modelsCar"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

const statusInside = "Inside"
const timeFormat = "2006-01-02 15:04:05"
const statusExited = "Exited"
const defaultImageURL = "example.com"

// CreateCar godoc
// @Summary Create a new car entry
// @Description Registers a new car entering the parking lot
// @Tags cars
// @Accept  json
// @Produce  json
// @Param parkno query string true "Parking spot number"
// @Param car body modelscar.Car_Model true "Car details"
// @Success 201 {object} map[string]interface{} "Created car details"
// @Failure 400 {object} ErrorResponse "Invalid request or car already inside"
// @Failure 500 {object} ErrorResponse "Database error"
// @Router /createcar [post]
func CreateCar(c *fiber.Ctx) error {
	var car modelscar.Car_Model

	now := time.Now().Format(timeFormat)
	parkno := c.Query("parkno")
	car.Start_time = now
	car.Status = statusInside
	car.ParkNo = parkno
	if err := c.BodyParser(&car); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}
	if err := database.DB.Order("id desc").First(&car, "car_number = ? AND status = ?", car.Car_number, statusInside).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Car is already inside the parking lot",
		})
	}
	if err := database.DB.Create(&car).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Database error",
			"error":   err.Error(),
		})
	}
	return c.Status(201).JSON(fiber.Map{
		"message": "Car created successfully",
		"car":     car,
	})
}

type GetCarsResponse struct {
	Cars       []modelscar.Car_Model `json:"cars"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalCount int64                 `json:"total_count"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
type UpdateCarResponse struct {
	Message string              `json:"message"`
	Car     modelscar.Car_Model `json:"car"`
}

// GetCars godoc
// @Summary Get list of cars
// @Description Get list of cars with pagination
// @Tags cars
// @Accept  json
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(5)
// @Param parkno query string false "Parking spot number"
// @Success 200 {object} GetCarsResponse
// @Failure 400 {object} ErrorResponse
// @Router /getallcars [get]
func GetCars(c *fiber.Ctx) error {
	var cars []modelscar.Car_Model
	var totalCount int64
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "5")
	parkno := c.Query("parkno")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid page number",
		})
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid limit number",
		})
	}

	query := database.DB.Model(&modelscar.Car_Model{})
	if parkno != "" {
		query = query.Where("park_no = ?", parkno)
	}

	query.Count(&totalCount)
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1

	offset := (page - 1) * limit
	query.Order("id desc").Limit(limit).Offset(offset).Find(&cars)

	if len(cars) == 0 {
		cars = []modelscar.Car_Model{}
	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	for i := range cars {
		cars[i].Image_Url = fmt.Sprintf("http://%s:%s/plate/%s", ip, port, cars[i].Image_Url)
	}
	return c.Status(200).JSON(fiber.Map{
		"cars":       cars,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
		"hasNext":    hasNext,
		"hasPrev":    hasPrev,
	})
}

// GetCar godoc
// @Summary Get a car by ID
// @Description Get a car by ID
// @Tags cars
// @Accept  json
// @Produce  json
// @Param id path int true "Car ID"
// @Success 200 {object} modelscar.Car_Model
// @Failure 404 {object} ErrorResponse
// @Router /getcar/{id} [get]
func GetCar(c *fiber.Ctx) error {
	id := c.Params("id")
	var car modelscar.Car_Model
	database.DB.Where("id = ?", id).First(&car)
	if car.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Car not found",
		})
	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	car.Image_Url = fmt.Sprintf("http://%s:%s/plate/%s", ip, port, car.Image_Url)

	c.Status(200)
	return c.JSON(car)
}

// UpdateCar godoc
// @Summary Update a car by plate number
// @Description Updates a car's status and calculates payment and duration based on start and end times.
// @Tags cars
// @Accept  json
// @Produce  json
// @Param plate path string true "Car plate number"
// @Param car body modelscar.Car_Model true "Car details to update"
// @Success 200 {object} map[string]interface{} "Updated car details"
// @Failure 400 {object} ErrorResponse "Car already exited or invalid request"
// @Failure 404 {object} ErrorResponse "Car not found"
// @Failure 500 {object} ErrorResponse "Error parsing time"
// @Router /updatecar/{plate} [put]
func UpdateCar(c *fiber.Ctx) error {
	plate := c.Params("plate")
	var car modelscar.Car_Model
	if err := database.DB.Where("car_number = ?", plate).Order("id desc").First(&car).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Car not found", "error": err.Error()})
	}

	var updatedCar modelscar.Car_Model
	if err := c.BodyParser(&updatedCar); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request", "error": err.Error()})
	}

	if car.Status == statusExited {
		return c.Status(400).JSON(fiber.Map{"message": "Car already exited"})
	}

	now := time.Now().Format(timeFormat)
	mapCarData(&car, &updatedCar, now)

	if car.Start_time != "" {
		startTime, err := time.Parse(timeFormat, car.Start_time)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Error parsing start time", "error": err.Error()})
		}
		endTime, err := time.Parse(timeFormat, updatedCar.End_time)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Error parsing end time", "error": err.Error()})
		}
		duration := endTime.Sub(startTime)
		updatedCar.Total_payment = math.Round(duration.Minutes() * 10)
		updatedCar.Duration = int(math.Round(duration.Minutes()))
	}

	getCookie := c.Cookies("user")
	updatedCar.User_id = getCookie
	if updatedCar.Reason == "" {
		updatedCar.Reason = "Toleg edildi"
	}

	if err := database.DB.Model(&car).Updates(updatedCar).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Database update failed", "error": err.Error()})
	}

	broadcast <- updatedCar
	return c.Status(200).JSON(fiber.Map{"message": "Car updated successfully", "car": updatedCar})
}

func mapCarData(source, target *modelscar.Car_Model, endTime string) {
	target.ID = source.ID
	target.Car_number = source.Car_number
	target.Start_time = source.Start_time
	target.End_time = endTime
	target.Status = statusExited
	target.Image_Url = defaultImageURL
	target.ParkNo = source.ParkNo
}

// SearchCar godoc
// @Summary Search for a car by plate number and optional filters
// @Description Search for a car by plate number, parking number, and other optional filters
// @Tags cars
// @Accept  json
// @Produce  json
// @Param car_number query string false "Car plate number"
// @Param enter_time query string false "Enter time (YYYY-MM-DD)"
// @Param end_time query string false "End time (YYYY-MM-DD)"
// @Param parkno query string false "Parking spot number"
// @Param status query string false "Car status (Inside, Exited)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(5)
// @Success 200 {object} GetCarsResponse
// @Failure 400 {object} ErrorResponse
// @Router /searchcar [get]
func SearchCar(c *fiber.Ctx) error {
	var cars []modelscar.Car_Model
	var totalCount int64

	carNumber := c.Query("car_number")
	enterTime := c.Query("enter_time")
	endTime := c.Query("end_time")
	parkNo := c.Query("parkno")
	status := c.Query("status")
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "5")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid page number",
		})
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid limit number",
		})
	}

	query := database.DB.Model(&modelscar.Car_Model{})

	if carNumber != "" {
		query = query.Where("car_number LIKE ?", "%"+carNumber+"%")
	}
	if enterTime != "" {
		if _, err := time.Parse("2006-01-02", enterTime); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid enter_time format. Use YYYY-MM-DD.",
			})
		}
		query = query.Where("DATE(start_time) = ?", enterTime)
	}
	if endTime != "" {
		if _, err := time.Parse("2006-01-02", endTime); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid end_time format. Use YYYY-MM-DD.",
			})
		}
		query = query.Where("DATE(end_time) = ?", endTime)
	}
	if parkNo != "" {
		query = query.Where("park_no = ?", parkNo)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Error counting cars",
			"error":   err.Error(),
		})
	}

	offset := (page - 1) * limit
	if err := query.Order("id desc").Limit(limit).Offset(offset).Find(&cars).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Error retrieving cars",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(GetCarsResponse{
		Cars:       cars,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
	})
}

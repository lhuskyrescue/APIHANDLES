package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Calculation struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}
type requestBody struct {
	Task string `json:"task"`
}
type CalculationRequest struct { //структура, отвечающая за вычисления
	Expression string `json:"expression"`
}

var calculations = []Calculation{} //инициализированный слайс
var task string

func calculateExpression(expression string) (string, error) { //функция, для вычисления
	expr, err := govaluate.NewEvaluableExpression(expression) //создание выражения для подсчета
	if err != nil {
		return "", err //если невалидное выражение, например ++,--, возвращаем пустую строку и ошибку
	}
	result, err := expr.Evaluate(nil) //передача результата из первого выражения
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), err

}
func getCalculations(c echo.Context) error {

	return c.JSON(http.StatusOK, calculations)
}
func postCalculations(c echo.Context) error {
	var req CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	result, err := calculateExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}
	calc := Calculation{
		ID:         uuid.NewString(),
		Expression: req.Expression,
		Result:     result,
	}
	calculations = append(calculations, calc)
	return c.JSON(http.StatusCreated, calc)
}
func postTask(c echo.Context) error {
	var reqBody requestBody
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}
	task = reqBody.Task
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Task saved successfully",
		"task":    task,
	})
}
func getHandler(c echo.Context) error {
	return c.String(http.StatusOK, "hello "+task)
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.GET("/calculations", getCalculations)
	e.GET("/", getHandler)
	e.POST("/task", postTask)
	e.POST("/calculations", postCalculations)
	e.Start("localhost:8080")
}

package main

import (
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

type CalculationRequest struct { //структура, отвечающая за вычисления
	Expression string `json:"expression"`
}
type requestBody struct {
	Task string `json:"task"`
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
	response := fmt.Sprintf("%v", task)
	return c.String(http.StatusOK, response)
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
func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.GET("/calculations", getCalculations)
	e.POST("/calculations", postCalculations)
	e.Start("localhost:8080")
}

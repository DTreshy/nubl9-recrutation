package api

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Result struct {
	StdDev float64   `json:"stddev"`
	Data   []float64 `json:"data"`
}

func calculateMean(data []float64) float64 {
	sum := 0.0

	for _, value := range data {
		sum += value
	}

	return sum / float64(len(data))
}

func (r *Result) calculateStandardDeviation() {
	mean := calculateMean(r.Data)

	squaredDifferencesSum := 0.0

	for _, value := range r.Data {
		difference := value - mean
		squaredDifferencesSum += difference * difference
	}

	meanSquaredDifferences := squaredDifferencesSum / float64(len(r.Data))
	r.StdDev = math.Sqrt(meanSquaredDifferences)
}

func Random(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)

	defer cancel()

	queryParams := c.Queries()
	length := "1"
	requests := 1

	if val, ok := queryParams["length"]; ok {
		if err := checkParam(val); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Errorf("length: %w", err).Error(),
			})
		}

		length = val
	}

	if val, ok := queryParams["requests"]; ok {
		if err := checkParam(val); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Errorf("requests: %w", err).Error(),
			})
		}

		num, _ := strconv.Atoi(val)
		requests = num
	}

	var (
		sumResult Result
		response  []Result
	)

	resultCh := make(chan Result)
	errCh := make(chan error)

	for i := 0; i < requests; i++ {
		go fetchData(length, resultCh, errCh)
	}

	for i := 0; i < requests; i++ {
		select {
		case result := <-resultCh:
			response = append(response, result)
			sumResult.Data = append(sumResult.Data, result.Data...)
		case err := <-errCh:
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": err.Error(),
				})
			}
		case <-ctx.Done():
			return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
				"message": "request timeout",
			})
		}
	}

	close(errCh)
	close(resultCh)

	sumResult.calculateStandardDeviation()

	response = append(response, sumResult)

	return c.Status(200).JSON(response)
}

func checkParam(val string) error {
	num, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("must be integer")
	}

	if num < 1 {
		return fmt.Errorf("minimum value is 1")
	}

	if num > 100 {
		return fmt.Errorf("maximum value is 100")
	}

	return nil
}

func fetchData(length string, resultCh chan<- Result, errCh chan<- error) {
	client := http.Client{}

	response, err := client.Get("https://www.random.org/integers/?num=" + length + "&min=1&max=100&col=1&base=10&format=plain&rnd=new")
	if err != nil {
		errCh <- err
		return
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		errCh <- err
		return
	}

	str, _ := strings.CutSuffix(string(body), "\n")
	decimalStrings := strings.Split(str, "\n")

	var result Result

	for _, val := range decimalStrings {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			errCh <- err
			return
		}

		result.Data = append(result.Data, f)
	}

	result.calculateStandardDeviation()

	resultCh <- result
}

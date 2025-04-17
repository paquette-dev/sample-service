package response_test

import (
	"net/http"
	"net/http/httptest"
	"sample-service/internal/response"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestResponse(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Response Suite")
}

var _ = ginkgo.Describe("Response", func() {
	var (
		e   *echo.Echo
		ctx echo.Context
		rec *httptest.ResponseRecorder
	)

	// Setup before each test
	ginkgo.BeforeEach(func() {
		// Create a new Echo instance
		e = echo.New()
		// Create a new recorder
		rec = httptest.NewRecorder()
		// Create a new request
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// Create a new context
		ctx = e.NewContext(req, rec)
	})

	ginkgo.Context("JSONSuccessResponse", func() {
		ginkgo.It("should return a success response with status 200", func() {
			// Test data
			message := "Success message"
			data := map[string]string{"key": "value"}

			// Call the function
			err := response.JSONSuccessResponse(ctx, message, data)

			// Assertions
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusOK))
			
			// Check response body contains expected values
			responseBody := rec.Body.String()
			gomega.Expect(responseBody).To(gomega.ContainSubstring(`"message":"Success message"`))
			gomega.Expect(responseBody).To(gomega.ContainSubstring(`"key":"value"`))
		})

		ginkgo.It("should work with nil data", func() {
			// Test data
			message := "Success with nil data"

			// Call the function
			err := response.JSONSuccessResponse(ctx, message, nil)

			// Assertions
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusOK))
			gomega.Expect(rec.Body.String()).To(gomega.ContainSubstring(`"message":"Success with nil data"`))
			gomega.Expect(rec.Body.String()).To(gomega.ContainSubstring(`"data":null`))
		})

		ginkgo.It("should work with array data", func() {
			// Test data
			message := "Success with array data"
			data := []string{"item1", "item2", "item3"}

			// Call the function
			err := response.JSONSuccessResponse(ctx, message, data)

			// Assertions
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusOK))
			gomega.Expect(rec.Body.String()).To(gomega.ContainSubstring(`"message":"Success with array data"`))
			gomega.Expect(rec.Body.String()).To(gomega.ContainSubstring(`"item1"`))
			gomega.Expect(rec.Body.String()).To(gomega.ContainSubstring(`"item2"`))
			gomega.Expect(rec.Body.String()).To(gomega.ContainSubstring(`"item3"`))
		})
	})

	ginkgo.Context("JSONErrorResponse", func() {
		ginkgo.It("should return an error response with status 500", func() {
			// Test data
			message := "Error message"
			errorStr := "Something went wrong"

			// Call the function
			err := response.JSONErrorResponse(ctx, message, errorStr)

			// Assertions
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusInternalServerError))
			gomega.Expect(rec.Body.String()).To(gomega.ContainSubstring(`"message":"Error message"`))
			gomega.Expect(rec.Body.String()).To(gomega.ContainSubstring(`"error":"Something went wrong"`))
		})

		ginkgo.It("should work with empty error string", func() {
			// Test data
			message := "Error with empty error string"
			errorStr := ""

			// Call the function
			err := response.JSONErrorResponse(ctx, message, errorStr)

			// Assertions
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusInternalServerError))
			gomega.Expect(rec.Body.String()).To(gomega.ContainSubstring(`"message":"Error with empty error string"`))
			gomega.Expect(rec.Body.String()).To(gomega.ContainSubstring(`"error":""`))
		})
	})
})
package controllers_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sample-service/internal/controllers"
	"sample-service/internal/model"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

type MockUserRepository struct {
	users []model.User
	err   error
	exists bool
}

func (m *MockUserRepository) GetAllUsers() ([]model.User, error) {
	return m.users, m.err
}

func (m *MockUserRepository) GetUserByID(id int) (*model.User, error) {
	for _, user := range m.users {
		if int(user.ID) == id {
			return &user, nil
		}
	}
	return nil, m.err
}

func (m *MockUserRepository) CheckIfUsernameExists(username string) (bool, error) {
	return m.exists, m.err
}

func (m *MockUserRepository) CreateUser(user model.User) (*model.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	newUser := user
	
	if newUser.ID == 0 {
		newUser.ID = 1
	}
	
	return &newUser, nil
}

func (m *MockUserRepository) UpdateUser(user model.User) (*model.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	
	updatedUser := user
	
	if updatedUser.ID == 0 {
		updatedUser.ID = 1
	}
	
	return &updatedUser, nil
}

func TestUserController(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "UserController Suite")
}

var _ = ginkgo.Describe("UserController", func() {
	var (
		e              *echo.Echo
		mockUserRepo   *MockUserRepository
		userController *controllers.UserController
		testUser       model.User
	)

	ginkgo.BeforeEach(func() {
		e = echo.New()
		mockUserRepo = &MockUserRepository{}
		userController = controllers.NewUserController(mockUserRepo)
		
		testUser = model.User{
			ID:         1,
			UserName:   "testuser",
			FirstName:  "Test",
			LastName:   "User",
			Email:      "testuser@example.com",
			Department: "IT",
			UserStatus: "A",
		}
	})

	ginkgo.Context("GetAllUsers", func() {
		ginkgo.It("should return all users successfully", func() {
			// Setup - success case
			mockUserRepo.users = []model.User{testUser}
			mockUserRepo.err = nil

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Execute
			err := userController.GetAllUsers(c)

			// Assert
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusOK))

			// Parse response
			var response struct {
				Message string      `json:"message"`
				Data    []model.User `json:"data"`
			}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())

			// Verify response content
			gomega.Expect(response.Message).To(gomega.Equal("Users retrieved successfully"))
			gomega.Expect(response.Data).To(gomega.HaveLen(1))
			gomega.Expect(response.Data[0].ID).To(gomega.Equal(int64(1)))
			gomega.Expect(response.Data[0].UserName).To(gomega.Equal("testuser"))
			gomega.Expect(response.Data[0].FirstName).To(gomega.Equal("Test"))
			gomega.Expect(response.Data[0].LastName).To(gomega.Equal("User"))
			gomega.Expect(response.Data[0].Email).To(gomega.Equal("testuser@example.com"))
			gomega.Expect(response.Data[0].Department).To(gomega.Equal("IT"))
			gomega.Expect(response.Data[0].UserStatus).To(gomega.Equal("A"))
		})

		ginkgo.It("should return error when repository fails", func() {
			// Setup - error case
			mockUserRepo.users = nil
			mockUserRepo.err = errors.New("database error")

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Execute
			err := userController.GetAllUsers(c)

			// Assert - in this case, the controller should have written the error response
			gomega.Expect(err).To(gomega.BeNil()) // Controller handles the error internally
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusInternalServerError))

			// Parse error response
			var response struct {
				Message string `json:"message"`
				Error   string `json:"error"`
			}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response.Message).To(gomega.Equal("Failed to retrieve users"))
			gomega.Expect(response.Error).To(gomega.Equal("database error"))
		})
	})

	ginkgo.Context("GetUserByID", func() {
		ginkgo.It("should return user by id successfully", func() {
			// Setup - success case
			mockUserRepo.users = []model.User{testUser}
			mockUserRepo.err = nil

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
    		c.SetParamValues("1")

			// Execute
			err := userController.GetUserByID(c)

			// Assert
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusOK))

			// Parse response
			var response struct {
				Message string      `json:"message"`
				Data    model.User `json:"data"`
			}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())

			// Verify response content
			gomega.Expect(response.Message).To(gomega.Equal("User retrieved successfully"))
			gomega.Expect(response.Data.ID).To(gomega.Equal(int64(1)))
			gomega.Expect(response.Data.UserName).To(gomega.Equal("testuser"))
			gomega.Expect(response.Data.FirstName).To(gomega.Equal("Test"))
			gomega.Expect(response.Data.LastName).To(gomega.Equal("User"))
			gomega.Expect(response.Data.Email).To(gomega.Equal("testuser@example.com"))
			gomega.Expect(response.Data.Department).To(gomega.Equal("IT"))
			gomega.Expect(response.Data.UserStatus).To(gomega.Equal("A"))
		})

		ginkgo.It("should return error when user not found", func() {
			// Setup - error case
			mockUserRepo.users = nil
			mockUserRepo.err = errors.New("database error")

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
    		c.SetParamValues("1")

			// Execute
			err := userController.GetUserByID(c)

			// Assert - in this case, the controller should have written the error response
			gomega.Expect(err).To(gomega.BeNil()) // Controller handles the error internally
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusInternalServerError))

			// Parse error response
			var response struct {
				Message string `json:"message"`
				Error   string `json:"error"`
			}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response.Message).To(gomega.Equal("User not found"))
			gomega.Expect(response.Error).To(gomega.Equal("database error"))
		})
	})

	ginkgo.Context("CreateUser", func() {
		ginkgo.It("should create user successfully", func() {
			// Setup - success case
			mockUserRepo.users = nil
			mockUserRepo.err = nil
			mockUserRepo.exists = false
			
			// Create a copy of testUser that will be returned by the mock
			returnedUser := testUser
			returnedUser.ID = 1
			
			// Set up the mock to return our user with ID=1
			mockUserRepo.users = []model.User{returnedUser}
			
			// Create request with JSON body
			requestBody := `{
				"user_name": "testuser",
				"first_name": "Test",
				"last_name": "User",
				"email": "testuser@example.com",
				"department": "IT",
				"user_status": "A"
			}`
			
			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			
			// Execute
			err := userController.CreateUser(c)
			
			// Assert
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusOK))
			
			// Parse response
			var response struct {
				Message string     `json:"message"`
				Data    model.User `json:"data"`
			}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			
			// Verify response content
			gomega.Expect(response.Message).To(gomega.Equal("User created successfully"))
			gomega.Expect(response.Data.ID).To(gomega.Equal(int64(1)))
			gomega.Expect(response.Data.UserName).To(gomega.Equal("testuser"))
			gomega.Expect(response.Data.FirstName).To(gomega.Equal("Test"))
			gomega.Expect(response.Data.LastName).To(gomega.Equal("User"))
			gomega.Expect(response.Data.Email).To(gomega.Equal("testuser@example.com"))
			gomega.Expect(response.Data.Department).To(gomega.Equal("IT"))
			gomega.Expect(response.Data.UserStatus).To(gomega.Equal("A"))
		})

		ginkgo.It("should return error when username already exists", func() {
			// Setup - error case
			mockUserRepo.users = nil
			mockUserRepo.err = fmt.Errorf("username 'testuser' already exists")
			mockUserRepo.exists = true
			
			// Create request with JSON body
			requestBody := `{
				"user_name": "testuser",
				"first_name": "Test",
				"last_name": "User",
				"email": "testuser@example.com",
				"department": "IT",
				"user_status": "A"
			}`
			
			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			
			// Execute
			err := userController.CreateUser(c)
			
			// Assert - in this case, the controller should have written the error response
			gomega.Expect(err).To(gomega.BeNil()) // Controller handles the error internally
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusInternalServerError))
			
			// Parse error response
			var response struct {
				Message string `json:"message"`
				Error   string `json:"error"`
			}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response.Message).To(gomega.Equal("Username already exists"))
			gomega.Expect(response.Error).To(gomega.Equal("username 'testuser' already exists"))
		})

		ginkgo.It("should return error when database error occurs", func() {
			// Setup - error case
			mockUserRepo.users = nil
			mockUserRepo.err = errors.New("database error")
			mockUserRepo.exists = false

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/users", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Execute
			err := userController.CreateUser(c)

			// Assert - in this case, the controller should have written the error response
			gomega.Expect(err).To(gomega.BeNil()) // Controller handles the error internally
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusInternalServerError))

			// Parse error response
			var response struct {
				Message string `json:"message"`
				Error   string `json:"error"`
			}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response.Message).To(gomega.Equal("Failed to create user"))
			gomega.Expect(response.Error).To(gomega.Equal("database error"))
		})
	})

	ginkgo.Context("UpdateUser", func() {
		ginkgo.It("should update user successfully", func() {
			// Setup - success case
			mockUserRepo.users = nil
			mockUserRepo.err = nil
			
			// Create a copy of testUser that will be returned by the mock
			returnedUser := testUser
			returnedUser.ID = 1
			
			// Set up the mock to return our user with ID=1
			mockUserRepo.users = []model.User{returnedUser}
			
			// Create request with JSON body
			requestBody := `{
				"user_name": "testuser",
				"first_name": "Test",
				"last_name": "User",
				"email": "testuser@example.com",
				"department": "IT",
				"user_status": "A"
			}`
			
			req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
    		c.SetParamValues("1")
			
			// Execute
			err := userController.UpdateUser(c)
			
			// Assert
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusOK))
			
			// Parse response
			var response struct {
				Message string     `json:"message"`
				Data    model.User `json:"data"`
			}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			
			// Verify response content
			gomega.Expect(response.Message).To(gomega.Equal("User updated successfully"))
			gomega.Expect(response.Data.ID).To(gomega.Equal(int64(1)))
			gomega.Expect(response.Data.UserName).To(gomega.Equal("testuser"))
			gomega.Expect(response.Data.FirstName).To(gomega.Equal("Test"))
			gomega.Expect(response.Data.LastName).To(gomega.Equal("User"))
			gomega.Expect(response.Data.Email).To(gomega.Equal("testuser@example.com"))
			gomega.Expect(response.Data.Department).To(gomega.Equal("IT"))
			gomega.Expect(response.Data.UserStatus).To(gomega.Equal("A"))
		})

		ginkgo.It("should return error when database error occurs", func() {
			// Setup - error case
			mockUserRepo.users = nil
			mockUserRepo.err = errors.New("database error")

			// Create request
			req := httptest.NewRequest(http.MethodPut, "/users/1", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Execute
			err := userController.UpdateUser(c)

			// Assert - in this case, the controller should have written the error response
			gomega.Expect(err).To(gomega.BeNil()) // Controller handles the error internally
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusInternalServerError))

			// Parse error response
			var response struct {
				Message string `json:"message"`
				Error   string `json:"error"`
			}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(response.Message).To(gomega.Equal("Failed to update user"))
			gomega.Expect(response.Error).To(gomega.Equal("database error"))
		})
	})
})
	


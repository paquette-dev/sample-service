package repository_test

import (
	"database/sql"
	"errors"
	"sample-service/internal/model"
	"sample-service/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestUserRepository(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "UserRepository Suite")
}

var expectedUsers = []model.User{
	{
		ID:         1,
		UserName:   "johndoe",
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@company.com",
		Department: "Engineering",
		UserStatus: "A",
	},
	{
		ID:         2,
		UserName:   "janesmith",
		FirstName:  "Jane",
		LastName:   "Smith",
		Email:      "jane.smith@company.com",
		Department: "Marketing",
		UserStatus: "A",
	},
}

var _ = ginkgo.Describe("UserRepository", func() {
	var (
		mockDB     *sql.DB
		mock       sqlmock.Sqlmock
		userRepo   repository.UserRepository
		err        error
	)

	ginkgo.BeforeEach(func() {
		mockDB, mock, err = sqlmock.New()
		if err != nil {
			ginkgo.Fail("Failed to create mock database: " + err.Error())
		}

		userRepo = repository.NewUserRepository(mockDB)
	})

	ginkgo.AfterEach(func() {
		mockDB.Close()
	})

	ginkgo.Context("GetAllUsers", func() {
		ginkgo.It("should return all users", func() {
			// Setup the expected query
			rows := sqlmock.NewRows([]string{"user_id", "user_name", "first_name", "last_name", "email", "department", "user_status"})
			
			// Add rows to the mock result
			for _, user := range expectedUsers {
				rows.AddRow(user.ID, user.UserName, user.FirstName, user.LastName, user.Email, user.Department, user.UserStatus)
			}

			// Expect the query to be executed
			mock.ExpectQuery("SELECT \\* FROM users").WillReturnRows(rows)

			// Call the function
			users, err := userRepo.GetAllUsers()

			// Assertions
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(users).To(gomega.HaveLen(2))
			gomega.Expect(users).To(gomega.Equal(expectedUsers))

			// Verify all expectations were met
			err = mock.ExpectationsWereMet()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

		})

		ginkgo.It("should return an error when the database query fails", func() {
			// Setup the expected query
			expectedError := errors.New("database query failed")
			mock.ExpectQuery("SELECT \\* FROM users").WillReturnError(expectedError)

			// Call the function
			users, err := userRepo.GetAllUsers()

			// Assertions
			gomega.Expect(err).To(gomega.Equal(expectedError))
			gomega.Expect(users).To(gomega.BeNil())

			// Verify all expectations were met
			err = mock.ExpectationsWereMet()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Context("GetUserByID", func() {
		ginkgo.It("should return a user by ID", func() {
			// Setup the expected query
			rows := sqlmock.NewRows([]string{"user_id", "user_name", "first_name", "last_name", "email", "department", "user_status"})
			
			// Add a single row for the expected user
			expectedUser := expectedUsers[0]
			rows.AddRow(expectedUser.ID, expectedUser.UserName, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email, expectedUser.Department, expectedUser.UserStatus)

			// Expect the query to be executed
			mock.ExpectQuery("SELECT \\* FROM users WHERE user_id = \\?").WithArgs(1).WillReturnRows(rows)

			// Call the function
			user, err := userRepo.GetUserByID(1)

			// Assertions
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(*user).To(gomega.Equal(expectedUser))

			// Verify all expectations were met
			err = mock.ExpectationsWereMet()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should return an error when the database query fails", func() {
			// Setup the expected query
			expectedError := errors.New("database query failed")
			mock.ExpectQuery("SELECT \\* FROM users WHERE user_id = \\?").WithArgs(1).WillReturnError(expectedError)

			// Call the function
			user, err := userRepo.GetUserByID(1)

			// Assertions
			gomega.Expect(err).To(gomega.Equal(expectedError))
			gomega.Expect(user).To(gomega.BeNil())

			// Verify all expectations were met
			err = mock.ExpectationsWereMet()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Context("CheckIfUsernameExists", func() {
		ginkgo.It("should return true if the username exists", func() {
			// Setup the expected query with a count column
			rows := sqlmock.NewRows([]string{"count"}).
				AddRow(1)

			// Expect the query to be executed
			mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE user_name = \\?").
				WithArgs("johndoe").
				WillReturnRows(rows)

			// Call the function
			exists, err := userRepo.CheckIfUsernameExists("johndoe")

			// Assertions
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(exists).To(gomega.BeTrue())

			// Verify all expectations were met
			err = mock.ExpectationsWereMet()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should return false if the username does not exist", func() {
			// Setup the expected query with a count column
			rows := sqlmock.NewRows([]string{"count"}).
				AddRow(0)

			// Expect the query to be executed
			mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE user_name = \\?").
				WithArgs("johndoe").
				WillReturnRows(rows)

			// Call the function
			exists, err := userRepo.CheckIfUsernameExists("johndoe")

			// Assertions
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(exists).To(gomega.BeFalse())

			// Verify all expectations were met
			err = mock.ExpectationsWereMet()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Context("CreateUser", func() {
		ginkgo.It("should create a new user", func() {
			// Setup the expected user
			expectedUser := expectedUsers[0]
    
			// First, mock the username check query
			mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE user_name = \\?").
				WithArgs(expectedUser.UserName).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			
			// Then, mock the insert query
			mock.ExpectExec("INSERT INTO users \\(user_name, first_name, last_name, email, department, user_status\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?\\)").
				WithArgs(
					expectedUser.UserName,
					expectedUser.FirstName,
					expectedUser.LastName,
					expectedUser.Email,
					expectedUser.Department,
					expectedUser.UserStatus,
				).
				WillReturnResult(sqlmock.NewResult(1, 1)) // id=1, affected=1
			
			// Call the function
			user, err := userRepo.CreateUser(expectedUser)
			
			// Assertions
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(user.ID).To(gomega.Equal(int64(1))) // ID should be set from LastInsertId
			
			// Create a copy of expectedUser with ID=1 for comparison
			expectedUserWithID := expectedUser
			expectedUserWithID.ID = 1
			gomega.Expect(*user).To(gomega.Equal(expectedUserWithID))
			
			// Verify all expectations were met
			err = mock.ExpectationsWereMet()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should return error when username already exists", func() {
			// Setup the expected user
			expectedUser := expectedUsers[0]

			// First, mock the username check query
			mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE user_name = \\?").
				WithArgs(expectedUser.UserName).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

			// Call the function
			_, err := userRepo.CreateUser(expectedUser)	
			
			// Assertions
			gomega.Expect(err).To(gomega.HaveOccurred())

			// Verify all expectations were met
			err = mock.ExpectationsWereMet()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should return an error when the database query fails", func() {
			// Setup the expected user
			expectedUser := expectedUsers[0]
    
			// First, mock the username check query
			mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE user_name = \\?").
				WithArgs(expectedUser.UserName).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

			// Setup the expected query
			expectedError := errors.New("database query failed")
			mock.ExpectExec("INSERT INTO users \\(user_name, first_name, last_name, email, department, user_status\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?\\)").
				WithArgs(
					expectedUser.UserName,
					expectedUser.FirstName,
					expectedUser.LastName,
					expectedUser.Email,
					expectedUser.Department,
					expectedUser.UserStatus,
				).
				WillReturnError(expectedError)

			// Call the function
			_, err := userRepo.CreateUser(expectedUser)

			// Assertions
			gomega.Expect(err).To(gomega.Equal(expectedError))

			// Verify all expectations were met
			err = mock.ExpectationsWereMet()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})
})	

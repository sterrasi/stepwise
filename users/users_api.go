package users

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/nu7hatch/gouuid"
	"github.com/sterrasi/stepwise/resource"
)

// UserRegistration data
type UserRegistration struct {
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName;omitempty"`
	LastName   string  `json:"lastName"`
	Nickname   string  `json:"nickname"`

	UserName string `json:"userName"`
	Password string `json:"password"`

	PrimaryEmail string `json:"primaryEmail"`
	Organization string `json:"organization"`
	Gender       string `json:"gender" valid:"in(m|M|F|f)"`

	VerificationCode string `json:"verificationCode"`
}

// Validate the UserRegistration
// this has to be in reference to fields.. find an elegant way to build that up
// for now just use an array of strings
func (reg *UserRegistration) Validate() ([]string, error) {

	response := make([]string, 0)

	if reg.FirstName == "" {
		response = append(response, "First Name is required")
	}

	if reg.LastName == "" {
		response = append(response, "Last Name is required")
	}

	// email
	if reg.PrimaryEmail == "" {
		response = append(response, "Primary Email is required")
	}
	if !govalidator.IsEmail(reg.PrimaryEmail) {
		response = append(response, "Primary Email is malformed")
	}

	// make sure it is unique
	safe, err := notAlreadyRegistered(reg.PrimaryEmail)
	if err != nil {
		return nil, fmt.Errorf("System Error: %s", err.Error())
	}
	if !safe {
		response = append(response, "Primary Email is taken")
	}

	// organization
	if reg.Organization == "" {
		response = append(response, "Organization is required")
	}

	// gender
	if reg.Gender == "" {
		response = append(response, "Gender is required")
	}
	if strings.ToLower(reg.Gender) != "m" || strings.ToLower(reg.Gender) != "f" {
		response = append(response, "Invalid Gender Specified")
	}
	return response, nil
}

// Config is the configuration for the user API
type Config struct {
	ResultsPerPage int `mapstructure:"default-results-per-page"`
}

// Register initializes the users package
func Register(e *echo.Group, database *gorm.DB, config *Config) {
	db = database

	resultsPerPage := string(config.ResultsPerPage)

	/*
	 * get users
	 *   offset - [int] (default: 0) offset into the index
	 *   limit  - [int] (default: 20) number of results to return
	 */
	e.GET("", func(c echo.Context) error {
		var offset, limit int

		if err := resource.Param("offset").Optional("0").Int(c, &offset); err != nil {
			return resource.BadRequest(err)
		}
		if err := resource.Param("limit").Optional(resultsPerPage).Int(c, &limit); err != nil {
			return resource.BadRequest(err)
		}

		users, err := GetUsers(offset, limit)
		if err != nil {
			return resource.InternalServerError(err)
		}
		return c.JSON(http.StatusOK, users)
	})

	/*
	 * register a new user
	 */
	e.POST("register", func(c echo.Context) error {

		// get registration
		registration := &UserRegistration{}
		if err := c.Bind(registration); err != nil {
			return resource.BadRequest(err)
		}
		issues, err := registration.Validate()
		if err != nil {
			return resource.InternalServerError(err)
		}
		if len(issues) > 0 {
			return resource.BadRequest(issues)
		}

		// create the verification code
		code, err := uuid.NewV4()
		if err != nil {
			return resource.InternalServerError(fmt.Errorf("Unable to create verification code :%s", err))
		}
		registration.VerificationCode = code.String()

		// register user
		// TODO: finish this
		//entity, err := RegisterUser(registration)
		// if err != nil {
		// 	return resource.InternalServerError(err)
		// }
		// return resource.Created(c, entity.GetID())

		return nil
	})

	resource.PatchMethod(e, newInstance, PatchUser)
	resource.GetMethod(e, GetUser)
	resource.UpdateMethod(e, newInstance, UpdateUser)
	resource.DeleteMethod(e, DeleteUser)
}

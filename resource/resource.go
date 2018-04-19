package resource

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sterrasi/stepwise/util"
)

type deleteFn func(int) error
type getByIDFn func(int) (util.Entity, error)
type newInstanceFn func() interface{}
type updateFn func(interface{}) error
type patchFn func(interface{}) error
type createFn func(interface{}) (util.Entity, error)

// DeleteMethod creates a standard delete method for a resource
func DeleteMethod(e *echo.Group, fn deleteFn) {
	e.DELETE("/:id", func(c echo.Context) error {
		var id int

		if err := Param("id").InPath().Int(c, &id); err != nil {
			return BadRequest(err)
		}
		if err := fn(id); err != nil {
			switch err {
			case util.ErrNotFound:
				return NotFound(err)
			default:
				return InternalServerError(err)
			}
		}
		return c.NoContent(http.StatusNoContent)
	})
}

// UpdateMethod creates a standard restful update method
func UpdateMethod(e *echo.Group, newFn newInstanceFn, upFn updateFn) {
	e.PUT("/:id", func(c echo.Context) error {
		var id int

		if err := Param("id").InPath().Int(c, &id); err != nil {
			return BadRequest(err)
		}
		resource := newFn()
		if err := c.Bind(resource); err != nil {
			return BadRequest(err)
		}

		if err := upFn(resource); err != nil {
			return InternalServerError(err)
		}

		return c.NoContent(http.StatusNoContent)
	})

}

// PatchMethod creates a standard restful patch method
func PatchMethod(e *echo.Group, newFn newInstanceFn, pFn patchFn) {
	e.PATCH("/:id", func(c echo.Context) error {
		var id int

		if err := Param("id").InPath().Int(c, &id); err != nil {
			return BadRequest(err)
		}
		resource := newFn()
		if err := c.Bind(resource); err != nil {
			return BadRequest(err)
		}

		if err := pFn(resource); err != nil {
			return InternalServerError(err)
		}

		return c.NoContent(http.StatusNoContent)
	})

}

// GetMethod creates a get method to retrieve a resource by ID
func GetMethod(e *echo.Group, fn getByIDFn) {

	e.GET("/:id", func(c echo.Context) error {
		var id int

		if err := Param("id").InPath().Int(c, &id); err != nil {
			return BadRequest(err)
		}

		resource, err := fn(id)
		if err != nil {
			return InternalServerError(err)
		}

		return c.JSON(http.StatusOK, resource)
	})
}

// CreateMethod creates an instance of a Resource
func CreateMethod(e *echo.Group, newFn newInstanceFn, crFn createFn) {
	e.POST("", func(c echo.Context) error {
		resource := newFn()
		if err := c.Bind(resource); err != nil {
			return BadRequest(err)
		}

		entity, err := crFn(resource)
		if err != nil {
			return InternalServerError(err)
		}
		return Created(c, entity.GetID())
	})
}

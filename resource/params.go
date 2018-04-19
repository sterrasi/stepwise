package resource

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Location type of parameter on an HTTP Request
type Location int

const (
	// Query parameter is on the HTTP query string
	Query Location = iota

	// Path parameter is part of the HTTP Request path
	Path
)

// ValueType The type of values allowed
type ValueType int

const (
	// Int type
	Int ValueType = iota

	// String type
	String
)

// HTTPParam descriptor for an HTTP parameter
type HTTPParam struct {
	name         string
	defaultValue string
	location     Location
	required     bool
	valueType    ValueType
}

// HTTPParamValue parameter value for http
type HTTPParamValue struct {
	IntValue    int
	UIntValue   uint
	StringValue string
}

// Param creates an HTTPParam
func Param(name string) *HTTPParam {
	return &HTTPParam{name: name, required: true, valueType: String, location: Query}
}

// Optional sets the param as optional
func (p *HTTPParam) Optional(defaultValue string) *HTTPParam {
	p.defaultValue = defaultValue
	p.required = false
	return p
}

// Required sets the param as required
func (p *HTTPParam) Required() *HTTPParam {
	p.required = true
	return p
}

// Int sets the param's type to int
func (p *HTTPParam) Int(c echo.Context, value *int) error {
	stringValue, err := p.getValue(c)
	if err != nil {
		return err
	}

	// convert to integer
	intValue, err := strconv.Atoi(stringValue)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("Cannot convert %s parameter value (%s) to integer",
				p.name, stringValue))
	}
	*value = intValue
	return nil
}

// Int sets the param's type to int
func (p *HTTPParam) String(c echo.Context, value *string) error {
	stringValue, err := p.getValue(c)
	if err != nil {
		return err
	}
	*value = stringValue
	return nil
}

// InPath identifies the parameter in the request path
func (p *HTTPParam) InPath() *HTTPParam {
	p.location = Path
	return p
}

func (p *HTTPParam) getValue(c echo.Context) (string, error) {
	// get value string
	var value string

	switch p.location {
	case Path:
		value = c.Param(p.name)

	default: // default to Query type
		value = c.QueryParam(p.name)
	}
	if value == "" {
		if p.required {
			return "", echo.NewHTTPError(http.StatusBadRequest,
				fmt.Sprintf("Required parameter %s not provided", p.name))
		}
		value = p.defaultValue
	}
	return value, nil
}

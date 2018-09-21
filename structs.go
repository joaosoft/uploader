package validator

import (
	"reflect"

	"github.com/joaosoft/logger"
)

func (v *Validator) init() {
	v.handlersBefore = v.NewDefaultBeforeHandlers()
	v.handlersMiddle = v.NewDefaultMiddleHandlers()
	v.handlersAfter = v.NewDefaultPosHandlers()
	v.activeHandlers = v.NewActiveHandlers()

}

type Validator struct {
	tag              string
	activeHandlers   map[string]bool
	handlersBefore   map[string]BeforeTagHandler
	handlersMiddle   map[string]MiddleTagHandler
	handlersAfter    map[string]AfterTagHandler
	errorCodeHandler ErrorCodeHandler
	callbacks        map[string]CallbackHandler
	sanitize         []string
	log              logger.ILogger
	validateAll      bool
}

type ErrorCodeHandler func(code string, arguments []interface{}, name string, value reflect.Value, expected interface{}, err *[]error) error
type CallbackHandler func(name string, value reflect.Value, expected interface{}, err *[]error) []error

type BeforeTagHandler func(name string, value reflect.Value, expected interface{}) []error
type MiddleTagHandler func(name string, value reflect.Value, expected interface{}, err *[]error) []error
type AfterTagHandler func(name string, value reflect.Value, expected interface{}, err *[]error) []error

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ValidatorHandler struct {
	validator *Validator
	values map[string]interface{}
}
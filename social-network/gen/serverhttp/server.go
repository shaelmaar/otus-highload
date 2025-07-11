// Package serverhttp provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package serverhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// BirthDate Дата рождения
type BirthDate = openapi_types.Date

// User defines model for User.
type User struct {
	// Biography Интересы
	Biography *string `json:"biography,omitempty"`

	// Birthdate Дата рождения
	Birthdate *BirthDate `json:"birthdate,omitempty"`

	// City Город
	City *string `json:"city,omitempty"`

	// FirstName Имя
	FirstName *string `json:"first_name,omitempty"`

	// Id Идентификатор пользователя
	Id *UserId `json:"id,omitempty"`

	// SecondName Фамилия
	SecondName *string `json:"second_name,omitempty"`
}

// UserId Идентификатор пользователя
type UserId = string

// N5xx defines model for 5xx.
type N5xx struct {
	// Code Код ошибки. Предназначен для классификации проблем и более быстрого решения проблем.
	Code *int `json:"code,omitempty"`

	// Message Описание ошибки
	Message string `json:"message"`

	// RequestId Идентификатор запроса. Предназначен для более быстрого поиска проблем.
	RequestId *string `json:"request_id,omitempty"`
}

// PostLoginJSONBody defines parameters for PostLogin.
type PostLoginJSONBody struct {
	// Id Идентификатор пользователя
	Id       *UserId `json:"id,omitempty"`
	Password *string `json:"password,omitempty"`
}

// PostUserRegisterJSONBody defines parameters for PostUserRegister.
type PostUserRegisterJSONBody struct {
	Biography *string `json:"biography,omitempty"`

	// Birthdate Дата рождения
	Birthdate  *BirthDate `json:"birthdate,omitempty"`
	City       *string    `json:"city,omitempty"`
	FirstName  *string    `json:"first_name,omitempty"`
	Password   *string    `json:"password,omitempty"`
	SecondName *string    `json:"second_name,omitempty"`
}

// PostLoginJSONRequestBody defines body for PostLogin for application/json ContentType.
type PostLoginJSONRequestBody PostLoginJSONBody

// PostUserRegisterJSONRequestBody defines body for PostUserRegister for application/json ContentType.
type PostUserRegisterJSONRequestBody PostUserRegisterJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /login)
	PostLogin(ctx echo.Context) error

	// (GET /user/get/{id})
	GetUserGetId(ctx echo.Context, id UserId) error

	// (POST /user/register)
	PostUserRegister(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostLogin converts echo context to params.
func (w *ServerInterfaceWrapper) PostLogin(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostLogin(ctx)
	return err
}

// GetUserGetId converts echo context to params.
func (w *ServerInterfaceWrapper) GetUserGetId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id UserId

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetUserGetId(ctx, id)
	return err
}

// PostUserRegister converts echo context to params.
func (w *ServerInterfaceWrapper) PostUserRegister(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostUserRegister(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/login", wrapper.PostLogin)
	router.GET(baseURL+"/user/get/:id", wrapper.GetUserGetId)
	router.POST(baseURL+"/user/register", wrapper.PostUserRegister)

}

type N5xxResponseHeaders struct {
	RetryAfter int
}
type N5xxJSONResponse struct {
	Body struct {
		// Code Код ошибки. Предназначен для классификации проблем и более быстрого решения проблем.
		Code *int `json:"code,omitempty"`

		// Message Описание ошибки
		Message string `json:"message"`

		// RequestId Идентификатор запроса. Предназначен для более быстрого поиска проблем.
		RequestId *string `json:"request_id,omitempty"`
	}

	Headers N5xxResponseHeaders
}

type PostLoginRequestObject struct {
	Body *PostLoginJSONRequestBody
}

type PostLoginResponseObject interface {
	VisitPostLoginResponse(w http.ResponseWriter) error
}

type PostLogin200JSONResponse struct {
	Token *string `json:"token,omitempty"`
}

func (response PostLogin200JSONResponse) VisitPostLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PostLogin400Response struct {
}

func (response PostLogin400Response) VisitPostLoginResponse(w http.ResponseWriter) error {
	w.WriteHeader(400)
	return nil
}

type PostLogin404Response struct {
}

func (response PostLogin404Response) VisitPostLoginResponse(w http.ResponseWriter) error {
	w.WriteHeader(404)
	return nil
}

type PostLogin500JSONResponse struct{ N5xxJSONResponse }

func (response PostLogin500JSONResponse) VisitPostLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Retry-After", fmt.Sprint(response.Headers.RetryAfter))
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response.Body)
}

type PostLogin503ResponseHeaders struct {
	RetryAfter int
}

type PostLogin503JSONResponse struct {
	Body struct {
		// Code Код ошибки. Предназначен для классификации проблем и более быстрого решения проблем.
		Code *int `json:"code,omitempty"`

		// Message Описание ошибки
		Message string `json:"message"`

		// RequestId Идентификатор запроса. Предназначен для более быстрого поиска проблем.
		RequestId *string `json:"request_id,omitempty"`
	}
	Headers PostLogin503ResponseHeaders
}

func (response PostLogin503JSONResponse) VisitPostLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Retry-After", fmt.Sprint(response.Headers.RetryAfter))
	w.WriteHeader(503)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetUserGetIdRequestObject struct {
	Id UserId `json:"id"`
}

type GetUserGetIdResponseObject interface {
	VisitGetUserGetIdResponse(w http.ResponseWriter) error
}

type GetUserGetId200JSONResponse User

func (response GetUserGetId200JSONResponse) VisitGetUserGetIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetUserGetId400Response struct {
}

func (response GetUserGetId400Response) VisitGetUserGetIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(400)
	return nil
}

type GetUserGetId404Response struct {
}

func (response GetUserGetId404Response) VisitGetUserGetIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(404)
	return nil
}

type GetUserGetId500JSONResponse struct{ N5xxJSONResponse }

func (response GetUserGetId500JSONResponse) VisitGetUserGetIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Retry-After", fmt.Sprint(response.Headers.RetryAfter))
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetUserGetId503ResponseHeaders struct {
	RetryAfter int
}

type GetUserGetId503JSONResponse struct {
	Body struct {
		// Code Код ошибки. Предназначен для классификации проблем и более быстрого решения проблем.
		Code *int `json:"code,omitempty"`

		// Message Описание ошибки
		Message string `json:"message"`

		// RequestId Идентификатор запроса. Предназначен для более быстрого поиска проблем.
		RequestId *string `json:"request_id,omitempty"`
	}
	Headers GetUserGetId503ResponseHeaders
}

func (response GetUserGetId503JSONResponse) VisitGetUserGetIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Retry-After", fmt.Sprint(response.Headers.RetryAfter))
	w.WriteHeader(503)

	return json.NewEncoder(w).Encode(response.Body)
}

type PostUserRegisterRequestObject struct {
	Body *PostUserRegisterJSONRequestBody
}

type PostUserRegisterResponseObject interface {
	VisitPostUserRegisterResponse(w http.ResponseWriter) error
}

type PostUserRegister200JSONResponse struct {
	UserId *string `json:"user_id,omitempty"`
}

func (response PostUserRegister200JSONResponse) VisitPostUserRegisterResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PostUserRegister400Response struct {
}

func (response PostUserRegister400Response) VisitPostUserRegisterResponse(w http.ResponseWriter) error {
	w.WriteHeader(400)
	return nil
}

type PostUserRegister500JSONResponse struct{ N5xxJSONResponse }

func (response PostUserRegister500JSONResponse) VisitPostUserRegisterResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Retry-After", fmt.Sprint(response.Headers.RetryAfter))
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response.Body)
}

type PostUserRegister503ResponseHeaders struct {
	RetryAfter int
}

type PostUserRegister503JSONResponse struct {
	Body struct {
		// Code Код ошибки. Предназначен для классификации проблем и более быстрого решения проблем.
		Code *int `json:"code,omitempty"`

		// Message Описание ошибки
		Message string `json:"message"`

		// RequestId Идентификатор запроса. Предназначен для более быстрого поиска проблем.
		RequestId *string `json:"request_id,omitempty"`
	}
	Headers PostUserRegister503ResponseHeaders
}

func (response PostUserRegister503JSONResponse) VisitPostUserRegisterResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Retry-After", fmt.Sprint(response.Headers.RetryAfter))
	w.WriteHeader(503)

	return json.NewEncoder(w).Encode(response.Body)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (POST /login)
	PostLogin(ctx context.Context, request PostLoginRequestObject) (PostLoginResponseObject, error)

	// (GET /user/get/{id})
	GetUserGetId(ctx context.Context, request GetUserGetIdRequestObject) (GetUserGetIdResponseObject, error)

	// (POST /user/register)
	PostUserRegister(ctx context.Context, request PostUserRegisterRequestObject) (PostUserRegisterResponseObject, error)
}

type StrictHandlerFunc = strictecho.StrictEchoHandlerFunc
type StrictMiddlewareFunc = strictecho.StrictEchoMiddlewareFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// PostLogin operation middleware
func (sh *strictHandler) PostLogin(ctx echo.Context) error {
	var request PostLoginRequestObject

	var body PostLoginJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	request.Body = &body

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.PostLogin(ctx.Request().Context(), request.(PostLoginRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostLogin")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(PostLoginResponseObject); ok {
		return validResponse.VisitPostLoginResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// GetUserGetId operation middleware
func (sh *strictHandler) GetUserGetId(ctx echo.Context, id UserId) error {
	var request GetUserGetIdRequestObject

	request.Id = id

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetUserGetId(ctx.Request().Context(), request.(GetUserGetIdRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetUserGetId")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(GetUserGetIdResponseObject); ok {
		return validResponse.VisitGetUserGetIdResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// PostUserRegister operation middleware
func (sh *strictHandler) PostUserRegister(ctx echo.Context) error {
	var request PostUserRegisterRequestObject

	var body PostUserRegisterJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	request.Body = &body

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.PostUserRegister(ctx.Request().Context(), request.(PostUserRegisterRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostUserRegister")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(PostUserRegisterResponseObject); ok {
		return validResponse.VisitPostUserRegisterResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

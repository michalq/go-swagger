// Code generated by go-swagger; DO NOT EDIT.

package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/michalq/go-swagger/examples/contributed-templates/stratoscale/restapi/operations"
	"github.com/michalq/go-swagger/examples/contributed-templates/stratoscale/restapi/operations/pet"
	"github.com/michalq/go-swagger/examples/contributed-templates/stratoscale/restapi/operations/store"
)

type contextKey string

const AuthKey contextKey = "Auth"

//go:generate mockery -name PetAPI -inpkg

// PetAPI
type PetAPI interface {
	PetCreate(ctx context.Context, params pet.PetCreateParams) middleware.Responder
	PetDelete(ctx context.Context, params pet.PetDeleteParams) middleware.Responder
	PetGet(ctx context.Context, params pet.PetGetParams) middleware.Responder
	PetList(ctx context.Context, params pet.PetListParams) middleware.Responder
	PetUpdate(ctx context.Context, params pet.PetUpdateParams) middleware.Responder
	PetUploadImage(ctx context.Context, params pet.PetUploadImageParams) middleware.Responder
}

//go:generate mockery -name StoreAPI -inpkg

// StoreAPI
type StoreAPI interface {
	InventoryGet(ctx context.Context, params store.InventoryGetParams) middleware.Responder
	OrderCreate(ctx context.Context, params store.OrderCreateParams) middleware.Responder
	// OrderDelete is For valid response try integer IDs with positive integer value. Negative or non-integer values will generate API errors
	OrderDelete(ctx context.Context, params store.OrderDeleteParams) middleware.Responder
	// OrderGet is For valid response try integer IDs with value >= 1 and <= 10. Other values will generated exceptions
	OrderGet(ctx context.Context, params store.OrderGetParams) middleware.Responder
}

// Config is configuration for Handler
type Config struct {
	PetAPI
	StoreAPI
	Logger func(string, ...interface{})
	// InnerMiddleware is for the handler executors. These do not apply to the swagger.json document.
	// The middleware executes after routing but before authentication, binding and validation
	InnerMiddleware func(http.Handler) http.Handler

	// Authorizer is used to authorize a request after the Auth function was called using the "Auth*" functions
	// and the principal was stored in the context in the "AuthKey" context value.
	Authorizer func(*http.Request) error

	// AuthRoles Applies when the "X-Auth-Roles" header is set
	AuthRoles func(token string) (interface{}, error)
}

// Handler returns an http.Handler given the handler configuration
// It mounts all the business logic implementers in the right routing.
func Handler(c Config) (http.Handler, error) {
	h, _, err := HandlerAPI(c)
	return h, err
}

// HandlerAPI returns an http.Handler given the handler configuration
// and the corresponding *Petstore instance.
// It mounts all the business logic implementers in the right routing.
func HandlerAPI(c Config) (http.Handler, *operations.PetstoreAPI, error) {
	spec, err := loads.Analyzed(swaggerCopy(SwaggerJSON), "")
	if err != nil {
		return nil, nil, fmt.Errorf("analyze swagger: %v", err)
	}
	api := operations.NewPetstoreAPI(spec)
	api.ServeError = errors.ServeError
	api.Logger = c.Logger

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()
	api.RolesAuth = func(token string) (interface{}, error) {
		if c.AuthRoles == nil {
			return token, nil
		}
		return c.AuthRoles(token)
	}

	api.APIAuthorizer = authorizer(c.Authorizer)
	api.StoreInventoryGetHandler = store.InventoryGetHandlerFunc(func(params store.InventoryGetParams, principal interface{}) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		ctx = storeAuth(ctx, principal)
		return c.StoreAPI.InventoryGet(ctx, params)
	})
	api.StoreOrderCreateHandler = store.OrderCreateHandlerFunc(func(params store.OrderCreateParams, principal interface{}) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		ctx = storeAuth(ctx, principal)
		return c.StoreAPI.OrderCreate(ctx, params)
	})
	api.StoreOrderDeleteHandler = store.OrderDeleteHandlerFunc(func(params store.OrderDeleteParams, principal interface{}) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		ctx = storeAuth(ctx, principal)
		return c.StoreAPI.OrderDelete(ctx, params)
	})
	api.StoreOrderGetHandler = store.OrderGetHandlerFunc(func(params store.OrderGetParams, principal interface{}) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		ctx = storeAuth(ctx, principal)
		return c.StoreAPI.OrderGet(ctx, params)
	})
	api.PetPetCreateHandler = pet.PetCreateHandlerFunc(func(params pet.PetCreateParams, principal interface{}) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		ctx = storeAuth(ctx, principal)
		return c.PetAPI.PetCreate(ctx, params)
	})
	api.PetPetDeleteHandler = pet.PetDeleteHandlerFunc(func(params pet.PetDeleteParams, principal interface{}) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		ctx = storeAuth(ctx, principal)
		return c.PetAPI.PetDelete(ctx, params)
	})
	api.PetPetGetHandler = pet.PetGetHandlerFunc(func(params pet.PetGetParams, principal interface{}) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		ctx = storeAuth(ctx, principal)
		return c.PetAPI.PetGet(ctx, params)
	})
	api.PetPetListHandler = pet.PetListHandlerFunc(func(params pet.PetListParams, principal interface{}) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		ctx = storeAuth(ctx, principal)
		return c.PetAPI.PetList(ctx, params)
	})
	api.PetPetUpdateHandler = pet.PetUpdateHandlerFunc(func(params pet.PetUpdateParams, principal interface{}) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		ctx = storeAuth(ctx, principal)
		return c.PetAPI.PetUpdate(ctx, params)
	})
	api.PetPetUploadImageHandler = pet.PetUploadImageHandlerFunc(func(params pet.PetUploadImageParams, principal interface{}) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		ctx = storeAuth(ctx, principal)
		return c.PetAPI.PetUploadImage(ctx, params)
	})
	api.ServerShutdown = func() {}
	return api.Serve(c.InnerMiddleware), api, nil
}

// swaggerCopy copies the swagger json to prevent data races in runtime
func swaggerCopy(orig json.RawMessage) json.RawMessage {
	c := make(json.RawMessage, len(orig))
	copy(c, orig)
	return c
}

// authorizer is a helper function to implement the runtime.Authorizer interface.
type authorizer func(*http.Request) error

func (a authorizer) Authorize(req *http.Request, principal interface{}) error {
	if a == nil {
		return nil
	}
	ctx := storeAuth(req.Context(), principal)
	return a(req.WithContext(ctx))
}

func storeAuth(ctx context.Context, principal interface{}) context.Context {
	return context.WithValue(ctx, AuthKey, principal)
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

type serverError struct {
	error
	err_msg string
}

func (err serverError) Error() string {
	return err.err_msg
}

type UserInputError struct {
	serverError
	err_msg string
}

type MiddlewareInterface interface {
	Handle(req *http.Request, body interface{}) error
}

type Middleware struct {
	MiddlewareInterface
	err_msg string
}

func (self Middleware) Handle(req *http.Request, body interface{}) error {
	return errors.New("not a middleware")
}

type MethodCheck struct {
	Middleware
	method string
}

type JsonBodyParser struct {
	Middleware
}

type BodyCheck struct {
	Middleware
}

func (self MethodCheck) Handle(req *http.Request, body interface{}) error {
	if req.Method != self.method {
		return errors.New(self.err_msg)
	}
	return nil
}

func (self JsonBodyParser) Handle(req *http.Request, body interface{}) error {
	if body == nil {
		return nil
	}
	decoder := json.NewDecoder(req.Body)

	decoder.DisallowUnknownFields()

	if err := decoder.Decode(body); err != nil {
		return err
	}
	return nil
}

func (self BodyCheck) Handle(req *http.Request, body interface{}) error {
	iter := reflect.ValueOf(body).MapRange()
	for iter.Next() {
		val, ok := iter.Value().Interface().(*interface{})
		if !ok {
			continue
		}
		if val == nil {
			return UserInputError{
				err_msg: fmt.Sprintf("key %s missing", iter.Key()),
			}
		}
	}
	return nil
}

func MiddlewareCheckBuilder(middlewareList ...MiddlewareInterface) func(req *http.Request, body interface{}) error {
	return func(req *http.Request, body interface{}) error {
		for _, middleware := range middlewareList {
			if err := middleware.Handle(req, body); err != nil {
				return err
			}
		}
		return nil
	}
}

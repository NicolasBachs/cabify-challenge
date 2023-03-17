package router

import (
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/interface-adapters/controller"
)

type Router interface {
	Serve(port string)
	Setup(*controller.AppController) error
}

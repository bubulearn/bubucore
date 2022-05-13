package ginsrv

import (
	"github.com/bubulearn/bubucore/di"
	"github.com/gin-gonic/gin"
)

// ControllerInterface is a controller interface
type ControllerInterface interface {

	// SetContainer sets DI container to the controller
	SetContainer(ctn *di.Container)

	// GetContainer returns controller's DI container instance
	GetContainer() *di.Container

	// Init initializes the controller's actions
	Init(group *gin.RouterGroup)
}

// ControllerDft is a basic controller realization
type ControllerDft struct {
	ctn *di.Container
}

// SetContainer sets DI container to the controller
func (c *ControllerDft) SetContainer(ctn *di.Container) {
	c.ctn = ctn
}

// GetContainer returns controller's DI container instance
func (c *ControllerDft) GetContainer() *di.Container {
	return c.ctn
}

// Init initializes the controller's actions
func (c *ControllerDft) Init(_ *gin.RouterGroup) {}

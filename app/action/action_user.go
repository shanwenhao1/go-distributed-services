package action

import (
	"github.com/gin-gonic/gin"
	"go-distributed-services/domain/model"
)

// user action request parameter
type UserJsonModel struct {
	model.RequestJsonModel
	Obj model.User `json:"obj"`
}

// login request
func (this UserJsonModel) Login(c *gin.Context) {
	//
}

package action

import "gin-web/dddProject/domain/model"

type UserJsonModel struct {
	model.RequestJsonModel
	Obj model.User `json:"obj"`
}

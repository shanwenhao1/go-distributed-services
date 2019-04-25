package model

// user 模型示例, UserId为Entity 唯一标识符
type User struct {
	UserId   string      `json:"user_id" gorm:"primary_key"`
	UserName string      `json:"user_name" gorm:"column:user_name"`
	RegDate  interface{} `json:"reg_date" gorm:"column:reg_date"`
	Sex      int32       `json:"sex" gorm:"-"`
}

func (user *User) ResetName(newName string) {
	user.UserName = newName
	user.Save()
}

func (user *User) Save() {
	// 保存User的修改
}

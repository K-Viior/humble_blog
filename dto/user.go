package dto

type UserDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegisterDto struct {
	Username   string `json:"username" binding:"required,min=3"`
	Age        int    `json:"age" binding:"isdefault=18,min=0,max=100"`
	Email      string `json:"email"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password"binding:"required,eqfield=Password"`
}

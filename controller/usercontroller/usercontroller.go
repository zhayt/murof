package usercontroller

type userManipulation interface {
	CreateUser(name string) int
	GetUser(name string) bool
}

type UserController struct {
	UserManipulation userManipulation
}

func NewUserController(manipulation userManipulation) *UserController {
	return &UserController{
		UserManipulation: manipulation,
	}
}

func (c *UserController) SignIn(name string) bool {
	return c.UserManipulation.GetUser(name)
}

func (c *UserController) SignUp(name string) int {
	return c.UserManipulation.CreateUser(name)
}

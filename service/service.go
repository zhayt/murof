package service

type storage interface {
	CreateUser(name string) int
	GetUser(name string) bool
	CreatePost(idUser int, text string) int
	GetPost(userId, postID int) string
	AllPosts(userId int) []string
}

type UserService struct {
	Storage storage
}

func NewUserService(storage storage) *UserService {
	return &UserService{Storage: storage}
}

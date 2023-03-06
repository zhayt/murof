package postcontroller

type postManipulation interface {
	CreatePost(idUser int, text string) int
	GetPost(userId, postID int) string
	AllPosts(userId int) []string
}

type PostController struct {
	PostManipulation postManipulation
}

func NewPostController(manipulation postManipulation) *PostController {
	return &PostController{PostManipulation: manipulation}
}

func (p *PostController) CreatePost(idUser int, text string) int {
	return p.PostManipulation.CreatePost(idUser, text)
}

func (p *PostController) GetPost(userId, postID int) string {
	return p.PostManipulation.GetPost(userId, postID)
}

func (p *PostController) GetAllPosts(userId int) []string {
	return p.PostManipulation.AllPosts(userId)
}

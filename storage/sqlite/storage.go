package sqlite

type Storage struct {
	users []string
	posts map[int][]string
}

func NewStorage() *Storage {
	return &Storage{users: make([]string, 0, 10), posts: make(map[int][]string)}
}

func (s *Storage) CreateUser(name string) int {
	// ...
	s.users = append(s.users, name)
	return len(s.users) - 1
}

func (s *Storage) GetUser(name string) bool {
	// ...
	for _, user := range s.users {
		if name == user {
			return true
		}
	}

	return false
}

func (s *Storage) CreatePost(idUser int, text string) int {
	s.posts[idUser] = append(s.posts[idUser], text)
	return len(s.posts[idUser]) - 1
}

func (s *Storage) GetPost(userId, postID int) string {
	return s.posts[userId][postID]
}

func (s *Storage) AllPosts(userId int) []string {
	return s.posts[userId]
}

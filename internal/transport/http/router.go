package http

import "net/http"

func (s *Server) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.mid.MiddleWare(s.handler.Home))
	mux.HandleFunc("/sign-up", s.handler.SignUp)
	mux.HandleFunc("/sign-in", s.handler.SignIn)
	mux.HandleFunc("/createpost", s.mid.MiddleWare(s.handler.CreatePost))
	mux.HandleFunc("/logout", s.mid.MiddleWare(s.handler.Logout))
	mux.HandleFunc("/post/", s.mid.MiddleWare(s.handler.Post))
	mux.HandleFunc("/like-post", s.mid.MiddleWare(s.handler.LikePost))
	mux.HandleFunc("/post/change/", s.mid.MiddleWare(s.handler.ChangePost))
	mux.HandleFunc("/post/delete/", s.mid.MiddleWare(s.handler.DeletePost))
	mux.HandleFunc("/myposts", s.mid.MiddleWare(s.handler.MyPosts))
	mux.HandleFunc("/my-comment-posts", s.mid.MiddleWare(s.handler.MyCommentPosts))
	mux.HandleFunc("/my-liked-posts", s.mid.MiddleWare(s.handler.MyLikedPosts))
	mux.Handle("/ui/css/", http.StripPrefix("/ui/css/", http.FileServer(http.Dir("./ui/css/"))))

	return mux
}

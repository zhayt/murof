package http

import "net/http"

func (s *Server) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.mid.AuthorizationRequired(s.handler.Home))
	mux.HandleFunc("/sign-up", s.handler.SignUp)
	mux.HandleFunc("/sign-in", s.mid.WithSessionBlocked(s.handler.SignIn))
	mux.HandleFunc("/createpost", s.mid.AuthorizationRequired(s.handler.CreatePost))
	mux.HandleFunc("/logout", s.mid.AuthorizationRequired(s.handler.Logout))
	mux.HandleFunc("/post/", s.mid.AuthorizationRequired(s.handler.Post))
	mux.HandleFunc("/like-post", s.mid.AuthorizationRequired(s.handler.LikePost))
	mux.HandleFunc("/post/change/", s.mid.AuthorizationRequired(s.handler.ChangePost))
	mux.HandleFunc("/post/delete/", s.mid.AuthorizationRequired(s.handler.DeletePost))
	mux.HandleFunc("/myposts", s.mid.AuthorizationRequired(s.handler.MyPosts))
	mux.HandleFunc("/my-comment-posts", s.mid.AuthorizationRequired(s.handler.MyCommentPosts))
	mux.HandleFunc("/my-liked-posts", s.mid.AuthorizationRequired(s.handler.MyLikedPosts))
	mux.Handle("/ui/css/", http.StripPrefix("/ui/css/", http.FileServer(http.Dir("./ui/css/"))))

	return mux
}

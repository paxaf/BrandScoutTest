package middleware

import "net/http"

func SimpleMiddleware(all, author, post http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			post.ServeHTTP(w, r)
			return
		}
		values := r.URL.Query()
		_, ok := values["author"]
		if ok {
			author.ServeHTTP(w, r)
			return
		}
		all.ServeHTTP(w, r)
	})
}

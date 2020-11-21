package router

// func LogMiddleware(l ilog.Logger) Middleware {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			start := time.Now()

// 			sw := &statusWriter{ResponseWriter: w}
// 			next.ServeHTTP(sw, r)

// 			duration := time.Now().Sub(start)
// 			l.WithFields(
// 				"Duration", duration,
// 				"Status", sw.status,
// 				"Method", r.Method,
// 				"RequestURI", r.RequestURI,
// 				"RequestId", sw.Header().Get("RequestID"),
// 			).Info("Access log")
// 		})
// 	}
// }

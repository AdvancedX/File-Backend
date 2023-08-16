package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	v1 "kratos-realworld/api/backend/v1"
	"kratos-realworld/internal/conf"
	"kratos-realworld/internal/service"
)

func NewSkipRoutersMatcher() selector.MatchFunc {

	skipRouters := make(map[string]struct{})
	skipRouters["/backend.v1.RealWorld/Login"] = struct{}{}
	skipRouters["/backend.v1.RealWorld/Register"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := skipRouters[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, jwtc *conf.JWT, backend *service.BackendService, logger log.Logger) *http.Server {

	var opts = []http.ServerOption{
		http.ErrorEncoder(errorEncoder),
		http.Middleware(
			recovery.Recovery(),
			//selector.Server(auth.JWTAuth(jwtc.Secret)).Match(NewSkipRoutersMatcher()).Build(),
			logging.Server(logger),
		),
		http.Filter(
			handlers.CORS(
				handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}),
				handlers.AllowedOrigins([]string{"*"}),
			),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterBackendHTTPServer(srv, backend)

	r := srv.Route("/")
	r.POST("/v1/file", CreateFileHandler(backend))
	r.PUT("/v1/file", UpdateFileHandler(backend))
	r.POST("/download", DownloadFileHandler(backend))
	return srv
}

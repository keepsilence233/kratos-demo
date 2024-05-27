package server

import (
	v1 "kratos-demo/api/helloworld/v1"
	stu "kratos-demo/api/student/v1"
	"kratos-demo/internal/conf"
	"kratos-demo/internal/service"
	"kratos-demo/internal/service_gin"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	greeter *service.GreeterService,
	student *service.StudentService,
	ginSrv *service_gin.GinService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
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
	v1.RegisterGreeterHTTPServer(srv, greeter)
	stu.RegisterStudentHTTPServer(srv, student)

	// 在浏览器中直接用GET请求下载服务器的文件
	// Notice 注意 这个路由要放到注册gin路由的上面 —— GET请求
	srv.HandleFunc("/downloadFile", ginSrv.DownloadFile)
	// 注册gin http server
	srv.HandlePrefix("/", service_gin.RegisterHTTPServer(ginSrv))
	return srv
}

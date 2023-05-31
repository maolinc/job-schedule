package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"job/internal/config"
	"job/internal/jobaction"
	"job/internal/server"
	"job/internal/svc"
	"job/internal/tools/middleware"
	"job/pb"
	"net/http"
	"os"
	"strings"
)

var configFile = flag.String("f", "etc/job.yaml", "the config file")

//go:embed static
var dir embed.FS

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterJobCronServer(grpcServer, server.NewJobCronServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	//rpc log,grpc的全局拦截器
	//s.AddUnaryInterceptors(errorx.LoggerInterceptor)

	gw := gateway.MustNewServer(c.Gateway)
	addOtherRouter(gw, &c)
	gw.Use(middleware.NewLogMiddleware().Handle)

	group := service.NewServiceGroup()
	group.Add(s)
	group.Add(gw)

	// 启动job
	jobaction.StartJob(ctx)

	defer group.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	fmt.Printf("Starting gateway at %s:%d...\n", c.Gateway.Host, c.Gateway.Port)
	fmt.Printf("Starting Scheduler server ...\n")
	group.Start()
}

func addOtherRouter(server *gateway.Server, c *config.Config) {
	host := fmt.Sprintf("%s:%d", c.IndexHost, c.Gateway.Port)
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-Type", "text/html")
			var (
				err error
				bs  []byte
			)
			if c.Mode == "dev" {
				bs, err = os.ReadFile("static/index.html")
			} else {
				file, er := dir.Open("static/index.html")
				if er != nil {
					return
				}
				bs, err = io.ReadAll(file)
			}
			if err != nil {
				return
			}
			bs = []byte(strings.Replace(string(bs), "${host}", host, 1))
			_, _ = writer.Write(bs)
		},
	})
	routes := fileSystem(&dir, "static", nil)
	server.AddRoutes(routes)
}

func fileSystem(fs *embed.FS, root string, routes []rest.Route) []rest.Route {
	if routes == nil {
		routes = make([]rest.Route, 0)
	}
	readDir, err := fs.ReadDir(root)
	if err != nil {
		return routes
	}
	for _, entry := range readDir {
		fileName := fmt.Sprintf("%s/%s", root, entry.Name())
		if entry.IsDir() {
			routes = fileSystem(fs, fileName, routes)
		} else {
			routes = append(routes, rest.Route{
				Method: http.MethodGet,
				Path:   "/" + fileName,
				Handler: func(writer http.ResponseWriter, request *http.Request) {
					file, err := fs.Open(fileName)
					if err != nil {
						writer.WriteHeader(http.StatusBadRequest)
						return
					}
					all, err := io.ReadAll(file)
					if err != nil {
						writer.WriteHeader(http.StatusBadRequest)
						return
					}
					_, _ = writer.Write(all)
				},
			})
		}
	}
	return routes
}

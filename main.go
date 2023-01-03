package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	grpc_middleware_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/auth"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	handler "github.com/zmzhang8/grpc_example/handler/v1"
	"github.com/zmzhang8/grpc_example/lib/auth"
	"github.com/zmzhang8/grpc_example/lib/log"
	middleware_logging "github.com/zmzhang8/grpc_example/middleware/logging"
	middleware_recovery "github.com/zmzhang8/grpc_example/middleware/recovery"
	middleware_skip "github.com/zmzhang8/grpc_example/middleware/skip"
	middleware_trace_id "github.com/zmzhang8/grpc_example/middleware/trace_id"
	pb "github.com/zmzhang8/grpc_example/proto/v1"
)

func main() {
	var (
		debug              = flag.Bool("debug", false, "Enable debug")
		port               = flag.Int("port", 8080, "Listen port")
		mode               = flag.String("mode", "grpc", "Server mode. Value should be one of grpc, gateway, gateway-hybrid and web-hybrid.\nIf gateway or gateway-hybrid is selected, grpc-server-endpoint must also be specified.")
		grpcServerEndpoint = flag.String("grpc-server-endpoint", "", "gRPC server endpoint")
		tlsCert            = flag.String("tls_cert", "", "TLS certificate")
		tlsKey             = flag.String("tls_key", "", "TLS key")
	)
	flag.Parse()

	logger := log.NewLogger(log.NewCore(false, os.Stdout, *debug))
	defer logger.Sync()
	if *debug {
		logger.Debug("Debug enabled")
	}

	var tlsConfig *tls.Config
	if *tlsCert != "" && *tlsKey != "" {
		logger.Info("TLS enabled")
		var err error
		if tlsConfig, err = loadTlsCert(*tlsCert, *tlsKey); err != nil {
			logger.Fatalw("Failed to load TLS cert", "error", err)
		}
	}

	if *mode == "grpc" {
		grpcServer := createGrpcServer(logger, tlsConfig, *debug)
		if err := runGrpcServer(logger, grpcServer, *port); err != nil {
			logger.Fatalw("gRPC server failed to serve", "error", err)
		}
	} else if *mode == "gateway" {
		if *grpcServerEndpoint == "" {
			logger.Fatal("grpc-server-endpoint must be specified")
		}
		if err := runGatewayServer(logger, *grpcServerEndpoint, *port, tlsConfig, *debug); err != nil {
			logger.Fatalw("gRPC-Gateway server failed to serve", "error", err)
		}
	} else if *mode == "gateway-hybrid" {
		if *grpcServerEndpoint == "" {
			logger.Fatal("grpc-server-endpoint must be specified")
		}
		grpcServer := createGrpcServer(logger, tlsConfig, *debug)
		if err := runGrpcGatewayHybridServer(logger, grpcServer, *grpcServerEndpoint, *port, tlsConfig, *debug); err != nil {
			logger.Fatal("gRPC and gRPC-Gateway Hybrid server failed to serve", "error", err)
		}
	} else if *mode == "web-hybrid" {
		grpcServer := createGrpcServer(logger, tlsConfig, *debug)
		if err := runGrpcWebHybridServer(logger, grpcServer, *port, tlsConfig); err != nil {
			logger.Fatalw("gRPC and gRPC-Web hybrid server failed to serve", "error", err)
		}
	} else {
		logger.Error("Invalid mode ", *mode)
	}
}

func loadTlsCert(tlsCert, tlsKey string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}, nil
}

// Run standalone gRPC server.
func runGrpcServer(
	logger log.Logger,
	grpcServer *grpc.Server,
	port int,
) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Error("Server failed to listen at port ", port)
		return err
	}

	logger.Info("gRPC server is listening at port ", port)
	return grpcServer.Serve(listener)
}

// Run standalone gRPC-Gateway server - a reverse-proxy server which translates a RESTful HTTP API into gRPC.
// The gateway server should be used with a gRPC server.
// https://github.com/grpc-ecosystem/grpc-gateway
func runGatewayServer(
	logger log.Logger,
	grpcServerEndpoint string,
	port int,
	tlsConfig *tls.Config,
	useSwagger bool,
) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcServerTlsEnabled := tlsConfig != nil
	gatewayMux, err := createGatewayMux(logger, grpcServerEndpoint, grpcServerTlsEnabled, ctx)
	if err != nil {
		logger.Error("Failed to create gateway mux")
		return err
	}

	var httpHandler http.Handler = gatewayMux
	if useSwagger {
		fileServer := http.FileServer(http.Dir("./third_party/swagger_ui"))
		httpHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/swagger/") {
				http.StripPrefix("/swagger/", fileServer).ServeHTTP(w, r)
			} else {
				gatewayMux.ServeHTTP(w, r)
			}
		})
	}

	logger.Info("gRPC-Gateway server is listening at port ", port)
	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   httpHandler,
		TLSConfig: tlsConfig,
	}
	if tlsConfig != nil {
		return server.ListenAndServeTLS("", "")
	} else {
		return server.ListenAndServe()
	}
}

// Run gRPC server and gRPC-Gateway server together on the same port using mux.
// https://github.com/philips/grpc-gateway-example
func runGrpcGatewayHybridServer(
	logger log.Logger,
	grpcServer *grpc.Server,
	grpcServerEndpoint string,
	port int,
	tlsConfig *tls.Config,
	useSwagger bool,
) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcServerTlsEnabled := tlsConfig != nil
	gatewayMux, err := createGatewayMux(logger, grpcServerEndpoint, grpcServerTlsEnabled, ctx)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", gatewayMux)

	if useSwagger {
		fileServer := http.FileServer(http.Dir("./third_party/swagger_ui"))
		mux.Handle("/swagger/", http.StripPrefix("/swagger/", fileServer))
	}

	httpHandler := func(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
				grpcServer.ServeHTTP(w, r)
			} else {
				httpHandler.ServeHTTP(w, r)
			}
		})
	}(grpcServer, mux)

	// https://stackoverflow.com/questions/69542087/why-am-i-getting-connection-connection-closed-before-server-preface-received-in
	httpHandler = h2c.NewHandler(httpHandler, &http2.Server{})

	logger.Info("gRPC and gRPC-Gateway Hybrid server is listening at port ", port)
	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   httpHandler,
		TLSConfig: tlsConfig,
	}
	if tlsConfig != nil {
		return server.ListenAndServeTLS("", "")
	} else {
		return server.ListenAndServe()
	}
}

// Run gRPC server and gRPC-Web server together on the same port using mux.
// Note that this server only supports unary calls and server-side streams.
// https://pkg.go.dev/github.com/improbable-eng/grpc-web/go/grpcweb
func runGrpcWebHybridServer(
	logger log.Logger,
	grpcServer *grpc.Server,
	port int,
	tlsConfig *tls.Config,
) error {
	grpcWebServer := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true // allow all origins
		}),
	)

	mux := http.NewServeMux()
	mux.Handle("/", grpcWebServer)

	httpHandler := func(wrappedGrpcServer *grpcweb.WrappedGrpcServer, httpHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if wrappedGrpcServer.IsGrpcWebRequest(r) {
				// handle gRPC-Web requests
				wrappedGrpcServer.ServeHTTP(w, r)
			} else if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
				// handle regular gRPC requests
				wrappedGrpcServer.ServeHTTP(w, r)
			} else {
				httpHandler.ServeHTTP(w, r)
			}
		})
	}(grpcWebServer, mux)
	// https://stackoverflow.com/questions/69542087/why-am-i-getting-connection-connection-closed-before-server-preface-received-in
	httpHandler = h2c.NewHandler(httpHandler, &http2.Server{})

	logger.Info("gRPC and gRPC-Web Hybrid server is listening at port ", port)
	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   httpHandler,
		TLSConfig: tlsConfig,
	}
	if tlsConfig != nil {
		return server.ListenAndServeTLS("", "")
	} else {
		return server.ListenAndServe()
	}
}

func createGrpcServer(
	logger log.Logger,
	tlsConfig *tls.Config,
	enableReflection bool,
) *grpc.Server {
	var credsOption grpc.ServerOption = grpc.EmptyServerOption{}
	if tlsConfig != nil {
		credsOption = grpc.Creds(credentials.NewTLS(tlsConfig))
	}

	loggerFunc := func(ctx context.Context, logger log.Logger) log.Logger {
		return logger.With("trace-id", middleware_trace_id.MustGetTraceID(ctx))
	}
	skipAuthFunc := func(ctx context.Context, service string, method string) bool {
		return service == grpc_reflection_v1alpha.ServerReflection_ServiceDesc.ServiceName ||
			service == grpc_health_v1.Health_ServiceDesc.ServiceName
	}
	server := grpc.NewServer(
		credsOption,
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			middleware_trace_id.StreamServerInterceptor(),
			middleware_logging.StreamServerInterceptor(logger, loggerFunc),
			middleware_recovery.StreamServerInterceptor(logger),
			middleware_skip.StreamServerInterceptor(
				grpc_middleware_auth.StreamServerInterceptor(auth.RejectAll),
				skipAuthFunc,
			),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			middleware_trace_id.UnaryServerInterceptor(),
			middleware_logging.UnaryServerInterceptor(logger, loggerFunc),
			middleware_recovery.UnaryServerInterceptor(logger),
			middleware_skip.UnaryServerInterceptor(
				grpc_middleware_auth.UnaryServerInterceptor(auth.RejectAll),
				skipAuthFunc,
			),
		)),
	)

	// Register reflection service
	if enableReflection {
		reflection.Register(server)
	}

	// Register health service
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// Register custom services
	pb.RegisterHealthServer(server, handler.NewHealthServer())
	pb.RegisterGreeterServer(server, handler.NewGreeterServer())
	pb.RegisterRouteGuideServer(server, handler.NewRouteGuideServer())
	pb.RegisterAccountServer(server, handler.NewAccountServer())

	return server
}

func createGatewayMux(
	logger log.Logger,
	grpcServerEndpoint string,
	grpcServerTlsEnabled bool,
	ctx context.Context,
) (*runtime.ServeMux, error) {
	credsOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	if grpcServerTlsEnabled {
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: false,
		})
		credsOption = grpc.WithTransportCredentials(creds)
	}

	clientConn, err := grpc.DialContext(ctx, grpcServerEndpoint, credsOption)
	if err != nil {
		logger.Error("Failed to dail ", grpcServerEndpoint)
		return nil, err
	}

	gatewayMux := runtime.NewServeMux()
	for _, f := range []func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error{
		pb.RegisterHealthHandler,
		pb.RegisterGreeterHandler,
		pb.RegisterRouteGuideHandler,
		pb.RegisterAccountHandler,
	} {
		if err := f(ctx, gatewayMux, clientConn); err != nil {
			return nil, err
		}
	}

	return gatewayMux, nil
}

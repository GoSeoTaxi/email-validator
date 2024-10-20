package main

import (
	"context"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/GoSeoTaxi/email-validator/internal/pb"
)

func main() {

	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()
	sugar := logger.Sugar()

	grpcServerEndpoint := getEnv("GRPC_SERVER_ENDPOINT", "email-validator:50051")
	httpPort := getEnv("HTTP_PORT", "8080")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaller runtime.Marshaler,
			w http.ResponseWriter, req *http.Request, err error) {
			sugar.Errorf("Error handling request %s %s: %v", req.Method, req.URL.Path, err)
			s, ok := status.FromError(err)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var httpStatus int
			switch s.Code() {
			case codes.Unavailable:
				httpStatus = http.StatusServiceUnavailable
			case codes.DeadlineExceeded:
				httpStatus = http.StatusGatewayTimeout
			default:
				httpStatus = runtime.HTTPStatusFromCode(s.Code())
			}
			w.WriteHeader(httpStatus)
		}),
	)

	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := pb.RegisterEmailValidatorHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		sugar.Fatalf("Failed to register service: %v", err)
	}

	sugar.Infof("Starting an HTTP server on port %s", httpPort)
	if err := http.ListenAndServe(":"+httpPort, mux); err != nil {
		sugar.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

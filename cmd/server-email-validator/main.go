package main

import (
	"context"
	"log"
	"net"

	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/GoSeoTaxi/email-validator/internal/adapter"
	"github.com/GoSeoTaxi/email-validator/internal/config"
	"github.com/GoSeoTaxi/email-validator/internal/domain"
	pb "github.com/GoSeoTaxi/email-validator/internal/pb"
	"github.com/GoSeoTaxi/email-validator/internal/service"
	"github.com/GoSeoTaxi/email-validator/internal/usecase"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.LoadConfig,
			NewRedisClient,
			NewRedisCache,
			NewDNSResolver,
			NewGRPCServer,
			usecase.NewEmailValidator,
			service.NewEmailValidatorService,
		),
		fx.Invoke(RegisterGRPCServer),
	)
	app.Run()
}

func NewRedisCache(client *redis.Client) domain.Cache {
	return adapter.NewRedisCache(client)
}

func NewDNSResolver(cfg *config.Config) domain.DNSResolver {
	netResolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, network, cfg.DNSHosts[0]+":53")
		},
	}
	return adapter.NewNetDNSResolver(netResolver)
}

func NewRedisClient(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisHost + ":" + cfg.RedisPort,
		DB:   cfg.RedisDB,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	err = client.ConfigSet(ctx, "maxmemory-policy", "allkeys-lru").Err()
	if err != nil {
		log.Fatalf("Failed to set maxmemory-policy: %v", err)
	}
	err = client.ConfigSet(ctx, "maxmemory", cfg.RedisMaxMemory).Err()
	if err != nil {
		log.Fatalf("Failed to set maxmemory: %v", err)
	}

	return client
}

func RegisterGRPCServer(lc fx.Lifecycle, cfg *config.Config, srv *grpc.Server, emailValidatorService *service.EmailValidatorService) {
	pb.RegisterEmailValidatorServer(srv, emailValidatorService)

	reflection.Register(srv)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Не удалось начать прослушивание: %v", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Запуск gRPC сервера на порту", cfg.GRPCPort)
			go func() {
				if err := srv.Serve(lis); err != nil {
					log.Fatalf("gRPC сервер завершился с ошибкой: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Остановка gRPC сервера")
			srv.GracefulStop()
			return nil
		},
	})
}

func NewGRPCServer() *grpc.Server {
	return grpc.NewServer()
}

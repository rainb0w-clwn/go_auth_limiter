package commands

import (
	"context"
	"log"
	"time"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/server/grpc/limiter"
	proto "github.com/rainb0w-clwn/go_auth_limiter/proto/limiter"
	"github.com/spf13/cobra"
)

var resetBucketCmd = &cobra.Command{
	Use:   "reset_bucket [login] [ip]",
	Short: "Сброс бакета",
	Run: func(_ *cobra.Command, args []string) {
		login := args[0]
		ip := args[0]

		log.Println("clear_bucket,", "login:", login, " ip:", ip)

		grpcClient, err := limiter.NewClient(cfg.GRPC.Host, cfg.GRPC.Port)
		if err != nil {
			log.Fatalf("failed to create gRPC client: %v", err)
		}
		defer grpcClient.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		ok, err := grpcClient.BucketReset(ctx, &proto.BucketResetRequest{
			Login: login,
			Ip:    ip,
		})
		if err != nil {
			log.Printf("BucketReset error: %v", err)
		} else {
			log.Printf("BucketReset success: %v", ok)
		}
	},
}

func init() {
	rootCmd.AddCommand(resetBucketCmd)
}

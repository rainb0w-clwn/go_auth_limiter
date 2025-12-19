package commands

import (
	"context"
	"log"
	"time"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/server/grpc/limiter"
	proto "github.com/rainb0w-clwn/go_auth_limiter/proto/limiter"
	"github.com/spf13/cobra"
)

var deleteCidrFromBlackListCmd = &cobra.Command{
	Use:   "delete_cidr_from_black_list [cidr]",
	Short: "Удалить подсеть из черного списка",
	Run: func(_ *cobra.Command, args []string) {
		cidr := args[0]

		grpcClient, err := limiter.NewClient(cfg.GRPC.Host, cfg.GRPC.Port)
		if err != nil {
			log.Fatalf("failed to create gRPC client: %v", err)
		}
		defer grpcClient.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		ok, err := grpcClient.BlackListDelete(ctx, &proto.BlackListDeleteRequest{
			IpNet: cidr,
		})
		if err != nil {
			log.Printf("BlackListDelete error: %v", err)
		} else {
			log.Printf("BlackListDelete success: %v", ok)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCidrFromBlackListCmd)
}

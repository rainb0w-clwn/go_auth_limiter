package commands

import (
	"context"
	"log"
	"time"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/server/grpc/limiter"
	proto "github.com/rainb0w-clwn/go_auth_limiter/proto/limiter"
	"github.com/spf13/cobra"
)

var addCidrToBlackListCmd = &cobra.Command{
	Use:   "add_cidr_to_black_list [cidr]",
	Short: "Добавить подсеть в черный список",
	Run: func(_ *cobra.Command, args []string) {
		cidr := args[0]

		grpcClient, err := limiter.NewClient(cfg.GRPC.Host, cfg.GRPC.Port)
		if err != nil {
			log.Fatalf("failed to create gRPC client: %v", err)
		}
		defer grpcClient.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		ok, err := grpcClient.BlackListAdd(ctx, &proto.BlackListAddRequest{
			IpNet: cidr,
		})
		if err != nil {
			log.Printf("BlackListAdd error: %v", err)
		} else {
			log.Printf("BlackListAdd success: %v", ok)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCidrToBlackListCmd)
}

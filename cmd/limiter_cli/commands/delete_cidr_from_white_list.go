package commands

import (
	"context"
	"log"
	"time"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/server/grpc/limiter"
	proto "github.com/rainb0w-clwn/go_auth_limiter/proto/limiter"
	"github.com/spf13/cobra"
)

var deleteCidrFromWhiteListCmd = &cobra.Command{
	Use:   "delete_cidr_from_white_list [cidr]",
	Short: "Удалить подсеть из белого списка",
	Run: func(_ *cobra.Command, args []string) {
		cidr := args[0]

		grpcClient, err := limiter.NewClient(cfg.GRPC.Host, cfg.GRPC.Port)
		if err != nil {
			log.Fatalf("failed to create gRPC client: %v", err)
		}
		defer grpcClient.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		ok, err := grpcClient.WhiteListDelete(ctx, &proto.WhiteListDeleteRequest{
			IpNet: cidr,
		})
		if err != nil {
			log.Printf("WhiteListDelete error: %v", err)
		} else {
			log.Printf("WhiteListDelete success: %v", ok)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCidrFromWhiteListCmd)
}

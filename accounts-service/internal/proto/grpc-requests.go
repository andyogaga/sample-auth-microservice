package requests

import (
	context "context"
	"fmt"

	services "accounts-service/internal/services"
)

var walletService services.WalletService

func SetupService(_walletService *services.WalletService) {
	walletService = *_walletService
}

type WalletsServer struct {
	UnimplementedWalletServiceServer
}

func (c *WalletsServer) CreateWallet(ctx context.Context, req *CreateWalletRequest) (*CreateWalletResponse, error) {
	fmt.Println(req)
	walletService.GetWallet(123)
	return nil, nil
}

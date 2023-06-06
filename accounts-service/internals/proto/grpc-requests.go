package requests

import (
	context "context"

	services "accounts-service/internals/services"
)

var walletService services.WalletService

func SetupService(_walletService *services.WalletService) {
	walletService = *_walletService
}

type WalletsServer struct {
	UnimplementedWalletServiceServer
}

func (c *WalletsServer) CreateWallet(ctx context.Context, req *CreateWalletRequest) (*CreateWalletResponse, error) {
	walletService.GetWallet(123)
	return nil, nil
}

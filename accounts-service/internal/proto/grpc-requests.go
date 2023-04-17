package requests

import (
	"accounts-service/internal/dto"
	context "context"
)

type WalletsServer struct {
	UnimplementedWalletServiceServer
	CreateWalletModel dto.CreateWallet
}

func (l *WalletsServer) CreateWallet(ctx context.Context, req *CreateWalletRequest) (*CreateWalletResponse, error) {

}

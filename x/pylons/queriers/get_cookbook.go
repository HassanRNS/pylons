package queriers

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/MikeSofaer/pylons/x/pylons/keep"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the nameservice Querier
const (
	KeyGetCookbook = "get_cookbook"
)

// GetCookbook returns a cookbook based on the cookbook id
func GetCookbook(ctx sdk.Context, path []string, req abci.RequestQuery, keeper keep.Keeper) ([]byte, sdk.Error) {
	cookbookID := path[0]
	cookbook, err := keeper.GetCookbook(ctx, cookbookID)

	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	// if we cannot find the value then it should return an error
	bz, err := codec.MarshalJSONIndent(keeper.Cdc, cookbook)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	return bz, nil

}
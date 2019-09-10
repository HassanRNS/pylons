package handlers

import (
	"encoding/json"

	"github.com/MikeSofaer/pylons/x/pylons/keep"
	"github.com/MikeSofaer/pylons/x/pylons/msgs"
	"github.com/MikeSofaer/pylons/x/pylons/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandlerMsgCreateRecipe is used to create cookbook by a developer
func HandlerMsgCreateRecipe(ctx sdk.Context, keeper keep.Keeper, msg msgs.MsgCreateRecipe) sdk.Result {

	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}
	cook, err2 := keeper.GetCookbook(ctx, msg.CookbookName)
	if !cook.Sender.Equals(msg.Sender) {
		return sdk.ErrUnauthorized("cookbook not owned by the sender").Result()
	}

	recipe := types.NewRecipe(msg.RecipeName, msg.CookbookName, msg.Description,
		msg.CoinInputs, msg.CoinOutputs, msg.ItemInputs, msg.ItemOutputs,
		0, msg.Sender)
	if err := keeper.SetRecipe(ctx, recipe); err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}

	mRecipe, err2 := json.Marshal(map[string]string{
		"RecipeID": recipe.ID,
	})

	if err != nil {
		return sdk.ErrInternal(err2.Error()).Result()
	}

	return sdk.Result{Data: mRecipe}
}
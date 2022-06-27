package app_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/client/cli"

	"github.com/Pylons-tech/pylons/testutil/network"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
)

const (
	errCode     = uint32(19)
	successCode = uint32(0)
)

func TestPreventSpamTx(t *testing.T) {
	// Set MaxTxsInBlock  = 2
	numberTxsinBlocks := 2
	config := network.ConfigWithMaxTxsInBlock(uint64(numberTxsinBlocks))
	net := network.New(t, config)

	val := net.Validators[0]
	ctx := val.ClientCtx

	res := make([]testutil.BufferWriter, numberTxsinBlocks+1)
	for i := 0; i < numberTxsinBlocks+1; i++ {
		res[i] = executeSendTx(ctx, net, net.Validators[0].Address.String())
	}
	net.WaitForNextBlock()
	var resp sdk.TxResponse
	require.NoError(t, ctx.Codec.UnmarshalJSON(res[numberTxsinBlocks].Bytes(), &resp))
	require.Equal(t, errCode, resp.Code)
}

func TestOtherTransactionIsValid(t *testing.T) {
	// Set MaxTxsInBlock  = 2
	numberTxsinBlocks := 2
	config := network.ConfigWithMaxTxsInBlock(uint64(numberTxsinBlocks))
	net := network.New(t, config)

	val := net.Validators[0]
	ctx := val.ClientCtx
	addr := GenerateAddressesInKeyring(ctx.Keyring, 1)

	//init balances
	executeOtherSendTx(ctx, net, net.Validators[0].Address.String(), addr[0].String())
	net.WaitForNextBlock()

	res := make([]testutil.BufferWriter, numberTxsinBlocks+1)
	for i := 0; i < numberTxsinBlocks+1; i++ {
		res[i] = executeSendTx(ctx, net, net.Validators[0].Address.String())
	}
	// other transaction by addr[0].String()
	validRes := executeOtherSendTx(ctx, net, addr[0].String(), addr[0].String())
	net.WaitForNextBlock()

	var resp sdk.TxResponse
	require.NoError(t, ctx.Codec.UnmarshalJSON(res[numberTxsinBlocks].Bytes(), &resp))
	require.Equal(t, errCode, resp.Code)

	require.NoError(t, ctx.Codec.UnmarshalJSON(validRes.Bytes(), &resp))
	require.Equal(t, successCode, resp.Code)

}
func TestPreventSpamTxWithTwoUsers(t *testing.T) {
	numberTxsinBlocks := 2
	config := network.ConfigWithMaxTxsInBlock(uint64(numberTxsinBlocks))
	net := network.New(t, config)
	val := net.Validators[0]
	ctx := val.ClientCtx
	address, _ := GenerateAddressWithAccount(ctx, t, net)
	net.WaitForNextBlock()
	OutRes := make([]testutil.BufferWriter, numberTxsinBlocks+1)
	for i := 0; i < numberTxsinBlocks+1; i++ {
		OutRes[i] = executeSendTxCMD(ctx, net.Validators[0].Address.String(), net)
		fmt.Println(OutRes[i])
	}
	res := executeSendTxCMD(ctx, address, net)
	net.WaitForNextBlock()
	var resp sdk.TxResponse
	require.NoError(t, ctx.Codec.UnmarshalJSON(OutRes[numberTxsinBlocks].Bytes(), &resp))
	require.Equal(t, errCode, resp.Code)

	require.NoError(t, ctx.Codec.UnmarshalJSON(res.Bytes(), &resp))
	require.Equal(t, successCode, resp.Code)

}

func GenerateAddressWithAccount(ctx client.Context, t *testing.T, net *network.Network) (string, error) {
	accs := GenerateAddressesInKeyring(ctx.Keyring, 1)
	common := CommonArgs(net)

	args := []string{net.Validators[0].Address.String(), accs[0].String(), "1000node0token"}
	args = append(args, common...)
	out, err := clitestutil.ExecTestCLICmd(ctx, bank.NewSendTxCmd(), args)
	if err != nil {
		return "", err
	}
	var resp sdk.TxResponse
	require.NoError(t, ctx.Codec.UnmarshalJSON(out.Bytes(), &resp))
	if uint32(0) != resp.Code {
		return "", fmt.Errorf("Error Code Not Success")
	}
	return accs[0].String(), nil
}

func CommonArgs(net *network.Network) []string {
	return []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastAsync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdk.NewInt(10))).String()),
	}
}

func executeSendTx(ctx client.Context, net *network.Network, addressSent string) testutil.BufferWriter {
	common := CommonArgs(net)
	args := []string{addressSent, net.Validators[0].Address.String(), "10stake"}
	args = append(args, common...)
	res, _ := clitestutil.ExecTestCLICmd(ctx, bank.NewSendTxCmd(), args)
	return res
}

// tx send coin from addressSent to
func executeOtherSendTx(ctx client.Context, net *network.Network, addressSent string, addressReceive string) testutil.BufferWriter {
	common := CommonArgs(net)
	args := []string{addressSent, addressReceive, "10stake"}
	args = append(args, common...)
	res, _ := clitestutil.ExecTestCLICmd(ctx, bank.NewSendTxCmd(), args)
	return res
}

func GenerateAddressesInKeyring(ring keyring.Keyring, n int) []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, n)
	for i := 0; i < n; i++ {
		info, _, _ := ring.NewMnemonic("NewUser"+strconv.Itoa(i), keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
		addrs[i] = info.GetAddress()
	}
	return addrs
}

func executeSendTxCMD(ctx client.Context, address string, net *network.Network) testutil.BufferWriter {
	common := CommonArgs(net)
	args := []string{address, address, "10stake"}
	args = append(args, common...)
	OutRes, _ := clitestutil.ExecTestCLICmd(ctx, bank.NewSendTxCmd(), args)
	return OutRes
}

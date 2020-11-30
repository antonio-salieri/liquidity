package liquidity

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/liquidity/x/liquidity/keeper"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// TODO: migrate to msg_server after rebase latest sdk 0.40.0 on milestone2
// TODO: emit events codes in milestone
// NewHandler returns a handler for all "liquidity" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateLiquidityPool:
			return handleMsgCreateLiquidityPool(ctx, k, msg)
		case *types.MsgDepositToLiquidityPool:
			return handleMsgDepositToLiquidityPool(ctx, k, msg)
		case *types.MsgWithdrawFromLiquidityPool:
			return handleMsgWithdrawFromLiquidityPool(ctx, k, msg)
		case *types.MsgSwap:
			return handleMsgSwap(ctx, k, msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func handleMsgCreateLiquidityPool(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateLiquidityPool) (*sdk.Result, error) {
	k.CreateLiquidityPool(ctx, msg)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			//types.EventTypeCreateLiquidityPool,
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.PoolCreatorAddress),
			sdk.NewAttribute(types.AttributeValueLiquidityPoolId, ""),
			sdk.NewAttribute(types.AttributeValueLiquidityPoolTypeIndex, fmt.Sprintf("%d", msg.PoolTypeIndex)),
			sdk.NewAttribute(types.AttributeValueReserveCoinDenoms, ""),
			sdk.NewAttribute(types.AttributeValueReserveAccount, ""),
			sdk.NewAttribute(types.AttributeValuePoolCoinDenom, ""),
			sdk.NewAttribute(types.AttributeValueSwapFeeRate, ""),
			sdk.NewAttribute(types.AttributeValueLiquidityPoolFeeRate, ""),
			sdk.NewAttribute(types.AttributeValueBatchSize, ""),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

func handleMsgDepositToLiquidityPool(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDepositToLiquidityPool) (*sdk.Result, error) {
	k.DepositLiquidityPoolToBatch(ctx, msg)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			//types.EventTypeDepositToLiquidityPoolToBatch,
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DepositorAddress),
			sdk.NewAttribute(types.AttributeValueBatchID, ""),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

func handleMsgWithdrawFromLiquidityPool(ctx sdk.Context, k keeper.Keeper, msg *types.MsgWithdrawFromLiquidityPool) (*sdk.Result, error) {
	k.WithdrawLiquidityPoolToBatch(ctx, msg)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			//types.EventTypeWithdrrawFromLiquidityPoolToBatch,
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.WithdrawerAddress),
			sdk.NewAttribute(types.AttributeValueBatchID, ""),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

func handleMsgSwap(ctx sdk.Context, k keeper.Keeper, msg *types.MsgSwap) (*sdk.Result, error) {
	// TODO: OfferCoinFee
	params := k.GetParams(ctx)

	// calculate reserve amount for OfferCoinFee
	offerCoinFeeReserve := sdk.NewCoin(msg.OfferCoin.Denom, msg.OfferCoin.Amount.ToDec().Mul(params.SwapFeeRate).TruncateInt())
	msg.OfferCoinFee = offerCoinFeeReserve
	// TODO EqualAlmost


	_, err := k.SwapLiquidityPoolToBatch(ctx, msg)
	if err != nil {
		return &sdk.Result{
			Events: ctx.EventManager().ABCIEvents(),
		}, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			//types.EventTypeSwapToBatch,
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SwapRequesterAddress),
			sdk.NewAttribute(types.AttributeValueBatchID, ""),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

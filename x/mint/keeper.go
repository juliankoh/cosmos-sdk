package mint

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// keeper of the staking store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramSpace params.Subspace
	sk         StakingKeeper
	fck        FeeCollectionKeeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey,
	paramSpace params.Subspace, sk StakingKeeper, fck FeeCollectionKeeper) Keeper {

	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		paramSpace: paramSpace.WithKeyTable(ParamKeyTable()),
		sk:         sk,
		fck:        fck,
	}
	return keeper
}

//______________________________________________________________________

// get the minter
func (k Keeper) GetMinter(ctx sdk.Context) (minter Minter) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(MinterKey)
	if b == nil {
		panic("stored minter should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &minter)
	return
}

// set the minter
func (k Keeper) SetMinter(ctx sdk.Context, minter Minter) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(minter)
	store.Set(MinterKey, b)
}

//______________________________________________________________________

// GetParams returns the total set of minting parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

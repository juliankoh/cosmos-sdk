package bank

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank/tags"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var _ Keeper = (*BaseKeeper)(nil)

// Keeper defines a module interface that facilitates the transfer of coins
// between accounts.
type Keeper interface {
	SendKeeper

	SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error)
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error)
	InputOutputCoins(ctx sdk.Context, inputs []types.Input, outputs []types.Output) (sdk.Tags, sdk.Error)

	DelegateCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
	UndelegateCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
}

// BaseKeeper manages transfers between accounts. It implements the Keeper interface.
type BaseKeeper struct {
	BaseSendKeeper

	ak         auth.AccountKeeper
	paramSpace params.Subspace
}

// NewBaseKeeper returns a new BaseKeeper
func NewBaseKeeper(ak auth.AccountKeeper,
	paramSpace params.Subspace,
	codespace sdk.CodespaceType) BaseKeeper {

	ps := paramSpace.WithKeyTable(types.ParamKeyTable())
	return BaseKeeper{
		BaseSendKeeper: NewBaseSendKeeper(ak, ps, codespace),
		ak:             ak,
		paramSpace:     ps,
	}
}

// SetCoins sets the coins at the addr.
func (keeper BaseKeeper) SetCoins(
	ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins,
) sdk.Error {

	if !amt.IsValid() {
		return sdk.ErrInvalidCoins(amt.String())
	}
	return setCoins(ctx, keeper.ak, addr, amt)
}

// SubtractCoins subtracts amt from the coins at the addr.
func (keeper BaseKeeper) SubtractCoins(
	ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins,
) (sdk.Coins, sdk.Error) {

	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}
	return subtractCoins(ctx, keeper.ak, addr, amt)
}

// AddCoins adds amt to the coins at the addr.
func (keeper BaseKeeper) AddCoins(
	ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins,
) (sdk.Coins, sdk.Error) {

	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}
	return addCoins(ctx, keeper.ak, addr, amt)
}

// InputOutputCoins handles a list of inputs and outputs
func (keeper BaseKeeper) InputOutputCoins(
	ctx sdk.Context, inputs []types.Input, outputs []types.Output,
) (sdk.Tags, sdk.Error) {

	return inputOutputCoins(ctx, keeper.ak, inputs, outputs)
}

// DelegateCoins performs delegation by deducting amt coins from an account with
// address addr. For vesting accounts, delegations amounts are tracked for both
// vesting and vested coins.
func (keeper BaseKeeper) DelegateCoins(
	ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins,
) (sdk.Tags, sdk.Error) {

	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}
	return delegateCoins(ctx, keeper.ak, addr, amt)
}

// UndelegateCoins performs undelegation by crediting amt coins to an account with
// address addr. For vesting accounts, undelegation amounts are tracked for both
// vesting and vested coins.
// If any of the undelegation amounts are negative, an error is returned.
func (keeper BaseKeeper) UndelegateCoins(
	ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins,
) (sdk.Tags, sdk.Error) {

	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}
	return undelegateCoins(ctx, keeper.ak, addr, amt)
}

// SendKeeper defines a module interface that facilitates the transfer of coins
// between accounts without the possibility of creating coins.
type SendKeeper interface {
	ViewKeeper

	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error

	GetSendEnabled(ctx sdk.Context) bool
	SetSendEnabled(ctx sdk.Context, enabled bool)
}

var _ SendKeeper = (*BaseSendKeeper)(nil)

// BaseSendKeeper only allows transfers between accounts without the possibility of
// creating coins. It implements the SendKeeper interface.
type BaseSendKeeper struct {
	BaseViewKeeper

	ak         auth.AccountKeeper
	paramSpace params.Subspace
}

// NewBaseSendKeeper returns a new BaseSendKeeper.
func NewBaseSendKeeper(ak auth.AccountKeeper,
	paramSpace params.Subspace, codespace sdk.CodespaceType) BaseSendKeeper {

	return BaseSendKeeper{
		BaseViewKeeper: NewBaseViewKeeper(ak, codespace),
		ak:             ak,
		paramSpace:     paramSpace,
	}
}

// SendCoins moves coins from one account to another
func (keeper BaseSendKeeper) SendCoins(
	ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins,
) sdk.Error {

	if !amt.IsValid() {
		return sdk.ErrInvalidCoins(amt.String())
	}
	return sendCoins(ctx, keeper.ak, fromAddr, toAddr, amt)
}

// GetSendEnabled returns the current SendEnabled
// nolint: errcheck
func (keeper BaseSendKeeper) GetSendEnabled(ctx sdk.Context) bool {
	var enabled bool
	keeper.paramSpace.Get(ctx, types.ParamStoreKeySendEnabled, &enabled)
	return enabled
}

// SetSendEnabled sets the send enabled
func (keeper BaseSendKeeper) SetSendEnabled(ctx sdk.Context, enabled bool) {
	keeper.paramSpace.Set(ctx, types.ParamStoreKeySendEnabled, &enabled)
}

var _ ViewKeeper = (*BaseViewKeeper)(nil)

// ViewKeeper defines a module interface that facilitates read only access to
// account balances.
type ViewKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool

	Codespace() sdk.CodespaceType
}

// BaseViewKeeper implements a read only keeper implementation of ViewKeeper.
type BaseViewKeeper struct {
	ak        auth.AccountKeeper
	codespace sdk.CodespaceType
}

// NewBaseViewKeeper returns a new BaseViewKeeper.
func NewBaseViewKeeper(ak auth.AccountKeeper, codespace sdk.CodespaceType) BaseViewKeeper {
	return BaseViewKeeper{ak: ak, codespace: codespace}
}

// GetCoins returns the coins at the addr.
func (keeper BaseViewKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return getCoins(ctx, keeper.ak, addr)
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper BaseViewKeeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return hasCoins(ctx, keeper.ak, addr, amt)
}

// Codespace returns the keeper's codespace.
func (keeper BaseViewKeeper) Codespace() sdk.CodespaceType {
	return keeper.codespace
}

func getCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress) sdk.Coins {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.NewCoins()
	}
	return acc.GetCoins()
}

func setCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	if !amt.IsValid() {
		return sdk.ErrInvalidCoins(amt.String())
	}
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		acc = am.NewAccountWithAddress(ctx, addr)
	}
	err := acc.SetCoins(amt)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	am.SetAccount(ctx, acc)
	return nil
}

// HasCoins returns whether or not an account has at least amt coins.
func hasCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) bool {
	return getCoins(ctx, am, addr).IsAllGTE(amt)
}

func getAccount(ctx sdk.Context, ak auth.AccountKeeper, addr sdk.AccAddress) auth.Account {
	return ak.GetAccount(ctx, addr)
}

func setAccount(ctx sdk.Context, ak auth.AccountKeeper, acc auth.Account) {
	ak.SetAccount(ctx, acc)
}

// subtractCoins subtracts amt coins from an account with the given address addr.
//
// CONTRACT: If the account is a vesting account, the amount has to be spendable.
func subtractCoins(ctx sdk.Context, ak auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error) {

	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}

	oldCoins, spendableCoins := sdk.NewCoins(), sdk.NewCoins()

	acc := getAccount(ctx, ak, addr)
	if acc != nil {
		oldCoins = acc.GetCoins()
		spendableCoins = acc.SpendableCoins(ctx.BlockHeader().Time)
	}

	// For non-vesting accounts, spendable coins will simply be the original coins.
	// So the check here is sufficient instead of subtracting from oldCoins.
	_, hasNeg := spendableCoins.SafeSub(amt)
	if hasNeg {
		return amt, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", spendableCoins, amt),
		)
	}

	newCoins := oldCoins.Sub(amt) // should not panic as spendable coins was already checked
	err := setCoins(ctx, ak, addr, newCoins)

	return newCoins, err
}

// AddCoins adds amt to the coins at the addr.
func addCoins(ctx sdk.Context, am auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error) {

	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}

	oldCoins := getCoins(ctx, am, addr)
	newCoins := oldCoins.Add(amt)

	if newCoins.IsAnyNegative() {
		return amt, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", oldCoins, amt),
		)
	}

	err := setCoins(ctx, am, addr, newCoins)

	return newCoins, err
}

// SendCoins moves coins from one account to another
// Returns ErrInvalidCoins if amt is invalid.
func sendCoins(ctx sdk.Context, am auth.AccountKeeper, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	// Safety check ensuring that when sending coins the keeper must maintain the
	if !amt.IsValid() {
		return sdk.ErrInvalidCoins(amt.String())
	}

	_, err := subtractCoins(ctx, am, fromAddr, amt)
	if err != nil {
		return err
	}

	_, err = addCoins(ctx, am, toAddr, amt)
	if err != nil {
		return err
	}

	return nil
}

// InputOutputCoins handles a list of inputs and outputs
// NOTE: Make sure to revert state changes from tx on error
func inputOutputCoins(ctx sdk.Context, am auth.AccountKeeper, inputs []types.Input, outputs []types.Output) (sdk.Tags, sdk.Error) {
	// Safety check ensuring that when sending coins the keeper must maintain the
	// Check supply invariant and validity of Coins.
	if err := types.ValidateInputsOutputs(inputs, outputs); err != nil {
		return nil, err
	}

	allTags := sdk.EmptyTags()

	for _, in := range inputs {
		_, err := subtractCoins(ctx, am, in.Address, in.Coins)
		if err != nil {
			return nil, err
		}

		allTags = allTags.AppendTag(
			tags.Sender, in.Address.String(),
		)
	}

	for _, out := range outputs {
		_, err := addCoins(ctx, am, out.Address, out.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTag(
			tags.Recipient, out.Address.String(),
		)
	}

	return allTags, nil
}

func delegateCoins(
	ctx sdk.Context, ak auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins,
) (sdk.Tags, sdk.Error) {

	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}

	acc := getAccount(ctx, ak, addr)
	if acc == nil {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}

	oldCoins := acc.GetCoins()

	_, hasNeg := oldCoins.SafeSub(amt)
	if hasNeg {
		return nil, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", oldCoins, amt),
		)
	}

	if err := trackDelegation(acc, ctx.BlockHeader().Time, amt); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to track delegation: %v", err))
	}

	setAccount(ctx, ak, acc)

	return sdk.NewTags(
		sdk.TagAction, tags.ActionDelegateCoins,
		sdk.TagDelegator, addr.String(),
	), nil
}

func undelegateCoins(
	ctx sdk.Context, ak auth.AccountKeeper, addr sdk.AccAddress, amt sdk.Coins,
) (sdk.Tags, sdk.Error) {

	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}

	acc := getAccount(ctx, ak, addr)
	if acc == nil {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}

	if err := trackUndelegation(acc, amt); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to track undelegation: %v", err))
	}

	setAccount(ctx, ak, acc)

	return sdk.NewTags(
		sdk.TagAction, tags.ActionUndelegateCoins,
		sdk.TagDelegator, addr.String(),
	), nil
}

// CONTRACT: assumes that amt is valid.
func trackDelegation(acc auth.Account, blockTime time.Time, amt sdk.Coins) error {
	vacc, ok := acc.(auth.VestingAccount)
	if ok {
		vacc.TrackDelegation(blockTime, amt)
		return nil
	}

	return acc.SetCoins(acc.GetCoins().Sub(amt))
}

// CONTRACT: assumes that amt is valid.
func trackUndelegation(acc auth.Account, amt sdk.Coins) error {
	vacc, ok := acc.(auth.VestingAccount)
	if ok {
		vacc.TrackUndelegation(amt)
		return nil
	}

	return acc.SetCoins(acc.GetCoins().Add(amt))
}

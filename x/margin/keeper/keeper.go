package keeper

import (
	"strings"

	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	bankKeeper types.BankKeeper
	clpKeeper  types.CLPKeeper
	paramStore paramtypes.Subspace
}

func NewKeeper(storeKey sdk.StoreKey,
	cdc codec.BinaryCodec,
	bankKeeper types.BankKeeper,
	clpKeeper types.CLPKeeper,
	ps paramtypes.Subspace) types.Keeper {

	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{bankKeeper: bankKeeper, clpKeeper: clpKeeper, paramStore: ps, storeKey: storeKey, cdc: cdc}
}

func (k Keeper) SetMTP(ctx sdk.Context, mtp *types.MTP) error {
	if err := mtp.Validate(); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetMTPKey(mtp.CollateralAsset, mtp.Address)
	store.Set(key, k.cdc.MustMarshal(mtp))
	return nil
}

func (k Keeper) GetMTP(ctx sdk.Context, symbol string, mtpAddress string) (types.MTP, error) {
	var mtp types.MTP
	key := types.GetMTPKey(symbol, mtpAddress)
	store := ctx.KVStore(k.storeKey)
	if !store.Has(key) {
		return mtp, types.ErrMTPDoesNotExist
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshal(bz, &mtp)
	return mtp, nil
}

func (k Keeper) GetMTPIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.MTPPrefix)
}

func (k Keeper) GetMTPs(ctx sdk.Context) []*types.MTP {
	var mtpList []*types.MTP
	iterator := k.GetMTPIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var mtp types.MTP
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &mtp)
		mtpList = append(mtpList, &mtp)
	}
	return mtpList
}

func (k Keeper) GetMTPsForAsset(ctx sdk.Context, asset string) []*types.MTP {
	var mtpList []*types.MTP
	iterator := k.GetMTPIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var mtp types.MTP
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &mtp)
		if mtp.CollateralAsset == asset {
			mtpList = append(mtpList, &mtp)
		}
	}
	return mtpList
}

func (k Keeper) GetAssetsForMTP(ctx sdk.Context, mtpAddress sdk.Address) []string {
	var assetList []string
	iterator := k.GetMTPIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var mtp types.MTP
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &mtp)
		if mtpAddress.String() == mtp.Address {
			assetList = append(assetList, mtp.CollateralAsset)
		}
	}
	return assetList
}

func (k Keeper) DestroyMTP(ctx sdk.Context, symbol string, mtpAddress string) error {
	key := types.GetMTPKey(symbol, mtpAddress)
	store := ctx.KVStore(k.storeKey)
	if !store.Has(key) {
		return types.ErrMTPDoesNotExist
	}
	store.Delete(key)
	return nil
}

func (k Keeper) ClpKeeper() types.CLPKeeper {
	return k.clpKeeper
}

func (k Keeper) BankKeeper() types.BankKeeper {
	return k.bankKeeper
}

func (k Keeper) GetLeverageParam(ctx sdk.Context) sdk.Uint {
	var leverageMax sdk.Uint
	k.paramStore.Get(ctx, types.KeyLeverageMaxParam, &leverageMax)
	return leverageMax
}

func (k Keeper) CustodySwap(ctx sdk.Context, pool clptypes.Pool, to string, sentAmount sdk.Uint) (sdk.Uint, error) {
	/*
	   calculate swap fee based on math spec
	   lambda_L = (0,1)
	   Notice this is NOT a simple hybrid of uniswap model and Thorchain slippaged based model
	   But a upgraded version that include swap, updating bouding curve (to be inside the old one)
	   One can think about this as a state jump:
	*/
	normalizationFactor, adjustExternalToken, err := k.ClpKeeper().GetNormalizationFactorForAsset(ctx, pool.ExternalAsset.Symbol)
	if err != nil {
		return sdk.ZeroUint(), err
	}

	X, XL, x, Y, YL, toRowan := SetInputs(sentAmount, to, pool)

	if !k.clpKeeper.ValidateZero([]sdk.Uint{X, x, Y}) {
		return sdk.ZeroUint(), nil
	}

	nf := sdk.NewUintFromBigInt(normalizationFactor.RoundInt().BigInt())
	if adjustExternalToken {
		if toRowan {
			X = X.Mul(nf)
			XL = XL.Mul(nf)
			x = x.Mul(nf)
		} else {
			Y = Y.Mul(nf)
			YL = YL.Mul(nf)
		}
	} else {
		if toRowan {
			Y = Y.Mul(nf)
			YL = YL.Mul(nf)
		} else {
			X = X.Mul(nf)
			XL = XL.Mul(nf)
			x = x.Mul(nf)
		}
	}

	minLen := k.clpKeeper.GetMinLen([]sdk.Uint{X, x, Y})
	Xd := k.clpKeeper.ReducePrecision(sdk.NewDecFromBigInt(X.BigInt()), minLen)
	XLd := k.clpKeeper.ReducePrecision(sdk.NewDecFromBigInt(XL.BigInt()), minLen)
	xd := k.clpKeeper.ReducePrecision(sdk.NewDecFromBigInt(x.BigInt()), minLen)
	Yd := k.clpKeeper.ReducePrecision(sdk.NewDecFromBigInt(Y.BigInt()), minLen)
	YLd := k.clpKeeper.ReducePrecision(sdk.NewDecFromBigInt(YL.BigInt()), minLen)

	numerator1 := xd.Mul(Yd.Add(YLd))
	denominator1 := xd.Add(Xd.Add(XLd))
	quotient1 := numerator1.Quo(denominator1)

	numerator2 := xd.Mul(Yd.Add(YLd)).Mul(Xd.Add(XLd))
	denominator2 := xd.Add(Xd.Add(XLd))
	denominator2 = denominator2.Mul(denominator2)
	quotient2 := numerator2.Quo(denominator2)

	y := quotient1.Add(quotient2)
	y = k.clpKeeper.IncreasePrecision(y, minLen)
	if !toRowan {
		y = y.Quo(normalizationFactor)
	}

	swapResult := sdk.NewUintFromBigInt(y.RoundInt().BigInt())

	if swapResult.GTE(Y) {
		return sdk.ZeroUint(), clptypes.ErrNotEnoughAssetTokens
	}

	return swapResult, nil
}

func SetInputs(sentAmount sdk.Uint, to string, pool clptypes.Pool) (sdk.Uint, sdk.Uint, sdk.Uint, sdk.Uint, sdk.Uint, bool) {
	var X sdk.Uint
	var XL sdk.Uint
	var Y sdk.Uint
	var YL sdk.Uint
	var x sdk.Uint
	toRowan := true
	if to == types.GetSettlementAsset() {
		Y = pool.NativeAssetBalance
		YL = pool.NativeLiabilities
		X = pool.ExternalAssetBalance
		XL = pool.ExternalLiabilities
	} else {
		X = pool.NativeAssetBalance
		XL = pool.NativeLiabilities
		Y = pool.ExternalAssetBalance
		YL = pool.ExternalLiabilities
		toRowan = false
	}
	x = sentAmount

	return X, XL, x, Y, YL, toRowan
}

func (k Keeper) Borrow(ctx sdk.Context, collateralAsset string, collateralAmount sdk.Uint, borrowAmount sdk.Uint, mtp types.MTP, pool clptypes.Pool, leverage sdk.Uint) error {
	mtp.CollateralAmount = mtp.CollateralAmount.Add(collateralAmount)
	mtp.LiabilitiesP = mtp.LiabilitiesP.Add(collateralAmount.Mul(leverage))
	mtp.CustodyAmount = mtp.CustodyAmount.Add(borrowAmount)
	mtp.Leverage = leverage

	h, err := k.UpdateMTPHealth(ctx, mtp, pool) // set mtp in func or return h?
	if err != nil {
		return err
	}
	mtp.MtpHealth = h

	mtpAddress, err := sdk.AccAddressFromBech32(mtp.Address)
	if err != nil {
		return err
	}
	collateralCoins := sdk.NewCoins(sdk.NewCoin(collateralAsset, sdk.NewIntFromBigInt(collateralAmount.BigInt())))
	err = k.BankKeeper().SendCoinsFromAccountToModule(ctx, mtpAddress, types.ModuleName, collateralCoins)
	if err != nil {
		return err
	}

	return k.SetMTP(ctx, &mtp)
}

func (k Keeper) UpdatePoolHealth(ctx sdk.Context, pool clptypes.Pool) error {
	// can be both X and Y
	ExternalAssetBalance := pool.ExternalAssetBalance
	ExternalLiabilities := pool.ExternalLiabilities
	NativeAssetBalance := pool.NativeAssetBalance
	NativeLiabilities := pool.NativeLiabilities

	mul1 := ExternalAssetBalance.Quo(ExternalAssetBalance.Add(ExternalLiabilities))
	mul2 := NativeAssetBalance.Quo(NativeAssetBalance.Add(NativeLiabilities))

	H := mul1.Mul(mul2)

	pool.Health = sdk.NewDecFromBigInt(H.BigInt())
	return k.ClpKeeper().SetPool(ctx, &pool)
}

// TODO Rename to CalcMTPHealth if not storing.
func (k Keeper) UpdateMTPHealth(ctx sdk.Context, mtp types.MTP, pool clptypes.Pool) (sdk.Dec, error) {
	// delta x in calculate in y currency
	nativeAsset := types.GetSettlementAsset()

	var err error
	var normalizedCollateral, normalizedLiabilities, normalizedCustody sdk.Uint
	if strings.EqualFold(mtp.CollateralAsset, nativeAsset) { // collateral is native
		normalizedCustody, err = k.CustodySwap(ctx, pool, mtp.CollateralAsset, mtp.CustodyAmount)
		if err != nil {
			return sdk.Dec{}, err
		}
		normalizedCollateral = mtp.CollateralAmount
		normalizedLiabilities = mtp.LiabilitiesP
	} else { // collateral is external
		normalizedCollateral, err = k.CustodySwap(ctx, pool, mtp.CustodyAsset, mtp.CollateralAmount)
		if err != nil {
			return sdk.Dec{}, err
		}
		normalizedLiabilities, err = k.CustodySwap(ctx, pool, mtp.CustodyAsset, mtp.LiabilitiesP)
		if err != nil {
			return sdk.Dec{}, err
		}
		normalizedCustody = mtp.CustodyAmount
	}

	if normalizedCollateral.Add(normalizedLiabilities).Add(normalizedCustody).Equal(sdk.ZeroUint()) {
		return sdk.Dec{}, types.ErrMTPInvalid
	}

	health := normalizedCollateral.Quo(normalizedCollateral.Add(normalizedLiabilities).Add(normalizedCustody))

	return sdk.NewDecFromBigInt(health.BigInt()), nil
}

func (k Keeper) TakeInCustody(ctx sdk.Context, mtp types.MTP, pool clptypes.Pool) error {
	nativeAsset := types.GetSettlementAsset()

	if strings.EqualFold(mtp.CollateralAsset, nativeAsset) {
		pool.NativeAssetBalance = pool.NativeAssetBalance.Sub(mtp.CollateralAmount)
		pool.NativeLiabilities = pool.NativeLiabilities.Add(mtp.CollateralAmount)
		pool.ExternalCustody = pool.ExternalCustody.Add(mtp.CustodyAmount)
	} else {
		pool.ExternalAssetBalance = pool.ExternalAssetBalance.Sub(mtp.CollateralAmount)
		pool.ExternalLiabilities = pool.NativeLiabilities.Add(mtp.CollateralAmount)
		pool.NativeCustody = pool.NativeCustody.Add(mtp.CustodyAmount)
	}

	return k.ClpKeeper().SetPool(ctx, &pool)
}

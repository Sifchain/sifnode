//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	"fmt"
	"math/big"
	"strings"

	adminkeeper "github.com/Sifchain/sifnode/x/admin/keeper"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	storeKey    sdk.StoreKey
	cdc         codec.BinaryCodec
	bankKeeper  types.BankKeeper
	clpKeeper   types.CLPKeeper
	adminKeeper adminkeeper.Keeper
	paramStore  paramtypes.Subspace
}

func NewKeeper(storeKey sdk.StoreKey,
	cdc codec.BinaryCodec,
	bankKeeper types.BankKeeper,
	clpKeeper types.CLPKeeper,
	adminKeeper adminkeeper.Keeper,
	ps paramtypes.Subspace) types.Keeper {

	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		bankKeeper:  bankKeeper,
		clpKeeper:   clpKeeper,
		adminKeeper: adminKeeper,
		paramStore:  ps,
		storeKey:    storeKey,
		cdc:         cdc,
	}
}

func (k Keeper) GetMTPCount(ctx sdk.Context) uint64 {
	var count uint64
	countBz := ctx.KVStore(k.storeKey).Get(types.MTPCountPrefix)
	if countBz == nil {
		count = 0
	} else {
		count = types.GetUint64FromBytes(countBz)
	}
	return count
}

func (k Keeper) GetOpenMTPCount(ctx sdk.Context) uint64 {
	var count uint64
	countBz := ctx.KVStore(k.storeKey).Get(types.OpenMTPCountPrefix)
	if countBz == nil {
		count = 0
	} else {
		count = types.GetUint64FromBytes(countBz)
	}
	return count
}

func (k Keeper) SetMTP(ctx sdk.Context, mtp *types.MTP) error {
	store := ctx.KVStore(k.storeKey)
	count := k.GetMTPCount(ctx)
	openCount := k.GetOpenMTPCount(ctx)

	if mtp.Id == 0 {
		// increment global id count
		count++
		mtp.Id = count
		store.Set(types.MTPCountPrefix, types.GetUint64Bytes(count))
		// increment open mtp count
		openCount++
		store.Set(types.OpenMTPCountPrefix, types.GetUint64Bytes(openCount))
	}

	if err := mtp.Validate(); err != nil {
		return err
	}
	key := types.GetMTPKey(mtp.Address, mtp.Id)
	store.Set(key, k.cdc.MustMarshal(mtp))
	return nil
}

func (k Keeper) GetMTP(ctx sdk.Context, mtpAddress string, id uint64) (types.MTP, error) {
	var mtp types.MTP
	key := types.GetMTPKey(mtpAddress, id)
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

func (k Keeper) GetMTPsForCollateralAsset(ctx sdk.Context, asset string) []*types.MTP {
	var mtpList []*types.MTP
	iterator := k.GetMTPIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var mtp types.MTP
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &mtp)
		if strings.EqualFold(mtp.CollateralAsset, asset) {
			mtpList = append(mtpList, &mtp)
		}
	}
	return mtpList
}

func (k Keeper) GetMTPsForCustodyAsset(ctx sdk.Context, asset string) []*types.MTP {
	var mtpList []*types.MTP
	iterator := k.GetMTPIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var mtp types.MTP
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &mtp)
		if strings.EqualFold(mtp.CustodyAsset, asset) {
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

func (k Keeper) GetMTPsForAddress(ctx sdk.Context, mtpAddress sdk.Address) []*types.MTP {
	var mtps []*types.MTP
	iterator := k.GetMTPIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var mtp types.MTP
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &mtp)
		if mtpAddress.String() == mtp.Address {
			mtps = append(mtps, &mtp)
		}
	}
	return mtps
}

func (k Keeper) DestroyMTP(ctx sdk.Context, mtpAddress string, id uint64) error {
	key := types.GetMTPKey(mtpAddress, id)
	store := ctx.KVStore(k.storeKey)
	if !store.Has(key) {
		return types.ErrMTPDoesNotExist
	}
	store.Delete(key)
	// decrement open mtp count
	openCount := k.GetOpenMTPCount(ctx)
	openCount--
	store.Set(types.OpenMTPCountPrefix, types.GetUint64Bytes(openCount))
	return nil
}

func (k Keeper) ClpKeeper() types.CLPKeeper {
	return k.clpKeeper
}

func (k Keeper) BankKeeper() types.BankKeeper {
	return k.bankKeeper
}

func (k Keeper) AdminKeeper() adminkeeper.Keeper {
	return k.adminKeeper
}

func (k Keeper) CustodySwap(ctx sdk.Context, pool clptypes.Pool, to string, sentAmount sdk.Uint) (sdk.Uint, error) {
	/*
	   calculate swap fee based on math spec
	   lambda_L = (0,1)
	   Notice this is NOT a simple hybrid of uniswap model and Thorchain slippaged based model
	   But a upgraded version that include swap, updating bouding curve (to be inside the old one)
	   One can think about this as a state jump:
	*/

	normalizationFactor, adjustExternalToken, err := k.ClpKeeper().GetNormalizationFactorFromAsset(ctx, *pool.ExternalAsset)
	if err != nil {
		return sdk.ZeroUint(), err
	}

	X, XL, x, Y, YL, toRowan := SetInputs(sentAmount, to, pool)

	if !k.clpKeeper.ValidateZero([]sdk.Uint{X, x, Y}) {
		return sdk.ZeroUint(), nil
	}

	pmtpCurrentRunningRate := k.clpKeeper.GetPmtpRateParams(ctx).PmtpCurrentRunningRate

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

	y = CalcCustodySwapPmtp(toRowan, y, pmtpCurrentRunningRate)

	swapResult := sdk.NewUintFromBigInt(y.RoundInt().BigInt())

	if swapResult.GTE(Y) {
		return sdk.ZeroUint(), clptypes.ErrNotEnoughAssetTokens
	}

	return swapResult, nil
}

func CalcCustodySwapPmtp(toRowan bool, y, pmtpCurrentRunningRate sdk.Dec) sdk.Dec {
	if toRowan {
		return y.Quo(sdk.NewDec(1).Add(pmtpCurrentRunningRate))
	}
	return y.Mul(sdk.NewDec(1).Add(pmtpCurrentRunningRate))
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

func (k Keeper) Borrow(ctx sdk.Context, collateralAsset string, collateralAmount sdk.Uint, custodyAmount sdk.Uint, mtp *types.MTP, pool *clptypes.Pool, eta sdk.Dec) error {
	mtpAddress, err := sdk.AccAddressFromBech32(mtp.Address)
	if err != nil {
		return err
	}
	collateralCoin := sdk.NewCoin(collateralAsset, sdk.NewIntFromBigInt(collateralAmount.BigInt()))

	if !k.bankKeeper.HasBalance(ctx, mtpAddress, collateralCoin) {
		return clptypes.ErrBalanceNotAvailable
	}

	collateralAmountDec := sdk.NewDecFromBigInt(collateralAmount.BigInt())
	liabilitiesPDec := collateralAmountDec.Mul(eta)

	mtp.CollateralAmount = mtp.CollateralAmount.Add(collateralAmount)

	mtp.LiabilitiesP = mtp.LiabilitiesP.Add(sdk.NewUintFromBigInt(liabilitiesPDec.TruncateInt().BigInt()))
	mtp.CustodyAmount = mtp.CustodyAmount.Add(custodyAmount)
	mtp.Leverage = eta.Add(sdk.OneDec())

	// print mtp.CustodyAmount
	ctx.Logger().Info(fmt.Sprintf("mtp.CustodyAmount: %s", mtp.CustodyAmount.String()))

	h, err := k.UpdateMTPHealth(ctx, *mtp, *pool) // set mtp in func or return h?
	if err != nil {
		return err
	}
	mtp.MtpHealth = h

	collateralCoins := sdk.NewCoins(collateralCoin)
	err = k.BankKeeper().SendCoinsFromAccountToModule(ctx, mtpAddress, clptypes.ModuleName, collateralCoins)
	if err != nil {
		return err
	}

	nativeAsset := types.GetSettlementAsset()

	if strings.EqualFold(mtp.CollateralAsset, nativeAsset) { // collateral is native
		pool.NativeAssetBalance = pool.NativeAssetBalance.Add(collateralAmount)
		pool.NativeLiabilities = pool.NativeLiabilities.Add(mtp.LiabilitiesP)
	} else { // collateral is external
		pool.ExternalAssetBalance = pool.ExternalAssetBalance.Add(collateralAmount)
		pool.ExternalLiabilities = pool.ExternalLiabilities.Add(mtp.LiabilitiesP)
	}
	err = k.ClpKeeper().SetPool(ctx, pool)
	if err != nil {
		return err
	}

	return k.SetMTP(ctx, mtp)
}

func (k Keeper) UpdatePoolHealth(ctx sdk.Context, pool *clptypes.Pool) error {
	pool.Health = k.CalculatePoolHealth(pool)
	return k.ClpKeeper().SetPool(ctx, pool)
}

func (k Keeper) CalculatePoolHealth(pool *clptypes.Pool) sdk.Dec {
	// can be both X and Y
	ExternalAssetBalance := sdk.NewDecFromBigInt(pool.ExternalAssetBalance.BigInt())
	ExternalLiabilities := sdk.NewDecFromBigInt(pool.ExternalLiabilities.BigInt())
	NativeAssetBalance := sdk.NewDecFromBigInt(pool.NativeAssetBalance.BigInt())
	NativeLiabilities := sdk.NewDecFromBigInt(pool.NativeLiabilities.BigInt())

	if ExternalAssetBalance.Add(ExternalLiabilities).IsZero() || NativeAssetBalance.Add(NativeLiabilities).IsZero() {
		return sdk.ZeroDec()
	}

	mul1 := ExternalAssetBalance.Quo(ExternalAssetBalance.Add(ExternalLiabilities))
	mul2 := NativeAssetBalance.Quo(NativeAssetBalance.Add(NativeLiabilities))

	H := mul1.Mul(mul2)

	return H
}

// TODO Rename to CalcMTPHealth if not storing.
func (k Keeper) UpdateMTPHealth(ctx sdk.Context, mtp types.MTP, pool clptypes.Pool) (sdk.Dec, error) {
	// delta x in calculate in y currency
	nativeAsset := types.GetSettlementAsset()

	var normalizedCollateral, normalizedLiabilitiesP, normalizedLiabilitiesI, normalizedCustody sdk.Dec
	if strings.EqualFold(mtp.CollateralAsset, nativeAsset) { // collateral is native
		normalizedCustodyInt, err := k.CustodySwap(ctx, pool, mtp.CollateralAsset, mtp.CustodyAmount)
		if err != nil {
			return sdk.Dec{}, err
		}
		normalizedCustody = sdk.NewDecFromBigInt(normalizedCustodyInt.BigInt())
		normalizedCollateral = sdk.NewDecFromBigInt(mtp.CollateralAmount.BigInt())
		normalizedLiabilitiesP = sdk.NewDecFromBigInt(mtp.LiabilitiesP.BigInt())
		normalizedLiabilitiesI = sdk.NewDecFromBigInt(mtp.LiabilitiesI.BigInt())
	} else { // collateral is external
		normalizedCollateralInt, err := k.CustodySwap(ctx, pool, mtp.CustodyAsset, mtp.CollateralAmount)
		if err != nil {
			return sdk.Dec{}, err
		}
		normalizedCollateral = sdk.NewDecFromBigInt(normalizedCollateralInt.BigInt())
		normalizedLiabilitiesPInt, err := k.CustodySwap(ctx, pool, mtp.CustodyAsset, mtp.LiabilitiesP)
		if err != nil {
			return sdk.Dec{}, err
		}
		normalizedLiabilitiesP = sdk.NewDecFromBigInt(normalizedLiabilitiesPInt.BigInt())
		normalizedLiabilitiesIInt, err := k.CustodySwap(ctx, pool, mtp.CustodyAsset, mtp.LiabilitiesI)
		if err != nil {
			return sdk.Dec{}, err
		}
		normalizedLiabilitiesI = sdk.NewDecFromBigInt(normalizedLiabilitiesIInt.BigInt())
		normalizedCustody = sdk.NewDecFromBigInt(mtp.CustodyAmount.BigInt())
	}

	if normalizedCollateral.Add(normalizedLiabilitiesP).Add(normalizedLiabilitiesI).Add(normalizedCustody).Equal(sdk.ZeroDec()) {
		return sdk.Dec{}, types.ErrMTPInvalid
	}

	health := normalizedCollateral.Quo(normalizedCollateral.Add(normalizedLiabilitiesP).Add(normalizedLiabilitiesI).Add(normalizedCustody))

	return health, nil
}

func (k Keeper) TakeInCustody(ctx sdk.Context, mtp types.MTP, pool *clptypes.Pool) error {
	nativeAsset := types.GetSettlementAsset()

	if strings.EqualFold(mtp.CustodyAsset, nativeAsset) {
		pool.NativeAssetBalance = pool.NativeAssetBalance.Sub(mtp.CustodyAmount)
		pool.NativeCustody = pool.NativeCustody.Add(mtp.CustodyAmount)
	} else {
		pool.ExternalAssetBalance = pool.ExternalAssetBalance.Sub(mtp.CustodyAmount)
		pool.ExternalCustody = pool.ExternalCustody.Add(mtp.CustodyAmount)
	}

	return k.ClpKeeper().SetPool(ctx, pool)
}

func (k Keeper) TakeOutCustody(ctx sdk.Context, mtp types.MTP, pool *clptypes.Pool) error {
	nativeAsset := types.GetSettlementAsset()

	if strings.EqualFold(mtp.CustodyAsset, nativeAsset) {
		pool.NativeCustody = pool.NativeCustody.Sub(mtp.CustodyAmount)
		pool.NativeAssetBalance = pool.NativeAssetBalance.Add(mtp.CustodyAmount)
	} else {
		pool.ExternalCustody = pool.ExternalCustody.Sub(mtp.CustodyAmount)
		pool.ExternalAssetBalance = pool.ExternalAssetBalance.Add(mtp.CustodyAmount)
	}

	return k.ClpKeeper().SetPool(ctx, pool)
}

func (k Keeper) Repay(ctx sdk.Context, mtp *types.MTP, pool clptypes.Pool, repayAmount sdk.Uint, takeInsurance bool) error {
	// nolint:ineffassign
	returnAmount, debtP, debtI := sdk.ZeroUint(), sdk.ZeroUint(), sdk.ZeroUint()
	CollateralAmount := mtp.CollateralAmount
	LiabilitiesP := mtp.LiabilitiesP
	LiabilitiesI := mtp.LiabilitiesI

	var err error
	mtp.MtpHealth, err = k.UpdateMTPHealth(ctx, *mtp, pool)
	if err != nil {
		return err
	}

	have := repayAmount.Add(CollateralAmount)
	owe := LiabilitiesP.Add(LiabilitiesI)

	fmt.Println("have:", have)
	fmt.Println("owe:", owe)
	fmt.Println("LiabilitiesP:", LiabilitiesP)

	if have.LT(LiabilitiesP) {
		//can't afford principle liability
		returnAmount = sdk.ZeroUint()
		debtP = LiabilitiesP.Sub(have)
		debtI = LiabilitiesI
	} else if have.LT(owe) {
		// v principle liability; x excess liability
		returnAmount = sdk.ZeroUint()
		debtP = sdk.ZeroUint()
		debtI = LiabilitiesP.Add(LiabilitiesI).Sub(have)
	} else {
		// can afford both
		returnAmount = have.Sub(LiabilitiesP).Sub(LiabilitiesI)
		debtP = sdk.ZeroUint()
		debtI = sdk.ZeroUint()
	}

	fmt.Println("returnAmount:", returnAmount)

	if !returnAmount.IsZero() {
		actualReturnAmount := returnAmount
		if takeInsurance {
			takePercentage := k.GetForceCloseFundPercentage(ctx)
			returnAmountDec := sdk.NewDecFromBigInt(returnAmount.BigInt())
			takeAmount := sdk.NewUintFromBigInt(takePercentage.Mul(returnAmountDec).TruncateInt().BigInt())
			actualReturnAmount = returnAmount.Sub(takeAmount)

			if !takeAmount.IsZero() {
				takeCoins := sdk.NewCoins(sdk.NewCoin(mtp.CollateralAsset, sdk.NewIntFromBigInt(takeAmount.BigInt())))
				fundAddr := k.GetInsuranceFundAddress(ctx)
				err = k.BankKeeper().SendCoinsFromModuleToAccount(ctx, clptypes.ModuleName, fundAddr, takeCoins)
				if err != nil {
					return err
				}
				k.EmitRepayInsuranceFund(ctx, mtp, takeAmount)
			}
		}

		if !actualReturnAmount.IsZero() {
			var coins sdk.Coins
			returnCoin := sdk.NewCoin(mtp.CollateralAsset, sdk.NewIntFromBigInt(actualReturnAmount.BigInt()))
			returnCoins := coins.Add(returnCoin)
			addr, err := sdk.AccAddressFromBech32(mtp.Address)
			if err != nil {
				return err
			}
			err = k.BankKeeper().MintCoins(ctx, clptypes.ModuleName, returnCoins)
			if err != nil {
				return err
			}
			err = k.BankKeeper().SendCoinsFromModuleToAccount(ctx, clptypes.ModuleName, addr, returnCoins)
			if err != nil {
				return err
			}
		}
	}

	nativeAsset := types.GetSettlementAsset()

	if strings.EqualFold(mtp.CollateralAsset, nativeAsset) {
		pool.NativeAssetBalance = pool.NativeAssetBalance.Sub(returnAmount).Sub(debtI).Sub(debtP)
		pool.NativeLiabilities = pool.NativeLiabilities.Sub(mtp.LiabilitiesP)
	} else {
		pool.ExternalAssetBalance = pool.ExternalAssetBalance.Sub(returnAmount).Sub(debtI).Sub(debtP)
		pool.ExternalLiabilities = pool.ExternalLiabilities.Sub(mtp.LiabilitiesP)
	}

	err = k.DestroyMTP(ctx, mtp.Address, mtp.Id)
	if err != nil {
		return err
	}

	return k.ClpKeeper().SetPool(ctx, &pool)
}

func (k Keeper) UpdateMTPInterestLiabilities(ctx sdk.Context, mtp *types.MTP, interestRate sdk.Dec) error {
	var liabilitiesIRat, liabilitiesRat, rate big.Rat

	rate.SetFloat64(interestRate.MustFloat64())

	liabilitiesRat.SetInt(mtp.LiabilitiesP.BigInt().Add(mtp.LiabilitiesP.BigInt(), mtp.LiabilitiesI.BigInt()))
	liabilitiesIRat.Mul(&rate, &liabilitiesRat)

	liabilitiesINew := liabilitiesIRat.Num().Quo(liabilitiesRat.Num(), liabilitiesIRat.Denom())
	mtp.LiabilitiesI = sdk.NewUintFromBigInt(liabilitiesINew.Add(liabilitiesINew, mtp.LiabilitiesI.BigInt()))

	return k.SetMTP(ctx, mtp)
}

func (k Keeper) InterestRateComputation(ctx sdk.Context, pool clptypes.Pool) (sdk.Dec, error) {
	interestRateMax := k.GetInterestRateMax(ctx)
	interestRateMin := k.GetInterestRateMin(ctx)
	interestRateIncrease := k.GetInterestRateIncrease(ctx)
	interestRateDecrease := k.GetInterestRateDecrease(ctx)
	healthGainFactor := k.GetHealthGainFactor(ctx)

	prevInterestRate := pool.InterestRate

	externalAssetBalance := sdk.NewDecFromBigInt(pool.ExternalAssetBalance.BigInt())
	ExternalLiabilities := sdk.NewDecFromBigInt(pool.ExternalLiabilities.BigInt())
	NativeAssetBalance := sdk.NewDecFromBigInt(pool.NativeAssetBalance.BigInt())
	NativeLiabilities := sdk.NewDecFromBigInt(pool.NativeLiabilities.BigInt())

	mul1 := externalAssetBalance.Add(ExternalLiabilities).Quo(externalAssetBalance)
	mul2 := NativeAssetBalance.Add(NativeLiabilities).Quo(NativeAssetBalance)

	targetInterestRate := healthGainFactor.Mul(mul1).Mul(mul2).Add(k.GetSQ(ctx, pool))

	interestRateChange := targetInterestRate.Sub(prevInterestRate)
	interestRate := prevInterestRate
	if interestRateChange.LTE(interestRateDecrease.Mul(sdk.NewDec(-1))) && interestRateChange.LTE(interestRateIncrease) {
		interestRate = targetInterestRate
	} else if interestRateChange.GT(interestRateIncrease) {
		interestRate = prevInterestRate.Add(interestRateIncrease)
	} else if interestRateChange.LT(interestRateDecrease.Mul(sdk.NewDec(-1))) {
		interestRate = prevInterestRate.Sub(interestRateDecrease)
	}

	newInterestRate := interestRate

	if interestRate.GT(interestRateMin) && interestRate.LT(interestRateMax) {
		newInterestRate = interestRate
	} else if interestRate.LTE(interestRateMin) {
		newInterestRate = interestRateMin
	} else if interestRate.GTE(interestRateMax) {
		newInterestRate = interestRateMax
	}

	return newInterestRate, nil
}

func (k Keeper) GetSQ(ctx sdk.Context, pool clptypes.Pool) sdk.Dec {
	q := k.ClpKeeper().GetRemovalQueue(ctx, pool.ExternalAsset.Symbol)
	if q.Count == 0 {
		return sdk.ZeroDec()
	}

	value := sdk.NewDecFromBigInt(q.TotalValue.BigInt())
	blocks := sdk.NewDec(ctx.BlockHeight() - q.StartHeight)
	modifier, err := sdk.NewDecFromStr("10000000000000000000000000")
	if err != nil {
		panic(err)
	}

	sq := value.Mul(blocks).Quo(modifier)

	return sq
}

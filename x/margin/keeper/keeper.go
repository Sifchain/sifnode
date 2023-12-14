package keeper

import (
	"fmt"
	"math"
	"math/big"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	adminkeeper "github.com/Sifchain/sifnode/x/admin/keeper"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const MaxPageLimit = 100

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
	ps paramtypes.Subspace) Keeper {

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

func (k Keeper) GetAllMTPS(ctx sdk.Context) []*types.MTP {
	var mtpList []*types.MTP
	iterator := k.GetMTPIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		var mtp types.MTP
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &mtp)
		mtpList = append(mtpList, &mtp)
	}
	return mtpList
}

func (k Keeper) GetMTPs(ctx sdk.Context, pagination *query.PageRequest) ([]*types.MTP, *query.PageResponse, error) {
	var mtpList []*types.MTP
	store := ctx.KVStore(k.storeKey)
	mtpStore := prefix.NewStore(store, types.MTPPrefix)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: math.MaxUint64 - 1,
		}
	}

	pageRes, err := query.Paginate(mtpStore, pagination, func(key []byte, value []byte) error {
		var mtp types.MTP
		k.cdc.MustUnmarshal(value, &mtp)
		mtpList = append(mtpList, &mtp)
		return nil
	})

	return mtpList, pageRes, err
}

func (k Keeper) GetMTPsForPool(ctx sdk.Context, asset string, pagination *query.PageRequest) ([]*types.MTP, *query.PageResponse, error) {
	var mtps []*types.MTP

	store := ctx.KVStore(k.storeKey)
	mtpStore := prefix.NewStore(store, types.MTPPrefix)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: math.MaxUint64 - 1,
		}
	}

	pageRes, err := query.FilteredPaginate(mtpStore, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var mtp types.MTP
		k.cdc.MustUnmarshal(value, &mtp)
		if accumulate && (types.StringCompare(mtp.CustodyAsset, asset) || types.StringCompare(mtp.CollateralAsset, asset)) {
			mtps = append(mtps, &mtp)
			return true, nil
		}

		return false, nil
	})

	return mtps, pageRes, err
}

func (k Keeper) GetMTPsForAddress(ctx sdk.Context, mtpAddress sdk.Address, pagination *query.PageRequest) ([]*types.MTP, *query.PageResponse, error) {
	var mtps []*types.MTP

	store := ctx.KVStore(k.storeKey)
	mtpStore := prefix.NewStore(store, types.GetMTPPrefixForAddress(mtpAddress.String()))

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: MaxPageLimit,
		}
	}

	if pagination.Limit > MaxPageLimit {
		return nil, nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", MaxPageLimit))
	}

	pageRes, err := query.Paginate(mtpStore, pagination, func(key []byte, value []byte) error {
		var mtp types.MTP
		k.cdc.MustUnmarshal(value, &mtp)
		mtps = append(mtps, &mtp)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return mtps, pageRes, nil
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

func (k Keeper) CLPSwap(ctx sdk.Context, sentAmount sdk.Uint, to string, pool clptypes.Pool) (sdk.Uint, error) {
	toAsset := ToAsset(to)

	marginEnabled := k.IsPoolEnabled(ctx, pool.ExternalAsset.Symbol)

	swapResult, err := k.ClpKeeper().CLPCalcSwap(ctx, sentAmount, toAsset, pool, marginEnabled)
	if err != nil {
		return sdk.Uint{}, err
	}
	if swapResult.IsZero() {
		return sdk.Uint{}, clptypes.ErrAmountTooLow
	}
	return swapResult, nil
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
	liabilitiesDec := collateralAmountDec.Mul(eta)

	mtp.CollateralAmount = mtp.CollateralAmount.Add(collateralAmount)

	mtp.Liabilities = mtp.Liabilities.Add(sdk.NewUintFromBigInt(liabilitiesDec.TruncateInt().BigInt()))
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

	if types.StringCompare(mtp.CollateralAsset, nativeAsset) { // collateral is native
		pool.NativeAssetBalance = pool.NativeAssetBalance.Add(collateralAmount)
		pool.NativeLiabilities = pool.NativeLiabilities.Add(mtp.Liabilities)
	} else { // collateral is external
		pool.ExternalAssetBalance = pool.ExternalAssetBalance.Add(collateralAmount)
		pool.ExternalLiabilities = pool.ExternalLiabilities.Add(mtp.Liabilities)
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

func (k Keeper) UpdateMTPHealth(ctx sdk.Context, mtp types.MTP, pool clptypes.Pool) (sdk.Dec, error) {
	xl := mtp.Liabilities

	if xl.IsZero() {
		return sdk.ZeroDec(), nil
	}
	// include unpaid interest in debt (from disabled incremental pay)
	if mtp.InterestUnpaidCollateral.GT(sdk.ZeroUint()) {
		xl = xl.Add(mtp.InterestUnpaidCollateral)
	}

	C, err := k.CLPSwap(ctx, mtp.CustodyAmount, mtp.CollateralAsset, pool)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	lr := sdk.NewDecFromBigInt(C.BigInt()).Quo(sdk.NewDecFromBigInt(xl.BigInt()))

	return lr, nil
}

func (k Keeper) TakeInCustody(ctx sdk.Context, mtp types.MTP, pool *clptypes.Pool) error {
	nativeAsset := types.GetSettlementAsset()

	if types.StringCompare(mtp.CustodyAsset, nativeAsset) {
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

	if types.StringCompare(mtp.CustodyAsset, nativeAsset) {
		pool.NativeCustody = pool.NativeCustody.Sub(mtp.CustodyAmount)
		pool.NativeAssetBalance = pool.NativeAssetBalance.Add(mtp.CustodyAmount)
	} else {
		pool.ExternalCustody = pool.ExternalCustody.Sub(mtp.CustodyAmount)
		pool.ExternalAssetBalance = pool.ExternalAssetBalance.Add(mtp.CustodyAmount)
	}

	return k.ClpKeeper().SetPool(ctx, pool)
}

func (k Keeper) Repay(ctx sdk.Context, mtp *types.MTP, pool *clptypes.Pool, repayAmount sdk.Uint, takeFundPayment bool) error {
	// nolint:staticcheck,ineffassign
	returnAmount, debtP, debtI := sdk.ZeroUint(), sdk.ZeroUint(), sdk.ZeroUint()
	Liabilities := mtp.Liabilities
	InterestUnpaidCollateral := mtp.InterestUnpaidCollateral

	var err error
	mtp.MtpHealth, err = k.UpdateMTPHealth(ctx, *mtp, *pool)
	if err != nil {
		return err
	}

	have := repayAmount
	owe := Liabilities.Add(InterestUnpaidCollateral)

	if have.LT(Liabilities) {
		//can't afford principle liability
		returnAmount = sdk.ZeroUint()
		debtP = Liabilities.Sub(have)
		debtI = InterestUnpaidCollateral
	} else if have.LT(owe) {
		// v principle liability; x excess liability
		returnAmount = sdk.ZeroUint()
		debtP = sdk.ZeroUint()
		debtI = Liabilities.Add(InterestUnpaidCollateral).Sub(have)
	} else {
		// can afford both
		returnAmount = have.Sub(Liabilities).Sub(InterestUnpaidCollateral)
		debtP = sdk.ZeroUint()
		debtI = sdk.ZeroUint()
	}
	if !returnAmount.IsZero() {
		actualReturnAmount := returnAmount
		if takeFundPayment {
			takePercentage := k.GetForceCloseFundPercentage(ctx)

			fundAddr := k.GetForceCloseFundAddress(ctx)
			takeAmount, err := k.TakeFundPayment(ctx, returnAmount, mtp.CollateralAsset, takePercentage, fundAddr)
			if err != nil {
				return err
			}
			actualReturnAmount = returnAmount.Sub(takeAmount)
			if !takeAmount.IsZero() {
				k.EmitFundPayment(ctx, mtp, takeAmount, mtp.CollateralAsset, types.EventRepayFund)
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
			err = k.BankKeeper().SendCoinsFromModuleToAccount(ctx, clptypes.ModuleName, addr, returnCoins)
			if err != nil {
				return err
			}
		}
	}

	nativeAsset := types.GetSettlementAsset()

	if types.StringCompare(mtp.CollateralAsset, nativeAsset) {
		pool.NativeAssetBalance = pool.NativeAssetBalance.Sub(returnAmount)
		pool.NativeLiabilities = pool.NativeLiabilities.Sub(mtp.Liabilities)
		pool.UnsettledNativeLiabilities = pool.UnsettledNativeLiabilities.Add(debtI).Add(debtP)
	} else {
		pool.ExternalAssetBalance = pool.ExternalAssetBalance.Sub(returnAmount)
		pool.ExternalLiabilities = pool.ExternalLiabilities.Sub(mtp.Liabilities)
		pool.UnsettledExternalLiabilities = pool.UnsettledExternalLiabilities.Add(debtI).Add(debtP)
	}
	err = k.DestroyMTP(ctx, mtp.Address, mtp.Id)
	if err != nil {
		return err
	}

	return k.ClpKeeper().SetPool(ctx, pool)
}

func (k Keeper) HandleInterestPayment(ctx sdk.Context, interestPayment sdk.Uint, mtp *types.MTP, pool *clptypes.Pool) sdk.Uint {
	incrementalInterestPaymentEnabled := k.GetIncrementalInterestPaymentEnabled(ctx)
	// if incremental payment on, pay interest
	if incrementalInterestPaymentEnabled {
		finalInterestPayment, err := k.IncrementalInterestPayment(ctx, interestPayment, mtp, pool)
		if err != nil {
			ctx.Logger().Error(sdkerrors.Wrap(err, "error executing incremental interest payment").Error())
		} else {
			return finalInterestPayment
		}
	} else { // else update unpaid mtp interest
		mtp.InterestUnpaidCollateral = interestPayment
	}
	return sdk.ZeroUint()
}

func (k Keeper) IncrementalInterestPayment(ctx sdk.Context, interestPayment sdk.Uint, mtp *types.MTP, pool *clptypes.Pool) (sdk.Uint, error) {
	// if mtp has unpaid interest, add to payment
	if mtp.InterestUnpaidCollateral.GT(sdk.ZeroUint()) {
		interestPayment = interestPayment.Add(mtp.InterestUnpaidCollateral)
	}

	// swap interest payment to custody asset for payment
	interestPaymentCustody, err := k.CLPSwap(ctx, interestPayment, mtp.CustodyAsset, *pool)
	if err != nil {
		return sdk.ZeroUint(), err
	}

	// if paying unpaid interest reset to 0
	mtp.InterestUnpaidCollateral = sdk.ZeroUint()

	// edge case, not enough custody to cover payment
	if interestPaymentCustody.GT(mtp.CustodyAmount) {
		// swap custody amount to collateral for updating interest unpaid
		custodyAmountCollateral, err := k.CLPSwap(ctx, mtp.CustodyAmount, mtp.CollateralAsset, *pool) // may need spot price here to not deduct fee
		if err != nil {
			return sdk.ZeroUint(), err
		}
		mtp.InterestUnpaidCollateral = interestPayment.Sub(custodyAmountCollateral)
		interestPayment = custodyAmountCollateral
		interestPaymentCustody = mtp.CustodyAmount
	}

	// add payment to total paid - collateral
	mtp.InterestPaidCollateral = mtp.InterestPaidCollateral.Add(interestPayment)

	// add payment to total paid - custody
	mtp.InterestPaidCustody = mtp.InterestPaidCustody.Add(interestPaymentCustody)

	// deduct interest payment from custody amount
	mtp.CustodyAmount = mtp.CustodyAmount.Sub(interestPaymentCustody)

	takePercentage := k.GetIncrementalInterestPaymentFundPercentage(ctx)
	fundAddr := k.GetIncrementalInterestPaymentFundAddress(ctx)
	takeAmount, err := k.TakeFundPayment(ctx, interestPaymentCustody, mtp.CustodyAsset, takePercentage, fundAddr)
	if err != nil {
		return sdk.ZeroUint(), err
	}
	actualInterestPaymentCustody := interestPaymentCustody.Sub(takeAmount)

	if !takeAmount.IsZero() {
		k.EmitFundPayment(ctx, mtp, takeAmount, mtp.CustodyAsset, types.EventIncrementalPayFund)
	}

	nativeAsset := types.GetSettlementAsset()

	if types.StringCompare(mtp.CustodyAsset, nativeAsset) { // custody is native
		pool.NativeCustody = pool.NativeCustody.Sub(interestPaymentCustody)
		pool.NativeAssetBalance = pool.NativeAssetBalance.Add(actualInterestPaymentCustody)
	} else { // custody is external
		pool.ExternalCustody = pool.ExternalCustody.Sub(interestPaymentCustody)
		pool.ExternalAssetBalance = pool.ExternalAssetBalance.Add(actualInterestPaymentCustody)
	}

	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return sdk.ZeroUint(), err
	}

	return actualInterestPaymentCustody, k.ClpKeeper().SetPool(ctx, pool)
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

	targetInterestRate := healthGainFactor.Mul(mul1).Mul(mul2)

	interestRateChange := targetInterestRate.Sub(prevInterestRate)
	interestRate := prevInterestRate
	if interestRateChange.GTE(interestRateDecrease.Mul(sdk.NewDec(-1))) && interestRateChange.LTE(interestRateIncrease) {
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

	sQ := k.GetSQFromBlocks(ctx, pool, newInterestRate)

	return newInterestRate.Add(sQ), nil
}

func (k Keeper) CheckMinLiabilities(ctx sdk.Context, collateralAmount sdk.Uint, eta sdk.Dec, pool clptypes.Pool, custodyAsset string) error {
	var interestRational, liabilitiesRational, rate big.Rat
	minInterestRate := k.GetInterestRateMin(ctx)

	collateralAmountDec := sdk.NewDecFromBigInt(collateralAmount.BigInt())
	liabilitiesDec := collateralAmountDec.Mul(eta)

	liabilities := sdk.NewUintFromBigInt(liabilitiesDec.TruncateInt().BigInt())

	rate.SetFloat64(minInterestRate.MustFloat64())
	liabilitiesRational.SetInt(liabilities.BigInt())
	interestRational.Mul(&rate, &liabilitiesRational)

	interestNew := interestRational.Num().Quo(interestRational.Num(), interestRational.Denom())

	samplePayment := sdk.NewUintFromBigInt(interestNew)

	if samplePayment.IsZero() && !minInterestRate.IsZero() {
		return types.ErrBorrowTooLow
	}

	// swap interest payment to custody asset
	_, err := k.CLPSwap(ctx, samplePayment, custodyAsset, pool)
	if err != nil {
		return types.ErrBorrowTooLow
	}

	return nil
}

func (k Keeper) GetSQBeginBlock(ctx sdk.Context, pool *clptypes.Pool) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetSQBeginBlockKey(pool))
	if bz == nil {
		return 0
	}
	return types.GetUint64FromBytes(bz)
}

func (k Keeper) SetSQBeginBlock(ctx sdk.Context, pool *clptypes.Pool, height uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetSQBeginBlockKey(pool), types.GetUint64Bytes(height))
}

func (k Keeper) TrackSQBeginBlock(ctx sdk.Context, pool *clptypes.Pool) {
	threshold := k.GetRemovalQueueThreshold(ctx)
	sqBeginBlock := k.GetSQBeginBlock(ctx, pool)
	if sqBeginBlock == 0 {
		if pool.Health.LTE(threshold) {
			k.SetSQBeginBlock(ctx, pool, uint64(ctx.BlockHeight()))
			k.EmitBelowRemovalThreshold(ctx, pool)
		}
	} else if pool.Health.GT(threshold) {
		k.SetSQBeginBlock(ctx, pool, 0)
		k.EmitAboveRemovalThreshold(ctx, pool)
	}
}

func (k Keeper) GetSQFromBlocks(ctx sdk.Context, pool clptypes.Pool, poolInterestRate sdk.Dec) sdk.Dec {
	beginBlock := k.GetSQBeginBlock(ctx, &pool)
	if beginBlock == 0 {
		return sdk.ZeroDec()
	}

	blocks := ctx.BlockHeight() - int64(beginBlock)
	maxInterestRate := k.GetInterestRateMax(ctx)
	poolInterestRateFloat, _ := poolInterestRate.Float64()
	minus := math.Pow(math.E, -1*poolInterestRateFloat*float64(blocks))
	minusDec, err := sdk.NewDecFromStr(fmt.Sprintf("%v", minus))
	if err != nil {
		minusDec = sdk.NewDec(0)
	}
	multipliedBy := sdk.NewDec(1).Sub(minusDec)
	return maxInterestRate.Mul(multipliedBy)
}

func (k Keeper) GetSQFromQueue(ctx sdk.Context, pool clptypes.Pool) sdk.Dec {
	q := k.ClpKeeper().GetRemovalQueue(ctx, pool.ExternalAsset.Symbol)
	if q.Count == 0 {
		return sdk.ZeroDec()
	}

	value := sdk.NewDecFromBigInt(q.TotalValue.BigInt())
	blocks := sdk.NewDec(ctx.BlockHeight() - q.StartHeight)
	modifier := k.GetSqModifier(ctx)
	sq := value.Mul(blocks).Quo(modifier)

	return sq
}

func ToAsset(asset string) clptypes.Asset {
	return clptypes.Asset{
		Symbol: asset,
	}
}

// get position of current block in epoch
func GetEpochPosition(ctx sdk.Context, epochLength int64) int64 {
	if epochLength <= 0 {
		epochLength = 1
	}
	currentHeight := ctx.BlockHeight()
	return currentHeight % epochLength
}

func (k Keeper) ForceCloseLong(ctx sdk.Context, mtp *types.MTP, pool *clptypes.Pool, isAdminClose bool, takeFundPayment bool) (sdk.Uint, error) {

	// check MTP health against threshold
	safetyFactor := k.GetSafetyFactor(ctx)

	epochLength := k.GetEpochLength(ctx)
	epochPosition := GetEpochPosition(ctx, epochLength)

	var err error
	if epochPosition > 0 {
		interestPayment := CalcMTPInterestLiabilities(mtp, pool.InterestRate, epochPosition, epochLength)

		finalInterestPayment := k.HandleInterestPayment(ctx, interestPayment, mtp, pool)

		nativeAsset := types.GetSettlementAsset()

		if types.StringCompare(mtp.CollateralAsset, nativeAsset) { // custody is external, payment is custody
			pool.BlockInterestExternal = pool.BlockInterestExternal.Add(finalInterestPayment)
		} else { // custody is native, payment is custody
			pool.BlockInterestNative = pool.BlockInterestNative.Add(finalInterestPayment)
		}

		mtp.MtpHealth, err = k.UpdateMTPHealth(ctx, *mtp, *pool)
		if err != nil {
			return sdk.ZeroUint(), err
		}
	}
	if !isAdminClose && mtp.MtpHealth.GT(safetyFactor) {
		return sdk.ZeroUint(), types.ErrMTPHealthy
	}

	err = k.TakeOutCustody(ctx, *mtp, pool)
	if err != nil {
		return sdk.ZeroUint(), err
	}

	repayAmount, err := k.CLPSwap(ctx, mtp.CustodyAmount, mtp.CollateralAsset, *pool)
	if err != nil {
		return sdk.ZeroUint(), err
	}

	err = k.Repay(ctx, mtp, pool, repayAmount, takeFundPayment)
	if err != nil {
		return sdk.ZeroUint(), err
	}

	return repayAmount, nil
}

func (k Keeper) TakeFundPayment(ctx sdk.Context, returnAmount sdk.Uint, returnAsset string, takePercentage sdk.Dec, fundAddr sdk.AccAddress) (sdk.Uint, error) {
	returnAmountDec := sdk.NewDecFromBigInt(returnAmount.BigInt())
	takeAmount := sdk.NewUintFromBigInt(takePercentage.Mul(returnAmountDec).TruncateInt().BigInt())

	if !takeAmount.IsZero() {
		takeCoins := sdk.NewCoins(sdk.NewCoin(returnAsset, sdk.NewIntFromBigInt(takeAmount.BigInt())))
		err := k.BankKeeper().SendCoinsFromModuleToAccount(ctx, clptypes.ModuleName, fundAddr, takeCoins)
		if err != nil {
			return sdk.ZeroUint(), err
		}
	}
	return takeAmount, nil
}

func (k Keeper) IsWhitelisted(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetWhitelistKey(address))
}

func (k Keeper) WhitelistAddress(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetWhitelistKey(address), []byte(address))
}

func (k Keeper) DewhitelistAddress(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetWhitelistKey(address))
}

func (k Keeper) GetWhitelist(ctx sdk.Context, pagination *query.PageRequest) ([]string, *query.PageResponse, error) {
	var list []string
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.WhitelistPrefix)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: math.MaxUint64 - 1,
		}
	}

	pageRes, err := query.Paginate(prefixStore, pagination, func(key []byte, value []byte) error {
		list = append(list, string(value))
		return nil
	})

	return list, pageRes, err
}

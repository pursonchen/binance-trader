package convert

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/pursonchen/go-binance/v2"
	"strconv"
)

type SpotClient struct {
	binanceSpotClient *binance.Client
}

func NewSpotClient(spotClient *binance.Client) *SpotClient {
	return &SpotClient{binanceSpotClient: spotClient}
}

type EstQuoteReq struct {
	BaseCcy  string `json:"baseCcy"`
	QuoteCcy string `json:"quoteCcy"`
}

type EstQuoteResp struct {
	Price string `json:"price"`
	Mins  int64  `json:"mins"`
}

var StableCoins = [2]string{"USDT", "BUSD"}

func (c *SpotClient) EstQuote(ctx context.Context, req *EstQuoteReq) (*EstQuoteResp, error) {

	// 目前支持稳定币兑换
	if req.BaseCcy != StableCoins[0] && req.QuoteCcy != StableCoins[0] {
		return nil, errors.New("coin not support")
	}

	res, err := c.binanceSpotClient.NewAveragePriceService().Symbol(req.BaseCcy + req.QuoteCcy).Do(ctx)

	if err != nil {
		return nil, err
	}

	return &EstQuoteResp{
		Mins:  res.Mins,
		Price: res.Price,
	}, nil
}

/*
订单类型 (orderTypes, type):

LIMIT 限价单
MARKET 市价单
STOP_LOSS 止损单
STOP_LOSS_LIMIT 限价止损单
TAKE_PROFIT 止盈单
TAKE_PROFIT_LIMIT 限价止盈单
LIMIT_MAKER 限价只挂单
订单返回类型 (newOrderRespType):

ACK
RESULT
FULL
订单方向 (方向 side):

BUY 买入
SELL 卖出
*/

type TradeReq struct {
	Symbol           string           `json:"symbol"`
	Side             binance.SideType `json:"side"` // BUY SELL
	Quantity         string           `json:"quantity"`
	NewClientOrderId string           `json:"newClientOrderId"`
}

type TradeResp struct {
	Symbol                   string                  `json:"symbol"`
	OrderID                  int64                   `json:"orderId"`
	ClientOrderID            string                  `json:"clientOrderId"`
	TransactTime             int64                   `json:"transactTime"`
	Price                    string                  `json:"price"`
	OrigQuantity             string                  `json:"origQty"`
	ExecutedQuantity         string                  `json:"executedQty"`
	CummulativeQuoteQuantity string                  `json:"cummulativeQuoteQty"`
	IsIsolated               bool                    `json:"isIsolated"`
	Status                   binance.OrderStatusType `json:"status"`
	TimeInForce              binance.TimeInForceType `json:"timeInForce"`
	Type                     binance.OrderType       `json:"type"`
	Side                     binance.SideType        `json:"side"`
	Fills                    []*binance.Fill         `json:"fills"`
	MarginBuyBorrowAmount    string                  `json:"marginBuyBorrowAmount"`
	MarginBuyBorrowAsset     string                  `json:"marginBuyBorrowAsset"`
}

func (c *SpotClient) Trade(ctx context.Context, req *TradeReq) (*TradeResp, error) {

	// check symbol quantity filters
	/*
		名义价值过滤器(NOTIONAL)定义了订单在一个交易对上可以下单的名义价值区间.

		applyMinToMarket 定义了 minNotional 是否适用于市价单(MARKET)
		applyMaxToMarket 定义了 maxNotional 是否适用于市价单(MARKET).

		要通过此过滤器, 订单的名义价值 (单价 x 数量, price * quantity) 需要满足如下条件:

		price * quantity <= maxNotional
		price * quantity >= minNotional
		对于市价单(MARKET), 用于计算的价格采用的是在 avgPriceMins 定义的时间之内的平均价.
		如果 avgPriceMins 为 0, 则采用最新的价格.
	*/
	exchangeInfo, err := c.binanceSpotClient.NewExchangeInfoService().Symbol(req.Symbol).Do(ctx)
	if err != nil {
		return nil, err
	}

	var minNotional string

	if len(exchangeInfo.Symbols) == 0 || exchangeInfo.Symbols == nil {
		return nil, errors.New("exchangeInfo.Symbols fail")
	}

	for _, filter := range exchangeInfo.Symbols[0].Filters {
		if filter["filterType"] == "MIN_NOTIONAL" {
			minNotional = filter["minNotional"].(string)
		}
	}

	baseAsset := exchangeInfo.Symbols[0].BaseAsset   // LUNC
	quoteAsset := exchangeInfo.Symbols[0].QuoteAsset // BUSD
	var quoteQuantity string

	if req.Side == "BUY" {
		fMinNotional, _ := strconv.ParseFloat(minNotional, 64)
		fQuantity, _ := strconv.ParseFloat(req.Quantity, 64)
		if fQuantity < fMinNotional {
			return nil, errors.New(fmt.Sprintf("quantity less than min notional : %s %s < %s %s", req.Quantity, quoteAsset, minNotional, quoteAsset))
		}
		quoteQuantity = req.Quantity

	} else {
		res, err := c.binanceSpotClient.NewAveragePriceService().Symbol(req.Symbol).Do(ctx)
		if err != nil {
			return nil, err
		}
		fMinNotional, _ := strconv.ParseFloat(minNotional, 64)
		fQuantity, _ := strconv.ParseFloat(req.Quantity, 64)
		fPrice, _ := strconv.ParseFloat(res.Price, 64)
		notional := fQuantity * fPrice
		if notional < fMinNotional {
			return nil, errors.New(fmt.Sprintf("quantity * avgPrice less than min notional : %f * %f %s/%s < %f %s", fQuantity, fPrice, baseAsset, quoteAsset, fMinNotional, quoteAsset))
		}

		quoteQuantity = strconv.FormatFloat(notional, 'f', 8, 64)
	}

	order, err := c.binanceSpotClient.NewCreateOrderService().Symbol(req.Symbol).
		Side(req.Side).Type(binance.OrderTypeMarket).NewClientOrderID(req.NewClientOrderId).
		QuoteOrderQty(quoteQuantity).Do(ctx)

	if err != nil {
		return nil, err
	}

	var resp TradeResp

	copier.Copy(&resp, order)

	return &resp, nil
}

type GetOrderReq struct {
	Symbol  string `json:"symbol"`
	OrderId int64  `json:"orderId"`
}

type GetOrderResp struct {
	Symbol                   string                  `json:"symbol"`
	OrderID                  int64                   `json:"orderId"`
	OrderListId              int64                   `json:"orderListId"`
	ClientOrderID            string                  `json:"clientOrderId"`
	Price                    string                  `json:"price"`
	OrigQuantity             string                  `json:"origQty"`
	ExecutedQuantity         string                  `json:"executedQty"`
	CummulativeQuoteQuantity string                  `json:"cummulativeQuoteQty"`
	Status                   binance.OrderStatusType `json:"status"`
	TimeInForce              binance.TimeInForceType `json:"timeInForce"`
	Type                     binance.OrderType       `json:"type"`
	Side                     binance.SideType        `json:"side"`
	StopPrice                string                  `json:"stopPrice"`
	IcebergQuantity          string                  `json:"icebergQty"`
	Time                     int64                   `json:"time"`
	UpdateTime               int64                   `json:"updateTime"`
	IsWorking                bool                    `json:"isWorking"`
	IsIsolated               bool                    `json:"isIsolated"`
	OrigQuoteOrderQuantity   string                  `json:"origQuoteOrderQty"`
}

func (c *SpotClient) GetOrder(ctx context.Context, req *GetOrderReq) (*GetOrderResp, error) {
	order, err := c.binanceSpotClient.NewGetOrderService().Symbol(req.Symbol).OrderID(req.OrderId).Do(ctx)
	if err != nil {
		return nil, err
	}

	var resp GetOrderResp

	copier.Copy(&resp, order)

	return &resp, err
}

type OrderListReq struct {
	Symbol    string `json:"symbol"`
	StartTime int64  `json:"startTime,omitempty"`
	EndTime   int64  `json:"endTime,omitempty"`
	Limit     int    `json:"limit"`
}

type OrderListResp struct {
	Data []*binance.Order `json:"data"`
}

func (c *SpotClient) OrderList(ctx context.Context, req *OrderListReq) (*OrderListResp, error) {

	list, err := c.binanceSpotClient.NewListOrdersService().Symbol(req.Symbol).Limit(req.Limit).Do(ctx)

	if err != nil {
		return nil, err
	}

	var resp []*binance.Order
	for _, li := range list {
		resp = append(resp, li)
	}

	return &OrderListResp{Data: resp}, nil
}

func (c *SpotClient) HangOrderList(ctx context.Context) (*OrderListResp, error) {

	list, err := c.binanceSpotClient.NewListOpenOrdersService().Do(ctx)

	if err != nil {
		return nil, err
	}

	var resp []*binance.Order
	for _, li := range list {
		resp = append(resp, li)
	}

	return &OrderListResp{Data: resp}, nil
}

type CancelReq struct {
	Symbol  string `json:"symbol"`
	OrderId int64  `json:"orderId"`
}

type CancelResp struct {
	Symbol              string `json:"symbol"`
	OrigClientOrderId   string `json:"origClientOrderId"`
	OrderId             int    `json:"orderId"`
	OrderListId         int    `json:"orderListId"`
	ClientOrderId       string `json:"clientOrderId"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
}

func (c *SpotClient) CancelOrder(ctx context.Context, req *CancelReq) (*CancelResp, error) {
	order, err := c.binanceSpotClient.NewCancelOrderService().Symbol(req.Symbol).OrderID(req.OrderId).Do(ctx)
	if err != nil {
		return nil, err
	}
	var resp CancelResp
	copier.Copy(&resp, order)

	return &resp, nil
}

type WithdrawReq struct {
	Coin            string `json:"coin"`
	Address         string `json:"address"`
	AddressTag      string `json:"addressTag"`
	Amount          string `json:"amount"`
	WithdrawOrderId string `json:"withdrawOrderId"`
}

type WithdrawResp struct {
	Id string
}

func (c *SpotClient) Withdraw(ctx context.Context, req *WithdrawReq) (*WithdrawResp, error) {
	withdrawId, err := c.binanceSpotClient.NewCreateWithdrawService().
		Coin(req.Coin).Address(req.Address).AddressTag(req.AddressTag).
		Amount(req.Amount).WithdrawOrderID(req.WithdrawOrderId).Do(ctx)

	if err != nil {
		return nil, err
	}

	return &WithdrawResp{Id: withdrawId.ID}, nil

}

type TradeFeeReq struct {
	Symbol string `json:"symbol,omitempty"`
}

type TradeFeeResp struct {
	Data []*binance.TradeFeeDetails `json:"data"`
}

func (c *SpotClient) TradeFee(ctx context.Context, req *TradeFeeReq) (*TradeFeeResp, error) {
	feeDetails, err := c.binanceSpotClient.NewTradeFeeService().Symbol(req.Symbol).Do(ctx)
	if err != nil {
		return nil, err
	}

	return &TradeFeeResp{Data: feeDetails}, nil
}

type WithdrawHistoryReq struct {
	Coin            string `json:"coin"`
	WithdrawOrderId string `json:"withdrawOrderId"`
}

type WithdrawHistoryResp struct {
	Data []*binance.Withdraw `json:"data"`
}

func (c *SpotClient) WithdrawHistory(ctx context.Context, req *WithdrawHistoryReq) (*WithdrawHistoryResp, error) {
	withdrawList, err := c.binanceSpotClient.NewListWithdrawsService().Coin(req.Coin).
		WithdrawOrderId(req.WithdrawOrderId).Do(ctx)
	if err != nil {
		return nil, err
	}

	var resp []*binance.Withdraw
	for _, withdraw := range withdrawList {
		resp = append(resp, withdraw)
	}

	return &WithdrawHistoryResp{Data: resp}, nil
}

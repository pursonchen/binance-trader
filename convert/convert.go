package convert

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/pursonchen/go-binance/v2"
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
	Symbol   string           `json:"symbol"`
	Side     binance.SideType `json:"side"` // BUY SELL
	Type     string           `json:"type"`
	Quantity string           `json:"quantity"`
}

type TradeResp struct {
	Symbol              string `json:"symbol"`
	OrderId             int    `json:"orderId"`
	OrderListId         int    `json:"orderListId"`
	ClientOrderId       string `json:"clientOrderId"`
	TransactTime        int64  `json:"transactTime"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	StrategyId          int    `json:"strategyId"`
	StrategyType        int    `json:"strategyType"`
	Fills               []struct {
		Price           string `json:"price"`
		Qty             string `json:"qty"`
		Commission      string `json:"commission"`
		CommissionAsset string `json:"commissionAsset"`
		TradeId         int    `json:"tradeId"`
	} `json:"fills"`
}

func (c *SpotClient) Trade(ctx context.Context, req *TradeReq) (*TradeResp, error) {

	order, err := c.binanceSpotClient.NewCreateOrderService().Symbol(req.Symbol).
		Side(req.Side).Type(binance.OrderTypeMarket).TimeInForce(binance.TimeInForceTypeGTC).
		Quantity(req.Quantity).Do(ctx)

	if err != nil {
		return nil, err
	}

	var resp *TradeResp

	copier.Copy(&resp, order)

	return resp, nil
}

type GetOrderReq struct {
	Symbol  string `json:"symbol"`
	OrderId int64  `json:"orderId"`
}

type GetOrderResp struct {
	Symbol              string `json:"symbol"`
	OrderId             int64  `json:"orderId"`
	OrderListId         int64  `json:"orderListId"`
	ClientOrderId       string `json:"clientOrderId"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	StopPrice           string `json:"stopPrice"`
	IcebergQty          string `json:"icebergQty"`
	Time                int64  `json:"time"`
	UpdateTime          int64  `json:"updateTime"`
	IsWorking           bool   `json:"isWorking"`
	OrigQuoteOrderQty   string `json:"origQuoteOrderQty"`
}

func (c *SpotClient) getOrder(ctx context.Context, req *GetOrderReq) (*GetOrderResp, error) {
	order, err := c.binanceSpotClient.NewGetOrderService().Symbol(req.Symbol).OrderID(req.OrderId).Do(ctx)
	if err != nil {
		return nil, err
	}
	var resp *GetOrderResp
	copier.Copy(&resp, order)

	return resp, err
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
	var resp *CancelResp
	copier.Copy(&resp, order)

	return resp, nil
}

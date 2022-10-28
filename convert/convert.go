package convert

import (
	"context"
	"errors"
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

func (c *SpotClient) EstQuote(ctx context.Context, params *EstQuoteReq) (*EstQuoteResp, error) {

	// 目前支持稳定币兑换
	if params.BaseCcy != StableCoins[0] && params.QuoteCcy != StableCoins[0] {
		return nil, errors.New("coin not support")
	}

	res, err := c.binanceSpotClient.NewAveragePriceService().Symbol(params.BaseCcy + params.QuoteCcy).Do(ctx)

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
	Symbol    string  `json:"symbol"`
	Side      string  `json:"side"`
	Type      string  `json:"type"`
	Quantity  float64 `json:"quantity"`
	Timestamp int64   `json:"timestamp"`
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

func (c *SpotClient) Trade(ctx context.Context, param *TradeReq) (*TradeResp, error) {
	return &TradeResp{}, nil
}

func (c *SpotClient) OrderList(ctx context.Context) {

}

func (c *SpotClient) CancelOrderList(ctx context.Context) {

}

func (c *SpotClient) CancelOrder(ctx context.Context) {

}

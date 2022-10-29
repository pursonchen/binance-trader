package convert

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/pursonchen/go-binance/v2"
	"github.com/smartystreets/goconvey/convey"
	"log"
	"os"
	"testing"
)

var sCli *SpotClient

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	//binance.UseTestnet = true
	binanceClient := binance.NewProxiedClient(os.Getenv("BINANCE_API_KEY"), os.Getenv("BINANCE_SECRET_KEY"), os.Getenv("PROXY_URL"))
	sCli = NewSpotClient(binanceClient)
	os.Exit(m.Run())
}

func TestEstQuote(t *testing.T) {

	convey.Convey("TestEstQuote", t, func(convCtx convey.C) {
		resp, err := sCli.EstQuote(context.Background(), &EstQuoteReq{
			BaseCcy:  "FIL",
			QuoteCcy: "USDT",
		})

		if err != nil {
			t.Fatalf(`test EstQuote fail %s`, err.Error())
		}

		log.Printf("result:%+v", resp)
		convCtx.So(resp.Price, convey.ShouldNotBeEmpty)
	})

}

func TestOrderList(t *testing.T) {
	convey.Convey("TestOrderList", t, func(convCtx convey.C) {
		resp, err := sCli.OrderList(context.Background(), &OrderListReq{
			Symbol:    "EOSBTC",
			StartTime: 0,
			EndTime:   0,
			Limit:     500,
		})

		if err != nil {
			t.Fatalf(`test TestOrderList fail %s`, err.Error())
		}
		for _, li := range resp.Data {
			log.Printf("result:%+v", li)
		}
		convCtx.So(len(resp.Data), convey.ShouldBeGreaterThan, 0)
	})
}

func TestGetOrder(t *testing.T) {
	convey.Convey("TestGetOrder", t, func(convCtx convey.C) {
		resp, err := sCli.getOrder(context.Background(), &GetOrderReq{
			Symbol:  "EOSBTC",
			OrderId: 659854195,
		})

		if err != nil {
			t.Fatalf(`test TestGetOrder fail %s`, err.Error())
		}

		log.Printf("result:%+v", resp)

		convCtx.So(resp.CummulativeQuoteQuantity, convey.ShouldEqual, "0.00384000")
	})
}

func TestTrade(t *testing.T) {
	convey.Convey("TestTrade", t, func(convCtx convey.C) {
		resp, err := sCli.Trade(context.Background(), &TradeReq{
			Symbol:   "LUNCBUSD",
			Side:     "BUY",
			Quantity: "10",
		})

		if err != nil {
			t.Fatalf(`test TestTrade BUY fail %s`, err.Error())
		}

		log.Printf("buy result:%+v", resp)

		sellQuantity := resp.ExecutedQuantity

		convCtx.So(resp.Status, convey.ShouldEqual, "FILLED")

		resp, err = sCli.Trade(context.Background(), &TradeReq{
			Symbol:   "LUNCBUSD",
			Side:     "SELL",
			Quantity: sellQuantity,
		})

		if err != nil {
			t.Fatalf(`test TestTrade SELL fail %s`, err.Error())
		}

		log.Printf("sell result:%+v", resp)

		convCtx.So(resp.Status, convey.ShouldEqual, "FILLED")

	})
}

func TestTradedFee(t *testing.T) {
	convey.Convey("TestTradedFee", t, func(convCtx convey.C) {
		resp, err := sCli.TradeFee(context.Background(), &TradeFeeReq{
			Symbol: "EOSBTC",
		})

		if err != nil {
			t.Fatalf(`test TestGetOrder fail %s`, err.Error())
		}

		log.Printf("result:%+v", resp.Data[0])

		convCtx.So(resp.Data[0].Symbol, convey.ShouldEqual, "EOSBTC")
	})
}

func TestWithdraw(t *testing.T) {
	convey.Convey("TestWithdraw", t, func(convCtx convey.C) {
		resp, err := sCli.Withdraw(context.Background(), &WithdrawReq{
			Coin:            "EOS",
			Address:         "pursonpurson",
			AddressTag:      "",
			Amount:          "10",
			WithdrawOrderId: "test",
		})

		if err != nil {
			t.Fatalf(`test TestWithdraw fail %s`, err.Error())
		}

		log.Printf("result:%+v", resp.Id)

		convCtx.So(resp.Id, convey.ShouldNotBeEmpty)
	})
}

func TestWithdrawList(t *testing.T) {
	convey.Convey("TestWithdrawList", t, func(convCtx convey.C) {
		resp, err := sCli.WithdrawHistory(context.Background(), &WithdrawHistoryReq{
			Coin:            "EOS",
			WithdrawOrderId: "",
		})

		if err != nil {
			t.Fatalf(`test TestWithdrawList fail %s`, err.Error())
		}

		log.Printf("result:%+v", resp.Data[0])

		convCtx.So(resp.Data[0].Coin, convey.ShouldEqual, "EOS")
	})
}

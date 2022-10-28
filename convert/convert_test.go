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

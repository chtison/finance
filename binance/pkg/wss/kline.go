package wss

import (
	"log"
	"os"
	"os/signal"

	"github.com/chtison/fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

const URL = "wss://stream.binance.com:9443"

type KlineEvent struct {
	Type   string `json:"e"`
	Time   int64  `json:"E"`
	Symbol string `json:"s"`
	Kline  Kline  `json:"k"`
}

type Kline struct {
	StartTime            int64  `json:"t"`
	EndTime              int64  `json:"T"`
	Symbol               string `json:"s"`
	Interval             string `json:"i"`
	FirstTradeID         int64  `json:"f"`
	LastTradeID          int64  `json:"L"`
	Open                 string `json:"o"`
	Close                string `json:"c"`
	High                 string `json:"h"`
	Low                  string `json:"l"`
	Volume               string `json:"v"`
	TradeNum             int64  `json:"n"`
	IsFinal              bool   `json:"x"`
	QuoteVolume          string `json:"q"`
	ActiveBuyVolume      string `json:"V"`
	ActiveBuyQuoteVolume string `json:"Q"`
}

func GetKlineStream(symbol string, interval string, logger *log.Logger) (<-chan *KlineEvent, <-chan error, chan bool, error) {
	u := URL + "/ws/" + symbol + "@kline_" + interval

	if logger != nil {
		logger.Printf("Connecting to %s", u)
	}

	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	if logger != nil {
		logger.Printf("Connected to %s", u)
	}

	msgs := make(chan *KlineEvent)
	errs := make(chan error)
	stop := make(chan bool, 1)

	go func() {
		<-stop
		errs = nil
		c.Close()
		stop <- true
	}()

	go func() {
		defer func() { stop <- true }()
		for {
			var msg *KlineEvent
			if err := c.ReadJSON(&msg); err != nil {
				if errs != nil {
					errs <- err
				}
				return
			}
			msgs <- msg
		}
	}()

	return msgs, errs, stop, nil
}

func getKlineStreamCommand() *cobra.Command {
	verbose := false
	runE := func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return cmd.Help()
		}
		var logger *log.Logger
		if verbose {
			logger = log.New(os.Stdout, "", log.LstdFlags)
		}
		msgs, errs, stop, err := GetKlineStream(args[0], args[1], logger)
		if err != nil {
			return err
		}

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		defer signal.Stop(c)
		for {
			select {
			case <-c:
				stop <- true
				<-stop
				return nil
			case err := <-errs:
				if err != nil {
					return err
				}
			case msg := <-msgs:
				fmt.Println(msg)
			}
		}
		return nil
	}
	cmd := &cobra.Command{
		Use:   "kline SYMBOL INTERVAL [FLAGS]",
		Short: "Connect to the kline/candlestick stream for the given SYMBOL and INTERVAL",
		Long: `Connect to the kline/candlestick stream for the given SYMBOL and INTERVAL

https://github.com/binance-exchange/binance-official-api-docs/blob/master/web-socket-streams.md#klinecandlestick-streams
		`,
		RunE:                  runE,
		DisableFlagsInUseLine: true,
	}
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Set verbose output")
	return cmd
}

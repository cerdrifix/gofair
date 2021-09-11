package streaming

import (
	"github.com/sirupsen/logrus"

	"github.com/belmegatron/gofair/streaming/models"
)


type MarketStream struct {

	listener             *Listener
	log                  *logrus.Logger

	IncomingMarketData   chan MarketBook
	SubscriptionResponse chan MarketSubscriptionResponse
	
	Cache                map[string]MarketCache
	InitialClk           string
	Clk                  string
}

func NewMarketStream(listener *Listener, log *logrus.Logger, marketUpdates *chan MarketBook) *MarketStream {
	marketStream := new(MarketStream)
	marketStream.log = log
	marketStream.IncomingMarketData = *marketUpdates
	marketStream.Cache = make(map[string]MarketCache)
	return marketStream
}

func (ms *MarketStream) Subscribe(marketFilter *models.MarketFilter, marketDataFilter *models.MarketDataFilter) {

	request := new(models.MarketSubscriptionMessage)
	request.SetID(ms.listener.uid)
	ms.listener.uid++
	request.MarketFilter = marketFilter
	request.MarketDataFilter = marketDataFilter

	ms.listener.marketSubscriptionRequest <- *request
}

func getMarketIDs(mcm models.MarketChangeMessage) []string {
	marketIDs := make([]string, 0)
	for _, marketChange := range mcm.Mc {
		marketIDs = append(marketIDs, marketChange.ID)
	}
	return marketIDs
}

func (ms *MarketStream) OnSubscribe(changeMessage models.MarketChangeMessage) {

	marketIDs := getMarketIDs(changeMessage)

	ms.log.WithFields(logrus.Fields{
		"marketIDs": marketIDs,
	}).Debug("Subscribed to Betfair Market Changes")

	for _, marketChange := range changeMessage.Mc {

		marketCache := CreateMarketCache(&changeMessage, marketChange)
		ms.Cache[marketChange.ID] = *marketCache
	}
}

func (ms *MarketStream) OnResubscribe(changeMessage models.MarketChangeMessage) {
	marketIDs := getMarketIDs(changeMessage)
	ms.log.WithFields(logrus.Fields{
		"marketIDs": marketIDs,
	}).Debug("Resubscribed to Betfair Market Changes")
}

func (ms *MarketStream) OnHeartbeat(changeMessage models.MarketChangeMessage) {
	ms.log.Debug("Heartbeat")
}

func (ms *MarketStream) OnUpdate(changeMessage models.MarketChangeMessage) {

	if ms.InitialClk == "" {
		ms.InitialClk = changeMessage.Clk
	}

	ms.Clk = changeMessage.Clk

	for _, marketChange := range changeMessage.Mc {

		if marketCache, ok := ms.Cache[marketChange.ID]; ok {
			marketCache.UpdateCache(&changeMessage, marketChange)
			ms.IncomingMarketData <- marketCache.Snap()
		} else {
			marketCache := CreateMarketCache(&changeMessage, marketChange)
			ms.Cache[marketChange.ID] = *marketCache
			ms.IncomingMarketData <- marketCache.Snap()
		}
	}
}

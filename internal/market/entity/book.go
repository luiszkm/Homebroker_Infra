package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Order           []*Order
	Transactions     []*Transaction
	OrdersChanel    chan *Order
	OrdersChanelOut chan *Order
	Wg              *sync.WaitGroup
}

func NewBook(orderChanel chan *Order, orderChanelOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Order:           []*Order{},
		Transactions:     []*Transaction{},
		OrdersChanel:    orderChanel,
		OrdersChanelOut: orderChanelOut,
		Wg:              wg,
	}
}

func (b *Book) Trade() {
	buyOrders := make(map[string]*OrderQueue)
	sellOrders := make(map[string]*OrderQueue)

	//buyOrders := NewOrderQueue()
	//	sellOrders := NewOrderQueue()

	//heap.Init(buyOrders)
	//heap.Init(sellOrders)

	for order := range b.OrdersChanel {
		asset:= order.Asset.ID
		if buyOrders[asset] == nil {
			buyOrders[asset] = NewOrderQueue()
			heap.Init(buyOrders[asset])
		}
		if sellOrders[asset] == nil {
			sellOrders[asset] = NewOrderQueue()
			heap.Init(sellOrders[asset])
		}

		if order.OrderType == "BUY" {
			buyOrders[asset].Push(order)
			if sellOrders[asset].Len() > 0 &&
				sellOrders[asset].Orders[0].Price <= order.Price {
				sellOrder := sellOrders[asset].Pop().(*Order)
				if sellOrder.PendingShares > 0 {
					trnsaction := NewTransaction(sellOrder, order, order.Shares, sellOrder.Price)
					b.AddTransaction(trnsaction, b.Wg)
					sellOrder.Transactions = append(sellOrder.Transactions, trnsaction)
					order.Transactions = append(order.Transactions, trnsaction)
					b.OrdersChanelOut <- sellOrder
					b.OrdersChanelOut <- order
					if sellOrder.PendingShares > 0 {
						sellOrders[asset].Push(sellOrders)
					}
				}
			}
		} else if order.OrderType == "SELL" {
			sellOrders[asset].Push(order)
			if buyOrders[asset].Len() > 0 && buyOrders[asset].Orders[0].Price >= order.Price {
				buyOrder := buyOrders[asset].Pop().(*Order)
				if buyOrder.PendingShares > 0 {
					transaction := NewTransaction(order, buyOrder, order.Shares, buyOrder.Price)
					b.AddTransaction(transaction, b.Wg)
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					buyOrder.Transactions = append(order.Transactions, transaction)
					b.OrdersChanelOut <- buyOrder
					b.OrdersChanelOut <- order
					if buyOrder.PendingShares > 0 {
						buyOrders[asset].Push(buyOrder)
					}
				}
			}
		}
	}
}

func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done()
	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares

	minShares := sellingShares

	if buyingShares < minShares {
		minShares = buyingShares
	}

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.AddSellOrderPendingShares(-minShares)

	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.ID, minShares)
	transaction.AddBuyOrderPendingShares(-minShares)

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)

	transaction.CloseBuyOrder()
	transaction.CloseSellOrder()

	b.Transactions = append(b.Transactions, transaction)
}

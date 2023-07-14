package entity

import (
	"github.com/google/uuid"
	"time"
	//"container/heap"
)

type Transaction struct {
	ID           string
	SellingOrder *Order
	BuyingOrder  *Order
	Shares       int
	Price        float64
	Total        float64
	DataTime     time.Time
}

func NewTransaction(sellingOrder *Order, buyOrder *Order, shares int, price float64) *Transaction {
	total := float64(shares) * price
	return &Transaction{
		ID:           uuid.New().String(),
		SellingOrder: sellingOrder,
		BuyingOrder:  buyOrder,
		Shares:       shares,
		Price:        price,
		Total:        total,
		DataTime:     time.Now(),
	}
}

func (t *Transaction) CalculateTotal(shares int, price float64) {
	t.Total = float64(shares) * price
}

func (t *Transaction) CloseBuyOrder() {
	if t.BuyingOrder.PendingShares == 0 {
		t.BuyingOrder.Status = "CLOSED"
	}
}
func (t *Transaction) CloseSellOrder() {
	if t.SellingOrder.PendingShares == 0 {
		t.SellingOrder.Status = "CLOSED"
	}
}

func (t *Transaction) AddBuyOrderPendingShares(shates int) {
	t.BuyingOrder.PendingShares += shates
}
func (t *Transaction) AddSellOrderPendingShares(shates int) {
	t.SellingOrder.PendingShares += shates
}

package entity

type Investor struct {
	ID string
	Nmae string
	AssetPosition []*InvestorAssetPosition
}
func NewInvestor (id string) *Investor {
	return &Investor{
		ID: id,
		AssetPosition: []*InvestorAssetPosition{},
	}
}
func(i *Investor) AddAssetPosition(assetPostion *InvestorAssetPosition){
	i.AssetPosition = append(i.AssetPosition, assetPostion)
}

func (i *Investor) UpdateAssetPosition(assetID string, qtdShares int) {
	assetPosition:= i.GetAssetPosition(assetID)
	if assetPosition == nil {
		i.AssetPosition = append(i.AssetPosition, NewInvestorAssetPosition(assetID, qtdShares))
	}else {
		assetPosition.Shares += qtdShares
	}
}

func (i *Investor) GetAssetPosition(assetID string) *InvestorAssetPosition {
	for _, assetPosition := range i.AssetPosition {
		if assetPosition.AssetID == assetID {
			return assetPosition
		}
	}
	return nil
}

type InvestorAssetPosition struct {
	AssetID string
	Shares int
}

func NewInvestorAssetPosition(assetID string, shares int) *InvestorAssetPosition {
	return &InvestorAssetPosition{
		AssetID: assetID,
		Shares: shares,
	}
}
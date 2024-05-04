package socrata

type Year int
type Town string

type RealEstateTransaction struct {
	Town        Town
	SaleYear    Year
	SalesRatio  float64
	SalesVolume float64
}

func NewTransaction(town Town, year Year, ratio, volume float64) *RealEstateTransaction {
	return &RealEstateTransaction{
		Town:        town,
		SaleYear:    year,
		SalesRatio:  ratio,
		SalesVolume: volume,
	}
}

func (r *RealEstateTransaction) GetSalesRatio() float64 {
	return r.SalesRatio
}

func (r *RealEstateTransaction) GetSalesVolume() float64 {
	return r.SalesVolume
}

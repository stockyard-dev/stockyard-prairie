package server

type Tier string

const (
	TierFree Tier = "free"
	TierPro  Tier = "pro"
)

type Limits struct {
	Tier        Tier
	Description string
}

func LimitsFor(tier string) Limits {
	if tier == "pro" {
		return Limits{Tier: TierPro, Description: "Unlimited sites and pageviews"}
	}
	return Limits{Tier: TierFree, Description: "1 site, 10k pageviews/mo"}
}

func (l Limits) IsPro() bool {
	return l.Tier == TierPro
}

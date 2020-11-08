package analyzer

type Candidate struct {
	isSplit      bool
	swarms       []string
	placementOne map[string]bool
	placementTwo map[string]bool
}

func (c *Candidate) IsSplit() bool                    { return c.isSplit }
func (c *Candidate) GetSwarmIDs() []string            { return c.swarms }
func (c *Candidate) GetPlacementOne() map[string]bool { return c.placementOne }
func (c *Candidate) GetPlacementTwo() map[string]bool { return c.placementTwo }

type swarmInfo struct {
	SwarmID   string
	SwarmSize int
	FitScore  float64
}

type swarmInfoList struct {
	Infos []swarmInfo
}

func (sl *swarmInfoList) Len() int { return len(sl.Infos) }
func (sl *swarmInfoList) Less(i, j int) bool {
	return sl.Infos[i].FitScore < sl.Infos[j].FitScore
}
func (sl *swarmInfoList) Swap(i, j int) {
	sl.Infos[i], sl.Infos[j] = sl.Infos[j], sl.Infos[i]
}

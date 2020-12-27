package analyzer

type Candidate struct {
	transfererID string
	transfereeID string
	transferSize int
}

func (c *Candidate) GetTransfererID() string { return c.transfererID }
func (c *Candidate) GetTransfereeID() string { return c.transfereeID }
func (c *Candidate) GetTransferSize() int    { return c.transferSize }

type swarmDistanceInfo struct {
	dataspace string
	distance  int
}

type swarmDistancesSlice []*swarmDistanceInfo

func (ss *swarmDistancesSlice) Len() int { return len(*ss) }
func (ss *swarmDistancesSlice) Less(i, j int) bool {
	return (*ss)[i].distance < (*ss)[j].distance
}
func (ss *swarmDistancesSlice) Swap(i, j int) {
	slice := (*ss)
	slice[i], slice[j] = slice[j], slice[i]
}

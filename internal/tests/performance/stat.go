package performance

import "time"

type Stat interface {
	StartTime() time.Time
	EndTime() time.Time
	EndState() string
	Success() bool
}

type BasicStat struct {
  start time.Time
  end time.Time
  state string
  success bool
}

func (bs *BasicStat) StartTime() time.Time { return bs.start }
func (bs *BasicStat) EndTime() time.Time { return bs.end }
func (bs *BasicStat) EndState() string { return bs.state }
func (bs *BasicStat) Success() bool { return bs.success }

func (bs *BasicStat) SetStartTime(t time.Time) { bs.start = t }
func (bs *BasicStat) SetEndTime(t time.Time) { bs.end = t }
func (bs *BasicStat) SetState(s string) { bs.state = s }
func (bs *BasicStat) SetSuccess(s bool) { bs.success = s }

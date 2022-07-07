package ctrl

// StallState keeps track of when a pipeline stage is stalled
// and when data can be moved from one stage to the next.
// First the `Stall` flag should be calculated for a stage to
// indicate data can't be moved forward. Then the `Calc*` methods
// can be called from
type StallState struct {
	// if true, we can't accept new data, nor pass on our data
	Stall bool

	// if true, reset the pipeline stage so its data is invalid
	Reset bool

	// should we latch our outputs into the next stage's inputs
	En bool

	// are we able to emit valid data to the next stage
	Valid bool

	// are we ready to accept data from the previous stage
	Ready bool
}

func (ss *StallState) CalcBeginning(next *StallState) {
	ss.En = next.Ready
	ss.Valid = !ss.Reset && !ss.Stall
	ss.Ready = !ss.Stall && ss.En
}

func (ss *StallState) CalcMiddle(prev *StallState, next *StallState) {
	ss.En = prev.Valid && next.Ready
	ss.Valid = !ss.Reset && !ss.Stall && prev.Valid
	ss.Ready = !ss.Stall && ss.En
}

func (ss *StallState) CalcEnding(prev *StallState) {
	ss.En = prev.Valid
	ss.Valid = !ss.Reset && !ss.Stall && prev.Valid
	ss.Ready = !ss.Stall && ss.En
}

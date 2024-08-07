package fsm

// EventName return event with src dst state
func (f *FSM) EventName(dst string) (eventName string, ok bool) {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	for key, transition := range f.transitions {
		if key.src == f.current && transition == dst {
			return key.event, true
		}
	}
	return "", false
}

// AvailableDstStates return all availabe dst state
func (f *FSM) AvailableDstStates() (states []string) {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	states = make([]string, 0)
	for key, transition := range f.transitions {
		if key.src == f.current {
			states = append(states, transition)
		}
	}
	return states
}

// AvailableSrcStates return all availabe src state
func (f *FSM) AvailableSrcStates() (states []string) {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	states = make([]string, 0)
	for key, transition := range f.transitions {
		if transition == f.current {
			states = append(states, key.src)
		}
	}
	return states
}

// CanConvertTo return can src to dst
func (f *FSM) CanConvertTo(dstState string) (ok bool) {
	_, ok = f.EventName(dstState)
	return ok
}

// IsReverseOrder return true if it is a reverse order event(the dest status in src status)
func (f *FSM) IsReverseOrder(beforStatus string) (ok bool) {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	for key, transition := range f.transitions {
		if transition == f.current { // find src where current as dst
			if key.src == beforStatus {
				return true
			}
		}
	}
	return false
}

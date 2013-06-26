package peek

import (
  "sort"
)

type eventSorter struct {
	events []*Event
	by     func(e1, e2 *Event) bool
}

func (s *eventSorter) Len() int {
	return len(s.events)
}

func (s *eventSorter) Less(i, j int) bool {
	return s.by(s.events[i], s.events[j])
}

func (s *eventSorter) Swap(i, j int) {
	s.events[i], s.events[j] = s.events[j], s.events[i]
}

type By func(e1, e2 *Event) bool

func (by By) Sort(events []*Event) {
	es := &eventSorter{
		events: events,
		by:     by,
	}

	sort.Sort(es)
}

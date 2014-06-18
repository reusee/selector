package selector

import "reflect"

type Selector struct {
	cases     []reflect.SelectCase
	selecting []reflect.SelectCase
	toSend    []func() interface{}
	cbs       []interface{}
}

type Case struct {
	index    int
	selector *Selector
	disabled bool
}

func New() *Selector {
	return &Selector{}
}

var emptyCase = reflect.SelectCase{
	Dir: reflect.SelectRecv,
}

func (s *Selector) Add(ch interface{}, cb interface{}, toSend func() interface{}) *Case {
	dir := reflect.SelectRecv
	if toSend != nil {
		dir = reflect.SelectSend
	}
	selectCase := reflect.SelectCase{
		Dir:  dir,
		Chan: reflect.ValueOf(ch),
	}
	s.cases = append(s.cases, selectCase)
	s.selecting = append(s.selecting, selectCase)
	s.toSend = append(s.toSend, toSend)
	s.cbs = append(s.cbs, cb)
	return &Case{
		index:    len(s.cases) - 1,
		selector: s,
	}
}

func (s *Selector) Select() {
	for i, send := range s.toSend {
		if send != nil {
			s.selecting[i].Send = reflect.ValueOf(send())
		}
	}
	n, recv, ok := reflect.Select(s.selecting)
	if s.selecting[n].Dir == reflect.SelectRecv {
		if ok {
			if s.cbs[n] != nil {
				switch cb := s.cbs[n].(type) {
				case func():
					cb()
				case func(interface{}):
					cb(recv.Interface())
				default:
					panic("unknown callback")
				}
			}
		}
	} else {
		if s.cbs[n] != nil {
			(s.cbs[n].(func()))()
		}
	}
}

func (c *Case) Disable() {
	if !c.disabled {
		c.selector.selecting[c.index] = emptyCase
		c.disabled = true
	}
}

func (c *Case) Enable() {
	if c.disabled {
		c.selector.selecting[c.index] = c.selector.cases[c.index]
		c.disabled = false
	}
}

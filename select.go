package selector

import "reflect"

type Selector struct {
	cases     []reflect.SelectCase
	selecting []reflect.SelectCase
	toSend    []*reflect.Value
	cbs       []interface{}
}

type Case struct {
	index    int
	selector *Selector
}

func New() *Selector {
	return &Selector{}
}

var emptyCase = reflect.SelectCase{
	Dir: reflect.SelectRecv,
}

func (s *Selector) Add(ch interface{}, cb interface{}, toSend interface{}) *Case {
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
	var toSendCb *reflect.Value
	if toSend != nil {
		v := reflect.ValueOf(toSend)
		toSendCb = &v
	}
	s.toSend = append(s.toSend, toSendCb)
	s.cbs = append(s.cbs, cb)
	return &Case{
		index:    len(s.cases) - 1,
		selector: s,
	}
}

func (s *Selector) Select() {
	for i, send := range s.toSend {
		if send != nil {
			ret := send.Call(nil)
			s.selecting[i].Send = ret[0]
		}
	}
	n, recv, ok := reflect.Select(s.selecting)
	if s.selecting[n].Dir == reflect.SelectRecv {
		if ok {
			if s.cbs[n] != nil {
				switch cb := s.cbs[n].(type) {
				case func():
					cb()
				default:
					reflect.ValueOf(cb).Call([]reflect.Value{
						recv,
					})
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
	c.selector.selecting[c.index] = emptyCase
}

func (c *Case) Enable() {
	c.selector.selecting[c.index] = c.selector.cases[c.index]
}

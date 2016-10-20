package promise

type ResolveFunc func(interface{})
type RejectFunc func(error)
type PromiseFunc func(ResolveFunc, RejectFunc)

type Promise struct {
	async   PromiseFunc
	resolve ResolveFunc
	reject  RejectFunc
}

func (p *Promise) Then(fn ResolveFunc) *Promise {
	p.resolve = fn
	return p
}

func (p *Promise) Catch(fn RejectFunc) *Promise {
	p.reject = fn
	return p
}

func (p *Promise) Do() {
	if p.resolve == nil {
		p.resolve = func(interface{}) {}
	}
	if p.reject == nil {
		p.reject = func(error) {}
	}
	go p.async(p.resolve, p.reject)
}

func NewPromise(fn PromiseFunc) *Promise {
	return &Promise{async: fn}
}

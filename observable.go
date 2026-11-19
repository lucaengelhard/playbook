package main

type Observable[T any] struct {
	Data        T
	Subscribers []*func(T)
}

func NewObservable[T any](initialValue T) *Observable[T] {
	return &Observable[T]{
		Data:        initialValue,
		Subscribers: []*func(T){},
	}
}

func (observable *Observable[T]) Update(updateFn func(T) T) {
	observable.Data = updateFn(observable.Data)
	for _, sub := range observable.Subscribers {
		fn := *sub
		fn(observable.Data)
	}
}

func (observable *Observable[T]) UpdateValue(value T) {
	observable.Update(func(_ T) T {
		return value
	})
}

func (observable *Observable[T]) OnUpdate(subscriber func(T)) func() {
	addr := &subscriber
	observable.Subscribers = append(observable.Subscribers, addr)

	return func() {
		newSubList := []*func(T){}
		for _, sub := range observable.Subscribers {
			if sub != addr {
				newSubList = append(newSubList, sub)
			}
		}
		observable.Subscribers = newSubList
	}
}

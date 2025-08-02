package algoritmos

import (
	"errors"
)

type Nulleable[T any] interface {
	Null() T
	Equal(T) bool
}

// Queue Cola de procesos o hilos

func (c *Cola[T]) Contains(t T) bool {
	for _, e := range c.elements {
		if e.Equal(t) {
			return true
		}
	}
	// Si llega acá es porque la queue está vacia => No contiene el elemento
	return false
}

func (c *Cola[T]) GetElements() []T {
	return c.elements
}

func (c *Cola[T]) Add(t T) {
	c.mutex.Lock()
	c.elements = append(c.elements, t)
	c.mutex.Unlock()
}

func (c *Cola[T]) GetAndRemoveNext() (T, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if len(c.elements) == 0 {
		return T.Null(*new(T)), errors.New("no hay elementos para quitar de la cola")
	}
	nextThread := c.elements[0]
	c.elements = c.elements[1:]

	return nextThread, nil
}

func (c *Cola[T]) IsEmpty() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.elements) == 0
}

func (c *Cola[T]) Size() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.elements)
}

func (c *Cola[T]) Do(f func(T)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for i := range c.elements {
		f(c.elements[i])
	}
}

func (c *Cola[T]) Remove(t T) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for i := range c.elements {
		if c.elements[i].Equal(t) { // Comparación de punteros directamente
			c.elements = append(c.elements[:i], c.elements[i+1:]...)
			return nil
		}
	}
	return errors.New("elemento no encontrado en la cola")
}

func (c *Cola[T]) First() T {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.elements) == 0 {
		var zero T
		return zero // o T.Null(*new(T)) si usás Null() como en otros métodos
	}
	return c.elements[0]
}

func (q *Cola[T]) Values() []T {
	return q.elements
}

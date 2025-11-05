package impl

func (c *ChanW[E]) Unwrap() chan<- E {
	if c == nil {
		return nil
	}
	return c.c
}

func (c *ChanW[E]) IsNil() bool {
	return c == nil || c.c == nil
}

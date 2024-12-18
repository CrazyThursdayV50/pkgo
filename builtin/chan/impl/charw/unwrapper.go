package impl

func (c *ChanW[E]) Unwrap() chan<- E {
	if c == nil {
		return nil
	}
	return c.c
}

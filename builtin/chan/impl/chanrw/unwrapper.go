package impl

func (c *ChanRW[E]) Unwrap() chan E {
	if c == nil {
		return nil
	}
	return c.c
}

package impl

func (c *ChanRW[E]) Unwrap() chan E {
	if c == nil {
		return nil
	}
	return c.c
}

func (c *ChanRW[E]) IsNil() bool {
	return c == nil || c.c == nil
}

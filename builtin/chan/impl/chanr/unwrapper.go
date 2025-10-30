package impl

func (c *ChanR[E]) Unwrap() <-chan E {
	if c == nil {
		return nil
	}
	return c.c
}

func (c *ChanR[E]) IsNil() bool {
	return c == nil || c.c == nil
}

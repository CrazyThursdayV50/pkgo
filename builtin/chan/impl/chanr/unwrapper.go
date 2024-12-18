package impl

func (c *ChanR[E]) Unwrap() <-chan E {
	if c == nil {
		return nil
	}
	return c.c
}

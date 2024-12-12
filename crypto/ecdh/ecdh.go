package ecdh

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
)

type Curve struct {
	curve ecdh.Curve
}

func NewCurve(curve ecdh.Curve) *Curve {
	var c Curve
	if curve == nil {
		curve = ecdh.P256()
	}

	c.curve = curve
	return &c
}

func (c *Curve) RandomEcdh() (*ecdh.PrivateKey, error) {
	return c.curve.GenerateKey(rand.Reader)
}

func (c *Curve) GenPubKey(prvkey []byte) ([]byte, error) {
	var buf = bytes.NewReader(prvkey)
	private, err := c.curve.GenerateKey(buf)
	if err != nil {
		return nil, err
	}

	return private.PublicKey().Bytes(), nil
}

func (c *Curve) ReadPrivate(key []byte) (*ecdh.PrivateKey, error) {
	return c.curve.NewPrivateKey(key)
}

func (c *Curve) ReadPublic(key []byte) (*ecdh.PublicKey, error) {
	return c.curve.NewPublicKey(key)
}

func (c *Curve) GenSharedPoint(prvA, prvB []byte) ([]byte, error) {
	privateA, err := c.curve.NewPrivateKey(prvA)
	if err != nil {
		return nil, err
	}

	privateB, err := c.curve.NewPrivateKey(prvB)
	if err != nil {
		return nil, err
	}

	return privateA.ECDH(privateB.PublicKey())
}

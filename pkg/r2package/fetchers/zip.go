package fetchers

import (
	"context"
	"errors"
)

type Zip struct {
	url string
}

func NewZip(url string) *Zip {
	return &Zip{url: url}
}

func (z *Zip) Fetch(ctx context.Context, dir string) error {
	return errors.New("not implemented")
}

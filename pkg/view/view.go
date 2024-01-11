package view

import (
	"github.com/edrlab/pubstore/pkg/conf"
	"github.com/edrlab/pubstore/pkg/stor"
)

type View struct {
	*conf.Config
	*stor.Store
}

func Init(c *conf.Config, s *stor.Store) View {
	return View{
		Config: c,
		Store:  s,
	}
}

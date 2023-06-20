package view

import "github.com/edrlab/pubstore/pkg/stor"

type View struct {
	stor *stor.Stor
}

func Init(s *stor.Stor) *View {
	return &View{stor: s}
}

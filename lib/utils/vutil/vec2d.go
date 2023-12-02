package vutil

import "fmt"

type Vec2d[T any] struct {
	NRows int
	NCols int
	Data  []T
}

func NewVec2d[T any](nRows, nCols int, data []T) (v Vec2d[T], err error) {
	if nRows*nCols == len(data) {
		v = Vec2d[T]{NRows: nRows, NCols: nCols, Data: data}
	} else {
		err = fmt.Errorf("incorrect vector dimensions: nRows=%v, nCols=%v, len(data)=%v", nRows, nCols, len(data))
	}
	return
}

func (v *Vec2d[T]) Get(iRow, iCol int) *T {
	return &v.Data[iRow*v.NCols+iCol]
}

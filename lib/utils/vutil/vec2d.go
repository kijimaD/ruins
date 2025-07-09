package vutil

import (
	"errors"
	"fmt"
)

// ジェネリック型Tの2次元ベクトル/行列
type Vec2d[T any] struct {
	NRows int
	NCols int
	Data  []T
}

// 指定された次元とデータで新しい2次元ベクトルを作成する
func NewVec2d[T any](nRows, nCols int, data []T) (v Vec2d[T], err error) {
	if nRows <= 0 || nCols <= 0 {
		return v, errors.New("dimensions must be positive")
	}
	if nRows*nCols != len(data) {
		return v, fmt.Errorf("incorrect vector dimensions: nRows=%v, nCols=%v, len(data)=%v", nRows, nCols, len(data))
	}
	return Vec2d[T]{NRows: nRows, NCols: nCols, Data: data}, nil
}

// 指定位置の要素へのポインタを返す
// 範囲外の場合はnilを返す
func (v *Vec2d[T]) Get(iRow, iCol int) *T {
	if !v.IsValidIndex(iRow, iCol) {
		return nil
	}
	return &v.Data[iRow*v.NCols+iCol]
}

// 境界チェック付きで指定位置の要素を返す
func (v *Vec2d[T]) GetSafe(iRow, iCol int) (T, error) {
	var zero T
	if !v.IsValidIndex(iRow, iCol) {
		return zero, fmt.Errorf("index out of bounds: row=%d, col=%d, bounds=[%d,%d]", iRow, iCol, v.NRows, v.NCols)
	}
	return v.Data[iRow*v.NCols+iCol], nil
}

// 指定位置の要素を設定する
func (v *Vec2d[T]) Set(iRow, iCol int, value T) error {
	if !v.IsValidIndex(iRow, iCol) {
		return fmt.Errorf("index out of bounds: row=%d, col=%d, bounds=[%d,%d]", iRow, iCol, v.NRows, v.NCols)
	}
	v.Data[iRow*v.NCols+iCol] = value
	return nil
}

// 指定されたインデックスが境界内かチェックする
func (v *Vec2d[T]) IsValidIndex(iRow, iCol int) bool {
	return iRow >= 0 && iRow < v.NRows && iCol >= 0 && iCol < v.NCols
}

// 総要素数を返す
func (v *Vec2d[T]) Size() int {
	return v.NRows * v.NCols
}

// ベクトルの次元を返す
func (v *Vec2d[T]) Dimensions() (rows, cols int) {
	return v.NRows, v.NCols
}

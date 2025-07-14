package raw

import "github.com/pkg/errors"

// ErrInvalidEnumType はenumに無効な値が指定された場合のエラー
var ErrInvalidEnumType = errors.New("enumに無効な値が指定された")

// ================
// 値タイプ

type ValueType string

const (
	PercentageType ValueType = "PERCENTAGE"
	NumeralType    ValueType = "NUMERAL"
)

func (enum ValueType) Valid() error {
	switch enum {
	case PercentageType, NumeralType:
		return nil
	}
	return errors.Wrapf(ErrInvalidEnumType, "get %s", enum)
}

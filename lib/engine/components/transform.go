package components

import (
	"fmt"

	"github.com/kijimaD/ruins/lib/engine/math"
	"github.com/kijimaD/ruins/lib/engine/utils"
)

// Transform origin variants
const (
	TransformOriginTopLeft      = "TopLeft"
	TransformOriginTopMiddle    = "TopMiddle"
	TransformOriginTopRight     = "TopRight"
	TransformOriginMiddleLeft   = "MiddleLeft"
	TransformOriginMiddle       = "Middle"
	TransformOriginMiddleRight  = "MiddleRight"
	TransformOriginBottomLeft   = "BottomLeft"
	TransformOriginBottomMiddle = "BottomMiddle"
	TransformOriginBottomRight  = "BottomRight"
)

// Transform component.
// The origin (0, 0) is the lower left part of screen.
// Image is first rotated, then scaled, and finally translated.
type Transform struct {
	// Scale1 vector defines image scaling. Contains scale value minus 1 so that zero value is identity.
	Scale1 math.Vector2 `toml:"scale_minus_1"`
	// Rotation angle is measured counterclockwise.
	Rotation float64
	// Translation defines the position of the image center relative to the origin.
	Translation math.Vector2
	// Origin defines the origin (0, 0) relative to the screen. Default is "BottomLeft".
	Origin string
	// Depth determines the drawing order on the screen. Images with higher depth are drawn above others.
	Depth float64
}

// NewTransform creates a new default transform, corresponding to identity.
func NewTransform() *Transform {
	return &Transform{}
}

// SetScale sets transform scale.
func (t *Transform) SetScale(sx, sy float64) *Transform {
	t.Scale1.X = sx - 1
	t.Scale1.Y = sy - 1
	return t
}

// SetRotation sets transform rotation.
func (t *Transform) SetRotation(angle float64) *Transform {
	t.Rotation = angle
	return t
}

// SetTranslation sets transform translation.
func (t *Transform) SetTranslation(tx, ty float64) *Transform {
	t.Translation.X = tx
	t.Translation.Y = ty
	return t
}

// SetDepth sets transform depth.
func (t *Transform) SetDepth(depth float64) *Transform {
	t.Depth = depth
	return t
}

// SetOrigin sets transform origin.
func (t *Transform) SetOrigin(origin string) *Transform {
	t.Origin = origin
	return t
}

// ComputeOriginOffset returns the transform origin offset.
func (t *Transform) ComputeOriginOffset(screenWidth, screenHeight float64) (offsetX, offsetY float64) {
	switch t.Origin {
	case TransformOriginTopLeft:
		offsetX, offsetY = 0, screenHeight
	case TransformOriginTopMiddle:
		offsetX, offsetY = screenWidth/2, screenHeight
	case TransformOriginTopRight:
		offsetX, offsetY = screenWidth, screenHeight
	case TransformOriginMiddleLeft:
		offsetX, offsetY = 0, screenHeight/2
	case TransformOriginMiddle:
		offsetX, offsetY = screenWidth/2, screenHeight/2
	case TransformOriginMiddleRight:
		offsetX, offsetY = screenWidth, screenHeight/2
	case TransformOriginBottomLeft:
		offsetX, offsetY = 0, 0
	case TransformOriginBottomMiddle:
		offsetX, offsetY = screenWidth/2, 0
	case TransformOriginBottomRight:
		offsetX, offsetY = screenWidth, 0
	case "": // TransformOriginBottomLeft
		offsetX, offsetY = 0, 0
	default:
		utils.LogError(fmt.Errorf("unknown transform origin value: %s", t.Origin))
	}
	return
}

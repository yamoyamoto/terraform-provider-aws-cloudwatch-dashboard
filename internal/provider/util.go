package provider

type widgetPosition struct {
	X int32
	Y int32
}

type widgetSize struct {
	Width  int32
	Height int32
}

const (
	MAX_WIDTH = 24
)

func calculatePosition(size widgetSize, beforeWidgetPosition *widgetPosition) widgetPosition {
	if beforeWidgetPosition == nil {
		return widgetPosition{
			X: 0,
			Y: 0,
		}
	}

	x := beforeWidgetPosition.X + size.Width
	y := beforeWidgetPosition.Y
	if x >= MAX_WIDTH {
		x = 0
		y += size.Height
	}

	return widgetPosition{
		X: x,
		Y: y,
	}
}

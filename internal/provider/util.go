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

	if beforeWidgetPosition.X+size.Width > MAX_WIDTH {
		return widgetPosition{
			X: 0,
			Y: beforeWidgetPosition.Y + size.Height,
		}
	}

	return widgetPosition{
		X: beforeWidgetPosition.X,
		Y: beforeWidgetPosition.Y,
	}
}

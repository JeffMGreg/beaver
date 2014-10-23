package beaver

var (
    // ANSI Color escap sequences
	cyan     paint = `1;36` // Debug
	green    paint = `1;32` // Info

	yellow   paint = `1;33` // Warning
	red      paint = `1;31` // Error
    purple   paint = `1;35` // Critical


    colors = []paint{cyan, green, yellow, red, purple}
)

type brush func(string) string
type paint string

type style struct {
    color paint
    code  string
}

func (s style) brush() brush {
    return func(text string) string {
        return s.code + text + "\033[0m"
    }
}

func newBrush(color paint) brush {
    return style{color, "\033[" + string(color) + "m" + ``}.brush()
}

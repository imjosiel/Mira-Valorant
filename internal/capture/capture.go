package capture

import (
	"image"
	"syscall"
	"unsafe"

	"github.com/kbinani/screenshot"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procGetCursorPos = user32.NewProc("GetCursorPos")
)

type Point struct {
	X, Y int32
}

func GetCursorPosition() (int, int, error) {
	var pt Point
	ret, _, err := procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	if ret == 0 {
		return 0, 0, err
	}
	return int(pt.X), int(pt.Y), nil
}

// CaptureRegion captures a rectangle of the screen
func CaptureRegion(x, y, width, height int) (*image.RGBA, error) {
	// screenshot.CaptureRect expects image.Rectangle
	rect := image.Rect(x, y, x+width, y+height)
	return screenshot.CaptureRect(rect)
}

// Helper to calculate the source rectangle based on cursor and zoom
func CalculateSourceRect(cursorX, cursorY int, scopeSize float64, zoomLevel float64) (int, int, int, int) {
	sourceSize := scopeSize / zoomLevel
	halfSize := int(sourceSize / 2)
	
	x := cursorX - halfSize
	y := cursorY - halfSize
	w := int(sourceSize)
	h := int(sourceSize)
	
	return x, y, w, h
}

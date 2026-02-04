package ui

import (
	"fmt"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"mira-valorant/internal/capture"
	"mira-valorant/internal/config"

	"github.com/lxn/win"
)

var (
	className  = syscall.StringToUTF16Ptr("MiraScopeClass")
	windowName = syscall.StringToUTF16Ptr("Mira Scope")

	// Dynamic Imports
	user32 = syscall.NewLazyDLL("user32.dll")
	gdi32  = syscall.NewLazyDLL("gdi32.dll")

	procSetLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")
	procSetWindowDisplayAffinity   = user32.NewProc("SetWindowDisplayAffinity")
	procGetAsyncKeyState           = user32.NewProc("GetAsyncKeyState")
	procCreatePen                  = gdi32.NewProc("CreatePen")
	procRectangle                  = gdi32.NewProc("Rectangle")
)

// Missing constants in lxn/win
const (
	LWA_ALPHA              = 0x00000002
	WDA_EXCLUDEFROMCAPTURE = 0x00000011
)

func GetAsyncKeyState(vKey int) uint16 {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(vKey))
	return uint16(ret)
}

func SetWindowDisplayAffinity(hwnd win.HWND, dwAffinity uint32) bool {
	ret, _, _ := procSetWindowDisplayAffinity.Call(
		uintptr(hwnd),
		uintptr(dwAffinity),
	)
	return ret != 0
}

func SetLayeredWindowAttributes(hwnd win.HWND, crKey win.COLORREF, bAlpha byte, dwFlags uint32) bool {
	ret, _, _ := procSetLayeredWindowAttributes.Call(
		uintptr(hwnd),
		uintptr(crKey),
		uintptr(bAlpha),
		uintptr(dwFlags),
	)
	return ret != 0
}

func CreatePen(iStyle int, cWidth int, color win.COLORREF) win.HPEN {
	ret, _, _ := procCreatePen.Call(
		uintptr(iStyle),
		uintptr(cWidth),
		uintptr(color),
	)
	return win.HPEN(ret)
}

// Renamed to WinRectangle to avoid collision with lxn/walk/declarative
func WinRectangle(hdc win.HDC, left, top, right, bottom int32) bool {
	ret, _, _ := procRectangle.Call(
		uintptr(hdc),
		uintptr(left),
		uintptr(top),
		uintptr(right),
		uintptr(bottom),
	)
	return ret != 0
}

func RunScopeWindow(state *config.AppState) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	hInst := win.GetModuleHandle(nil)

	var wc win.WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.LpfnWndProc = syscall.NewCallback(wndProc)
	wc.HInstance = hInst
	wc.LpszClassName = className
	wc.HbrBackground = win.HBRUSH(win.GetStockObject(win.BLACK_BRUSH))

	if atom := win.RegisterClassEx(&wc); atom == 0 {
		fmt.Println("RegisterClassEx failed")
		return
	}

	exStyle := uint32(win.WS_EX_LAYERED | win.WS_EX_TRANSPARENT | win.WS_EX_TOPMOST | win.WS_EX_TOOLWINDOW)
	style := uint32(win.WS_POPUP)

	hwnd := win.CreateWindowEx(
		exStyle,
		className,
		windowName,
		style,
		0, 0, 250, 250,
		0, 0, hInst, nil,
	)

	if hwnd == 0 {
		fmt.Println("CreateWindowEx failed")
		return
	}

	SetLayeredWindowAttributes(hwnd, 0, 255, LWA_ALPHA)
	// Prevent the scope window from being captured by BitBlt/screen capture
	SetWindowDisplayAffinity(hwnd, WDA_EXCLUDEFROMCAPTURE)

	// High Performance Timer (~60 FPS)
	// 33ms = ~30 FPS
	// 16ms = ~60 FPS
	// 8ms  = ~120 FPS
	// Using 10ms target for smoother experience (~100 FPS cap)
	ticker := time.NewTicker(time.Millisecond * 10)
	defer ticker.Stop()

	var msg win.MSG
	lastKeyState := false

	for {
		if win.PeekMessage(&msg, 0, 0, 0, win.PM_REMOVE) {
			if msg.Message == win.WM_QUIT {
				break
			}
			win.TranslateMessage(&msg)
			win.DispatchMessage(&msg)
		}

		// Hotkey Logic
		hotkey := state.Hotkey()
		mode := state.HotkeyMode()
		active := state.IsActive()

		if hotkey != 0 {
			keyState := (GetAsyncKeyState(hotkey) & 0x8000) != 0

			if mode == "Hold" {
				if keyState != active {
					state.SetIsActive(keyState)
					active = keyState
				}
			} else { // Toggle
				if keyState && !lastKeyState {
					newState := !active
					state.SetIsActive(newState)
					active = newState
				}
			}
			lastKeyState = keyState
		}

		if !active {
			if win.IsWindowVisible(hwnd) {
				win.ShowWindow(hwnd, win.SW_HIDE)
			}
			time.Sleep(10 * time.Millisecond) // Faster poll for hotkey
			continue
		} else {
			if !win.IsWindowVisible(hwnd) {
				win.ShowWindow(hwnd, win.SW_SHOW)
			}
		}

		sizeVal := state.ScopeSize()
		size := int32(sizeVal)

		follow := state.FollowCursor()

		var targetX, targetY int
		if follow {
			tx, ty, _ := capture.GetCursorPosition()
			targetX, targetY = int(tx), int(ty)
		} else {
			// Center of screen
			targetX = int(win.GetSystemMetrics(win.SM_CXSCREEN) / 2)
			targetY = int(win.GetSystemMetrics(win.SM_CYSCREEN) / 2)
		}

		// Window Position (Top-Left)
		winX := int32(targetX) - size/2
		winY := int32(targetY) - size/2

		win.SetWindowPos(hwnd, win.HWND_TOPMOST, winX, winY, size, size, win.SWP_NOACTIVATE|win.SWP_NOREDRAW)

		zoom := state.ZoomLevel()

		// Render
		hdcWindow := win.GetDC(hwnd)

		// Create Mem DC for double buffering
		hdcMem := win.CreateCompatibleDC(hdcWindow)
		hBitmap := win.CreateCompatibleBitmap(hdcWindow, size, size)
		hOld := win.SelectObject(hdcMem, win.HGDIOBJ(hBitmap))

		// 1. Capture Screen (StretchBlt from Screen DC to Mem DC)
		// We capture around the target point (Cursor or Screen Center)
		srcSize := float64(size) / zoom
		srcX := targetX - int(srcSize/2)
		srcY := targetY - int(srcSize/2)

		hdcScreen := win.GetDC(0)

		// Set stretch mode to COLORONCOLOR (simpler/faster) or HALFTONE (better quality)
		// HALFTONE is slower but looks better. COLORONCOLOR is fastest.
		// Trying HALFTONE first for quality, if slow switch to COLORONCOLOR (3)
		win.SetStretchBltMode(hdcMem, win.HALFTONE)
		win.SetBrushOrgEx(hdcMem, 0, 0, nil)

		win.StretchBlt(
			hdcMem,
			0, 0, size, size,
			hdcScreen,
			int32(srcX), int32(srcY), int32(srcSize), int32(srcSize),
			win.SRCCOPY,
		)
		win.ReleaseDC(0, hdcScreen)

		// 2. Draw Border
		borderEnabled := state.BorderEnabled()
		if borderEnabled {
			thick := state.BorderThickness()
			colorIdx := state.BorderColor()
			var c win.COLORREF
			switch colorIdx {
			case 0:
				c = win.RGB(255, 0, 0)
			case 1:
				c = win.RGB(0, 255, 0)
			case 2:
				c = win.RGB(0, 0, 255)
			case 3:
				c = win.RGB(255, 255, 0)
			default:
				c = win.RGB(255, 0, 0)
			}

			pen := CreatePen(win.PS_SOLID, int(thick), c)
			oldPen := win.SelectObject(hdcMem, win.HGDIOBJ(pen))
			oldBrush := win.SelectObject(hdcMem, win.HGDIOBJ(win.GetStockObject(win.NULL_BRUSH)))

			WinRectangle(hdcMem, 0, 0, size, size)

			win.SelectObject(hdcMem, oldBrush)
			win.SelectObject(hdcMem, oldPen)
			win.DeleteObject(win.HGDIOBJ(pen))
		}

		// 3. Blt to Window
		win.BitBlt(hdcWindow, 0, 0, size, size, hdcMem, 0, 0, win.SRCCOPY)

		// Cleanup
		win.SelectObject(hdcMem, hOld)
		win.DeleteObject(win.HGDIOBJ(hBitmap))
		win.DeleteDC(hdcMem)
		win.ReleaseDC(hwnd, hdcWindow)

		<-ticker.C
	}
}

func wndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_DESTROY:
		win.PostQuitMessage(0)
		return 0
	case win.WM_NCHITTEST:
		// Fix constant overflow properly for 32/64 bit safety
		// win.HTTRANSPARENT is -1
		val := int32(win.HTTRANSPARENT)
		return uintptr(val)
	}
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

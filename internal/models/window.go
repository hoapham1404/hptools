package models

// WindowInfo represents window position and size information
type WindowInfo struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// RECT structure for window coordinates
type RECT struct {
	Left, Top, Right, Bottom int32
}

package Utilities

import (
	"github.com/fatih/color"
)

// ErrorCheckedColourPrint Helper function to simplify error handling
func ErrorCheckedColourPrint(colour *color.Color, message string) {
	_, err := colour.Println(message)
	if err != nil {
		return
	}
}

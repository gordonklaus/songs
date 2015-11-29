package main

import (
	"github.com/gordonklaus/audio"
	"github.com/gordonklaus/audiogui"
)

var reverb_pattern = audiogui.NewPattern([]*audio.Note{}, map[string][]*audio.ControlPoint{
	"Sustain": {
		{0, -16},
	},
	"Dry": {
		{0, -12},
	},
	"Wet": {
		{0, 0},
		{173, 0},
		{184, -16},
	},
	"Decay": {
		{0, 8},
	},
})

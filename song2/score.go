package main

import "github.com/gordonklaus/audio"

var score = &audio.Score{[]*audio.Part{
	{"Sines", []*audio.PatternEvent{
		{0, sines_pattern},
	}},
}}

package main

import "github.com/gordonklaus/audio"

var score = &audio.Score{[]*audio.Part{
	{"Voices", []*audio.PatternEvent{
		{0, part_pattern},
	}},
}}

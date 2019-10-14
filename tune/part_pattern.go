package main

import (
	"github.com/gordonklaus/audio"
	"github.com/gordonklaus/audiogui"
)

var part_pattern = audiogui.NewPattern([]*audio.Note{
	{0, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 8},
			{2, 8},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{1, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 8.415037499278844},
			{2, 8.415037499278844},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{2, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 8.093109404391482},
			{2, 8.093109404391482},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{3, map[string][]*audio.ControlPoint{
		"Amplitude": {
			{0, 0},
		},
		"Pitch": {
			{0, 8.678071905112638},
			{2, 8.678071905112638},
		},
	}},
	{4, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 7.830074998557689},
			{4, 7.830074998557689},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{5, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 8.567040592723895},
			{3, 8.567040592723895},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{6, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 7.263034405833794},
			{2, 7.263034405833794},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{7, map[string][]*audio.ControlPoint{
		"Amplitude": {
			{0, 0},
		},
		"Pitch": {
			{0, 9.070389327891398},
			{1, 9.070389327891398},
		},
	}},
	{8, map[string][]*audio.ControlPoint{
		"Amplitude": {
			{0, 0},
		},
		"Pitch": {
			{0, 7.678071905112638},
			{3, 7.678071905112638},
		},
	}},
	{8, map[string][]*audio.ControlPoint{
		"Amplitude": {
			{0, 0},
		},
		"Pitch": {
			{0, 9},
			{2, 9},
		},
	}},
	{9, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 8.415037499278844},
			{2, 8.415037499278844},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{10, map[string][]*audio.ControlPoint{
		"Amplitude": {
			{0, 0},
		},
		"Pitch": {
			{0, 9.093109404391482},
			{2, 9.093109404391482},
			{5, 9.070389327891398},
		},
	}},
	{11, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 8.678071905112638},
			{1, 8.678071905112638},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{12, map[string][]*audio.ControlPoint{
		"Amplitude": {
			{0, 0},
		},
		"Pitch": {
			{0, 8.584962500721156},
			{2, 8.584962500721156},
		},
	}},
	{13, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 8.84799690655495},
			{3, 8.84799690655495},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{14, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 8.263034405833794},
			{2, 8.263034405833794},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{15, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 8.070389327891398},
			{1, 8.070389327891398},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{16, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 8},
			{2, 8},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{16, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 8.678071905112638},
			{2, 8.678071905112638},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{23, map[string][]*audio.ControlPoint{
		"Amplitude": {
			{0, 0},
		},
		"Pitch": {
			{0, 9},
		},
	}},
	{23, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 8.321928094887362},
			{3, 8.321928094887362},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{23, map[string][]*audio.ControlPoint{
		"Amplitude": {
			{0, 0},
		},
		"Pitch": {
			{0, 9},
			{1.5, 9},
		},
	}},
	{23, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 7},
			{3, 7},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{23, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 7.584962500721156},
			{3, 7.584962500721156},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{24.5, map[string][]*audio.ControlPoint{
		"Amplitude": {
			{0, 0},
		},
		"Pitch": {
			{0, 8.906890595608518},
			{1.5, 8.906890595608518},
		},
	}},
	{26, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 7.07681559705083},
			{1.5, 7.07681559705083},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{26, map[string][]*audio.ControlPoint{
		"Amplitude": {
			{0, 0},
		},
		"Pitch": {
			{0, 8.584962500721156},
			{1.5, 8.584962500721156},
		},
	}},
	{26, map[string][]*audio.ControlPoint{
		"Amplitude": {
			{0, 0},
		},
		"Pitch": {
			{0, 7.754887502163468},
			{3, 7.754887502163468},
		},
	}},
	{27.5, map[string][]*audio.ControlPoint{
		"Amplitude": {
			{0, 0},
		},
		"Pitch": {
			{0, 8.491853096329674},
			{1.5, 8.491853096329674},
		},
	}},
	{27.5, map[string][]*audio.ControlPoint{
		"Amplitude": {
			{0, 0},
		},
		"Pitch": {
			{0, 7.169925001442312},
			{1.5, 7.169925001442312},
		},
	}},
	{29, map[string][]*audio.ControlPoint{
		"Amplitude": {
			{0, 0},
		},
		"Pitch": {
			{0, 8.321928094887362},
			{1.5, 8.321928094887362},
		},
	}},
	{29, map[string][]*audio.ControlPoint{
		"Amplitude": {
			{0, 0},
		},
		"Pitch": {
			{0, 7.584962500721156},
			{3, 7.584962500721156},
		},
	}},
	{29, map[string][]*audio.ControlPoint{
		"Amplitude": {
			{0, 0},
		},
		"Pitch": {
			{0, 7},
			{3, 7},
		},
	}},
	{30.5, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 8.169925001442312},
			{0.75, 8.169925001442312},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{31.25, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 8},
			{0.75, 8},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{32, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 8.169925001442312},
			{3, 8.169925001442312},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{32, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 7.754887502163468},
			{1.5, 7.754887502163468},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{32, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 7.491853096329675},
			{1.5, 7.491853096329675},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{33.5, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 7.906890595608519},
			{1.5, 7.906890595608519},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
	{33.5, map[string][]*audio.ControlPoint{
		"Pitch": {
			{0, 7.584962500721156},
			{1.5, 7.584962500721156},
		},
		"Amplitude": {
			{0, 0},
		},
	}},
}, map[string][]*audio.ControlPoint{})

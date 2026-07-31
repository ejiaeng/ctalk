package main

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
)

func initTextInput() textinput.Model {
	ti := textinput.New()
	ti.Prompt = "Send a message:"
	ti.Placeholder = "Write a message ..."
	return ti
}

func initViewPort() viewport.Model {
	vp := viewport.New(viewport.WithHeight(20), viewport.WithWidth(80))
	vp.SoftWrap = true
	return vp
}

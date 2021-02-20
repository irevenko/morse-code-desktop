package main

import (
	m "./morse"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hajimehoshi/oto"
	"log"
)

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	m.SetWindowOptions(win)

	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	m.SetGridOptions(grid)

	toLabel, _ := gtk.LabelNew("To Morse")
	fromLabel, _ := gtk.LabelNew("From Morse")
	toLabelResult, _ := gtk.LabelNew("To Morse Translation:")
	fromLabelResult, _ := gtk.LabelNew("From Morse Translation:")

	grid.Attach(toLabel, 0,0,1,1)
	grid.Attach(fromLabel, 0,1,1,1)
	grid.Attach(toLabelResult, 0,3,1,1)
	grid.Attach(fromLabelResult, 0,4,1,1)

	toEntry, _ := gtk.EntryNew()
	fromEntry, _ := gtk.EntryNew()

	grid.Attach(toEntry, 1,0,2,1)
	grid.Attach(fromEntry, 1,1,2,1)

	translateToBox := m.SetupBox(gtk.ORIENTATION_HORIZONTAL)
	translateToTextView := m.SetupTview()
	translateToBox.PackStart(translateToTextView, true, true, 0)

	translateFromBox := m.SetupBox(gtk.ORIENTATION_HORIZONTAL)
	translateFromTextView := m.SetupTview()
	translateFromBox.PackStart(translateFromTextView, true, true, 0)

	grid.Attach(translateToBox, 1,3,3,1)
	grid.Attach(translateFromBox, 1,4,3,1)

	toButton := m.SetupButton("Translate to", func() {
		text := m.GetTextFromEntry(toEntry)
		translation := m.TranslateToMorse(text)

		buffer, err := translateToTextView.GetBuffer()
		if err != nil {
			log.Fatal("Unable to get the buffer", err)
		}
		buffer.SetText(translation)
	})

	fromButton := m.SetupButton("Translate from", func() {
		text := m.GetTextFromEntry(fromEntry)
		translation := m.TranslateFromMorse(text)

		buffer, err := translateFromTextView.GetBuffer()
		if err != nil {
			log.Fatal("Unable to get the buffer", err)
		}
		buffer.SetText(translation)
	})


	c, err := oto.NewContext(20000, 2, 2, 4096)
	if err != nil {
		log.Fatal(err)
	}
	p := c.NewPlayer()

	playButton := m.SetupButton("Play Morse Code", func() {
		morseSequence := m.GetTextFromTview(translateToTextView)
		m.MorseToSound(morseSequence, c, p)
	})

	grid.Attach(toButton, 3, 0, 1, 1)
	grid.Attach(fromButton, 3, 1, 1, 1)
	grid.Attach(playButton, 1,5,1,1)

	win.Add(grid)
	win.Connect("destroy", gtk.MainQuit)
	win.ShowAll()
	gtk.Main()
}

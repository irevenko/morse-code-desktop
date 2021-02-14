package morse

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
)

func GetBufferFromEntry(entry *gtk.Entry) *gtk.EntryBuffer {
	buffer, err := entry.GetBuffer()
	if err != nil {
		log.Fatal("Unable to get the buffer", err)
	}
	return buffer
}

func GetTextFromEntry(entry *gtk.Entry) string {
	buffer := GetBufferFromEntry(entry)
	text, err := buffer.GetText()
	if err != nil {
		log.Fatal("Unable to get the text", err)
	}
	return text
}

func SetupButton(label string, onClick func()) *gtk.Button {
	btn, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	btn.Connect("clicked", onClick)
	return btn
}

func SetWindowOptions(w *gtk.Window) {
	w.SetPosition(gtk.WIN_POS_CENTER)
	w.SetDefaultSize(800, 300)
	w.SetTitle("Morse Code Translator")
	w.SetResizable(false)
	w.SetBorderWidth(10)
}

func SetGridOptions(g *gtk.Grid) {
	g.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	g.SetRowSpacing(20)
	g.SetColumnSpacing(10)
}

func SetupBox(orient gtk.Orientation) *gtk.Box {
	box, err := gtk.BoxNew(orient, 200)
	if err != nil {
		log.Fatal("Unable to create box:", err)
	}
	return box
}

func SetupTview() *gtk.TextView {
	tv, err := gtk.TextViewNew()
	if err != nil {
		log.Fatal("Unable to create TextView:", err)
	}
	return tv
}
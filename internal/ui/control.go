package ui

import (
	"fmt"
	"log"

	"mira-valorant/internal/config"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func RunControlWindow(state *config.AppState) {
	mw := new(walk.MainWindow)

	var zoomLabel, sizeLabel, thickLabel *walk.Label
	var activeCheckBox, followCheckBox, borderCheckBox *walk.CheckBox
	var zoomSlider, sizeSlider, thickSlider *walk.Slider
	var colorComboBox, keyComboBox, modeComboBox *walk.ComboBox

	// Key Map for Dropdown
	keyMap := []struct {
		Name string
		Code int
	}{
		{"Nenhum", 0},
		{"Clique Direito (Mouse)", 0x02},  // VK_RBUTTON
		{"Botão Lateral 1 (Mouse)", 0x05}, // VK_XBUTTON1
		{"Botão Lateral 2 (Mouse)", 0x06}, // VK_XBUTTON2
		{"Botão do Meio (Mouse)", 0x04},   // VK_MBUTTON
		{"Shift Esquerdo", 0xA0},          // VK_LSHIFT
		{"Ctrl Esquerdo", 0xA2},           // VK_LCONTROL
		{"Alt Esquerdo", 0xA4},            // VK_LMENU
		{"Caps Lock", 0x14},               // VK_CAPITAL
		{"V", 0x56},
		{"C", 0x43},
		{"F", 0x46},
		{"X", 0x58},
		{"Z", 0x5A},
	}

	if err := (MainWindow{
		AssignTo: &mw,
		Title:    "Mira Controller (Native)",
		MinSize:  Size{Width: 300, Height: 500},
		Size:     Size{Width: 300, Height: 500},
		Layout:   VBox{},
		Children: []Widget{
			GroupBox{
				Title:  "Ativação",
				Layout: VBox{},
				Children: []Widget{
					CheckBox{
						AssignTo: &activeCheckBox,
						Text:     "Ativar Luneta (Manual)",
						OnCheckedChanged: func() {
							state.SetIsActive(activeCheckBox.Checked())
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							Label{Text: "Tecla:"},
							ComboBox{
								AssignTo: &keyComboBox,
								Model: func() []string {
									var names []string
									for _, k := range keyMap {
										names = append(names, k.Name)
									}
									return names
								}(),
								CurrentIndex: 0,
								OnCurrentIndexChanged: func() {
									if keyComboBox != nil {
										idx := keyComboBox.CurrentIndex()
										if idx >= 0 && idx < len(keyMap) {
											state.SetHotkey(keyMap[idx].Code)
										}
									}
								},
							},
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							Label{Text: "Modo:"},
							ComboBox{
								AssignTo:     &modeComboBox,
								Model:        []string{"Alternar (Toggle)", "Segurar (Hold)"},
								CurrentIndex: 0,
								OnCurrentIndexChanged: func() {
									if modeComboBox != nil {
										idx := modeComboBox.CurrentIndex()
										if idx == 0 {
											state.SetHotkeyMode("Toggle")
										} else {
											state.SetHotkeyMode("Hold")
										}
									}
								},
							},
						},
					},
				},
			},
			GroupBox{
				Title:  "Visualização",
				Layout: VBox{},
				Children: []Widget{
					Label{AssignTo: &zoomLabel, Text: "Zoom: 2.0x"},
					Slider{
						AssignTo: &zoomSlider,
						MinValue: 10, // 1.0x * 10
						MaxValue: 80, // 8.0x * 10
						Value:    20,
						OnValueChanged: func() {
							val := zoomSlider.Value()
							f := float64(val) / 10.0
							state.SetZoomLevel(f)
							zoomLabel.SetText(fmt.Sprintf("Zoom: %.1fx", f))
						},
					},
					Label{AssignTo: &sizeLabel, Text: "Tamanho: 250px"},
					Slider{
						AssignTo: &sizeSlider,
						MinValue: 100,
						MaxValue: 600,
						Value:    250,
						OnValueChanged: func() {
							val := sizeSlider.Value()
							state.SetScopeSize(float64(val))
							sizeLabel.SetText(fmt.Sprintf("Tamanho: %dpx", val))
						},
					},
					CheckBox{
						AssignTo: &followCheckBox,
						Text:     "Seguir Cursor",
						Checked:  false,
						OnCheckedChanged: func() {
							state.SetFollowCursor(followCheckBox.Checked())
						},
					},
				},
			},
			GroupBox{
				Title:  "Borda",
				Layout: VBox{},
				Children: []Widget{
					CheckBox{
						AssignTo: &borderCheckBox,
						Text:     "Mostrar Borda",
						Checked:  true,
						OnCheckedChanged: func() {
							state.SetBorderEnabled(borderCheckBox.Checked())
						},
					},
					Label{AssignTo: &thickLabel, Text: "Espessura: 2"},
					Slider{
						AssignTo: &thickSlider,
						MinValue: 1,
						MaxValue: 10,
						Value:    2,
						OnValueChanged: func() {
							val := thickSlider.Value()
							state.SetBorderThickness(float64(val))
							thickLabel.SetText(fmt.Sprintf("Espessura: %d", val))
						},
					},
					ComboBox{
						AssignTo: &colorComboBox,
						Model:    []string{"Vermelho", "Verde", "Azul", "Amarelo"},
						Value:    "Vermelho",
						OnCurrentIndexChanged: func() {
							if colorComboBox != nil {
								state.SetBorderColor(colorComboBox.CurrentIndex())
							}
						},
					},
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	mw.Run()
}

package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"net/http"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Language represents a programming language
type Language struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Quick Reference")
	myWindow.Resize(fyne.NewSize(400, 800))

	// Header
	title := canvas.NewText("Quick Reference", color.Black)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Description text
	description := widget.NewLabel("Here are some cheat sheets and quick references contributed by open source angels.")
	description.Wrapping = fyne.TextWrapWord

	// Fetch languages from API
	languages, err := fetchLanguages()
	if err != nil {
		fmt.Println("Error fetching languages:", err)
		// Use some default languages if API fails
		languages = []Language{
			{ID: 1, Name: "Python", Color: "4F75A1"},
			{ID: 2, Name: "JavaScript", Color: "F1DC5D"},
			{ID: 3, Name: "Swift", Color: "EB735F"},
			{ID: 4, Name: "Rust", Color: "8E6CE1"},
			{ID: 5, Name: "Kotlin", Color: "8E6CE1"},
			{ID: 6, Name: "Dart", Color: "74B8DF"},
			{ID: 7, Name: "Go", Color: "6BAFC6"},
		}
	}

	// Create content container
	content := container.NewVBox(
		container.NewPadded(title),
		container.NewPadded(description),
	)

	// Add language cards
	for _, lang := range languages {
		// Parse hex color
		hexColor := lang.Color
		if len(hexColor) == 6 {
			r, _ := strconv.ParseUint(hexColor[0:2], 16, 8)
			g, _ := strconv.ParseUint(hexColor[2:4], 16, 8)
			b, _ := strconv.ParseUint(hexColor[4:6], 16, 8)
			bgColor := color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
			
			// Create a local copy of the ID for closure
			langID := lang.ID
			langName := lang.Name
			
			langCard := createLanguageCard(langName, bgColor, func() {
				fmt.Printf("Tapped on language %s (ID: %d)\n", langName, langID)
				// Here you could navigate to a detail page or call another API
			})
			
			content.Add(container.NewPadded(langCard))
		}
	}

	// Add some padding around the content
	paddedContent := container.NewPadded(content)
	
	// Set the window content
	myWindow.SetContent(paddedContent)
	myWindow.ShowAndRun()
}

func fetchLanguages() ([]Language, error) {
	resp, err := http.Get("http://localhost:5000/api/languages")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var languages []Language
	err = json.NewDecoder(resp.Body).Decode(&languages)
	if err != nil {
		return nil, err
	}

	return languages, nil
}

func createLanguageCard(language string, bgColor color.Color, onTapped func()) *fyne.Container {
	// Create a tap button that covers the whole card area
	tapButton := widget.NewButton("", onTapped)
	// Make the button transparent
	tapButton.Importance = widget.LowImportance
	
	// Create background rectangle
	background := canvas.NewRectangle(bgColor)
	background.SetMinSize(fyne.NewSize(380, 80))

	// Language name
	langLabel := canvas.NewText(language, color.Black)
	langLabel.TextSize = 18
	langLabel.TextStyle = fyne.TextStyle{Bold: true}
	langLabel.Alignment = fyne.TextAlignCenter

	// Icon (placeholder)
	icon := widget.NewIcon(theme.InfoIcon())

	// Arrow icon
	arrowIcon := widget.NewIcon(theme.NavigateNextIcon())

	// Create layout for content
	contentLayout := container.New(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		container.NewPadded(icon),
		layout.NewSpacer(),
		langLabel,
		layout.NewSpacer(),
		arrowIcon,
		layout.NewSpacer(),
	)

	// Stack the content over the background
	cardContent := container.NewStack(background, contentLayout)
	
	// Create a container with the button overlaid on top of the content
	return container.NewStack(cardContent, tapButton)
}

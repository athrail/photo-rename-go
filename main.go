package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	exif "github.com/dsoprea/go-exif/v3"
	// "github.com/charmbracelet/lipgloss"
)

type model struct {
	path  string
	table table.Model
}

type renameEntry struct {
	filename string
	date     string
	result   string
}

func initialModel(renames []renameEntry) model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Original file", Width: 40},
		{Title: "Date taken", Width: 15},
		{Title: "New name", Width: 40},
	}

	rows := make([]table.Row, len(renames))
	for idx, entry := range renames {
		rows[idx] = table.Row{
			fmt.Sprintf("%v", idx),
			entry.filename,
			entry.date,
			entry.result,
		}
	}

	return model{
		path:  "test path",
		table: table.New(table.WithColumns(columns), table.WithRows(rows), table.WithFocused(false)),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Welcome to PhotoRename tool\n"
	s += "It shall go through the files in a provided path and add suffix to their name based on when was the photo taken.\nIt reads EXIF data that is encoded in the images to do that.\n"

	s += m.table.View() + "\n"

	s += "\nPress q to quit\n"

	return s
}

func grabExifData(image string) (string, error) {
	data, err := os.ReadFile(image)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	exifData, err := exif.SearchAndExtractExif(data)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	tags, med, err := exif.GetFlatExifData(exifData, nil)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	fmt.Println(tags)
	fmt.Println(med)

	return "", nil
}

func getImagesData(imagesPath string) ([]renameEntry, error) {
	if _, err := os.Stat(imagesPath); os.IsNotExist(err) {
		return []renameEntry{}, err
	}

	files, err := os.ReadDir(imagesPath)
	if err != nil {
		return []renameEntry{}, err
	}

	var entries []renameEntry
	for _, file := range files {
		if ext := strings.ToLower(path.Ext(file.Name())); ext == ".jpg" || ext == ".jpeg" {
			entries = append(entries,
				renameEntry{filename: file.Name(), date: "unknown", result: "will calculate"})
		}
	}

	grabExifData(path.Join(imagesPath, entries[0].filename))

	return entries, nil
}

func main() {
	var lookupPath string

	flag.StringVar(&lookupPath, "images", "", "A path where to look for images to parse")
	flag.Parse()

	if lookupPath == "" {
		fmt.Println("Not enough parameters provided")
		os.Exit(1)
	}

	renames, err := getImagesData("/home/athrail/Downloads/exif-samples-master/jpg/")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(renames))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

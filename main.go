package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type todoItem struct {
	Id          int32
	Title       string
	Description string
}

var todoItems = []todoItem{
	{Id: 1, Title: "First item", Description: "My fancy descrtiption"},
	{Id: 2, Title: "My important todo", Description: "It is really important"},
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	todoItemList []todoItem       // items on the to-do list
	cursor       int              // which to-do list item our cursor is pointing at
	selected     map[int]struct{} // which to-do items are selected
	newItemInput textinput.Model
	err          error
}

func initialModel() model {

	ti := textinput.New()
	ti.Placeholder = "New todo item - add"
	ti.Width = 100

	return model{
		// Our to-do list is a grocery list
		todoItemList: todoItems,

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),

		newItemInput: ti,
		err:          nil,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	// return nil

	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.todoItemList)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		case "a":
			fmt.Println("a presssed")
			m.newItemInput.Focus()
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	for i, item := range m.todoItemList {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, item.Title)
	}

	m.newItemInput.View()

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

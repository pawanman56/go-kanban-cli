package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const divisor = 4

const (
	todo status = iota
	inProgress
	done
)

/* CUSTOM ITEM */

type Task struct {
	status      status
	title       string
	description string
}

// implement the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

/* MAIN MODEL */

type Model struct {
	focused status
	lists   []list.Model
	loaded  bool
	err     error
}

func New() *Model {
	return &Model{}
}

// TODO: call this on tea.WindowSizeMsg
func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	// Init To Do
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "Buy milk", description: "Strawberry Milk"},
		Task{status: todo, title: "Eat Sushi", description: "Or momo"},
		Task{status: todo, title: "Do laundry", description: "Dry and Fold them"},
	})

	// Init in progress
	m.lists[inProgress].Title = "In progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: todo, title: "Write code", description: "Just follow the tutorial"},
	})

	// Init done
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: todo, title: "Stay cool", description: "Cool as a cucumber"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loaded {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.lists[todo].View(),
			m.lists[inProgress].View(),
			m.lists[done].View(),
		)

	} else {
		return "loading..."
	}
}

func main() {
	m := New()
	p := tea.NewProgram(m)

	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

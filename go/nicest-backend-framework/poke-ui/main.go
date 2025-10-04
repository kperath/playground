package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"playground/nicest-backend-framework/pokemon"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func main() {
	resp, err := http.Get("http://localhost:8080/pokemon/party/1")
	if err != nil {
		panic(err)
	}

	type PokemonResponse struct {
		Pokemon pokemon.Pokemon `json:"pokemon"`
	}

	respData := PokemonResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", respData.Pokemon)
	// return

	columns := []table.Column{
		{Title: "pokedexnumbe", Width: 10},
		{Title: "islegendary", Width: 10},
		{Title: "megaevolved", Width: 10},
		{Title: "generation", Width: 10},
		{Title: "nickname", Width: 10},
		{Title: "name", Width: 10},
		{Title: "japanesename", Width: 10},
		{Title: "classification", Width: 10},
		{Title: "type1", Width: 10},
		{Title: "type2", Width: 10},
		{Title: "abilities", Width: 10},
		{Title: "experiencegrowth", Width: 10},
		{Title: "experience", Width: 10},
		{Title: "basetotal", Width: 10},
		{Title: "attack", Width: 10},
		{Title: "defense", Width: 10},
		{Title: "speed", Width: 10},
		{Title: "spattack", Width: 10},
		{Title: "spdefense", Width: 10},
		{Title: "hp", Width: 10},
		{Title: "level", Width: 10},
		{Title: "item", Width: 10},
		{Title: "nature", Width: 10},
		{Title: "attacks", Width: 10},
		{Title: "image", Width: 10},
	}

	pkmn := respData.Pokemon
	rows := []table.Row{
		{
			fmt.Sprintf("%d", pkmn.PokedexNumber),
			fmt.Sprintf("%t", pkmn.IsLegendary),
			fmt.Sprintf("%t", pkmn.MegaEvolved),
			pkmn.Generation,
			pkmn.NickName,
			pkmn.Name,
			pkmn.JapaneseName,
			pkmn.Classification,
			pkmn.Type1,
			pkmn.Type2,
			fmt.Sprintf("%v", pkmn.Abilities),
			fmt.Sprintf("%v", pkmn.ExperienceGrowth),
			fmt.Sprintf("%v", pkmn.Experience),
			fmt.Sprintf("%v", pkmn.BaseTotal),
			fmt.Sprintf("%v", pkmn.Attack),
			fmt.Sprintf("%v", pkmn.Defense),
			fmt.Sprintf("%v", pkmn.Speed),
			fmt.Sprintf("%v", pkmn.SpAttack),
			fmt.Sprintf("%v", pkmn.SpDefense),
			fmt.Sprintf("%v", pkmn.HP),
			fmt.Sprintf("%v", pkmn.Level),
			pkmn.Item,
			pkmn.Nature,
			fmt.Sprintf("%v", pkmn.Attacks),
			pkmn.Image,
		},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	// s := table.DefaultStyles()
	// s.Header = s.Header.
	// 	BorderStyle(lipgloss.NormalBorder()).
	// 	BorderForeground(lipgloss.Color("240")).
	// 	BorderBottom(true).
	// 	Bold(false)
	// s.Selected = s.Selected.
	// 	Foreground(lipgloss.Color("229")).
	// 	Background(lipgloss.Color("57")).
	// 	Bold(false)
	// t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

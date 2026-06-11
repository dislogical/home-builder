package bubbles

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"

	tea "charm.land/bubbletea/v2"

	homebuilder "github.com/dislogical/home-builder/pkg"
)

type ContextModel struct {
	HbCtx     *homebuilder.Context
	Resources []homebuilder.Resource

	view    tea.View
	list    list.Model
	spinner spinner.Model
}

func (m *ContextModel) Init() tea.Cmd {
	m.spinner = spinner.New(spinner.WithSpinner(spinner.MiniDot))

	cmds := make([]tea.Cmd, len(m.Resources), len(m.Resources)+1)

	items := make([]list.Item, len(m.Resources))
	for idx, resource := range m.Resources {
		items[idx] = &resourceItem{
			hbctx:    m.HbCtx,
			resource: &resource,
			spinner:  &m.spinner,

			status: homebuilder.StatusUnknown,
			err:    nil,
		}
		cmds[idx] = initResource(&resource)
	}

	m.view = tea.View{
		WindowTitle: "home-builder",
		AltScreen:   true,
	}

	del := list.NewDefaultDelegate()
	del.SetSpacing(0)

	m.list = list.New(items, del, 0, 0)
	m.list.Title = "Home Builder"
	m.list.SetShowStatusBar(false)

	cmds = append(cmds, m.spinner.Tick)

	return tea.Batch(cmds...)
}

func (m *ContextModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case resourceStateMsg:
		for idx, res := range m.Resources {
			if res.Meta.String() == msg.resource.String() {
				item := m.list.Items()[idx].(*resourceItem) //nolint:forcetypeassert
				item.status = msg.status
				item.err = msg.err

				return m, nil
			}
		}
	}

	cmds := make([]tea.Cmd, 2) //nolint:mnd
	m.list, cmds[0] = m.list.Update(msg)
	m.spinner, cmds[1] = m.spinner.Update(msg)

	return m, tea.Batch(cmds...)
}

func (m *ContextModel) View() tea.View {
	m.view.SetContent(m.list.View())

	return m.view
}

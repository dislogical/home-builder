package bubbles

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/lipgloss/v2"

	"github.com/enescakir/emoji"

	tea "charm.land/bubbletea/v2"

	homebuilder "github.com/dislogical/home-builder/pkg"
)

type resourceStateMsg struct {
	resource homebuilder.Metadata

	status homebuilder.ResourceStatus
	err    error
}

type resourceItem struct {
	hbctx    *homebuilder.Context
	resource *homebuilder.Resource
	spinner  *spinner.Model

	status homebuilder.ResourceStatus
	err    error
}

var _ list.DefaultItem = (*resourceItem)(nil)

var resourceStyles = map[homebuilder.ResourceStatus]lipgloss.Style{
	homebuilder.StatusUpToDate: lipgloss.NewStyle().
		Foreground(lipgloss.Green).
		SetString(emoji.CheckMark.String()),

	homebuilder.StatusNeedsUpdate: lipgloss.NewStyle().
		Foreground(lipgloss.Yellow).
		SetString(emoji.CounterclockwiseArrowsButton.String()),

	homebuilder.StatusMissing: lipgloss.NewStyle().
		Faint(true).
		SetString(emoji.CrossMark.String()),

	homebuilder.StatusError: lipgloss.NewStyle().
		Foreground(lipgloss.Red).
		SetString(emoji.ExclamationMark.String()),

	homebuilder.StatusUnknown: lipgloss.NewStyle(),
}

func (i *resourceItem) Title() string {
	switch i.status {
	case homebuilder.StatusUnknown:
		return resourceStyles[i.status].Render(i.spinner.View(), i.resource.Meta.Name)
	default:
		return resourceStyles[i.status].Render(i.resource.Meta.Name)
	}
}

func (i *resourceItem) Description() string {
	if i.err != nil {
		return i.err.Error()
	}

	return i.resource.Meta.Type
}

func (i *resourceItem) FilterValue() string {
	return i.resource.Meta.String()
}

func initResource(res *homebuilder.Resource) tea.Cmd {
	return func() tea.Msg {
		err := res.Prepare()
		if err != nil {
			return resourceStateMsg{
				resource: res.Meta,
				status:   homebuilder.StatusError,
				err:      fmt.Errorf("preparing resource: %w", err),
			}
		}

		status, err := res.Backend.GetStatus()
		if err != nil {
			return resourceStateMsg{
				resource: res.Meta,
				status:   homebuilder.StatusError,
				err:      fmt.Errorf("retrieving status for %s: %w", res.Meta, err),
			}
		}

		return resourceStateMsg{
			resource: res.Meta,
			status:   status,
			err:      nil,
		}
	}
}

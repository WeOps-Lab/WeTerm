package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"weterm/constants"
	"weterm/index"
	"weterm/model"
	"weterm/pages"
	"weterm/pages/example"
)

type App struct {
	model *model.AppModel
}

func NewApp() *App {
	coreApp := tview.NewApplication()
	corePage := tview.NewPages()
	coreList := tview.NewList()

	coreApp.SetRoot(corePage, true)

	model := &model.AppModel{
		CoreApp:   coreApp,
		CorePages: corePage,
		CoreList:  coreList,
	}
	return &App{model: model}
}

func (receiver *App) setupPages() {
	index.SetUpMenuPage(receiver.model)
	pages.SetUpStatusPage(receiver.model)
	example.SetUpFormSamplePage(receiver.model)
}

func (receiver *App) setupInputCapture() {
	receiver.model.CoreApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			if receiver.model.CancelFunc != nil {
				receiver.model.CancelFunc()
			}
			receiver.model.CorePages.SwitchToPage(constants.MAIN_PAGE_NAME)
			return nil
		}
		return event
	})
}

func (receiver *App) Start() {
	receiver.setupPages()
	receiver.setupInputCapture()

	if err := receiver.model.CoreApp.Run(); err != nil {
		panic(err)
	}
}

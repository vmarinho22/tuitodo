package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	modalPageName = "modal"
	errorPageName = "error"
)

type ModalService struct {
	application  *tview.Application
	pages        *tview.Pages
	defaultFocus tview.Primitive
	lastFocused  tview.Primitive
	open         bool
}

func newModalService(application *tview.Application, pages *tview.Pages, defaultFocus tview.Primitive) *ModalService {
	return &ModalService{
		application:  application,
		pages:        pages,
		defaultFocus: defaultFocus,
		lastFocused:  defaultFocus,
	}
}

func (service *ModalService) IsOpen() bool {
	return service.open
}

func (service *ModalService) IsBlocking() bool {
	return service.open || service.frontPage() == errorPageName
}

func (service *ModalService) FilterKey(event *tcell.EventKey) *tcell.EventKey {
	if service.frontPage() == errorPageName {
		if isEscapeKey(event) {
			service.CloseError()
			return nil
		}
		return event
	}
	if service.open && isEscapeKey(event) {
		service.Close()
		return nil
	}
	return event
}

func (service *ModalService) CaptureEscape(event *tcell.EventKey) *tcell.EventKey {
	if isEscapeKey(event) {
		service.Close()
		return nil
	}
	return event
}

func (service *ModalService) Remember(primitive tview.Primitive) {
	if primitive != nil {
		service.lastFocused = primitive
	}
}

func (service *ModalService) OpenPage(content tview.Primitive) {
	service.rememberFocus()
	service.open = true
	if form, ok := content.(*tview.Form); ok {
		form.SetCancelFunc(service.Close)
		form.SetInputCapture(service.CaptureEscape)
	}
	overlay := newFittedCenter(content)
	overlay.SetInputCapture(service.CaptureEscape)
	service.pages.AddPage(modalPageName, overlay, true, true)
	service.application.SetFocus(content)
}

func (service *ModalService) OpenDialog(modal *tview.Modal) {
	service.rememberFocus()
	service.open = true
	modal.SetInputCapture(service.CaptureEscape)
	service.pages.AddPage(modalPageName, modal, true, true)
	service.application.SetFocus(modal)
}

func (service *ModalService) Close() {
	service.open = false
	service.pages.RemovePage(modalPageName)
	service.restoreFocus()
}

func (service *ModalService) ShowError(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			service.CloseError()
		})
	modal.SetTitle(" Erro ")
	service.pages.AddPage(errorPageName, modal, true, true)
	service.application.SetFocus(modal)
}

func (service *ModalService) CloseError() {
	service.pages.RemovePage(errorPageName)
	if service.open {
		if _, primitive := service.pages.GetFrontPage(); primitive != nil {
			service.application.SetFocus(primitive)
		}
		return
	}
	service.restoreFocus()
}

func (service *ModalService) frontPage() string {
	name, _ := service.pages.GetFrontPage()
	return name
}

func (service *ModalService) rememberFocus() {
	service.Remember(service.application.GetFocus())
}

func (service *ModalService) restoreFocus() {
	focus := service.lastFocused
	if focus == nil {
		focus = service.defaultFocus
	}
	service.application.SetFocus(focus)
}

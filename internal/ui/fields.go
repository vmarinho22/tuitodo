package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	inputBorderColor = tcell.NewRGBColor(200, 200, 200)
	dialogBackground = tcell.ColorBlack
	dialogLabelColor = tview.Styles.SecondaryTextColor
	dialogFieldColor = tview.Styles.PrimaryTextColor
)

type stackedField struct {
	*tview.Flex
	inner       tview.FormItem
	caption     *tview.TextView
	innerHeight int
}

func (field *stackedField) GetLabel() string {
	return ""
}

func (field *stackedField) GetFieldWidth() int {
	return 0
}

func (field *stackedField) GetFieldHeight() int {
	return field.innerHeight + 1
}

func (field *stackedField) SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) tview.FormItem {
	field.SetBackgroundColor(dialogBackground)
	field.inner.SetFormAttributes(0, labelColor, dialogBackground, fieldTextColor, tcell.ColorDefault)
	return field
}

func (field *stackedField) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
	field.inner.SetFinishedFunc(handler)
	return field
}

func (field *stackedField) SetDisabled(disabled bool) tview.FormItem {
	field.inner.SetDisabled(disabled)
	return field
}

func styleForm(form *tview.Form) {
	form.SetBackgroundColor(dialogBackground)
	form.SetFieldBackgroundColor(tcell.ColorDefault)
	form.SetFieldTextColor(dialogFieldColor)
	form.SetLabelColor(dialogLabelColor)
}

func newStyledForm(title string) *tview.Form {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" " + strings.TrimSpace(title) + " ")
	form.SetButtonsAlign(tview.AlignCenter)
	styleForm(form)
	return form
}

type messageField struct {
	*tview.TextView
	finished func(key tcell.Key)
	height   int
}

func (field *messageField) GetLabel() string { return "" }

func (field *messageField) GetFieldWidth() int { return 0 }

func (field *messageField) GetFieldHeight() int { return field.height }

func (field *messageField) SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) tview.FormItem {
	field.SetBackgroundColor(dialogBackground)
	field.SetTextColor(dialogFieldColor)
	return field
}

func (field *messageField) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
	field.finished = handler
	return field
}

func (field *messageField) SetDisabled(disabled bool) tview.FormItem {
	return field
}

func (field *messageField) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return field.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyTab, tcell.KeyBacktab, tcell.KeyEnter, tcell.KeyEscape:
			if field.finished != nil {
				field.finished(event.Key())
			}
		}
	})
}

func addStyledMessage(form *tview.Form, message string) {
	view := tview.NewTextView().
		SetText(message).
		SetTextColor(dialogFieldColor).
		SetWrap(true)
	view.SetBackgroundColor(dialogBackground)
	height := strings.Count(message, "\n") + 1
	if height < 2 {
		height = 2
	}
	form.AddFormItem(&messageField{TextView: view, height: height})
}

func focusFormButtons(form *tview.Form) {
	form.SetFocus(form.GetFormItemCount())
}

func bindFormButtonArrows(form *tview.Form, application *tview.Application) {
	count := form.GetButtonCount()
	if count == 0 || application == nil {
		return
	}
	itemCount := form.GetFormItemCount()
	for i := 0; i < count; i++ {
		index := i
		button := form.GetButton(index)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			offset := 0
			switch event.Key() {
			case tcell.KeyRight, tcell.KeyDown:
				offset = 1
			case tcell.KeyLeft, tcell.KeyUp:
				offset = -1
			default:
				return event
			}
			next := (index + offset + count) % count
			form.SetFocus(itemCount + next)
			application.SetFocus(form.GetButton(next))
			return nil
		})
	}
}

func addStyledInputField(form *tview.Form, label, value string) *tview.InputField {
	input := tview.NewInputField().
		SetText(value).
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetFieldTextColor(dialogFieldColor)
	styleInputBox(input.Box)
	form.AddFormItem(newStackedField(label, input, 3))
	return input
}

type categoryPicker struct {
	*tview.List
	finished func(key tcell.Key)
	caption  *tview.TextView
}

func (picker *categoryPicker) SetCaption(text string) {
	if picker.caption != nil {
		picker.caption.SetText(text)
	}
}

func (picker *categoryPicker) GetLabel() string { return "" }

func (picker *categoryPicker) GetFieldWidth() int { return 0 }

func (picker *categoryPicker) GetFieldHeight() int { return 6 }

func (picker *categoryPicker) SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) tview.FormItem {
	picker.SetBackgroundColor(dialogBackground)
	picker.SetMainTextColor(fieldTextColor)
	picker.SetSelectedTextColor(dialogBackground)
	picker.SetSelectedBackgroundColor(fieldTextColor)
	return picker
}

func (picker *categoryPicker) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
	picker.finished = handler
	return picker
}

func (picker *categoryPicker) SetDisabled(disabled bool) tview.FormItem {
	return picker
}

func (picker *categoryPicker) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	listHandler := picker.List.InputHandler()
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if isEscapeKey(event) {
			if picker.finished != nil {
				picker.finished(tcell.KeyEscape)
			}
			return
		}
		switch event.Key() {
		case tcell.KeyTab, tcell.KeyBacktab:
			if picker.finished != nil {
				picker.finished(event.Key())
			}
			return
		case tcell.KeyEnter:
			if picker.finished != nil {
				picker.finished(tcell.KeyTab)
			}
			return
		}
		if shortcutRune(event) == 'j' {
			event = tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		} else if shortcutRune(event) == 'k' {
			event = tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		}
		if listHandler != nil {
			listHandler(event, setFocus)
		}
	}
}

func addStyledCategoryPicker(form *tview.Form, label string, options []string, selected int) *categoryPicker {
	list := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true).
		SetWrapAround(true).
		SetSelectedFocusOnly(false)
	list.SetMainTextColor(dialogFieldColor)
	list.SetSelectedTextColor(dialogBackground)
	list.SetSelectedBackgroundColor(dialogFieldColor)
	styleInputBox(list.Box)
	for _, option := range options {
		list.AddItem(option, "", 0, nil)
	}
	if selected >= 0 && selected < len(options) {
		list.SetCurrentItem(selected)
	}
	picker := &categoryPicker{List: list}
	field := newStackedField(label, picker, 6)
	picker.caption = field.caption
	form.AddFormItem(field)
	return picker
}

func newStackedField(label string, inner tview.FormItem, innerHeight int) *stackedField {
	caption := tview.NewTextView().
		SetText(label).
		SetTextColor(dialogLabelColor)
	caption.SetBackgroundColor(dialogBackground)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(caption, 1, 0, false).
		AddItem(inner, innerHeight, 0, true)
	layout.SetBackgroundColor(dialogBackground)

	return &stackedField{Flex: layout, inner: inner, caption: caption, innerHeight: innerHeight}
}

func styleInputBox(box *tview.Box) {
	box.SetBorder(true)
	box.SetBorderColor(inputBorderColor)
	box.SetBorderAttributes(tcell.AttrDim)
	box.SetBackgroundColor(dialogBackground)
}

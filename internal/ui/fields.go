package ui

import (
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
	inner tview.FormItem
}

func (field *stackedField) GetLabel() string {
	return ""
}

func (field *stackedField) GetFieldWidth() int {
	return 0
}

func (field *stackedField) GetFieldHeight() int {
	return 4
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

func addStyledInputField(form *tview.Form, label, value string) *tview.InputField {
	input := tview.NewInputField().
		SetText(value).
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetFieldTextColor(dialogFieldColor)
	styleInputBox(input.Box)
	form.AddFormItem(newStackedField(label, input))
	return input
}

func addStyledDropDown(form *tview.Form, label string, options []string, selected int) *tview.DropDown {
	dropdown := tview.NewDropDown().
		SetOptions(options, nil).
		SetCurrentOption(selected).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetFieldTextColor(dialogFieldColor).
		SetFocusedStyle(tcell.StyleDefault.Foreground(dialogFieldColor).Background(tcell.ColorDefault))
	styleInputBox(dropdown.Box)
	form.AddFormItem(newStackedField(label, dropdown))
	return dropdown
}

func newStackedField(label string, inner tview.FormItem) *stackedField {
	caption := tview.NewTextView().
		SetText(label).
		SetTextColor(dialogLabelColor)
	caption.SetBackgroundColor(dialogBackground)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(caption, 1, 0, false).
		AddItem(inner, 3, 0, true)
	layout.SetBackgroundColor(dialogBackground)

	return &stackedField{Flex: layout, inner: inner}
}

func styleInputBox(box *tview.Box) {
	box.SetBorder(true)
	box.SetBorderColor(inputBorderColor)
	box.SetBorderAttributes(tcell.AttrDim)
	box.SetBackgroundColor(dialogBackground)
}

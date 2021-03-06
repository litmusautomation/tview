package tview

import (
	"github.com/gdamore/tcell"
)

// DefaultFormFieldWidth is the default field screen width of form elements
// whose field width is flexible (0). This is used in the Form class for
// horizontal layouts.
var DefaultFormFieldWidth = 10

// FormItem is the interface all form items must implement to be able to be
// included in a form.
type FormItem interface {
	Primitive

	// GetLabel returns the item's label text.
	GetLabel() string
	// GetLabelWidth returns the item's label width.
	GetLabelWidth() int

	// GetFieldWidth returns the item's field width.
	GetFieldWidth() int

	GetFieldAlign() (align int)

	GetBorderPadding() (top, bottom, left, right int)

	// SetFormAttributes sets a number of item attributes at once.
	SetFormAttributes(labelWidth, fieldWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) FormItem

	// GetFieldWidth returns the width of the form item's field (the area which
	// is manipulated by the user) in number of screen cells. A value of 0
	// indicates the the field width is flexible and may use as much space as
	// required.

	// SetEnteredFunc sets the handler function for when the user finished
	// entering data into the item. The handler may receive events for the
	// Enter key (we're done), the Escape key (cancel input), the Tab key (move to
	// next field), and the Backtab key (move to previous field).
	SetFinishedFunc(handler func(key tcell.Key)) FormItem

	GetID() int
}

// Form allows you to combine multiple one-line form elements into a vertical
// or horizontal layout. Form elements include types such as InputField or
// Checkbox. These elements can be optionally followed by one or more buttons
// for which you can define form-wide actions (e.g. Save, Clear, Cancel).
//
// See https://github.com/rivo/tview/wiki/Form for an example.
type Form struct {
	*Box

	// The items of the form (one row per item).
	items []FormItem

	// The buttons of the form.
	buttons []*Button

	activeButtons []*Button

	// If set to true, instead of position items and buttons from top to bottom,
	// they are positioned from left to right.
	horizontal bool

	// The alignment of the buttons.
	buttonsAlign int

	buttonsPaddingTop int
	buttonsIndent     int

	lastItem, lastButton int

	// The number of empty rows between items.
	itemPadding int

	// The index of the item or button which has focus. (Items are counted first,
	// buttons are counted last.)
	focusedElement int

	// The label color.
	labelColor tcell.Color

	// The background color of the input area.
	fieldBackgroundColor tcell.Color

	// The text color of the input area.
	fieldTextColor tcell.Color

	// The background color of the buttons.
	buttonBackgroundColor tcell.Color

	// The color of the button text.
	buttonTextColor tcell.Color

	// An optional function which is called when the user hits Escape.
	cancel func()

	// The content alignment, one of AlignLeft, AlignCenter, or AlignRight.
	align int

	itemsColumn []int

	columnPadding int
}

// NewForm returns a new form.
func NewForm() *Form {
	box := NewBox().SetBorderPadding(1, 1, 1, 1)

	f := &Form{
		Box:                   box,
		align:                 AlignLeft,
		columnPadding:         1,
		itemPadding:           1,
		buttonsPaddingTop:     2,
		buttonsIndent:         4,
		labelColor:            Styles.LabelTextColor,
		fieldBackgroundColor:  Styles.FieldBackgroundColor,
		fieldTextColor:        Styles.FieldTextColor,
		buttonBackgroundColor: Styles.ButtonBackgroundColor,
		buttonTextColor:       Styles.ButtonTextColor,
	}

	f.width = 0
	f.height = 0
	f.focus = f

	return f
}

// Buttons returns active buttons
func (f *Form) Buttons() []*Button {
	var buttons []*Button
	for i := 0; i < len(f.buttons); i++ {
		if !f.buttons[i].GetHidden() {
			buttons = append(buttons, f.buttons[i])
		}
	}
	return buttons
}

// SetAlign sets the content alignment within the flex. This must be
// either AlignLeft, AlignCenter, or AlignRight.
func (f *Form) SetAlign(align int) *Form {
	f.align = align
	return f
}

// SetItemPadding sets the number of empty rows between form items for vertical
// layouts and the number of empty cells between form items for horizontal
// layouts.
func (f *Form) SetItemPadding(padding int) *Form {
	f.itemPadding = padding
	return f
}

// SetColumnPadding sets the number of empty rows between form columns.
func (f *Form) SetColumnPadding(padding int) *Form {
	f.columnPadding = padding
	return f
}

// SetButtonPadding sets the number of empty rows between fields.
func (f *Form) SetButtonPadding(padding int) *Form {
	f.buttonsPaddingTop = padding
	return f
}

// SetButtonIndent makes indent between buttons
func (f *Form) SetButtonIndent(padding int) *Form {
	f.buttonsIndent = padding
	return f
}

// SetHorizontal sets the direction the form elements are laid out. If set to
// true, instead of positioning them from top to bottom (the default), they are
// positioned from left to right, moving into the next row if there is not
// enough space.
func (f *Form) SetHorizontal(horizontal bool) *Form {
	f.horizontal = horizontal
	return f
}

// SetLabelColor sets the color of the labels.
func (f *Form) SetLabelColor(color tcell.Color) *Form {
	f.labelColor = color
	return f
}

// SetFieldBackgroundColor sets the background color of the input areas.
func (f *Form) SetFieldBackgroundColor(color tcell.Color) *Form {
	f.fieldBackgroundColor = color
	return f
}

// SetFieldTextColor sets the text color of the input areas.
func (f *Form) SetFieldTextColor(color tcell.Color) *Form {
	f.fieldTextColor = color
	return f
}

// SetButtonsAlign sets how the buttons align horizontally, one of AlignLeft
// (the default), AlignCenter, and AlignRight. This is only
func (f *Form) SetButtonsAlign(align int) *Form {
	f.buttonsAlign = align
	return f
}

// SetButtonBackgroundColor sets the background color of the buttons.
func (f *Form) SetButtonBackgroundColor(color tcell.Color) *Form {
	f.buttonBackgroundColor = color
	return f
}

// SetButtonTextColor sets the color of the button texts.
func (f *Form) SetButtonTextColor(color tcell.Color) *Form {
	f.buttonTextColor = color
	return f
}

// AddInputField adds an input field to the form. It has a label, an optional
// initial value, a field width (a value of 0 extends it as far as possible),
// an optional accept function to validate the item's value (set to nil to
// accept any text), and an (optional) callback function which is invoked when
// the input field's text has changed.
func (f *Form) AddInputField(label, value string, fieldWidth int, accept func(textToCheck string, lastChar rune) bool, changed func(text string)) *Form {
	item := NewInputField().
		SetLabel(label).
		SetText(value).
		SetFieldWidth(fieldWidth).
		SetAcceptanceFunc(accept).
		SetChangedFunc(changed)
	f.items = append(f.items, item)
	return f
}

// AddPasswordField adds a password field to the form. This is similar to an
// input field except that the user's input not shown. Instead, a "mask"
// character is displayed. The password field has a label, an optional initial
// value, a field width (a value of 0 extends it as far as possible), and an
// (optional) callback function which is invoked when the input field's text has
// changed.
func (f *Form) AddPasswordField(label, value string, fieldWidth int, mask rune, changed func(text string)) *Form {
	if mask == 0 {
		mask = '*'
	}
	f.items = append(f.items, NewInputField().
		SetLabel(label).
		SetText(value).
		SetFieldWidth(fieldWidth).
		SetMaskCharacter(mask).
		SetChangedFunc(changed))
	return f
}

// AddDropDown adds a drop-down element to the form. It has a label, options,
// and an (optional) callback function which is invoked when an option was
// selected. The initial option may be a negative value to indicate that no
// option is currently selected.
func (f *Form) AddDropDown(label string, options []*DropDownOption, initialOption int, selected func(option *DropDownOption, optionIndex int)) *Form {
	f.items = append(f.items, NewDropDown().
		SetLabel(label).
		SetCurrentOption(initialOption).
		SetOptions(options, selected))
	return f
}

// AddRadioButton adds a radio button element to the form
func (f *Form) AddRadioButton(label string, options []*RadioOption, initialOption int, selected func(option *RadioOption, optionIndex int)) *Form {
	f.items = append(f.items, NewRadioButtons().
		SetLabel(label).
		SetCurrentOption(initialOption).
		SetOptions(options))
	return f
}

// AddCheckbox adds a checkbox to the form. It has a label, an initial state,
// and an (optional) callback function which is invoked when the state of the
// checkbox was changed by the user.
func (f *Form) AddCheckbox(label string, checked bool, changed func(checked bool)) *Form {
	f.items = append(f.items, NewCheckbox().
		SetLabel(label).
		SetChecked(checked).
		SetChangedFunc(changed))
	return f
}

// AddButton adds a new button to the form. The "selected" function is called
// when the user selects this button. It may be nil.
func (f *Form) AddButton(label string, selected func()) *Form {
	f.buttons = append(f.buttons, NewButton(label).SetSelectedFunc(selected))
	return f
}

// HiddenButton hides button by index
func (f *Form) HiddenButton(index int, state bool) *Form {
	if len(f.buttons) > index {
		f.buttons[index].SetHidden(state)
	}
	return f
}

// Clear removes all input elements from the form, including the buttons if
// specified.
func (f *Form) Clear(includeButtons bool) *Form {
	f.items = nil
	if includeButtons {
		f.buttons = nil
	}
	f.focusedElement = 0
	return f
}

// ClearButtons removes all buttons.
func (f *Form) ClearButtons() *Form {
	f.buttons = nil
	return f
}

// AddFormItem adds a new item to the form. This can be used to add your own
// objects to the form. Note, however, that the Form class will override some
// of its attributes to make it work in the form context. Specifically, these
// are:
//
//   - The label width
//   - The label color
//   - The background color
//   - The field text color
//   - The field background color
func (f *Form) AddFormItem(item FormItem) *Form {
	return f.AddFormItemWithColumn(item, 0)
}

// AddFormItemWithColumn adds a new item to the form and sets the column.
func (f *Form) AddFormItemWithColumn(item FormItem, column int) *Form {
	f.items = append(f.items, item)
	f.itemsColumn = append(f.itemsColumn, column)
	f.lastButton = len(f.items)
	return f
}

// GetFormItem returns the form element at the given position, starting with
// index 0. Elements are referenced in the order they were added. Buttons are
// not included.
func (f *Form) GetFormItem(index int) FormItem {
	return f.items[index]
}

// GetFormButton returns the form element at the given position, starting with
// index 0. Elements are referenced in the order they were added. Buttons are
// not included.
func (f *Form) GetFormButton(index int) *Button {
	return f.buttons[index]
}

// GetFormItemByLabel returns the first form element with the given label. If
// no such element is found, nil is returned. Buttons are not searched and will
// therefore not be returned.
func (f *Form) GetFormItemByLabel(label string) FormItem {
	for _, item := range f.items {
		if item.GetLabel() == label {
			return item
		}
	}
	return nil
}

// SetCancelFunc sets a handler which is called when the user hits the Escape
// key.
func (f *Form) SetCancelFunc(callback func()) *Form {
	f.cancel = callback
	return f
}

// GetRect returns the current position of the rectangle, x, y, width, and
// height.
func (f *Form) GetRect() (int, int, int, int) {
	x, y, width, height := f.Box.GetRect()

	maxColumns := f.getColoumnsCount()
	if width == 0 {
		maxWidth, _, _ := f.getMaxWidthItems()
		for column := 0; column < maxColumns; column++ {
			if width < maxWidth[column] {
				width = maxWidth[column]
			}
		}
		width += f.paddingLeft + f.paddingRight
		if f.border {
			width += 2
		}
	}

	if height == 0 {
		height = f.getMaxHeightColumn()

		if len(f.Buttons()) > 0 {
			height += 1 + f.buttonsPaddingTop
		}

		if f.border {
			height += 2
		}
		height += f.paddingTop + f.paddingBottom
	}
	return x, y, width, height
}

func (f *Form) getColoumnsCount() int {
	var maxColumns int
	for _, num := range f.itemsColumn {
		if maxColumns < num {
			maxColumns = num
		}
	}
	return maxColumns + 1
}

func (f *Form) getMaxHeightColumn() (height int) {
	maxColumns := f.getColoumnsCount()
	maxHeight := make([]int, maxColumns)

	if len(f.items) > 0 {
		for i := 0; i < len(f.items); i++ {
			column := f.itemsColumn[i]
			_, _, _, h := f.items[i].GetRect()
			maxHeight[column] += h + f.itemPadding
		}

		for column := 0; column < maxColumns; column++ {
			if height < maxHeight[column] {
				height = maxHeight[column]
			}
		}
		height -= f.itemPadding
	}

	return
}

func (f *Form) getMaxWidthItems() (maxWidth, maxLabelWidth, maxFieldWidth []int) {
	maxColumns := f.getColoumnsCount()

	// Find the longest label.
	maxLabelWidth = make([]int, maxColumns)
	maxFieldWidth = make([]int, maxColumns)

	for index, item := range f.items {
		_, _, leftPadding, rightPadding := item.GetBorderPadding()
		labelWidth := item.GetLabelWidth() + leftPadding
		fieldWidth := item.GetFieldWidth() + rightPadding
		column := f.itemsColumn[index]

		if labelWidth > 0 && labelWidth > maxLabelWidth[column]-1 && item.GetFieldAlign() == AlignCenter {
			maxLabelWidth[column] = labelWidth + 1
		}
		if fieldWidth > maxFieldWidth[column] {
			maxFieldWidth[column] = fieldWidth
		}
	}

	maxWidth = make([]int, maxColumns)
	for index, item := range f.items {
		_, _, leftPadding, rightPadding := item.GetBorderPadding()
		labelWidth := item.GetLabelWidth() + leftPadding
		fieldWidth := item.GetFieldWidth() + rightPadding
		column := f.itemsColumn[index]
		if labelWidth+fieldWidth > maxWidth[column] {
			maxWidth[column] = labelWidth + fieldWidth
		}
		if item.GetFieldAlign() == AlignCenter {
			maxWidth[column] = maxLabelWidth[column] + maxFieldWidth[column]
		}
	}

	return
}

// Draw draws this primitive onto the screen.
func (f *Form) Draw(screen tcell.Screen) {
	f.Box.Draw(screen)

	// Determine the dimensions.
	x, y, boxWidth, boxHeight := f.GetInnerRect()
	topLimit := y
	bottomLimit := y + boxHeight
	rightLimit := x + boxWidth
	startX := x

	maxColumns := f.getColoumnsCount()
	maxColWidth, maxColLabelWidth, _ := f.getMaxWidthItems()

	var maxWidth int
	for width, column := 0, 0; column < maxColumns; column++ {
		width += maxColWidth[column]
		maxWidth = width
		width += 1 + f.columnPadding
	}

	switch f.align {
	case AlignCenter:
		x += (boxWidth - maxWidth) / 2
	case AlignRight:
		x += boxWidth - maxWidth
	}

	var colX, colY []int
	for width, column := 0, 0; column < maxColumns; column++ {
		colX = append(colX, x+width)
		colY = append(colY, topLimit)
		width += maxColWidth[column] + 1 + f.columnPadding
	}

	// Calculate positions of form items.
	positions := make([]struct{ x, y, width, height int }, len(f.items)+len(f.Buttons()))
	var focusedPosition struct{ x, y, width, height int }
	for index, item := range f.items {
		column := f.itemsColumn[index]
		x := colX[column]
		y := colY[column]
		_, _, leftPadding, rightPadding := item.GetBorderPadding()
		labelWidth := item.GetLabelWidth()
		fieldWidth := item.GetFieldWidth()
		var itemWidth int
		if f.horizontal {
			if fieldWidth == 0 {
				fieldWidth = DefaultFormFieldWidth
			}
			itemWidth = labelWidth + fieldWidth
		} else {
			itemWidth = boxWidth
			// Implement alignment of field
			switch item.GetFieldAlign() {
			case AlignCenter:
				labelWidth = maxColLabelWidth[column] - leftPadding
				fieldWidth = maxColWidth[column] - leftPadding - labelWidth - rightPadding
			case AlignRight:
				labelWidth = maxColWidth[column] - leftPadding - fieldWidth - rightPadding
			case AlignLeft:
				fieldWidth = maxColWidth[column] - leftPadding - labelWidth - rightPadding
			}
		}

		// Advance to next line if there is no space.
		if f.horizontal && x+labelWidth+1 >= rightLimit {
			x = startX
			y += 2
		}

		// Adjust the item's attributes.
		if x+itemWidth >= rightLimit {
			itemWidth = rightLimit - x
		}
		item.SetFormAttributes(
			labelWidth,
			fieldWidth,
			f.labelColor,
			tcell.ColorBlack,
			f.fieldTextColor,
			f.fieldBackgroundColor,
		)

		// Save position.
		positions[index].x = x
		positions[index].y = y
		positions[index].width = maxColWidth[column]
		_, _, _, positions[index].height = item.GetRect()
		if item.GetFocusable().HasFocus() {
			focusedPosition = positions[index]
		}

		// Advance to next item.
		if f.horizontal {
			x += itemWidth + f.itemPadding
		} else {
			y += positions[index].height + f.itemPadding
		}
		colY[column] = y
	}

	y = topLimit + f.getMaxHeightColumn()

	// How wide are the buttons?
	buttonWidths := make([]int, len(f.Buttons()))
	buttonsWidth := 0
	for index, button := range f.Buttons() {
		w := StringWidth(button.GetLabel()) + 4
		buttonWidths[index] = w
		buttonsWidth += w + f.buttonsIndent
	}
	buttonsWidth--

	// Where do we place them?
	if !f.horizontal && x+buttonsWidth < rightLimit {
		if f.buttonsAlign == AlignRight {
			x = rightLimit - buttonsWidth
		} else if f.buttonsAlign == AlignCenter {
			x = (boxWidth-buttonsWidth)/2 + startX
		}

		// In vertical layouts, buttons always appear after an empty line.
		// if f.itemPadding == 0 {
		// 	y++
		// }

	}

	if len(f.Buttons()) > 0 {
		y += f.buttonsPaddingTop
	}

	// Calculate positions of buttons.
	for index, button := range f.Buttons() {
		space := rightLimit - x
		buttonWidth := buttonWidths[index]
		if f.horizontal {
			if space < buttonWidth-4 {
				x = startX
				y += 2
				space = boxWidth
			}
		} else {
			if space < 1 {
				break // No space for this button anymore.
			}
		}
		if buttonWidth > space {
			buttonWidth = space
		}
		button.SetLabelColor(f.buttonTextColor).
			SetLabelColorActivated(f.buttonBackgroundColor).
			SetBackgroundColorActivated(f.buttonTextColor).
			SetBackgroundColor(f.buttonBackgroundColor)

		buttonIndex := index + len(f.items)
		positions[buttonIndex].x = x
		positions[buttonIndex].y = y
		positions[buttonIndex].width = buttonWidth
		positions[buttonIndex].height = 1

		if button.HasFocus() {
			focusedPosition = positions[buttonIndex]
		}

		x += buttonWidth + f.buttonsIndent
	}

	// Determine vertical offset based on the position of the focused item.
	var offset int
	if focusedPosition.y+focusedPosition.height > bottomLimit {
		offset = focusedPosition.y + focusedPosition.height - bottomLimit
		if focusedPosition.y-offset < topLimit {
			offset = focusedPosition.y - topLimit
		}
	}

	// Draw items.
	for index, item := range f.items {
		// Set position.
		y := positions[index].y - offset
		height := positions[index].height
		item.SetRect(positions[index].x, y, positions[index].width, height)

		// Is this item visible?
		if y+height <= topLimit || (y > bottomLimit) {
			continue
		}

		// Draw items with focus last (in case of overlaps).
		if item.GetFocusable().HasFocus() {
			defer item.Draw(screen)
		} else {
			item.Draw(screen)
		}
	}

	// Draw buttons.
	for index, button := range f.Buttons() {
		// Set position.
		buttonIndex := index + len(f.items)
		y := positions[buttonIndex].y - offset
		height := positions[buttonIndex].height
		button.SetRect(positions[buttonIndex].x, y, positions[buttonIndex].width, height)

		// Is this button visible?
		if y+height <= topLimit || (y >= bottomLimit && boxHeight > 0) {
			continue
		}

		// Draw button.
		button.Draw(screen)
	}
}

// ResetFocus sets focus on the first element
func (f *Form) ResetFocus() {
	for i, item := range f.items {
		if item.IsDisable() {
			continue
		}
		f.focusedElement = i
		break
	}
	f.lastItem = 0
	f.lastButton = len(f.items)
}

// Focus is called by the application when the primitive receives focus.
func (f *Form) Focus(delegate func(p Primitive)) {
	if len(f.items)+len(f.Buttons()) == 0 {
		return
	}

	// Hand on the focus to one of our child elements.
	if f.focusedElement < 0 || f.focusedElement >= len(f.items)+len(f.Buttons()) {
		f.focusedElement = 0
	}

	var (
		nextStep = func(indexes ...int) {
			for i := 1; i < len(indexes); i++ {
				if indexes[i-1] < indexes[i] && f.focusedElement < indexes[i] || indexes[i-1] > indexes[i] && f.focusedElement > indexes[i] {
					indexes = append(indexes[i:], indexes[0:i]...)
					break
				}
			}

			current := f.focusedElement
			for i := 0; i < len(indexes); i++ {
				if indexes[i] >= len(f.items)+len(f.Buttons()) {
					continue
				}
				f.focusedElement = indexes[i]
				if f.focusedElement < len(f.items) && current < len(f.items) && f.items[current].GetID() == f.items[f.focusedElement].GetID() {
					continue
				}
				if f.focusedElement >= len(f.items) && current >= len(f.items) && f.Buttons()[current-len(f.items)].GetID() == f.Buttons()[f.focusedElement-len(f.items)].GetID() {
					continue
				}

				f.Focus(delegate)
				if f.focusedElement < len(f.items) && f.items[f.focusedElement].GetFocusable().HasFocus() {
					f.lastItem = f.focusedElement
					break
				}
				if f.focusedElement >= len(f.items) && f.Buttons()[len(f.Buttons())-1-(f.focusedElement-len(f.items))].GetFocusable().HasFocus() {
					f.lastButton = f.focusedElement
					break
				}
			}
		}

		makeRange = func(start, end int) []int {
			var a []int
			if start < end {
				a = make([]int, end-start+1)
				for i := range a {
					a[i] = start + i
				}
			} else {
				a = make([]int, start-end+1)
				for i := range a {
					a[i] = start - i
				}
			}
			return a
		}
	)

	itemHandler := func(key tcell.Key) {
		switch key {
		case tcell.KeyTab:
			nextStep(makeRange(0, len(f.items))...)
		case tcell.KeyBacktab:
			nextStep(makeRange(len(f.items), 0)...)
		case tcell.KeyEnter:
			nextStep(makeRange(0, len(f.items))...)
		case tcell.KeyUp:
			nextStep(makeRange(len(f.items)-1, 0)...)
		case tcell.KeyDown:
			nextStep(makeRange(0, len(f.items)-1)...)
		case tcell.KeyEscape:
			if f.cancel != nil {
				f.cancel()
			} else {
				f.focusedElement = 0
				f.Focus(delegate)
			}
		}
	}

	buttonHandler := func(key tcell.Key) {
		switch key {
		case tcell.KeyTab:
			nextStep(makeRange(0, len(f.items))...)
		case tcell.KeyBacktab:
			nextStep(makeRange(len(f.items), 0)...)
		case tcell.KeyUp, tcell.KeyLeft:
			nextStep(makeRange(len(f.items), len(f.items)+len(f.Buttons())-1)...)
		case tcell.KeyDown, tcell.KeyRight:
			nextStep(makeRange(len(f.items)+len(f.Buttons())-1, len(f.items))...)
		case tcell.KeyEscape:
			if f.cancel != nil {
				f.cancel()
			} else {
				f.focusedElement = 0
				f.Focus(delegate)
			}
		}
	}

	if f.focusedElement < len(f.items) {
		// We're selecting an item.
		item := f.items[f.focusedElement]
		item.SetFinishedFunc(itemHandler)
		delegate(item)
	} else {
		// We're selecting a button.
		//		fmt.Println(len(f.buttons) - 1 - (f.focusedElement - len(f.items)))
		button := f.Buttons()[len(f.Buttons())-1-(f.focusedElement-len(f.items))]
		button.SetBlurFunc(buttonHandler)
		delegate(button)
	}
}

// HasFocus returns whether or not this primitive has focus.
func (f *Form) HasFocus() bool {
	for _, item := range f.items {
		if item.GetFocusable().HasFocus() {
			return true
		}
	}
	for _, button := range f.Buttons() {
		if button.focus.HasFocus() {
			return true
		}
	}
	return false
}

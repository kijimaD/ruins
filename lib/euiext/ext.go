package euiext

import (
	img "image"
	"image/color"
	"math"
	"slices"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type List struct {
	EntrySelectedEvent *event.Event

	containerOpts               []widget.ContainerOpt
	scrollContainerOpts         []widget.ScrollContainerOpt
	sliderOpts                  []widget.SliderOpt
	entries                     []any
	entryLabelFunc              ListEntryLabelFunc
	entryFace                   text.Face
	entryUnselectedColor        *widget.ButtonImage
	entrySelectedColor          *widget.ButtonImage
	entryUnselectedTextColor    *widget.ButtonTextColor
	entryTextColor              *widget.ButtonTextColor
	entryTextPadding            widget.Insets
	entryTextHorizontalPosition widget.TextPosition
	entryTextVerticalPosition   widget.TextPosition
	controlWidgetSpacing        int
	hideHorizontalSlider        bool
	hideVerticalSlider          bool
	allowReselect               bool
	selectFocus                 bool

	init            *widget.MultiOnce
	container       *widget.Container
	listContent     *widget.Container
	scrollContainer *widget.ScrollContainer
	vSlider         *widget.Slider
	hSlider         *widget.Slider
	buttons         []*widget.Button
	selectedEntry   any

	disableDefaultKeys bool
	focused            bool
	tabOrder           int
	justMoved          bool
	focusIndex         int
	prevFocusIndex     int

	focusMap map[widget.FocusDirection]widget.Focuser

	// 独自拡張
	entryEnterFunc ListEntryEnterFunc
}

type ListOpt func(l *List)

type ListEntryLabelFunc func(e any) string
type ListEntryEnterFunc func(e any)

type ListEntryColor struct {
	Unselected                 color.Color
	Selected                   color.Color
	DisabledUnselected         color.Color
	DisabledSelected           color.Color
	SelectingBackground        color.Color
	SelectedBackground         color.Color
	FocusedBackground          color.Color
	SelectingFocusedBackground color.Color
	SelectedFocusedBackground  color.Color
	DisabledSelectedBackground color.Color
}

type ListEntrySelectedEventArgs struct {
	List          *List
	Entry         any
	PreviousEntry any
}

type ListEntrySelectedHandlerFunc func(args *ListEntrySelectedEventArgs)

type ListOptions struct{}

var ListOpts ListOptions

func NewList(listOpts ...ListOpt) *List {
	l := &List{
		EntrySelectedEvent: &event.Event{},

		entryTextHorizontalPosition: widget.TextPositionCenter,
		entryTextVerticalPosition:   widget.TextPositionCenter,

		init:           &widget.MultiOnce{},
		focusIndex:     0,
		prevFocusIndex: -1,
		focusMap:       make(map[widget.FocusDirection]widget.Focuser),
	}

	l.init.Append(l.createWidget)

	for _, o := range listOpts {
		o(l)
	}

	l.resetFocusIndex()

	l.validate()

	return l
}

func (l *List) validate() {
	if len(l.scrollContainerOpts) == 0 {
		panic("List: ScrollContainerOpts are required.")
	}
	if len(l.sliderOpts) == 0 {
		panic("List: SliderOpts are required.")
	}
	if l.entryFace == nil {
		panic("List: EntryFontFace is required.")
	}
	if l.entryLabelFunc == nil {
		panic("List: EntryLabelFunc is required.")
	}
	if l.entryEnterFunc == nil {
		panic("List: EntryEnterFunc is required.")
	}
	if l.entryTextColor == nil || l.entryTextColor.Idle == nil {
		panic("List: ListEntryColor.Selected is required.")
	}
	if l.entryUnselectedTextColor == nil || l.entryUnselectedTextColor.Idle == nil {
		panic("List: ListEntryColor.Unselected is required.")
	}
}

func (o ListOptions) ContainerOpts(opts ...widget.ContainerOpt) ListOpt {
	return func(l *List) {
		l.containerOpts = append(l.containerOpts, opts...)
	}
}

func (o ListOptions) ScrollContainerOpts(opts ...widget.ScrollContainerOpt) ListOpt {
	return func(l *List) {
		l.scrollContainerOpts = append(l.scrollContainerOpts, opts...)
	}
}

func (o ListOptions) SliderOpts(opts ...widget.SliderOpt) ListOpt {
	return func(l *List) {
		l.sliderOpts = append(l.sliderOpts, opts...)
	}
}

func (o ListOptions) ControlWidgetSpacing(s int) ListOpt {
	return func(l *List) {
		l.controlWidgetSpacing = s
	}
}

func (o ListOptions) HideHorizontalSlider() ListOpt {
	return func(l *List) {
		l.hideHorizontalSlider = true
	}
}

func (o ListOptions) HideVerticalSlider() ListOpt {
	return func(l *List) {
		l.hideVerticalSlider = true
	}
}

func (o ListOptions) Entries(e []any) ListOpt {
	return func(l *List) {
		l.entries = slices.CompactFunc(e, func(a any, b any) bool { return a == b })
	}
}

func (o ListOptions) EntryLabelFunc(f ListEntryLabelFunc) ListOpt {
	return func(l *List) {
		l.entryLabelFunc = f
	}
}

func (o ListOptions) EntryEnterFunc(f ListEntryEnterFunc) ListOpt {
	return func(l *List) {
		l.entryEnterFunc = f
	}
}

func (o ListOptions) EntryFontFace(f text.Face) ListOpt {
	return func(l *List) {
		l.entryFace = f
	}
}

func (o ListOptions) DisableDefaultKeys(val bool) ListOpt {
	return func(l *List) {
		l.disableDefaultKeys = val
	}
}

func (o ListOptions) EntryColor(c *ListEntryColor) ListOpt {
	return func(l *List) {
		l.entryUnselectedColor = &widget.ButtonImage{
			Idle:     image.NewNineSliceColor(color.Transparent),
			Disabled: image.NewNineSliceColor(color.Transparent),
			Hover:    image.NewNineSliceColor(c.FocusedBackground),
			Pressed:  image.NewNineSliceColor(c.SelectingBackground),
		}

		l.entrySelectedColor = &widget.ButtonImage{
			Idle:     image.NewNineSliceColor(c.SelectedBackground),
			Disabled: image.NewNineSliceColor(c.DisabledSelectedBackground),
			Hover:    image.NewNineSliceColor(c.SelectedFocusedBackground),
			Pressed:  image.NewNineSliceColor(c.SelectingFocusedBackground),
		}

		l.entryUnselectedTextColor = &widget.ButtonTextColor{
			Idle:     c.Unselected,
			Disabled: c.DisabledUnselected,
		}

		l.entryTextColor = &widget.ButtonTextColor{
			Idle:     c.Selected,
			Disabled: c.DisabledSelected,
		}
	}
}

func (o ListOptions) EntryTextPadding(i widget.Insets) ListOpt {
	return func(l *List) {
		l.entryTextPadding = i
	}
}

// EntryTextPosition sets the position of the text for entries.
// Defaults to both TextPositionCenter.
func (o ListOptions) EntryTextPosition(h widget.TextPosition, v widget.TextPosition) ListOpt {
	return func(l *List) {
		l.entryTextHorizontalPosition = h
		l.entryTextVerticalPosition = v
	}
}

func (o ListOptions) EntrySelectedHandler(f ListEntrySelectedHandlerFunc) ListOpt {
	return func(l *List) {
		l.EntrySelectedEvent.AddHandler(func(args any) {
			f(args.(*ListEntrySelectedEventArgs))
		})
	}
}

func (o ListOptions) AllowReselect() ListOpt {
	return func(l *List) {
		l.allowReselect = true
	}
}

// SelectFocus automatically selects each focused entry.
func (o ListOptions) SelectFocus() ListOpt {
	return func(l *List) {
		l.selectFocus = true
	}
}

func (o ListOptions) TabOrder(tabOrder int) ListOpt {
	return func(l *List) {
		l.tabOrder = tabOrder
	}
}

func (l *List) GetWidget() *widget.Widget {
	l.init.Do()
	return l.container.GetWidget()
}

func (l *List) PreferredSize() (int, int) {
	l.init.Do()
	return l.container.PreferredSize()
}

func (l *List) SetLocation(rect img.Rectangle) {
	l.init.Do()
	l.container.GetWidget().Rect = rect
}

func (l *List) RequestRelayout() {
	l.init.Do()
	l.container.RequestRelayout()
}

func (l *List) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	l.init.Do()
	l.container.SetupInputLayer(def)
}

func (l *List) Render(screen *ebiten.Image) {
	l.init.Do()

	d := l.container.GetWidget().Disabled
	if l.vSlider != nil {
		l.vSlider.DrawTrackDisabled = d
	}
	if l.hSlider != nil {
		l.hSlider.DrawTrackDisabled = d
	}
	l.scrollContainer.GetWidget().Disabled = d

	if l.focusIndex != l.prevFocusIndex && l.focusIndex >= 0 && l.focusIndex < len(l.buttons) {
		l.scrollVisible(l.buttons[l.focusIndex])
	}

	if l.selectFocus {
		l.SelectFocused()
	}

	l.container.Render(screen)
}

func (l *List) Update() {
	l.init.Do()

	l.handleInput()
	l.container.Update()
}

/** Focuser Interface - Start **/

func (l *List) Focus(focused bool) {
	l.init.Do()
	l.GetWidget().FireFocusEvent(l, focused, img.Point{-1, -1})
	l.focused = focused
}

func (l *List) IsFocused() bool {
	return l.focused
}

func (l *List) TabOrder() int {
	return l.tabOrder
}

func (l *List) GetFocus(direction widget.FocusDirection) widget.Focuser {
	return l.focusMap[direction]
}

func (l *List) AddFocus(direction widget.FocusDirection, focus widget.Focuser) {
	l.focusMap[direction] = focus
}

/** Focuser Interface - End **/

func (l *List) handleInput() {
	if l.focused && !l.GetWidget().Disabled && len(l.buttons) > 0 {
		if !l.disableDefaultKeys && (input.KeyPressed(ebiten.KeyUp) || input.KeyPressed(ebiten.KeyDown)) {
			if !l.justMoved {
				if input.KeyPressed(ebiten.KeyDown) {
					l.FocusNext()
				} else {
					l.FocusPrevious()
				}
			}
		} else {
			l.justMoved = false
		}
		l.buttons[l.focusIndex].Focus(true)
	} else if len(l.buttons) > 0 && l.focusIndex <= len(l.buttons) {
		l.buttons[l.focusIndex].Focus(false)
	}
}

func (l *List) FocusNext() {
	if len(l.buttons) > 0 {
		direction := 1
		l.buttons[l.focusIndex].Focus(false)
		l.prevFocusIndex = l.focusIndex
		l.focusIndex += direction
		if l.focusIndex < 0 {
			l.focusIndex = len(l.buttons) - 1
		}
		if l.focusIndex >= len(l.buttons) {
			l.focusIndex = 0
		}
		l.justMoved = true
		l.buttons[l.focusIndex].Focus(true)
	}
}

func (l *List) FocusPrevious() {
	if len(l.buttons) > 0 {
		direction := -1
		l.buttons[l.focusIndex].Focus(false)
		l.prevFocusIndex = l.focusIndex
		l.focusIndex += direction
		if l.focusIndex < 0 {
			l.focusIndex = len(l.buttons) - 1
		}
		if l.focusIndex >= len(l.buttons) {
			l.focusIndex = 0
		}
		l.justMoved = true
		l.buttons[l.focusIndex].Focus(true)
	}
}

func (l *List) SelectFocused() {
	if l.focusIndex >= 0 && l.focusIndex < len(l.buttons) {
		l.scrollVisible(l.buttons[l.focusIndex])
		l.setSelectedEntry(l.entries[l.focusIndex], false)
	}
}

func (l *List) resetFocusIndex() {
	if len(l.buttons) > 0 {
		if l.focusIndex != -1 && l.focusIndex < len(l.buttons) {
			l.buttons[l.focusIndex].Focus(false)
		}
		for i := 0; i < len(l.entries); i++ {
			if l.entries[i] == l.selectedEntry {
				if i != l.focusIndex {
					l.prevFocusIndex = l.focusIndex
					l.focusIndex = i
				}
				return
			}
		}
		l.focusIndex = 0
	}
}

func (l *List) createWidget() {
	var cols int
	if l.hideVerticalSlider {
		cols = 1
	} else {
		cols = 2
	}

	l.container = widget.NewContainer(
		append([]widget.ContainerOpt{
			widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.TrackHover(true)),
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(cols),
				widget.GridLayoutOpts.Stretch([]bool{true, false}, []bool{true, false}),
				widget.GridLayoutOpts.Spacing(l.controlWidgetSpacing, l.controlWidgetSpacing),
			))}, l.containerOpts...)...,
	)

	l.listContent = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical))),
		widget.ContainerOpts.AutoDisableChildren(),
	)

	l.buttons = make([]*widget.Button, 0, len(l.entries))
	for _, e := range l.entries {
		e := e
		but := l.createEntry(e)

		l.buttons = append(l.buttons, but)
		l.listContent.AddChild(but)
	}

	l.scrollContainer = widget.NewScrollContainer(append(l.scrollContainerOpts, []widget.ScrollContainerOpt{
		widget.ScrollContainerOpts.Content(l.listContent),
		widget.ScrollContainerOpts.StretchContentWidth(),
	}...)...)

	l.container.AddChild(l.scrollContainer)

	if !l.hideVerticalSlider {
		pageSizeFunc := func() int {
			return int(math.Round(float64(l.scrollContainer.ViewRect().Dy()) / float64(l.listContent.GetWidget().Rect.Dy()) * 1000))
		}

		l.vSlider = widget.NewSlider(append(l.sliderOpts, []widget.SliderOpt{
			widget.SliderOpts.Direction(widget.DirectionVertical),
			widget.SliderOpts.MinMax(0, 1000),
			widget.SliderOpts.PageSizeFunc(pageSizeFunc),
			widget.SliderOpts.DisableDefaultKeys(l.disableDefaultKeys),
			widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
				current := args.Slider.Current
				if pageSizeFunc() >= 1000 {
					current = 0
				}
				l.scrollContainer.ScrollTop = float64(current) / 1000
			}),
		}...)...)
		l.container.AddChild(l.vSlider)

		l.scrollContainer.GetWidget().ScrolledEvent.AddHandler(func(args any) {
			a := args.(*widget.WidgetScrolledEventArgs)
			p := pageSizeFunc() / 3
			if p < 1 {
				p = 1
			}
			l.vSlider.Current -= int(math.Round(a.Y * float64(p)))
		})
	}

	if !l.hideHorizontalSlider {
		l.hSlider = widget.NewSlider(append(l.sliderOpts, []widget.SliderOpt{
			widget.SliderOpts.Direction(widget.DirectionHorizontal),
			widget.SliderOpts.MinMax(0, 1000),
			widget.SliderOpts.PageSizeFunc(func() int {
				return int(math.Round(float64(l.scrollContainer.ViewRect().Dx()) / float64(l.listContent.GetWidget().Rect.Dx()) * 1000))
			}),
			widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
				l.scrollContainer.ScrollLeft = float64(args.Slider.Current) / 1000
			}),
		}...)...)
		l.container.AddChild(l.hSlider)
	}

}

// Updates the entries in the list.
// Note: Duplicates will be removed.
func (l *List) SetEntries(newEntries []any) {
	//Remove old entries
	for i := range l.entries {
		but := l.buttons[i]
		l.listContent.RemoveChild(but)
	}
	l.entries = nil
	l.buttons = nil

	//Add new Entries
	for idx := range newEntries {
		if !slices.ContainsFunc(l.entries, func(cmp any) bool {
			return cmp == newEntries[idx]
		}) {
			l.entries = append(l.entries, newEntries[idx])
			but := l.createEntry(newEntries[idx])
			l.buttons = append(l.buttons, but)
			l.listContent.AddChild(but)
		}
	}
	l.selectedEntry = nil
	l.resetFocusIndex()
}

// Remove the passed in entry from the list if it exists
func (l *List) RemoveEntry(entry any) {
	l.init.Do()

	if len(l.entries) > 0 && entry != nil {
		for i, e := range l.entries {
			if e == entry {
				but := l.buttons[i]
				l.entries = append(l.entries[:i], l.entries[i+1:]...)
				l.buttons = append(l.buttons[:i], l.buttons[i+1:]...)
				l.listContent.RemoveChild(but)

				entryLen := len(l.entries)
				if l.focusIndex >= entryLen {
					l.focusIndex = i - 1
				}

				if l.focusIndex >= 0 && l.focusIndex < entryLen {
					l.setSelectedEntry(l.entries[l.focusIndex], false)
				}
				break
			}
		}
		l.resetFocusIndex()
	}
}

// Add a new entry to the end of the list
// Note: Duplicates will not be added
func (l *List) AddEntry(entry any) {
	l.init.Do()
	if !l.checkForDuplicates(l.entries, entry) {
		l.entries = append(l.entries, entry)
		but := l.createEntry(entry)
		l.buttons = append(l.buttons, but)
		l.listContent.AddChild(but)
	}
	l.resetFocusIndex()
}

// Return the current entries in the list
func (l *List) Entries() []any {
	l.init.Do()
	return l.entries
}

// Return the currently selected entry in the list
func (l *List) SelectedEntry() any {
	l.init.Do()
	return l.selectedEntry
}

// Set the Selected Entry to e if it is found.
func (l *List) SetSelectedEntry(entry any) {
	l.setSelectedEntry(entry, false)
}

func (l *List) setSelectedEntry(e any, user bool) {
	if e != l.selectedEntry || (user && l.allowReselect) {
		l.init.Do()

		prev := l.selectedEntry
		l.selectedEntry = e
		l.resetFocusIndex()
		for i, b := range l.buttons {
			if l.entries[i] == e {
				b.Image = l.entrySelectedColor
				b.TextColor = l.entryTextColor
			} else {
				b.Image = l.entryUnselectedColor
				b.TextColor = l.entryUnselectedTextColor
			}
		}

		l.EntrySelectedEvent.Fire(&ListEntrySelectedEventArgs{
			Entry:         e,
			PreviousEntry: prev,
		})
	}
}

func (l *List) checkForDuplicates(entries []any, entry any) bool {
	return slices.ContainsFunc(entries, func(cmp any) bool {
		return cmp == entry
	})
}

func (l *List) createEntry(entry any) *widget.Button {
	but := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.Image(l.entryUnselectedColor),
		widget.ButtonOpts.Text(l.entryLabelFunc(entry), l.entryFace, l.entryUnselectedTextColor),
		widget.ButtonOpts.TextPadding(l.entryTextPadding),
		widget.ButtonOpts.TextPosition(l.entryTextHorizontalPosition, l.entryTextVerticalPosition),
		widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
			l.setSelectedEntry(entry, true)
		}),
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			l.entryEnterFunc(entry)
		}),
	)

	return but
}

func (l *List) setScrollTop(t float64) {
	l.init.Do()
	if l.vSlider != nil {
		l.vSlider.Current = int(math.Round(t * 1000))
	}
	l.scrollContainer.ScrollTop = t
}

func (l *List) setScrollLeft(left float64) {
	l.init.Do()
	if l.hSlider != nil {
		l.hSlider.Current = int(math.Round(left * 1000))
	}
	l.scrollContainer.ScrollLeft = left
}

func (l *List) scrollVisible(w widget.HasWidget) {
	vrect := l.scrollContainer.ViewRect()
	wrect := w.GetWidget().Rect
	if !wrect.In(vrect) {
		crect := l.scrollContainer.ContentRect()
		scrollTop := l.scrollContainer.ScrollTop
		scrollHeight := crect.Dy() - vrect.Dy()
		if wrect.Max.Y > vrect.Max.Y {
			scrollTop = float64(wrect.Max.Y-vrect.Dy()-crect.Min.Y) / float64(scrollHeight)
		} else if wrect.Min.Y < vrect.Min.Y {
			scrollTop = float64(wrect.Min.Y-crect.Min.Y) / float64(scrollHeight)
		}
		scrollLeft := l.scrollContainer.ScrollLeft
		scrollWidth := crect.Dx() - vrect.Dx()
		if wrect.Max.X > vrect.Max.X {
			scrollLeft = float64(wrect.Max.X-vrect.Dx()-crect.Min.X) / float64(scrollWidth)
		} else if wrect.Min.X < vrect.Min.X {
			scrollLeft = float64(wrect.Min.X-crect.Min.X) / float64(scrollWidth)
		}
		l.setScrollTop(scrollClamp(scrollTop, l.scrollContainer.ScrollTop))
		l.setScrollLeft(scrollClamp(scrollLeft, l.scrollContainer.ScrollLeft))
	} else if wrect != vrect {
		l.prevFocusIndex = l.focusIndex
	}
}

func scrollClamp(targetScroll, currentScroll float64) float64 {
	const maxScrollStep = 0.1
	minScroll := currentScroll - maxScrollStep
	maxScroll := currentScroll + maxScrollStep
	return math.Max(minScroll, math.Min(targetScroll, maxScroll))
}

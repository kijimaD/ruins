package euiext

import (
	"image/color"
	"math"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type List struct {
	EntrySelectedEvent *event.Event

	containerOpts               []widget.ContainerOpt
	scrollContainerOpts         []widget.ScrollContainerOpt
	sliderOpts                  []widget.SliderOpt
	entries                     []any
	entryLabelFunc              widget.ListEntryLabelFunc
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
}

type ListOpt func(l *List)

type ListEntryLabelFunc func(e any) string

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
		}))

	return but
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
	if l.entryTextColor == nil || l.entryTextColor.Idle == nil {
		panic("List: ListEntryColor.Selected is required.")
	}
	if l.entryUnselectedTextColor == nil || l.entryUnselectedTextColor.Idle == nil {
		panic("List: ListEntryColor.Unselected is required.")
	}
}

func NewList(opts ...ListOpt) *List {
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

	for _, o := range opts {
		o(l)
	}

	l.resetFocusIndex()

	l.validate()

	return l
}

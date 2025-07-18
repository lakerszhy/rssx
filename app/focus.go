package app

const (
	focusFeed focus = iota
	focusItem
	focusPreview
)

type focus int

func (f focus) prev() focus {
	switch f {
	case focusFeed:
		f = focusPreview
	case focusItem:
		f = focusFeed
	case focusPreview:
		f = focusItem
	}
	return f
}

func (f focus) next() focus {
	switch f {
	case focusFeed:
		f = focusItem
	case focusItem:
		f = focusPreview
	case focusPreview:
		f = focusFeed
	}
	return f
}

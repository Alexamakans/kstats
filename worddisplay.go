package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const defaultMaxLineLength = 80
const themeSizeName = theme.SizeNameSubHeadingText

type wordDisplay struct {
	first  *widget.RichText
	second *widget.RichText
	third  *widget.RichText
	lines  []*widget.RichText

	maxLineLength int
}

func newWordDisplay() *wordDisplay {
	first := widget.NewRichText()
	second := widget.NewRichText()
	third := widget.NewRichText()
	return &wordDisplay{
		first,
		second,
		third,
		[]*widget.RichText{
			first,
			second,
			third,
		},
		defaultMaxLineLength,
	}
}

func (w *wordDisplay) refresh(rawIndex int) {
	lineIndex, _ := w.rawIndexToLineAndSegmentIndex(rawIndex)
	w.lines[lineIndex].Refresh()
}

func (w *wordDisplay) refreshAll() {
	w.first.Refresh()
	w.second.Refresh()
	w.third.Refresh()
}

func (w *wordDisplay) setIndexToColorName(index int, colorName fyne.ThemeColorName) {
	lineIndex, segmentIndex := w.rawIndexToLineAndSegmentIndex(index)
	w.setColorName(lineIndex, segmentIndex, colorName)
}

func (w *wordDisplay) setIndexToColorNameAndRune(index int, r rune, colorName fyne.ThemeColorName) {
	lineIndex, segmentIndex := w.rawIndexToLineAndSegmentIndex(index)
	w.setColorNameAndRune(lineIndex, segmentIndex, r, colorName)
}

func (w *wordDisplay) setColorNameAndRune(lineIndex, segmentIndex int, r rune, colorName fyne.ThemeColorName) {
	line := w.lines[lineIndex]
	if r == ' ' && colorName == theme.ColorNameError {
		r = '_'
	}
	newSegment := widget.TextSegment{
		Text: string(r),
		Style: widget.RichTextStyle{
			Inline:    true,
			ColorName: colorName,
			Alignment: fyne.TextAlignCenter,
			SizeName:  themeSizeName,
		},
	}

	line.Segments[segmentIndex] = &newSegment
}

func (w *wordDisplay) getRune(rawIndex int) rune {
	lineIndex, segmentIndex := w.rawIndexToLineAndSegmentIndex(rawIndex)
	return rune(w.lines[lineIndex].Segments[segmentIndex].Textual()[0])
}

func (w *wordDisplay) setColorName(lineIndex, segmentIndex int, colorName fyne.ThemeColorName) {
	line := w.lines[lineIndex]
	w.setColorNameAndRune(lineIndex, segmentIndex, rune(line.Segments[segmentIndex].Textual()[0]), colorName)
}

// pushNewLine returns the number of words used from the passed in parameter.
func (w *wordDisplay) pushNewLine(words []string) int {
	w.first.Segments = w.second.Segments
	w.second.Segments = w.third.Segments

	var wordsUsed int
	w.third.Segments, wordsUsed = w.buildSegments(words)
	return wordsUsed
}

// buildSegments returns the built segments and the number of words used from the passed in parameter.
func (w *wordDisplay) buildSegments(words []string) ([]widget.RichTextSegment, int) {
	var segments []widget.RichTextSegment
	wordIndex := 0
	wordsUsed := 0
	extra := 1
	for wordIndex < len(words) && len(segments)+len(words[wordIndex])+extra <= w.maxLineLength {
		wordsUsed++
		for _, char := range words[wordIndex] {
			segments = append(segments, &widget.TextSegment{
				Text: string(char),
				Style: widget.RichTextStyle{
					Inline:    true,
					ColorName: theme.ColorNameForeground,
					Alignment: fyne.TextAlignCenter,
					SizeName:  themeSizeName,
				},
			})
		}
		if extra == 1 {
			segments = append(segments, &widget.TextSegment{
				Text: " ",
				Style: widget.RichTextStyle{
					Inline:    true,
					ColorName: theme.ColorNameForeground,
					Alignment: fyne.TextAlignCenter,
					SizeName:  themeSizeName,
				},
			})
		}
		wordIndex++
		if wordIndex == len(words)-1 {
			extra = 0
		}
	}
	return segments, wordsUsed
}

func (w *wordDisplay) rawIndexToLineAndSegmentIndex(index int) (int, int) {
	firstLen := len(w.first.Segments)
	secondLen := len(w.second.Segments)
	thirdLen := len(w.third.Segments)
	if index >= firstLen && index < firstLen+secondLen {
		return 1, index - firstLen
	} else if index >= firstLen+secondLen && index < firstLen+secondLen+thirdLen {
		return 2, index - firstLen - secondLen
	}
	return 0, index
}

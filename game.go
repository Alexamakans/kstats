package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const defaultRefreshRate = 250

func getWords(path string) ([]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	words := strings.Split(strings.ReplaceAll(string(raw), "\r", ""), "\n")
	return words, nil
}

type game struct {
	words         []string
	rawIndex      int
	wordDisplay   *wordDisplay
	statsDisplay  *widget.Label
	start         time.Time
	refreshRateMs int

	wrongPresses   int
	correctPresses int
	presses        int
	statCollector  statCollector
}

func newGame(path string) (game, error) {
	words, err := getWords(path)
	if err != nil {
		return game{}, err
	}

	g := game{
		words:         words,
		refreshRateMs: defaultRefreshRate,
		wordDisplay:   newWordDisplay(),
		statCollector: newStatCollector(),
		statsDisplay:  widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}),
	}

	g.nextLine()
	g.nextLine()
	g.nextLine()
	g.rawIndex = 0
	g.wordDisplay.refreshAll()

	return g, nil
}

func (g *game) calculateCpm(presses int) float64 {
	return float64(presses) / time.Since(g.start).Minutes()
}

func (g *game) calculateWpm(presses int) float64 {
	return g.calculateCpm(presses) / 5.0
}

func (g *game) onTypedRune(r rune) {
	if g.rawIndex == 0 {
		g.start = time.Now()
	}
	next := g.peekNextRune()
	g.statCollector.collect(r, next)
	if r == next {
		g.setCurrentColorName(theme.ColorNameSuccess)
		g.wordDisplay.refresh(g.rawIndex)
		g.correctPresses++
	} else {
		g.setCurrentColorNameAndRune(r, theme.ColorNameError)
		g.wordDisplay.refresh(g.rawIndex)
		g.wrongPresses++
	}
	g.presses++
	_ = g.nextRune()

	lineIndex, _ := g.wordDisplay.rawIndexToLineAndSegmentIndex(g.rawIndex)
	if lineIndex == 2 {
		g.nextLine()
		g.wordDisplay.refreshAll()
	}

	g.updateStats(g.statsDisplay)
}

func (g *game) setCurrentColorName(colorName fyne.ThemeColorName) {
	g.wordDisplay.setIndexToColorName(g.rawIndex, colorName)
}

func (g *game) setCurrentColorNameAndRune(r rune, colorName fyne.ThemeColorName) {
	g.wordDisplay.setIndexToColorNameAndRune(g.rawIndex, r, colorName)
}

func (g *game) nextLine() {
	if len(g.words) == 0 {
		return
	}

	g.rawIndex -= len(g.wordDisplay.first.String())
	wordsUsed := g.wordDisplay.pushNewLine(g.words)
	g.words = g.words[wordsUsed:]
}

func (g *game) peekNextRune() rune {
	return g.wordDisplay.getRune(g.rawIndex)
}

func (g *game) nextRune() rune {
	// dont worry about it we ballin'
	g.rawIndex++
	return g.wordDisplay.getRune(g.rawIndex - 1)
}

func (g *game) updateStatsLoop(statsDisplay *widget.Label) {
	for {
		g.updateStats(statsDisplay)
		time.Sleep(time.Duration(g.refreshRateMs) * time.Millisecond)
	}
}

func (g *game) updateStats(statsDisplay *widget.Label) {
	s := g.statCollector.calculateStats()
	text := ""
	text += fmt.Sprintf("%.1f wpm\n", g.calculateWpm(g.correctPresses))
	text += fmt.Sprintf("%d errors (%.2f%%)\n", g.wrongPresses, 100*float64(g.correctPresses)/float64(g.presses))

	buf := strings.Builder{}
	w := tabwriter.NewWriter(&buf, 12, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Type\tMin\tMax\tMedian\tMean\tCount")
	fmt.Fprintf(w, "Left -> Right:\t%s\n", s.leftHandToRightHand.String())
	fmt.Fprintf(w, "Right -> Left:\t%s\n", s.rightHandToLeftHand.String())
	fmt.Fprintf(w, "Left -> Right (no space):\t%s\n", s.leftHandToRightHandIgnoreSpace.String())
	fmt.Fprintf(w, "Right -> Left (no space):\t%s\n", s.rightHandToLeftHandIgnoreSpace.String())

	fmt.Fprintf(w, "Lpinky (Same):\t%s\n", s.leftPinkySameFinger.String())
	fmt.Fprintf(w, "Lring  (Same):\t%s\n", s.leftRingSameFinger.String())
	fmt.Fprintf(w, "Lmiddl (Same):\t%s\n", s.leftMiddleSameFinger.String())
	fmt.Fprintf(w, "Lindex (Same):\t%s\n", s.leftIndexSameFinger.String())

	fmt.Fprintf(w, "Rpinky (Same):\t%s\n", s.rightPinkySameFinger.String())
	fmt.Fprintf(w, "Rring  (Same):\t%s\n", s.rightRingSameFinger.String())
	fmt.Fprintf(w, "Rmiddl (Same):\t%s\n", s.rightMiddleSameFinger.String())
	fmt.Fprintf(w, "Rindex (Same):\t%s\n", s.rightIndexSameFinger.String())

	fmt.Fprintf(w, "Lpinky (Diff):\t%s\n", s.leftPinkyDifferentFinger.String())
	fmt.Fprintf(w, "Lring  (Diff):\t%s\n", s.leftRingDifferentFinger.String())
	fmt.Fprintf(w, "Lmiddl (Diff):\t%s\n", s.leftMiddleDifferentFinger.String())
	fmt.Fprintf(w, "Lindex (Diff):\t%s\n", s.leftIndexDifferentFinger.String())

	fmt.Fprintf(w, "Rpinky (Diff):\t%s\n", s.rightPinkyDifferentFinger.String())
	fmt.Fprintf(w, "Rring  (Diff):\t%s\n", s.rightRingDifferentFinger.String())
	fmt.Fprintf(w, "Rmiddl (Diff):\t%s\n", s.rightMiddleDifferentFinger.String())
	fmt.Fprintf(w, "Rindex (Diff):\t%s\n", s.rightIndexDifferentFinger.String())
	w.Flush()
	text += buf.String()
	statsDisplay.SetText(text)
	statsDisplay.Refresh()
}

package main

import (
	"fmt"
	"slices"
	"time"
)

type statCollector struct {
	transitions            []transition
	fingers                []*finger
	previousCollectionTime time.Time
	previousChar           rune
}

func newStatCollector() statCollector {
	return statCollector{
		fingers: []*finger{
			&leftPinky,
			&leftRing,
			&leftMiddle,
			&leftIndex,

			&rightPinky,
			&rightRing,
			&rightMiddle,
			&rightIndex,
			&rightThumb,
		},
	}
}

var validCharacters = []rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l',
	'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
	'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L',
	'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
	'Y', 'Z',
	' ',
}

func (s *statCollector) collect(actual, expected rune) {
	if !slices.Contains(validCharacters, actual) {
		return
	}
	dur := time.Since(s.previousCollectionTime)
	s.previousCollectionTime = time.Now()
	if s.previousChar != 0 {
		t := newTransition(actual, expected, s.previousChar, dur)
		s.transitions = append(s.transitions, t)
	}

	s.previousChar = actual
}

func (s *statCollector) calculateStats() stats {
	st := stats{
		leftHandToRightHand:            newStat(s.getLeftToRightHandTransitions()),
		leftHandToRightHandIgnoreSpace: newStat(s.getLeftToRightHandTransitionsIgnoreSpace()),
		rightHandToLeftHand:            newStat(s.getRightToLeftHandTransitions()),
		rightHandToLeftHandIgnoreSpace: newStat(s.getRightToLeftHandTransitionsIgnoreSpace()),
		leftPinkySameFinger:            newStat(s.getSameFingerTransitions(&leftPinky)),
		leftRingSameFinger:             newStat(s.getSameFingerTransitions(&leftRing)),
		leftMiddleSameFinger:           newStat(s.getSameFingerTransitions(&leftMiddle)),
		leftIndexSameFinger:            newStat(s.getSameFingerTransitions(&leftIndex)),
		leftThumbSameFinger:            newStat(s.getSameFingerTransitions(&leftThumb)),
		rightPinkySameFinger:           newStat(s.getSameFingerTransitions(&rightPinky)),
		rightRingSameFinger:            newStat(s.getSameFingerTransitions(&rightRing)),
		rightMiddleSameFinger:          newStat(s.getSameFingerTransitions(&rightMiddle)),
		rightIndexSameFinger:           newStat(s.getSameFingerTransitions(&rightIndex)),
		rightThumbSameFinger:           newStat(s.getSameFingerTransitions(&rightThumb)),
		leftPinkyDifferentFinger:       newStat(s.getDifferentFingerSameHandTransitions(&leftPinky)),
		leftRingDifferentFinger:        newStat(s.getDifferentFingerSameHandTransitions(&leftRing)),
		leftMiddleDifferentFinger:      newStat(s.getDifferentFingerSameHandTransitions(&leftMiddle)),
		leftIndexDifferentFinger:       newStat(s.getDifferentFingerSameHandTransitions(&leftIndex)),
		leftThumbDifferentFinger:       newStat(s.getDifferentFingerSameHandTransitions(&leftThumb)),
		rightPinkyDifferentFinger:      newStat(s.getDifferentFingerSameHandTransitions(&rightPinky)),
		rightRingDifferentFinger:       newStat(s.getDifferentFingerSameHandTransitions(&rightRing)),
		rightMiddleDifferentFinger:     newStat(s.getDifferentFingerSameHandTransitions(&rightMiddle)),
		rightIndexDifferentFinger:      newStat(s.getDifferentFingerSameHandTransitions(&rightIndex)),
		rightThumbDifferentFinger:      newStat(s.getDifferentFingerSameHandTransitions(&rightThumb)),
	}

	return st
}

func (s *statCollector) getSameFingerTransitions(f *finger) []transition {
	var transitions []transition
	for _, t := range s.transitions {
		if t.correct && s.getFinger(t.from) == f && s.getFinger(t.toActual) == f {
			transitions = append(transitions, t)
		}
	}

	return transitions
}

func (s *statCollector) getDifferentFingerSameHandTransitions(f *finger) []transition {
	var transitions []transition
	for _, t := range s.transitions {
		from := s.getFinger(t.from)
		to := s.getFinger(t.toActual)
		if t.correct && from.left != to.left && from != to && to == f {
			transitions = append(transitions, t)
		}
	}

	return transitions
}

func (s *statCollector) getLeftToRightHandTransitions() []transition {
	var transitions []transition
	for _, t := range s.transitions {
		if t.correct && s.getFinger(t.from).left && !s.getFinger(t.toActual).left {
			transitions = append(transitions, t)
		}
	}

	return transitions
}

func (s *statCollector) getLeftToRightHandTransitionsIgnoreSpace() []transition {
	var transitions []transition
	for _, t := range s.transitions {
		if t.from == ' ' || t.toActual == ' ' {
			continue
		}
		if t.correct && s.getFinger(t.from).left && !s.getFinger(t.toActual).left {
			transitions = append(transitions, t)
		}
	}

	return transitions
}

func (s *statCollector) getRightToLeftHandTransitions() []transition {
	var transitions []transition
	for _, t := range s.transitions {
		if t.correct && !s.getFinger(t.from).left && s.getFinger(t.toActual).left {
			transitions = append(transitions, t)
		}
	}

	return transitions
}

func (s *statCollector) getRightToLeftHandTransitionsIgnoreSpace() []transition {
	var transitions []transition
	for _, t := range s.transitions {
		if t.from == ' ' || t.toActual == ' ' {
			continue
		}
		if t.correct && !s.getFinger(t.from).left && s.getFinger(t.toActual).left {
			transitions = append(transitions, t)
		}
	}

	return transitions
}

func (s *statCollector) getFinger(r rune) *finger {
	for _, f := range s.fingers {
		if slices.Contains(f.chars, r) {
			return f
		}
	}

	panic(fmt.Sprintf("'%s' (%v) is not assigned to a finger\n", string(r), r))
}

var (
	leftPinky  = newFinger(true, 'q', 'a', 'z', 'Q', 'A', 'Z')
	leftRing   = newFinger(true, 'w', 's', 'x', 'W', 'S', 'X')
	leftMiddle = newFinger(true, 'e', 'd', 'c', 'E', 'D', 'C')
	leftIndex  = newFinger(true, 'r', 'f', 'v', 't', 'g', 'b', 'R', 'F', 'V', 'T', 'G', 'B')
	leftThumb  = newFinger(true)

	rightPinky  = newFinger(false, 'p', 'P')
	rightRing   = newFinger(false, 'o', 'l', 'O', 'L')
	rightMiddle = newFinger(false, 'i', 'k', 'I', 'K')
	rightIndex  = newFinger(false, 'u', 'j', 'm', 'U', 'J', 'M', 'y', 'h', 'n', 'Y', 'H', 'N')
	rightThumb  = newFinger(false, ' ')
)

type finger struct {
	left  bool
	chars []rune
}

func newFinger(left bool, keys ...rune) finger {
	return finger{
		left:  left,
		chars: keys,
	}
}

type transition struct {
	from       rune
	toExpected rune
	toActual   rune
	duration   time.Duration
	correct    bool
}

func newTransition(actual, expected, prev rune, duration time.Duration) transition {
	return transition{
		from:       prev,
		toExpected: expected,
		toActual:   actual,
		duration:   duration,
		correct:    actual == expected,
	}
}

type stats struct {
	leftHandToRightHand            stat
	leftHandToRightHandIgnoreSpace stat
	rightHandToLeftHand            stat
	rightHandToLeftHandIgnoreSpace stat

	leftPinkySameFinger  stat
	leftRingSameFinger   stat
	leftMiddleSameFinger stat
	leftIndexSameFinger  stat
	leftThumbSameFinger  stat

	rightPinkySameFinger  stat
	rightRingSameFinger   stat
	rightMiddleSameFinger stat
	rightIndexSameFinger  stat
	rightThumbSameFinger  stat

	leftPinkyDifferentFinger  stat
	leftRingDifferentFinger   stat
	leftMiddleDifferentFinger stat
	leftIndexDifferentFinger  stat
	leftThumbDifferentFinger  stat

	rightPinkyDifferentFinger  stat
	rightRingDifferentFinger   stat
	rightMiddleDifferentFinger stat
	rightIndexDifferentFinger  stat
	rightThumbDifferentFinger  stat
}

type stat struct {
	durationMin    time.Duration
	durationMax    time.Duration
	durationMedian time.Duration
	durationMean   time.Duration

	minT transition
	maxT transition

	count int
}

func newStat(transitions []transition) stat {
	if len(transitions) == 0 {
		return stat{}
	}
	minT := slices.MinFunc(transitions, func(a, b transition) int {
		return int(a.duration) - int(b.duration)
	})
	minDuration := minT.duration
	maxT := slices.MaxFunc(transitions, func(a, b transition) int {
		return int(a.duration) - int(b.duration)
	})
	maxDuration := maxT.duration
	var sum time.Duration
	for _, t := range transitions {
		sum += t.duration
	}
	mean := time.Duration(float64(sum) / float64(len(transitions)))

	slices.SortFunc(transitions, func(a, b transition) int {
		if a.duration < b.duration {
			return -1
		}
		return 1
	})

	var median time.Duration
	if len(transitions) > 0 {
		median = transitions[len(transitions)/2].duration
	}

	return stat{
		durationMin:    minDuration,
		durationMax:    maxDuration,
		durationMedian: median,
		durationMean:   mean,
		minT:           minT,
		maxT:           maxT,
		count:          len(transitions),
	}
}

func (s *stat) String() string {
	minS := fmt.Sprintf("%d ms", s.durationMin.Milliseconds())
	if s.minT.from != 0 {
		minS = fmt.Sprintf("%v ms (%s)",
			s.durationMin.Milliseconds(),
			fmt.Sprintf("%s -> %s", string(s.printableChar(s.minT.from)), string(s.printableChar(s.minT.toActual))))
	}
	maxS := fmt.Sprintf("%d ms", s.durationMax.Milliseconds())
	if s.maxT.from != 0 {
		maxS = fmt.Sprintf("%v ms (%s)",
			s.durationMax.Milliseconds(),
			fmt.Sprintf("%s -> %s", string(s.printableChar(s.maxT.from)), string(s.printableChar(s.maxT.toActual))))
	}
	return fmt.Sprintf("%s\t%s\t%d ms\t%d ms\t%d",
		minS,
		maxS,
		s.durationMedian.Milliseconds(),
		s.durationMean.Milliseconds(),
		s.count,
	)
}

func (s *stat) printableChar(r rune) rune {
	if r == ' ' {
		return '_'
	}

	return r
}

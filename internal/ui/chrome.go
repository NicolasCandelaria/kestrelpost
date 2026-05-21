package ui

import (
	"fmt"
	"strings"

	lg "charm.land/lipgloss/v2"
	"kestrelpost/internal/game"
)

// terminal.shop–inspired chrome: grid header, framed body, shortcut footer.
// Aesthetic reference: https://github.com/IsaiahPapa/terminal.shop

const minFrameWidth = 56

var (
	colFG       = lg.Color("#d4d4d4")
	colMuted    = lg.Color("#737373")
	colLine     = lg.Color("#404040")
	colActiveBG = lg.Color("#262626")
	colActiveFG = lg.Color("#fafafa")

	root = lg.NewStyle().Background(lg.Color("#000000")).Foreground(colFG)

	tabLo = root.Copy().
		Padding(0, 1).
		BorderStyle(lg.NormalBorder()).
		BorderLeft(true).
		BorderForeground(colLine)

	tabHi = root.Copy().
		Padding(0, 1).
		BorderStyle(lg.NormalBorder()).
		BorderLeft(true).
		BorderForeground(colLine).
		Background(colActiveBG).
		Foreground(colActiveFG).
		Bold(true)

	keyHi = lg.NewStyle().Background(colActiveBG).Foreground(colActiveFG).Bold(true)
	keyLo = root.Copy().Foreground(colActiveFG).Bold(true)

	frame = root.Copy().
		BorderStyle(lg.NormalBorder()).
		BorderForeground(colLine).
		BorderTop(true).BorderLeft(true).BorderRight(true).BorderBottom(true)

	inner = root.Copy().Padding(1, 2)

	muted = root.Copy().Foreground(colMuted)
)

func effectiveWidth(w int) int {
	if w < minFrameWidth {
		return minFrameWidth
	}
	return w
}

func headerBar(w int, phase game.Phase, night, fuel int) string {
	cellW := w / 4
	if cellW < 14 {
		cellW = 14
	}

	relayOn := phase == game.PhaseIntro || phase == game.PhaseNight

	cell := func(on bool, keyLetter, title, value string) string {
		st := tabLo
		if on {
			st = tabHi
		}
		keyPart := keyLo.Render(keyLetter)
		if on {
			keyPart = keyHi.Render(keyLetter)
		}
		rest := " " + title
		if value != "" {
			rest += " " + value
		}
		var restStyled string
		if on {
			restStyled = tabHi.Copy().Bold(false).Render(rest)
		} else {
			restStyled = tabLo.Copy().Bold(false).Render(rest)
		}
		return st.Width(cellW).Render(keyPart + restStyled)
	}

	c1 := cell(relayOn, "r", "radio", "")
	c2 := cell(phase == game.PhaseNight, "s", "scan", "")
	c3 := cell(false, "n", "night", fmt.Sprintf("%d", night))
	c4 := cell(false, "g", "power", gaugeWord(fuel))

	return root.Width(w).Render(lg.JoinHorizontal(lg.Top, c1, c2, c3, c4))
}

func promoLine(w int) string {
	s := "Kestrel Post · Northern Manitoba · Operator Seven"
	return muted.Width(w).Align(lg.Center).Render(s)
}

func footerBar(w int, phase game.Phase) string {
	var left, right string
	switch phase {
	case game.PhaseIntro:
		left = keyLo.Render("t") + root.Render(" tutorial  ·  ") + keyLo.Render("enter") + root.Render(" start")
		right = keyLo.Render("q") + root.Render(" quit")
	case game.PhaseTutorial:
		left = keyLo.Render("enter") + root.Render(" start night 1")
		right = keyLo.Render("q") + root.Render(" quit")
	case game.PhaseNight:
		left = keyLo.Render("1/2/3") + root.Render(" choose  ·  ") + keyLo.Render("s") + root.Render(" scan  ·  ") + keyLo.Render("f") + root.Render(" feed dog")
		right = keyLo.Render("q") + root.Render(" quit")
	default:
		left = muted.Render("run ended")
		right = keyLo.Render("q") + root.Render(" quit")
	}
	lw, rw := lg.Width(left), lg.Width(right)
	pad := w - lw - rw
	if pad < 2 {
		pad = 2
	}
	return root.Width(w).Render(left + strings.Repeat(" ", pad) + right)
}

func gaugeWord(fuel int) string {
	switch {
	case fuel <= 0:
		return "empty"
	case fuel <= 20:
		return "critical"
	case fuel <= 45:
		return "low"
	case fuel <= 75:
		return "steady"
	default:
		return "full"
	}
}

func shopFrame(width int, phase game.Phase, night, fuel int, body string) string {
	w := effectiveWidth(width)
	head := headerBar(w, phase, night, fuel)

	innerW := w - 4
	if innerW < 24 {
		innerW = 24
	}
	bodyBlock := inner.Width(innerW).Align(lg.Left).Render(body)
	mid := frame.Width(w).Render(bodyBlock)

	rule := lg.NewStyle().Foreground(colLine).Width(w).Render(strings.Repeat("─", w))
	foot := footerBar(w, phase)
	bottom := root.Width(w).Render(promoLine(w) + "\n" + rule + "\n" + foot)

	return root.Width(w).Render(head + "\n" + mid + "\n" + bottom)
}

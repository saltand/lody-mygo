package app

// The controls several surfaces share: buttons, a checkbox, a progress
// bar and the header's status pill, as @lody/ui draws them.

import (
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/transfer"
	"github.com/egoist/mygo/ui"

	"github.com/lbr77/lody-mygo/locales"
)

// progressBar is Lody's: a bar filled to the percent, or, without one, a
// 40% band crossing the track every 1.4s.
func progressBar(c *ui.Context, p *palette, determinate bool, percent int) {
	track := ui.Row(c.Key("progress")).FillWidth().Height(6).Radius(3).Background(p.Text.Alpha(0.06)).ClipX()
	track.Children(func() {
		if determinate {
			w := track.AnimateWith("fill", float32(percent), durRegular, easeLody)
			ui.Box(c).WidthPercent(w).Height(6).Radius(3).Background(p.Primary)
			return
		}
		band := ui.Box(c.Key("band")).WidthPercent(40).Height(6).Radius(3).Background(p.Primary)
		k := band.Keyframes("sweep", []ui.Keyframe{{Value: 0}, {At: 1400 * time.Millisecond, Value: 1, Ease: easeLody}}, ui.KeyframeOptions{Repeat: ui.RepeatForever})
		w := track.Bounds().W * 0.4
		band.Translate(w*(-1+3.5*k), 0)
	})
}

// lodyCheckbox is @lody/ui's check box: a 16px well, the accent with an
// ink mark once checked, or a dash when some of what it stands for are,
// fading between them over 120ms.
func lodyCheckbox(c *ui.Context, p *palette, key string, on, some bool) ui.Element {
	b := ui.Box(c.Key(key)).Size(16, 16).Radius(5).Shrink(0).Center().Transition(colorFade)
	if on || some {
		top := p.Primary.Mix(ui.Hex("#FFFFFF"), 0.1)
		bottom := p.Primary
		if p.Dark {
			top, bottom = p.Primary, p.Primary.Mix(ui.Hex("#000000"), 0.08)
		}
		b.Gradient(top, bottom, 180)
		if p.Dark {
			b.InsetShadow(0, 1, 0, 0, ui.RGBA(255, 255, 255, 0.55)).Shadow(0, 1, 2, 0, ui.RGBA(0, 0, 0, 0.5))
		} else {
			b.InsetShadow(0, 1, 0, 0, ui.RGBA(255, 255, 255, 0.18)).Shadow(0, 1, 1.5, 0, ui.RGBA(28, 28, 28, 0.18))
		}
	} else {
		b.Background(wellGround(p))
		if p.Dark {
			b.InsetShadow(0, 1, 2, 0, ui.RGBA(0, 0, 0, 0.4)).InsetShadow(0, 0, 0, 1, ui.RGBA(255, 255, 255, 0.06)).Shadow(0, 1, 0, 0, ui.RGBA(255, 255, 255, 0.04))
		} else {
			b.InsetShadow(0, 1, 1.5, 0, ui.RGBA(28, 28, 28, 0.08)).InsetShadow(0, 0, 0, 1, ui.RGBA(28, 28, 28, 0.06)).Shadow(0, 1, 0, 0, ui.RGBA(255, 255, 255, 0.8))
		}
	}
	b.Children(func() {
		switch {
		case on:
			ui.Icon(c, shareCheckMark).Size(10, 10).TextColor(p.Canvas)
		case some:
			ui.Box(c).Size(8, 1.6).Radius(1).Background(p.Canvas)
		}
	})
	return b
}

var shareCheckMark = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10" fill="none"><path d="M1.6 5.2 3.9 7.5 8.4 2.7" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/></svg>`))

// lodyButton is @lody/ui's Button at the small step: primary in ink,
// secondary raised, ghost bare; a busy one shows its spinner in place of
// its icon.
func lodyButton(c *ui.Context, p *palette, variant, key string, icon *ui.SVG, label string, enabled bool, busy ...bool) ui.Element {
	b, ink := lodyButtonLook(c, p, variant, key, label, enabled)
	b.Children(func() {
		switch {
		case len(busy) > 0 && busy[0]:
			spinner(c).Size(16, 16)
		case icon != nil:
			ui.Icon(c, icon).Size(16, 16).TextColor(ink)
		}
		if label != "" {
			ui.Text(c, label).FontSize(sizeSub).FontWeight(500).LetterSpacing(-0.13).TextColor(ink).SingleLine()
		}
	})
	return b
}

// lodyButtonLook is lodyButton without its content: the button, and the
// ink what it holds takes.
func lodyButtonLook(c *ui.Context, p *palette, variant, key, label string, enabled bool) (ui.Element, ui.Color) {
	b := ui.ButtonBase(c.Key("btn/" + key)).Height(28).Radius(8).Gap(6).AlignItems(ui.Center).Justify(ui.Center).Shrink(0).
		Label(label).Disabled(!enabled).Transition(colorFade)
	if label == "" {
		b.Width(28)
	} else {
		b.Padding(0, 10)
	}
	hovered := b.Hovered() && enabled
	pressed := b.Pressed() && enabled
	// Its fill, ink, shadow, press and opacity all move over 120ms on
	// ease.standard, as @lody/ui's Button transitions them.
	motion := ui.Animation{Duration: durFast, Ease: easeStandard}
	ink := p.Text
	switch variant {
	case "primary":
		fill := p.Text
		if hovered {
			fill = fill.Mix(p.Canvas, 0.12)
		}
		fill = b.AnimateColor("fill", fill, motion)
		ink = p.Canvas
		switch {
		case pressed:
			b.Background(fill)
		case p.Dark:
			b.Gradient(fill, fill.Mix(ui.Hex("#000000"), 0.08), 180).InsetShadow(0, 1, 0, 0, ui.RGBA(255, 255, 255, 0.55)).Shadow(0, 1, 2, 0, ui.RGBA(0, 0, 0, 0.5))
		default:
			b.Gradient(fill.Mix(ui.Hex("#FFFFFF"), 0.1), fill, 180).InsetShadow(0, 1, 0, 0, ui.RGBA(255, 255, 255, 0.18)).Shadow(0, 1, 1.5, 0, ui.RGBA(28, 28, 28, 0.18))
		}
	case "secondary":
		fill := raisedGround(p)
		if p.Dark {
			fill = p.Popover
		}
		if hovered {
			fill = fill.Mix(p.Text, 0.04)
		}
		fill = b.AnimateColor("fill", fill, motion)
		if pressed {
			b.Background(fill).Shadow(0, 0, 0, 0.5, p.Text.Alpha(0.14)).InsetShadow(0, 1, 1, 0, p.Text.Alpha(0.14))
		} else {
			raisedLook(b, p, fill)
		}
	default:
		ink = p.Secondary
		if hovered {
			b.Background(p.Text.Alpha(0.06))
			ink = p.Text
		}
	}
	ink = b.AnimateColor("ink", ink, motion)
	b.Opacity(b.AnimateWith("dim", map[bool]float32{true: 1, false: 0.45}[enabled], durFast, easeStandard))
	if enabled {
		b.Cursor(ui.CursorPointer)
		// Only the filled variants sink as they are pressed; a ghost stays.
		if variant != "ghost" {
			pressable(b, 0.97)
		}
	}
	return b, ink
}

// raisedLook stands a fill on the surface: the raised sheen over it, and
// the raised shadow under it.
func raisedLook(e ui.Element, p *palette, fill ui.Color) { raisedLookK(e, p, fill, 1) }

// raisedLookK is raisedLook k of the way in: its fill, sheen and shadow
// at k of their strength, for one fading in or out.
func raisedLookK(e ui.Element, p *palette, fill ui.Color, k float32) {
	if p.Dark {
		e.Gradient(fill.Mix(ui.Hex("#FFFFFF"), 0.04).Alpha(k), fill.Alpha(k), 180).InsetShadow(0, 1, 0, 0, ui.RGBA(255, 255, 255, 0.08*k)).
			Shadow(0, 0, 0, 0.5, ui.RGBA(0, 0, 0, 0.7*k)).Shadow(0, 1, 2, 0, ui.RGBA(0, 0, 0, 0.5*k))
		return
	}
	e.Gradient(fill.Alpha(k), fill.Mix(ui.Hex("#1C1C1C"), 0.03).Alpha(k), 180).Shadow(0, 0, 0, 0.5, ui.RGBA(28, 28, 28, 0.16*k)).
		Shadow(0, 1, 1, 0, ui.RGBA(28, 28, 28, 0.06*k)).Shadow(0, 2, 4, -1, ui.RGBA(28, 28, 28, 0.07*k))
}

// headerPill is a status pill of the conversation header, as Lody's.
func headerPill(c *ui.Context, p *palette, key string, icon *ui.SVG, label string, open bool) ui.Element {
	// Tailwind's transition-colors: the fill and the ink fade over 150ms.
	b := ui.ButtonBase(c.Key(key)).Height(28).Padding(0, 8).Gap(6).Radius(6).AlignItems(ui.Center).Shrink(0).Transition(twFade)
	ink := p.Text.Alpha(0.8)
	switch {
	case open:
		b.Background(p.Text.Alpha(0.07))
		ink = p.Text
	case b.Hovered():
		b.Background(p.Text.Alpha(0.05))
		ink = p.Text
	}
	ink = b.AnimateColor("ink", ink, ui.Animation{Duration: 150 * time.Millisecond, Ease: twEase})
	b.Children(func() {
		ui.Icon(c, icon).Size(14, 14).TextColor(ink)
		ui.Text(c, label).FontSize(11.2).FontWeight(500).TextColor(ink).SingleLine()
	})
	return b
}

func abs32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

// archivedBadge marks an archived conversation in the header.
func archivedBadge(c *ui.Context, p *palette) {
	b := headerPill(c, p, "archived", iconArchive, locales.AccessArchived(), false).Label(locales.AccessArchived() + ": " + locales.AccessArchivedDescription())
	ui.TooltipBase(c, b, func(tip ui.Element) {
		tip.Padding(8, 10).Radius(8).MaxWidth(288).Background(p.Popover).Border(0.5, p.CardRim).Gap(2).
			Shadow(0, 6, 16, -4, ui.RGBA(0, 0, 0, 0.18)).
			Transition(ui.ElementTransition{Duration: durFast, Ease: easeStandard, Enter: &ui.Motion{Y: 4}, Exit: &ui.Motion{Y: 4}})
		tip.Children(func() {
			ui.Text(c, locales.AccessArchived()).FontSize(sizeFoot).FontWeight(500).TextColor(p.Text)
			ui.Text(c, locales.AccessArchivedDescription()).FontSize(sizeFoot).TextColor(p.Secondary)
		})
	})
}

// clipImageClipboard puts PNG bytes on the clipboard; tests replace it.
var clipImageClipboard = func(data []byte) error {
	return mygo.Clipboard.Write(transfer.New(transfer.NewItem(transfer.Bytes(transfer.PNG, data))))
}

// saveFileDialog asks where to save, "" when the user backed out; tests
// replace it.
var saveFileDialog = func(win *mygo.Window, name string) (string, error) {
	return mygo.Dialog.Save(mygo.SaveDialogOptions{Parent: win, DefaultPath: name})
}

// note is a small line of prose, as the panels foot with one.
func note(c *ui.Context, p *palette, text string, ink ui.Color) {
	ui.Text(c, text).FontSize(sizeFoot).LineHeight(1.45).TextColor(ink).Selectable()
}

// groupThousands writes a count with its commas: 12,345.
func groupThousands(n int64) string {
	s := fmt.Sprint(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

// stripRecommended is a model's name without Lody's "(recommended)".
var stripRecommended = regexp.MustCompile(`(?i)\s*\(recommended\)`)
// webBilling opens the workspace's billing on lody.ai.
func (a *App) webBilling() {
	slug := a.workspaceID
	if ws := a.workspace(); ws != nil && ws.Slug != "" {
		slug = ws.Slug
	}
	mygo.Shell.OpenExternal("https://lody.ai/" + url.PathEscape(slug) + "/settings/billing")
}

// webLink is a link-styled button that opens lody.ai.
func webLink(c *ui.Context, p *palette, label string, open func()) {
	b := ui.ButtonBase(c.Key(label)).Height(24).Padding(0, 4).Radius(6).Gap(4).AlignItems(ui.Center).Cursor(ui.CursorPointer).Label(label).OnClick(open).Transition(colorFade)
	ink := p.Primary
	if b.Hovered() {
		b.Background(p.Hover)
	}
	b.Children(func() {
		ui.Text(c, label).FontSize(sizeFoot).TextColor(ink)
		ui.Icon(c, iconExternal).Size(12, 12).TextColor(ink)
	})
}


// formatDay writes a date as the app's language does: 2026年11月27日 or
// Nov 27, 2026.
func formatDay(ms float64) string {
	if ms <= 0 {
		return ""
	}
	t := time.UnixMilli(int64(ms))
	if zhLocale() {
		return t.Format("2006年1月2日")
	}
	return t.Format("Jan 2, 2006")
}

// billingOverview is the workspace's plan as the checkout check reads it.
type billingOverview struct {
	BillingAccountID  *string  `json:"billingAccountId"`
	PlanTier          string   `json:"planTier"`
	EffectivePlanTier string   `json:"effectivePlanTier"`
	EntitlementSource string   `json:"entitlementSource"`
	BillingInterval   *string  `json:"billingInterval"`
	CancelAtPeriodEnd bool     `json:"cancelAtPeriodEnd"`
	CurrentPeriodEnd  *float64 `json:"currentPeriodEnd"`
	NextBillingAt     *float64 `json:"nextBillingAt"`
	GiftEndsAt        *float64 `json:"giftEndsAt"`
	PromotionEndsAt   *float64 `json:"promotionalEntitlementEndsAt"`
	SeatCount         float64  `json:"seatCount"`
	CanManageBilling  bool     `json:"canManageBilling"`
	ReferralRewards   bool     `json:"referralRewardsSupported"`
}

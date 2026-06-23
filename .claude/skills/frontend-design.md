---
name: frontend-design
description: Guidance for building distinctive, accessible web interfaces. Covers design tokens, typography, component design, accessibility, and visual identity. Synthesized from Material Design 3, IBM Carbon, Shopify Polaris, WCAG 2.2, and Practical Typography.
---

# Frontend Design

**Trigger:** User asks to build, design, review, or refactor a web UI — or any frontend component, layout, or visual system.

Apply automatically. Work top-down: token system first, then typography, then components, then accessibility pass.

---

## 1. Design Tokens — The Foundation

Token architecture from Material Design 3, IBM Carbon, and Shopify Polaris independently converges on the same three-tier model. Use it.

### Three tiers

```
reference   →  concrete values: hex colors, px sizes, font family names
system      →  design-character decisions: which reference values mean "primary", "surface", "on-surface"
component   →  per-element attributes: what color a Button's container uses at rest vs. pressed
```

Reference tokens are never used in component code directly. System tokens are the design language. Component tokens reference system tokens.

### CSS implementation

Implement as CSS custom properties. Scope at any selector level — never require global overrides.

```css
/* reference */
:root {
  --color-blue-500: #1a73e8;
  --color-neutral-900: #1c1c1c;
  --size-4: 4px;
}

/* system */
:root {
  --color-primary: var(--color-blue-500);
  --color-on-surface: var(--color-neutral-900);
  --space-xs: var(--size-4);
}

/* component — scoped, not global */
.btn-primary {
  --btn-bg: var(--color-primary);
  --btn-text: var(--color-on-surface-inverse);
  background: var(--btn-bg);
  color: var(--btn-text);
}

/* per-subtree theme override — no !important, no global change */
.dark-panel {
  --color-primary: var(--color-blue-300);
}
```

### Semantic naming rule

Name tokens by **function**, not by visual appearance. This is Polaris's model and it's correct.

```
✓  --color-bg-surface-hover      (what it does, where, in what state)
✓  --color-text-link-active
✓  --color-border-critical

✗  --color-blue-hover             (visual, not semantic)
✗  --blue                         (raw, unusable across themes)
```

Categories: `bg` (backgrounds), `text`, `border`, `icon`, `fill` (interactive/filled elements). Each gets its own semantic namespace.

---

## 2. Typography — Deliberate Constraints

Typography carries personality. The defaults are wrong for most designs.

### Line length

**45–90 characters per line.** This is a layout decision — the browser default lets text span the full viewport, which is always too wide.

```css
/* enforce at container level */
.content {
  max-width: 65ch;   /* ~66 chars — the typographic sweet spot */
}
```

Never skip this. A design without line-length control looks unfinished regardless of type choices.

### Type scale

Set a deliberate scale with a clear ratio (1.25 Major Third, 1.333 Perfect Fourth, 1.5 Major Sixth). Do not set sizes ad hoc.

```css
:root {
  --text-xs:   0.75rem;   /* 12px */
  --text-sm:   0.875rem;  /* 14px */
  --text-base: 1rem;      /* 16px — browser default, do not fight it */
  --text-lg:   1.25rem;   /* 20px */
  --text-xl:   1.5rem;    /* 24px */
  --text-2xl:  2rem;      /* 32px */
  --text-3xl:  3rem;      /* 48px */
}
```

### Typeface roles

Assign faces to roles, not to elements. Minimum two roles:

| Role | Use | Character |
|------|-----|-----------|
| Display | Headlines, hero text | Characterful — this is where personality lives |
| Body | Prose, UI labels, microcopy | Legible at small sizes, neutral enough not to compete |
| Utility (optional) | Captions, data, timestamps | Monospace or tabular figures if numbers need to align |

Do not reach for the same display face every time. The subject's world determines the face. A logistics tool and a creative portfolio have different display personalities.

### Pairing rule

Display + Body must be distinct enough to justify using two faces. If they look like the same family at a glance, replace one. Contrast in weight, width, or construction is what makes a pairing work.

---

## 3. Accessibility — Non-Negotiable Floor

WCAG 2.2 POUR framework is the normative standard. AA conformance is the minimum for production UI.

### POUR principles

Every accessibility decision maps to one of four:
- **Perceivable** — users can receive the content (contrast, alt text, captions)
- **Operable** — users can navigate and interact (keyboard, focus, no traps)
- **Understandable** — content and UI behavior are clear (labels, errors, language)
- **Robust** — works across assistive technologies (valid semantics, ARIA)

### Contrast — hard numbers

These are pass/fail, not guidelines:

| Text type | Minimum ratio | WCAG criterion |
|-----------|--------------|----------------|
| Normal text (< 18pt / < 14pt bold) | **4.5 : 1** | SC 1.4.3 Level AA |
| Large text (≥ 18pt / ≥ 14pt bold) | **3 : 1** | SC 1.4.3 Level AA |
| UI components, focus indicators | **3 : 1** | SC 1.4.11 Level AA |

Check every color combination before shipping. Use a contrast checker — do not eyeball it.

Note: WCAG 3.0 (in development) proposes the APCA algorithm to replace these ratios. Do not adopt APCA yet — it is not normative. Apply the 2.x ratios now.

### Accessible names

Every interactive element must have an accessible name. This is "among the most important responsibilities" per the WAI-ARIA Authoring Practices Guide.

```html
<!-- visible label — preferred -->
<button>Save changes</button>

<!-- icon-only — needs aria-label -->
<button aria-label="Close dialog">✕</button>

<!-- input — use <label>, not placeholder -->
<label for="email">Email address</label>
<input id="email" type="email">
```

Accessible descriptions (`aria-describedby`) add supplementary context — error messages, help text. They are secondary to names, not a substitute.

### Keyboard interaction model

Two distinct navigation modes — apply the right one per component:

| Action | Key | When |
|--------|-----|------|
| Move between widgets | `Tab` / `Shift+Tab` | Always — between interactive components |
| Navigate within a widget | Arrow keys | Inside: menus, listboxes, tabs, sliders |
| Exit a widget | `Tab` | Moves focus out of the component entirely |

```
Tab  →  [Button]  →  [Input]  →  [Menu▾]  →  [Button]
                                    ↕ ↕
                               Arrow keys navigate menu items
```

This model is not optional for ARIA widgets. Getting it wrong produces components that are keyboard-accessible in appearance but broken in practice.

### Semantic HTML first

Prefer native elements over ARIA. `<button>` over `<div role="button">`. `<nav>` over `<div role="navigation">`. Native elements carry role, state, and keyboard behavior for free. ARIA adds what HTML cannot express — it does not replace HTML.

---

## 4. Component Design Principles

From IBM Carbon's model: patterns are bottom-up, not top-down. Harvest from real UI problems, not invented scenarios.

### Anatomy rule

Every component has a clear anatomy: container, content slots, interactive states. Define all states before writing code:

```
default → hover → focus → active → disabled → error (where applicable)
```

Missing a state is a bug, not an aesthetic choice.

### State via tokens, not hardcoded values

```css
/* wrong — hardcoded, no theming path */
.btn:hover { background: #1557b0; }

/* right — state expressed through token override */
.btn:hover { --btn-bg: var(--color-primary-hover); }
```

### Composition over configuration

A component that accepts 12 boolean props is a smell. Split into composable primitives. Prefer slot-based composition (children, named slots) over prop-driven variants beyond a small set.

---

## 5. Visual Identity

**Note:** Visual identity principles (logo clearspace, brand color hierarchy, scale ratios) were not confirmed by the research pass. The following is editorial guidance based on design system patterns, not verified primary source claims.

### The signature element rule

One distinctive element per page or component. Let it be the thing the design is remembered by. Keep everything around it quiet. Decoration that does not serve the brief is subtraction from budget, not addition of value.

### Palette construction

Start with purpose, not with color preference:

1. **Background** — 1 value, near-neutral (light or dark)
2. **Primary** — 1 value, highest brand energy, used sparingly
3. **Surface** — 1–2 values for cards, modals, elevated elements
4. **On-[color]** — for every background, define its foreground text color explicitly
5. **Accent** — 1 value, used for emphasis only (not decorative fill)
6. **Critical / Warning / Success** — semantic feedback colors, separate from brand palette

Six named hex values is enough for most UIs. More is not more expressive — it is more inconsistency risk.

### What makes a palette distinctive

Not the hues — the relationships. The ratio of primary to neutral, the temperature (warm neutral vs. cool neutral), the luminance range between surface and background. Two palettes can share a blue and read completely differently.

---

## 6. Motion

**Note:** The research pass produced no surviving verified claims on motion principles. The web.dev prefers-reduced-motion article and Josh Comeau's CSS transition guides were fetched but yielded no adversarially-confirmed claims. Apply motion conservatively and verify your approach against current sources.

Minimum known-safe practice:

```css
/* respect user system preferences */
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

## Review Checklist

Run before shipping any frontend component:

```
Tokens
[ ] Three-tier token hierarchy in use (reference / system / component)
[ ] All color tokens named by function, not by visual appearance
[ ] Component states expressed via token overrides, not hardcoded values

Typography
[ ] Body text container has max-width (target ~65ch)
[ ] Type scale uses a consistent ratio — no ad hoc sizes
[ ] Display and body faces are meaningfully distinct

Accessibility
[ ] All interactive elements have accessible names
[ ] Contrast: 4.5:1 normal text, 3:1 large text, 3:1 UI components
[ ] Keyboard: Tab moves between widgets, arrow keys within widgets
[ ] Semantic HTML used wherever possible; ARIA only where HTML is insufficient

Components
[ ] All states defined (default, hover, focus, active, disabled, error)
[ ] No hardcoded color or spacing values — all through tokens
[ ] Composition preferred over prop proliferation
```

---

## Sources (verified claims only)

| Source | Confirmed claims |
|--------|-----------------|
| [Material Design 3 — Tokens](https://m3.material.io/foundations/design-tokens/overview) | 3-tier token hierarchy, CSS scoping |
| [IBM Carbon — Get Started](https://carbondesignsystem.com/designing/get-started/) | Token pattern corroboration, bottom-up component model |
| [Shopify Polaris — Color Tokens](https://polaris.shopify.com/tokens/color) | Semantic function-based naming |
| [Practical Typography — Line Length](https://practicaltypography.com/line-length.html) | 45–90 char constraint |
| [WCAG 2.2](https://www.w3.org/TR/WCAG22/) | POUR / A–AA–AAA structure |
| [WCAG SC 1.4.3](https://www.w3.org/WAI/WCAG21/Understanding/contrast-minimum.html) | 4.5:1 / 3:1 contrast ratios |
| [WAI-ARIA APG — Names and Descriptions](https://www.w3.org/WAI/ARIA/apg/practices/names-and-descriptions/) | Accessible name responsibility |
| [Inclusive Components](https://inclusive-components.design/) | Tab/arrow keyboard model |

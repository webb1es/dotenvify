# DotEnvify Landing Page — AI Agent Prompt

## Objective

Build a single-viewport marketing landing page for **DotEnvify**, a JetBrains IDE plugin that pulls environment
variables from Azure DevOps Variable Groups and converts raw key-value data into properly formatted `.env` files — all
without leaving the IDE.

**Critical constraint:** On desktop (1440px+), all content must fit in a single viewport — no scrolling. On mobile,
vertical scroll is acceptable.

**Target URL:** dotenvify.webbies.dev (or similar)

---

## About the Product

DotEnvify is a native JetBrains plugin (Kotlin) that works across IntelliJ IDEA, GoLand, WebStorm, PyCharm, Rider,
CLion, and RubyMine. It eliminates the manual copy-paste workflow of managing `.env` files by integrating directly with
Azure DevOps and providing smart formatting tools inside the IDE.

**Plugin ID:** `dev.webbies.dotenvify`
**Version:** 0.1.0
**Author:** Webster Muchefa / Webbies (hello@webbies.dev, https://webbies.dev)
**License:** MIT
**Source:** https://github.com/webb1es/dotenvify

---

## Brand & Design Direction — JetBrains-Inspired

Follow the JetBrains marketing design language closely. The page should feel like it belongs in the JetBrains ecosystem.

### Color Palette — Dark Mode (default)

| Role             | Color                          | Usage                                                            |
|------------------|--------------------------------|------------------------------------------------------------------|
| Background       | `#000000`                      | Primary page background (pure black, like jetbrains.com)         |
| Surface          | `#1E1E1E`                      | Cards, feature blocks, elevated elements                         |
| Border           | `#2D2D2D`                      | Subtle dividers and card edges                                   |
| Text primary     | `#FFFFFF`                      | Headlines, primary body text                                     |
| Text secondary   | `#A0A0A0`                      | Descriptions, supporting text                                    |
| Gradient primary | `#3B74EE → #510EB2`            | Hero background glow, headline accent (JetBrains blue-to-purple) |
| Gradient accent  | `#087CFA → #21D789`            | Feature highlights (JetBrains blue-to-green)                     |
| CTA primary      | `#087CFA`                      | "Get Plugin" button — solid filled (JetBrains blue)              |
| CTA secondary    | transparent + `#FFFFFF` border | "GitHub" button — outline style                                  |
| Accent warm      | `#FF318C`                      | Badges, highlights, hover states (JetBrains pink)                |
| Success          | `#21D789`                      | Status indicators (JetBrains green)                              |

### Color Palette — Light Mode

| Role             | Color                          | Usage                                      |
|------------------|--------------------------------|--------------------------------------------|
| Background       | `#FFFFFF`                      | Primary page background                    |
| Surface          | `#F5F5F5`                      | Cards, feature blocks, elevated elements   |
| Border           | `#E0E0E0`                      | Subtle dividers and card edges             |
| Text primary     | `#1A1A1A`                      | Headlines, primary body text               |
| Text secondary   | `#6B6B6B`                      | Descriptions, supporting text              |
| Gradient primary | `#3B74EE → #510EB2`            | Same gradients — these work on both themes |
| Gradient accent  | `#087CFA → #21D789`            | Same accent gradient                       |
| CTA primary      | `#087CFA`                      | Same blue — strong contrast on white       |
| CTA secondary    | transparent + `#1A1A1A` border | Outline button with dark border            |
| Accent warm      | `#FF318C`                      | Same pink accent                           |
| Success          | `#21D789`                      | Same green                                 |

### Theme Toggle

- Default to dark mode, but detect `prefers-color-scheme` and respect OS preference
- Provide a sun/moon toggle button in the header bar (right side, before the CTA)
- Store preference in `localStorage` so it persists across visits
- Use CSS custom properties (variables) for all colors so toggling is a single class swap on `<html>`

### Typography

| Element             | Font           | Size    | Weight         |
|---------------------|----------------|---------|----------------|
| Headline            | Inter          | 40-48px | 700 (Bold)     |
| Subheadline         | Inter          | 18-20px | 400 (Regular)  |
| Feature title       | Inter          | 16px    | 600 (Semibold) |
| Feature description | Inter          | 14px    | 400 (Regular)  |
| Code / monospace    | JetBrains Mono | 13px    | 400 (Regular)  |
| CTA buttons         | Inter          | 15px    | 600 (Semibold) |
| Footer              | Inter          | 12px    | 400 (Regular)  |

### Visual Principles

- **Pure black background** with color coming from gradients and accent elements — matches jetbrains.com
- **Gradient glows** behind hero area and feature cards (soft, radial, not sharp edges)
- **No stock photos** — only code snippets, plugin UI mockups, and JetBrains IDE logos
- **Compact density** — tight spacing, no wasted vertical space (everything must fit one viewport)
- **JetBrains button style** — primary buttons are solid filled with rounded corners (~8px), secondary buttons are
  outline with white border
- **Card style** — dark surface (`#1E1E1E`) with subtle `#2D2D2D` border, no shadows, slight hover glow effect
- **Header bar** — fixed top bar (48-56px height) with logo left, nav links center (if any), CTA button right.
  Semi-transparent black with backdrop blur.

---

## Page Layout — Single Viewport (Desktop)

The entire page fits in one screen on desktop. Use a compact, information-dense layout.

```
┌─────────────────────────────────────────────────────────┐
│  [Logo] DotEnvify              [GitHub] [Get Plugin →]  │  ← Header bar (48px)
├─────────────────────────────────────────────────────────┤
│                                                         │
│         Headline (gradient text)                        │
│         Subheadline (one line)                          │  ← Hero (compact, ~20% height)
│         [Get Plugin]  [View on GitHub]                  │
│                                                         │
├────────────┬────────────┬───────────────────────────────┤
│            │            │                               │
│  Azure     │  Paste &   │                               │
│  DevOps    │  Format    │    IDE Mockup / Code          │  ← Main area (~55% height)
│            │            │    Before → After visual      │
│  3-4 lines │  3-4 lines │                               │
│            │            │                               │
├────────────┴──────┬─────┴───────────────────────────────┤
│                   │                                     │
│  Diagnostics      │   IDE logos row + "Works across     │  ← Bottom row (~15% height)
│  3-4 lines        │    all JetBrains IDEs"              │
│                   │   [IntelliJ] [GoLand] [WebStorm]    │
│                   │   [PyCharm] [Rider] [CLion]         │
│                   │                                     │
├───────────────────┴─────────────────────────────────────┤
│  MIT · Webbies · GitHub · hello@webbies.dev    © 2026   │  ← Footer (32px)
└─────────────────────────────────────────────────────────┘
```

### Section Details

#### Header Bar

- Left: Plugin icon + "DotEnvify" wordmark
- Right: "GitHub" (icon + text, outline button) and "Get Plugin" (solid blue `#087CFA` button)
- Background: `rgba(0,0,0,0.8)` with `backdrop-filter: blur(12px)`

#### Hero (compact)

- **Headline** with gradient text (`#3B74EE → #510EB2`): "Azure DevOps variables → .env in one click."
- **Subheadline** (one line, `#A0A0A0`): "A JetBrains plugin that pulls, formats, and manages your environment
  variables — without leaving your IDE."
- Two CTA buttons side by side (JetBrains button style)
- Subtle radial gradient glow behind the text area
- No hero image here — keep it tight

#### Feature Cards (3 columns, left side)

Three compact cards on the left ~40% of the main area:

**Card 1: Azure DevOps**

- Small icon or emoji-free badge
- "Fetch variables from Azure DevOps Variable Groups. Secure Device Code OAuth. Smart merge with per-key conflict
  resolution."

**Card 2: Paste & Format**

- "Parse any format — `KEY=VALUE`, `KEY VALUE`, line pairs. Smart quoting, sorting, export prefix. Live preview as you
  type."

**Card 3: Diagnostics**

- "Find missing and unused env keys across 10 languages. Click to navigate to source. Auto-watch for .env changes."

#### Right Side Visual (~60% of main area)

A single compelling visual — choose one:

- **Option A:** Stylized code block showing before → after transformation (dark editor theme, JetBrains Mono font, with
  syntax-colored output)
- **Option B:** A mockup of the plugin tool window inside IntelliJ (screenshot placeholder with clear `[SCREENSHOT]`
  marker)
- **Option C:** Animated terminal-style demo showing keys being fetched and formatted

#### Bottom Row

- Left: Any additional feature highlight or the before/after code snippet (if not used above)
- Right: Row of JetBrains IDE logos/icons with text: "Works across all JetBrains IDEs · 2024.1+"
- IDE names in small text below each logo

#### Footer (single line)

- `MIT License · Built by Webbies · GitHub · hello@webbies.dev · © 2026`
- All inline, `#A0A0A0` text, `12px`

---

## Feature Copy (condensed for single-viewport)

Keep feature text extremely concise. Each feature card gets 2-3 short sentences max.

**Azure DevOps Integration:**
Fetch variables from Azure DevOps Variable Groups with one click. Secure OAuth sign-in via Device Code Flow. Smart merge
when .env already exists — per-key conflict resolution.

**Paste & Format:**
Parse any input format — KEY=VALUE, quoted, space-separated, line pairs. Smart quoting, alphabetical sorting, export
prefix. Live preview updates as you type.

**Diagnostics:**
Detect missing and unused env keys across your codebase. Supports JS, Python, Go, Java, Kotlin, Ruby, PHP, Rust, C#,
YAML. Click to navigate to source.

---

## Code Example (for the visual area)

Use JetBrains Mono, dark theme styling. Show transformation:

**Input (messy):**

```
API_KEY abc123
DATABASE_URL="postgres://localhost:5432/db"
secret_token mytoken
REDIS_HOST
redis.example.com
export NODE_ENV=production
```

**Output (clean):**

```env
API_KEY=abc123
DATABASE_URL="postgres://localhost:5432/db"
NODE_ENV=production
REDIS_HOST=redis.example.com
```

Show a subtle arrow or transition between them. Add a small caption: "Sorted · Lowercase filtered · Auto-quoted · Export
stripped"

---

## Technical Requirements

- **Framework:** Astro or plain HTML/CSS/JS. Keep it minimal — a single-page site does not need React/Next.js overhead.
- **Styling:** Tailwind CSS. Use CSS custom properties for all theme colors so dark/light toggle is a class swap.
- **Layout:** CSS Grid for the main layout. `height: 100dvh` (dynamic viewport height) on desktop with
  `overflow: hidden`.
- **Animations:** Minimal — a subtle gradient pulse on the hero glow, hover effects on cards and buttons. No
  scroll-based animations (there is no scroll on desktop).
- **Performance:** Lighthouse 95+. Inline critical CSS. Minimal JS (only for theme toggle and any interactive elements).
- **SEO:** Meta tags, Open Graph, Twitter card. Keywords: "JetBrains plugin", "env file", "Azure DevOps", "environment
  variables", "dotenv".
- **Hosting:** Static site — Vercel, Netlify, or Cloudflare Pages.
- **Fonts:** Load Inter (400, 600, 700) and JetBrains Mono (400) from Google Fonts or self-host.
- **Assets:**
    - Plugin icon: `src/main/resources/META-INF/pluginIcon.svg`
    - JetBrains IDE logos: use official SVGs from JetBrains brand resources or simple text labels as fallback
    - Screenshot placeholders: use a styled `div` with dashed border and `[SCREENSHOT]` text

---

## Responsiveness

The page must work well across all screen sizes. Use `100dvh` (not `100vh`) to account for mobile browser chrome.

| Breakpoint           | Layout                                                         | Scroll                  |
|----------------------|----------------------------------------------------------------|-------------------------|
| Desktop (1280px+)    | Full grid, single viewport, no scroll                          | `overflow: hidden`      |
| Laptop (1024-1279px) | Slightly compressed grid, smaller fonts, still single viewport | `overflow: hidden`      |
| Tablet (768-1023px)  | 2-column grid, hero stacked above features                     | Vertical scroll allowed |
| Mobile (<768px)      | Single column, all sections stacked                            | Vertical scroll allowed |

- Use `clamp()` for font sizes so they scale fluidly (e.g., `clamp(28px, 4vw, 48px)` for headline)
- Feature cards stack vertically on mobile, 2-column on tablet, 3-column on desktop
- Code example block hides or collapses on small screens (feature cards take priority)
- CTA buttons go full-width on mobile
- Header collapses to logo + hamburger or just logo + single CTA on mobile
- IDE logo row wraps to 2 rows on mobile if needed
- Test at: 360px, 390px, 768px, 1024px, 1280px, 1440px, 1920px

---

## Accessibility

This is not optional. The page must meet **WCAG 2.1 AA** compliance.

### Color & Contrast

- All text must meet **4.5:1 contrast ratio** against its background (AA standard)
- Large text (24px+ or 18.66px+ bold) must meet **3:1 ratio**
- The secondary text color `#A0A0A0` on `#000000` passes (10.6:1) — but verify `#6B6B6B` on `#FFFFFF` in light mode (
  passes at 5.7:1)
- Do not rely on color alone to convey information — use icons, labels, or patterns alongside color
- Gradient text must have a fallback solid color for readability (`-webkit-background-clip: text` with `color` fallback)

### Keyboard Navigation

- All interactive elements (buttons, links, theme toggle) must be focusable and operable via keyboard
- Visible focus indicators — use a `2px` outline with offset in the accent color (`#087CFA`), not just browser default
- Logical tab order following visual layout (header → hero CTAs → feature cards → IDE logos → footer links)
- Theme toggle must be a `<button>` with `aria-label="Switch to light/dark mode"`

### Semantic HTML

- Use proper heading hierarchy: single `<h1>` for the headline, `<h2>` for feature card titles
- Use `<nav>` for header navigation, `<main>` for content, `<footer>` for footer
- Feature cards should use `<section>` or `<article>` with appropriate headings
- CTA links that look like buttons: use `<a>` with `role="button"` if they navigate, or `<button>` if they perform an
  action
- Code examples: use `<pre><code>` with appropriate `aria-label` or preceding heading for context

### Screen Readers

- All images and icons must have meaningful `alt` text (or `alt=""` and `aria-hidden="true"` for decorative ones)
- IDE logos: `alt="IntelliJ IDEA"`, `alt="GoLand"`, etc.
- Plugin icon: `alt="DotEnvify plugin icon"`
- Theme toggle: `aria-label="Switch to dark mode"` / `"Switch to light mode"` (update dynamically)
- Skip-to-content link: add a visually hidden "Skip to main content" link as the first focusable element
- Code blocks: wrap in a region with `aria-label="Code example: before and after formatting"`

### Motion & Preferences

- Respect `prefers-reduced-motion` — disable gradient pulse, transitions, and hover animations when set
- Respect `prefers-color-scheme` for initial theme (override with manual toggle)
- Respect `prefers-contrast` — if `more`, increase border visibility and reduce gradient subtlety

### Touch & Mobile

- All tap targets must be at least **44x44px** (WCAG 2.5.5)
- Adequate spacing between interactive elements to prevent mis-taps
- No hover-only interactions — anything revealed on hover must also be accessible via tap/focus

---

## Content Guidelines

- Write for developers. Be direct.
- Every word must earn its place — single viewport means zero wasted space.
- No buzzwords, no filler paragraphs.
- Code blocks use real examples with JetBrains Mono.
- Feature descriptions: benefit first, mechanism second.

---

## Deliverables

1. Complete single-page source code (ready to deploy)
2. Desktop: single viewport, no scroll. Mobile/tablet: stacked, scrollable.
3. Dark mode default with light mode toggle (respects OS preference)
4. WCAG 2.1 AA compliant
5. Responsive across all breakpoints (360px to 1920px)
6. All copy finalized and inline
7. Screenshot placeholder slots clearly marked
8. SEO meta tags configured
9. Static build output, deploy-ready

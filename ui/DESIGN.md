---
name: Neo-Console HUD
colors:
  surface: '#18181b'
  surface-dim: '#131316'
  surface-bright: '#39393c'
  surface-container-lowest: '#0e0e11'
  surface-container-low: '#1b1b1e'
  surface-container: '#27272a'
  surface-container-high: '#2a2a2d'
  surface-container-highest: '#353438'
  on-surface: '#e4e1e6'
  on-surface-variant: '#bec8d2'
  inverse-surface: '#e4e1e6'
  inverse-on-surface: '#303033'
  outline: '#88929b'
  outline-variant: '#3e4850'
  surface-tint: '#89ceff'
  primary: '#89ceff'
  on-primary: '#00344d'
  primary-container: '#0ea5e9'
  on-primary-container: '#003751'
  inverse-primary: '#006591'
  secondary: '#7bd0ff'
  on-secondary: '#00354a'
  secondary-container: '#00a6e0'
  on-secondary-container: '#00374d'
  tertiary: '#7bd1fa'
  on-tertiary: '#003547'
  tertiary-container: '#4ba5cc'
  on-tertiary-container: '#00384b'
  error: '#ffb4ab'
  on-error: '#690005'
  error-container: '#93000a'
  on-error-container: '#ffdad6'
  primary-fixed: '#c9e6ff'
  primary-fixed-dim: '#89ceff'
  on-primary-fixed: '#001e2f'
  on-primary-fixed-variant: '#004c6e'
  secondary-fixed: '#c4e7ff'
  secondary-fixed-dim: '#7bd0ff'
  on-secondary-fixed: '#001e2c'
  on-secondary-fixed-variant: '#004c69'
  tertiary-fixed: '#c0e8ff'
  tertiary-fixed-dim: '#7bd1fa'
  on-tertiary-fixed: '#001e2b'
  on-tertiary-fixed-variant: '#004d66'
  background: '#131316'
  on-background: '#e4e1e6'
  surface-variant: '#353438'
  background-base: '#09090b'
  border-subtle: '#27272a'
  border-frost: rgba(56, 189, 248, 0.4)
  text-muted: '#a1a1aa'
  text-foreground: '#ffffff'
  frost-dark: '#0284c7'
  frost-glow: rgba(56, 189, 248, 0.25)
typography:
  display-hero:
    fontFamily: Sora
    fontSize: 3.5rem
    fontWeight: '700'
    lineHeight: 4rem
    letterSpacing: -0.03em
  display-hero-mobile:
    fontFamily: Sora
    fontSize: 2.25rem
    fontWeight: '700'
    lineHeight: 2.75rem
    letterSpacing: -0.02em
  headline-xl:
    fontFamily: Sora
    fontSize: 2rem
    fontWeight: '600'
    lineHeight: 2.5rem
    letterSpacing: -0.02em
  headline-xl-mobile:
    fontFamily: Sora
    fontSize: 1.5rem
    fontWeight: '600'
    lineHeight: 2rem
    letterSpacing: -0.015em
  headline-lg:
    fontFamily: Sora
    fontSize: 1.5rem
    fontWeight: '600'
    lineHeight: 2rem
    letterSpacing: -0.02em
  headline-md:
    fontFamily: Sora
    fontSize: 1.25rem
    fontWeight: '500'
    lineHeight: 1.75rem
    letterSpacing: -0.01em
  body-lg:
    fontFamily: Sora
    fontSize: 1rem
    fontWeight: '400'
    lineHeight: 1.5rem
    letterSpacing: 0em
  body-md:
    fontFamily: Sora
    fontSize: 0.875rem
    fontWeight: '400'
    lineHeight: 1.375rem
    letterSpacing: 0em
  body-sm:
    fontFamily: Sora
    fontSize: 0.75rem
    fontWeight: '400'
    lineHeight: 1.125rem
    letterSpacing: 0.01em
  telemetry-metric:
    fontFamily: Space Grotesk
    fontSize: 1.125rem
    fontWeight: '600'
    lineHeight: 1.5rem
    letterSpacing: -0.01em
  telemetry-label:
    fontFamily: Space Grotesk
    fontSize: 0.75rem
    fontWeight: '500'
    lineHeight: 1rem
    letterSpacing: 0.05em
  label-hud:
    fontFamily: Space Grotesk
    fontSize: 0.6875rem
    fontWeight: '700'
    lineHeight: 0.875rem
    letterSpacing: 0.08em
rounded:
  sm: 0.25rem
  DEFAULT: 0.5rem
  md: 0.75rem
  lg: 1rem
  xl: 1.5rem
  full: 9999px
spacing:
  gutter: 1.5rem
  gutter-mobile: 1rem
  margin: 2rem
  margin-mobile: 1rem
  space-xs: 0.25rem
  space-sm: 0.5rem
  space-md: 1rem
  space-lg: 1.5rem
  space-xl: 2.5rem
---

## Brand & Style

This design system establishes a high-performance "Clean HUD / Neo-Console" aesthetic tailored for a self-hosted game launcher and streaming client. The visual direction merges cutting-edge console firmware ergonomics with telemetry-grade heads-up displays (HUD), evoking the immersion, speed, and precision of high-end gaming hardware like Steam Deck, Moonlight, and next-gen console shells.

Key stylistic pillars:
- **Neo-Console Ergonomics**: Deep Zinc bases with soft ambient backdrops, pill capsule tags, and smooth squircles crafted for dual 10-foot gamepad navigation and desktop mouse-keyboard fidelity.
- **Luminescent Frost Accents**: Vibrant electric frost cyan highlights that slice through deep charcoal layers, serving as focused interactive affordances and low-latency feedback indicators.
- **Frosted HUD Translucency**: Subtle glassmorphism (`backdrop-blur-md`) with soft inner ring glows, providing depth without sluggish visual overhead.
- **Telemetry Precision**: Clean, tabular telemetry readouts for stream quality, bitrates, hardware status, and latency benchmarks directly embedded into the user interface hierarchy.

## Colors

The palette operates in a dedicated, immersive dark mode. High-contrast typography and vivid frost cyan accents pierce deep zinc layers to prioritize game art, active stream telemetry, and rapid navigational recognition.

- **Background Baseline (`#09090b`)**: True Zinc 950 backdrop. Provides maximum contrast for high-res hero art and launcher posters.
- **Surfaces (`#18181b` / `#27272a`)**: Zinc 900 base containers with Zinc 800 elevated cards, HUD bars, and floating utility panels.
- **Borders (`#27272a`)**: Low-contrast boundary lines preserving spatial structure without distracting from cover art or video feeds.
- **Brand Frost Continuum**:
  - `primary_color_hex` (`#0ea5e9` - Frost 500): Core call-to-actions, active launcher tabs, and primary controller focus rings.
  - `secondary_color_hex` (`#38bdf8` - Frost 400): Interactive highlights, hover flares, and stream indicators.
  - `tertiary_color_hex` (`#7dd3fc` - Frost 300): Subtle telemetry values, active status chips, and crisp text badges.
  - `frost-dark` (`#0284c7` - Frost 600): Pressed button depths and gradient base anchors.
- **Typography Foregrounds**: `#ffffff` for dominant headers, titles, and active inputs; `#a1a1aa` (Zinc 400) for timestamps, secondary specs, and passive launcher telemetry.

## Typography

Typography establishes an intentional balance between cutting-edge console firmware and technical performance telemetry:

- **Sora (Primary Headings & Core UI)**: Imparts a sleek, high-tech geometric curve profile evocative of contemporary console UI shells. Sora handles primary game titles, hero headers, body descriptions, and button controls.
- **Space Grotesk (Telemetry, Stats & HUD Badges)**: Deployed strictly for technical indicators, performance stats (FPS, bitrate, network jitter, ping), hardware overlays, and platform identifiers. Space Grotesk features wide, technological glyph proportions and tabular-friendly numbers that align seamlessly across dynamic HUD data counters.
- Subpixel text anti-aliasing and optical tracking tuning must be active across all dark surfaces to eliminate blur on lower-resolution TV displays and portable handheld screens.

## Layout & Spacing

The layout is built around an expansive, fluid console canvas supporting both mouse pointers and directional controller focus rings:

- **Shell Architecture**: Full-bleed application shell with a pinned top horizontal HUD bar and a fluid game grid spanning the viewport.
- **Desktop & Large Displays (>= 1280px)**: 12-column adaptive grid, margins set to `2rem` (`32px`), and gutters at `1.5rem` (`24px`). Poster cards form a 4-to-6 column fluid arrangement.
- **Handheld & Tablet Formats (768px - 1279px)**: 8-column layout with `1.25rem` gutters and `1.5rem` outer margins. Quick-access HUD bar elements collapse into expandable pill drawers.
- **Compact Mobile / Remote Play Handhelds (< 768px)**: 4-column layout with `1rem` gutters and outer margins, transitioning library grids to 2-column or 1-column scroll carousels with horizontal swipe affordances.
- **Spacing Rhythm**: Element gaps leverage `space-xs` (4px) and `space-sm` (8px) for tightly grouped telemetry pairs and button icons. Major content cards and HUD modules utilize `space-md` (16px) and `space-lg` (24px).

## Elevation & Depth

Spatial layering relies on glassmorphic translucency, ambient backdrop diffusion, and high-tech frost halos:

- **Canvas Base**: `#09090b` flat ground plane.
- **Elevated Translucent HUD (Surface 1)**: Glass-sheen layers using `background: rgba(24, 24, 27, 0.8)` paired with `backdrop-blur-md` and a faint hairline ring `1px solid rgba(255, 255, 255, 0.08)`.
- **Game Poster Elevation (Surface 2)**: Zinc 900 base (`#18181b`) with soft ambient drop shadows `0 12px 30px -8px rgba(0, 0, 0, 0.8)`. On gamepad selection or cursor hover, the card scales up smoothly (`transform: scale(1.04)`) with an elevated shadow of `0 20px 40px -10px rgba(0, 0, 0, 0.9)`.
- **Frost Glow Effect**: Applied to focused, active, and streaming states:
  - `frost-glow-sm`: `0 0 15px -3px rgba(56, 189, 248, 0.2)`
  - `frost-glow-md`: `0 0 20px -5px rgba(56, 189, 248, 0.25)`
  - `frost-glow-lg`: `0 0 35px -5px rgba(14, 165, 233, 0.4)`
- **Glassmorphic Metadata Plates**: Bottom overlays on poster cards employ gradient dark-to-frosted transitions with `backdrop-blur-sm` to ensure metadata stays readable regardless of background game keyart.

## Shapes

The interface embraces a neo-console geometry emphasizing capsules and smooth squircles:

- **Buttons & HUD Badges**: Pill-shaped capsules (`rounded-full` / 9999px) provide an ergonomic, touchable and controller-ready profile.
- **Poster Cards & Panels**: Smooth squircles (`rounded-2xl` / 16px to 24px) that soften the visual field and frame video game cover art gracefully.
- **Modals & Flyout Sheets**: Deep rounded corners (`rounded-3xl` / 32px) for overlay sheets and quick-access side consoles.
- **Inputs & Dropdowns**: Squircle rounded corners (`rounded-xl` / 12px) balancing structure with comfortable targets.

## Components

### Buttons
- **Primary Action (Capsule Frost)**: Pill shape (`rounded-full`), deep dark base (`rgba(14, 165, 233, 0.15)`), frosted border `ring-1 ring-frost-500/40`, text `#ffffff` in `Sora` semi-bold. On hover or controller focus: switches to solid `#0ea5e9` with dark text `#09090b` and triggers `0 0 20px -5px rgba(56, 189, 248, 0.25)`.
- **Secondary Squircle Button**: `rounded-2xl`, background `#18181b/80` with `backdrop-blur-md`, border `1px solid #27272a`. Hover: border brightens to `#38bdf8/40` with subtle scale transition (`scale-[1.02]`).
- **Icon Utility Controls**: Pill or circular capsules (`w-10 h-10 rounded-full`) with Zinc 800 surfaces, housing controller glyphs or quick toggles.

### Game Poster Cards
- Built with `rounded-2xl` overflow-hidden clipping and high-aspect game box ratios (typically 2:3 or 16:9 for banners).
- Base card: ambient drop shadow and `1px solid #27272a` boundary.
- Interaction: Smooth spring scale transition on hover/focus (`transition-transform duration-300 ease-out transform hover:scale-[1.03]`), accompanied by an outer cyan rim glow (`ring-1 ring-frost-400/50`).
- Frosted Bottom Metadata Plate: Translucent gradient overlay (`bg-gradient-to-t from-zinc-950 via-zinc-950/80 to-transparent backdrop-blur-sm`) displaying game title in `Sora` and stream/platform compatibility in `Space Grotesk`.

### Badges & HUD Chips
- **Status Indicators**: Pill-shaped (`rounded-full px-3 py-1 bg-surface/80 border border-frost-500/30 text-frost-300`).
- **Telemetry Readouts**: Formatted in `Space Grotesk` (`telemetry-label`), combining an animated 6px pulse dot with real-time stats (e.g., `60 FPS`, `3.2 ms`, `120 Mbps`, `Moonlight Active`).
- **Controller/Hardware Tags**: Compact pills with micro icon and uppercase label in `#a1a1aa`.

### Navigation (Top Horizontal HUD Bar)
- Persistent horizontal shell (`backdrop-blur-md bg-zinc-950/80 border-b border-zinc-800/80 h-16 px-6`).
- Navigation Tabs: Expansive pill-shaped indicators (`rounded-full px-5 py-2 text-sm font-medium transition-all`). Active tab features a sliding frosted cyan background (`bg-frost-500/15 text-frost-300 border border-frost-500/40 shadow-[0_0_15px_-3px_rgba(56,189,248,0.25)]`).
- Right HUD Cluster: Quick-glance system telemetry pills, active host server status (Sunshine/Moonlight), connection ping indicator, and gamepad battery icons.

### Input Fields & Search
- Rounded squircle text bars (`rounded-full` for search inputs, `rounded-xl` for config forms).
- Background `#09090b` with border `#27272a`. Focused state blooms with `border-frost-500` and frost halo ring.
- Integrated controller shortcut glyphs (e.g., `(Y) Search` or `[/] Filter`) rendered in muted `Space Grotesk`.

### Stream Telemetry Overlay (OSD)
- Detachable floating HUD card: `rounded-2xl bg-zinc-900/90 backdrop-blur-md border border-frost-500/20 p-4 shadow-2xl`.
- Grid arrangement displaying real-time stream bitrate graphs, frame times, packet loss, and host GPU thermals powered by `Space Grotesk` tabular figures.
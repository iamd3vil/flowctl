The idiomatic approach for theming in Tailwind v4 uses **CSS variables with the `@theme` directive**, combined with the `@variant` directive for dark mode. This pattern is flexible enough to support multiple themes beyond just light/dark.

## The Core Pattern

Here's the recommended approach:

### 1. Define Semantic Theme Variables

```css
/* styles.css */
@import "tailwindcss";

/* Override the dark variant to use a selector (instead of media query) */
@variant dark (&:where(.dark, .dark *));

/* Define your semantic color tokens that reference CSS variables */
@theme {
  --color-background: var(--bg);
  --color-foreground: var(--fg);
  --color-primary: var(--primary);
  --color-primary-foreground: var(--primary-fg);
  --color-muted: var(--muted);
  --color-muted-foreground: var(--muted-fg);
  --color-border: var(--border);
  --color-card: var(--card);
  --color-card-foreground: var(--card-fg);
}

/* Define the actual values per theme in @layer base */
@layer base {
  /* Light theme (default) */
  :root {
    --bg: #ffffff;
    --fg: #020817;
    --primary: #1e40af;
    --primary-fg: #eff6ff;
    --muted: #f1f5f9;
    --muted-fg: #64748b;
    --border: #e2e8f0;
    --card: #ffffff;
    --card-fg: #020817;
  }

  /* Dark theme */
  .dark {
    --bg: #020817;
    --fg: #f8fafc;
    --primary: #3b82f6;
    --primary-fg: #eff6ff;
    --muted: #1e293b;
    --muted-fg: #94a3b8;
    --border: #1e293b;
    --card: #0f172a;
    --card-fg: #f8fafc;
  }
}
```

### 2. Usage in HTML

```html
<div class="bg-background text-foreground">
  <div class="bg-card text-card-foreground border-border rounded-lg p-4">
    <button class="bg-primary text-primary-foreground px-4 py-2 rounded">
      Click me
    </button>
  </div>
</div>
```

## Extending to Multiple Themes

The beauty of this pattern is that you can easily add more themes:

```css
@layer base {
  /* ... existing light/dark themes ... */

  /* Ocean theme */
  .theme-ocean {
    --bg: #0c4a6e;
    --fg: #f0f9ff;
    --primary: #22d3ee;
    --primary-fg: #083344;
    /* ... */
  }

  /* Forest theme */
  .theme-forest {
    --bg: #14532d;
    --fg: #f0fdf4;
    --primary: #4ade80;
    --primary-fg: #052e16;
    /* ... */
  }

  /* Combine with dark mode */
  .dark.theme-ocean {
    --bg: #082f49;
    --primary: #67e8f9;
    /* ... */
  }
}
```

## Theme Toggle (JavaScript)

```javascript
// Simple theme manager
const ThemeManager = {
  setTheme(theme) {
    const root = document.documentElement;
    
    // Clear existing theme classes
    root.classList.remove('dark', 'theme-ocean', 'theme-forest');
    
    // Apply new theme
    if (theme === 'dark') {
      root.classList.add('dark');
    } else if (theme === 'system') {
      if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
        root.classList.add('dark');
      }
    } else if (theme !== 'light') {
      root.classList.add(theme); // e.g., 'theme-ocean'
    }
    
    localStorage.setItem('theme', theme);
  },

  init() {
    const saved = localStorage.getItem('theme') || 'system';
    this.setTheme(saved);

    // Listen for system preference changes
    window.matchMedia('(prefers-color-scheme: dark)')
      .addEventListener('change', (e) => {
        if (localStorage.getItem('theme') === 'system') {
          this.setTheme('system');
        }
      });
  }
};
```

## Using `data-` Attributes Instead of Classes

If you prefer data attributes (cleaner separation):

```css
@variant dark (&:where([data-theme="dark"], [data-theme="dark"] *));

@layer base {
  :root, [data-theme="light"] {
    --bg: #ffffff;
    /* ... */
  }

  [data-theme="dark"] {
    --bg: #020817;
    /* ... */
  }

  [data-theme="ocean"] {
    --bg: #0c4a6e;
    /* ... */
  }
}
```

Then toggle with: `document.documentElement.dataset.theme = 'dark'`

## Key Takeaways

1. **Use `@theme` for semantic tokens** that map CSS variables to Tailwind utility classes
2. **Use `@layer base`** to define the actual color values per theme
3. **Use `@variant dark`** to customize how dark mode is triggered
4. **Use indirection** (semantic variables referencing raw variables) so values can be swapped by changing a class or attribute
5. **Keep theme names semantic** (e.g., `--color-primary`, not `--color-blue-500`)

This approach gives you the flexibility to add any number of themes while keeping your component markup clean and theme-agnostic.

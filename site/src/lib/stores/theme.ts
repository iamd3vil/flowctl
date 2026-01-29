import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';

export type Theme = 'light' | 'dark' | 'system';
export type ResolvedTheme = 'light' | 'dark';

const getInitialTheme = (): Theme => {
    if (!browser) return 'system';
    return (localStorage.getItem('theme') as Theme) || 'system';
};

function createThemeStore() {
    const { subscribe, set } = writable<Theme>(getInitialTheme());

    return {
        subscribe,
        set: (theme: Theme) => {
            if (browser) localStorage.setItem('theme', theme);
            set(theme);
        },
    };
}

export const theme = createThemeStore();

export const resolvedTheme = derived<typeof theme, ResolvedTheme>(
    theme,
    ($theme, set) => {
        if (!browser) { set('light'); return; }

        if ($theme === 'system') {
            const mq = window.matchMedia('(prefers-color-scheme: dark)');
            const update = () => set(mq.matches ? 'dark' : 'light');
            update();
            mq.addEventListener('change', update);
            return () => mq.removeEventListener('change', update);
        }
        set($theme);
    },
    'light'
);

export function applyTheme(resolved: ResolvedTheme) {
    if (!browser) return;
    document.documentElement.classList.remove('light', 'dark');
    document.documentElement.classList.add(resolved);
}

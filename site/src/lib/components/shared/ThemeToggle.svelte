<script lang="ts">
    import { theme, type Theme } from '$lib/stores/theme';
    import { IconSun, IconMoon, IconDeviceDesktop } from '@tabler/icons-svelte';

    let { collapsed = false }: { collapsed?: boolean } = $props();

    const icons = { light: IconSun, dark: IconMoon, system: IconDeviceDesktop };
    const order: Theme[] = ['light', 'dark', 'system'];

    function cycle() {
        const next = order[(order.indexOf($theme) + 1) % 3];
        theme.set(next);
    }
</script>

<button
    type="button"
    onclick={cycle}
    class="flex items-center justify-center p-2 text-muted-foreground
           hover:bg-subtle rounded-lg transition-colors cursor-pointer
           {collapsed ? '' : 'w-full px-4'}"
    title="{$theme} theme"
>
    <svelte:component this={icons[$theme]} size={20} />
    {#if !collapsed}
        <span class="ml-3 capitalize">{$theme}</span>
    {/if}
</button>

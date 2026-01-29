<script lang="ts">
	import { beforeNavigate, afterNavigate } from '$app/navigation';
	import { page } from '$app/stores';
	import { isLoading } from '$lib/stores/auth';
	import Logo from './Logo.svelte';
	import LoadingSpinner from './LoadingSpinner.svelte';

	let navigating = $state(false);

	// Use $derived for reactive store values (Svelte 5 pattern)
	let initialLoading = $derived($isLoading);
	let currentPath = $derived($page.url.pathname);
	let isRootPage = $derived(currentPath === '/');
	let isLoginPage = $derived(currentPath === '/login' || currentPath.startsWith('/login/'));

	beforeNavigate(() => {
		navigating = true;
	});

	afterNavigate(() => {
		navigating = false;
	});
</script>

{#if !isRootPage && !isLoginPage && initialLoading}
	<!-- Splash screen for initial page load -->
	<div
		role="status"
		aria-live="polite"
		aria-label="Loading application"
		class="fixed inset-0 z-50 flex items-center justify-center bg-card"
	>
		<div class="flex flex-col items-center gap-6">
			<Logo height="h-16" />
			<LoadingSpinner label="Loading..." />
		</div>
	</div>
{:else if !isRootPage && !isLoginPage && navigating}
	<!-- Loading indicator for navigation between pages -->
	<div
		role="status"
		aria-live="polite"
		aria-label="Loading page"
		class="fixed inset-0 z-50 flex items-center justify-center bg-card/60 backdrop-blur-sm transition-opacity duration-200"
	>
		<div class="flex flex-col items-center gap-3">
			<LoadingSpinner size="lg" />
		</div>
	</div>
{/if}

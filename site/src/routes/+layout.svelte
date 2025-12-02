<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { currentUser, isAuthenticated, isLoading } from '$lib/stores/auth';
	import NotificationPopup from '$lib/components/shared/NotificationPopup.svelte';
	import GlobalLoadingIndicator from '$lib/components/shared/GlobalLoadingIndicator.svelte';

	let { children, data } = $props();

	// Update stores when userPromise resolves
	$effect(() => {
		let cancelled = false;

		data.userPromise
			.then((user) => {
				if (!cancelled) {
					currentUser.set(user);
					isAuthenticated.set(!!user);
					isLoading.set(false);
				}
			})
			.catch((error) => {
				if (!cancelled) {
					console.error('[Auth] Failed to load user profile:', error);
					currentUser.set(null);
					isAuthenticated.set(false);
					isLoading.set(false);
				}
			});

		return () => {
			cancelled = true;
		};
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
	<link href="https://fonts.googleapis.com/css2?family=Bitter:ital,wght@0,100..900;1,100..900&family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&display=swap" rel="stylesheet">
</svelte:head>

{@render children?.()}

<NotificationPopup />
<GlobalLoadingIndicator />

<script lang="ts">
  import { apiClient } from '$lib/apiClient';
  import { goto, invalidateAll } from '$app/navigation';
  import { page } from '$app/stores';
  import { handleInlineError } from '$lib/utils/errorHandling';
  import { getDefaultNamespace } from '$lib/utils/navigation';
  import { isLoading } from '$lib/stores/auth';
  import { onMount } from 'svelte';
  import Logo from '$lib/components/shared/Logo.svelte';
  import LoginCard from '$lib/components/login/LoginCard.svelte';
  import Footer from '$lib/components/login/Footer.svelte';
  import ThemeToggle from '$lib/components/shared/ThemeToggle.svelte';

  let username = $state('');
  let password = $state('');
  let loading = $state(false);
  let error = $state('');

  const redirectUrl = $derived($page.url.searchParams.get('redirect_url'));

  // Reset loading state when landing on login page
  onMount(() => {
    isLoading.set(false);
  });

  const submit = async (event: SubmitEvent) => {
    event.preventDefault();
    if (!username || !password) {
      return;
    }

    loading = true;

    try {
      await apiClient.auth.login({ username, password });
      await invalidateAll();
      // Keep loading state active until navigation completes
      if (redirectUrl && redirectUrl.startsWith('/')) {
        goto(redirectUrl);
      } else {
        const namespace = await getDefaultNamespace();
        goto(`/view/${namespace}/flows`);
      }
    } catch (err) {
      handleInlineError(err, 'Unable to Sign In');
      loading = false;
    }
  };

</script>

<svelte:head>
  <title>Login - Flowctl</title>
</svelte:head>

<main class="min-h-screen flex items-center justify-center bg-muted px-4 relative">
  <section class="w-full max-w-md">
    <div class="mb-8 flex justify-center p-4">
      <Logo height="h-14" />
    </div>
    <LoginCard
      onSubmit={submit}
      {loading}
      {error}
      bind:username
      bind:password
      {redirectUrl}
    />
    <Footer />
  </section>
  <div class="absolute bottom-4 right-4">
    <ThemeToggle collapsed={true} />
  </div>
</main>

<script lang="ts">
  import { apiClient } from '$lib/apiClient';
  import { goto, invalidateAll } from '$app/navigation';
  import Logo from '$lib/components/login/Logo.svelte';
  import LoginCard from '$lib/components/login/LoginCard.svelte';
  import Footer from '$lib/components/login/Footer.svelte';

  let username = $state('');
  let password = $state('');
  let loading = $state(false);
  let error = $state('');

  const submit = async (event: SubmitEvent) => {
    event.preventDefault();
    if (!username || !password) {
      error = 'Please fill in all fields';
      return;
    }

    loading = true;
    error = '';

    try {
      await apiClient.auth.login({ username, password });
      await invalidateAll(); // Refresh all load functions including layout
      goto('/view/default/flows');
    } catch (err) {
      console.error('Login failed:', err);
      error = 'Invalid username or password';
    } finally {
      loading = false;
    }
  };

</script>

<svelte:head>
  <title>Login - Flowctl</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@tabler/icons-webfont@latest/tabler-icons.min.css">
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-slate-50">
  <div class="w-full max-w-md">
    <Logo />
    <LoginCard 
      onSubmit={submit} 
      {error} 
      {loading} 
      bind:username 
      bind:password 
    />
    <Footer />
  </div>
</div>
<script lang="ts">
  let { 
    name,
    description,
    icon = `<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"/>
    </svg>`
  }: {
    name: string,
    description?: string,
    icon?: string
  } = $props();
</script>

<style>
  .gradient-bg {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }
</style>

<!-- Hero Section -->
<div class="gradient-bg px-6 py-8">
  <div class="max-w-4xl mx-auto text-center">
    <div class="flex items-center justify-center mb-4">
      <div class="w-12 h-12 bg-white/20 backdrop-blur-sm rounded-xl flex items-center justify-center">
        {@html icon}
      </div>
    </div>
    <h1 class="text-3xl font-bold text-white mb-2">{name}</h1>
    {#if description}
      <p class="text-lg text-white/90 max-w-xl mx-auto">{description}</p>
    {/if}
  </div>
</div>
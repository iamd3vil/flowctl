<script lang="ts">
	import { page } from '$app/state';
	import type { PageData } from './$types';
	import PageHeader from '$lib/components/shared/PageHeader.svelte';
	import Header from '$lib/components/shared/Header.svelte';
	import NamespaceSecretsTab from '$lib/components/namespace-secrets/NamespaceSecretsTab.svelte';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
  <title>Secrets - {page.params.namespace} - Flowctl</title>
</svelte:head>

<Header breadcrumbs={[
  { label: page.params.namespace!, url: `/view/${page.params.namespace}/flows` },
  { label: "Secrets" }
]}>
  {#snippet children()}
    <div class="mb-10"></div>
  {/snippet}
</Header>

<div class="p-12">
	<NamespaceSecretsTab
		namespace={data.namespace}
		disabled={!data.permissions?.canCreate}
	/>
</div>

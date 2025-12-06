<script lang="ts">
	import { handleInlineError } from '$lib/utils/errorHandling';
	import { autofocus } from '$lib/utils/autofocus';
	import type { FlowSecretReq, FlowSecretUpdateReq, FlowSecretResp } from '$lib/types';
	import { IconEye, IconEyeOff } from '@tabler/icons-svelte';

	interface Props {
		isEditMode?: boolean;
		secretData?: FlowSecretResp | null;
		onSave: (secretData: FlowSecretReq | FlowSecretUpdateReq) => void;
		onClose: () => void;
	}

	let {
		isEditMode = false,
		secretData = null,
		onSave,
		onClose
	}: Props = $props();

	// Form state
	let formData = $state({
		key: '',
		value: '',
		description: ''
	});

	let loading = $state(false);
	let showValue = $state(false);

	// Initialize form data when secretData changes
	$effect(() => {
		if (isEditMode && secretData) {
			formData = {
				key: secretData.key || '',
				value: '', // Don't load existing secret value for security
				description: secretData.description || ''
			};
		} else if (!isEditMode) {
			// Reset form for new secret
			formData = {
				key: '',
				value: '',
				description: ''
			};
		}
	});

	async function handleSubmit() {
		try {
			loading = true;

			const secretFormData: FlowSecretReq | FlowSecretUpdateReq = isEditMode
				? {
					value: formData.value,
					description: formData.description || undefined
				}
				: {
					key: formData.key,
					value: formData.value,
					description: formData.description || undefined
				};

			await onSave(secretFormData);
		} catch (error) {
			handleInlineError(error);
		} finally {
			loading = false;
		}
	}
</script>

<!-- Modal overlay -->
<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="fixed inset-0 z-50 flex items-center justify-center bg-gray-900/60" onclick={onClose} role="dialog" aria-modal="true">
	<!-- Modal content -->
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div class="bg-white rounded-lg shadow-xl w-full max-w-md mx-4" onclick={(e) => e.stopPropagation()}>
		<div class="p-6">
			<h2 class="text-xl font-semibold mb-4">
				{isEditMode ? 'Edit Secret' : 'Add New Secret'}
			</h2>
			
			<form onsubmit={handleSubmit} class="space-y-4">
				<!-- Secret Key -->
				<div>
					<label for="key" class="block text-sm font-medium text-gray-700 mb-1">
						Key <span class="text-red-500">*</span>
					</label>
					<input
						type="text"
						id="key"
						bind:value={formData.key}
						required
						disabled={loading || isEditMode}
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent disabled:bg-gray-100 disabled:text-gray-500"
						placeholder="SECRET_KEY"
						use:autofocus
					/>
					{#if isEditMode}
						<p class="mt-1 text-xs text-gray-500">Key cannot be changed. Delete and recreate the secret to use a different key.</p>
					{/if}
				</div>

				<!-- Secret Value -->
				<div>
					<label for="value" class="block text-sm font-medium text-gray-700 mb-1">
						Value <span class="text-red-500">*</span>
					</label>
					<div class="relative">
						<input
							type={showValue ? 'text' : 'password'}
							id="value"
							bind:value={formData.value}
							required
							disabled={loading}
							class="w-full px-3 py-2 pr-10 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
							placeholder={isEditMode ? 'Enter new value to update' : 'Enter secret value'}
						/>
						<button
							type="button"
							onclick={() => showValue = !showValue}
							class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 cursor-pointer"
							title={showValue ? 'Hide value' : 'Show value'}
						>
							{#if showValue}
								<IconEyeOff size={20} />
							{:else}
								<IconEye size={20} />
							{/if}
						</button>
					</div>
					{#if isEditMode}
						<p class="mt-1 text-xs text-gray-500">Enter a new value to update. Previous value is not shown for security.</p>
					{/if}
				</div>

				<!-- Description -->
				<div>
					<label for="description" class="block text-sm font-medium text-gray-700 mb-1">
						Description
					</label>
					<textarea
						id="description"
						bind:value={formData.description}
						disabled={loading}
						rows="3"
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
						placeholder="Optional description"
					></textarea>
				</div>

				<!-- Action buttons -->
				<div class="flex justify-end gap-2 pt-4">
					<button
						type="button"
						onclick={onClose}
						disabled={loading}
						class="px-4 py-2 text-gray-600 bg-gray-100 rounded hover:bg-gray-200 disabled:opacity-50 cursor-pointer"
					>
						Cancel
					</button>
					<button
						type="submit"
						disabled={loading || !formData.key || !formData.value}
						class="px-4 py-2 bg-primary-500 text-white rounded hover:bg-primary-600 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
					>
						{loading ? 'Saving...' : isEditMode ? 'Update' : 'Create'}
					</button>
				</div>
			</form>
		</div>
	</div>
</div>
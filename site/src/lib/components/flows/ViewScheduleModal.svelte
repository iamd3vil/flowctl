<script lang="ts">
	import { autofocus } from '$lib/utils/autofocus';
	import JsonDisplay from '$lib/components/shared/JsonDisplay.svelte';
	import type { UserSchedule } from '$lib/types';

	interface Props {
		schedule: UserSchedule;
		onClose: () => void;
	}

	let { schedule, onClose }: Props = $props();

	function handleClose() {
		onClose();
	}

	// Close on Escape key
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			handleClose();
		}
	}
</script>

<svelte:window on:keydown={handleKeydown} />

<!-- Modal Backdrop -->
<div
	class="fixed inset-0 z-50 flex items-center justify-center bg-overlay"
	onclick={handleClose}
	onkeydown={(e) => e.key === 'Escape' && handleClose()}
	role="dialog"
	aria-modal="true"
	tabindex="-1"
>
	<!-- Modal Content -->
	<div
		class="bg-card rounded-lg shadow-lg w-full max-w-lg p-6 m-4"
		onclick={(e) => e.stopPropagation()}
	>
		<h3 class="font-bold text-lg mb-4 text-foreground">Schedule Details</h3>

		<!-- Schedule Info Box -->
		<div class="bg-muted border border-border rounded-lg p-4 mb-4">
			<div class="grid grid-cols-2 gap-3">
				<div>
					<div class="text-xs font-medium text-muted-foreground uppercase mb-1">Cron Expression</div>
					<code class="bg-card px-2 py-1 rounded text-sm font-mono border border-input text-foreground">{schedule.cron}</code>
				</div>
				<div>
					<div class="text-xs font-medium text-muted-foreground uppercase mb-1">Timezone</div>
					<div class="text-sm text-foreground">{schedule.timezone}</div>
				</div>
			</div>
		</div>

		<!-- Schedule Inputs -->
		<div class="mb-4">
			<JsonDisplay data={schedule.inputs} title="Inputs" expanded={true} />
		</div>

		<!-- Footer -->
		<div class="flex justify-end">
			<button
				type="button"
				class="inline-flex items-center px-5 py-2.5 text-sm font-medium text-foreground bg-subtle rounded-lg hover:bg-subtle-hover cursor-pointer"
				onclick={handleClose}
				use:autofocus
			>
				Close
			</button>
		</div>
	</div>
</div>

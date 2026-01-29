<script lang="ts">
	import type { FlowInput } from '$lib/types';

	let {
		inputs = [],
		values = $bindable({}),
		errors = {},
		useFormData = false
	}: {
		inputs?: FlowInput[];
		values?: Record<string, any>;
		errors?: Record<string, string>;
		useFormData?: boolean;
	} = $props();
</script>

{#if inputs && inputs.length > 0}
	<div class="space-y-4">
	{#each inputs as input (input.name)}
		<div>
			<label for={input.name} class="block text-sm font-medium text-foreground mb-2">
				{input.label || input.name}
				{#if input.required}
					<span class="text-red-500">*</span>
				{/if}
			</label>

			{#if errors[input.name]}
				<p class="text-sm text-danger-600 mb-2">{errors[input.name]}</p>
			{/if}

			{#if input.type === 'string' || input.type === 'number'}
				{#if useFormData}
					<input
						type={input.type === 'string' ? 'text' : 'number'}
						id={input.name}
						name={input.name}
						value={values[input.name] ?? input.default ?? ''}
						placeholder={input.description || ''}
						required={input.required}
						class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
					/>
				{:else}
					<input
						type={input.type === 'string' ? 'text' : 'number'}
						bind:value={values[input.name]}
						placeholder={input.description || ''}
						required={input.required}
						class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
					/>
				{/if}
			{:else if input.type === 'checkbox'}
				<div class="flex items-center">
					{#if useFormData}
						<input
							type="checkbox"
							id={input.name}
							name={input.name}
							value="true"
							checked={values[input.name] ?? input.default === 'true'}
							class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-input rounded"
						/>
					{:else}
						<input
							type="checkbox"
							bind:checked={values[input.name]}
							class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-input rounded"
						/>
					{/if}
				</div>
			{:else if input.type === 'select' && input.options}
				{#if useFormData}
					<select
						id={input.name}
						name={input.name}
						required={input.required}
						value={values[input.name] ?? input.default ?? ''}
						class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
					>
						<option value="">Select an option</option>
						{#each input.options as option}
							<option value={option} selected={option === (values[input.name] ?? input.default)}
								>{option}</option
							>
						{/each}
					</select>
				{:else}
					<select
						bind:value={values[input.name]}
						required={input.required}
						class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
					>
						<option value="">Select an option</option>
						{#each input.options as option}
							<option value={option}>{option}</option>
						{/each}
					</select>
				{/if}
			{:else if input.type === 'file'}
				<div class="flex flex-col">
					<input
						type="file"
						id={input.name}
						name={input.name}
						required={input.required}
						class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-medium file:bg-primary-50 file:text-primary-700 hover:file:bg-primary-100"
					/>
					{#if input.max_file_size}
						<p class="text-xs text-muted-foreground mt-1">
							Max size: {Math.round(input.max_file_size / (1024 * 1024))}MB
						</p>
					{/if}
				</div>
			{:else if input.type === 'datetime'}
				<div class="flex items-center">
					{#if useFormData}
						<input
							type="datetime-local"
							id={input.name}
							name={input.name}
							value={values[input.name] ?? input.default ?? ''}
							required={input.required}
							class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
						/>
					{:else}
						<input
							type="datetime-local"
							bind:value={values[input.name]}
							required={input.required}
							class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
						/>
					{/if}
				</div>
			{:else if input.type === 'password'}
				<div class="flex items-center">
					{#if useFormData}
						<input
							type="password"
							id={input.name}
							name={input.name}
							value={values[input.name] ?? input.default ?? ''}
							required={input.required}
							class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
						/>
					{:else}
						<input
							type="password"
							bind:value={values[input.name]}
							required={input.required}
							class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
						/>
					{/if}
				</div>
			{:else}
				<!-- Fallback for other input types -->
				{#if useFormData}
					<input
						type="text"
						id={input.name}
						name={input.name}
						value={values[input.name] ?? input.default ?? ''}
						placeholder={input.description || ''}
						required={input.required}
						class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
					/>
				{:else}
					<input
						type="text"
						bind:value={values[input.name]}
						placeholder={input.description || ''}
						required={input.required}
						class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
					/>
				{/if}
			{/if}

			{#if input.description}
				<p class="text-sm text-muted-foreground mt-1">{input.description}</p>
			{/if}
		</div>
	{/each}
	</div>
{/if}

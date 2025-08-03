<script lang="ts">
	interface Props {
		title: string;
		value: number | string;
		icon: string;
		color?: 'blue' | 'green' | 'purple' | 'red' | 'yellow' | 'gray';
	}

	let { title, value, icon, color = 'blue' }: Props = $props();

	const colorClasses = {
		blue: 'bg-blue-100 text-blue-600',
		green: 'bg-green-100 text-green-600',
		purple: 'bg-purple-100 text-purple-600',
		red: 'bg-red-100 text-red-600',
		yellow: 'bg-yellow-100 text-yellow-600',
		gray: 'bg-gray-100 text-gray-600'
	};
</script>

<div class="bg-white rounded-lg border border-gray-200 p-6">
	<div class="flex items-center justify-between">
		<div>
			<p class="text-sm font-medium text-gray-600">{title}</p>
			<p class="text-2xl font-bold text-gray-900">{value}</p>
		</div>
		<div class="w-12 h-12 {colorClasses[color]} rounded-lg flex items-center justify-center">
			{@html icon}
		</div>
	</div>
</div>
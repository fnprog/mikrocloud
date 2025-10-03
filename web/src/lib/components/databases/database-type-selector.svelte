<script lang="ts">
	import type { DatabaseType } from '$lib/api/databases';

	interface Props {
		selected?: DatabaseType;
		onSelect: (type: DatabaseType) => void;
	}

	let { selected = $bindable(), onSelect }: Props = $props();

	const databaseTypes: Array<{
		type: DatabaseType;
		label: string;
		icon: string;
	}> = [
		{ type: 'postgresql', label: 'PostgreSQL', icon: '🐘' },
		{ type: 'mysql', label: 'MySQL', icon: '🐬' },
		{ type: 'mariadb', label: 'MariaDB', icon: '🦭' },
		{ type: 'redis', label: 'Redis', icon: '📦' },
		{ type: 'keydb', label: 'KeyDB', icon: '🔑' },
		{ type: 'dragonfly', label: 'Dragonfly', icon: '🔷' },
		{ type: 'mongodb', label: 'MongoDB', icon: '🍃' },
		{ type: 'clickhouse', label: 'ClickHouse', icon: '⚡' }
	];
</script>

<div class="grid grid-cols-3 gap-4">
	{#each databaseTypes as db}
		<button
			type="button"
			class="flex flex-col items-center gap-3 rounded-lg border-2 bg-card p-6 transition-all hover:border-primary/50 {selected ===
			db.type
				? 'border-primary'
				: 'border-border'}"
			onclick={() => {
				selected = db.type;
				onSelect(db.type);
			}}
		>
			<div class="flex h-12 w-12 items-center justify-center text-4xl">
				{db.icon}
			</div>
			<span class="text-sm font-medium">{db.label}</span>
			{#if selected === db.type}
				<div
					class="flex h-5 w-5 items-center justify-center rounded-full bg-primary text-primary-foreground"
				>
					<svg class="h-3 w-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"
						></path>
					</svg>
				</div>
			{/if}
		</button>
	{/each}
</div>

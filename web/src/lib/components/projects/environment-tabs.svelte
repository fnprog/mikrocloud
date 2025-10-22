<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Plus } from 'lucide-svelte';
	import type { Environment } from '$lib/api';

	interface Props {
		environments: Environment[];
		selectedEnvironmentId?: string;
		onSelect: (environmentId: string | undefined) => void;
		onAdd: () => void;
		counts?: Record<string, number>;
	}

	let {
		environments,
		selectedEnvironmentId = $bindable(),
		onSelect,
		onAdd,
		counts = {}
	}: Props = $props();

	function handleTabClick(envId: string | undefined) {
		selectedEnvironmentId = envId;
		onSelect(envId);
	}
</script>

<div class="flex items-center gap-2">
	<div class="flex items-center justify-stretch gap-1 overflow-x-auto">
		{#each environments as env (env.id)}
			<Button
				variant={selectedEnvironmentId === env.id ? 'secondary' : 'ghost'}
				class="transition-colors hover:text-foreground {selectedEnvironmentId === env.id
					? 'font-semibold'
					: 'text-muted-foreground'}"
				onclick={() => handleTabClick(env.id)}
			>
				{env.name}
				{#if counts[env.id] !== undefined}
					<Badge class="ml-2" variant={selectedEnvironmentId === env.id ? 'default' : 'secondary'}
						>{counts[env.id]}</Badge
					>
				{/if}
			</Button>
		{/each}
		<Button size="sm" variant="outline" onclick={onAdd} class="border-dashed">
			<Plus class="mr-1 h-4 w-4" />
			Add
		</Button>
	</div>
</div>

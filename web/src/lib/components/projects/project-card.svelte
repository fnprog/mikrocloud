<script lang="ts">
	import type { Project } from '$lib/features/projects/types';
	import { Card } from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import {
		DropdownMenu,
		DropdownMenuContent,
		DropdownMenuItem,
		DropdownMenuTrigger
	} from '$lib/components/ui/dropdown-menu';
	import { EllipsisVertical, Trash2, Settings } from 'lucide-svelte';
	import { IconStackFront } from '@tabler/icons-svelte';
	import { formatTimeAgo } from '$lib/utils/dates';

	interface Props {
		project: Project;
		onDelete?: (id: string) => void;
		onEdit?: (id: string) => void;
		onclick?: () => void;
	}

	let { project, onDelete, onEdit, onclick }: Props = $props();
</script>

<Card
	class="group relative flex flex-col gap-10 p-6 transition-shadow hover:shadow-md {onclick
		? 'cursor-pointer'
		: ''}"
	{onclick}
>
	<div class="flex items-start justify-between">
		<div class="flex items-center gap-3">
			<div class="bg-primary/10 text-primary flex size-12 items-center justify-center rounded-lg">
				<IconStackFront class="size-6" />
			</div>
			<div>
				<h3 class="font-semibold text-lg">{project.name}</h3>
				{#if project.description}
					<p class="text-muted-foreground text-sm">{project.description}</p>
				{/if}
			</div>
		</div>

		<DropdownMenu>
			<DropdownMenuTrigger>
				<Button variant="ghost" size="icon" class="size-8">
					<EllipsisVertical class="size-4" />
					<span class="sr-only">Open menu</span>
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent align="end">
				<DropdownMenuItem onclick={() => onEdit?.(project.id)}>
					<Settings class="mr-2 size-4" />
					<span>Settings</span>
				</DropdownMenuItem>
				<DropdownMenuItem onclick={() => onDelete?.(project.id)} class="text-destructive">
					<Trash2 class="mr-2 size-4" />
					<span>Delete</span>
				</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
	</div>

	<div class="flex items-center gap-2 text-sm">
		<span class="text-muted-foreground">Last deploy:</span>
		<span>{formatTimeAgo(project.created_at)}</span>
	</div>
</Card>

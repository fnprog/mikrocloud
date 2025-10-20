<script lang="ts">
	import {
		Activity,
		Database as DatabaseIcon,
		Globe,
		Layers,
		Settings,
		Terminal
	} from 'lucide-svelte';
	import { Position, Handle, type NodeProps } from '@xyflow/svelte';
	import { Button } from '../button';
	import { Badge } from '../badge';

	interface NodeData {
		id: string;
		name: string;
		description?: string;
		status: string;
		nodeType: 'application' | 'database';
		domain?: string;
		environmentId?: string;
		environmentName?: string;
		dbType?: string;
		connectionString?: string;
		onViewLogs?: (id: string) => void;
		onOpenTerminal?: (id: string) => void;
		onOpenSettings?: (id: string) => void;
	}

	let { data }: NodeProps<NodeData> = $props();
</script>

<div
	class="bg-white border-2 rounded-lg p-4 shadow-md min-w-[300px] max-w-[350px] {data.nodeType ===
	'application'
		? 'border-blue-500'
		: 'border-green-500'}"
>
	<Handle type="target" position={Position.Left} class="!bg-gray-400" />
	<Handle type="source" position={Position.Right} class="!bg-gray-400" />

	<div class="flex items-start justify-between mb-3">
		<div class="flex items-center gap-2 flex-1">
			{#if data.nodeType === 'application'}
				<div class="p-2 bg-blue-100 rounded-lg">
					<Layers class="size-5 text-blue-600" />
				</div>
			{:else}
				<div class="p-2 bg-green-100 rounded-lg">
					<DatabaseIcon class="size-5 text-green-600" />
				</div>
			{/if}
			<div class="flex-1 min-w-0">
				<h3 class="font-semibold text-sm truncate">{data.name}</h3>
				<p class="text-xs text-muted-foreground truncate">{data.description || 'No description'}</p>
			</div>
		</div>
		<Badge variant={data.status === 'running' ? 'default' : 'secondary'} class="text-xs">
			{data.status}
		</Badge>
	</div>

	<div class="space-y-2 text-xs">
		{#if data.nodeType === 'application'}
			{#if data.domain}
				<div class="flex items-center gap-2 p-2 bg-gray-50 rounded">
					<Globe class="size-3 text-gray-600 shrink-0" />
					<a
						href="https://{data.domain}"
						target="_blank"
						rel="noopener noreferrer"
						class="text-blue-600 hover:underline truncate"
					>
						{data.domain}
					</a>
				</div>
			{/if}
			{#if data.environmentId}
				<div class="flex items-center justify-between p-2 bg-gray-50 rounded">
					<span class="text-gray-600">Environment</span>
					<code class="font-mono text-xs bg-gray-200 px-1.5 py-0.5 rounded"
						>{data.environmentName || data.environmentId}</code
					>
				</div>
			{/if}
		{:else}
			<div class="flex items-center justify-between p-2 bg-gray-50 rounded">
				<span class="text-gray-600">Type</span>
				<Badge variant="outline" class="text-xs uppercase">{data.dbType}</Badge>
			</div>
			{#if data.connectionString}
				<div class="p-2 bg-gray-50 rounded">
					<span class="text-gray-600 block mb-1">Connection</span>
					<code class="font-mono text-[10px] break-all text-gray-800">
						{data.connectionString.substring(0, 50)}...
					</code>
				</div>
			{/if}
		{/if}
	</div>

	<div class="flex gap-1.5 mt-3 pt-3 border-t">
		{#if data.nodeType === 'application'}
			<Button
				variant="ghost"
				size="sm"
				class="flex-1 h-8 text-xs"
				onclick={() => data.onViewLogs?.(data.id)}
			>
				<Activity class="size-3" />
			</Button>
			<Button
				variant="ghost"
				size="sm"
				class="flex-1 h-8 text-xs"
				onclick={() => data.onOpenTerminal?.(data.id)}
			>
				<Terminal class="size-3" />
			</Button>
		{/if}
		<Button
			variant="ghost"
			size="sm"
			class="flex-1 h-8 text-xs"
			onclick={() => data.onOpenSettings?.(data.id)}
		>
			<Settings class="size-3" />
		</Button>
	</div>
</div>

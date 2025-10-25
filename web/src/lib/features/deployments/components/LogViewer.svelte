<script lang="ts">
	import type { StructuredLog, LogLevel } from '../types';

	interface Props {
		logs: StructuredLog[];
		isStreaming?: boolean;
		class?: string;
	}

	let { logs, isStreaming = false, class: className }: Props = $props();

	let selectedFilter = $state<LogLevel | 'all'>('all');
	let logsContainer: HTMLDivElement | null = null;

	const logCounts = $derived({
		all: logs.length,
		error: logs.filter((log) => log.level === 'error').length,
		warn: logs.filter((log) => log.level === 'warn').length,
		success: logs.filter((log) => log.level === 'success').length,
		info: logs.filter((log) => log.level === 'info').length
	});

	const filteredLogs = $derived(
		selectedFilter === 'all' ? logs : logs.filter((log) => log.level === selectedFilter)
	);

	function getLevelColor(level: LogLevel): string {
		switch (level) {
			case 'error':
				return 'text-red-500';
			case 'warn':
				return 'text-orange-500';
			case 'success':
				return 'text-green-500';
			default:
				return 'text-muted-foreground';
		}
	}

	function getLevelBadge(level: LogLevel): string {
		switch (level) {
			case 'error':
				return 'ERR';
			case 'warn':
				return 'WRN';
			case 'success':
				return 'SUC';
			default:
				return 'INF';
		}
	}

	$effect(() => {
		if (isStreaming && logsContainer) {
			logsContainer.scrollTop = logsContainer.scrollHeight;
		}
	});
</script>

<div class={className}>
	<div class="flex items-center gap-2 py-4 border-b">
		<div class="flex gap-2 px-3">
			<button
				onclick={() => (selectedFilter = 'all')}
				class="px-3 py-1 rounded-md text-sm transition-colors {selectedFilter === 'all'
					? 'bg-primary text-primary-foreground'
					: 'bg-muted hover:bg-muted/80'}"
			>
				All ({logCounts.all})
			</button>
			<button
				onclick={() => (selectedFilter = 'error')}
				class="px-3 py-1 rounded-md text-sm transition-colors {selectedFilter === 'error'
					? 'bg-red-500 text-white'
					: 'bg-muted hover:bg-muted/80'}"
			>
				Errors ({logCounts.error})
			</button>
			<button
				onclick={() => (selectedFilter = 'warn')}
				class="px-3 py-1 rounded-md text-sm transition-colors {selectedFilter === 'warn'
					? 'bg-orange-500 text-white'
					: 'bg-muted hover:bg-muted/80'}"
			>
				Warnings ({logCounts.warn})
			</button>
		</div>
	</div>

	<div
		bind:this={logsContainer}
		class="relative r overflow-x-auto overflow-y-auto max-h-[600px] font-mono text-sm px-3"
	>
		<!-- <div -->
		<!-- 	class="sticky top-0 left-0 right-0 h-8 bg-linear-to-b from-muted to-transparent pointer-events-none" -->
		<!-- ></div> -->

		{#if filteredLogs.length > 0}
			<div class="space-y-0.5 -mt-8 pt-8">
				{#each filteredLogs as log (log)}
					<div class="flex gap-3 hover:bg-background/50 px-2 py-0.5 rounded">
						<span class="text-muted-foreground/70 shrink-0 select-none">
							{log.timestamp}
						</span>
						<span class="{getLevelColor(log.level)} font-bold shrink-0 select-none w-8">
							{getLevelBadge(log.level)}
						</span>
						<span class="whitespace-pre-wrap break-all">{log.message}</span>
					</div>
				{/each}
			</div>
		{:else}
			<p class="text-muted-foreground text-center py-8">No logs available...</p>
		{/if}
	</div>
</div>

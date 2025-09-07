<script>
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import {
		Search,
		HelpCircle,
		ChevronDown,
		ArrowLeft,
		Download,
		Pause,
		Play,
		RotateCcw
	} from 'lucide-svelte';

	let isStreaming = $state(true);
	let selectedLevel = $state('all');

	// Mock log data
	let logs = $state([
		{
			id: 1,
			timestamp: '2024-01-31 14:23:45',
			level: 'info',
			message: 'Server started on port 3000',
			source: 'app'
		},
		{
			id: 2,
			timestamp: '2024-01-31 14:23:46',
			level: 'info',
			message: 'Database connection established',
			source: 'database'
		},
		{
			id: 3,
			timestamp: '2024-01-31 14:24:12',
			level: 'warn',
			message: 'Deprecated API endpoint accessed: /api/v1/users',
			source: 'api'
		},
		{
			id: 4,
			timestamp: '2024-01-31 14:24:15',
			level: 'error',
			message: 'Failed to process payment for user 12345',
			source: 'payment'
		},
		{
			id: 5,
			timestamp: '2024-01-31 14:24:18',
			level: 'info',
			message: 'User authentication successful for user@example.com',
			source: 'auth'
		},
		{
			id: 6,
			timestamp: '2024-01-31 14:24:22',
			level: 'debug',
			message: 'Cache hit for key: user_profile_12345',
			source: 'cache'
		}
	]);

	function setActiveSection(section) {
		activeSection = section;
		if (section === 'Overview') {
			goto(`/project/${projectId}`);
		} else {
			goto(`/project/${projectId}/${section.toLowerCase().replace(' ', '-')}`);
		}
	}

	function getLevelColor(level) {
		switch (level) {
			case 'error':
				return 'text-red-600 bg-red-50';
			case 'warn':
				return 'text-yellow-600 bg-yellow-50';
			case 'info':
				return 'text-blue-600 bg-blue-50';
			case 'debug':
				return 'text-gray-600 bg-gray-50';
			default:
				return 'text-gray-600 bg-gray-50';
		}
	}

	function toggleStreaming() {
		isStreaming = !isStreaming;
	}

	function downloadLogs() {
		console.log('Downloading logs...');
	}

	function clearLogs() {
		logs = [];
	}

	const filteredLogs = $derived.by(() => {
		if (selectedLevel === 'all') return logs;
		return logs.filter((log) => log.level === selectedLevel);
	});
</script>

<svelte:head>
	<!-- <title>Logs - {project.name}</title> -->
</svelte:head>

<div class="flex-1 p-6">
	<!-- Header -->
	<div class="flex items-center justify-between mb-6">
		<div class="flex items-center space-x-4">
			<!-- <Button variant="ghost" size="sm" onclick={() => goto(`/project/${projectId}`)}> -->
			<!-- 	<ArrowLeft class="w-4 h-4" /> -->
			<!-- </Button> -->
			<div>
				<h1 class="text-2xl font-semibold text-gray-900">Logs</h1>
				<!-- <p class="text-sm text-gray-500 mt-1">Real-time application logs for {project.name}.</p> -->
			</div>
		</div>
		<div class="flex items-center space-x-3">
			<!-- Log Level Filter -->
			<select
				bind:value={selectedLevel}
				class="px-3 py-2 border border-gray-300 rounded-md text-sm"
			>
				<option value="all">All Levels</option>
				<option value="error">Error</option>
				<option value="warn">Warning</option>
				<option value="info">Info</option>
				<option value="debug">Debug</option>
			</select>

			<Button variant="outline" onclick={downloadLogs}>
				<Download class="w-4 h-4 mr-2" />
				Download
			</Button>
			<Button variant="outline" onclick={toggleStreaming}>
				{#if isStreaming}
					<Pause class="w-4 h-4 mr-2" />
					Pause
				{:else}
					<Play class="w-4 h-4 mr-2" />
					Resume
				{/if}
			</Button>
			<Button variant="outline" onclick={clearLogs}>
				<RotateCcw class="w-4 h-4 mr-2" />
				Clear
			</Button>
		</div>
	</div>

	<!-- Logs Container -->
	<Card class="h-[calc(100vh-200px)]">
		<CardContent class="p-0 h-full">
			<div class="h-full overflow-auto bg-gray-900 text-green-400 font-mono text-sm">
				<div class="p-4 space-y-1">
					{#each filteredLogs as log (log.id)}
						<div class="flex items-start space-x-4 py-1 hover:bg-gray-800 px-2 rounded">
							<span class="text-gray-500 text-xs w-32 flex-shrink-0 mt-0.5">
								{log.timestamp}
							</span>
							<Badge
								variant="outline"
								class="text-xs {getLevelColor(log.level)} border-0 w-16 justify-center"
							>
								{log.level.toUpperCase()}
							</Badge>
							<span class="text-blue-400 text-xs w-20 flex-shrink-0 mt-0.5">
								[{log.source}]
							</span>
							<span class="text-gray-300 flex-1">
								{log.message}
							</span>
						</div>
					{/each}

					{#if isStreaming}
						<div class="flex items-center space-x-2 py-2 px-2">
							<div class="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
							<span class="text-gray-500 text-xs">Live streaming...</span>
						</div>
					{/if}
				</div>
			</div>
		</CardContent>
	</Card>
</div>

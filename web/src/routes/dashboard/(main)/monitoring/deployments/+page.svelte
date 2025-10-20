<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import { Search, RefreshCw, Clock, GitBranch, User, Calendar, TrendingUp, TrendingDown } from 'lucide-svelte';

	let searchQuery = $state('');
	let statusFilter = $state('all');
	let projectFilter = $state('all');
	let isRefreshing = $state(false);

	const statusFilters = [
		{ value: 'all', label: 'All Statuses' },
		{ value: 'success', label: 'Success' },
		{ value: 'failed', label: 'Failed' },
		{ value: 'pending', label: 'Pending' },
		{ value: 'running', label: 'Running' }
	];

	const projectFilters = [
		{ value: 'all', label: 'All Projects' },
		{ value: 'mikrocloud-app', label: 'mikrocloud-app' },
		{ value: 'web-dashboard', label: 'web-dashboard' },
		{ value: 'api-service', label: 'api-service' }
	];

	const deployments = [
		{
			id: 'dep_001',
			project: 'mikrocloud-app',
			environment: 'production',
			status: 'success',
			version: 'v2.4.1',
			branch: 'main',
			commit: 'a1b2c3d',
			commitMessage: 'feat: Add monitoring dashboard',
			author: 'John Doe',
			startedAt: '2024-10-10T14:25:00Z',
			completedAt: '2024-10-10T14:27:30Z',
			duration: 150,
			buildTime: 85,
			deployTime: 65
		},
		{
			id: 'dep_002',
			project: 'web-dashboard',
			environment: 'staging',
			status: 'running',
			version: 'v1.8.3',
			branch: 'develop',
			commit: 'e4f5g6h',
			commitMessage: 'fix: Update dependencies',
			author: 'Jane Smith',
			startedAt: '2024-10-10T14:30:00Z',
			completedAt: null,
			duration: null,
			buildTime: null,
			deployTime: null
		},
		{
			id: 'dep_003',
			project: 'api-service',
			environment: 'production',
			status: 'failed',
			version: 'v3.1.0',
			branch: 'main',
			commit: 'i7j8k9l',
			commitMessage: 'refactor: Database connection pool',
			author: 'Bob Johnson',
			startedAt: '2024-10-10T14:15:00Z',
			completedAt: '2024-10-10T14:16:45Z',
			duration: 105,
			buildTime: 90,
			deployTime: 15
		},
		{
			id: 'dep_004',
			project: 'mikrocloud-app',
			environment: 'staging',
			status: 'success',
			version: 'v2.4.0',
			branch: 'develop',
			commit: 'm0n1p2q',
			commitMessage: 'chore: Update logging configuration',
			author: 'Alice Williams',
			startedAt: '2024-10-10T13:45:00Z',
			completedAt: '2024-10-10T13:48:20Z',
			duration: 200,
			buildTime: 120,
			deployTime: 80
		},
		{
			id: 'dep_005',
			project: 'web-dashboard',
			environment: 'production',
			status: 'success',
			version: 'v1.8.2',
			branch: 'main',
			commit: 'r3s4t5u',
			commitMessage: 'feat: Add new analytics widgets',
			author: 'Charlie Brown',
			startedAt: '2024-10-10T13:20:00Z',
			completedAt: '2024-10-10T13:22:45Z',
			duration: 165,
			buildTime: 95,
			deployTime: 70
		},
		{
			id: 'dep_006',
			project: 'api-service',
			environment: 'production',
			status: 'success',
			version: 'v3.0.9',
			branch: 'main',
			commit: 'v6w7x8y',
			commitMessage: 'fix: Handle edge case in authentication',
			author: 'Diana Prince',
			startedAt: '2024-10-10T12:50:00Z',
			completedAt: '2024-10-10T12:52:20Z',
			duration: 140,
			buildTime: 80,
			deployTime: 60
		}
	];

	const filteredDeployments = $derived(
		deployments.filter((dep) => {
			const matchesSearch =
				dep.project.toLowerCase().includes(searchQuery.toLowerCase()) ||
				dep.commitMessage.toLowerCase().includes(searchQuery.toLowerCase()) ||
				dep.author.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesStatus = statusFilter === 'all' || dep.status === statusFilter;
			const matchesProject = projectFilter === 'all' || dep.project === projectFilter;
			return matchesSearch && matchesStatus && matchesProject;
		})
	);

	const stats = $derived({
		total: deployments.length,
		successful: deployments.filter((d) => d.status === 'success').length,
		failed: deployments.filter((d) => d.status === 'failed').length,
		running: deployments.filter((d) => d.status === 'running').length,
		avgDuration: Math.round(
			deployments.filter((d) => d.duration).reduce((sum, d) => sum + (d.duration || 0), 0) /
				deployments.filter((d) => d.duration).length
		),
		successRate: Math.round(
			(deployments.filter((d) => d.status === 'success').length / deployments.length) * 100
		)
	});

	function getStatusBadgeVariant(status: string): 'default' | 'secondary' | 'destructive' {
		switch (status) {
			case 'success':
				return 'default';
			case 'failed':
				return 'destructive';
			case 'running':
			case 'pending':
				return 'secondary';
			default:
				return 'secondary';
		}
	}

	function getStatusColor(status: string) {
		switch (status) {
			case 'success':
				return 'text-green-500';
			case 'failed':
				return 'text-red-500';
			case 'running':
				return 'text-blue-500';
			case 'pending':
				return 'text-yellow-500';
			default:
				return 'text-gray-500';
		}
	}

	function formatDuration(seconds: number | null) {
		if (seconds === null) return '-';
		const minutes = Math.floor(seconds / 60);
		const secs = seconds % 60;
		return `${minutes}m ${secs}s`;
	}

	function formatTime(timestamp: string) {
		const date = new Date(timestamp);
		return date.toLocaleTimeString();
	}

	function formatDate(timestamp: string) {
		const date = new Date(timestamp);
		return date.toLocaleDateString();
	}

	async function handleRefresh() {
		isRefreshing = true;
		await new Promise((resolve) => setTimeout(resolve, 1000));
		isRefreshing = false;
	}
</script>

<div class="flex-1 overflow-y-auto">
	<div class="border-b border-border bg-background sticky top-0 z-10">
		<div class="px-8 py-6">
			<div class="flex items-center justify-between mb-4">
				<div>
					<h1 class="text-3xl font-bold">Deployments</h1>
					<p class="text-muted-foreground mt-1">Track deployment history, timings, and success rates</p>
				</div>
			<Button variant="outline" size="icon" onclick={handleRefresh} disabled={isRefreshing}>
				<RefreshCw class={`h-4 w-4 ${isRefreshing ? 'animate-spin' : ''}`} />
			</Button>
			</div>

			<div class="grid grid-cols-5 gap-4 mb-4">
				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Total Deployments</CardDescription>
						<CardTitle class="text-2xl">{stats.total}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Successful</CardDescription>
						<CardTitle class="text-2xl text-green-500">{stats.successful}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Failed</CardDescription>
						<CardTitle class="text-2xl text-red-500">{stats.failed}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Avg Duration</CardDescription>
						<CardTitle class="text-2xl">{formatDuration(stats.avgDuration)}</CardTitle>
					</CardHeader>
				</Card>

				<Card>
					<CardHeader class="pb-2">
						<CardDescription class="text-xs">Success Rate</CardDescription>
						<CardTitle class="text-2xl flex items-center gap-2">
							{stats.successRate}%
							{#if stats.successRate >= 90}
								<TrendingUp class="h-5 w-5 text-green-500" />
							{:else}
								<TrendingDown class="h-5 w-5 text-yellow-500" />
							{/if}
						</CardTitle>
					</CardHeader>
				</Card>
			</div>

			<div class="grid grid-cols-3 gap-3">
				<div class="relative">
					<Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
					<Input
						type="text"
						placeholder="Search deployments..."
						class="pl-10"
						bind:value={searchQuery}
					/>
				</div>

				<Select.Root
					selected={{ value: statusFilter, label: statusFilter }}
					onSelectedChange={(v) => v && (statusFilter = v.value)}
				>
					<Select.Trigger>
						<Select.Value placeholder="Filter by status" />
					</Select.Trigger>
					<Select.Content>
						{#each statusFilters as filter}
							<Select.Item value={filter.value}>{filter.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>

				<Select.Root
					selected={{ value: projectFilter, label: projectFilter }}
					onSelectedChange={(v) => v && (projectFilter = v.value)}
				>
					<Select.Trigger>
						<Select.Value placeholder="Filter by project" />
					</Select.Trigger>
					<Select.Content>
						{#each projectFilters as filter}
							<Select.Item value={filter.value}>{filter.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		</div>
	</div>

	<div class="p-8 space-y-4">
		<div class="text-sm text-muted-foreground mb-2">
			Showing {filteredDeployments.length} of {deployments.length} deployments
		</div>

		{#each filteredDeployments as deployment (deployment.id)}
			<Card>
				<CardHeader>
					<div class="flex items-start justify-between">
						<div class="flex-1">
							<div class="flex items-center gap-3 mb-2">
								<CardTitle class="text-lg">{deployment.project}</CardTitle>
								<Badge variant="outline">{deployment.environment}</Badge>
								<Badge variant={getStatusBadgeVariant(deployment.status)}>
									{deployment.status}
								</Badge>
								<span class="text-sm text-muted-foreground">{deployment.version}</span>
							</div>
							<CardDescription class="flex items-center gap-4 text-xs">
								<span class="flex items-center gap-1">
									<Calendar class="h-3 w-3" />
									{formatDate(deployment.startedAt)} at {formatTime(deployment.startedAt)}
								</span>
								<span class="flex items-center gap-1">
									<User class="h-3 w-3" />
									{deployment.author}
								</span>
								<span class="flex items-center gap-1">
									<GitBranch class="h-3 w-3" />
									{deployment.branch}
								</span>
							</CardDescription>
						</div>
					</div>
				</CardHeader>
				<CardContent>
					<div class="grid gap-4 lg:grid-cols-2">
						<div class="space-y-3">
							<div class="p-3 bg-muted rounded-lg">
								<div class="text-xs text-muted-foreground mb-1">Commit</div>
								<div class="font-mono text-sm font-medium">{deployment.commit}</div>
								<div class="text-sm mt-1">{deployment.commitMessage}</div>
							</div>

							{#if deployment.status === 'running'}
								<div class="p-3 bg-blue-500/10 border border-blue-500/20 rounded-lg">
									<div class="flex items-center gap-2">
										<div class="h-2 w-2 bg-blue-500 rounded-full animate-pulse"></div>
										<span class="text-sm {getStatusColor(deployment.status)}">
											Deployment in progress...
										</span>
									</div>
								</div>
							{/if}
						</div>

						<div class="grid grid-cols-3 gap-3">
							<div class="border rounded-lg p-3">
								<div class="text-xs text-muted-foreground mb-1">Total Duration</div>
								<div class="text-lg font-bold flex items-center gap-1">
									<Clock class="h-4 w-4" />
									{formatDuration(deployment.duration)}
								</div>
							</div>

							<div class="border rounded-lg p-3">
								<div class="text-xs text-muted-foreground mb-1">Build Time</div>
								<div class="text-lg font-bold">
									{formatDuration(deployment.buildTime)}
								</div>
							</div>

							<div class="border rounded-lg p-3">
								<div class="text-xs text-muted-foreground mb-1">Deploy Time</div>
								<div class="text-lg font-bold">
									{formatDuration(deployment.deployTime)}
								</div>
							</div>
						</div>
					</div>

					{#if deployment.status === 'success' || deployment.status === 'failed'}
						<div class="mt-4 pt-4 border-t flex items-center justify-between">
							<div class="text-sm text-muted-foreground">
								{#if deployment.completedAt}
									Completed at {formatTime(deployment.completedAt)}
								{/if}
							</div>
							<div class="flex items-center gap-2">
								<Button variant="outline" size="sm">View Logs</Button>
								<Button variant="outline" size="sm">Rollback</Button>
							</div>
						</div>
					{/if}
				</CardContent>
			</Card>
		{/each}

		{#if filteredDeployments.length === 0}
			<Card>
				<CardContent class="py-12">
					<div class="text-center text-muted-foreground">
						No deployments found matching your criteria
					</div>
				</CardContent>
			</Card>
		{/if}
	</div>
</div>

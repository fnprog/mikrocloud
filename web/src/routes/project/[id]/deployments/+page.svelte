<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent } from '$lib/components/ui/card';
	import {
		Github,
		CheckCircle,
		GitBranch,
		Clock,
		AlertCircle,
		XCircle,
		Eye,
		RotateCcw
	} from 'lucide-svelte';

	const projectId = page.params.id;

	// Mock project data
	let project = $state({
		name: 'Focalpoint Dashboard',
		repo: 'focalpoint/dashboard',
		branch: 'main',
		workspace: 'Focalpoint',
		category: 'Applications',
		deploymentEnabled: true,
		domain: 'app.focalpoint.com'
	});

	// Mock deployments data
	let deployments = $state([
		{
			id: 1,
			status: 'success',
			message: 'Deployment successful',
			time: '2 minutes ago',
			duration: '1m 23s',
			commit: 'feat: add user dashboard',
			commitHash: 'a1b2c3d',
			branch: 'main',
			author: 'John Doe',
			buildLogs: 'Build completed successfully...'
		},
		{
			id: 2,
			status: 'success',
			message: 'Deployment successful',
			time: '1 hour ago',
			duration: '2m 15s',
			commit: 'fix: resolve API timeout issue',
			commitHash: 'e4f5g6h',
			branch: 'main',
			author: 'Jane Smith',
			buildLogs: 'Build completed successfully...'
		},
		{
			id: 3,
			status: 'failed',
			message: 'Build failed',
			time: '3 hours ago',
			duration: '45s',
			commit: 'refactor: update API endpoints',
			commitHash: 'i7j8k9l',
			branch: 'develop',
			author: 'Bob Johnson',
			buildLogs: 'Error: Module not found...'
		},
		{
			id: 4,
			status: 'cancelled',
			message: 'Deployment cancelled',
			time: '1 day ago',
			duration: '30s',
			commit: 'chore: update dependencies',
			commitHash: 'm0n1o2p',
			branch: 'main',
			author: 'Alice Brown',
			buildLogs: 'Build cancelled by user...'
		},
		{
			id: 5,
			status: 'success',
			message: 'Deployment successful',
			time: '2 days ago',
			duration: '1m 45s',
			commit: 'feat: implement user authentication',
			commitHash: 'q3r4s5t',
			branch: 'main',
			author: 'John Doe',
			buildLogs: 'Build completed successfully...'
		}
	]);

	function getStatusIcon(status: string) {
		switch (status) {
			case 'success':
				return CheckCircle;
			case 'failed':
				return AlertCircle;
			case 'cancelled':
				return XCircle;
			case 'running':
				return Clock;
			default:
				return Clock;
		}
	}

	function getStatusColor(status: string) {
		switch (status) {
			case 'success':
				return 'text-green-500';
			case 'failed':
				return 'text-red-500';
			case 'cancelled':
				return 'text-gray-500';
			case 'running':
				return 'text-blue-500';
			default:
				return 'text-gray-500';
		}
	}

	function getStatusBadgeVariant(status: string) {
		switch (status) {
			case 'success':
				return 'default';
			case 'failed':
				return 'destructive';
			case 'cancelled':
				return 'secondary';
			case 'running':
				return 'default';
			default:
				return 'secondary';
		}
	}

	function redeployCommit(deploymentId: number) {
		console.log('Redeploying commit:', deploymentId);
	}

	function viewLogs(deploymentId: number) {
		console.log('Viewing logs for:', deploymentId);
	}
</script>

<svelte:head>
	<title>Deployments - {project.name}</title>
</svelte:head>

<!-- Main Content -->
<div class="flex-1 p-6">
	<!-- Header -->
	<div class="flex items-center justify-between mb-6">
		<div class="flex items-center space-x-4">
			<div>
				<h1 class="text-2xl font-semibold text-gray-900">Deployments</h1>
				<p class="text-sm text-gray-500 mt-1">
					View and manage all deployments for {project.name}.
				</p>
			</div>
		</div>
		<Button class="bg-gray-900 hover:bg-gray-800">Deploy now</Button>
	</div>

	<!-- Deployments List -->
	<div class="space-y-4">
		{#each deployments as deployment (deployment.id)}
			<Card class="hover:shadow-sm transition-shadow">
				<CardContent class="p-6">
					<div class="flex items-center justify-between">
						<div class="flex items-center space-x-4">
							{#if deployment.status === 'success'}
								<CheckCircle class="w-5 h-5 text-green-500" />
							{:else if deployment.status === 'failed'}
								<AlertCircle class="w-5 h-5 text-red-500" />
							{:else if deployment.status === 'cancelled'}
								<XCircle class="w-5 h-5 text-gray-500" />
							{:else}
								<Clock class="w-5 h-5 text-blue-500" />
							{/if}

							<div>
								<div class="flex items-center space-x-3 mb-1">
									<h3 class="font-medium text-gray-900">{deployment.message}</h3>
									<Badge variant={getStatusBadgeVariant(deployment.status)} class="text-xs">
										{deployment.status}
									</Badge>
								</div>
								<div class="flex items-center space-x-4 text-sm text-gray-600">
									<div class="flex items-center space-x-1">
										<Github class="w-4 h-4" />
										<span>{deployment.commitHash}</span>
									</div>
									<div class="flex items-center space-x-1">
										<GitBranch class="w-4 h-4" />
										<span>{deployment.branch}</span>
									</div>
									<span>by {deployment.author}</span>
									<div class="flex items-center space-x-1">
										<Clock class="w-4 h-4" />
										<span>{deployment.duration}</span>
									</div>
								</div>
								<p class="text-sm text-gray-700 mt-1">{deployment.commit}</p>
							</div>
						</div>
						<div class="flex items-center space-x-4">
							<span class="text-sm text-gray-500">{deployment.time}</span>
							<div class="flex items-center space-x-2">
								<Button size="sm" variant="outline" onclick={() => viewLogs(deployment.id)}>
									<Eye class="w-4 h-4 mr-1" />
									Logs
								</Button>
								{#if deployment.status === 'success'}
									<Button size="sm" variant="outline" onclick={() => redeployCommit(deployment.id)}>
										<RotateCcw class="w-4 h-4 mr-1" />
										Redeploy
									</Button>
								{/if}
							</div>
						</div>
					</div>
				</CardContent>
			</Card>
		{/each}
	</div>
</div>

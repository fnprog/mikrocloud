<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Plus, Ellipsis, Database, Grid3x3, List, ExternalLink } from 'lucide-svelte';
	import { SiGithub } from '@icons-pack/svelte-simple-icons';

	// Mock data for projects
	let projects = $state([
		{
			id: 1,
			name: 'sveltekit-boilerplate',
			status: 'success',
			statusText: 'Production deployment successful 22h ago',
			services: ['github', 'postgres', 'redis'],
			lastDeploy: '22h ago'
		},
		{
			id: 2,
			name: 'expert-legs',
			status: 'deploying',
			statusText: 'Deployment in progress',
			services: ['github', 'postgres'],
			lastDeploy: '5m ago'
		},
		{
			id: 3,
			name: 'api-gateway',
			status: 'failed',
			statusText: 'Deployment failed 1h ago',
			services: ['github', 'mongodb'],
			lastDeploy: '1h ago'
		},
		{
			id: 4,
			name: 'frontend-app',
			status: 'paused',
			statusText: 'Project paused',
			services: ['github'],
			lastDeploy: '3d ago'
		}
	]);

	let viewMode = $state('grid'); // 'grid' or 'list'

	function toggleViewMode() {
		viewMode = viewMode === 'grid' ? 'list' : 'grid';
	}

	function getStatusColor(status: string) {
		switch (status) {
			case 'success':
				return 'bg-green-500';
			case 'deploying':
				return 'bg-yellow-500 animate-pulse';
			case 'failed':
				return 'bg-red-500';
			case 'paused':
				return 'bg-gray-400';
			default:
				return 'bg-gray-400';
		}
	}

	function getServiceIcon(service: string) {
		switch (service) {
			case 'github':
				return SiGithub;
			case 'postgres':
				return Database;
			case 'redis':
				return Database;
			case 'mongodb':
				return Database;
			default:
				return Database;
		}
	}

	function goToProject(projectId: number) {
		goto(`/project/${projectId}`);
	}
</script>

<svelte:head>
	<title>Dashboard - Projects Overview</title>
</svelte:head>

<!-- Main Content -->
<div class="flex-1 flex flex-col overflow-hidden">
	<!-- Header -->
	<div class="px-6 py-4">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-2xl font-semibold text-gray-900">Dashboard</h1>
				<p class="text-sm text-gray-500 mt-1">Overview of all your projects and their status.</p>
			</div>
			<div class="flex items-center space-x-3">
				<div class="flex items-center bg-gray-100 rounded-md p-1">
					<button
						onclick={toggleViewMode}
						class="p-1 rounded {viewMode === 'grid' ? 'bg-white shadow-sm' : ''}"
					>
						<Grid3x3 class="w-4 h-4" />
					</button>
					<button
						onclick={toggleViewMode}
						class="p-1 rounded {viewMode === 'list' ? 'bg-white shadow-sm' : ''}"
					>
						<List class="w-4 h-4" />
					</button>
				</div>
				<Button>
					<Plus class="w-4 h-4 mr-2" />
					New Project
				</Button>
			</div>
		</div>
	</div>

	<!-- Projects Grid/List -->
	<div class="flex-1 overflow-auto p-6">
		{#if viewMode === 'grid'}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each projects as project (project.id)}
					<Card
						class="cursor-pointer hover:shadow-md transition-shadow"
						onclick={() => goToProject(project.id)}
					>
						<CardHeader class="pb-3">
							<div class="flex items-center justify-between">
								<CardTitle class="text-lg">{project.name}</CardTitle>
								<Button variant="ghost" size="sm" onclick={(e) => e.stopPropagation()}>
									<Ellipsis class="w-4 h-4" />
								</Button>
							</div>
							<div class="flex items-center space-x-2">
								<div class="w-2 h-2 rounded-full {getStatusColor(project.status)}"></div>
								<span class="text-sm text-gray-600">{project.statusText}</span>
							</div>
						</CardHeader>
						<CardContent>
							<div class="flex items-center space-x-2 mb-3">
								{#each project.services as service (service)}
									{@const IconComponent = getServiceIcon(service)}
									<div class="w-6 h-6 bg-gray-100 rounded p-1">
										<IconComponent class="w-4 h-4 text-gray-600" />
									</div>
								{/each}
							</div>
							<div class="flex items-center justify-between text-sm text-gray-500">
								<span>Last deploy: {project.lastDeploy}</span>
								<ExternalLink class="w-4 h-4" />
							</div>
						</CardContent>
					</Card>
				{/each}
			</div>
		{:else}
			<div class="space-y-3">
				{#each projects as project (project.id)}
					<Card
						class="cursor-pointer hover:shadow-sm transition-shadow"
						onclick={() => goToProject(project.id)}
					>
						<CardContent class="p-4">
							<div class="flex items-center justify-between">
								<div class="flex items-center space-x-4">
									<div class="w-2 h-2 rounded-full {getStatusColor(project.status)}"></div>
									<div>
										<h3 class="font-medium text-gray-900">{project.name}</h3>
										<p class="text-sm text-gray-600">{project.statusText}</p>
									</div>
								</div>
								<div class="flex items-center space-x-4">
									<div class="flex items-center space-x-2">
										{#each project.services as service (service)}
											{@const IconComponent = getServiceIcon(service)}
											<div class="w-5 h-5 bg-gray-100 rounded p-0.5">
												<IconComponent class="w-4 h-4 text-gray-600" />
											</div>
										{/each}
									</div>
									<span class="text-sm text-gray-500">{project.lastDeploy}</span>
									<Button variant="ghost" size="sm" onclick={(e) => e.stopPropagation()}>
										<Ellipsis class="w-4 h-4" />
									</Button>
								</div>
							</div>
						</CardContent>
					</Card>
				{/each}
			</div>
		{/if}
	</div>
</div>

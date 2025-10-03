<script>
	import { Button } from '$lib/components/ui/button';
	import FocalpointNode from '$lib/components/FocalpointNode.svelte';

	import { ExternalLink, Github, CircleCheckBig, GitBranch } from 'lucide-svelte';

	import { SvelteFlow, Controls, Background, Position } from '@xyflow/svelte';

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

	// Custom node types
	const nodeTypes = {
		focalpoint: FocalpointNode
	};

	const nodeDefaults = {
		sourcePosition: Position.Right,
		targetPosition: Position.Left
	};

	// SvelteFlow nodes and edges - flowing right to left
	let nodes = $state.raw([
		{
			id: '1',
			type: 'default',
			position: { x: 50, y: 100 },
			data: {
				label: 'Client'
			},
			style:
				'background: #f3f4f6; color: #374151; border: 1px solid #d1d5db; border-radius: 8px; padding: 12px; width: 120px; text-align: center; font-weight: 500;',
			...nodeDefaults
		},
		{
			id: '2',
			type: 'default',
			position: { x: 300, y: 100 },
			data: {
				label: 'Cloudflare'
			},
			style:
				'background: #fef3e2; color: #92400e; border: 1px solid #fbbf24; border-radius: 8px; padding: 16px; width: 200px;',

			...nodeDefaults
		},
		{
			id: '3',
			type: 'focalpoint',
			position: { x: 650, y: 0 },
			data: {
				name: 'Focalpoint Dashboard',
				type: 'Web process',
				startCommand: 'npm run start',
				resources: 'S1 (0.5 CPU, 1 GB RAM)',
				instances: '1',
				domain: 'app.focalpoint.com'
			},

			...nodeDefaults
		}
	]);

	let edges = $state.raw([
		{
			id: 'e1-2',
			source: '1',
			target: '2',
			label: 'REQUESTS',
			style:
				'stroke: #9ca3af; stroke-width: 2; stroke-dasharray: 8,4; animation: dash 2s linear infinite;',
			labelStyle:
				'font-size: 10px; font-weight: 600; color: #6b7280; text-transform: uppercase; letter-spacing: 0.05em;'
		},
		{
			id: 'e2-3',
			source: '2',
			target: '3',
			label: 'PORT: 8080',
			style:
				'stroke: #fb923c; stroke-width: 2; stroke-dasharray: 8,4; animation: dash 2s linear infinite;',
			labelStyle:
				'font-size: 10px; font-weight: 600; color: #c2410c; background: #fed7aa; padding: 2px 8px; border-radius: 4px;'
		}
	]);

	// Mock deployments
	let deployments = $state([
		{
			id: 1,
			status: 'success',
			message: 'Deployment successful',
			time: '2 minutes ago',
			commit: 'feat: add user dashboard',
			branch: 'main'
		},
		{
			id: 2,
			status: 'success',
			message: 'Deployment successful',
			time: '1 hour ago',
			commit: 'fix: resolve API timeout',
			branch: 'main'
		}
	]);

	function visitApp() {
		window.open(`https://${project.domain}`, '_blank');
	}

	function deployNow() {
		// Mock deploy action
		console.log('Deploying...');
	}

	function viewAllDeployments() {
		// Navigate to deployments page
		console.log('View all deployments');
	}
</script>

<svelte:head>
	<title>{project.name} - Project Overview</title>
	<style>
		@keyframes dash {
			to {
				stroke-dashoffset: -12;
			}
		}

		.svelte-flow .react-flow__edge-path {
			animation: dash 2s linear infinite;
		}
	</style>
</svelte:head>

<div class="min-h-screen bg-gray-50">
	<div class="flex-1 flex">
		<!-- Center Content -->
		<div class="flex-1 p-6">
			<!-- Header -->
			<div class="flex items-center justify-between mb-6">
				<div>
					<h1 class="text-2xl font-semibold text-gray-900 mb-2">{project.name}</h1>
					<div class="flex items-center space-x-4">
						<div class="flex items-center space-x-2">
							<CircleCheckBig class="w-4 h-4 text-green-500" />
							<span class="text-sm text-green-600">Automatic deployment enabled</span>
						</div>
						<div class="flex items-center space-x-2 text-sm text-gray-600">
							<Github class="w-4 h-4" />
							<span>{project.repo}</span>
							<GitBranch class="w-4 h-4" />
							<span>{project.branch}</span>
							<ExternalLink class="w-4 h-4" />
						</div>
					</div>
				</div>
				<Button onclick={visitApp} class="bg-gray-900 hover:bg-gray-800">
					<ExternalLink class="w-4 h-4 mr-2" />
					Visit App
				</Button>
			</div>

			<!-- Architecture Diagram -->
			<div class="bg-white rounded-lg border border-gray-200 mb-6 relative" style="height: 400px;">
				<SvelteFlow
					bind:nodes
					bind:edges
					{nodeTypes}
					fitView
					class="bg-gray-50"
					nodesDraggable={false}
					nodesConnectable={false}
					elementsSelectable={false}
				>
					<Background bgColor="#e5e7eb" gap={20} />
					<Controls />
				</SvelteFlow>
			</div>

			<!-- Latest Deployments -->
			<div class="bg-white rounded-lg border border-gray-200 p-6">
				<div class="flex items-center justify-between mb-4">
					<h2 class="text-lg font-semibold text-gray-900">Latest deployments</h2>
					<div class="flex items-center space-x-3">
						<Button variant="outline" onclick={viewAllDeployments}>View all deployments</Button>
						<Button onclick={deployNow} class="bg-gray-900 hover:bg-gray-800">Deploy now</Button>
					</div>
				</div>

				<div class="space-y-3">
					{#each deployments as deployment (deployment.id)}
						<div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
							<div class="flex items-center space-x-3">
								<CircleCheckBig class="w-5 h-5 text-green-500" />
								<div>
									<div class="font-medium text-gray-900">{deployment.message}</div>
									<div class="text-sm text-gray-600">{deployment.commit}</div>
								</div>
							</div>
							<div class="text-sm text-gray-500">{deployment.time}</div>
						</div>
					{/each}
				</div>
			</div>
		</div>
	</div>
</div>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Ellipsis, Globe, Plus } from 'lucide-svelte';

	let projects = $state([
		{
			id: 1,
			name: 'Ultra mega cool project',
			lastDeploy: '2h ago',
			icons: ['🌐', '❤️', '👤']
		},
		{
			id: 2,
			name: 'Nice product',
			lastDeploy: '7d ago',
			icons: ['🌐', '⚙️']
		},
		{
			id: 3,
			name: 'Big Project',
			lastDeploy: '09/09/25',
			icons: ['🌐', '✅']
		},
		{
			id: 4,
			name: 'Todo List SaaS',
			lastDeploy: '08/04/2024',
			icons: ['🌐']
		},
		{
			id: 5,
			name: 'Tik Tak Toe',
			lastDeploy: '01/01/1999',
			icons: ['🌐']
		}
	]);

	let activities = $state([
		{
			id: 1,
			user: 'John',
			avatar: '👤',
			action: 'created a project',
			project: 'ultra mega cool project',
			time: '2h ago',
			type: 'create'
		},
		{
			id: 2,
			user: 'Mikrocloud',
			avatar: '🔧',
			action: 'Deployment b33423 failed',
			project: 'ultra mega cool project',
			time: '5h ago',
			type: 'failed'
		},
		{
			id: 3,
			user: 'John',
			avatar: '👤',
			action: 'Paused a service',
			project: 'Nice Product/web-server',
			time: '2d ago',
			type: 'pause'
		},
		{
			id: 4,
			user: 'John',
			avatar: '👤',
			action: 'created a project',
			project: 'Todo list saas',
			time: '23d ago',
			type: 'create'
		}
	]);

	let servers = $state([
		{
			id: 1,
			name: 'Localhost',
			provider: 'DigitalOcean',
			cpu: 82,
			mem: 82,
			lastCheck: '30s ago'
		}
	]);

	function goToProject(projectId: number) {
		goto(`/project/${projectId}`);
	}

	function getActivityColor(type: string) {
		switch (type) {
			case 'create':
				return 'bg-green-500';
			case 'failed':
				return 'bg-red-500';
			case 'pause':
				return 'bg-yellow-500';
			default:
				return 'bg-blue-500';
		}
	}
</script>

<svelte:head>
	<title>Dashboard - mikrocloud</title>
</svelte:head>

<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
	<div class="lg:col-span-2 space-y-6">
		<div>
			<h2 class="text-lg font-semibold mb-4">Projects</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				{#each projects as project (project.id)}
					<div
						class="bg-white/5 border border-white/10 rounded-lg p-4 hover:bg-white/10 transition-colors cursor-pointer"
						onclick={() => goToProject(project.id)}
					>
						<div class="flex items-center justify-between mb-8">
							<div class="flex items-center space-x-2">
								<Globe class="w-4 h-4 text-gray-400" />
								<span class="text-sm">{project.name}</span>
							</div>
							<button class="text-gray-400 hover:text-white" onclick={(e) => e.stopPropagation()}>
								<Ellipsis class="w-5 h-5" />
							</button>
						</div>
						<div class="flex items-center justify-between">
							<span class="text-xs text-gray-500">Last deploy : {project.lastDeploy}</span>
							<div class="flex items-center space-x-1">
								{#each project.icons as icon}
									<div
										class="w-6 h-6 bg-white/10 rounded-full flex items-center justify-center text-xs"
									>
										{icon}
									</div>
								{/each}
							</div>
						</div>
					</div>
				{/each}
			</div>
		</div>

		<div>
			<h2 class="text-lg font-semibold mb-4">Servers</h2>
			<div class="bg-white/5 border border-white/10 rounded-lg p-4">
				{#each servers as server (server.id)}
					<div>
						<div class="flex items-center space-x-2 mb-2">
							<span class="text-sm">🖥️</span>
							<span class="text-sm font-medium">{server.name}</span>
						</div>
						<div class="text-xs text-gray-500 mb-3">{server.provider}</div>
						<div class="space-y-2 mb-3">
							<div>
								<div class="flex items-center justify-between text-xs text-gray-400 mb-1">
									<span>CPU :</span>
									<span>{server.cpu}%</span>
								</div>
								<div class="w-full bg-white/10 rounded-full h-1.5">
									<div class="bg-white h-1.5 rounded-full" style="width: {server.cpu}%"></div>
								</div>
							</div>
							<div>
								<div class="flex items-center justify-between text-xs text-gray-400 mb-1">
									<span>MEM:</span>
									<span>{server.mem}%</span>
								</div>
								<div class="w-full bg-white/10 rounded-full h-1.5">
									<div class="bg-white h-1.5 rounded-full" style="width: {server.mem}%"></div>
								</div>
							</div>
						</div>
						<div class="text-xs text-gray-500">Last check: {server.lastCheck}</div>
					</div>
				{/each}
			</div>
		</div>
	</div>

	<div>
		<h2 class="text-lg font-semibold mb-4">Activity</h2>
		<div class="bg-white/5 border border-white/10 rounded-lg p-4">
			<div class="space-y-6">
				{#each activities as activity (activity.id)}
					<div class="flex items-start space-x-3">
						<div class="relative">
							<div
								class="w-8 h-8 rounded-full {getActivityColor(
									activity.type
								)} flex items-center justify-center text-xs flex-shrink-0"
							>
								{activity.type === 'create' ? '✓' : activity.type === 'failed' ? '✕' : '⏸'}
							</div>
							{#if activity.id !== activities[activities.length - 1].id}
								<div class="absolute top-8 left-4 w-px h-8 bg-white/10"></div>
							{/if}
						</div>
						<div class="flex-1 min-w-0">
							<div class="flex items-center justify-between mb-1">
								<div class="flex items-center space-x-2">
									<span class="text-xs">{activity.avatar}</span>
									<span class="text-sm font-medium">{activity.user}</span>
								</div>
								<span class="text-xs text-gray-500">{activity.time}</span>
							</div>
							<div class="text-xs text-gray-400 mb-1">{activity.action}</div>
							<div class="text-xs bg-white/5 border border-white/10 rounded px-2 py-1 inline-block">
								{activity.project}
							</div>
						</div>
					</div>
				{/each}
			</div>
		</div>
	</div>
</div>

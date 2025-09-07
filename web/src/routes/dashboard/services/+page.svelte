<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import {
		Plus,
		Search,
		ChevronDown,
		MoreHorizontal,
		Github,
		Database,
		Globe,
		HelpCircle,
		FileText,
		Settings,
		LogOut,
		Server,
		GitBranch,
		HardDrive,
		Key,
		Bell,
		Terminal,
		Users,
		User,
		Grid3X3,
		List,
		ExternalLink,
		Play,
		Pause,
		RotateCcw
	} from 'lucide-svelte';

	// Mock data for services
	let services = [
		{
			id: 1,
			name: 'web-frontend',
			type: 'Web Service',
			status: 'running',
			instances: 2,
			cpu: '45%',
			memory: '512MB',
			lastDeploy: '2h ago',
			repo: 'github.com/company/frontend'
		},
		{
			id: 2,
			name: 'api-backend',
			type: 'Backend Service',
			status: 'running',
			instances: 3,
			cpu: '23%',
			memory: '1.2GB',
			lastDeploy: '4h ago',
			repo: 'github.com/company/backend'
		},
		{
			id: 3,
			name: 'postgres-db',
			type: 'Database',
			status: 'running',
			instances: 1,
			cpu: '12%',
			memory: '2GB',
			lastDeploy: '1d ago',
			repo: 'postgres:14'
		},
		{
			id: 4,
			name: 'redis-cache',
			type: 'Cache',
			status: 'stopped',
			instances: 0,
			cpu: '0%',
			memory: '0MB',
			lastDeploy: '3d ago',
			repo: 'redis:alpine'
		}
	];

	let activeSection = 'Services';
	let viewMode = 'grid';
	let currentWorkspace = 'Personal Workspace';

	function toggleViewMode() {
		viewMode = viewMode === 'grid' ? 'list' : 'grid';
	}

	function getStatusColor(status: string) {
		switch (status) {
			case 'running':
				return 'bg-green-500';
			case 'stopped':
				return 'bg-red-500';
			case 'starting':
				return 'bg-yellow-500 animate-pulse';
			default:
				return 'bg-gray-400';
		}
	}

	function getServiceIcon(type: string) {
		switch (type.toLowerCase()) {
			case 'web service':
				return Globe;
			case 'backend service':
				return Server;
			case 'database':
				return Database;
			case 'cache':
				return Database;
			default:
				return Server;
		}
	}

	function startService(serviceId: number) {
		const service = services.find((s) => s.id === serviceId);
		if (service) {
			service.status = 'starting';
			setTimeout(() => {
				service.status = 'running';
				service.instances = 1;
			}, 2000);
		}
	}

	function stopService(serviceId: number) {
		const service = services.find((s) => s.id === serviceId);
		if (service) {
			service.status = 'stopped';
			service.instances = 0;
			service.cpu = '0%';
			service.memory = '0MB';
		}
	}

	function restartService(serviceId: number) {
		const service = services.find((s) => s.id === serviceId);
		if (service) {
			service.status = 'starting';
			setTimeout(() => {
				service.status = 'running';
			}, 2000);
		}
	}
</script>

<svelte:head>
	<title>Services - Dashboard</title>
</svelte:head>

<!-- Main Content -->
<div class="flex-1 flex flex-col overflow-hidden">
	<!-- Header -->
	<div class="bg-white border-b border-gray-200 px-6 py-4">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-2xl font-semibold text-gray-900">Services</h1>
				<p class="text-sm text-gray-500 mt-1">Manage and monitor all your running services.</p>
			</div>
			<div class="flex items-center space-x-3">
				<div class="flex items-center bg-gray-100 rounded-md p-1">
					<button
						onclick={toggleViewMode}
						class="p-1 rounded {viewMode === 'grid' ? 'bg-white shadow-sm' : ''}"
					>
						<Grid3X3 class="w-4 h-4" />
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
					New Service
				</Button>
			</div>
		</div>
	</div>

	<!-- Services Grid/List -->
	<div class="flex-1 overflow-auto p-6">
		{#if viewMode === 'grid'}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each services as service}
					{@const IconComponent = getServiceIcon(service.type)}
					<Card class="hover:shadow-md transition-shadow">
						<CardHeader class="pb-3">
							<div class="flex items-center justify-between">
								<div class="flex items-center space-x-2">
									<IconComponent class="w-5 h-5 text-gray-600" />
									<CardTitle class="text-lg">{service.name}</CardTitle>
								</div>
								<div class="flex items-center space-x-2">
									<div class="w-2 h-2 rounded-full {getStatusColor(service.status)}"></div>
									<Badge variant="outline" class="text-xs">
										{service.status}
									</Badge>
								</div>
							</div>
							<p class="text-sm text-gray-600">{service.type}</p>
						</CardHeader>
						<CardContent>
							<div class="space-y-3">
								<div class="grid grid-cols-2 gap-4 text-sm">
									<div>
										<span class="text-gray-500">Instances</span>
										<p class="font-medium">{service.instances}</p>
									</div>
									<div>
										<span class="text-gray-500">CPU</span>
										<p class="font-medium">{service.cpu}</p>
									</div>
									<div>
										<span class="text-gray-500">Memory</span>
										<p class="font-medium">{service.memory}</p>
									</div>
									<div>
										<span class="text-gray-500">Last Deploy</span>
										<p class="font-medium">{service.lastDeploy}</p>
									</div>
								</div>
								<div class="text-xs text-gray-500 truncate">
									{service.repo}
								</div>
								<div class="flex space-x-2">
									{#if service.status === 'stopped'}
										<Button size="sm" variant="outline" onclick={() => startService(service.id)}>
											<Play class="w-3 h-3 mr-1" />
											Start
										</Button>
									{:else}
										<Button size="sm" variant="outline" onclick={() => stopService(service.id)}>
											<Pause class="w-3 h-3 mr-1" />
											Stop
										</Button>
									{/if}
									<Button size="sm" variant="outline" onclick={() => restartService(service.id)}>
										<RotateCcw class="w-3 h-3 mr-1" />
										Restart
									</Button>
								</div>
							</div>
						</CardContent>
					</Card>
				{/each}
			</div>
		{:else}
			<div class="space-y-3">
				{#each services as service}
					{@const IconComponent = getServiceIcon(service.type)}
					<Card>
						<CardContent class="p-4">
							<div class="flex items-center justify-between">
								<div class="flex items-center space-x-4">
									<div class="w-2 h-2 rounded-full {getStatusColor(service.status)}"></div>
									<div class="flex items-center space-x-2">
										<IconComponent class="w-5 h-5 text-gray-600" />
										<div>
											<h3 class="font-medium text-gray-900">{service.name}</h3>
											<p class="text-sm text-gray-600">{service.type}</p>
										</div>
									</div>
								</div>
								<div class="flex items-center space-x-6">
									<div class="text-sm">
										<span class="text-gray-500">Instances:</span>
										<span class="font-medium ml-1">{service.instances}</span>
									</div>
									<div class="text-sm">
										<span class="text-gray-500">CPU:</span>
										<span class="font-medium ml-1">{service.cpu}</span>
									</div>
									<div class="text-sm">
										<span class="text-gray-500">Memory:</span>
										<span class="font-medium ml-1">{service.memory}</span>
									</div>
									<div class="flex space-x-2">
										{#if service.status === 'stopped'}
											<Button size="sm" variant="outline" onclick={() => startService(service.id)}>
												<Play class="w-3 h-3" />
												Start
											</Button>
										{:else}
											<Button size="sm" variant="outline" onclick={() => stopService(service.id)}>
												<Pause class="w-3 h-3" />
												Stop
											</Button>
										{/if}
										<Button size="sm" variant="outline" onclick={() => restartService(service.id)}>
											<RotateCcw class="w-3 h-3" />
											Restart
										</Button>
									</div>
								</div>
							</div>
						</CardContent>
					</Card>
				{/each}
			</div>
		{/if}
	</div>
</div>

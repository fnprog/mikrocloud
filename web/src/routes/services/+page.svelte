<script>
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
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
		ArrowUpDown
	} from 'lucide-svelte';

	// Mock data for services
	let services = $state([
		{
			name: 'App',
			type: 'APP',
			tag: 'PERSONAL',
			status: 'Deploying',
			statusColor: 'text-red-500',
			indicator: 'bg-red-500',
			lastDeploy: null
		},
		{
			name: 'Test - 002',
			type: 'DATABASE',
			tag: 'STATIC',
			status: 'Active',
			statusColor: 'text-green-500',
			indicator: 'bg-green-500',
			lastDeploy: '4:02PM'
		},
		{
			name: 'Dashboard',
			type: null,
			tag: null,
			status: 'Active',
			statusColor: 'text-green-500',
			indicator: 'bg-green-500',
			lastDeploy: '16.06.2024'
		},
		{
			name: 'Front-end',
			type: 'STATIC',
			tag: null,
			status: 'Paused',
			statusColor: 'text-gray-500',
			indicator: 'bg-gray-400',
			lastDeploy: '13.06.2024'
		},
		{
			name: 'New project',
			type: 'PERSONAL',
			tag: null,
			status: 'Paused',
			statusColor: 'text-gray-500',
			indicator: 'bg-gray-400',
			lastDeploy: '13.06.2024'
		},
		{
			name: 'Optima',
			type: 'PROJECT',
			tag: 'DASHBOARD',
			status: 'Disconnected',
			statusColor: 'text-gray-500',
			indicator: 'bg-gray-400',
			lastDeploy: '12.06.2024'
		},
		{
			name: 'Closed',
			type: 'APP',
			tag: null,
			status: 'Disconnected',
			statusColor: 'text-gray-500',
			indicator: 'bg-gray-400',
			lastDeploy: '14.04.2024'
		}
	]);

	let activeSection = $state('Overview');
	let showPrevious = $state(false);

	function handleLogout() {
		goto('/login');
	}

	function setActiveSection(section) {
		activeSection = section;
	}

	function togglePrevious() {
		showPrevious = !showPrevious;
	}

	// Filter services based on showPrevious
	const visibleServices = $derived.by(() => {
		if (showPrevious) {
			return services;
		}
		return services.filter((service) => service.status !== 'Disconnected');
	});

	const currentServices = $derived.by(() => {
		return visibleServices.filter(
			(service) => service.status === 'Active' || service.status === 'Deploying'
		);
	});

	const lastWeekServices = $derived.by(() => {
		return visibleServices.filter((service) => service.status === 'Paused');
	});

	const previousServices = $derived.by(() => {
		return visibleServices.filter((service) => service.status === 'Disconnected');
	});
</script>

<svelte:head>
	<title>Dashboard - Services</title>
</svelte:head>

<div class="flex h-screen bg-gray-50">
	<!-- Sidebar -->
	<div class="w-80 bg-white border-r border-gray-200 flex flex-col">
		<!-- User Profile -->
		<div class="p-4 border-b border-gray-200">
			<div class="flex items-center space-x-3">
				<div class="w-8 h-8 bg-gray-800 rounded-full flex items-center justify-center">
					<span class="text-white text-sm font-medium">R</span>
				</div>
				<div class="flex-1">
					<div class="flex items-center space-x-2">
						<span class="font-medium text-gray-900">Raymond</span>
						<ChevronDown class="w-4 h-4 text-gray-500" />
					</div>
					<div class="text-sm text-gray-500">rayspace</div>
				</div>
				<Search class="w-5 h-5 text-gray-400" />
			</div>
		</div>

		<!-- Run Service Button -->
		<div class="p-4 border-b border-gray-200">
			<Button class="w-full bg-gray-800 hover:bg-gray-700 text-white justify-start">
				<Plus class="w-4 h-4 mr-2" />
				Run Service
				<span class="ml-auto text-xs bg-gray-700 px-2 py-1 rounded">⌘ /</span>
			</Button>
		</div>

		<!-- Navigation -->
		<nav class="flex-1 p-4">
			<div class="space-y-1">
				<button
					onclick={() => setActiveSection('Overview')}
					class="w-full flex items-center px-3 py-2 text-sm font-medium rounded-md {activeSection ===
					'Overview'
						? 'bg-gray-100 text-gray-900'
						: 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}"
				>
					<div class="w-5 h-5 mr-3 flex items-center justify-center">
						<div class="w-2 h-2 bg-gray-400 rounded-full"></div>
					</div>
					Overview
				</button>

				<button
					onclick={() => setActiveSection('Domains')}
					class="w-full flex items-center px-3 py-2 text-sm font-medium rounded-md {activeSection ===
					'Domains'
						? 'bg-gray-100 text-gray-900'
						: 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}"
				>
					<Globe class="w-5 h-5 mr-3" />
					Domains
				</button>

				<button
					onclick={() => setActiveSection('Databases')}
					class="w-full flex items-center px-3 py-2 text-sm font-medium rounded-md {activeSection ===
					'Databases'
						? 'bg-gray-100 text-gray-900'
						: 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}"
				>
					<Database class="w-5 h-5 mr-3" />
					Databases
				</button>

				<button
					onclick={() => setActiveSection('Upgrade')}
					class="w-full flex items-center px-3 py-2 text-sm font-medium rounded-md {activeSection ===
					'Upgrade'
						? 'bg-gray-100 text-gray-900'
						: 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}"
				>
					<div class="w-5 h-5 mr-3 flex items-center justify-center">
						<div class="w-3 h-3 border-2 border-gray-400 rounded-sm"></div>
					</div>
					Upgrade
				</button>
			</div>

			<!-- Other Section -->
			<div class="mt-8">
				<div class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-3">Other</div>
				<div class="space-y-1">
					<button
						class="w-full flex items-center px-3 py-2 text-sm font-medium text-gray-600 rounded-md hover:bg-gray-50 hover:text-gray-900"
					>
						<HelpCircle class="w-5 h-5 mr-3" />
						Help & Support
					</button>
					<button
						class="w-full flex items-center px-3 py-2 text-sm font-medium text-gray-600 rounded-md hover:bg-gray-50 hover:text-gray-900"
					>
						<FileText class="w-5 h-5 mr-3" />
						Documentation
					</button>
					<button
						class="w-full flex items-center px-3 py-2 text-sm font-medium text-gray-600 rounded-md hover:bg-gray-50 hover:text-gray-900"
					>
						<Settings class="w-5 h-5 mr-3" />
						Settings
					</button>
				</div>
			</div>

			<!-- Tags Section -->
			<div class="mt-8">
				<div class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-3">Tags</div>
				<div class="space-y-1">
					<button
						class="w-full flex items-center px-3 py-2 text-sm font-medium text-gray-600 rounded-md hover:bg-gray-50 hover:text-gray-900"
					>
						<div class="w-2 h-2 bg-gray-400 rounded-full mr-3"></div>
						App
					</button>
					<button
						class="w-full flex items-center px-3 py-2 text-sm font-medium text-gray-600 rounded-md hover:bg-gray-50 hover:text-gray-900"
					>
						<div class="w-2 h-2 bg-gray-400 rounded-full mr-3"></div>
						Database
					</button>
					<button
						class="w-full flex items-center px-3 py-2 text-sm font-medium text-gray-600 rounded-md hover:bg-gray-50 hover:text-gray-900"
					>
						<div class="w-2 h-2 bg-gray-400 rounded-full mr-3"></div>
						Personal
					</button>
					<button
						class="w-full flex items-center px-3 py-2 text-sm font-medium text-gray-600 rounded-md hover:bg-gray-50 hover:text-gray-900"
					>
						<div class="w-2 h-2 bg-gray-400 rounded-full mr-3"></div>
						Static
					</button>
				</div>
			</div>
		</nav>

		<!-- Logout -->
		<div class="p-4 border-t border-gray-200">
			<button
				onclick={handleLogout}
				class="w-full flex items-center px-3 py-2 text-sm font-medium text-gray-600 rounded-md hover:bg-gray-50 hover:text-gray-900"
			>
				<LogOut class="w-5 h-5 mr-3" />
				Log out
			</button>
		</div>
	</div>

	<!-- Main Content -->
	<div class="flex-1 flex flex-col overflow-hidden">
		<!-- Header -->
		<div class="bg-white border-b border-gray-200 px-6 py-4">
			<div class="flex items-center justify-between">
				<div>
					<h1 class="text-2xl font-semibold text-gray-900">Services</h1>
					<p class="text-sm text-gray-500 mt-1">
						Run service from your private or any public repository.
					</p>
				</div>
				<div class="flex items-center space-x-2">
					<Button variant="outline" size="sm">
						<ArrowUpDown class="w-4 h-4 mr-2" />
						<ChevronDown class="w-4 h-4" />
					</Button>
					<Button variant="ghost" size="sm">
						<MoreHorizontal class="w-4 h-4" />
					</Button>
				</div>
			</div>
		</div>

		<!-- Services List -->
		<div class="flex-1 overflow-auto p-6">
			<div class="space-y-6">
				<!-- Current Services -->
				{#if currentServices.length > 0}
					<div class="space-y-3">
						{#each currentServices as service}
							<div class="bg-white rounded-lg border border-gray-200 p-4">
								<div class="flex items-center justify-between">
									<div class="flex items-center space-x-3">
										<div class="w-2 h-2 rounded-full {service.indicator}"></div>
										<div>
											<div class="flex items-center space-x-2">
												<span class="font-medium text-gray-900">{service.name}</span>
												{#if service.type}
													<Badge variant="secondary" class="text-xs">{service.type}</Badge>
												{/if}
												{#if service.tag}
													<Badge variant="outline" class="text-xs">{service.tag}</Badge>
												{/if}
											</div>
											{#if service.lastDeploy}
												<div class="text-sm text-gray-500 mt-1">
													Last deploy {service.lastDeploy}
												</div>
											{/if}
										</div>
									</div>
									<div class="flex items-center space-x-3">
										<span class="text-sm font-medium {service.statusColor}">{service.status}</span>
										<Github class="w-4 h-4 text-gray-400" />
										<Button variant="ghost" size="sm">
											<MoreHorizontal class="w-4 h-4" />
										</Button>
									</div>
								</div>
							</div>
						{/each}
					</div>
				{/if}

				<!-- Last Week Services -->
				{#if lastWeekServices.length > 0}
					<div>
						<div class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-3">
							Last Week
						</div>
						<div class="space-y-3">
							{#each lastWeekServices as service}
								<div class="bg-white rounded-lg border border-gray-200 p-4">
									<div class="flex items-center justify-between">
										<div class="flex items-center space-x-3">
											<div class="w-2 h-2 rounded-full {service.indicator}"></div>
											<div>
												<div class="flex items-center space-x-2">
													<span class="font-medium text-gray-900">{service.name}</span>
													{#if service.type}
														<Badge variant="secondary" class="text-xs">{service.type}</Badge>
													{/if}
													{#if service.tag}
														<Badge variant="outline" class="text-xs">{service.tag}</Badge>
													{/if}
												</div>
												{#if service.lastDeploy}
													<div class="text-sm text-gray-500 mt-1">
														Last deploy {service.lastDeploy}
													</div>
												{/if}
											</div>
										</div>
										<div class="flex items-center space-x-3">
											<span class="text-sm font-medium {service.statusColor}">{service.status}</span
											>
											<Github class="w-4 h-4 text-gray-400" />
											<Button variant="ghost" size="sm">
												<MoreHorizontal class="w-4 h-4" />
											</Button>
										</div>
									</div>
								</div>
							{/each}
						</div>
					</div>
				{/if}

				<!-- Previous Services -->
				{#if showPrevious && previousServices.length > 0}
					<div>
						<div class="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-3">
							Previous
						</div>
						<div class="space-y-3">
							{#each previousServices as service}
								<div class="bg-white rounded-lg border border-gray-200 p-4">
									<div class="flex items-center justify-between">
										<div class="flex items-center space-x-3">
											<div class="w-2 h-2 rounded-full {service.indicator}"></div>
											<div>
												<div class="flex items-center space-x-2">
													<span class="font-medium text-gray-900">{service.name}</span>
													{#if service.type}
														<Badge variant="secondary" class="text-xs">{service.type}</Badge>
													{/if}
													{#if service.tag}
														<Badge variant="outline" class="text-xs">{service.tag}</Badge>
													{/if}
												</div>
												{#if service.lastDeploy}
													<div class="text-sm text-gray-500 mt-1">
														Last deploy {service.lastDeploy}
													</div>
												{/if}
											</div>
										</div>
										<div class="flex items-center space-x-3">
											<span class="text-sm font-medium {service.statusColor}">{service.status}</span
											>
											<Github class="w-4 h-4 text-gray-400" />
											<Button variant="ghost" size="sm">
												<MoreHorizontal class="w-4 h-4" />
											</Button>
										</div>
									</div>
								</div>
							{/each}
						</div>
					</div>
				{/if}

				<!-- Show Previous Button -->
				{#if !showPrevious && previousServices.length > 0}
					<div class="text-center">
						<Button variant="outline" onclick={togglePrevious}>Show previous</Button>
					</div>
				{/if}
			</div>
		</div>
	</div>
</div>

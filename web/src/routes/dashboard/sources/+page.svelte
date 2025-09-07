<script>
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import {
		Plus,
		Search,
		ChevronDown,
		MoreHorizontal,
		Github,
		GitBranch,
		ExternalLink,
		Settings,
		LogOut,
		Server,
		HardDrive,
		Key,
		Bell,
		Terminal,
		Users,
		User,
		Grid3X3,
		FileText,
		Clock,
		Star
	} from 'lucide-svelte';

	// Mock data for sources
	let sources = $state([
		{
			id: 1,
			name: 'frontend-app',
			provider: 'GitHub',
			repository: 'company/frontend-app',
			branch: 'main',
			lastCommit: 'feat: add user dashboard',
			lastCommitTime: '2 hours ago',
			author: 'John Doe',
			status: 'connected',
			private: true,
			stars: 24
		},
		{
			id: 2,
			name: 'api-backend',
			provider: 'GitHub',
			repository: 'company/api-backend',
			branch: 'develop',
			lastCommit: 'fix: resolve authentication issue',
			lastCommitTime: '5 hours ago',
			author: 'Jane Smith',
			status: 'connected',
			private: true,
			stars: 18
		},
		{
			id: 3,
			name: 'mobile-app',
			provider: 'GitLab',
			repository: 'company/mobile-app',
			branch: 'main',
			lastCommit: 'refactor: update API endpoints',
			lastCommitTime: '1 day ago',
			author: 'Bob Johnson',
			status: 'disconnected',
			private: false,
			stars: 42
		},
		{
			id: 4,
			name: 'documentation',
			provider: 'GitHub',
			repository: 'company/docs',
			branch: 'main',
			lastCommit: 'docs: update deployment guide',
			lastCommitTime: '3 days ago',
			author: 'Alice Brown',
			status: 'connected',
			private: false,
			stars: 7
		}
	]);

	let activeSection = $state('Sources');
	let currentWorkspace = $state('Personal Workspace');

	function handleLogout() {
		goto('/login');
	}

	function setActiveSection(section) {
		activeSection = section;
		if (section === 'Dashboard') {
			goto('/dashboard');
		} else {
			goto(`/dashboard/${section.toLowerCase().replace(' ', '-')}`);
		}
	}

	function getStatusColor(status) {
		switch (status) {
			case 'connected':
				return 'bg-green-500';
			case 'disconnected':
				return 'bg-red-500';
			case 'syncing':
				return 'bg-yellow-500 animate-pulse';
			default:
				return 'bg-gray-400';
		}
	}

	function getProviderIcon(provider) {
		switch (provider.toLowerCase()) {
			case 'github':
				return Github;
			case 'gitlab':
				return GitBranch;
			default:
				return GitBranch;
		}
	}

	function connectSource(sourceId) {
		const source = sources.find((s) => s.id === sourceId);
		if (source) {
			source.status = 'syncing';
			setTimeout(() => {
				source.status = 'connected';
			}, 2000);
		}
	}

	function disconnectSource(sourceId) {
		const source = sources.find((s) => s.id === sourceId);
		if (source) {
			source.status = 'disconnected';
		}
	}
</script>

<svelte:head>
	<title>Sources - Dashboard</title>
</svelte:head>

<div class="flex h-screen bg-gray-50">
	<!-- Sidebar -->
	<div class="w-80 bg-white border-r border-gray-200 flex flex-col">
		<!-- Workspace Switcher -->
		<div class="p-4 border-b border-gray-200">
			<div class="flex items-center space-x-3">
				<div class="w-8 h-8 bg-gray-800 rounded-full flex items-center justify-center">
					<span class="text-white text-sm font-medium">P</span>
				</div>
				<div class="flex-1">
					<div class="flex items-center space-x-2">
						<span class="font-medium text-gray-900">{currentWorkspace}</span>
						<ChevronDown class="w-4 h-4 text-gray-500" />
					</div>
				</div>
				<Search class="w-5 h-5 text-gray-400" />
			</div>
		</div>

		<!-- Navigation -->
		<nav class="flex-1 p-4 overflow-y-auto">
			<div class="space-y-1">
				{#each ['Dashboard', 'Services', 'Servers', 'Sources', 'S3 Storage', 'Shared Variables', 'Notifications', 'Keys and Tokens', 'Terminal'] as section}
					<button
						onclick={() => setActiveSection(section)}
						class="w-full flex items-center px-3 py-2 text-sm font-medium rounded-md {activeSection ===
						section
							? 'bg-gray-100 text-gray-900'
							: 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}"
					>
						{#if section === 'Dashboard'}
							<Grid3X3 class="w-5 h-5 mr-3" />
						{:else if section === 'Services'}
							<div class="w-5 h-5 mr-3 flex items-center justify-center">
								<div class="w-2 h-2 bg-gray-400 rounded-full"></div>
							</div>
						{:else if section === 'Servers'}
							<Server class="w-5 h-5 mr-3" />
						{:else if section === 'Sources'}
							<GitBranch class="w-5 h-5 mr-3" />
						{:else if section === 'S3 Storage'}
							<HardDrive class="w-5 h-5 mr-3" />
						{:else if section === 'Shared Variables'}
							<FileText class="w-5 h-5 mr-3" />
						{:else if section === 'Notifications'}
							<Bell class="w-5 h-5 mr-3" />
						{:else if section === 'Keys and Tokens'}
							<Key class="w-5 h-5 mr-3" />
						{:else if section === 'Terminal'}
							<Terminal class="w-5 h-5 mr-3" />
						{/if}
						{section}
					</button>
				{/each}
			</div>
		</nav>

		<!-- Bottom Section -->
		<div class="p-4 border-t border-gray-200 space-y-1">
			<button
				class="w-full flex items-center px-3 py-2 text-sm font-medium text-gray-600 rounded-md hover:bg-gray-50 hover:text-gray-900"
			>
				<User class="w-5 h-5 mr-3" />
				Profile
			</button>
			<button
				class="w-full flex items-center px-3 py-2 text-sm font-medium text-gray-600 rounded-md hover:bg-gray-50 hover:text-gray-900"
			>
				<Users class="w-5 h-5 mr-3" />
				Team
			</button>
			<button
				class="w-full flex items-center px-3 py-2 text-sm font-medium text-gray-600 rounded-md hover:bg-gray-50 hover:text-gray-900"
			>
				<Settings class="w-5 h-5 mr-3" />
				Settings
			</button>
			<button
				onclick={handleLogout}
				class="w-full flex items-center px-3 py-2 text-sm font-medium text-gray-600 rounded-md hover:bg-gray-50 hover:text-gray-900"
			>
				<LogOut class="w-5 h-5 mr-3" />
				Logout
			</button>
		</div>
	</div>

	<!-- Main Content -->
	<div class="flex-1 flex flex-col overflow-hidden">
		<!-- Header -->
		<div class="bg-white border-b border-gray-200 px-6 py-4">
			<div class="flex items-center justify-between">
				<div>
					<h1 class="text-2xl font-semibold text-gray-900">Sources</h1>
					<p class="text-sm text-gray-500 mt-1">
						Manage your connected repositories and source code.
					</p>
				</div>
				<Button>
					<Plus class="w-4 h-4 mr-2" />
					Connect Repository
				</Button>
			</div>
		</div>

		<!-- Sources Grid -->
		<div class="flex-1 overflow-auto p-6">
			<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
				{#each sources as source}
					<Card class="hover:shadow-md transition-shadow">
						<CardHeader class="pb-3">
							<div class="flex items-center justify-between">
								<div class="flex items-center space-x-2">
									{#if source.provider.toLowerCase() === 'github'}
										<Github class="w-5 h-5 text-gray-600" />
									{:else}
										<GitBranch class="w-5 h-5 text-gray-600" />
									{/if}
									<CardTitle class="text-lg">{source.name}</CardTitle>
									{#if source.private}
										<Badge variant="secondary" class="text-xs">Private</Badge>
									{:else}
										<Badge variant="outline" class="text-xs">Public</Badge>
									{/if}
								</div>
								<div class="flex items-center space-x-2">
									<div class="w-2 h-2 rounded-full {getStatusColor(source.status)}"></div>
									<Badge variant="outline" class="text-xs">
										{source.status}
									</Badge>
								</div>
							</div>
							<div class="flex items-center space-x-4 text-sm text-gray-600">
								<span>{source.repository}</span>
								<span>•</span>
								<span>{source.branch}</span>
								{#if source.stars > 0}
									<div class="flex items-center space-x-1">
										<Star class="w-3 h-3" />
										<span>{source.stars}</span>
									</div>
								{/if}
							</div>
						</CardHeader>
						<CardContent>
							<div class="space-y-4">
								<!-- Last Commit -->
								<div class="bg-gray-50 rounded-lg p-3">
									<div class="flex items-center space-x-2 mb-1">
										<Clock class="w-4 h-4 text-gray-500" />
										<span class="text-sm font-medium text-gray-900">Latest Commit</span>
									</div>
									<p class="text-sm text-gray-700 mb-1">{source.lastCommit}</p>
									<div class="flex items-center space-x-2 text-xs text-gray-500">
										<span>by {source.author}</span>
										<span>•</span>
										<span>{source.lastCommitTime}</span>
									</div>
								</div>

								<!-- Actions -->
								<div class="flex space-x-2">
									{#if source.status === 'disconnected'}
										<Button
											size="sm"
											variant="outline"
											onclick={() => connectSource(source.id)}
											class="flex-1"
										>
											Connect
										</Button>
									{:else}
										<Button
											size="sm"
											variant="outline"
											onclick={() => disconnectSource(source.id)}
											class="flex-1"
										>
											Disconnect
										</Button>
									{/if}
									<Button size="sm" variant="outline">
										<ExternalLink class="w-4 h-4" />
									</Button>
									<Button size="sm" variant="outline">
										<MoreHorizontal class="w-4 h-4" />
									</Button>
								</div>
							</div>
						</CardContent>
					</Card>
				{/each}
			</div>
		</div>
	</div>
</div>

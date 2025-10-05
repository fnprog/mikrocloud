<script>
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import {
		Settings,
		Save,
		Trash2,
		AlertTriangle,
		Shield,
		Bell,
		Users,
		Database,
		Code,
		Globe,
		Lock
	} from 'lucide-svelte';

	let project = $state({
		name: 'Focalpoint Dashboard',
		description: 'A comprehensive dashboard application for managing business operations',
		workspace: 'Focalpoint',
		category: 'Applications',
		visibility: 'private',
		region: 'us-east-1',
		framework: 'Next.js',
		buildCommand: 'npm run build',
		startCommand: 'npm start',
		nodeVersion: '18.x',
		autoDeployEnabled: true,
		branch: 'main',
		rootDirectory: '/',
		outputDirectory: '.next'
	});

	let activeSection = $state('Settings');
	let currentUser = $state('Adam Smith');
	let activeTab = $state('general');
	let showDeleteModal = $state(false);
	let deleteConfirmation = $state('');

	// Mock collaborators
	let collaborators = $state([
		{
			id: 1,
			name: 'John Doe',
			email: 'john@focalpoint.com',
			role: 'owner',
			avatar: 'JD',
			joinedAt: '2024-01-15'
		},
		{
			id: 2,
			name: 'Jane Smith',
			email: 'jane@focalpoint.com',
			role: 'admin',
			avatar: 'JS',
			joinedAt: '2024-01-20'
		},
		{
			id: 3,
			name: 'Bob Johnson',
			email: 'bob@focalpoint.com',
			role: 'developer',
			avatar: 'BJ',
			joinedAt: '2024-01-25'
		}
	]);

	// Mock integrations
	let integrations = $state([
		{
			id: 1,
			name: 'GitHub',
			description: 'Source code repository',
			status: 'connected',
			icon: Code,
			connectedAt: '2024-01-15'
		},
		{
			id: 2,
			name: 'Slack',
			description: 'Team notifications',
			status: 'connected',
			icon: Bell,
			connectedAt: '2024-01-20'
		},
		{
			id: 3,
			name: 'MongoDB Atlas',
			description: 'Database service',
			status: 'disconnected',
			icon: Database,
			connectedAt: null
		}
	]);

	function saveGeneralSettings() {
		console.log('Saving general settings...', project);
		// Mock save
	}

	function saveBuildSettings() {
		console.log('Saving build settings...', project);
		// Mock save
	}

	function deleteProject() {
		if (deleteConfirmation === project.name) {
			console.log('Deleting project...');
			goto('/dashboard');
		}
	}

	function getRoleColor(role) {
		switch (role) {
			case 'owner':
				return 'text-purple-600 bg-purple-50';
			case 'admin':
				return 'text-blue-600 bg-blue-50';
			case 'developer':
				return 'text-green-600 bg-green-50';
			default:
				return 'text-gray-600 bg-gray-50';
		}
	}

	function removeCollaborator(collaboratorId) {
		collaborators = collaborators.filter((c) => c.id !== collaboratorId);
	}

	function toggleIntegration(integrationId) {
		const integration = integrations.find((i) => i.id === integrationId);
		if (integration) {
			integration.status = integration.status === 'connected' ? 'disconnected' : 'connected';
			integration.connectedAt =
				integration.status === 'connected' ? new Date().toISOString().split('T')[0] : null;
		}
	}
</script>

<svelte:head>
	<title>Settings - {project.name}</title>
</svelte:head>

<!-- Main Content -->
<div class="flex-1 p-6">
	<!-- Header -->
	<div class="flex items-center justify-between mb-6">
		<div class="flex items-center space-x-4">
			<div>
				<h1 class="text-2xl font-semibold text-gray-900">Settings</h1>
				<p class="text-sm text-gray-500 mt-1">
					Configure project settings and preferences for {project.name}.
				</p>
			</div>
		</div>
	</div>

	<!-- Settings Tabs -->
	<div class="mb-6">
		<div class="border-b border-gray-200">
			<nav class="-mb-px flex space-x-8">
				{#each [{ id: 'general', label: 'General', icon: Settings }, { id: 'build', label: 'Build & Deploy', icon: Code }, { id: 'collaborators', label: 'Collaborators', icon: Users }, { id: 'integrations', label: 'Integrations', icon: Globe }, { id: 'security', label: 'Security', icon: Shield }, { id: 'danger', label: 'Danger Zone', icon: AlertTriangle }] as tab (tab.id)}
					<button
						onclick={() => (activeTab = tab.id)}
						class="flex items-center space-x-2 py-2 px-1 border-b-2 font-medium text-sm {activeTab ===
						tab.id
							? 'border-blue-500 text-blue-600'
							: 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
					>
						<tab.icon class="w-4 h-4" />
						<span>{tab.label}</span>
					</button>
				{/each}
			</nav>
		</div>
	</div>

	<!-- Tab Content -->
	{#if activeTab === 'general'}
		<Card>
			<CardHeader>
				<CardTitle>General Settings</CardTitle>
			</CardHeader>
			<CardContent class="space-y-6">
				<div class="grid grid-cols-2 gap-6">
					<div>
						<Label for="project-name">Project Name</Label>
						<Input id="project-name" bind:value={project.name} />
					</div>
					<div>
						<Label for="visibility">Visibility</Label>
						<select
							id="visibility"
							bind:value={project.visibility}
							class="w-full px-3 py-2 border border-gray-300 rounded-md"
						>
							<option value="private">Private</option>
							<option value="public">Public</option>
						</select>
					</div>
				</div>
				<div>
					<Label for="description">Description</Label>
					<Input id="description" bind:value={project.description} />
				</div>
				<div class="grid grid-cols-2 gap-6">
					<div>
						<Label for="region">Region</Label>
						<select
							id="region"
							bind:value={project.region}
							class="w-full px-3 py-2 border border-gray-300 rounded-md"
						>
							<option value="us-east-1">US East (N. Virginia)</option>
							<option value="us-west-2">US West (Oregon)</option>
							<option value="eu-west-1">Europe (Ireland)</option>
							<option value="ap-southeast-1">Asia Pacific (Singapore)</option>
						</select>
					</div>
					<div>
						<Label for="framework">Framework</Label>
						<select
							id="framework"
							bind:value={project.framework}
							class="w-full px-3 py-2 border border-gray-300 rounded-md"
						>
							<option value="Next.js">Next.js</option>
							<option value="React">React</option>
							<option value="Vue.js">Vue.js</option>
							<option value="Svelte">Svelte</option>
							<option value="Static">Static HTML</option>
						</select>
					</div>
				</div>
				<div class="flex justify-end">
					<Button onclick={saveGeneralSettings}>
						<Save class="w-4 h-4 mr-2" />
						Save Changes
					</Button>
				</div>
			</CardContent>
		</Card>
	{:else if activeTab === 'build'}
		<Card>
			<CardHeader>
				<CardTitle>Build & Deploy Settings</CardTitle>
			</CardHeader>
			<CardContent class="space-y-6">
				<div class="grid grid-cols-2 gap-6">
					<div>
						<Label for="build-command">Build Command</Label>
						<Input id="build-command" bind:value={project.buildCommand} />
					</div>
					<div>
						<Label for="start-command">Start Command</Label>
						<Input id="start-command" bind:value={project.startCommand} />
					</div>
				</div>
				<div class="grid grid-cols-3 gap-6">
					<div>
						<Label for="node-version">Node.js Version</Label>
						<select
							id="node-version"
							bind:value={project.nodeVersion}
							class="w-full px-3 py-2 border border-gray-300 rounded-md"
						>
							<option value="20.x">20.x (Latest)</option>
							<option value="18.x">18.x (LTS)</option>
							<option value="16.x">16.x</option>
						</select>
					</div>
					<div>
						<Label for="root-directory">Root Directory</Label>
						<Input id="root-directory" bind:value={project.rootDirectory} />
					</div>
					<div>
						<Label for="output-directory">Output Directory</Label>
						<Input id="output-directory" bind:value={project.outputDirectory} />
					</div>
				</div>
				<div class="flex items-center space-x-2">
					<input
						type="checkbox"
						id="auto-deploy"
						bind:checked={project.autoDeployEnabled}
						class="rounded"
					/>
					<Label for="auto-deploy">Enable automatic deployments from {project.branch} branch</Label>
				</div>
				<div class="flex justify-end">
					<Button onclick={saveBuildSettings}>
						<Save class="w-4 h-4 mr-2" />
						Save Changes
					</Button>
				</div>
			</CardContent>
		</Card>
	{:else if activeTab === 'collaborators'}
		<Card>
			<CardHeader>
				<div class="flex items-center justify-between">
					<CardTitle>Collaborators</CardTitle>
					<Button>
						<Users class="w-4 h-4 mr-2" />
						Invite Collaborator
					</Button>
				</div>
			</CardHeader>
			<CardContent>
				<div class="space-y-4">
					{#each collaborators as collaborator (collaborator.id)}
						<div class="flex items-center justify-between p-4 border rounded-lg">
							<div class="flex items-center space-x-4">
								<div class="w-10 h-10 bg-blue-500 rounded-full flex items-center justify-center">
									<span class="text-white text-sm font-medium">{collaborator.avatar}</span>
								</div>
								<div>
									<h4 class="font-medium text-gray-900">{collaborator.name}</h4>
									<p class="text-sm text-gray-600">{collaborator.email}</p>
								</div>
								<Badge variant="outline" class="text-xs {getRoleColor(collaborator.role)}">
									{collaborator.role}
								</Badge>
							</div>
							<div class="flex items-center space-x-2">
								<span class="text-sm text-gray-500">Joined {collaborator.joinedAt}</span>
								{#if collaborator.role !== 'owner'}
									<Button
										size="sm"
										variant="outline"
										onclick={() => removeCollaborator(collaborator.id)}
									>
										<Trash2 class="w-4 h-4" />
									</Button>
								{/if}
							</div>
						</div>
					{/each}
				</div>
			</CardContent>
		</Card>
	{:else if activeTab === 'integrations'}
		<Card>
			<CardHeader>
				<CardTitle>Integrations</CardTitle>
			</CardHeader>
			<CardContent>
				<div class="space-y-4">
					{#each integrations as integration (integration.id)}
						<div class="flex items-center justify-between p-4 border rounded-lg">
							<div class="flex items-center space-x-4">
								<integration.icon class="w-8 h-8 text-gray-600" />
								<div>
									<h4 class="font-medium text-gray-900">{integration.name}</h4>
									<p class="text-sm text-gray-600">{integration.description}</p>
								</div>
							</div>
							<div class="flex items-center space-x-4">
								<Badge
									variant="outline"
									class="text-xs {integration.status === 'connected'
										? 'text-green-600 bg-green-50'
										: 'text-gray-600 bg-gray-50'}"
								>
									{integration.status}
								</Badge>
								<Button
									size="sm"
									variant="outline"
									onclick={() => toggleIntegration(integration.id)}
								>
									{integration.status === 'connected' ? 'Disconnect' : 'Connect'}
								</Button>
							</div>
						</div>
					{/each}
				</div>
			</CardContent>
		</Card>
	{:else if activeTab === 'security'}
		<Card>
			<CardHeader>
				<CardTitle>Security Settings</CardTitle>
			</CardHeader>
			<CardContent class="space-y-6">
				<div class="flex items-center justify-between p-4 border rounded-lg">
					<div class="flex items-center space-x-4">
						<Lock class="w-6 h-6 text-gray-600" />
						<div>
							<h4 class="font-medium text-gray-900">Two-Factor Authentication</h4>
							<p class="text-sm text-gray-600">Add an extra layer of security to your account</p>
						</div>
					</div>
					<Button variant="outline">Enable 2FA</Button>
				</div>

				<div class="flex items-center justify-between p-4 border rounded-lg">
					<div class="flex items-center space-x-4">
						<Shield class="w-6 h-6 text-gray-600" />
						<div>
							<h4 class="font-medium text-gray-900">Access Logs</h4>
							<p class="text-sm text-gray-600">View recent access attempts and login history</p>
						</div>
					</div>
					<Button variant="outline">View Logs</Button>
				</div>

				<div class="flex items-center justify-between p-4 border rounded-lg">
					<div class="flex items-center space-x-4">
						<Bell class="w-6 h-6 text-gray-600" />
						<div>
							<h4 class="font-medium text-gray-900">Security Notifications</h4>
							<p class="text-sm text-gray-600">Get notified about security events</p>
						</div>
					</div>
					<div class="flex items-center space-x-2">
						<input type="checkbox" checked class="rounded" />
						<span class="text-sm">Enabled</span>
					</div>
				</div>
			</CardContent>
		</Card>
	{:else if activeTab === 'danger'}
		<Card class="border-red-200">
			<CardHeader>
				<div class="flex items-center space-x-2">
					<AlertTriangle class="w-5 h-5 text-red-500" />
					<CardTitle class="text-red-700">Danger Zone</CardTitle>
				</div>
			</CardHeader>
			<CardContent class="space-y-6">
				<div class="p-4 bg-red-50 border border-red-200 rounded-lg">
					<h4 class="font-medium text-red-800 mb-2">Delete Project</h4>
					<p class="text-sm text-red-700 mb-4">
						Once you delete a project, there is no going back. Please be certain. This action will
						permanently delete the project, all deployments, and associated data.
					</p>
					<Button variant="destructive" onclick={() => (showDeleteModal = true)}>
						<Trash2 class="w-4 h-4 mr-2" />
						Delete Project
					</Button>
				</div>
			</CardContent>
		</Card>
	{/if}
</div>

<!-- Delete Confirmation Modal -->
{#if showDeleteModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<Card class="w-full max-w-md">
			<CardHeader>
				<div class="flex items-center space-x-2">
					<AlertTriangle class="w-5 h-5 text-red-500" />
					<CardTitle class="text-red-700">Delete Project</CardTitle>
				</div>
			</CardHeader>
			<CardContent class="space-y-4">
				<p class="text-sm text-gray-600">
					This action cannot be undone. This will permanently delete the <strong
						>{project.name}</strong
					> project and all of its data.
				</p>
				<div>
					<Label for="delete-confirmation">Type <strong>{project.name}</strong> to confirm:</Label>
					<Input
						id="delete-confirmation"
						bind:value={deleteConfirmation}
						placeholder={project.name}
						class="mt-1"
					/>
				</div>
				<div class="flex space-x-2 pt-4">
					<Button
						variant="destructive"
						onclick={deleteProject}
						disabled={deleteConfirmation !== project.name}
						class="flex-1"
					>
						<Trash2 class="w-4 h-4 mr-2" />
						Delete Project
					</Button>
					<Button
						variant="outline"
						onclick={() => {
							showDeleteModal = false;
							deleteConfirmation = '';
						}}
						class="flex-1"
					>
						Cancel
					</Button>
				</div>
			</CardContent>
		</Card>
	</div>
{/if}

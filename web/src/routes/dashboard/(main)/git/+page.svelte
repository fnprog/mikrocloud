<script lang="ts">
	import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import {
		Sheet,
		SheetContent,
		SheetDescription,
		SheetHeader,
		SheetTitle
	} from '$lib/components/ui/sheet';
	import { Select, SelectContent, SelectItem, SelectTrigger } from '$lib/components/ui/select';
	import { Plus, Trash2, Github, GitBranch, Server } from 'lucide-svelte';
	import { gitApi, type GitProvider, type GitSource } from '$lib/api/git';

	let isCreateSheetOpen = $state(false);
	let isDeleteModalOpen = $state(false);
	let selectedSource = $state<GitSource | null>(null);

	let formData = $state({
		name: '',
		provider: 'github' as GitProvider,
		github_type: 'cloud' as 'cloud' | 'enterprise',
		access_token: '',
		refresh_token: '',
		custom_url: '',
		webhook_endpoint_type: 'ip' as 'ip' | 'domain',
		instance_ip: '127.0.0.1',
		instance_domain: '',
		allow_preview_deployments: false,
		manual_setup: false
	});

	const queryClient = useQueryClient();

	const sourcesQuery = createQuery(() => ({
		queryKey: ['git-sources'],
		queryFn: () => gitApi.listGitSources(),
		staleTime: 30 * 1000
	}));

	const createFn = createMutation(() => ({
		mutationFn: (data: {
			name: string;
			provider: GitProvider;
			access_token: string;
			refresh_token?: string;
			custom_url?: string;
		}) => gitApi.createGitSource(data),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['git-sources'] });
			isCreateSheetOpen = false;
			resetForm();
		}
	}));

	const deleteFn = createMutation(() => ({
		mutationFn: (sourceId: string) => gitApi.deleteGitSource(sourceId),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['git-sources'] });
			isDeleteModalOpen = false;
			selectedSource = null;
		}
	}));

	function createSource() {
		createFn.mutate({
			name: formData.name,
			provider: formData.provider,
			access_token: formData.access_token,
			refresh_token: formData.refresh_token || undefined,
			custom_url: formData.provider === 'custom' ? formData.custom_url : undefined
		});
	}

	function deleteSource() {
		if (!selectedSource) return;
		deleteMutation.mutate(selectedSource.id);
	}

	function resetForm() {
		formData = {
			name: '',
			provider: 'github',
			github_type: 'cloud',
			access_token: '',
			refresh_token: '',
			custom_url: '',
			webhook_endpoint_type: 'ip',
			instance_ip: '127.0.0.1',
			instance_domain: '',
			allow_preview_deployments: false,
			manual_setup: false
		};
	}

	function startOAuthRegistration() {
		const webhookUrl =
			formData.webhook_endpoint_type === 'domain' && formData.instance_domain
				? `https://${formData.instance_domain}/webhooks/git`
				: `http://${formData.instance_ip}:3000/webhooks/git`;

		const params = new URLSearchParams({
			name: formData.name,
			provider: formData.provider,
			webhook_url: webhookUrl,
			allow_preview: formData.allow_preview_deployments.toString()
		});

		if (
			formData.provider === 'github' &&
			formData.github_type === 'enterprise' &&
			formData.custom_url
		) {
			params.append('custom_url', formData.custom_url);
		}

		window.location.href = `/api/auth/git/oauth/start?${params.toString()}`;
	}

	function getProviderIcon(provider: GitProvider) {
		switch (provider) {
			case 'github':
				return Github;
			case 'gitlab':
			case 'bitbucket':
				return GitBranch;
			case 'custom':
				return Server;
			default:
				return GitBranch;
		}
	}

	function getProviderBadgeColor(provider: GitProvider) {
		switch (provider) {
			case 'github':
				return 'bg-gray-900 text-white';
			case 'gitlab':
				return 'bg-orange-600 text-white';
			case 'bitbucket':
				return 'bg-blue-600 text-white';
			case 'custom':
				return 'bg-purple-600 text-white';
			default:
				return 'bg-gray-500 text-white';
		}
	}

	function formatDate(dateStr: string) {
		return new Date(dateStr).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}
</script>

<svelte:head>
	<title>Git Sources - Dashboard</title>
</svelte:head>

<div class="flex-1 flex flex-col overflow-hidden">
	<div class=" px-6 py-4">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-2xl font-semibold">Git Sources</h1>
				<p class="text-sm text-muted-foreground mt-1">
					Manage Git provider connections and access tokens
				</p>
			</div>
			<Button onclick={() => (isCreateSheetOpen = true)}>
				<Plus class="h-4 w-4 mr-2" />
				Add Source
			</Button>
		</div>
	</div>

	<Sheet bind:open={isCreateSheetOpen}>
		<SheetContent class="overflow-y-auto sm:max-w-xl">
			<SheetHeader>
				<SheetTitle>Connect Git Provider</SheetTitle>
				<SheetDescription>
					Connect your Git provider to deploy applications from repositories
				</SheetDescription>
			</SheetHeader>
			<div class="space-y-6 py-6">
				<div class="space-y-2">
					<Label for="name">Source Name</Label>
					<Input id="name" placeholder="My GitHub Account" bind:value={formData.name} />
					<p class="text-xs text-muted-foreground">A friendly name to identify this connection</p>
				</div>

				<div class="space-y-2">
					<Label for="provider">Provider</Label>
					<Select value={formData.provider} onValueChange={(v) => v && (formData.provider = v)}>
						<SelectTrigger id="provider">
							{formData.provider === 'github'
								? 'GitHub'
								: formData.provider === 'gitlab'
									? 'GitLab'
									: formData.provider === 'bitbucket'
										? 'Bitbucket'
										: 'Custom/Self-hosted'}
						</SelectTrigger>
						<SelectContent>
							<SelectItem value="github">GitHub</SelectItem>
							<SelectItem value="gitlab">GitLab</SelectItem>
							<SelectItem value="bitbucket">Bitbucket</SelectItem>
							<SelectItem value="custom">Custom/Self-hosted</SelectItem>
						</SelectContent>
					</Select>
				</div>

				{#if formData.provider === 'github'}
					<div class="space-y-2">
						<Label for="github_type">GitHub Type</Label>
						<Select
							value={formData.github_type}
							onValueChange={(v) => v && (formData.github_type = v)}
						>
							<SelectTrigger id="github_type">
								{formData.github_type === 'cloud' ? 'GitHub Cloud' : 'GitHub Enterprise Server'}
							</SelectTrigger>
							<SelectContent>
								<SelectItem value="cloud">GitHub Cloud (github.com)</SelectItem>
								<SelectItem value="enterprise">GitHub Enterprise Server</SelectItem>
							</SelectContent>
						</Select>
					</div>

					{#if formData.github_type === 'enterprise'}
						<div class="space-y-2">
							<Label for="enterprise_url">Enterprise Server URL</Label>
							<Input
								id="enterprise_url"
								placeholder="https://github.company.com"
								bind:value={formData.custom_url}
							/>
							<p class="text-xs text-muted-foreground">URL of your GitHub Enterprise Server</p>
						</div>
					{/if}
				{:else if formData.provider === 'custom'}
					<div class="space-y-2">
						<Label for="custom_url">Git Server URL</Label>
						<Input
							id="custom_url"
							placeholder="https://git.example.com"
							bind:value={formData.custom_url}
						/>
						<p class="text-xs text-muted-foreground">Base URL of your self-hosted Git instance</p>
					</div>
				{/if}

				<div class="space-y-2">
					<Label>Webhook Endpoint</Label>
					<Select
						value={formData.webhook_endpoint_type}
						onValueChange={(v) => v && (formData.webhook_endpoint_type = v)}
					>
						<SelectTrigger>
							{formData.webhook_endpoint_type === 'ip' ? 'Use IP Address' : 'Use Domain Name'}
						</SelectTrigger>
						<SelectContent>
							<SelectItem value="ip">Use IP Address</SelectItem>
							<SelectItem value="domain" disabled={!formData.instance_domain}>
								Use Domain Name {!formData.instance_domain ? '(Configure in Settings)' : ''}
							</SelectItem>
						</SelectContent>
					</Select>
					<div class="text-xs text-muted-foreground bg-muted/50 p-3 rounded-md">
						<p class="font-medium mb-1">Webhook URL:</p>
						<code class="text-xs">
							{formData.webhook_endpoint_type === 'domain' && formData.instance_domain
								? `https://${formData.instance_domain}/webhooks/git`
								: formData.instance_ip
									? `http://${formData.instance_ip}:3000/webhooks/git`
									: 'Not configured'}
						</code>
						<p class="mt-2">All Git webhooks will be sent to this endpoint.</p>
						{#if formData.webhook_endpoint_type === 'domain' && !formData.instance_domain}
							<p class="text-yellow-600 mt-1">
								⚠️ Set your instance's FQDN in Settings to use domain name.
							</p>
						{/if}
					</div>
				</div>

				<div class="space-y-3">
					<div class="flex items-center space-x-2">
						<input
							type="checkbox"
							id="preview_deployments"
							bind:checked={formData.allow_preview_deployments}
							class="h-4 w-4 rounded border-gray-300"
						/>
						<Label for="preview_deployments" class="text-sm font-normal">
							Allow preview deployment permissions
						</Label>
					</div>
					<p class="text-xs text-muted-foreground pl-6">
						When enabled, deployments from pull requests will have access to environment variables
						and secrets.
					</p>
				</div>

				<div class="border-t pt-4">
					<h4 class="text-sm font-medium mb-2">Permissions</h4>
					<p class="text-xs text-muted-foreground mb-3">
						The following permissions will be requested:
					</p>
					<div class="space-y-1 text-xs text-muted-foreground pl-4">
						<div class="flex items-center gap-2">
							<div class="w-1 h-1 rounded-full bg-muted-foreground"></div>
							<span>Repository contents (read)</span>
						</div>
						<div class="flex items-center gap-2">
							<div class="w-1 h-1 rounded-full bg-muted-foreground"></div>
							<span>Metadata (read)</span>
						</div>
						<div class="flex items-center gap-2">
							<div class="w-1 h-1 rounded-full bg-muted-foreground"></div>
							<span>Webhooks (read & write)</span>
						</div>
						<div class="flex items-center gap-2">
							<div class="w-1 h-1 rounded-full bg-muted-foreground"></div>
							<span>Pull requests (read)</span>
						</div>
						<div class="flex items-center gap-2">
							<div class="w-1 h-1 rounded-full bg-muted-foreground"></div>
							<span>Email addresses (read)</span>
						</div>
					</div>
				</div>

				{#if !formData.manual_setup}
					<div class="space-y-3">
						<Button
							class="w-full"
							onclick={startOAuthRegistration}
							disabled={!formData.name || createFn.isPending}
						>
							{#if formData.provider === 'github'}
								<Github class="h-4 w-4 mr-2" />
								Register with GitHub
							{:else if formData.provider === 'gitlab'}
								<GitBranch class="h-4 w-4 mr-2" />
								Register with GitLab
							{:else if formData.provider === 'bitbucket'}
								<GitBranch class="h-4 w-4 mr-2" />
								Register with Bitbucket
							{:else}
								<Server class="h-4 w-4 mr-2" />
								Register with Git Provider
							{/if}
						</Button>
						<Button
							variant="ghost"
							class="w-full text-xs"
							onclick={() => (formData.manual_setup = true)}
						>
							Setup manually instead (Advanced)
						</Button>
					</div>
				{:else}
					<div class="space-y-4 border-t pt-4">
						<div class="bg-yellow-50 border border-yellow-200 p-3 rounded-md">
							<p class="text-xs text-yellow-800">
								<strong>Advanced:</strong> Manual setup requires you to create and configure the OAuth
								app or access token yourself.
							</p>
						</div>

						<div class="space-y-2">
							<Label for="access_token">Access Token / OAuth Token</Label>
							<Input
								id="access_token"
								type="password"
								placeholder="ghp_xxxxxxxxxxxx"
								bind:value={formData.access_token}
							/>
							<p class="text-xs text-muted-foreground">
								{#if formData.provider === 'github'}
									Create a personal access token or OAuth app at GitHub Settings → Developer
									settings
								{:else if formData.provider === 'gitlab'}
									Create an access token at GitLab Settings → Access Tokens
								{:else if formData.provider === 'bitbucket'}
									Create an app password at Bitbucket Settings → App passwords
								{:else}
									Create an access token in your Git provider's settings
								{/if}
							</p>
						</div>

						<div class="space-y-2">
							<Label for="refresh_token">Refresh Token (Optional)</Label>
							<Input
								id="refresh_token"
								type="password"
								placeholder="Optional refresh token"
								bind:value={formData.refresh_token}
							/>
							<p class="text-xs text-muted-foreground">If your provider supports token refresh</p>
						</div>

						<div class="flex gap-2">
							<Button
								variant="outline"
								class="flex-1"
								onclick={() => (formData.manual_setup = false)}
							>
								Back
							</Button>
							<Button
								class="flex-1"
								onclick={createSource}
								disabled={!formData.name || !formData.access_token || $createMutation.isPending}
							>
								Add Source
							</Button>
						</div>
					</div>
				{/if}
			</div>
			{#if !formData.manual_setup}
				<div class="flex justify-end">
					<Button variant="outline" onclick={() => (isCreateSheetOpen = false)}>Cancel</Button>
				</div>
			{/if}
		</SheetContent>
	</Sheet>

	<div class="flex-1 overflow-y-auto p-6">
		{#if createFn.error || deleteFn.error}
			<Card class="border-red-200 bg-red-50 mb-4">
				<CardContent class="pt-6">
					<p class="text-red-800">
						{createFn.error instanceof Error
							? createFn.error.message
							: deleteFn.error instanceof Error
								? deleteFn.error.message
								: 'An error occurred'}
					</p>
				</CardContent>
			</Card>
		{/if}

		{#if sourcesQuery.isLoading}
			<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
				{#each [1, 2, 3] as _}
					<Card>
						<CardHeader>
							<div class="h-4 bg-gray-200 rounded animate-pulse w-1/2"></div>
							<div class="h-3 bg-gray-200 rounded animate-pulse w-1/3 mt-2"></div>
						</CardHeader>
						<CardContent>
							<div class="h-3 bg-gray-200 rounded animate-pulse w-full"></div>
						</CardContent>
					</Card>
				{/each}
			</div>
		{:else if sourcesQuery.error}
			<Card class="border-red-200 bg-red-50">
				<CardContent class="pt-6">
					<p class="text-red-800">
						Error loading sources: {sourcesQuery.error instanceof Error
							? sourcesQuery.error.message
							: 'Unknown error'}
					</p>
				</CardContent>
			</Card>
		{:else if sourcesQuery.data && sourcesQuery.data.length === 0}
			<Card>
				<CardContent class="pt-6 text-center py-12">
					<GitBranch class="h-12 w-12 mx-auto text-gray-400 mb-4" />
					<h3 class="text-lg font-medium text-gray-900 mb-2">No Git sources yet</h3>
					<p class="text-gray-500 mb-4">
						Connect your Git providers to deploy applications from repositories
					</p>
					<Button onclick={() => (isCreateSheetOpen = true)}>
						<Plus class="h-4 w-4 mr-2" />
						Add Your First Source
					</Button>
				</CardContent>
			</Card>
		{:else if sourcesQuery.data}
			<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
				{#each sourcesQuery.data as source (source.id)}
					{@const Icon = getProviderIcon(source.provider)}
					<Card>
						<CardHeader>
							<div class="flex items-start justify-between">
								<div class="flex items-center gap-3">
									<div class="p-2 rounded-lg {getProviderBadgeColor(source.provider)}">
										<Icon class="h-5 w-5" />
									</div>
									<div>
										<CardTitle class="text-lg">{source.name}</CardTitle>
										<Badge variant="secondary" class="mt-1">
											{source.provider}
										</Badge>
									</div>
								</div>
								<Button
									variant="ghost"
									size="sm"
									onclick={() => {
										selectedSource = source;
										isDeleteModalOpen = true;
									}}
								>
									<Trash2 class="h-4 w-4 text-red-600" />
								</Button>
							</div>
						</CardHeader>
						<CardContent>
							{#if source.custom_url}
								<p class="text-sm text-gray-600 mb-2">
									<span class="font-medium">URL:</span>
									{source.custom_url}
								</p>
							{/if}
							<p class="text-xs text-gray-500">
								Created {formatDate(source.created_at)}
							</p>
						</CardContent>
					</Card>
				{/each}
			</div>
		{/if}
	</div>
</div>

{#if isDeleteModalOpen && selectedSource}
	<div
		class="fixed inset-0 z-50 bg-black/80 flex items-center justify-center"
		onclick={() => (isDeleteModalOpen = false)}
	>
		<div
			class="bg-white rounded-lg shadow-lg p-6 max-w-md w-full mx-4"
			onclick={(e) => e.stopPropagation()}
		>
			<h2 class="text-xl font-semibold mb-2">Delete Git Source</h2>
			<p class="text-gray-600 mb-6">
				Are you sure you want to delete "{selectedSource.name}"? This action cannot be undone.
			</p>
			<div class="flex justify-end gap-2">
				<Button variant="outline" onclick={() => (isDeleteModalOpen = false)}>Cancel</Button>
				<Button variant="destructive" onclick={deleteSource}>Delete</Button>
			</div>
		</div>
	</div>
{/if}

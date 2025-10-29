<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as InputGroup from '$lib/components/ui/input-group/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { RefreshCw, Check, AlertCircle, Lock, Github } from 'lucide-svelte';
	import { gitApi } from '$lib/features/git-sources/api';
	import { Search } from 'lucide-svelte';

	interface Props {
		source_repository_url: string;
		public_repository_url: string;
	}

	let { source_repository_url = $bindable(), public_repository_url = $bindable() }: Props =
		$props();

	let isValidating = $state(false);
	let validationStatus = $state<'idle' | 'valid' | 'invalid'>('idle');
	let validationMessage = $state('');

	let gitSources = $state<{ id: string; name: string; created_at?: string; provider?: string; custom_url?: string }[]>([]);
	let selectedSourceId = $state<string | undefined>(undefined);

	// When a source is selected from the dropdown, populate the parent-bound repository fields
	$effect(() => {
		if (!selectedSourceId) return;
		const src = gitSources.find((s) => s.id === selectedSourceId);
		if (src) {
			const repoUrl = src.custom_url ?? src.name;
			// update both possible bound props so the parent receives the selection
			source_repository_url = repoUrl;
			public_repository_url = repoUrl;
		}
	});

	const repositories = [
		{ id: '1', name: 'mikrocloud-web', isPrivate: true, date: 'Oct 18' },
		{ id: '2', name: 'trotroways', isPrivate: true, date: 'Oct 12' },
		{ id: '3', name: 'toffnet', isPrivate: true, date: 'Oct 11', icon: 'N' },
		{ id: '4', name: 'mikrocloud', isPrivate: false, date: 'Oct 8' },
		{ id: '5', name: 'vault', isPrivate: true, date: 'Sep 27' }
	];

	import { onMount } from 'svelte';

	onMount(async () => {
		try {
			const list = await gitApi.listGitSources();
			gitSources = list.map((s) => ({ id: s.id, name: s.name, created_at: s.created_at, provider: s.provider }));
			if (gitSources.length > 0 && !selectedSourceId) selectedSourceId = gitSources[0].id;
		} catch (err) {
			// ignore errors for now
		}
	});

	function detectProvider(url: string): 'github' | 'gitlab' | 'bitbucket' | 'custom' {
		const lower = (url || '').toLowerCase();
		if (lower.includes('github.com')) return 'github';
		if (lower.includes('gitlab.com')) return 'gitlab';
		if (lower.includes('bitbucket.org')) return 'bitbucket';
		// assume 'owner/repo' style or default to github
		if (lower.includes('/') && !lower.includes('.')) return 'github';
		return 'custom';
	}

	async function validateRepository() {
		if (!public_repository_url) return;

		isValidating = true;
		validationStatus = 'idle';

		try {
			const provider = detectProvider(public_repository_url);
			const result = await gitApi.validateRepository({
				provider,
				repository: public_repository_url
			});

			if (result.valid) {
				validationStatus = 'valid';
				validationMessage = result.message || 'Repository is accessible';
			} else {
				validationStatus = 'invalid';
				validationMessage = result.message || 'Unable to access repository';
			}
		} catch (error) {
			validationStatus = 'invalid';
			validationMessage = error instanceof Error ? error.message : 'Unable to access repository';
		} finally {
			isValidating = false;
		}
	}
</script>

<Tabs.Root value="source">
	<Tabs.List>
		<Tabs.Trigger value="source">Git Sources</Tabs.Trigger>
		<Tabs.Trigger value="public">Public Git repository</Tabs.Trigger>
	</Tabs.List>
	<Tabs.Content value="source">
		<div class="flex gap-4 my-4">
			<InputGroup.Root class="flex-2">
				<InputGroup.Input placeholder="Search" />
				<InputGroup.Addon>
					<Search />
				</InputGroup.Addon>
			</InputGroup.Root>

			<DropdownMenu.Root>
				<DropdownMenu.Trigger>
					{#snippet child({ props })}
						<Button {...props} variant="outline" class="w-full flex-1">Open</Button>
					{/snippet}
				</DropdownMenu.Trigger>
				<DropdownMenu.Content>
					<DropdownMenu.RadioGroup bind:value={selectedSourceId}>
						{#if gitSources.length > 0}
							{#each gitSources as source}
								<DropdownMenu.RadioItem value={source.id}>{source.name}</DropdownMenu.RadioItem>
							{/each}
						{:else}
							<DropdownMenu.RadioItem value="none">No sources</DropdownMenu.RadioItem>
						{/if}
					</DropdownMenu.RadioGroup>
					<DropdownMenu.Item>Add New source</DropdownMenu.Item>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</div>
		<div class="space-y-px border border-border rounded-lg overflow-hidden">
			{#if gitSources.length > 0}
				{#each gitSources as repo, index (repo.id)}
					<div
						class={`flex items-center justify-between px-6 py-4 bg-card hover:bg-muted transition-colors ${
							index !== gitSources.length - 1 ? 'border-b border-border' : ''
						}`}
					>
						<div class="flex items-center gap-3 flex-1">
							<div class="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
								<Github size={16} className="text-muted-foreground" />
							</div>
							<div class="flex items-center gap-2">
								<span class="text-foreground font-medium">{repo.name}</span>
								<span class="text-muted-foreground text-sm">{repo.created_at}</span>
							</div>
						</div>
						<Button variant="outline" size="sm">Import</Button>
					</div>
				{/each}
			{:else}
				{#each repositories as repo, index (repo.id)}
					<div
						class={`flex items-center justify-between px-6 py-4 bg-card hover:bg-muted transition-colors ${
							index !== repositories.length - 1 ? 'border-b border-border' : ''
						}`}
					>
						<div class="flex items-center gap-3 flex-1">
							{#if repo.icon}
								<div
									class="w-8 h-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-semibold text-sm"
								>
									{repo.icon}
								</div>
							{:else}
								<div class="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
									<Github size={16} className="text-muted-foreground" />
								</div>
							{/if}
							<div class="flex items-center gap-2">
								<span class="text-foreground font-medium">{repo.name}</span>
								{#if repo.isPrivate}
									<Lock size={14} class="text-muted-foreground" />
								{/if}
								<span class="text-muted-foreground text-sm">· {repo.date}</span>
							</div>
						</div>
						<Button variant="outline" size="sm">Import</Button>
					</div>
				{/each}
			{/if}
		</div>
	</Tabs.Content>
	<Tabs.Content value="public">
		<div class="space-y-2 mt-10">
			<Label for="repository">Repository URL</Label>
			<div class="flex gap-2">
				<div class="flex-1 relative">
					<Input
						id="repository"
						placeholder={'Leave empty if URL contains path'}
						bind:value={public_repository_url}
						class="pr-8"
					/>
					{#if validationStatus === 'valid'}
						<Check class="h-4 w-4 text-green-600 absolute right-2 top-1/2 -translate-y-1/2" />
					{:else if validationStatus === 'invalid'}
						<AlertCircle
							class="h-4 w-4 text-destructive absolute right-2 top-1/2 -translate-y-1/2"
						/>
					{/if}
				</div>
				<Button
					variant="outline"
					size="icon"
					onclick={validateRepository}
					disabled={isValidating || !public_repository_url}
				>
					<RefreshCw class="h-4 w-4 {isValidating ? 'animate-spin' : ''}" />
				</Button>
			</div>
			{#if validationMessage}
				<p
					class="text-xs"
					class:text-green-600={validationStatus === 'valid'}
					class:text-destructive={validationStatus === 'invalid'}
				>
					{validationMessage}
				</p>
			{:else}
				<p class="text-xs text-muted-foreground">
					Enter the repository path (e.g., username/repository-name)
				</p>
			{/if}
		</div>
	</Tabs.Content>
</Tabs.Root>

<!-- <div class="flex gap-5 space-x-5"> -->
<!-- 	<div class="space-y-6 flex-1"> -->
<!-- 		<h2 class="text-2xl font-semibold">Import From Git Source</h2> -->


<!-- 		<div class="space-y-2"> -->
<!-- 			<div class="flex items-center justify-between"> -->
<!-- 				<Label for="branch">Default branch</Label> -->
<!-- 				<Button -->
<!-- 					variant="ghost" -->
<!-- 					size="sm" -->
<!-- 					onclick={fetchBranches} -->
<!-- 					disabled={isFetchingBranches || !repository} -->
<!-- 					class="h-7 text-xs" -->
<!-- 				> -->
<!-- 					<RefreshCw class="h-3 w-3 mr-1 {isFetchingBranches ? 'animate-spin' : ''}" /> -->
<!-- 					Fetch branches -->
<!-- 				</Button> -->
<!-- 			</div> -->
<!-- 			<Select type="single" bind:value={branch}> -->
<!-- 				<SelectTrigger id="branch"> -->
<!-- 					{branch || 'Select a branch'} -->
<!-- 				</SelectTrigger> -->
<!-- 				<SelectContent> -->
<!-- 					{#each branches as branchOption} -->
<!-- 						<SelectItem value={branchOption}>{branchOption}</SelectItem> -->
<!-- 					{/each} -->
<!-- 				</SelectContent> -->
<!-- 			</Select> -->
<!-- 		</div> -->
<!-- 	</div> -->
<!-- 	<Separator orientation="vertical" /> -->

<!-- 	<div class="space-y-6 flex-1 mt-8"> -->
<!-- 		<h2 class="text-2xl font-semibold">Import From Public Repository</h2> -->


<!-- 		<div class="space-y-2"> -->
<!-- 			<div class="flex items-center justify-between"> -->
<!-- 				<Label for="branch">Default branch</Label> -->
<!-- 				<Button -->
<!-- 					variant="ghost" -->
<!-- 					size="sm" -->
<!-- 					onclick={fetchBranches} -->
<!-- 					disabled={isFetchingBranches || !repository} -->
<!-- 					class="h-7 text-xs" -->
<!-- 				> -->
<!-- 					<RefreshCw class="h-3 w-3 mr-1 {isFetchingBranches ? 'animate-spin' : ''}" /> -->
<!-- 					Fetch branches -->
<!-- 				</Button> -->
<!-- 			</div> -->
<!-- 			<Select type="single" bind:value={branch}> -->
<!-- 				<SelectTrigger id="branch"> -->
<!-- 					{branch || 'Select a branch'} -->
<!-- 				</SelectTrigger> -->
<!-- 				<SelectContent> -->
<!-- 					{#each branches as branchOption} -->
<!-- 						<SelectItem value={branchOption}>{branchOption}</SelectItem> -->
<!-- 					{/each} -->
<!-- 				</SelectContent> -->
<!-- 			</Select> -->
<!-- 		</div> -->
<!-- 	</div> -->
<!-- </div> -->

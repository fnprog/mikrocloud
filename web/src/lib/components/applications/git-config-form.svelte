<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import { Tabs, TabsList, TabsTrigger } from '$lib/components/ui/tabs';
	import { RefreshCw } from 'lucide-svelte';
	import {
		Select,
		SelectContent,
		SelectItem,
		SelectTrigger
	} from '$lib/components/ui/select/index';

	interface Props {
		provider: 'github' | 'gitlab' | 'bitbucket';
		onProviderChange: (provider: 'github' | 'gitlab' | 'bitbucket') => void;
		repository: string;
		onRepositoryChange: (repository: string) => void;
		branch: string;
		onBranchChange: (branch: string) => void;
		autoDeploy: boolean;
		onAutoDeployChange: (autoDeploy: boolean) => void;
		isPrivate: boolean;
		onIsPrivateChange: (isPrivate: boolean) => void;
	}

	let {
		provider,
		onProviderChange,
		repository,
		onRepositoryChange,
		branch,
		onBranchChange,
		autoDeploy,
		onAutoDeployChange,
		isPrivate,
		onIsPrivateChange
	}: Props = $props();

	const branches = ['main', 'master', 'develop', 'staging', 'production'];
</script>

<div class="space-y-6">
	<div class="flex gap-4">
		<Button
			variant={!isPrivate ? 'default' : 'outline'}
			onclick={() => onIsPrivateChange(false)}
			class="flex-1"
		>
			Public repository
		</Button>
		<Button
			variant={isPrivate ? 'default' : 'outline'}
			onclick={() => onIsPrivateChange(true)}
			class="flex-1"
		>
			Private repository
		</Button>
	</div>

	<Tabs
		value={provider}
		onValueChange={(v) => v && onProviderChange(v as 'github' | 'gitlab' | 'bitbucket')}
	>
		<TabsList class="grid w-full grid-cols-3">
			<TabsTrigger value="github">GitHub</TabsTrigger>
			<TabsTrigger value="gitlab">GitLab</TabsTrigger>
			<TabsTrigger value="bitbucket">Bitbucket</TabsTrigger>
		</TabsList>
	</Tabs>

	<div class="space-y-2">
		<Label for="repository">Repository</Label>
		<div class="flex gap-2">
			<Input
				id="repository"
				placeholder="username/repo-name"
				value={repository}
				oninput={(e) => onRepositoryChange(e.currentTarget.value)}
				class="flex-1"
			/>
			<Button variant="outline" size="icon">
				<RefreshCw class="h-4 w-4" />
			</Button>
		</div>
		<p class="text-xs text-muted-foreground">
			Enter the repository path (e.g., username/repository-name)
		</p>
	</div>

	<div class="space-y-2">
		<Label for="branch">Default branch</Label>
		<Select value={branch} onValueChange={(v) => v && onBranchChange(v)}>
			<SelectTrigger id="branch">
				{branch || 'Select a branch'}
			</SelectTrigger>
			<SelectContent>
				{#each branches as branchOption}
					<SelectItem value={branchOption}>{branchOption}</SelectItem>
				{/each}
			</SelectContent>
		</Select>
	</div>

	<div class="flex items-center space-x-2">
		<input
			type="checkbox"
			id="auto-deploy"
			checked={autoDeploy}
			onchange={(e) => onAutoDeployChange(e.currentTarget.checked)}
			class="h-4 w-4 rounded border-gray-300"
		/>
		<Label for="auto-deploy" class="font-normal cursor-pointer">
			Automatic deployment on commit
		</Label>
	</div>
</div>

<script lang="ts">
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Label } from '$lib/components/ui/label';
	import { Card } from '$lib/components/ui/card';
	import { Switch } from '$lib/components/ui/switch';
	import { Badge } from '$lib/components/ui/badge';
	import { Dialog, DialogFooter, DialogContent, DialogTitle, DialogHeader, DialogDescription } from '$lib/components/ui/dialog';
	import * as Table from '$lib/components/ui/table';
	import {
		Plus,
		Trash2,
		Copy,
		Code,
		Eye,
		EyeOff,
		Upload,
		Download,
		AlertCircle
	} from 'lucide-svelte';
	import { Alert, AlertTitle, AlertDescription } from '$lib/components/ui/alert';

	const projectId = $derived(page.params.id!);

	interface EnvVar {
		id: string;
		key: string;
		value: string;
		is_multiline: boolean;
		is_secret: boolean;
		created_at: string;
	}

	let envVars = $state<EnvVar[]>([
		{
			id: '1',
			key: 'DATABASE_URL',
			value: 'postgresql://user:pass@localhost:5432/db',
			is_multiline: false,
			is_secret: true,
			created_at: new Date().toISOString()
		},
		{
			id: '2',
			key: 'API_CONFIG',
			value: '{\n  "timeout": 30000,\n  "retries": 3\n}',
			is_multiline: true,
			is_secret: false,
			created_at: new Date().toISOString()
		}
	]);

	let showAddModal = $state(false);
	let showBulkModal = $state(false);
	let newVarKey = $state('');
	let newVarValue = $state('');
	let newVarIsMultiline = $state(false);
	let newVarIsSecret = $state(false);
	let bulkEnvContent = $state('');
	let visibleSecrets = $state<Set<string>>(new Set());
	let searchQuery = $state('');

	const filteredVars = $derived(
		envVars.filter((v) => v.key.toLowerCase().includes(searchQuery.toLowerCase()))
	);

	function handleAddVariable() {
		if (!newVarKey.trim() || !newVarValue.trim()) {
			alert('Key and value are required');
			return;
		}

		const newVar: EnvVar = {
			id: Math.random().toString(36).substring(7),
			key: newVarKey.trim(),
			value: newVarValue,
			is_multiline: newVarIsMultiline,
			is_secret: newVarIsSecret,
			created_at: new Date().toISOString()
		};

		envVars.push(newVar);
		resetAddForm();
		showAddModal = false;
	}

	function handleBulkAdd() {
		if (!bulkEnvContent.trim()) {
			alert('Please provide environment variables');
			return;
		}

		const lines = bulkEnvContent.split('\n').filter((line) => line.trim() && !line.startsWith('#'));

		lines.forEach((line) => {
			const index = line.indexOf('=');
			if (index > 0) {
				const key = line.substring(0, index).trim();
				const value = line.substring(index + 1).trim();

				const existing = envVars.find((v) => v.key === key);
				if (existing) {
					existing.value = value;
				} else {
					envVars.push({
						id: Math.random().toString(36).substring(7),
						key,
						value,
						is_multiline: value.includes('\n'),
						is_secret: false,
						created_at: new Date().toISOString()
					});
				}
			}
		});

		bulkEnvContent = '';
		showBulkModal = false;
	}

	function handleDeleteVar(id: string) {
		if (confirm('Are you sure you want to delete this variable?')) {
			envVars = envVars.filter((v) => v.id !== id);
		}
	}

	function toggleSecretVisibility(id: string) {
		if (visibleSecrets.has(id)) {
			visibleSecrets.delete(id);
		} else {
			visibleSecrets.add(id);
		}
		visibleSecrets = visibleSecrets;
	}

	function copyValue(value: string) {
		navigator.clipboard.writeText(value);
	}

	function resetAddForm() {
		newVarKey = '';
		newVarValue = '';
		newVarIsMultiline = false;
		newVarIsSecret = false;
	}

	function exportEnvFile() {
		const content = envVars.map((v) => `${v.key}=${v.value}`).join('\n');
		const blob = new Blob([content], { type: 'text/plain' });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = '.env';
		a.click();
		URL.revokeObjectURL(url);
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="font-bold text-3xl">Shared Environment Variables</h1>
		<p class="text-muted-foreground mt-1">
			Manage project-level environment variables accessible across all environments
		</p>
	</div>

	<Alert>
		<AlertCircle class="size-4" />
		<AlertTitle>Variable Syntax</AlertTitle>
		<AlertDescription>
			Use <code class="bg-muted px-1.5 py-0.5 rounded font-mono text-xs"
				>{'{{'} environment.VARIABLENAME {'}}'}</code
			> to reference these variables in your applications.
		</AlertDescription>
	</Alert>

	<div class="flex items-center justify-between gap-4">
		<div class="relative flex-1 max-w-md">
			<Input
				type="text"
				placeholder="Search variables..."
				bind:value={searchQuery}
				class="w-full"
			/>
		</div>
		<div class="flex gap-2">
			<Button variant="outline" onclick={exportEnvFile}>
				<Download class="size-4" />
				Export
			</Button>
			<Button variant="outline" onclick={() => (showBulkModal = true)}>
				<Upload class="size-4" />
				Bulk Import
			</Button>
			<Button onclick={() => (showAddModal = true)}>
				<Plus class="size-4" />
				Add Variable
			</Button>
		</div>
	</div>

	{#if filteredVars.length === 0}
		<Card class="p-12">
			<div class="text-center">
				<Code class="size-12 mx-auto text-muted-foreground mb-4" />
				<h3 class="font-semibold text-lg mb-2">No environment variables</h3>
				<p class="text-muted-foreground text-sm mb-4">
					{searchQuery
						? 'No variables match your search'
						: 'Add your first environment variable to get started'}
				</p>
				{#if !searchQuery}
					<Button onclick={() => (showAddModal = true)}>
						<Plus class="size-4" />
						Add Variable
					</Button>
				{/if}
			</div>
		</Card>
	{:else}
		<Card>
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head class="w-[250px]">Key</Table.Head>
						<Table.Head>Value</Table.Head>
						<Table.Head class="w-[100px]">Type</Table.Head>
						<Table.Head class="w-[120px] text-right">Actions</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each filteredVars as envVar (envVar.id)}
						<Table.Row>
							<Table.Cell class="font-mono text-sm">{envVar.key}</Table.Cell>
							<Table.Cell>
								<div class="flex items-center gap-2">
									{#if envVar.is_secret && !visibleSecrets.has(envVar.id)}
										<code class="font-mono text-xs text-muted-foreground">••••••••</code>
									{:else if envVar.is_multiline}
										<code
											class="font-mono text-xs bg-muted px-2 py-1 rounded block max-w-md overflow-hidden text-ellipsis whitespace-nowrap"
										>
											{envVar.value}
										</code>
									{:else}
										<code
											class="font-mono text-xs bg-muted px-2 py-1 rounded max-w-md overflow-hidden text-ellipsis whitespace-nowrap block"
										>
											{envVar.value}
										</code>
									{/if}
									{#if envVar.is_secret}
										<Button
											variant="ghost"
											size="sm"
											class="h-6 w-6 p-0"
											onclick={() => toggleSecretVisibility(envVar.id)}
										>
											{#if visibleSecrets.has(envVar.id)}
												<EyeOff class="size-3" />
											{:else}
												<Eye class="size-3" />
											{/if}
										</Button>
									{/if}
								</div>
							</Table.Cell>
							<Table.Cell>
								<div class="flex gap-1">
									{#if envVar.is_multiline}
										<Badge variant="secondary" class="text-xs">Multiline</Badge>
									{/if}
									{#if envVar.is_secret}
										<Badge variant="outline" class="text-xs">Secret</Badge>
									{/if}
								</div>
							</Table.Cell>
							<Table.Cell class="text-right">
								<div class="flex justify-end gap-1">
									<Button
										variant="ghost"
										size="sm"
										class="h-8 w-8 p-0"
										onclick={() => copyValue(envVar.value)}
									>
										<Copy class="size-4" />
									</Button>
									<Button
										variant="ghost"
										size="sm"
										class="h-8 w-8 p-0 text-destructive hover:text-destructive"
										onclick={() => handleDeleteVar(envVar.id)}
									>
										<Trash2 class="size-4" />
									</Button>
								</div>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</Card>
	{/if}
</div>

<Dialog bind:open={showAddModal}>
	<DialogContent class="sm:max-w-[500px]">
		<DialogHeader>
			<DialogTitle>Add Environment Variable</DialogTitle>
			<DialogDescription>
				Add a new shared environment variable for this project
			</DialogDescription>
		</DialogHeader>

		<div class="space-y-4">
			<div class="space-y-2">
				<Label for="key">Key</Label>
				<Input id="key" bind:value={newVarKey} placeholder="DATABASE_URL" />
			</div>

			<div class="space-y-2">
				<Label for="value">Value</Label>
				{#if newVarIsMultiline}
					<Textarea
						id="value"
						bind:value={newVarValue}
						placeholder="Enter value..."
						rows={6}
						class="font-mono text-sm"
					/>
				{:else}
					<Input
						id="value"
						bind:value={newVarValue}
						placeholder="Enter value..."
						class="font-mono text-sm"
					/>
				{/if}
			</div>

			<div class="flex items-center justify-between">
				<Label for="multiline" class="text-sm font-normal">Multiline value</Label>
				<Switch id="multiline" bind:checked={newVarIsMultiline} />
			</div>

			<div class="flex items-center justify-between">
				<Label for="secret" class="text-sm font-normal">Mark as secret</Label>
				<Switch id="secret" bind:checked={newVarIsSecret} />
			</div>
		</div>

		<DialogFooter>
			<Button variant="outline" onclick={() => (showAddModal = false)}>Cancel</Button>
			<Button onclick={handleAddVariable}>Add Variable</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>

<Dialog bind:open={showBulkModal}>
	<DialogContent class="sm:max-w-[600px]">
		<DialogHeader>
			<DialogTitle>Bulk Import Environment Variables</DialogTitle>
			<DialogDescription>
				Paste your .env file content. One variable per line in KEY=VALUE format
			</DialogDescription>
		</DialogHeader>

		<div class="space-y-4">
			<Textarea
				bind:value={bulkEnvContent}
				placeholder={'DATABASE_URL=postgresql://...\nAPI_KEY=your_api_key\nPORT=3000'}
				rows={12}
				class="font-mono text-sm"
			/>
			<p class="text-muted-foreground text-xs">
				Lines starting with # will be ignored. Existing variables will be updated.
			</p>
		</div>

		<DialogFooter>
			<Button variant="outline" onclick={() => (showBulkModal = false)}>Cancel</Button>
			<Button onclick={handleBulkAdd}>Import Variables</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>

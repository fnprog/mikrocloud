<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Textarea } from '$lib/components/ui/textarea';
	import * as Card from '$lib/components/ui/card';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Table from '$lib/components/ui/table';
	import Empty from '$lib/components/ui/empty/empty.svelte';
	import { Plus, Trash2, Database } from 'lucide-svelte';

	type BackupDestination = {
		id: string;
		name: string;
		description: string;
		endpoint: string;
		bucket: string;
		region: string;
		access_key: string;
		secret_key: string;
	};

	let destinations = $state<BackupDestination[]>([]);
	let isDialogOpen = $state(false);

	let name = $state('');
	let description = $state('');
	let endpoint = $state('');
	let bucket = $state('');
	let region = $state('');
	let accessKey = $state('');
	let secretKey = $state('');
	let isValidating = $state(false);
	let validationMessage = $state<{ type: 'success' | 'error'; text: string } | null>(null);

	onMount(async () => {
		await loadDestinations();
	});

	async function loadDestinations() {
		try {
			const response = await fetch('/api/settings/backup-locations');
			if (response.ok) {
				const data = await response.json();
				destinations = data || [];
			}
		} catch (error) {
			console.error('Failed to load backup destinations:', error);
		}
	}

	function resetForm() {
		name = '';
		description = '';
		endpoint = '';
		bucket = '';
		region = '';
		accessKey = '';
		secretKey = '';
		validationMessage = null;
	}

	async function handleValidate() {
		isValidating = true;
		validationMessage = null;

		try {
			const response = await fetch('/api/settings/backup-locations/validate', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					endpoint,
					bucket,
					region,
					access_key: accessKey,
					secret_key: secretKey
				})
			});

			if (response.ok) {
				validationMessage = { type: 'success', text: 'Connection validated successfully!' };
			} else {
				const error = await response.json();
				validationMessage = {
					type: 'error',
					text: error.message || 'Failed to validate connection'
				};
			}
		} catch (error) {
			validationMessage = {
				type: 'error',
				text: 'Failed to validate connection. Please check your settings.'
			};
		} finally {
			isValidating = false;
		}
	}

	async function handleCreate() {
		try {
			const response = await fetch('/api/settings/backup-locations', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					name,
					description,
					endpoint,
					bucket,
					region,
					access_key: accessKey,
					secret_key: secretKey
				})
			});

			if (!response.ok) {
				throw new Error('Failed to create backup destination');
			}

			await loadDestinations();
			resetForm();
			isDialogOpen = false;
		} catch (error) {
			console.error('Failed to create backup destination:', error);
		}
	}

	async function handleDelete(id: string) {
		try {
			const response = await fetch(`/api/settings/backup-locations/${id}`, {
				method: 'DELETE'
			});

			if (!response.ok) {
				throw new Error('Failed to delete backup destination');
			}

			destinations = destinations.filter((d) => d.id !== id);
		} catch (error) {
			console.error('Failed to delete backup destination:', error);
		}
	}
</script>

<div class="space-y-6">
	<Card.Root>
		<Card.Header>
			<div class="flex items-center justify-between">
				<div>
					<Card.Title>S3 Backup Destinations</Card.Title>
					<Card.Description>
						Configure S3-compatible storage for automatic backup uploads.
					</Card.Description>
				</div>
				<Dialog.Root bind:open={isDialogOpen}>
					<Dialog.Trigger>
						{#snippet child({ props })}
							<Button onclick={() => resetForm()} {...props}>
								<Plus class="h-4 w-4 mr-2" />
								Add Destination
							</Button>
						{/snippet}
					</Dialog.Trigger>
					<Dialog.Portal>
						<Dialog.Overlay />
						<Dialog.Content class="max-w-2xl max-h-[90vh] overflow-y-auto">
							<Dialog.Header>
								<Dialog.Title>Create S3 Backup Destination</Dialog.Title>
								<Dialog.Description>
									Configure a new S3-compatible storage destination for backups.
								</Dialog.Description>
							</Dialog.Header>

							<div class="space-y-4 py-4">
								<div class="space-y-2">
									<Label for="name">Name</Label>
									<Input id="name" bind:value={name} placeholder="My S3 Backup" />
									<p class="text-xs text-muted-foreground">
										A friendly name for this backup destination.
									</p>
								</div>

								<div class="space-y-2">
									<Label for="description">Description</Label>
									<Textarea
										id="description"
										bind:value={description}
										placeholder="Primary backup storage location"
										rows={3}
									/>
									<p class="text-xs text-muted-foreground">
										Optional description for this destination.
									</p>
								</div>

								<div class="space-y-2">
									<Label for="endpoint">Endpoint</Label>
									<Input
										id="endpoint"
										bind:value={endpoint}
										placeholder="https://s3.amazonaws.com"
									/>
									<p class="text-xs text-muted-foreground">
										S3-compatible endpoint URL (e.g., s3.amazonaws.com, minio.example.com).
									</p>
								</div>

								<div class="grid grid-cols-2 gap-4">
									<div class="space-y-2">
										<Label for="bucket">Bucket</Label>
										<Input id="bucket" bind:value={bucket} placeholder="my-backup-bucket" />
										<p class="text-xs text-muted-foreground">S3 bucket name.</p>
									</div>

									<div class="space-y-2">
										<Label for="region">Region</Label>
										<Input id="region" bind:value={region} placeholder="us-east-1" />
										<p class="text-xs text-muted-foreground">S3 bucket region.</p>
									</div>
								</div>

								<div class="space-y-2">
									<Label for="access-key">Access Key</Label>
									<Input
										id="access-key"
										bind:value={accessKey}
										placeholder="AKIAIOSFODNN7EXAMPLE"
									/>
									<p class="text-xs text-muted-foreground">S3 access key ID.</p>
								</div>

								<div class="space-y-2">
									<Label for="secret-key">Secret Key</Label>
									<Input
										id="secret-key"
										type="password"
										bind:value={secretKey}
										placeholder="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
									/>
									<p class="text-xs text-muted-foreground">S3 secret access key.</p>
								</div>

								{#if validationMessage}
									{@const isSuccess = validationMessage.type === 'success'}
									<div
										class="rounded-lg p-4 border"
										class:bg-green-50={isSuccess}
										class:border-green-200={isSuccess}
										class:dark:bg-green-950={isSuccess}
										class:dark:border-green-900={isSuccess}
										class:bg-red-50={!isSuccess}
										class:border-red-200={!isSuccess}
										class:dark:bg-red-950={!isSuccess}
										class:dark:border-red-900={!isSuccess}
									>
										<p
											class="text-sm font-medium"
											class:text-green-900={isSuccess}
											class:dark:text-green-100={isSuccess}
											class:text-red-900={!isSuccess}
											class:dark:text-red-100={!isSuccess}
										>
											{isSuccess ? '✓' : '✗'}
											{validationMessage.text}
										</p>
									</div>
								{/if}
							</div>

							<Dialog.Footer>
								<Button variant="outline" onclick={handleValidate} disabled={isValidating}>
									{isValidating ? 'Validating...' : 'Validate Connection'}
								</Button>
								<Button onclick={handleCreate}>Create Destination</Button>
							</Dialog.Footer>
						</Dialog.Content>
					</Dialog.Portal>
				</Dialog.Root>
			</div>
		</Card.Header>
		<Card.Content>
			{#if destinations.length === 0}
				<Empty class="border">
					<Database class="h-12 w-12 text-muted-foreground/50" />
					<div class="space-y-2">
						<h3 class="text-lg font-semibold">No backup destinations</h3>
						<p class="text-sm text-muted-foreground">
							Add your first S3-compatible storage destination to enable automatic backups.
						</p>
					</div>
				</Empty>
			{:else}
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>Name</Table.Head>
							<Table.Head>Endpoint</Table.Head>
							<Table.Head>Bucket</Table.Head>
							<Table.Head>Region</Table.Head>
							<Table.Head class="text-right">Actions</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each destinations as destination}
							<Table.Row>
								<Table.Cell class="font-medium">
									<div>
										<p>{destination.name}</p>
										{#if destination.description}
											<p class="text-xs text-muted-foreground">{destination.description}</p>
										{/if}
									</div>
								</Table.Cell>
								<Table.Cell>{destination.endpoint}</Table.Cell>
								<Table.Cell>{destination.bucket}</Table.Cell>
								<Table.Cell>{destination.region}</Table.Cell>
								<Table.Cell class="text-right">
									<Button variant="ghost" size="sm" onclick={() => handleDelete(destination.id)}>
										<Trash2 class="h-4 w-4" />
									</Button>
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			{/if}
		</Card.Content>
	</Card.Root>
</div>

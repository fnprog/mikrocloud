<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Label } from '$lib/components/ui/label';
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Select from '$lib/components/ui/select';
	import * as Alert from '$lib/components/ui/alert';
	import { AlertTriangle, Archive, Copy, Trash2, Users } from 'lucide-svelte';
	import { createProjectQuery } from '$lib/features/projects/queries';

	const projectId = $derived(page.params.id!);
	const projectQuery = $derived(createProjectQuery(projectId));

	let projectName = $state('');
	let projectDescription = $state('');

	let showDeleteModal = $state(false);
	let showTransferModal = $state(false);
	let showArchiveModal = $state(false);
	let deleteConfirmation = $state('');

	$effect(() => {
		if (projectQuery.data) {
			projectName = projectQuery.data.name;
			projectDescription = projectQuery.data.description || '';
		}
	});

	function handleUpdateGeneral() {
		console.log('Updating general settings:', { projectName, projectDescription });
	}

	function handleDeleteProject() {
		if (deleteConfirmation !== projectName) {
			alert('Project name does not match');
			return;
		}
		console.log('Deleting project:', projectId);
		goto('/dashboard/projects');
	}

	function handleTransferProject() {
		console.log('Transferring project:', projectId);
		showTransferModal = false;
	}

	function handleArchiveProject() {
		console.log('Archiving project:', projectId);
		showArchiveModal = false;
	}

	function handleDuplicateProject() {
		console.log('Duplicating project:', projectId);
	}
</script>

<div class="flex flex-col gap-6 p-6">
	<div>
		<h1 class="font-bold text-3xl">Project Settings</h1>
		<p class="text-muted-foreground mt-1">Manage your project configuration and preferences</p>
	</div>

	<Card.Root>
		<Card.Header>
			<Card.Title>General</Card.Title>
			<Card.Description>Basic project information and metadata</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			<div class="space-y-2">
				<Label for="name">Project Name</Label>
				<Input id="name" bind:value={projectName} placeholder="My Awesome Project" />
			</div>

			<div class="space-y-2">
				<Label for="description">Description</Label>
				<Textarea
					id="description"
					bind:value={projectDescription}
					placeholder="A brief description of your project..."
					rows={3}
				/>
			</div>

			<div class="flex items-center gap-2">
				<Badge variant="secondary">
					<Copy class="size-3 mr-1" />
					ID: {projectId}
				</Badge>
				<Badge variant="outline">
					Created: {projectQuery.data?.created_at
						? new Date(projectQuery.data.created_at).toLocaleDateString()
						: 'N/A'}
				</Badge>
			</div>
		</Card.Content>
		<Card.Footer>
			<Button onclick={handleUpdateGeneral}>Save Changes</Button>
		</Card.Footer>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title>Project Actions</Card.Title>
			<Card.Description>Manage project lifecycle and ownership</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-3">
			<div class="flex items-center justify-between p-3 border rounded-lg">
				<div class="flex items-center gap-3">
					<div class="p-2 bg-blue-100 rounded-lg">
						<Copy class="size-4 text-blue-600" />
					</div>
					<div>
						<p class="font-medium text-sm">Duplicate Project</p>
						<p class="text-muted-foreground text-xs">Create a copy with all configurations</p>
					</div>
				</div>
				<Button variant="outline" onclick={handleDuplicateProject}>Duplicate</Button>
			</div>

			<div class="flex items-center justify-between p-3 border rounded-lg">
				<div class="flex items-center gap-3">
					<div class="p-2 bg-purple-100 rounded-lg">
						<Users class="size-4 text-purple-600" />
					</div>
					<div>
						<p class="font-medium text-sm">Transfer Ownership</p>
						<p class="text-muted-foreground text-xs">Move project to another organization</p>
					</div>
				</div>
				<Button variant="outline" onclick={() => (showTransferModal = true)}>Transfer</Button>
			</div>
		</Card.Content>
	</Card.Root>

	<Card.Root class="border-destructive">
		<Card.Header>
			<Card.Title class="text-destructive">Danger Zone</Card.Title>
			<Card.Description>Irreversible actions that will affect your project</Card.Description>
		</Card.Header>
		<Card.Content>
			<Alert.Root variant="destructive">
				<AlertTriangle class="size-4" />
				<Alert.Title>Delete Project</Alert.Title>
				<Alert.Description>
					Permanently delete this project and all its resources. This action cannot be undone.
				</Alert.Description>
			</Alert.Root>
		</Card.Content>
		<Card.Footer>
			<Button variant="destructive" onclick={() => (showDeleteModal = true)}>
				<Trash2 class="size-4" />
				Delete Project
			</Button>
		</Card.Footer>
	</Card.Root>
</div>

<Dialog.Root bind:open={showDeleteModal}>
	<Dialog.Content class="sm:max-w-[500px]">
		<Dialog.Header>
			<Dialog.Title>Delete Project</Dialog.Title>
			<Dialog.Description>
				This action cannot be undone. All resources, deployments, and data will be permanently
				deleted.
			</Dialog.Description>
		</Dialog.Header>

		<div class="space-y-4">
			<Alert.Root variant="destructive">
				<AlertTriangle class="size-4" />
				<Alert.Title>Warning</Alert.Title>
				<Alert.Description>
					Type <strong>{projectName}</strong> to confirm deletion
				</Alert.Description>
			</Alert.Root>

			<div class="space-y-2">
				<Label for="confirm">Project Name</Label>
				<Input
					id="confirm"
					bind:value={deleteConfirmation}
					placeholder={projectName}
					class="font-mono"
				/>
			</div>
		</div>

		<Dialog.Footer>
			<Button variant="outline" onclick={() => (showDeleteModal = false)}>Cancel</Button>
			<Button
				variant="destructive"
				onclick={handleDeleteProject}
				disabled={deleteConfirmation !== projectName}
			>
				Delete Project
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<Dialog.Root bind:open={showTransferModal}>
	<Dialog.Content class="sm:max-w-[500px]">
		<Dialog.Header>
			<Dialog.Title>Transfer Project</Dialog.Title>
			<Dialog.Description>Transfer this project to another organization or user</Dialog.Description>
		</Dialog.Header>

		<div class="space-y-4">
			<div class="space-y-2">
				<Label for="target-org">Target Organization</Label>
				<Select.Root type="single">
					<Select.Trigger>Select organization</Select.Trigger>
					<Select.Content>
						<Select.Item value="org1">Organization 1</Select.Item>
						<Select.Item value="org2">Organization 2</Select.Item>
						<Select.Item value="personal">Personal Account</Select.Item>
					</Select.Content>
				</Select.Root>
			</div>

			<Alert.Root>
				<AlertTriangle class="size-4" />
				<Alert.Description>
					You will lose access to this project after the transfer is complete.
				</Alert.Description>
			</Alert.Root>
		</div>

		<Dialog.Footer>
			<Button variant="outline" onclick={() => (showTransferModal = false)}>Cancel</Button>
			<Button onclick={handleTransferProject}>Transfer Project</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<Dialog.Root bind:open={showArchiveModal}>
	<Dialog.Content class="sm:max-w-[500px]">
		<Dialog.Header>
			<Dialog.Title>Archive Project</Dialog.Title>
			<Dialog.Description>
				Make this project read-only and stop all active deployments
			</Dialog.Description>
		</Dialog.Header>

		<Alert.Root>
			<Archive class="size-4" />
			<Alert.Description>
				Archived projects can be restored later, but all deployments will be stopped.
			</Alert.Description>
		</Alert.Root>

		<Dialog.Footer>
			<Button variant="outline" onclick={() => (showArchiveModal = false)}>Cancel</Button>
			<Button onclick={handleArchiveProject}>Archive Project</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Label } from '$lib/components/ui/label';
	import * as Card from '$lib/components/ui/card';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Alert from '$lib/components/ui/alert';
	import { AlertTriangle, Trash2 } from 'lucide-svelte';
	import { createProjectQuery } from '$lib/features/projects/queries';

	const projectId = $derived(page.params.id!);
	const projectQuery = $derived(createProjectQuery(projectId));

	let projectName = $state('');
	let projectDescription = $state('');
	let showDeleteModal = $state(false);
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
</script>

<div class="border-b">
	<div class="max-w-7xl mx-auto px-6">
		<h1 class="text-3xl font-semibold my-[40px] w-full mx-auto">Settings</h1>
	</div>
</div>
<div class="mt-[46px]"></div>

<div class="flex flex-col gap-6 max-w-7xl mx-auto px-6">
	<Card.Root class="gap-3">
		<Card.Header>
			<Card.Title class="text-xl">Project Metadata</Card.Title>
		</Card.Header>
		<Card.Content class="space-y-4">
			<Field.Set>
				<Field.Group>
					<Field.Field>
						<Field.Label for="name">Project Name</Field.Label>
						<Input id="name" bind:value={projectName} placeholder="My Awesome Project" />
					</Field.Field>
					<Field.Field>
						<Field.Label for="description">Description</Field.Label>

						<Textarea
							id="description"
							bind:value={projectDescription}
							placeholder="A brief description of your project..."
							rows={3}
						/>
					</Field.Field>
				</Field.Group>
			</Field.Set>
		</Card.Content>
		<Card.Footer class="justify-end">
			<Button size="sm" onclick={handleUpdateGeneral}>Save Changes</Button>
		</Card.Footer>
	</Card.Root>

	<Card.Root class="gap-3">
		<Card.Header>
			<Card.Title class="text-xl">Project ID</Card.Title>
		</Card.Header>
		<Card.Content class="space-y-4">
			<p class="text-xs">Will be used in the future for api work</p>
			<Input
				disabled={true}
				id="domain"
				type="text"
				value={projectId}
				placeholder="https://mikrocloud.example.com"
			/>
		</Card.Content>
	</Card.Root>

	<Card.Root class="border-destructive/50 gap-3">
		<Card.Header>
			<Card.Title>Delete Project</Card.Title>
		</Card.Header>
		<Card.Content>
			<p class="text-xs">
				Permanently delete this project and all its resources. This action cannot be undone.
			</p>
		</Card.Content>
		<Card.Footer class="justify-end">
			<Button size="sm" variant="destructive" onclick={() => (showDeleteModal = true)}>
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

<!-- <Dialog.Root bind:open={showTransferModal}> -->
<!-- 	<Dialog.Content class="sm:max-w-[500px]"> -->
<!-- 		<Dialog.Header> -->
<!-- 			<Dialog.Title>Transfer Project</Dialog.Title> -->
<!-- 			<Dialog.Description>Transfer this project to another organization or user</Dialog.Description> -->
<!-- 		</Dialog.Header> -->
<!---->
<!-- 		<div class="space-y-4"> -->
<!-- 			<div class="space-y-2"> -->
<!-- 				<Label for="target-org">Target Organization</Label> -->
<!-- 				<Select.Root type="single"> -->
<!-- 					<Select.Trigger>Select organization</Select.Trigger> -->
<!-- 					<Select.Content> -->
<!-- 						<Select.Item value="org1">Organization 1</Select.Item> -->
<!-- 						<Select.Item value="org2">Organization 2</Select.Item> -->
<!-- 						<Select.Item value="personal">Personal Account</Select.Item> -->
<!-- 					</Select.Content> -->
<!-- 				</Select.Root> -->
<!-- 			</div> -->
<!---->
<!-- 			<Alert.Root> -->
<!-- 				<AlertTriangle class="size-4" /> -->
<!-- 				<Alert.Description> -->
<!-- 					You will lose access to this project after the transfer is complete. -->
<!-- 				</Alert.Description> -->
<!-- 			</Alert.Root> -->
<!-- 		</div> -->
<!---->
<!-- 		<Dialog.Footer> -->
<!-- 			<Button variant="outline" onclick={() => (showTransferModal = false)}>Cancel</Button> -->
<!-- 			<Button onclick={handleTransferProject}>Transfer Project</Button> -->
<!-- 		</Dialog.Footer> -->
<!-- 	</Dialog.Content> -->
<!-- </Dialog.Root> -->

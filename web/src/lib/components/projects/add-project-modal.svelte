<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Field from '$lib/components/ui/field/index.js';
	import {
		Dialog,
		DialogContent,
		DialogDescription,
		DialogFooter,
		DialogHeader,
		DialogTitle
	} from '$lib/components/ui/dialog';
	import type { CreateProjectRequest } from '$lib/features/projects/types';

	interface Props {
		open?: boolean;
		onClose?: () => void;
		onSubmit?: (data: CreateProjectRequest) => Promise<void>;
	}

	let { open = $bindable(false), onClose, onSubmit }: Props = $props();

	let name = $state('');
	let description = $state('');
	let loading = $state(false);
	let error = $state('');

	const handleSubmit = async (e: Event) => {
		e.preventDefault();
		error = '';

		if (!name.trim()) {
			error = 'Project name is required';
			return;
		}

		loading = true;
		try {
			await onSubmit?.({ name: name.trim(), description: description.trim() || undefined });
			name = '';
			description = '';
			open = false;
			onClose?.();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create project';
		} finally {
			loading = false;
		}
	};

	const handleOpenChange = (isOpen: boolean) => {
		open = isOpen;
		if (!isOpen) {
			name = '';
			description = '';
			error = '';
			onClose?.();
		}
	};
</script>

<Dialog bind:open onOpenChange={handleOpenChange}>
	<DialogContent>
		<DialogHeader>
			<DialogTitle>Create New Project</DialogTitle>
			<DialogDescription>Add a new project to your workspace.</DialogDescription>
		</DialogHeader>
		<form onsubmit={handleSubmit} class="py-4" id="project_form">
			<Field.Set>
				<Field.Group>
					<Field.Field>
						<Field.Label for="name">Project name</Field.Label>
						<Input
							id="name"
							type="text"
							bind:value={name}
							placeholder="my-awesome-project"
							disabled={loading}
							required
						/>
					</Field.Field>

					<Field.Field>
						<Field.Label for="description">Description (optional)</Field.Label>
						<Input
							id="description"
							type="text"
							bind:value={description}
							placeholder="A brief description of your project"
							disabled={loading}
						/>
					</Field.Field>

					{#if error}
						<Field.Error>{error}</Field.Error>
					{/if}
				</Field.Group>
			</Field.Set>
		</form>
		<DialogFooter>
			<!-- <Button type="button" variant="outline" onclick={() => (open = false)} disabled={loading}> -->
			<!-- 	Cancel -->
			<!-- </Button> -->
			<Button form="project_form" type="submit" disabled={loading}>
				{loading ? 'Creating...' : 'Create Project'}
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>

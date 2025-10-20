<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Button } from '$lib/components/ui/button';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { createOrganizationsListQuery } from '$lib/features/organizations/queries';
	import { updateOrganizationMutationQuery } from '$lib/features/organizations/mutations';
	import { toast } from 'svelte-sonner';

	const organizationsQuery = createOrganizationsListQuery();
	const updateMutation = updateOrganizationMutationQuery({
		onSuccess: () => {
			toast.success('Organization updated successfully');
		},
		onError: (error) => {
			toast.error(`Failed to update organization: ${error.message}`);
		}
	});

	const currentOrg = $derived(organizationsQuery.data?.[0]);

	let organizationName = $state('');
	let description = $state('');

	$effect(() => {
		if (currentOrg) {
			organizationName = currentOrg.name;
			description = currentOrg.description || '';
		}
	});

	function handleSave() {
		if (!currentOrg) return;

		updateMutation.mutate({
			id: currentOrg.id,
			data: {
				name: organizationName,
				description: description || undefined,
			}
		});
	}
</script>

{#if organizationsQuery.isLoading}
	<div class="flex items-center justify-center py-12">
		<p class="text-muted-foreground">Loading organization...</p>
	</div>
{:else if organizationsQuery.isError}
	<div class="flex items-center justify-center py-12">
		<p class="text-destructive">Failed to load organization</p>
	</div>
{:else if !currentOrg}
	<div class="flex items-center justify-center py-12">
		<p class="text-muted-foreground">No organization found</p>
	</div>
{:else}
	<div class="space-y-6">
		<div>
			<h1 class="text-3xl font-bold mb-2">General Settings</h1>
			<p class="text-muted-foreground">
				Manage your organization's basic information and settings.
			</p>
		</div>

		<Card>
			<CardHeader>
				<CardTitle>Organization Information</CardTitle>
				<CardDescription>Update your organization's name and description</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="space-y-2">
					<Label for="org-name">Organization Name</Label>
					<Input id="org-name" bind:value={organizationName} placeholder="My Organization" />
				</div>
				<div class="space-y-2">
					<Label for="org-slug">Organization Slug</Label>
					<Input id="org-slug" value={currentOrg.slug} readonly class="opacity-60" />
					<p class="text-sm text-muted-foreground">
						Used in URLs and API endpoints. Cannot be changed after creation.
					</p>
				</div>
				<div class="space-y-2">
					<Label for="description">Description</Label>
					<Textarea
						id="description"
						bind:value={description}
						placeholder="A brief description of your organization"
						rows={3}
					/>
				</div>
			</CardContent>
		</Card>

		<Card>
			<CardHeader>
				<CardTitle>Organization ID</CardTitle>
				<CardDescription>Your unique organization identifier</CardDescription>
			</CardHeader>
			<CardContent>
				<div class="flex items-center gap-2">
					<Input value={currentOrg.id} readonly class="font-mono text-sm" />
					<Button
						variant="outline"
						size="sm"
						onclick={() => {
							navigator.clipboard.writeText(currentOrg.id);
							toast.success('Organization ID copied to clipboard');
						}}
					>
						Copy
					</Button>
				</div>
			</CardContent>
		</Card>

		<Button onclick={handleSave} disabled={updateMutation.isPending}>
			{updateMutation.isPending ? 'Saving...' : 'Save Changes'}
		</Button>
	</div>
{/if}

<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Select from '$lib/components/ui/select';
	import { Badge } from '$lib/components/ui/badge';

	let transferEmail = $state('');
	let confirmDeleteText = $state('');
	let showTransferDialog = $state(false);
	let showDeleteDialog = $state(false);

	const organizationOwners = [
		{ id: 1, name: 'John Doe', email: 'john@example.com', avatarUrl: null },
		{ id: 2, name: 'Jane Smith', email: 'jane@example.com', avatarUrl: null }
	];

	let selectedOwner = $state<typeof organizationOwners[0] | null>(null);

	function handleTransferOwnership() {
		console.log('Transferring ownership to:', selectedOwner);
		showTransferDialog = false;
		selectedOwner = null;
	}

	function handleDeleteOrganization() {
		console.log('Deleting organization...');
		showDeleteDialog = false;
		confirmDeleteText = '';
	}

	const canDelete = $derived(confirmDeleteText === 'DELETE');
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-3xl font-bold mb-2">Organization Settings</h1>
		<p class="text-muted-foreground">
			Advanced organization settings and dangerous actions.
		</p>
	</div>

	<Card>
		<CardHeader>
			<CardTitle>Transfer Ownership</CardTitle>
			<CardDescription>
				Transfer ownership of this organization to another member with Admin or Owner role
			</CardDescription>
		</CardHeader>
		<CardContent class="space-y-4">
			<div class="p-4 border border-yellow-500/50 bg-yellow-500/10 rounded-lg">
				<div class="flex items-start gap-3">
					<svg
						class="h-5 w-5 text-yellow-500 mt-0.5 shrink-0"
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
						/>
					</svg>
					<div class="space-y-1">
						<p class="text-sm font-medium text-yellow-700 dark:text-yellow-400">
							Warning: This action cannot be undone
						</p>
						<p class="text-sm text-yellow-600 dark:text-yellow-500">
							Transferring ownership will immediately revoke your owner privileges. You will become
							an admin unless the new owner changes your role.
						</p>
					</div>
				</div>
			</div>

			<Button variant="outline" onclick={() => (showTransferDialog = true)}>
				Transfer Ownership
			</Button>
		</CardContent>
	</Card>

	<Card class="border-destructive">
		<CardHeader>
			<CardTitle class="text-destructive">Danger Zone</CardTitle>
			<CardDescription>Irreversible and destructive actions</CardDescription>
		</CardHeader>
		<CardContent class="space-y-4">
			<div class="flex items-center justify-between p-4 border rounded-lg">
				<div>
					<p class="font-medium">Delete Organization</p>
					<p class="text-sm text-muted-foreground">
						Permanently delete this organization and all associated data
					</p>
				</div>
				<Button variant="destructive" onclick={() => (showDeleteDialog = true)}>
					Delete Organization
				</Button>
			</div>

			<div class="p-4 border border-destructive/50 bg-destructive/10 rounded-lg">
				<div class="flex items-start gap-3">
					<svg
						class="h-5 w-5 text-destructive mt-0.5 shrink-0"
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
						/>
					</svg>
					<div class="space-y-1">
						<p class="text-sm font-medium">Deleting this organization will:</p>
						<ul class="text-sm text-muted-foreground space-y-1 list-disc list-inside">
							<li>Remove all projects and deployments</li>
							<li>Delete all databases and their backups</li>
							<li>Remove all team members</li>
							<li>Cancel any active subscriptions</li>
							<li>Delete all configuration and settings</li>
						</ul>
						<p class="text-sm font-medium mt-2">This action cannot be undone.</p>
					</div>
				</div>
			</div>
		</CardContent>
	</Card>
</div>

<Dialog.Root bind:open={showTransferDialog}>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Transfer Organization Ownership</Dialog.Title>
			<Dialog.Description>
				Select a member to transfer ownership to. This action cannot be undone.
			</Dialog.Description>
		</Dialog.Header>
		<div class="space-y-4 py-4">
			<div class="space-y-2">
				<Label for="new-owner">New Owner</Label>
				<Select.Root
					onSelectedChange={(selected) => {
						selectedOwner = selected?.value || null;
					}}
				>
					<Select.Trigger id="new-owner">
						<Select.Value placeholder="Select a member" />
					</Select.Trigger>
					<Select.Content>
						{#each organizationOwners as owner}
							<Select.Item value={owner}>
								<div class="flex items-center gap-2">
									<div class="w-6 h-6 rounded-full bg-primary/10 flex items-center justify-center">
										<span class="text-xs font-medium">{owner.name.charAt(0)}</span>
									</div>
									<div>
										<p class="text-sm font-medium">{owner.name}</p>
										<p class="text-xs text-muted-foreground">{owner.email}</p>
									</div>
								</div>
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
			</div>

			<div class="p-4 border border-yellow-500/50 bg-yellow-500/10 rounded-lg">
				<p class="text-sm text-yellow-700 dark:text-yellow-400">
					After transferring ownership, you will no longer have owner privileges. The new owner can
					change your role or remove you from the organization.
				</p>
			</div>
		</div>
		<Dialog.Footer>
			<Button variant="outline" onclick={() => (showTransferDialog = false)}>Cancel</Button>
			<Button onclick={handleTransferOwnership} disabled={!selectedOwner}>
				Transfer Ownership
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<Dialog.Root bind:open={showDeleteDialog}>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Delete Organization</Dialog.Title>
			<Dialog.Description>
				This action cannot be undone. All data associated with this organization will be
				permanently deleted.
			</Dialog.Description>
		</Dialog.Header>
		<div class="space-y-4 py-4">
			<div class="space-y-2">
				<Label for="confirm-delete">
					Type <code class="text-destructive font-mono">DELETE</code> to confirm
				</Label>
				<Input
					id="confirm-delete"
					bind:value={confirmDeleteText}
					placeholder="DELETE"
					class="font-mono"
				/>
			</div>

			<div class="p-4 border border-destructive/50 bg-destructive/10 rounded-lg">
				<p class="text-sm font-medium mb-2">This will permanently delete:</p>
				<ul class="text-sm text-muted-foreground space-y-1 list-disc list-inside">
					<li>All projects and deployments</li>
					<li>All databases and backups</li>
					<li>All team members and invitations</li>
					<li>All billing and subscription data</li>
				</ul>
			</div>
		</div>
		<Dialog.Footer>
			<Button variant="outline" onclick={() => (showDeleteDialog = false)}>Cancel</Button>
			<Button variant="destructive" onclick={handleDeleteOrganization} disabled={!canDelete}>
				Delete Organization
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

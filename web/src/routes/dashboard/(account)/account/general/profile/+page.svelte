<script lang="ts">
	import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { authApi } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Copy, Trash2 } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	const queryClient = useQueryClient();

	const userQuery = createQuery(() => ({
		queryKey: ['user-profile'],
		queryFn: () => authApi.getProfile()
	}));

	let avatarFile: File | null = $state(null);
	let avatarPreview: string | null = $state(null);
	let displayName = $state('');
	let username = $state('');
	let isDeleteModalOpen = $state(false);

	$effect(() => {
		if (userQuery.data) {
			displayName = userQuery.data.name || '';
			username = userQuery.data.username || '';
		}
	});

	function handleAvatarChange(e: Event) {
		const target = e.target as HTMLInputElement;
		const file = target.files?.[0];

		if (!file) return;

		if (file.size > 2 * 1024 * 1024) {
			toast.error('File size must be less than 2MB');
			return;
		}

		if (!file.type.startsWith('image/')) {
			toast.error('File must be an image');
			return;
		}

		avatarFile = file;

		const reader = new FileReader();
		reader.onload = (e) => {
			avatarPreview = e.target?.result as string;
		};
		reader.readAsDataURL(file);
	}

	const updateProfileMutation = createMutation(() => ({
		mutationFn: async (data: { name?: string; username?: string; avatar?: string }) => {
			return fetch('/api/auth/profile', {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${authApi.getToken()}`
				},
				body: JSON.stringify(data)
			});
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['user-profile'] });
			toast.success('Profile updated successfully');
		},
		onError: () => {
			toast.error('Failed to update profile');
		}
	}));

	async function handleSaveDisplayName() {
		if (!displayName || displayName.length > 32) {
			toast.error('Display name must be between 1 and 32 characters');
			return;
		}

		let avatarData: string | undefined;
		if (avatarPreview) {
			avatarData = avatarPreview;
		}

		updateProfileMutation.mutate({ name: displayName, avatar: avatarData });
	}

	async function handleSaveUsername() {
		if (!username) {
			toast.error('Username is required');
			return;
		}

		updateProfileMutation.mutate({ username });
	}

	function copyUserId() {
		if (userQuery.data?.id) {
			navigator.clipboard.writeText(userQuery.data.id);
			toast.success('User ID copied to clipboard');
		}
	}

	function handleDeleteAccount() {
		isDeleteModalOpen = true;
	}

	const deleteAccountMutation = createMutation(() => ({
		mutationFn: async () => {
			return fetch('/api/auth/profile', {
				method: 'DELETE',
				headers: {
					Authorization: `Bearer ${authApi.getToken()}`
				}
			});
		},
		onSuccess: () => {
			authApi.logout();
			window.location.href = '/auth/login';
		},
		onError: () => {
			toast.error('Failed to delete account');
		}
	}));

	function confirmDeleteAccount() {
		deleteAccountMutation.mutate();
	}
</script>

<div class="space-y-8">
	<div>
		<h1 class="text-3xl font-bold">Profile</h1>
		<p class="text-muted-foreground mt-2">Manage your personal information and preferences.</p>
	</div>

	<div class="space-y-6">
		<div class="border border-border rounded-lg p-6">
			<div class="flex items-start justify-between">
				<div class="space-y-1">
					<h3 class="text-lg font-semibold">Avatar</h3>
					<p class="text-sm text-muted-foreground">This is your avatar.</p>
					<p class="text-sm text-muted-foreground">
						Click on the avatar to upload a custom one from your files.
					</p>
				</div>

				<label for="avatar-upload" class="cursor-pointer">
					{#if avatarPreview}
						<img
							src={avatarPreview}
							alt="Avatar preview"
							class="w-20 h-20 rounded-full object-cover"
						/>
					{:else if userQuery.data?.name}
						<div
							class="w-20 h-20 rounded-full bg-primary flex items-center justify-center text-2xl font-bold text-primary-foreground"
						>
							{userQuery.data.name.charAt(0).toUpperCase()}
						</div>
					{/if}
					<input
						id="avatar-upload"
						type="file"
						accept="image/*"
						class="hidden"
						onchange={handleAvatarChange}
					/>
				</label>
			</div>

			<div class="mt-4 p-4 bg-muted/50 rounded-lg">
				<p class="text-sm text-muted-foreground">An avatar is optional but strongly recommended.</p>
			</div>
		</div>

		<div class="border border-border rounded-lg p-6">
			<div class="space-y-4">
				<div class="space-y-1">
					<h3 class="text-lg font-semibold">Display Name</h3>
					<p class="text-sm text-muted-foreground">
						Please enter your full name, or a display name you are comfortable with.
					</p>
				</div>

				<div class="space-y-2">
					<Input
						bind:value={displayName}
						placeholder="Enter display name"
						maxlength={32}
						disabled={userQuery.isLoading}
					/>
				</div>

				<div class="flex items-center justify-between">
					<p class="text-sm text-muted-foreground">Please use 32 characters at maximum.</p>
					<Button onclick={handleSaveDisplayName} disabled={updateProfileMutation.isPending}>
						{updateProfileMutation.isPending ? 'Saving...' : 'Save'}
					</Button>
				</div>
			</div>
		</div>

		<div class="border border-border rounded-lg p-6">
			<div class="space-y-4">
				<div class="space-y-1">
					<h3 class="text-lg font-semibold">Username</h3>
					<p class="text-sm text-muted-foreground">This is your URL namespace within MikroCloud.</p>
				</div>

				<div class="space-y-2">
					<Input
						bind:value={username}
						placeholder="Enter username"
						disabled={userQuery.isLoading}
					/>
				</div>

				<div class="flex items-center justify-end">
					<Button onclick={handleSaveUsername} disabled={updateProfileMutation.isPending}>
						{updateProfileMutation.isPending ? 'Saving...' : 'Save'}
					</Button>
				</div>
			</div>
		</div>

		<div class="border border-border rounded-lg p-6">
			<div class="space-y-4">
				<div class="space-y-1">
					<h3 class="text-lg font-semibold">User ID</h3>
					<p class="text-sm text-muted-foreground">This is your user ID within MikroCloud.</p>
				</div>

				<div class="flex items-center gap-2">
					<Input value={userQuery.data?.id || ''} readonly class="font-mono" />
					<Button variant="outline" size="icon" onclick={copyUserId}>
						<Copy class="w-4 h-4" />
					</Button>
				</div>

				<p class="text-sm text-muted-foreground">Used when interacting with the MikroCloud API.</p>
			</div>
		</div>

		<div class="border border-border rounded-lg p-6">
			<div class="space-y-4">
				<div class="space-y-1">
					<h3 class="text-lg font-semibold">Reset Tips</h3>
					<p class="text-sm text-muted-foreground">See onboarding tips you might have missed.</p>
				</div>

				<div class="flex items-center justify-end">
					<Button variant="outline">Reset</Button>
				</div>
			</div>
		</div>

		<div class="border border-destructive rounded-lg p-6 bg-destructive/5">
			<div class="space-y-4">
				<div class="space-y-1">
					<h3 class="text-lg font-semibold text-destructive">Delete Account</h3>
					<p class="text-sm text-muted-foreground">
						Permanently remove your Personal Account and all of its contents from the MikroCloud
						platform. This action is not reversible, so please continue with caution.
					</p>
				</div>

				<div class="flex items-center justify-end">
					<Button variant="destructive" onclick={handleDeleteAccount}>
						<Trash2 class="w-4 h-4 mr-2" />
						Delete Personal Account
					</Button>
				</div>
			</div>
		</div>
	</div>
</div>

{#if isDeleteModalOpen}
	<div class="fixed inset-0 z-50 bg-black/80 flex items-center justify-center">
		<div
			class="bg-card border border-border rounded-lg p-6 max-w-md w-full mx-4"
			onclick={(e) => e.stopPropagation()}
		>
			<h3 class="text-xl font-semibold mb-4">Delete Account</h3>
			<p class="text-muted-foreground mb-6">
				Are you sure you want to delete your account? This action cannot be undone. All your data
				will be permanently removed.
			</p>
			<div class="flex gap-3 justify-end">
				<Button variant="outline" onclick={() => (isDeleteModalOpen = false)}>Cancel</Button>
				<Button
					variant="destructive"
					onclick={confirmDeleteAccount}
					disabled={deleteAccountMutation.isPending}
				>
					{deleteAccountMutation.isPending ? 'Deleting...' : 'Delete Account'}
				</Button>
			</div>
		</div>
	</div>
{/if}

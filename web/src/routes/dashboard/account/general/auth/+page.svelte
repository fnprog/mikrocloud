<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import { Mail, Lock, Key, Github, Gitlab, Container, Shield, Plus, Trash2 } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	import { createQuery, createMutation } from '@tanstack/svelte-query';
	import { authApi } from '$lib/api/auth';

	let emailForm = $state({
		newEmail: '',
		currentPassword: ''
	});

	let passwordForm = $state({
		currentPassword: '',
		newPassword: '',
		confirmPassword: ''
	});

	let twoFactorEnabled = $state(false);

	const userQuery = createQuery(() => ({
		queryKey: ['user', 'profile'],
		queryFn: () => authApi.getProfile()
	}));

	let user = $derived(userQuery.data);

	const updateEmailMutation = createMutation(() => ({
		mutationFn: async (data: { email: string; password: string }) => {
			const response = await fetch('/api/auth/email', {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${authApi.getToken()}`
				},
				body: JSON.stringify(data)
			});
			if (!response.ok) {
				const error = await response.json();
				throw new Error(error.message || 'Failed to update email');
			}
			return response.json();
		},
		onSuccess: () => {
			toast.success('Email updated successfully');
			emailForm.newEmail = '';
			emailForm.currentPassword = '';
		},
		onError: (error: Error) => {
			toast.error('Failed to update email', { description: error.message });
		}
	}));

	let isEmailPending = $derived(updateEmailMutation.isPending);

	const updatePasswordMutation = createMutation(() => ({
		mutationFn: async (data: {
			currentPassword: string;
			newPassword: string;
			confirmPassword: string;
		}) => {
			const response = await fetch('/api/auth/password', {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${authApi.getToken()}`
				},
				body: JSON.stringify(data)
			});
			if (!response.ok) {
				const error = await response.json();
				throw new Error(error.message || 'Failed to update password');
			}
			return response.json();
		},
		onSuccess: () => {
			toast.success('Password updated successfully');
			passwordForm.currentPassword = '';
			passwordForm.newPassword = '';
			passwordForm.confirmPassword = '';
		},
		onError: (error: Error) => {
			toast.error('Failed to update password', { description: error.message });
		}
	}));

	let isPasswordPending = $derived(updatePasswordMutation.isPending);

	function handleEmailUpdate() {
		if (!emailForm.newEmail || !emailForm.currentPassword) {
			toast.error('Please fill in all fields');
			return;
		}

		const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		if (!emailRegex.test(emailForm.newEmail)) {
			toast.error('Please enter a valid email address');
			return;
		}

		updateEmailMutation.mutate({
			email: emailForm.newEmail,
			password: emailForm.currentPassword
		});
	}

	function handlePasswordUpdate() {
		if (
			!passwordForm.currentPassword ||
			!passwordForm.newPassword ||
			!passwordForm.confirmPassword
		) {
			toast.error('Please fill in all fields');
			return;
		}

		if (passwordForm.newPassword.length < 8) {
			toast.error('New password must be at least 8 characters');
			return;
		}

		if (passwordForm.newPassword !== passwordForm.confirmPassword) {
			toast.error('New passwords do not match');
			return;
		}

		updatePasswordMutation.mutate(passwordForm);
	}

	function handleToggle2FA() {
		if (!twoFactorEnabled) {
			toast.info('2FA setup not implemented yet');
		} else {
			toast.info('2FA disable not implemented yet');
		}
	}

	function handleAddPasskey() {
		toast.info('Passkey registration not implemented yet');
	}

	function handleConnectOAuth(provider: string) {
		toast.info(`${provider} OAuth connection not implemented yet`);
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold">Authentication</h1>
		<p class="text-sm text-muted-foreground mt-1">
			Manage your email, password, and sign-in methods
		</p>
	</div>

	<div class="border border-border rounded-lg p-6 space-y-6">
		<div>
			<h2 class="text-lg font-semibold mb-4 flex items-center gap-2">
				<Mail class="w-5 h-5" />
				Change Email
			</h2>
			<div class="space-y-4 max-w-md">
				<div class="space-y-2">
					<Label for="current-email">Current Email</Label>
					<Input
						id="current-email"
						type="email"
						value={user?.email || ''}
						disabled
						class="bg-muted"
					/>
				</div>
				<div class="space-y-2">
					<Label for="new-email">New Email</Label>
					<Input
						id="new-email"
						type="email"
						bind:value={emailForm.newEmail}
						placeholder="new.email@example.com"
					/>
				</div>
				<div class="space-y-2">
					<Label for="email-password">Current Password</Label>
					<Input
						id="email-password"
						type="password"
						bind:value={emailForm.currentPassword}
						placeholder="Enter your current password"
					/>
				</div>
				<Button onclick={handleEmailUpdate} disabled={isEmailPending}>
					{isEmailPending ? 'Updating...' : 'Update Email'}
				</Button>
			</div>
		</div>

		<div class="border-t border-border pt-6">
			<h2 class="text-lg font-semibold mb-4 flex items-center gap-2">
				<Lock class="w-5 h-5" />
				Change Password
			</h2>
			<div class="space-y-4 max-w-md">
				<div class="space-y-2">
					<Label for="current-password">Current Password</Label>
					<Input
						id="current-password"
						type="password"
						bind:value={passwordForm.currentPassword}
						placeholder="Enter your current password"
					/>
				</div>
				<div class="space-y-2">
					<Label for="new-password">New Password</Label>
					<Input
						id="new-password"
						type="password"
						bind:value={passwordForm.newPassword}
						placeholder="Enter new password (min. 8 characters)"
					/>
				</div>
				<div class="space-y-2">
					<Label for="confirm-password">Confirm New Password</Label>
					<Input
						id="confirm-password"
						type="password"
						bind:value={passwordForm.confirmPassword}
						placeholder="Confirm new password"
					/>
				</div>
				<Button onclick={handlePasswordUpdate} disabled={isPasswordPending}>
					{isPasswordPending ? 'Updating...' : 'Update Password'}
				</Button>
			</div>
		</div>
	</div>

	<!-- <div class="border border-border rounded-lg p-6 space-y-6"> -->
	<!-- 	<div> -->
	<!-- 		<h2 class="text-lg font-semibold mb-4 flex items-center gap-2"> -->
	<!-- 			<Key class="w-5 h-5" /> -->
	<!-- 			Sign-in Methods -->
	<!-- 		</h2> -->
	<!---->
	<!-- 		<div class="space-y-4"> -->
	<!-- 			<div class="flex items-center justify-between p-4 border border-border rounded-lg"> -->
	<!-- 				<div class="flex items-center gap-3"> -->
	<!-- 					<Mail class="w-5 h-5 text-muted-foreground" /> -->
	<!-- 					<div> -->
	<!-- 						<p class="font-medium">Email</p> -->
	<!-- 						<p class="text-sm text-muted-foreground">{user?.email || 'Loading...'}</p> -->
	<!-- 					</div> -->
	<!-- 				</div> -->
	<!-- 				<Button variant="outline" size="sm">Manage</Button> -->
	<!-- 			</div> -->
	<!---->
	<!-- 			<div class="flex items-center justify-between p-4 border border-border rounded-lg"> -->
	<!-- 				<div class="flex items-center gap-3"> -->
	<!-- 					<Key class="w-5 h-5 text-muted-foreground" /> -->
	<!-- 					<div> -->
	<!-- 						<p class="font-medium">Passkeys</p> -->
	<!-- 						<p class="text-sm text-muted-foreground">0 passkeys registered</p> -->
	<!-- 					</div> -->
	<!-- 				</div> -->
	<!-- 				<Button variant="outline" size="sm" onclick={handleAddPasskey}> -->
	<!-- 					<Plus class="w-4 h-4 mr-1" /> -->
	<!-- 					Add -->
	<!-- 				</Button> -->
	<!-- 			</div> -->
	<!---->
	<!-- 			<div class="border-t border-border pt-4"> -->
	<!-- 				<p class="text-sm font-medium mb-3">OAuth Providers</p> -->
	<!-- 				<div class="space-y-2"> -->
	<!-- 					<div class="flex items-center justify-between p-3 border border-border rounded-lg"> -->
	<!-- 						<div class="flex items-center gap-3"> -->
	<!-- 							<Github class="w-5 h-5" /> -->
	<!-- 							<span class="text-sm">GitHub</span> -->
	<!-- 						</div> -->
	<!-- 						<Button -->
	<!-- 							variant="outline" -->
	<!-- 							size="sm" -->
	<!-- 							disabled -->
	<!-- 							onclick={() => handleConnectOAuth('GitHub')} -->
	<!-- 						> -->
	<!-- 							Connect -->
	<!-- 						</Button> -->
	<!-- 					</div> -->
	<!---->
	<!-- 					<div class="flex items-center justify-between p-3 border border-border rounded-lg"> -->
	<!-- 						<div class="flex items-center gap-3"> -->
	<!-- 							<Gitlab class="w-5 h-5" /> -->
	<!-- 							<span class="text-sm">GitLab</span> -->
	<!-- 						</div> -->
	<!-- 						<Button -->
	<!-- 							variant="outline" -->
	<!-- 							size="sm" -->
	<!-- 							disabled -->
	<!-- 							onclick={() => handleConnectOAuth('GitLab')} -->
	<!-- 						> -->
	<!-- 							Connect -->
	<!-- 						</Button> -->
	<!-- 					</div> -->
	<!---->
	<!-- 					<div class="flex items-center justify-between p-3 border border-border rounded-lg"> -->
	<!-- 						<div class="flex items-center gap-3"> -->
	<!-- 							<Container class="w-5 h-5" /> -->
	<!-- 							<span class="text-sm">Bitbucket</span> -->
	<!-- 						</div> -->
	<!-- 						<Button -->
	<!-- 							variant="outline" -->
	<!-- 							size="sm" -->
	<!-- 							disabled -->
	<!-- 							onclick={() => handleConnectOAuth('Bitbucket')} -->
	<!-- 						> -->
	<!-- 							Connect -->
	<!-- 						</Button> -->
	<!-- 					</div> -->
	<!-- 				</div> -->
	<!-- 			</div> -->
	<!-- 		</div> -->
	<!-- 	</div> -->
	<!-- </div> -->

	<!-- <div class="border border-border rounded-lg p-6"> -->
	<!-- 	<div> -->
	<!-- 		<h2 class="text-lg font-semibold mb-4 flex items-center gap-2"> -->
	<!-- 			<Shield class="w-5 h-5" /> -->
	<!-- 			Two-Factor Authentication -->
	<!-- 		</h2> -->
	<!-- 		<div class="space-y-4"> -->
	<!-- 			<div class="flex items-center justify-between p-4 border border-border rounded-lg"> -->
	<!-- 				<div> -->
	<!-- 					<p class="font-medium">Enable 2FA</p> -->
	<!-- 					<p class="text-sm text-muted-foreground"> -->
	<!-- 						Add an extra layer of security to your account -->
	<!-- 					</p> -->
	<!-- 				</div> -->
	<!-- 				<Switch checked={twoFactorEnabled} onCheckedChange={handleToggle2FA} /> -->
	<!-- 			</div> -->
	<!---->
	<!-- 			{#if twoFactorEnabled} -->
	<!-- 				<div class="space-y-3 pl-4 border-l-2 border-primary"> -->
	<!-- 					<div class="flex items-center justify-between p-3 bg-muted rounded-lg"> -->
	<!-- 						<div class="flex items-center gap-3"> -->
	<!-- 							<Key class="w-4 h-4 text-muted-foreground" /> -->
	<!-- 							<div> -->
	<!-- 								<p class="text-sm font-medium">Authenticator App</p> -->
	<!-- 								<p class="text-xs text-muted-foreground">Use TOTP codes from an app</p> -->
	<!-- 							</div> -->
	<!-- 						</div> -->
	<!-- 						<Button variant="outline" size="sm">Setup</Button> -->
	<!-- 					</div> -->
	<!---->
	<!-- 					<div class="flex items-center justify-between p-3 bg-muted rounded-lg"> -->
	<!-- 						<div class="flex items-center gap-3"> -->
	<!-- 							<Key class="w-4 h-4 text-muted-foreground" /> -->
	<!-- 							<div> -->
	<!-- 								<p class="text-sm font-medium">Security Keys</p> -->
	<!-- 								<p class="text-xs text-muted-foreground">Use passkeys for 2FA</p> -->
	<!-- 							</div> -->
	<!-- 						</div> -->
	<!-- 						<Button variant="outline" size="sm"> -->
	<!-- 							<Plus class="w-3 h-3 mr-1" /> -->
	<!-- 							Add -->
	<!-- 						</Button> -->
	<!-- 					</div> -->
	<!-- 				</div> -->
	<!-- 			{/if} -->
	<!-- 		</div> -->
	<!-- 	</div> -->
	<!-- </div> -->
</div>

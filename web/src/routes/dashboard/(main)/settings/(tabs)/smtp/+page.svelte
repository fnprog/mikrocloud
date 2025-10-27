<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import * as Card from '$lib/components/ui/card';
	import { createSMTPSettingsQuery } from '$lib/features/settings/queries/settings';
	import { createUpdateSMTPSettingsMutation } from '$lib/features/settings/mutations';

	const settingsQuery = createSMTPSettingsQuery();
	const updateSettingsMutation = createUpdateSMTPSettingsMutation();

	let enabled = $state(false);
	let host = $state('');
	let port = $state(587);
	let username = $state('');
	let password = $state('');
	let fromEmail = $state('');
	let fromName = $state('Mikrocloud');

	$effect(() => {
		if (settingsQuery.data) {
			enabled = settingsQuery.data.enabled ?? false;
			host = settingsQuery.data.host || '';
			port = settingsQuery.data.port || 587;
			username = settingsQuery.data.username || '';
			password = settingsQuery.data.password || '';
			fromEmail = settingsQuery.data.from_email || '';
			fromName = settingsQuery.data.from_name || 'Mikrocloud';
		}
	});

	async function handleSave() {
		await updateSettingsMutation.mutateAsync({
			enabled,
			host,
			port,
			username,
			password,
			from_email: fromEmail,
			from_name: fromName
		});
	}
</script>

<div class="space-y-6">
	<Card.Root>
		<Card.Header>
			<Card.Title>SMTP Configuration</Card.Title>
			<Card.Description>
				Configure SMTP settings for sending emails (password reset, notifications, etc.)
			</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-6">
			<div class="flex items-center justify-between">
				<div class="space-y-1">
					<Label for="smtp-enabled">Enable SMTP</Label>
					<p class="text-xs text-muted-foreground">
						Enable SMTP for sending emails from the application.
					</p>
				</div>
				<Switch id="smtp-enabled" bind:checked={enabled} />
			</div>

			{#if enabled}
				<div class="grid grid-cols-2 gap-4">
					<div class="space-y-2">
						<Label for="smtp-host">SMTP Host</Label>
						<Input id="smtp-host" bind:value={host} placeholder="smtp.gmail.com" />
					</div>

					<div class="space-y-2">
						<Label for="smtp-port">SMTP Port</Label>
						<Input id="smtp-port" type="number" bind:value={port} placeholder="587" />
					</div>
				</div>

				<div class="grid grid-cols-2 gap-4">
					<div class="space-y-2">
						<Label for="smtp-username">Username</Label>
						<Input id="smtp-username" bind:value={username} placeholder="your-email@gmail.com" />
					</div>

					<div class="space-y-2">
						<Label for="smtp-password">Password</Label>
						<Input id="smtp-password" type="password" bind:value={password} placeholder="your-app-password" />
					</div>
				</div>

				<div class="grid grid-cols-2 gap-4">
					<div class="space-y-2">
						<Label for="smtp-from-email">From Email</Label>
						<Input id="smtp-from-email" type="email" bind:value={fromEmail} placeholder="noreply@mikrocloud.com" />
					</div>

					<div class="space-y-2">
						<Label for="smtp-from-name">From Name</Label>
						<Input id="smtp-from-name" bind:value={fromName} placeholder="Mikrocloud" />
					</div>
				</div>
			{/if}
		</Card.Content>
	</Card.Root>

	<div class="flex justify-end">
		<Button onclick={handleSave} disabled={updateSettingsMutation.isPending}>
			{#if updateSettingsMutation.isPending}
				Saving...
			{:else}
				Save Changes
			{/if}
		</Button>
	</div>
</div>
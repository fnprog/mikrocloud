<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import * as Card from '$lib/components/ui/card';

	let updateCheckCron = $state('0 * * * *');
	let autoUpdate = $state(true);
	let autoUpdateCron = $state('0 0 * * *');

	onMount(async () => {
		try {
			const response = await fetch('/api/settings/updates');
			if (response.ok) {
				const data = await response.json();
				updateCheckCron = data.update_check_cron || '0 * * * *';
				autoUpdate = data.auto_update ?? true;
				autoUpdateCron = data.auto_update_cron || '0 0 * * *';
			}
		} catch (error) {
			console.error('Failed to load settings:', error);
		}
	});

	async function handleSave() {
		try {
			const response = await fetch('/api/settings/updates', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					update_check_cron: updateCheckCron,
					auto_update: autoUpdate,
					auto_update_cron: autoUpdateCron
				})
			});

			if (!response.ok) {
				throw new Error('Failed to save settings');
			}
		} catch (error) {
			console.error('Failed to save settings:', error);
		}
	}
</script>

<div class="space-y-6">
	<Card.Root>
		<Card.Header>
			<Card.Title>Version Checks</Card.Title>
		</Card.Header>
		<Card.Content class="space-y-6">
			<div class="space-y-2">
				<Label for="update-check">Update Check Cron Expression</Label>
				<Input id="update-check" bind:value={updateCheckCron} placeholder="0 * * * *" />
				<p class="text-xs text-muted-foreground">
					Cron expression to check for new versions and pull new Service Templates from CDN. Default
					is <code class="rounded bg-muted px-1 py-0.5">0 * * * *</code> (every hour).
				</p>
			</div>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title>Automatic Updates</Card.Title>
		</Card.Header>
		<Card.Content class="space-y-6">
			<div class="flex items-center justify-between">
				<div class="space-y-1">
					<Label for="auto-update">Enable Automatic Updates</Label>
					<p class="text-xs text-muted-foreground">
						Automatically install updates when they become available (coming soon).
					</p>
				</div>
				<Switch id="auto-update" bind:checked={autoUpdate} />
			</div>

			{#if autoUpdate}
				<div class="space-y-2">
					<Label for="auto-update-cron">Auto Update Cron Expression</Label>
					<Input id="auto-update-cron" bind:value={autoUpdateCron} placeholder="0 0 * * *" />
					<p class="text-xs text-muted-foreground">
						Cron expression for automatic updates. Default is <code
							class="rounded bg-muted px-1 py-0.5">0 0 * * *</code
						> (daily at midnight).
					</p>
				</div>
			{/if}
		</Card.Content>
	</Card.Root>

	<div class="flex justify-end">
		<Button onclick={handleSave}>Save Changes</Button>
	</div>
</div>

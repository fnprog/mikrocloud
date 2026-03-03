<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import * as Field from '$lib/components/ui/field/index.js';
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
			<Field.Set>
				<Field.Field>
					<Field.Label for="update-check">Update Check Cron Expression</Field.Label>
					<div class="flex gap-2">
						<Input id="update-check" bind:value={updateCheckCron} placeholder="0 * * * *" />
					</div>
					<Field.Description>
						Cron expression to check for new versions and pull new Service Templates from CDN.
						Default is <code class="rounded bg-muted px-1 py-0.5">0 * * * *</code> (every hour).
					</Field.Description>
				</Field.Field>
			</Field.Set>
		</Card.Content>
		<Card.Footer
			class="flex flex-col-reverse gap-2 sm:mt-4 sm:-mb-6 sm:flex-row sm:justify-end sm:rounded-b-xl sm:border-t sm:bg-muted/50 sm:px-6 sm:py-2"
		>
			<Button size="sm">save</Button>
		</Card.Footer>
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
		<Card.Footer
			class="flex flex-col-reverse gap-2 sm:mt-4 sm:-mb-6 sm:flex-row sm:justify-end sm:rounded-b-xl sm:border-t sm:bg-muted/50 sm:px-6 sm:py-2"
		>
			<Button size="sm">save</Button>
		</Card.Footer>
	</Card.Root>
</div>

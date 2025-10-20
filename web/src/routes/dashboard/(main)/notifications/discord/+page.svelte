<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import { Switch } from '$lib/components/ui/switch';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Separator } from '$lib/components/ui/separator';

	let webhookUrl = $state('');
	let botUsername = $state('MikroCloud');

	let notifyDeployments = $state(true);
	let notifyContainerEvents = $state(true);
	let notifyBackups = $state(false);
	let notifyScheduledTasks = $state(false);
	let notifyServerEvents = $state(true);

	function handleSave() {
		console.log('Saving Discord notification settings...');
	}

	function handleTest() {
		console.log('Testing Discord notification...');
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-3xl font-bold mb-2">Discord Notifications</h1>
		<p class="text-muted-foreground">
			Configure Discord webhook notifications for your account.
		</p>
	</div>

	<Card>
		<CardHeader>
			<CardTitle>Discord Webhook</CardTitle>
			<CardDescription>Set up your Discord webhook URL to receive notifications</CardDescription>
		</CardHeader>
		<CardContent class="space-y-4">
			<div class="space-y-2">
				<Label for="webhook-url">Webhook URL</Label>
				<Input
					id="webhook-url"
					bind:value={webhookUrl}
					placeholder="https://discord.com/api/webhooks/..."
					type="url"
				/>
				<p class="text-sm text-muted-foreground">
					Create a webhook in your Discord server settings → Integrations → Webhooks
				</p>
			</div>
			<div class="space-y-2">
				<Label for="bot-username">Bot Username (Optional)</Label>
				<Input id="bot-username" bind:value={botUsername} placeholder="MikroCloud" />
			</div>
		</CardContent>
	</Card>

	<Card>
		<CardHeader>
			<CardTitle>Notification Events</CardTitle>
			<CardDescription>Choose which events trigger Discord notifications</CardDescription>
		</CardHeader>
		<CardContent class="space-y-4">
			<div class="flex items-center justify-between">
				<div class="space-y-0.5">
					<Label for="notify-deployments">Deployment Events</Label>
					<p class="text-sm text-muted-foreground">Notify on deployment success, failure, or rollback</p>
				</div>
				<Switch id="notify-deployments" bind:checked={notifyDeployments} />
			</div>
			<Separator />
			<div class="flex items-center justify-between">
				<div class="space-y-0.5">
					<Label for="notify-container">Container Events</Label>
					<p class="text-sm text-muted-foreground">Notify when containers start, stop, or restart</p>
				</div>
				<Switch id="notify-container" bind:checked={notifyContainerEvents} />
			</div>
			<Separator />
			<div class="flex items-center justify-between">
				<div class="space-y-0.5">
					<Label for="notify-backups">Backup Events</Label>
					<p class="text-sm text-muted-foreground">Notify on backup completion or failure</p>
				</div>
				<Switch id="notify-backups" bind:checked={notifyBackups} />
			</div>
			<Separator />
			<div class="flex items-center justify-between">
				<div class="space-y-0.5">
					<Label for="notify-scheduled">Scheduled Task Events</Label>
					<p class="text-sm text-muted-foreground">Notify when scheduled tasks complete or fail</p>
				</div>
				<Switch id="notify-scheduled" bind:checked={notifyScheduledTasks} />
			</div>
			<Separator />
			<div class="flex items-center justify-between">
				<div class="space-y-0.5">
					<Label for="notify-server">Server Events</Label>
					<p class="text-sm text-muted-foreground">Notify on server health, resource alerts, or downtime</p>
				</div>
				<Switch id="notify-server" bind:checked={notifyServerEvents} />
			</div>
		</CardContent>
	</Card>

	<div class="flex gap-3">
		<Button onclick={handleSave}>Save Configuration</Button>
		<Button variant="outline" onclick={handleTest}>Send Test Message</Button>
	</div>
</div>

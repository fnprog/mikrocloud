<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import { Switch } from '$lib/components/ui/switch';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Separator } from '$lib/components/ui/separator';

	let botToken = $state('');
	let chatId = $state('');

	let notifyDeployments = $state(true);
	let notifyContainerEvents = $state(true);
	let notifyBackups = $state(false);
	let notifyScheduledTasks = $state(false);
	let notifyServerEvents = $state(true);

	function handleSave() {
		console.log('Saving Telegram notification settings...');
	}

	function handleTest() {
		console.log('Testing Telegram notification...');
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-3xl font-bold mb-2">Telegram Notifications</h1>
		<p class="text-muted-foreground">
			Configure Telegram bot notifications for your account.
		</p>
	</div>

	<Card>
		<CardHeader>
			<CardTitle>Telegram Bot Configuration</CardTitle>
			<CardDescription>Set up your Telegram bot to receive notifications</CardDescription>
		</CardHeader>
		<CardContent class="space-y-4">
			<div class="space-y-2">
				<Label for="bot-token">Bot Token</Label>
				<Input
					id="bot-token"
					bind:value={botToken}
					placeholder="123456789:ABCdefGHIjklMNOpqrsTUVwxyz"
					type="password"
				/>
				<p class="text-sm text-muted-foreground">
					Create a bot via <a href="https://t.me/BotFather" target="_blank" class="text-primary hover:underline">@BotFather</a> on Telegram
				</p>
			</div>
			<div class="space-y-2">
				<Label for="chat-id">Chat ID</Label>
				<Input id="chat-id" bind:value={chatId} placeholder="123456789" />
				<p class="text-sm text-muted-foreground">
					Your Telegram user or group chat ID (use <a href="https://t.me/userinfobot" target="_blank" class="text-primary hover:underline">@userinfobot</a> to find it)
				</p>
			</div>
		</CardContent>
	</Card>

	<Card>
		<CardHeader>
			<CardTitle>Notification Events</CardTitle>
			<CardDescription>Choose which events trigger Telegram notifications</CardDescription>
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

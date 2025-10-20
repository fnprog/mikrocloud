<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import { Switch } from '$lib/components/ui/switch';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Separator } from '$lib/components/ui/separator';
	import { Select, SelectContent, SelectItem, SelectTrigger } from '$lib/components/ui/select';

	let emailProvider = $state('smtp');
	let smtpHost = $state('');
	let smtpPort = $state('587');
	let smtpUsername = $state('');
	let smtpPassword = $state('');
	let smtpFromEmail = $state('');
	let smtpFromName = $state('');

	let resendApiKey = $state('');
	let resendFromEmail = $state('');

	let sendgridApiKey = $state('');
	let sendgridFromEmail = $state('');

	let notifyDeployments = $state(true);
	let notifyContainerEvents = $state(true);
	let notifyBackups = $state(false);
	let notifyScheduledTasks = $state(false);
	let notifyServerEvents = $state(true);

	function handleSave() {
		console.log('Saving email notification settings...');
	}

	function handleTest() {
		console.log('Testing email notification...');
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-3xl font-bold mb-2">Email Notifications</h1>
		<p class="text-muted-foreground">
			Configure email notifications for your account using SMTP, Resend, or SendGrid.
		</p>
	</div>

	<Card>
		<CardHeader>
			<CardTitle>Email Provider</CardTitle>
			<CardDescription>Choose your email delivery provider</CardDescription>
		</CardHeader>
		<CardContent class="space-y-4">
			<div class="space-y-2">
				<Label for="provider">Provider</Label>
				<Select type="single" bind:value={emailProvider}>
					<SelectTrigger id="provider">
						{emailProvider === 'smtp' ? 'SMTP' : emailProvider === 'resend' ? 'Resend' : 'SendGrid'}
					</SelectTrigger>
					<SelectContent>
						<SelectItem value="smtp">SMTP</SelectItem>
						<SelectItem value="resend">Resend</SelectItem>
						<SelectItem value="sendgrid">SendGrid</SelectItem>
					</SelectContent>
				</Select>
			</div>

			{#if emailProvider === 'smtp'}
				<div class="space-y-4">
					<div class="grid grid-cols-2 gap-4">
						<div class="space-y-2">
							<Label for="smtp-host">SMTP Host</Label>
							<Input id="smtp-host" bind:value={smtpHost} placeholder="smtp.example.com" />
						</div>
						<div class="space-y-2">
							<Label for="smtp-port">Port</Label>
							<Input id="smtp-port" bind:value={smtpPort} placeholder="587" />
						</div>
					</div>
					<div class="space-y-2">
						<Label for="smtp-username">Username</Label>
						<Input id="smtp-username" bind:value={smtpUsername} placeholder="user@example.com" />
					</div>
					<div class="space-y-2">
						<Label for="smtp-password">Password</Label>
						<Input id="smtp-password" type="password" bind:value={smtpPassword} />
					</div>
					<div class="space-y-2">
						<Label for="smtp-from-email">From Email</Label>
						<Input id="smtp-from-email" bind:value={smtpFromEmail} placeholder="noreply@example.com" />
					</div>
					<div class="space-y-2">
						<Label for="smtp-from-name">From Name</Label>
						<Input id="smtp-from-name" bind:value={smtpFromName} placeholder="MikroCloud" />
					</div>
				</div>
			{:else if emailProvider === 'resend'}
				<div class="space-y-4">
					<div class="space-y-2">
						<Label for="resend-api-key">Resend API Key</Label>
						<Input id="resend-api-key" type="password" bind:value={resendApiKey} placeholder="re_..." />
					</div>
					<div class="space-y-2">
						<Label for="resend-from-email">From Email</Label>
						<Input id="resend-from-email" bind:value={resendFromEmail} placeholder="noreply@yourdomain.com" />
					</div>
				</div>
			{:else if emailProvider === 'sendgrid'}
				<div class="space-y-4">
					<div class="space-y-2">
						<Label for="sendgrid-api-key">SendGrid API Key</Label>
						<Input id="sendgrid-api-key" type="password" bind:value={sendgridApiKey} placeholder="SG..." />
					</div>
					<div class="space-y-2">
						<Label for="sendgrid-from-email">From Email</Label>
						<Input id="sendgrid-from-email" bind:value={sendgridFromEmail} placeholder="noreply@yourdomain.com" />
					</div>
				</div>
			{/if}
		</CardContent>
	</Card>

	<Card>
		<CardHeader>
			<CardTitle>Notification Events</CardTitle>
			<CardDescription>Choose which events trigger email notifications</CardDescription>
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
		<Button variant="outline" onclick={handleTest}>Send Test Email</Button>
	</div>
</div>

<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import * as Card from '$lib/components/ui/card';
	import { createGeneralSettingsQuery } from '$lib/features/settings/queries/settings';
	import { createUpdateGeneralSettingsMutation } from '$lib/features/settings/mutations';
	import { settingsApi } from '$lib/features/settings/api';

	const settingsQuery = createGeneralSettingsQuery();
	const updateSettingsMutation = createUpdateGeneralSettingsMutation();

	let domain = $state('');
	let timezone = $state('UTC');
	let ipv4 = $state('');
	let ipv6 = $state('');
	let allowRegistrations = $state(true);
	let doNotTrack = $state(false);
	let isDetecting = $state(false);

	const timezones = [
		'UTC',
		'America/New_York',
		'America/Chicago',
		'America/Denver',
		'America/Los_Angeles',
		'Europe/London',
		'Europe/Paris',
		'Europe/Berlin',
		'Asia/Tokyo',
		'Asia/Shanghai',
		'Australia/Sydney'
	];

	$effect(() => {
		if (settingsQuery.data) {
			domain = settingsQuery.data.domain || '';
			timezone = settingsQuery.data.timezone || 'UTC';
			ipv4 = settingsQuery.data.ipv4 || '';
			ipv6 = settingsQuery.data.ipv6 || '';
			allowRegistrations = settingsQuery.data.allow_registrations ?? true;
			doNotTrack = settingsQuery.data.do_not_track ?? false;
		}
	});

	function handleSave() {
		updateSettingsMutation.mutate({
			domain,
			timezone,
			ipv4,
			ipv6,
			allow_registrations: allowRegistrations,
			do_not_track: doNotTrack
		});
	}

	async function handleDetectIPs() {
		isDetecting = true;
		try {
			const detected = await settingsApi.detectIPs();
			if (detected.ipv4) {
				ipv4 = detected.ipv4;
			}
			if (detected.ipv6) {
				ipv6 = detected.ipv6;
			}
		} catch (error) {
			console.error('Failed to detect IP addresses:', error);
		} finally {
			isDetecting = false;
		}
	}
</script>

<div class="space-y-6">
	<Card.Root>
		<Card.Header>
			<Card.Title>Instance Configuration</Card.Title>
		</Card.Header>
		<Card.Content class="space-y-6">
			<div class="space-y-2">
				<Label for="domain">Domain</Label>
				<Input
					id="domain"
					type="text"
					bind:value={domain}
					placeholder="https://mikrocloud.example.com"
				/>
				<p class="text-xs text-muted-foreground">
					Enter the full domain name (FQDN) of the instance, including 'https://' if you want to
					secure the dashboard with HTTPS. Setting this will make the dashboard accessible via this
					domain, secured by HTTPS, instead of just the IP address.
				</p>
			</div>

			<div class="space-y-2">
				<Label for="timezone">Timezone</Label>
				<select
					id="timezone"
					bind:value={timezone}
					class="border-input-new bg-secondary-new selection:bg-primary dark:bg-input/30 selection:text-primary-foreground ring-offset-background shadow-xs flex h-9 w-full min-w-0 rounded-md border px-3 py-1 text-base outline-none transition-[color,box-shadow] disabled:cursor-not-allowed disabled:opacity-50 md:text-sm focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]"
				>
					{#each timezones as tz}
						<option value={tz}>{tz}</option>
					{/each}
				</select>
				<p class="text-xs text-muted-foreground">
					This is used for the update check and automatic update frequency.
				</p>
			</div>

			<div class="space-y-2">
				<Label for="ipv4">Public IPv4 Address</Label>
				<div class="flex gap-2">
					<Input id="ipv4" type="text" bind:value={ipv4} placeholder="192.168.1.100" class="flex-1" />
					<Button variant="outline" onclick={handleDetectIPs} disabled={isDetecting} size="sm">
						{isDetecting ? 'Detecting...' : 'Auto-detect'}
					</Button>
				</div>
				<p class="text-xs text-muted-foreground">
					Enter the IPv4 address of the instance. It is useful if you have several IPv4 addresses.
				</p>
			</div>

			<div class="space-y-2">
				<Label for="ipv6">Public IPv6 Address</Label>
				<div class="flex gap-2">
					<Input id="ipv6" type="text" bind:value={ipv6} placeholder="2001:db8::1" class="flex-1" />
					<Button variant="outline" onclick={handleDetectIPs} disabled={isDetecting} size="sm">
						{isDetecting ? 'Detecting...' : 'Auto-detect'}
					</Button>
				</div>
				<p class="text-xs text-muted-foreground">
					Enter the IPv6 address of the instance. It is useful if you have several IPv6 addresses.
				</p>
			</div>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title>Access Control</Card.Title>
		</Card.Header>
		<Card.Content class="space-y-6">
			<div class="flex items-center justify-between">
				<div class="space-y-1">
					<Label for="allow-registrations">Allow Registrations</Label>
					<p class="text-xs text-muted-foreground">
						If disabled, the signup button will be removed and the registration route will be
						blocked.
					</p>
				</div>
				<Switch id="allow-registrations" bind:checked={allowRegistrations} />
			</div>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title>Privacy</Card.Title>
		</Card.Header>
		<Card.Content class="space-y-6">
			<div class="flex items-center justify-between">
				<div class="space-y-1">
					<Label for="do-not-track">Do Not Track</Label>
					<p class="text-xs text-muted-foreground">Disable telemetry and analytics collection.</p>
				</div>
				<Switch id="do-not-track" bind:checked={doNotTrack} />
			</div>
		</Card.Content>
	</Card.Root>

	<div class="flex justify-end">
		<Button onclick={handleSave} disabled={updateSettingsMutation.isPending}>
			{updateSettingsMutation.isPending ? 'Saving...' : 'Save Changes'}
		</Button>
	</div>
</div>

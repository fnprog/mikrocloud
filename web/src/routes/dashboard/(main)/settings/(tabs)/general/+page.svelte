<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import * as Card from '$lib/components/ui/card';
	import * as Command from '$lib/components/ui/command';
	import * as Popover from '$lib/components/ui/popover';
	import { Check, ChevronsUpDown } from 'lucide-svelte';
	import { cn } from '$lib/utils';
	import { tick } from 'svelte';
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
	let timezoneOpen = $state(false);
	let timezoneSearch = $state('');
	let triggerRef = $state<HTMLButtonElement>(null!);

	const allTimezones = $derived.by(() => {
		try {
			return Intl.supportedValuesOf('timeZone');
		} catch {
			return ['UTC'];
		}
	});

	const filteredTimezones = $derived(
		timezoneSearch
			? allTimezones.filter((tz) => tz.toLowerCase().includes(timezoneSearch.toLowerCase()))
			: allTimezones
	);

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

	function closeAndFocusTrigger() {
		timezoneOpen = false;
		tick().then(() => {
			triggerRef?.focus();
		});
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
				<Popover.Root bind:open={timezoneOpen}>
					<Popover.Trigger bind:ref={triggerRef}>
						{#snippet child({ props })}
							<Button
								variant="outline"
								class="w-full justify-between"
								{...props}
								role="combobox"
								aria-expanded={timezoneOpen}
							>
								{timezone || 'Select timezone...'}
								<ChevronsUpDown class="ml-2 size-4 shrink-0 opacity-50" />
							</Button>
						{/snippet}
					</Popover.Trigger>
					<Popover.Content class="w-[400px] p-0">
						<Command.Root>
							<Command.Input placeholder="Search timezone..." bind:value={timezoneSearch} />
							<Command.List>
								<Command.Empty>No timezone found.</Command.Empty>
								<Command.Group>
									{#each filteredTimezones as tz}
										<Command.Item
											value={tz}
											onSelect={() => {
												timezone = tz;
												closeAndFocusTrigger();
											}}
										>
											<Check class={cn('mr-2 size-4', timezone !== tz && 'text-transparent')} />
											{tz}
										</Command.Item>
									{/each}
								</Command.Group>
							</Command.List>
						</Command.Root>
					</Popover.Content>
				</Popover.Root>
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

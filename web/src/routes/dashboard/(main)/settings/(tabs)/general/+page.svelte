<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import * as Card from '$lib/components/ui/card';
	import * as Command from '$lib/components/ui/command';
	import * as Popover from '$lib/components/ui/popover';
	import * as Field from '$lib/components/ui/field/index.js';
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
		}
	});

	function handleSave() {
		updateSettingsMutation.mutate({
			domain,
			timezone,
			ipv4,
			ipv6,
			allow_registrations: allowRegistrations,
			do_not_track: true
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
	<Card.Root class="gap-3">
		<Card.Header>
			<Card.Title class="text-xl">Domain</Card.Title>
		</Card.Header>
		<Card.Content class="space-y-4">
			<p class="text-xs">
				Enter the full domain name (FQDN) of the instance, including 'https://' if you want to
				secure the dashboard with HTTPS. Setting this will make the dashboard accessible via this
				domain, secured by HTTPS, instead of just the IP address.
			</p>
			<Input
				id="domain"
				type="text"
				bind:value={domain}
				placeholder="https://mikrocloud.example.com"
			/>
		</Card.Content>
	</Card.Root>

	<Card.Root class="gap-3">
		<Card.Header>
			<Card.Title class="text-xl">Timezone</Card.Title>
		</Card.Header>
		<Card.Content class="space-y-4">
			<p class="text-xs">This is used for the update check and automatic update frequency.</p>
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
		</Card.Content>
	</Card.Root>

	<Card.Root class="gap-3">
		<Card.Header>
			<Card.Title class="text-xl">Addressing</Card.Title>
		</Card.Header>
		<Card.Content class="space-y-4">
			<Field.Set>
				<Field.Group>
					<Field.Field>
						<Field.Label for="ipv4">Public IPv4 Address</Field.Label>
						<div class="flex gap-2">
							<Input
								id="ipv4"
								type="text"
								bind:value={ipv4}
								placeholder="192.168.1.100"
								class="flex-1"
							/>
							<Button variant="outline" onclick={handleDetectIPs} disabled={isDetecting} size="sm">
								{isDetecting ? 'Detecting...' : 'Auto-detect'}
							</Button>
						</div>
						<Field.Description>
							Enter the IPv4 address of the instance. It is useful if you have several IPv4
							addresses.
						</Field.Description>
					</Field.Field>

					<Field.Field>
						<Field.Label for="ipv6">Public IPv6 Address</Field.Label>

						<div class="flex gap-2">
							<Input
								id="ipv6"
								type="text"
								bind:value={ipv6}
								placeholder="2001:db8::1"
								class="flex-1"
							/>
							<Button variant="outline" onclick={handleDetectIPs} disabled={isDetecting} size="sm">
								{isDetecting ? 'Detecting...' : 'Auto-detect'}
							</Button>
						</div>
						<Field.Description>
							Enter the IPv6 address of the instance. It is useful if you have several IPv6
							addresses.
						</Field.Description>
					</Field.Field>
				</Field.Group>
			</Field.Set>
		</Card.Content>
	</Card.Root>

	<Card.Root class="gap-3">
		<Card.Header>
			<Card.Title class="text-xl">Allow Registrations</Card.Title>
		</Card.Header>
		<Card.Content class="space-y-4">
			<p class="text-xs">
				If disabled, the signup button will be removed and the registration route will be blocked.
			</p>
			<Label for="allow-registrations" class="hidden sr-only">Allow Registrations</Label>
			<div class="flex items-center gap-3">
				<Switch id="allow-registrations" bind:checked={allowRegistrations} />
				{#if allowRegistrations}
					<span>Enabled</span>
				{:else}
					<span>Disabled</span>
				{/if}
			</div>
		</Card.Content>
	</Card.Root>

	<!-- <div class="flex justify-end"> -->
	<!-- 	<Button onclick={handleSave} disabled={updateSettingsMutation.isPending}> -->
	<!-- 		{updateSettingsMutation.isPending ? 'Saving...' : 'Save Changes'} -->
	<!-- 	</Button> -->
	<!-- </div> -->
</div>

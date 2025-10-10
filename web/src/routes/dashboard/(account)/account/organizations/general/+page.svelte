<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';

	let organizationName = $state('Default Organization');
	let organizationSlug = $state('default');
	let description = $state('');
	let billingEmail = $state('admin@example.com');

	function handleSave() {
		console.log('Saving organization settings...');
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-3xl font-bold mb-2">General Settings</h1>
		<p class="text-muted-foreground">
			Manage your organization's basic information and settings.
		</p>
	</div>

	<Card>
		<CardHeader>
			<CardTitle>Organization Information</CardTitle>
			<CardDescription>Update your organization's name and description</CardDescription>
		</CardHeader>
		<CardContent class="space-y-4">
			<div class="space-y-2">
				<Label for="org-name">Organization Name</Label>
				<Input id="org-name" bind:value={organizationName} placeholder="My Organization" />
			</div>
			<div class="space-y-2">
				<Label for="org-slug">Organization Slug</Label>
				<Input id="org-slug" bind:value={organizationSlug} placeholder="my-organization" />
				<p class="text-sm text-muted-foreground">
					Used in URLs and API endpoints. Must be unique and lowercase.
				</p>
			</div>
			<div class="space-y-2">
				<Label for="description">Description</Label>
				<Textarea
					id="description"
					bind:value={description}
					placeholder="A brief description of your organization"
					rows={3}
				/>
			</div>
		</CardContent>
	</Card>

	<Card>
		<CardHeader>
			<CardTitle>Billing Information</CardTitle>
			<CardDescription>Configure billing and payment settings</CardDescription>
		</CardHeader>
		<CardContent class="space-y-4">
			<div class="space-y-2">
				<Label for="billing-email">Billing Email</Label>
				<Input
					id="billing-email"
					type="email"
					bind:value={billingEmail}
					placeholder="billing@example.com"
				/>
				<p class="text-sm text-muted-foreground">
					Invoices and billing notifications will be sent to this email.
				</p>
			</div>
			<div class="space-y-2">
				<Label>Current Plan</Label>
				<div class="flex items-center justify-between p-4 border rounded-lg">
					<div>
						<p class="font-medium">Free Plan</p>
						<p class="text-sm text-muted-foreground">Limited to 5 projects</p>
					</div>
					<Button variant="outline">Upgrade Plan</Button>
				</div>
			</div>
		</CardContent>
	</Card>

	<Card>
		<CardHeader>
			<CardTitle>Organization ID</CardTitle>
			<CardDescription>Your unique organization identifier</CardDescription>
		</CardHeader>
		<CardContent>
			<div class="flex items-center gap-2">
				<Input value="org-00000000-0000-0000-0000-000000000000" readonly class="font-mono text-sm" />
				<Button variant="outline" size="sm">Copy</Button>
			</div>
		</CardContent>
	</Card>

	<Button onclick={handleSave}>Save Changes</Button>
</div>

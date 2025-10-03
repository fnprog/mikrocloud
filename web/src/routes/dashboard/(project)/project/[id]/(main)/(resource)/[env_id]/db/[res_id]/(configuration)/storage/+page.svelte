<script lang="ts">
	import { page } from '$app/stores';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Badge } from '$lib/components/ui/badge';
	import { HardDrive, Plus, Trash2 } from 'lucide-svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { databasesApi } from '$lib/api/databases';

	const projectId = $derived($page.params.id);
	const resId = $derived($page.params.res_id);

	const databaseQuery = createQuery({
		queryKey: ['database', projectId, resId],
		queryFn: () => databasesApi.get(projectId, resId),
		enabled: !!projectId && !!resId
	});

	const database = $derived($databaseQuery.data);

	const mountedVolumes = $state([
		{ id: '1', name: '/var/lib/postgresql/data', size: '10 GB', used: '2.3 GB' }
	]);
</script>

	<div class="space-y-6">
		<div>
			<h2 class="text-2xl font-bold tracking-tight">Persistent Storage</h2>
			<p class="text-muted-foreground">Manage persistent storage volumes for {database?.name || 'database'}</p>
		</div>

	<Card>
		<CardHeader>
			<CardTitle>Mounted Volumes</CardTitle>
			<CardDescription>
				Persistent volumes attached to this database container
			</CardDescription>
		</CardHeader>
		<CardContent>
			<div class="space-y-4">
				{#if mountedVolumes.length === 0}
					<div class="text-center py-8 text-muted-foreground">
						<HardDrive class="mx-auto h-12 w-12 mb-2 opacity-50" />
						<p>No persistent volumes mounted</p>
					</div>
				{:else}
					<div class="space-y-3">
						{#each mountedVolumes as volume}
							<div class="flex items-center justify-between p-4 border rounded-lg">
								<div class="flex items-center gap-3">
									<HardDrive class="h-5 w-5 text-muted-foreground" />
									<div>
										<p class="font-medium">{volume.name}</p>
										<p class="text-sm text-muted-foreground">
											{volume.used} / {volume.size} used
										</p>
									</div>
								</div>
								<div class="flex items-center gap-2">
									<Badge variant="secondary">Active</Badge>
									<Button variant="ghost" size="icon">
										<Trash2 class="h-4 w-4" />
									</Button>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		</CardContent>
	</Card>

	<Card>
		<CardHeader>
			<CardTitle>Add New Volume</CardTitle>
			<CardDescription>Mount a new persistent volume to this database</CardDescription>
		</CardHeader>
		<CardContent>
			<form class="space-y-4">
				<div class="grid gap-4 sm:grid-cols-2">
					<div class="space-y-2">
						<Label for="mount-path">Mount Path</Label>
						<Input
							id="mount-path"
							placeholder="/data"
							type="text"
						/>
					</div>
					<div class="space-y-2">
						<Label for="volume-size">Size (GB)</Label>
						<Input
							id="volume-size"
							placeholder="10"
							type="number"
							min="1"
						/>
					</div>
				</div>
				<Button type="submit" class="w-full sm:w-auto">
					<Plus class="h-4 w-4 mr-2" />
					Add Volume
				</Button>
			</form>
		</CardContent>
	</Card>
</div>

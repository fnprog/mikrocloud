<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import * as InputGroup from '$lib/components/ui/input-group/index.js';
	import ProjectCard from '$lib/components/projects/project-card.svelte';
	import AddProjectModal from '$lib/components/projects/add-project-modal.svelte';
	import { Plus, Search } from 'lucide-svelte';
	import { createProjectsQuery } from '$lib/features/projects/queries';

	import {
		createProjectMutationQuery,
		deleteProjectMutationQuery
	} from '$lib/features/projects/mutations';
	import type { Project } from '$lib/features/projects/types';
	import { Skeleton } from '$lib/components/ui/skeleton';

	const projectsQuery = createProjectsQuery();
	let projects = $state<Project[]>([]);
	let searchQuery = $state('');
	let showAddModal = $state(false);

	$effect(() => {
		if (!projectsQuery.isLoading && !projectsQuery.isError) {
			projects = projectsQuery.data ?? [];
		}
	});

	const onProjectDeletionError = (err: Error) => {
		alert(err?.message ?? 'Failed to delete project');
	};

	const createProjectMutation = createProjectMutationQuery();
	const deleteProjectMutation = deleteProjectMutationQuery({ onError: onProjectDeletionError });

	const handleCreateProject = async (data: { name: string; description?: string }) => {
		console.log(data);
		createProjectMutation.mutate(data);
	};

	const handleDeleteProject = async (id: string) => {
		if (!confirm('Are you sure you want to delete this project?')) return;
		deleteProjectMutation.mutate(id);
	};

	const filteredProjects = $derived(
		projects.filter(
			(p) =>
				p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				p.description?.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);
</script>

<div class="min-h-screen">
	<div class="border-b">
		<div class="max-w-7xl mx-auto px-6 flex items-center justify-between">
			<div class="my-[40px] w-full mx-auto">
				<h1 class="text-3xl font-semibold">Projects</h1>
			</div>
			<Button size="lg" onclick={() => (showAddModal = true)}>
				<Plus class="size-4" />
				Add Project
			</Button>
		</div>
	</div>
	<div class="mt-[46px]"></div>

	<div class="max-w-7xl mx-auto px-6 space-y-6">
		<InputGroup.Root>
			<InputGroup.Input
				placeholder="Search projects..."
				bind:value={searchQuery}
				disabled={projectsQuery.isLoading}
			/>
			<InputGroup.Addon>
				<Search />
			</InputGroup.Addon>
		</InputGroup.Root>

		<!-- Loading skeleton: show a grid of pulsing card placeholders -->
		{#if projectsQuery.isLoading}

			<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
				{#each Array(6) as _}
					<div class="bg-card border border-border rounded-lg p-6">
						<div class="flex items-center gap-3 mb-4">
							<Skeleton class="w-12 h-12 rounded-lg" />
							<div class="flex-1">
								<Skeleton class="h-4 rounded mb-2 w-3/4" />
								<Skeleton class="h-3 rounded w-1/2" />
							</div>
						</div>
						<Skeleton class="h-3 rounded w-1/3" />
					</div>
				{/each}
			</div>
		{:else if projectsQuery.isError}
			<!-- Improved error state with a clear message and retry CTA -->
			<div class="flex min-h-[400px] items-center justify-center">
				<div class="max-w-lg w-full">
					<div class="border border-destructive/10 bg-destructive/5 rounded-lg p-6 text-center">
						<h2 class="text-destructive font-semibold mb-2">Failed to load projects</h2>
						<p class="text-muted-foreground mb-4">
							{projectsQuery.error ?? 'An unexpected error occurred while fetching projects.'}
						</p>
						<div class="flex items-center justify-center gap-2">
							<Button variant="outline" onclick={projectsQuery.refetch}>Retry</Button>
							<Button onclick={() => (showAddModal = true)}>Create Project</Button>
						</div>
					</div>
				</div>
			</div>
		{:else if filteredProjects.length === 0}
			<!-- Improved empty state with CTA -->
			<div class="flex min-h-[400px] items-center justify-center">
				<div class="text-center">
					{#if searchQuery}
						<p class="text-muted-foreground">No projects found matching "{searchQuery}"</p>
					{:else}
						<div class="max-w-md">
							<h3 class="text-lg font-semibold mb-2">No projects yet</h3>
							<p class="text-muted-foreground mb-4">
								Add your first project to organize your work and invite collaborators.
							</p>
							<div class="flex items-center justify-center gap-2">
								<Button onclick={() => (showAddModal = true)}>
									<Plus class="size-4 mr-2" />
									Create Project
								</Button>
								<Button variant="ghost" onclick={() => goto('/dashboard/project/explore')}>
									Explore Templates
								</Button>
							</div>
						</div>
					{/if}
				</div>
			</div>
		{:else}
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
				{#each filteredProjects as project (project.id)}
					<ProjectCard
						{project}
						onDelete={handleDeleteProject}
						onclick={() => goto(`/dashboard/project/${project.id}/overview`)}
					/>
				{/each}
			</div>
		{/if}
	</div>
</div>

<AddProjectModal bind:open={showAddModal} onSubmit={handleCreateProject} />

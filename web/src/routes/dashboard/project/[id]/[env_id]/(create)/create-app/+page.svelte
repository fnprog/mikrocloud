<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

	import { Button, buttonVariants } from '$lib/components/ui/button';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { ChevronDown, ChevronLeft } from 'lucide-svelte';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as Collapsible from '$lib/components/ui/collapsible/index.js';
	import SourceTypeSelector from '$lib/components/applications/source-type-selector.svelte';
	import GitConfigForm from '$lib/components/applications/git-config-form.svelte';
	import DockerConfigForm from '$lib/components/applications/docker-config-form.svelte';
	import EnvVarsForm from '$lib/components/applications/env-vars-form.svelte';
	import { createApplicationMutationQuery } from '$lib/features/applications/mutations';
	import type { CreateApplicationRequest } from '$lib/features/applications/types';

	let projectId = page.params.id!;
	let envId = page.params.env_id!;
	let showSelectPhase = $state(true);

	let sourceType = $state<'git' | 'docker' | 'zip'>('git');
	let buildType = $state<'auto' | 'nixpacks' | 'static' | 'dockerfile' | 'compose'>('auto');

	let publishDirectory = $state('dist');
	let isStatic = $state(false);
	let isSpa = $state(false);

	let appName = $state('');
	let gitBranch = $state('main');

	let dockerfilePath = $state('Dockerfile');
	let composePath = $state('docker-compose.yml');
	let basePath = $state('/');

	let source_repository = $state('');
	let public_repository = $state('');
	let branch = $state('main');

	let dockerType = $state<'dockerfile' | 'compose'>('dockerfile');
	let dockerContent = $state('');
	let dockerfileUploadMethod = $state<'paste' | 'upload'>('paste');
	let zipFile = $state<File | null>(null);

	let registryImage = $state('');
	let registryTag = $state('latest');

	interface EnvVar {
		key: string;
		value: string;
	}

	let envVars = $state<EnvVar[]>([]);

	const createApplicationMutation = createApplicationMutationQuery({
		onSuccess: () => {
			goto(`/dashboard/project/${projectId}`);
		}
	});

	const buildPacks = [
		{ value: 'auto', label: 'Auto' },
		{ value: 'nixpacks', label: 'NixPacks' },
		{ value: 'static', label: 'Static' },
		{ value: 'containerfile', label: 'ContainerFile' },
		{ value: 'compose', label: 'Compose' }
	];

	//HELPER: Will be changed in the future
	const branches = [{ value: 'main', label: 'Main' }];

	const buildPackLabel = $derived(
		buildPacks.find((d) => d.value === buildType)?.label ?? 'Choose department'
	);

	function handleSubmit() {
		const deploymentSource =
			sourceType === 'git'
				? {
						type: 'git' as const,
						git_repo: {
							url: public_repository,
							branch,
							path: basePath
						}
					}
				: {
						type: 'upload' as const,
						upload: {
							filename:
								sourceType === 'docker'
									? dockerType === 'dockerfile'
										? 'Dockerfile'
										: 'docker-compose.yml'
									: zipFile?.name || 'upload.zip',
							file_path: ''
						}
					};

		const buildpack =
			buildType === 'static'
				? {
						type: 'static',
						config: {
							output_dir: publishDirectory,
							is_static: isStatic,
							is_spa: isSpa
						}
					}
				: buildType === 'dockerfile'
					? {
							type: 'dockerfile',
							config: {
								dockerfile_path: dockerfilePath
							}
						}
					: buildType === 'compose'
						? {
								type: 'docker-compose',
								config: {
									compose_file: composePath
								}
							}
						: {
								type: buildType,
								config: {}
							};

		const data: CreateApplicationRequest = {
			name: appName,
			environment_id: envId,
			project_id: projectId,
			deployment_source: deploymentSource,
			buildpack,
			env_vars: envVars
		};

		createApplicationMutation.mutate(data);
	}
</script>

<div class="container mx-auto max-w-7xl py-8 px-4">
	<div class="space-y-8">
		<h1 class="text-3xl font-bold">Time to Ship</h1>

		{#if showSelectPhase}
			<!-- Phase 1: Select a source type -->
			<div class="space-y-6">
				<p class="text-sm text-muted-foreground">Select the source for your new web service</p>
				<SourceTypeSelector bind:selected={sourceType} />
				<div class="flex justify-end">
					<Button onclick={() => (showSelectPhase = false)}>Continue</Button>
				</div>
			</div>
		{:else}
			<!-- Phase 2: Detailed deployment layout -->
			<div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
				<!-- Left column: contextual info / compact nav -->
				<div class="lg:col-span-4">
					<div class="sticky top-24 space-y-6">
						<div>
							<p class="text-sm font-medium">Source</p>
							<p class="text-xs text-muted-foreground">
								{sourceType === 'git'
									? 'Git repository'
									: sourceType === 'docker'
										? 'Dockerfile / Compose / Registry Image'
										: 'Upload archive'}
							</p>
						</div>

						<div class="space-y-4 text-sm text-muted-foreground">
							<div>
								<p class="font-medium">Name</p>
								<p class="text-xs">
									{appName === '' ? 'A unique name for your web service' : appName}
								</p>
							</div>

							{#if sourceType !== 'docker'}
								<div>
									<p class="font-medium">Build Pack</p>
									<p class="text-xs">
										{buildType}
									</p>
								</div>
							{/if}
							{#if sourceType === 'git'}
								<div>
									<p class="font-medium">Branch</p>
									<p class="text-xs">{gitBranch}</p>
								</div>
							{/if}
							<div>
								<p class="font-medium">Root</p>
								<p class="text-xs">{basePath}</p>
							</div>
						</div>

						<div class="pt-4">
							<Button variant="secondary" onclick={() => (showSelectPhase = true)}
								><ChevronLeft size={16} /> Back</Button
							>
						</div>
					</div>
				</div>

				<!-- Right column: forms -->
				<div class="lg:col-span-8 space-y-10">
					{#if sourceType === 'git'}
						<GitConfigForm
							bind:source_repository_url={source_repository}
							bind:public_repository_url={public_repository}
						/>
					{:else if sourceType === 'docker'}
						<DockerConfigForm
							bind:method={dockerfileUploadMethod}
							bind:fileType={dockerType}
							bind:content={dockerContent}
						/>
					{:else if sourceType === 'zip'}
						<Card>
							<CardHeader>
								<CardTitle>Upload file</CardTitle>
								<CardDescription
									>Upload a zipped archive containing your application code</CardDescription
								>
							</CardHeader>
							<CardContent>
								<div class="space-y-2">
									<Label for="zip-file">Zip file</Label>
									<Input
										id="zip-file"
										type="file"
										accept=".zip,.tar,.tar.gz,.tgz"
										onchange={(e) => {
											const files = e.currentTarget.files;
											if (files && files.length > 0) {
												zipFile = files[0];
											}
										}}
										required
									/>
									<p class="text-xs text-muted-foreground">
										Supported formats: .zip, .tar, .tar.gz, .tgz
									</p>
								</div>
							</CardContent>
						</Card>
					{/if}

					<div class="w-full">
						<Field.Field>
							<Field.Label for="app-name">Application name</Field.Label>
							<Input id="app-name" placeholder="my-app" bind:value={appName} required />
							<Field.Description>A unique name to identify your application.</Field.Description>
						</Field.Field>
					</div>

					{#if sourceType === 'git' || sourceType === 'zip'}
						<div class="w-full">
							<Field.Field>
								<Field.Label for="build-type">Build pack</Field.Label>
								<Select.Root type="single" bind:value={buildType}>
									<Select.Trigger id="build-type">
										{buildPackLabel}
									</Select.Trigger>
									<Select.Content>
										{#each buildPacks as buildpack (buildpack.value)}
											<Select.Item {...buildpack} />
										{/each}
									</Select.Content>
								</Select.Root>
							</Field.Field>
						</div>
					{/if}

					{#if sourceType === 'git'}
						<div class="w-full">
							<Field.Field>
								<Field.Label for="branch">Git Branch</Field.Label>
								<Select.Root type="single" bind:value={buildType}>
									<Select.Trigger id="branch">
										{branch || 'Select a branch'}
									</Select.Trigger>
									<Select.Content>
										{#each branches as branchOption}
											<Select.Item value={branchOption.value}>{branchOption.label}</Select.Item>
										{/each}
									</Select.Content>
								</Select.Root>
							</Field.Field>
						</div>
					{/if}

					<div class="w-full">
						<Field.Field>
							<Field.Label for="root">Root directory</Field.Label>
							<Input id="root" placeholder="e.g. src" bind:value={basePath} required />
						</Field.Field>
					</div>

					<Collapsible.Root
						class="w-full rounded-md  border border-input bg-background dark:bg-input/30"
					>
						<div class="flex items-center w-full">
							<Collapsible.Trigger
								class={buttonVariants({
									variant: 'secondary',
									size: 'lg',
									class:
										'rounded-md items-center justify-start space-x-4 w-full border-input bg-background dark:bg-transparent px-3 py-1'
								})}
							>
								<ChevronDown />
								<h4 class="text-sm font-semibold">Environment Variables</h4>
								<span class="sr-only">Toggle</span>
							</Collapsible.Trigger>
						</div>
						<Collapsible.Content class="space-y-2">
							<EnvVarsForm bind:envVars />
						</Collapsible.Content>
					</Collapsible.Root>

					<Button onclick={handleSubmit}>Create Application</Button>

					{#if createApplicationMutation.isError}
						<div class="text-sm text-destructive mt-4">
							Failed to create application: {createApplicationMutation.error?.message ||
								'Unknown error'}
						</div>
					{/if}
				</div>
			</div>
		{/if}
	</div>
</div>

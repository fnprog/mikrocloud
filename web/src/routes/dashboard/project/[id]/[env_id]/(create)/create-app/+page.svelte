<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

	import { Button } from '$lib/components/ui/button';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Select, SelectContent, SelectItem, SelectTrigger } from '$lib/components/ui/select';
	import { Loader2 } from 'lucide-svelte';

	import SourceTypeSelector from '$lib/components/applications/source-type-selector.svelte';
	import GitConfigForm from '$lib/components/applications/git-config-form.svelte';
	import DockerConfigForm from '$lib/components/applications/docker-config-form.svelte';
	import BuildTypeSelector from '$lib/components/applications/build-type-selector.svelte';
	import { createApplicationMutationQuery } from '$lib/features/applications/mutations';
	import type { CreateApplicationRequest } from '$lib/features/applications/types';

	let projectId = page.params.id!;
	let envId = page.params.env_id!;

	let sourceType = $state<'git' | 'docker' | 'zip'>('git');
	let buildType = $state<'nixpacks' | 'heroku' | 'paketo' | 'static' | 'dockerfile' | 'compose'>(
		'nixpacks'
	);
	let publishDirectory = $state('dist');

	let appName = $state('');
	let appDescription = $state('');
	let location = $state('local');

	let gitProvider = $state<'github' | 'gitlab' | 'bitbucket' | 'custom'>('github');
	let dockerfilePath = $state('Dockerfile');
	let composePath = $state('docker-compose.yml');
	let basePath = $state('/');
	let customGitUrl = $state('');
	let repository = $state('');
	let branch = $state('main');
	let autoDeploy = $state(true);
	let isPrivate = $state(false);

	let dockerType = $state<'dockerfile' | 'compose'>('dockerfile');
	let dockerfileContent = $state('');
	let dockerfileUploadMethod = $state<'paste' | 'upload'>('paste');
	let zipFile = $state<File | null>(null);

	const createApplicationMutation = createApplicationMutationQuery({
		onSuccess: () => {
			goto(`/dashboard/project/${projectId}`);
		}
	});

	function handleSubmit() {
		const deploymentSource =
			sourceType === 'git'
				? {
						type: 'git' as const,
						git_repo: {
							url:
								gitProvider === 'custom'
									? customGitUrl
									: `https://${gitProvider}.com/${repository}.git`,
							branch,
							path: basePath
						}
					}
				: sourceType === 'docker'
					? {
							type: 'docker' as const,
							config: {
								type: dockerType,
								content: dockerType === 'dockerfile' ? dockerfileContent : composeContent
							}
						}
					: {
							type: 'zip' as const,
							config: {
								file: zipFile
							}
						};

		const buildpack =
			buildType === 'static'
				? {
						type: 'static',
						config: {
							output_dir: publishDirectory
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
			description: appDescription,
			environment_id: envId,
			project_id: projectId,
			deployment_source: deploymentSource,
			buildpack
		};

		createApplicationMutation.mutate(data);
	}
</script>

<div class="container mx-auto max-w-4xl py-8 px-4">
	<div class="space-y-8">
		<div>
			<h1 class="text-3xl font-bold">Create new application</h1>
			<p class="text-muted-foreground mt-2">
				Deploy your application from a Git repository or Docker image
			</p>
		</div>

		<Card>
			<CardHeader>
				<CardTitle>Deployment source</CardTitle>
				<CardDescription>Choose where your application code comes from</CardDescription>
			</CardHeader>
			<CardContent>
				<SourceTypeSelector selected={sourceType} onSelect={(type) => (sourceType = type)} />
			</CardContent>
		</Card>

		{#if sourceType === 'git'}
			<Card>
				<CardHeader>
					<CardTitle>Git configuration</CardTitle>
					<CardDescription>Configure your repository settings</CardDescription>
				</CardHeader>
				<CardContent>
					<GitConfigForm
						provider={gitProvider}
						onProviderChange={(p) => (gitProvider = p)}
						{repository}
						onRepositoryChange={(r) => (repository = r)}
						{branch}
						onBranchChange={(b) => (branch = b)}
						{autoDeploy}
						onAutoDeployChange={(a) => (autoDeploy = a)}
						{isPrivate}
						onIsPrivateChange={(p) => (isPrivate = p)}
						{customGitUrl}
						onCustomGitUrlChange={(u) => (customGitUrl = u)}
						{basePath}
						onBasePathChange={(p) => (basePath = p)}
					/>
				</CardContent>
			</Card>

			<Card>
				<CardHeader>
					<CardTitle>Build configuration</CardTitle>
					<CardDescription>Choose how to build your application</CardDescription>
				</CardHeader>
				<CardContent>
					<BuildTypeSelector
						selected={buildType}
						onSelect={(type) => (buildType = type)}
						{publishDirectory}
						onPublishDirectoryChange={(dir) => (publishDirectory = dir)}
					/>
				</CardContent>
			</Card>
		{:else if sourceType === 'docker'}
			<Card>
				<CardHeader>
					<CardTitle>Docker configuration</CardTitle>
					<CardDescription>Provide your Docker configuration</CardDescription>
				</CardHeader>
				<CardContent>
					<DockerConfigForm
						method={dockerfileUploadMethod}
						fileType={dockerType}
						content={dockerfileContent}
					/>
				</CardContent>
			</Card>
		{:else if sourceType === 'zip'}
			<Card>
				<CardHeader>
					<CardTitle>Upload file</CardTitle>
					<CardDescription>Upload a zipped archive containing your application code</CardDescription
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

			<Card>
				<CardHeader>
					<CardTitle>Build configuration</CardTitle>
					<CardDescription>Choose how to build your application</CardDescription>
				</CardHeader>
				<CardContent>
					<BuildTypeSelector
						selected={buildType}
						onSelect={(type) => (buildType = type)}
						{publishDirectory}
						onPublishDirectoryChange={(dir) => (publishDirectory = dir)}
					/>
				</CardContent>
			</Card>
		{/if}

		<Card>
			<CardHeader>
				<CardTitle>Application details</CardTitle>
				<CardDescription>Configure your application settings</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="space-y-2">
					<Label for="app-name">Application name</Label>
					<Input id="app-name" placeholder="my-app" bind:value={appName} required />
				</div>

				<div class="space-y-2">
					<Label for="app-description">Description (optional)</Label>
					<Input
						id="app-description"
						placeholder="A brief description of your application"
						bind:value={appDescription}
					/>
				</div>

				<div class="space-y-2">
					<Label for="location">Location</Label>
					<Select type="single" bind:value={location}>
						<SelectTrigger id="location">
							{location || 'Select location'}
						</SelectTrigger>
						<SelectContent>
							<SelectItem value="local">Local</SelectItem>
							<SelectItem value="us-east">US East</SelectItem>
							<SelectItem value="us-west">US West</SelectItem>
							<SelectItem value="eu-west">EU West</SelectItem>
						</SelectContent>
					</Select>
				</div>
			</CardContent>
		</Card>

		<div class="flex gap-4">
			<Button
				onclick={handleSubmit}
				disabled={!appName ||
					(sourceType === 'git'
						? !repository
						: sourceType === 'docker'
							? !dockerfileContent && !composeContent
							: !zipFile) ||
					createApplicationMutation.isPending}
				class="flex-1"
			>
				{#if createApplicationMutation.isPending}
					<Loader2 class="h-4 w-4 mr-2 animate-spin" />
					Creating...
				{:else}
					Create application
				{/if}
			</Button>
		</div>

		{#if createApplicationMutation.isError}
			<div class="text-sm text-destructive">
				Failed to create application: {createApplicationMutation.error?.message || 'Unknown error'}
			</div>
		{/if}
	</div>
</div>

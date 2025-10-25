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
	import { Loader2, ChevronLeft, ChevronRight } from 'lucide-svelte';

	import Stepper from '$lib/components/ui/stepper/stepper.svelte';
	import SourceTypeSelector from '$lib/components/applications/source-type-selector.svelte';
	import GitConfigForm from '$lib/components/applications/git-config-form.svelte';
	import DockerConfigForm from '$lib/components/applications/docker-config-form.svelte';
	import RegistryConfigForm from '$lib/components/applications/registry-config-form.svelte';
	import BuildTypeSelector from '$lib/components/applications/build-type-selector.svelte';
	import EnvVarsForm from '$lib/components/applications/env-vars-form.svelte';
	import PortsForm from '$lib/components/applications/ports-form.svelte';
	import AppDetailsForm from '$lib/components/applications/app-details-form.svelte';
	import { createApplicationMutationQuery } from '$lib/features/applications/mutations';
	import type { CreateApplicationRequest } from '$lib/features/applications/types';

	let projectId = page.params.id!;
	let envId = page.params.env_id!;

	let currentStep = $state(0);

	let sourceType = $state<'git' | 'docker' | 'registry' | 'zip'>('git');
	let buildType = $state<'nixpacks' | 'static' | 'dockerfile' | 'compose'>('nixpacks');
	let publishDirectory = $state('dist');
	let isStatic = $state(false);
	let isSpa = $state(false);

	let appName = $state('');
	let appDescription = $state('');

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
	let dockerContent = $state('');
	let dockerfileUploadMethod = $state<'paste' | 'upload'>('paste');
	let zipFile = $state<File | null>(null);

	let registryImage = $state('');
	let registryTag = $state('latest');

	interface EnvVar {
		key: string;
		value: string;
	}

	interface PortMapping {
		containerPort: number;
		hostPort: number;
		protocol: 'tcp' | 'udp';
	}

	let buildTimeVars = $state<EnvVar[]>([]);
	let runtimeVars = $state<EnvVar[]>([]);
	let portMappings = $state<PortMapping[]>([]);

	const steps = [
		{ title: 'Source Type', description: 'Choose deployment source' },
		{ title: 'Source Config', description: 'Configure source details' },
		{ title: 'Build Config', description: 'Build configuration' },
		{ title: 'Environment', description: 'Environment variables' },
		{ title: 'Ports', description: 'Exposed ports' },
		{ title: 'Details', description: 'Application details' },
		{ title: 'Review', description: 'Review & create' }
	];

	const createApplicationMutation = createApplicationMutationQuery({
		onSuccess: () => {
			goto(`/dashboard/project/${projectId}`);
		}
	});

	function canProceed(step: number): boolean {
		switch (step) {
			case 0:
				return true;
			case 1:
				if (sourceType === 'git') return !!repository;
				if (sourceType === 'registry') return !!registryImage;
				if (sourceType === 'docker') return !!dockerContent;
				if (sourceType === 'zip') return !!zipFile;
				return false;
			case 2:
				return true;
			case 3:
				return true;
			case 4:
				return true;
			case 5:
				return !!appName;
			case 6:
				return !!appName;
			default:
				return false;
		}
	}

	function shouldShowStep(step: number): boolean {
		if (step === 2) {
			return sourceType === 'git' || sourceType === 'zip';
		}
		if (step === 4) {
			return sourceType === 'registry' || sourceType === 'docker';
		}
		return true;
	}

	function nextStep() {
		let next = currentStep + 1;
		while (next < steps.length && !shouldShowStep(next)) {
			next++;
		}
		if (next < steps.length) {
			currentStep = next;
		}
	}

	function previousStep() {
		let prev = currentStep - 1;
		while (prev >= 0 && !shouldShowStep(prev)) {
			prev--;
		}
		if (prev >= 0) {
			currentStep = prev;
		}
	}

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
				: sourceType === 'registry'
					? {
							type: 'registry' as const,
							registry: {
								image: registryImage,
								tag: registryTag
							}
						}
					: {
							type: 'upload' as const,
							upload: {
								filename: sourceType === 'docker' 
									? (dockerType === 'dockerfile' ? 'Dockerfile' : 'docker-compose.yml')
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

		const allEnvVars: Record<string, string> = {};
		[...buildTimeVars, ...runtimeVars].forEach((v) => {
			if (v.key) {
				allEnvVars[v.key] = v.value;
			}
		});

		const data: CreateApplicationRequest = {
			name: appName,
			description: appDescription,
			environment_id: envId,
			project_id: projectId,
			deployment_source: deploymentSource,
			buildpack,
			env_vars: allEnvVars
		};

		createApplicationMutation.mutate(data);
	}
</script>

<div class="container mx-auto max-w-4xl py-8 px-4">
	<div class="space-y-8">
		<div>
			<h1 class="text-3xl font-bold">Create new application</h1>
			<p class="text-muted-foreground mt-2">
				Deploy your application with our step-by-step wizard
			</p>
		</div>

		<Stepper {steps} {currentStep} />

		<div class="min-h-96">
			{#if currentStep === 0}
				<Card>
					<CardHeader>
						<CardTitle>Deployment source</CardTitle>
						<CardDescription>Choose where your application code comes from</CardDescription>
					</CardHeader>
					<CardContent>
						<SourceTypeSelector selected={sourceType} onSelect={(type) => (sourceType = type)} />
					</CardContent>
				</Card>
			{:else if currentStep === 1}
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
				{:else if sourceType === 'docker'}
					<Card>
						<CardHeader>
							<CardTitle>Docker configuration</CardTitle>
							<CardDescription>Provide your Docker configuration</CardDescription>
						</CardHeader>
						<CardContent>
							<DockerConfigForm
								bind:method={dockerfileUploadMethod}
								bind:fileType={dockerType}
								bind:content={dockerContent}
							/>
						</CardContent>
					</Card>
				{:else if sourceType === 'registry'}
					<Card>
						<CardHeader>
							<CardTitle>Container registry</CardTitle>
							<CardDescription>Pull and deploy a container image from a registry</CardDescription>
						</CardHeader>
						<CardContent>
							<RegistryConfigForm
								image={registryImage}
								onImageChange={(img) => (registryImage = img)}
								tag={registryTag}
								onTagChange={(t) => (registryTag = t)}
							/>
						</CardContent>
					</Card>
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
			{:else if currentStep === 2}
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
							{isStatic}
							onIsStaticChange={(val) => (isStatic = val)}
							{isSpa}
							onIsSpaChange={(val) => (isSpa = val)}
						/>
					</CardContent>
				</Card>
			{:else if currentStep === 3}
				<Card>
					<CardHeader>
						<CardTitle>Environment variables</CardTitle>
						<CardDescription>Configure build-time and runtime environment variables</CardDescription>
					</CardHeader>
					<CardContent>
						<EnvVarsForm
							{buildTimeVars}
							onBuildTimeVarsChange={(vars) => (buildTimeVars = vars)}
							{runtimeVars}
							onRuntimeVarsChange={(vars) => (runtimeVars = vars)}
						/>
					</CardContent>
				</Card>
			{:else if currentStep === 4}
				<Card>
					<CardHeader>
						<CardTitle>Port configuration</CardTitle>
						<CardDescription>Configure which ports to expose from your container</CardDescription>
					</CardHeader>
					<CardContent>
						<PortsForm
							{portMappings}
							onPortMappingsChange={(mappings) => (portMappings = mappings)}
						/>
					</CardContent>
				</Card>
			{:else if currentStep === 5}
				<Card>
					<CardHeader>
						<CardTitle>Application details</CardTitle>
						<CardDescription>Name and describe your application</CardDescription>
					</CardHeader>
					<CardContent>
						<AppDetailsForm
							name={appName}
							onNameChange={(n) => (appName = n)}
							description={appDescription}
							onDescriptionChange={(d) => (appDescription = d)}
						/>
					</CardContent>
				</Card>
			{:else if currentStep === 6}
				<Card>
					<CardHeader>
						<CardTitle>Review & create</CardTitle>
						<CardDescription>Review your configuration before creating the application</CardDescription>
					</CardHeader>
					<CardContent class="space-y-6">
						<div class="space-y-4">
							<div>
								<h3 class="font-semibold mb-2">Application</h3>
								<div class="text-sm space-y-1">
									<p><span class="text-muted-foreground">Name:</span> {appName}</p>
									{#if appDescription}
										<p><span class="text-muted-foreground">Description:</span> {appDescription}</p>
									{/if}
								</div>
							</div>

							<div>
								<h3 class="font-semibold mb-2">Source</h3>
								<div class="text-sm space-y-1">
									<p><span class="text-muted-foreground">Type:</span> {sourceType}</p>
									{#if sourceType === 'git'}
										<p><span class="text-muted-foreground">Repository:</span> {repository}</p>
										<p><span class="text-muted-foreground">Branch:</span> {branch}</p>
									{:else if sourceType === 'registry'}
										<p><span class="text-muted-foreground">Image:</span> {registryImage}:{registryTag}</p>
									{/if}
								</div>
							</div>

							{#if sourceType === 'git' || sourceType === 'zip'}
								<div>
									<h3 class="font-semibold mb-2">Build</h3>
									<div class="text-sm space-y-1">
										<p><span class="text-muted-foreground">Type:</span> {buildType}</p>
										{#if buildType === 'static'}
											<p><span class="text-muted-foreground">Output directory:</span> {publishDirectory}</p>
											{#if isStatic}
												<p class="text-muted-foreground">• Static files only (no build)</p>
											{/if}
											{#if isSpa}
												<p class="text-muted-foreground">• Single Page Application mode</p>
											{/if}
										{/if}
									</div>
								</div>
							{/if}

							{#if buildTimeVars.length > 0 || runtimeVars.length > 0}
								<div>
									<h3 class="font-semibold mb-2">Environment Variables</h3>
									<div class="text-sm space-y-1">
										<p class="text-muted-foreground">
											{buildTimeVars.length + runtimeVars.length} variable(s) configured
										</p>
									</div>
								</div>
							{/if}

							{#if portMappings.length > 0}
								<div>
									<h3 class="font-semibold mb-2">Port Mappings</h3>
									<div class="text-sm space-y-1">
										{#each portMappings as mapping}
											<p class="text-muted-foreground">
												• Container {mapping.containerPort} → Host {mapping.hostPort || 'auto'} ({mapping.protocol})
											</p>
										{/each}
									</div>
								</div>
							{/if}
						</div>
					</CardContent>
				</Card>
			{/if}
		</div>

		<div class="flex gap-4">
			<Button
				variant="outline"
				onclick={previousStep}
				disabled={currentStep === 0}
				class="flex-1"
			>
				<ChevronLeft class="h-4 w-4 mr-2" />
				Previous
			</Button>

			{#if currentStep === steps.length - 1}
				<Button
					onclick={handleSubmit}
					disabled={!canProceed(currentStep) || createApplicationMutation.isPending}
					class="flex-1"
				>
					{#if createApplicationMutation.isPending}
						<Loader2 class="h-4 w-4 mr-2 animate-spin" />
						Creating...
					{:else}
						Create application
					{/if}
				</Button>
			{:else}
				<Button
					onclick={nextStep}
					disabled={!canProceed(currentStep)}
					class="flex-1"
				>
					Next
					<ChevronRight class="h-4 w-4 ml-2" />
				</Button>
			{/if}
		</div>

		{#if createApplicationMutation.isError}
			<div class="text-sm text-destructive">
				Failed to create application: {createApplicationMutation.error?.message || 'Unknown error'}
			</div>
		{/if}
	</div>
</div>

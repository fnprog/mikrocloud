<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import { Tabs, TabsContent, TabsList, TabsTrigger } from '$lib/components/ui/tabs';
	import { Upload } from 'lucide-svelte';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';

	interface Props {
		content: string;
		method: 'paste' | 'upload';
		fileType: 'dockerfile' | 'compose' | 'registry';
		image?: string;
		tag?: string;
	}

	let {
		content = $bindable(),
		method = $bindable(),
		fileType = $bindable(),
		image = $bindable(),
		tag = $bindable()
	}: Props = $props();

	let imageInput = $state('');
	let tagInput = $state('latest');

	$effect(() => {
		if (image !== imageInput) imageInput = image ?? '';
		if (tag !== tagInput) tagInput = tag ?? 'latest';
	});

	$effect(() => {
		image = imageInput;
		tag = tagInput;
	});

	let fileInput: HTMLInputElement;

	function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		const file = target.files?.[0];
		if (file) {
			const reader = new FileReader();
			reader.onload = (e) => {
				content = e.target?.result as string;
			};
			reader.readAsText(file);
		}
	}
</script>

<div class="space-y-6">
	<Tabs bind:value={fileType}>
		<TabsList class="grid w-full grid-cols-3">
			<TabsTrigger value="dockerfile">Dockerfile</TabsTrigger>
			<TabsTrigger value="compose">Docker Compose</TabsTrigger>
			<TabsTrigger value="registry">Existing Image</TabsTrigger>
		</TabsList>

		<!-- Registry tab content -->
		<TabsContent value="registry" class="space-y-4 mt-4">
			<Field.Set>
				<Field.Group>
					<Field.Field>
						<Field.Label for="registry_image">Project name</Field.Label>
						<Input id="image" placeholder="e.g. nginx" bind:value={imageInput} />
					</Field.Field>
				</Field.Group>
			</Field.Set>

			<DropdownMenu.Root>
				<DropdownMenu.Trigger>
					{#snippet child({ props })}
						<Button {...props} variant="outline" class="w-full flex-1">Open</Button>
					{/snippet}
				</DropdownMenu.Trigger>
				<DropdownMenu.Content>
					<DropdownMenu.Item>Add New source</DropdownMenu.Item>
				</DropdownMenu.Content>
			</DropdownMenu.Root>

			<!-- <div class="flex items-center justify-between"> -->
			<!-- 	<p class="text-sm text-muted-foreground">Select credentials</p> -->
			<!-- 	<Button variant="ghost" size="sm">Add credential</Button> -->
			<!-- </div> -->
			<!-- <div class="text-xs text-muted-foreground"> -->
			<!-- 	Note: You can use registry credentials to access private images. -->
			<!-- </div> -->
		</TabsContent>
	</Tabs>

	{#if fileType !== 'registry'}
		<Tabs bind:value={method}>
			<TabsList class="grid w-full grid-cols-2">
				<TabsTrigger value="paste">Paste content</TabsTrigger>
				<TabsTrigger value="upload">Upload file</TabsTrigger>
			</TabsList>

			<TabsContent value="paste" class="space-y-2 mt-4">
				<Label for="docker-content" class="hidden sr-only">
					{fileType === 'dockerfile' ? 'Dockerfile' : 'Docker Compose'} content
				</Label>
				<Textarea
					id="docker-content"
					placeholder={fileType === 'dockerfile'
						? 'FROM node:18\nWORKDIR /app\nCOPY . .\nRUN npm install\nCMD ["npm", "start"]'
						: 'version: "3.8"\nservices:\n  app:\n    build: .\n    ports:\n      - "3000:3000"'}
					bind:value={content}
					rows={12}
					class="font-mono text-sm"
				/>
			</TabsContent>

			<TabsContent value="upload" class="space-y-4 mt-4">
				<div
					class="border-2 border-dashed rounded-lg p-8 text-center hover:border-primary/50 transition-colors"
				>
					<Upload class="h-8 w-8 mx-auto mb-4 text-muted-foreground" />
					<p class="text-sm text-muted-foreground mb-4">
						Click to upload or drag and drop your {fileType === 'dockerfile'
							? 'Dockerfile'
							: 'docker-compose.yml'}
					</p>
					<input
						type="file"
						bind:this={fileInput}
						onchange={handleFileSelect}
						accept={fileType === 'dockerfile' ? '*' : '.yml,.yaml'}
						class="hidden"
					/>
					<Button type="button" variant="outline" onclick={() => fileInput?.click()}>
						Select File
					</Button>
				</div>
				{#if content}
					<div class="bg-muted p-4 rounded-lg">
						<p class="text-sm font-medium mb-2">File content preview:</p>
						<pre
							class="text-xs overflow-auto max-h-48 bg-background p-3 rounded border">{content}</pre>
					</div>
				{/if}
			</TabsContent>
		</Tabs>
	{/if}
</div>

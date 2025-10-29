<script lang="ts">
	import { Github, FileArchive } from 'lucide-svelte';
	import { IconBrandDocker } from '@tabler/icons-svelte';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as RadioGroup from '$lib/components/ui/radio-group/index.js';

	interface Props {
		selected: 'git' | 'docker' | 'zip';
	}

	let { selected = $bindable('git') }: Props = $props();

	let hoveredId = $state<string | null>(null);
</script>

<Field.Group>
	<Field.Set>
		<RadioGroup.Root bind:value={selected}>
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
				<Field.Label
					for="git"
					class="cursor-pointer relative min-h-48  overflow-hidden"
					onmouseenter={() => (hoveredId = 'git')}
					onmouseleave={() => (hoveredId = null)}
				>
					<Field.Field orientation="horizontal">
						<Field.Content>
							<Field.Title>Git repository</Field.Title>
							<Field.Description>Deploy from a Git repository</Field.Description>
						</Field.Content>
						<RadioGroup.Item value="git" id="git" />
					</Field.Field>
					<div class="absolute -bottom-15 -right-15 w-50 h-50 flex items-center justify-center">
						<img
							src="/GitHub.svg"
							alt="git"
							width="200"
							height="200"
							class={`w-full h-full object-contain transition-all duration-300 ${
								hoveredId === 'git' ? '-rotate-3' : '-rotate-12'
							}`}
						/>
					</div>
				</Field.Label>
				<Field.Label
					for="docker"
					class="cursor-pointer relative min-h-48  overflow-hidden"
					onmouseenter={() => (hoveredId = 'docker')}
					onmouseleave={() => (hoveredId = null)}
				>
					<Field.Field orientation="horizontal">
						<Field.Content>
							<Field.Title>Docker image</Field.Title>
							<Field.Description
								>Use a Dockerfile, Docker Compose or Pull from a registry</Field.Description
							>
						</Field.Content>
						<RadioGroup.Item value="docker" id="docker" />
					</Field.Field>
					<div
						class="absolute -bottom-15 -right-15 w-50 h-50 flex items-center justify-center overflow-hidden"
					>
						<img
							src="/docker.png"
							alt="git"
							width="200"
							height="200"
							class={`w-full h-full object-contain transition-all duration-300 ${
								hoveredId === 'docker' ? '-rotate-3' : '-rotate-12'
							}`}
						/>
					</div>
				</Field.Label>
				<Field.Label
					for="zip"
					class="cursor-pointer relative min-h-48  overflow-hidden"
					onmouseenter={() => (hoveredId = 'zip')}
					onmouseleave={() => (hoveredId = null)}
				>
					<Field.Field orientation="horizontal">
						<Field.Content>
							<Field.Title>Zipped file</Field.Title>
							<Field.Description>Upload a compressed archive</Field.Description>
						</Field.Content>
						<RadioGroup.Item value="zip" id="zip" />
					</Field.Field>
					<div
						class="absolute -bottom-15 -right-15 w-50 h-50 flex items-center justify-center overflow-hidden"
					>
						<img
							src="/zip.png"
							alt="git"
							width="200"
							height="200"
							class={`w-full h-full object-contain transition-all duration-300 ${
								hoveredId === 'zip' ? '-rotate-3' : '-rotate-12'
							}`}
						/>
					</div>
				</Field.Label>
			</div>
		</RadioGroup.Root>
	</Field.Set>
</Field.Group>

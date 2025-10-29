<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Plus, Trash2 } from 'lucide-svelte';

	interface EnvVar {
		key: string;
		value: string;
	}

	interface Props {
		envVars: EnvVar[];
	}

	let { envVars = $bindable([]) }: Props = $props();

	function addEnvVar() {
		envVars = [...envVars, { key: '', value: '' }];
	}

	function removeBuildTimeVar(index: number) {
		const newVars = envVars.filter((_, i) => i !== index);
		envVars = newVars;
	}

	function updateBuildTimeVar(index: number, field: 'key' | 'value', newValue: string) {
		const newVars = envVars.map((v, i) => (i === index ? { ...v, [field]: newValue } : v));
		envVars = newVars;
	}
</script>

<div class="space-y-6">
	<div class="space-y-4 p-2 px-3">
		{#if envVars.length === 0}
			<p class="text-sm text-center text-muted-foreground">No environment variables defined</p>
		{:else}
			{#each envVars as envVar, index}
				<div class="flex gap-2">
					<div class="flex-1">
						<Input
							placeholder="KEY"
							value={envVar.key}
							oninput={(e) => updateBuildTimeVar(index, 'key', e.currentTarget.value)}
						/>
					</div>
					<div class="flex-1">
						<Input
							placeholder="value"
							value={envVar.value}
							oninput={(e) => updateBuildTimeVar(index, 'value', e.currentTarget.value)}
						/>
					</div>
					<Button
						variant="ghost"
						size="icon"
						onclick={() => removeBuildTimeVar(index)}
						class="shrink-0"
					>
						<Trash2 class="h-4 w-4" />
					</Button>
				</div>
			{/each}
		{/if}
		<Button variant="outline" size="sm" onclick={addEnvVar}>
			<Plus class="h-4 w-4 mr-2" />
			Add variable
		</Button>
	</div>
</div>

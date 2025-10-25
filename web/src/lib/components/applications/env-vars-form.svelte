<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Plus, Trash2 } from 'lucide-svelte';

	interface EnvVar {
		key: string;
		value: string;
	}

	interface Props {
		buildTimeVars: EnvVar[];
		onBuildTimeVarsChange: (vars: EnvVar[]) => void;
		runtimeVars: EnvVar[];
		onRuntimeVarsChange: (vars: EnvVar[]) => void;
	}

	let { buildTimeVars, onBuildTimeVarsChange, runtimeVars, onRuntimeVarsChange }: Props = $props();

	function addBuildTimeVar() {
		onBuildTimeVarsChange([...buildTimeVars, { key: '', value: '' }]);
	}

	function removeBuildTimeVar(index: number) {
		const newVars = buildTimeVars.filter((_, i) => i !== index);
		onBuildTimeVarsChange(newVars);
	}

	function updateBuildTimeVar(index: number, field: 'key' | 'value', newValue: string) {
		const newVars = buildTimeVars.map((v, i) =>
			i === index ? { ...v, [field]: newValue } : v
		);
		onBuildTimeVarsChange(newVars);
	}

	function addRuntimeVar() {
		onRuntimeVarsChange([...runtimeVars, { key: '', value: '' }]);
	}

	function removeRuntimeVar(index: number) {
		const newVars = runtimeVars.filter((_, i) => i !== index);
		onRuntimeVarsChange(newVars);
	}

	function updateRuntimeVar(index: number, field: 'key' | 'value', newValue: string) {
		const newVars = runtimeVars.map((v, i) =>
			i === index ? { ...v, [field]: newValue } : v
		);
		onRuntimeVarsChange(newVars);
	}
</script>

<div class="space-y-6">
	<Card>
		<CardHeader>
			<CardTitle>Build-time variables</CardTitle>
			<CardDescription>
				These variables are available during the build process
			</CardDescription>
		</CardHeader>
		<CardContent>
			<div class="space-y-4">
				{#if buildTimeVars.length === 0}
					<p class="text-sm text-muted-foreground">No build-time variables defined</p>
				{:else}
					{#each buildTimeVars as envVar, index}
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
								class="flex-shrink-0"
							>
								<Trash2 class="h-4 w-4" />
							</Button>
						</div>
					{/each}
				{/if}
				<Button variant="outline" size="sm" onclick={addBuildTimeVar}>
					<Plus class="h-4 w-4 mr-2" />
					Add variable
				</Button>
			</div>
		</CardContent>
	</Card>

	<Card>
		<CardHeader>
			<CardTitle>Runtime variables</CardTitle>
			<CardDescription>
				These variables are available when your application is running
			</CardDescription>
		</CardHeader>
		<CardContent>
			<div class="space-y-4">
				{#if runtimeVars.length === 0}
					<p class="text-sm text-muted-foreground">No runtime variables defined</p>
				{:else}
					{#each runtimeVars as envVar, index}
						<div class="flex gap-2">
							<div class="flex-1">
								<Input
									placeholder="KEY"
									value={envVar.key}
									oninput={(e) => updateRuntimeVar(index, 'key', e.currentTarget.value)}
								/>
							</div>
							<div class="flex-1">
								<Input
									placeholder="value"
									value={envVar.value}
									oninput={(e) => updateRuntimeVar(index, 'value', e.currentTarget.value)}
								/>
							</div>
							<Button
								variant="ghost"
								size="icon"
								onclick={() => removeRuntimeVar(index)}
								class="flex-shrink-0"
							>
								<Trash2 class="h-4 w-4" />
							</Button>
						</div>
					{/each}
				{/if}
				<Button variant="outline" size="sm" onclick={addRuntimeVar}>
					<Plus class="h-4 w-4 mr-2" />
					Add variable
				</Button>
			</div>
		</CardContent>
	</Card>
</div>

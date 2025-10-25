<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Select, SelectContent, SelectItem, SelectTrigger } from '$lib/components/ui/select';
	import { Plus, Trash2 } from 'lucide-svelte';

	interface PortMapping {
		containerPort: number;
		hostPort: number;
		protocol: 'tcp' | 'udp';
	}

	interface Props {
		portMappings: PortMapping[];
		onPortMappingsChange: (mappings: PortMapping[]) => void;
	}

	let { portMappings, onPortMappingsChange }: Props = $props();

	function addPortMapping() {
		onPortMappingsChange([
			...portMappings,
			{ containerPort: 80, hostPort: 0, protocol: 'tcp' }
		]);
	}

	function removePortMapping(index: number) {
		const newMappings = portMappings.filter((_, i) => i !== index);
		onPortMappingsChange(newMappings);
	}

	function updatePortMapping(
		index: number,
		field: 'containerPort' | 'hostPort' | 'protocol',
		value: number | string
	) {
		const newMappings = portMappings.map((m, i) =>
			i === index ? { ...m, [field]: field === 'protocol' ? value : Number(value) } : m
		);
		onPortMappingsChange(newMappings);
	}
</script>

<div class="space-y-4">
	<div>
		<Label class="text-base font-semibold">Port mappings</Label>
		<p class="text-sm text-muted-foreground mt-1">
			Configure which container ports to expose and their host port mappings
		</p>
	</div>

	{#if portMappings.length === 0}
		<p class="text-sm text-muted-foreground">No port mappings defined</p>
	{:else}
		<div class="space-y-3">
			{#each portMappings as mapping, index}
				<div class="flex gap-2 items-end">
					<div class="flex-1 space-y-2">
						<Label for={`container-port-${index}`} class="text-xs">Container port</Label>
						<Input
							id={`container-port-${index}`}
							type="number"
							min="1"
							max="65535"
							value={mapping.containerPort}
							oninput={(e) =>
								updatePortMapping(index, 'containerPort', e.currentTarget.value)}
						/>
					</div>
					<div class="flex-1 space-y-2">
						<Label for={`host-port-${index}`} class="text-xs">Host port (0 = auto)</Label>
						<Input
							id={`host-port-${index}`}
							type="number"
							min="0"
							max="65535"
							value={mapping.hostPort}
							oninput={(e) => updatePortMapping(index, 'hostPort', e.currentTarget.value)}
						/>
					</div>
					<div class="w-24 space-y-2">
						<Label for={`protocol-${index}`} class="text-xs">Protocol</Label>
						<Select
							type="single"
							value={mapping.protocol}
							onValueChange={(val) => {
								if (val) updatePortMapping(index, 'protocol', val);
							}}
						>
							<SelectTrigger id={`protocol-${index}`}>
								{mapping.protocol.toUpperCase()}
							</SelectTrigger>
							<SelectContent>
								<SelectItem value="tcp">TCP</SelectItem>
								<SelectItem value="udp">UDP</SelectItem>
							</SelectContent>
						</Select>
					</div>
					<Button
						variant="ghost"
						size="icon"
						onclick={() => removePortMapping(index)}
						class="flex-shrink-0"
					>
						<Trash2 class="h-4 w-4" />
					</Button>
				</div>
			{/each}
		</div>
	{/if}

	<Button variant="outline" size="sm" onclick={addPortMapping}>
		<Plus class="h-4 w-4 mr-2" />
		Add port mapping
	</Button>

	<p class="text-xs text-muted-foreground mt-2">
		Set host port to 0 to automatically assign an available port
	</p>
</div>

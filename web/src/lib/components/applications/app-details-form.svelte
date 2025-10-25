<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';

	interface Props {
		name: string;
		onNameChange: (value: string) => void;
		description: string;
		onDescriptionChange: (value: string) => void;
	}

	let { name, onNameChange, description, onDescriptionChange }: Props = $props();
</script>

<div class="space-y-4">
	<div>
		<Label class="text-base font-semibold">Application details</Label>
		<p class="text-sm text-muted-foreground mt-1">
			Configure your application name and description
		</p>
	</div>

	<div class="space-y-2">
		<Label for="app-name">Application name</Label>
		<Input
			id="app-name"
			placeholder="my-app"
			value={name}
			oninput={(e) => onNameChange(e.currentTarget.value)}
			required
		/>
		<p class="text-xs text-muted-foreground">
			A unique name to identify your application
		</p>
	</div>

	<div class="space-y-2">
		<Label for="app-description">Description (optional)</Label>
		<Input
			id="app-description"
			placeholder="A brief description of your application"
			value={description}
			oninput={(e) => onDescriptionChange(e.currentTarget.value)}
		/>
	</div>
</div>

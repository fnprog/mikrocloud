<script lang="ts" module>
	export interface StepperPanelProps {
		class?: string;
		children: Snippet;
	}
</script>

<script lang="ts">
	import { getContext } from 'svelte';
	import type { Snippet } from 'svelte';
	import type { StepperContext } from './stepper.svelte';

	let { class: className = '', children }: StepperPanelProps = $props();

	const ctx = getContext<StepperContext>('stepper');
	if (!ctx) throw new Error('StepperPanel must be used within Stepper');
</script>

<div data-slot="stepper-panel" data-state={ctx.activeStep} class="w-full {className}">
	{@render children()}
</div>

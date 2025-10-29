<script lang="ts" module>
	export interface StepperNavProps {
		class?: string;
		children: Snippet;
	}
</script>

<script lang="ts">
	import { getContext } from 'svelte';
	import type { Snippet } from 'svelte';
	import type { StepperContext } from './stepper.svelte';

	let { class: className = '', children }: StepperNavProps = $props();

	const ctx = getContext<StepperContext>('stepper');
	if (!ctx) throw new Error('StepperNav must be used within Stepper');
</script>

<nav
	data-slot="stepper-nav"
	data-state={ctx.activeStep}
	data-orientation={ctx.orientation}
	class="group/stepper-nav inline-flex data-[orientation=horizontal]:w-full data-[orientation=horizontal]:flex-row data-[orientation=vertical]:flex-col {className}"
>
	{@render children()}
</nav>

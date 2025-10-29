<script lang="ts" module>
	export interface StepperItemProps {
		step: number;
		completed?: boolean;
		disabled?: boolean;
		loading?: boolean;
		class?: string;
		children: Snippet;
	}
</script>

<script lang="ts">
	import { getContext, setContext } from 'svelte';
	import type { Snippet } from 'svelte';
	import type { StepperContext, StepState, StepItemContext } from './stepper.svelte';

	let {
		step,
		completed = false,
		disabled = false,
		loading = false,
		class: className = '',
		children
	}: StepperItemProps = $props();

	const stepperCtx = getContext<StepperContext>('stepper');
	if (!stepperCtx) throw new Error('StepperItem must be used within Stepper');

	const state = $derived<StepState>(
		completed || step < stepperCtx.activeStep
			? 'completed'
			: stepperCtx.activeStep === step
				? 'active'
				: 'inactive'
	);

	const isLoading = $derived(loading && step === stepperCtx.activeStep);

	const itemContext: StepItemContext = {
		step,
		get state() {
			return state;
		},
		isDisabled: disabled,
		get isLoading() {
			return isLoading;
		}
	};

	setContext('stepItem', itemContext);
</script>

<div
	data-slot="stepper-item"
	data-state={state}
	data-loading={isLoading || undefined}
	class="group/step flex items-center justify-center group-data-[orientation=horizontal]/stepper-nav:flex-row group-data-[orientation=vertical]/stepper-nav:flex-col not-last:flex-1 {className}"
>
	{@render children()}
</div>

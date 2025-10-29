<script lang="ts" module>
	export interface StepperTriggerProps {
		class?: string;
		tabIndex?: number;
		children: Snippet;
	}
</script>

<script lang="ts">
	import { getContext, onMount, onDestroy } from 'svelte';
	import type { Snippet } from 'svelte';
	import type { StepperContext, StepItemContext } from './stepper.svelte';

	let { class: className = '', tabIndex, children }: StepperTriggerProps = $props();

	const stepperCtx = getContext<StepperContext>('stepper');
	const itemCtx = getContext<StepItemContext>('stepItem');

	if (!stepperCtx || !itemCtx) {
		throw new Error('StepperTrigger must be used within StepperItem');
	}

	let btnRef: HTMLButtonElement;

	const isSelected = $derived(stepperCtx.activeStep === itemCtx.step);
	const id = `stepper-tab-${itemCtx.step}`;
	const panelId = `stepper-panel-${itemCtx.step}`;

	const myIdx = $derived(stepperCtx.triggerNodes.findIndex((n) => n === btnRef));

	const handleKeyDown = (e: KeyboardEvent) => {
		if (stepperCtx.disableTrigger) return;
		
		switch (e.key) {
			case 'ArrowRight':
			case 'ArrowDown':
				e.preventDefault();
				if (myIdx !== -1) stepperCtx.focusNext(myIdx);
				break;
			case 'ArrowLeft':
			case 'ArrowUp':
				e.preventDefault();
				if (myIdx !== -1) stepperCtx.focusPrev(myIdx);
				break;
			case 'Home':
				e.preventDefault();
				stepperCtx.focusFirst();
				break;
			case 'End':
				e.preventDefault();
				stepperCtx.focusLast();
				break;
			case 'Enter':
			case ' ':
				e.preventDefault();
				stepperCtx.setActiveStep(itemCtx.step);
				break;
		}
	};

	onMount(() => {
		if (btnRef) {
			stepperCtx.registerTrigger(btnRef);
		}
	});

	onDestroy(() => {
		if (btnRef) {
			stepperCtx.unregisterTrigger(btnRef);
		}
	});
</script>

	<button
		bind:this={btnRef}
		role="tab"
		{id}
		aria-selected={isSelected}
		aria-controls={panelId}
		tabindex={typeof tabIndex === 'number' ? tabIndex : isSelected ? 0 : -1}
		data-slot="stepper-trigger"
		data-state={itemCtx.state}
		data-loading={itemCtx.isLoading || undefined}
		disabled={itemCtx.isDisabled || stepperCtx.disableTrigger}
		class="cursor-pointer focus-visible:border-ring focus-visible:ring-ring/50 inline-flex items-center gap-3 rounded-full outline-none focus-visible:z-10 focus-visible:ring-[3px] disabled:pointer-events-none disabled:opacity-60 {className}"
		onclick={() => !stepperCtx.disableTrigger && stepperCtx.setActiveStep(itemCtx.step)}
		onkeydown={handleKeyDown}
	>
	{@render children()}
</button>

<script lang="ts">
	import { setContext, getContext } from 'svelte';
	import type { Snippet } from 'svelte';

	// Types
	export type StepperOrientation = 'horizontal' | 'vertical';
	export type StepState = 'active' | 'completed' | 'inactive' | 'loading';

	export interface StepIndicators {
		active?: Snippet;
		completed?: Snippet;
		inactive?: Snippet;
		loading?: Snippet;
	}

	export interface StepperContext {
		activeStep: number;
		setActiveStep: (step: number) => void;
		stepsCount: number;
		orientation: StepperOrientation;
		registerTrigger: (node: HTMLButtonElement) => void;
		unregisterTrigger: (node: HTMLButtonElement) => void;
		triggerNodes: HTMLButtonElement[];
		focusNext: (currentIdx: number) => void;
		focusPrev: (currentIdx: number) => void;
		focusFirst: () => void;
		focusLast: () => void;
		indicators: StepIndicators;
		disableTrigger: boolean;
	}

	export interface StepItemContext {
		step: number;
		state: StepState;
		isDisabled: boolean;
		isLoading: boolean;
	}

	interface Props {
		defaultValue?: number;
		value?: number;
		onValueChange?: (value: number) => void;
		orientation?: StepperOrientation;
		class?: string;
		indicators?: StepIndicators;
		disableTrigger?: boolean;
		children: Snippet;
	}

	let {
		defaultValue = 1,
		value = $bindable(defaultValue),
		onValueChange,
		orientation = 'horizontal',
		class: className = '',
		indicators = {},
		disableTrigger = false,
		children
	}: Props = $props();

	let triggerNodes = $state<HTMLButtonElement[]>([]);

	const registerTrigger = (node: HTMLButtonElement) => {
		if (!triggerNodes.includes(node)) {
			triggerNodes = [...triggerNodes, node];
		}
	};

	const unregisterTrigger = (node: HTMLButtonElement) => {
		triggerNodes = triggerNodes.filter((n) => n !== node);
	};

	const setActiveStep = (step: number) => {
		value = step;
		onValueChange?.(step);
	};

	const focusTrigger = (idx: number) => {
		if (triggerNodes[idx]) triggerNodes[idx].focus();
	};

	const focusNext = (currentIdx: number) => {
		focusTrigger((currentIdx + 1) % triggerNodes.length);
	};

	const focusPrev = (currentIdx: number) => {
		focusTrigger((currentIdx - 1 + triggerNodes.length) % triggerNodes.length);
	};

	const focusFirst = () => focusTrigger(0);
	const focusLast = () => focusTrigger(triggerNodes.length - 1);

	const context: StepperContext = {
		get activeStep() {
			return value;
		},
		setActiveStep,
		stepsCount: 0,
		orientation,
		registerTrigger,
		unregisterTrigger,
		get triggerNodes() {
			return triggerNodes;
		},
		focusNext,
		focusPrev,
		focusFirst,
		focusLast,
		indicators,
		get disableTrigger() {
			return disableTrigger;
		}
	};

	setContext('stepper', context);
</script>

<div
	role="tablist"
	aria-orientation={orientation}
	data-slot="stepper"
	data-orientation={orientation}
	class="w-full {className}"
>
	{@render children()}
</div>

<script lang="ts">
	import { cn } from '$lib/utils';
	import { Check } from 'lucide-svelte';

	interface Props {
		steps: Array<{ id: number; title: string; description?: string }>;
		currentStep: number;
		class?: string;
	}

	let { steps, currentStep, class: className }: Props = $props();
</script>

<div class={cn('w-full', className)}>
	<nav aria-label="Progress">
		<ol class="flex items-center justify-between">
			{#each steps as step, index}
				{@const isCompleted = currentStep > step.id}
				{@const isCurrent = currentStep === step.id}
				{@const isUpcoming = currentStep < step.id}

				<li class="flex-1 {index !== steps.length - 1 ? 'pr-8' : ''} relative">
					<div class="flex flex-col items-center">
						<div class="flex items-center w-full">
							<div
								class={cn(
									'flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full border-2 transition-all',
									isCompleted && 'border-primary bg-primary text-primary-foreground',
									isCurrent && 'border-primary bg-background text-primary',
									isUpcoming && 'border-muted-foreground/30 bg-background text-muted-foreground'
								)}
							>
								{#if isCompleted}
									<Check class="h-5 w-5" />
								{:else}
									<span class="text-sm font-medium">{step.id}</span>
								{/if}
							</div>
							{#if index !== steps.length - 1}
								<div
									class={cn(
										'absolute top-5 left-[calc(50%+20px)] right-[-32px] h-0.5 transition-all',
										currentStep > step.id ? 'bg-primary' : 'bg-muted-foreground/30'
									)}
								></div>
							{/if}
						</div>
						<div class="mt-2 text-center">
							<p
								class={cn(
									'text-sm font-medium transition-colors',
									isCurrent && 'text-primary',
									isCompleted && 'text-foreground',
									isUpcoming && 'text-muted-foreground'
								)}
							>
								{step.title}
							</p>
							{#if step.description}
								<p class="text-xs text-muted-foreground mt-0.5">{step.description}</p>
							{/if}
						</div>
					</div>
				</li>
			{/each}
		</ol>
	</nav>
</div>

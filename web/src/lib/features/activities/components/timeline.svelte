<script lang="ts">
	import { Clock, Package, Zap, Pause, AlertCircle } from 'lucide-svelte';
	import { Avatar, AvatarImage, AvatarFallback } from '$lib/components/ui/avatar';
	import { cn } from '$lib/utils';
	import type { Activity, ActivityLevel } from '../types';

	interface ActivityLink {
		text: string;
		href?: string;
	}

	interface Props {
		items: Activity[];
	}

	let { items = [] }: Props = $props();

  const getErrorLevelColor = (level: ActivityLevel) => {
    switch (level) {
      case 'error':
        return 'border-red-500';
      case 'warn':
        return 'border-orange-500';
      case 'success':
        return 'border-emerald-500';
      case 'info':
      default:
        return 'border-gray-500';
    }
  };

  const getErrorLevelIconColor = (level: ActivityLevel) => {
    switch (level) {
      case 'error':
        return 'text-red-500';
      case 'warn':
        return 'text-orange-500';
      case 'success':
        return 'text-emerald-500';
      case 'info':
      default:
        return 'text-gray-500';
    }
  };

	const getResourceIcon = (resourceType?: string) => {
		switch (resourceType) {
			case 'deployment':
				return Package;
			case 'service':
				return Zap;
			case 'app':
				return Zap;
			case 'pause':
				return Pause;
			default:
				return AlertCircle;
		}
	};

	const isActivityLink = (value: string | ActivityLink): value is ActivityLink => {
		return typeof value === 'object' && value !== null;
	};
</script>

<div class="w-full">
	<!-- Timeline Container -->
	<div class="border border-gray-700 rounded-2xl p-6 bg-black/40 backdrop-blur overflow-hidden">
		<div class="relative">
			<!-- Timeline line using svelte-cn for conditional styling -->
			<div
				class={cn(
					'absolute left-5 top-0 bottom-0 w-0.5 bg-linear-to-b from-red-500 via-teal-500 to-orange-500 -z-10 -my-6'
				)}
			></div>

			<!-- Activity Items -->
			<div class="space-y-8">
				{#each items as item, index (item.id)}
					{@const IconComponent = getResourceIcon(undefined)}
					<div class="relative pl-20">
						<!-- Circle with icon using svelte-cn for dynamic classes -->
						<div
							class={cn(
								'absolute left-0 top-0 w-10 h-10 rounded-full border-2 flex items-center justify-center bg-black/80 z-10',
								getErrorLevelColor(item.level)
							)}
						>
							<IconComponent class={cn('w-5 h-5', getErrorLevelIconColor(item.level))} />
						</div>

						<!-- Content -->
						<div class="space-y-2">
							<div class="flex items-start justify-between">
								<div class="flex items-center gap-3">
									<Avatar class="w-8 h-8">
										<AvatarImage
											src={item.initiator?.avatar || '/placeholder.svg'}
											alt={item.initiator?.name}
										/>
										<AvatarFallback class="text-xs">{item.initiator?.initials}</AvatarFallback>
									</Avatar>
									<div>
										<h3 class="text-white font-semibold">{item.initiator?.name}</h3>
									</div>
								</div>
								<div class="flex items-center gap-2 text-gray-400 text-sm">
									<Clock class="w-4 h-4" />
									{item.timestamp}
								</div>
							</div>

							<p class="text-gray-300">
								{item.eventType}
								{#if item.description}
									<span class="text-gray-400"> {item.description}</span>
								{/if}
							</p>

							<!-- Metadata rendering with svelte-cn for conditional styling -->
							{#if item.metadata && Object.entries(item.metadata).length > 0}
								<div class="mt-3 space-y-2">
									{#each Object.entries(item.metadata) as [key, value] (key)}
										<div
											class={cn(
												'inline-flex items-center gap-2 px-3 py-2 bg-gray-900/50 border border-gray-700 rounded-lg mr-2'
											)}
										>
											<div
												class={cn(
													'w-5 h-5 rounded bg-teal-500/20 flex items-center justify-center text-teal-400 text-xs font-bold'
												)}
											>
												⬚
											</div>
											{#if isActivityLink(value)}
												<a
													href={value.href || '#'}
													class="text-emerald-400 hover:text-emerald-300 text-sm font-mono transition-colors"
												>
													{value.text}
												</a>
											{:else}
												<span class="text-gray-300 text-sm font-mono">{value}</span>
											{/if}
										</div>
									{/each}
								</div>
							{/if}
						</div>

						<!-- Divider -->
						{#if index < items.length - 1}
							<div class="mt-6 border-t border-gray-700"></div>
						{/if}
					</div>
				{/each}
			</div>
		</div>
	</div>
</div>

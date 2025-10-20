<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Switch } from '$lib/components/ui/switch';
	import { Palette, Moon, Sun, Monitor, Sparkles, Eye } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	type Theme = 'dark' | 'light' | 'auto';

	let selectedTheme = $state<Theme>('dark');
	let reducedMotion = $state(false);
	let compactMode = $state(false);

	$effect(() => {
		const savedTheme = localStorage.getItem('theme') as Theme | null;
		if (savedTheme) {
			selectedTheme = savedTheme;
		}

		const savedReducedMotion = localStorage.getItem('reducedMotion');
		if (savedReducedMotion) {
			reducedMotion = savedReducedMotion === 'true';
		}

		const savedCompactMode = localStorage.getItem('compactMode');
		if (savedCompactMode) {
			compactMode = savedCompactMode === 'true';
		}
	});

	function handleThemeChange(theme: Theme) {
		selectedTheme = theme;
		localStorage.setItem('theme', theme);

		if (theme === 'auto') {
			const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
			applyTheme(prefersDark ? 'dark' : 'light');
		} else {
			applyTheme(theme);
		}

		toast.success(`Theme changed to ${theme}`);
	}

	function applyTheme(theme: 'dark' | 'light') {
		if (theme === 'dark') {
			document.documentElement.classList.add('dark');
		} else {
			document.documentElement.classList.remove('dark');
		}
	}

	function handleReducedMotionToggle() {
		reducedMotion = !reducedMotion;
		localStorage.setItem('reducedMotion', String(reducedMotion));

		if (reducedMotion) {
			document.documentElement.style.setProperty('--transition-duration', '0ms');
		} else {
			document.documentElement.style.removeProperty('--transition-duration');
		}

		toast.success(reducedMotion ? 'Reduced motion enabled' : 'Reduced motion disabled');
	}

	function handleCompactModeToggle() {
		compactMode = !compactMode;
		localStorage.setItem('compactMode', String(compactMode));

		toast.success(compactMode ? 'Compact mode enabled' : 'Compact mode disabled');
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold">Appearance</h1>
		<p class="text-sm text-muted-foreground mt-1">Customize the look and feel of your dashboard</p>
	</div>

	<div class="border border-border rounded-lg p-6 space-y-6">
		<div>
			<h2 class="text-lg font-semibold mb-4 flex items-center gap-2">
				<Palette class="w-5 h-5" />
				Theme
			</h2>
			<p class="text-sm text-muted-foreground mb-4">
				Select the theme that best suits your preference
			</p>

			<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
				<button
					class="relative flex flex-col items-center gap-3 p-6 border-2 rounded-lg transition-all hover:border-primary {selectedTheme ===
					'dark'
						? 'border-primary bg-primary/5'
						: 'border-border'}"
					onclick={() => handleThemeChange('dark')}
				>
					<div
						class="w-full h-24 rounded-md bg-gradient-to-br from-slate-900 to-slate-800 border border-slate-700 flex items-center justify-center"
					>
						<Moon class="w-8 h-8 text-slate-400" />
					</div>
					<div class="text-center">
						<p class="font-medium">Dark</p>
						<p class="text-xs text-muted-foreground">Dark mode theme</p>
					</div>
					{#if selectedTheme === 'dark'}
						<div
							class="absolute top-2 right-2 w-5 h-5 rounded-full bg-primary flex items-center justify-center"
						>
							<svg class="w-3 h-3 text-primary-foreground" fill="currentColor" viewBox="0 0 12 12">
								<path d="M4.5 9L1.5 6L2.55 4.95L4.5 6.9L9.45 1.95L10.5 3L4.5 9Z" />
							</svg>
						</div>
					{/if}
				</button>

				<button
					class="relative flex flex-col items-center gap-3 p-6 border-2 rounded-lg transition-all hover:border-primary {selectedTheme ===
					'light'
						? 'border-primary bg-primary/5'
						: 'border-border'}"
					onclick={() => handleThemeChange('light')}
				>
					<div
						class="w-full h-24 rounded-md bg-gradient-to-br from-slate-50 to-slate-100 border border-slate-300 flex items-center justify-center"
					>
						<Sun class="w-8 h-8 text-slate-600" />
					</div>
					<div class="text-center">
						<p class="font-medium">Light</p>
						<p class="text-xs text-muted-foreground">Light mode theme</p>
					</div>
					{#if selectedTheme === 'light'}
						<div
							class="absolute top-2 right-2 w-5 h-5 rounded-full bg-primary flex items-center justify-center"
						>
							<svg class="w-3 h-3 text-primary-foreground" fill="currentColor" viewBox="0 0 12 12">
								<path d="M4.5 9L1.5 6L2.55 4.95L4.5 6.9L9.45 1.95L10.5 3L4.5 9Z" />
							</svg>
						</div>
					{/if}
				</button>

				<button
					class="relative flex flex-col items-center gap-3 p-6 border-2 rounded-lg transition-all hover:border-primary {selectedTheme ===
					'auto'
						? 'border-primary bg-primary/5'
						: 'border-border'}"
					onclick={() => handleThemeChange('auto')}
				>
					<div
						class="w-full h-24 rounded-md bg-gradient-to-r from-slate-900 via-slate-500 to-slate-50 border border-slate-400 flex items-center justify-center"
					>
						<Monitor class="w-8 h-8 text-slate-200" />
					</div>
					<div class="text-center">
						<p class="font-medium">Auto</p>
						<p class="text-xs text-muted-foreground">Follow system theme</p>
					</div>
					{#if selectedTheme === 'auto'}
						<div
							class="absolute top-2 right-2 w-5 h-5 rounded-full bg-primary flex items-center justify-center"
						>
							<svg class="w-3 h-3 text-primary-foreground" fill="currentColor" viewBox="0 0 12 12">
								<path d="M4.5 9L1.5 6L2.55 4.95L4.5 6.9L9.45 1.95L10.5 3L4.5 9Z" />
							</svg>
						</div>
					{/if}
				</button>
			</div>
		</div>
	</div>

	<div class="border border-border rounded-lg p-6 space-y-6">
		<div>
			<h2 class="text-lg font-semibold mb-4 flex items-center gap-2">
				<Eye class="w-5 h-5" />
				Display Preferences
			</h2>

			<div class="space-y-4">
				<div class="flex items-center justify-between p-4 border border-border rounded-lg">
					<div class="flex items-center gap-3">
						<Sparkles class="w-5 h-5 text-muted-foreground" />
						<div>
							<p class="font-medium">Reduced Motion</p>
							<p class="text-sm text-muted-foreground">
								Minimize animations and transitions for better performance
							</p>
						</div>
					</div>
					<Switch checked={reducedMotion} onCheckedChange={handleReducedMotionToggle} />
				</div>

				<div class="flex items-center justify-between p-4 border border-border rounded-lg">
					<div class="flex items-center gap-3">
						<Monitor class="w-5 h-5 text-muted-foreground" />
						<div>
							<p class="font-medium">Compact Mode</p>
							<p class="text-sm text-muted-foreground">
								Reduce spacing and padding for a more condensed interface
							</p>
						</div>
					</div>
					<Switch checked={compactMode} onCheckedChange={handleCompactModeToggle} />
				</div>
			</div>
		</div>
	</div>

	<div class="border border-border rounded-lg p-6">
		<div>
			<h2 class="text-lg font-semibold mb-2 flex items-center gap-2">
				<Palette class="w-5 h-5" />
				Color Customization
			</h2>
			<p class="text-sm text-muted-foreground mb-4">
				Advanced theme customization options coming soon
			</p>
			<Button variant="outline" disabled>
				<Palette class="w-4 h-4 mr-2" />
				Customize Colors
			</Button>
		</div>
	</div>
</div>

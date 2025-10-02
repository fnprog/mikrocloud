<script lang="ts">
	let { children } = $props();
</script>

<nav class="flex items-center space-x-8 text-sm">
	<a href="/dashboard" class="text-white border-b-2 border-white pb-3">Overview</a>
	<a href="/dashboard/projects" class="text-gray-400 hover:text-white pb-3">Projects</a>
	<a href="/dashboard/monitoring" class="text-gray-400 hover:text-white pb-3">Monitoring</a>
	<a href="/dashboard/settings" class="text-gray-400 hover:text-white pb-3">Settings</a>
</nav>

<div class="p-6">
	{@render children()}
</div>

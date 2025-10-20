<script lang="ts">
	type Props = {
		data: any;
		unit?: string;
		height?: number;
		color?: string;
	};

	let { data, unit = '', height = 200, color = '#000' }: Props = $props();
</script>

<div style="height: {height}px;" class="flex items-center justify-center border rounded bg-muted/20">
	<p class="text-sm text-muted-foreground">Chart placeholder</p>
</div>

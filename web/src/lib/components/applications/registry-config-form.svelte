<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';

	interface Props {
		image: string;
		onImageChange: (image: string) => void;
		tag: string;
		onTagChange: (tag: string) => void;
	}

	let { image, onImageChange, tag, onTagChange }: Props = $props();
</script>

<div class="space-y-6">
	<div class="space-y-2">
		<Label for="registry-image">Image name</Label>
		<Input
			id="registry-image"
			placeholder="nginx, postgres, node, myregistry.com/myimage"
			value={image}
			oninput={(e) => onImageChange(e.currentTarget.value)}
			required
		/>
		<p class="text-xs text-muted-foreground">
			Enter the image name (e.g., nginx, postgres, or myregistry.com/myimage)
		</p>
	</div>

	<div class="space-y-2">
		<Label for="registry-tag">Tag</Label>
		<Input
			id="registry-tag"
			placeholder="latest"
			value={tag}
			oninput={(e) => onTagChange(e.currentTarget.value)}
		/>
		<p class="text-xs text-muted-foreground">Image tag or version (default: latest)</p>
	</div>
</div>

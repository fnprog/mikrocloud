<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { createSetupStatusQuery } from '$lib/features/auth/queries';

	const setupStatusQuery = createSetupStatusQuery();

	onMount(async () => {
		try {
			const status = await setupStatusQuery.refetch();
			if (status.data?.is_setup === false) {
				goto('/register');
			} else {
				goto('/login');
			}
		} catch (error) {
			// If we can't check setup status, default to login
			goto('/login');
		}
	});
</script>

<div>Redirecting...</div>

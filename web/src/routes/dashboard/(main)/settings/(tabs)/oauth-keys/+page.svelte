<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Card from '$lib/components/ui/card';
	import { Github, GitlabIcon as Gitlab, Chrome } from 'lucide-svelte';

	type OAuthProvider = {
		client_id: string;
		secret: string;
		redirect_url: string;
		base_url?: string;
		tenant?: string;
	};

	let github = $state<OAuthProvider>({
		client_id: '',
		secret: '',
		redirect_url: ''
	});

	let gitlab = $state<OAuthProvider>({
		client_id: '',
		secret: '',
		redirect_url: '',
		base_url: ''
	});

	let bitbucket = $state<OAuthProvider>({
		client_id: '',
		secret: '',
		redirect_url: ''
	});

	let google = $state<OAuthProvider>({
		client_id: '',
		secret: '',
		redirect_url: '',
		tenant: ''
	});

	onMount(async () => {
		try {
			const response = await fetch('/api/settings/oauth');
			if (response.ok) {
				const data = await response.json();
				github = data.github || github;
				gitlab = data.gitlab || gitlab;
				bitbucket = data.bitbucket || bitbucket;
				google = data.google || google;
			}
		} catch (error) {
			console.error('Failed to load OAuth settings:', error);
		}
	});

	async function handleSaveProvider(provider: string, data: OAuthProvider) {
		try {
			const response = await fetch(`/api/settings/oauth/${provider}`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(data)
			});

			if (!response.ok) {
				throw new Error(`Failed to save ${provider} settings`);
			}
		} catch (error) {
			console.error(`Failed to save ${provider} settings:`, error);
		}
	}
</script>

<div class="space-y-6">
	<Card.Root>
		<Card.Header>
			<div class="flex items-center gap-2">
				<Github class="h-5 w-5" />
				<Card.Title>GitHub OAuth</Card.Title>
			</div>
			<Card.Description>Configure GitHub OAuth authentication.</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			<div class="space-y-2">
				<Label for="github-client-id">Client ID</Label>
				<Input id="github-client-id" bind:value={github.client_id} placeholder="Iv1.abc123..." />
			</div>

			<div class="space-y-2">
				<Label for="github-secret">Client Secret</Label>
				<Input
					id="github-secret"
					type="password"
					bind:value={github.secret}
					placeholder="secret123..."
				/>
			</div>

			<div class="space-y-2">
				<Label for="github-redirect">Redirect URL</Label>
				<Input
					id="github-redirect"
					bind:value={github.redirect_url}
					placeholder="https://example.com/auth/github/callback"
				/>
			</div>
		</Card.Content>
		<Card.Footer class="flex justify-end">
			<Button onclick={() => handleSaveProvider('github', github)}>Save GitHub</Button>
		</Card.Footer>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<div class="flex items-center gap-2">
				<Gitlab class="h-5 w-5" />
				<Card.Title>GitLab OAuth</Card.Title>
			</div>
			<Card.Description>Configure GitLab OAuth authentication.</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			<div class="space-y-2">
				<Label for="gitlab-client-id">Client ID</Label>
				<Input id="gitlab-client-id" bind:value={gitlab.client_id} placeholder="abc123..." />
			</div>

			<div class="space-y-2">
				<Label for="gitlab-secret">Client Secret</Label>
				<Input
					id="gitlab-secret"
					type="password"
					bind:value={gitlab.secret}
					placeholder="secret123..."
				/>
			</div>

			<div class="space-y-2">
				<Label for="gitlab-redirect">Redirect URL</Label>
				<Input
					id="gitlab-redirect"
					bind:value={gitlab.redirect_url}
					placeholder="https://example.com/auth/gitlab/callback"
				/>
			</div>

			<div class="space-y-2">
				<Label for="gitlab-base-url">Base URL</Label>
				<Input id="gitlab-base-url" bind:value={gitlab.base_url} placeholder="https://gitlab.com" />
				<p class="text-xs text-muted-foreground">
					For self-hosted GitLab instances. Leave empty for gitlab.com.
				</p>
			</div>
		</Card.Content>
		<Card.Footer class="flex justify-end">
			<Button onclick={() => handleSaveProvider('gitlab', gitlab)}>Save GitLab</Button>
		</Card.Footer>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<div class="flex items-center gap-2">
				<svg
					class="h-5 w-5"
					viewBox="0 0 24 24"
					fill="currentColor"
					xmlns="http://www.w3.org/2000/svg"
				>
					<path
						d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm0 2.162c.735 0 1.332.597 1.332 1.332v5.838l3.419-3.419a1.332 1.332 0 0 1 1.886 1.886l-3.419 3.419h5.838a1.332 1.332 0 1 1 0 2.664h-5.838l3.419 3.419a1.332 1.332 0 0 1-1.886 1.886l-3.419-3.419v5.838a1.332 1.332 0 1 1-2.664 0v-5.838l-3.419 3.419a1.332 1.332 0 0 1-1.886-1.886l3.419-3.419H2.944a1.332 1.332 0 1 1 0-2.664h5.838L5.363 6.799a1.332 1.332 0 0 1 1.886-1.886l3.419 3.419V2.494c0-.735.597-1.332 1.332-1.332z"
					/>
				</svg>
				<Card.Title>Bitbucket OAuth</Card.Title>
			</div>
			<Card.Description>Configure Bitbucket OAuth authentication.</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			<div class="space-y-2">
				<Label for="bitbucket-client-id">Client ID</Label>
				<Input id="bitbucket-client-id" bind:value={bitbucket.client_id} placeholder="abc123..." />
			</div>

			<div class="space-y-2">
				<Label for="bitbucket-secret">Client Secret</Label>
				<Input
					id="bitbucket-secret"
					type="password"
					bind:value={bitbucket.secret}
					placeholder="secret123..."
				/>
			</div>

			<div class="space-y-2">
				<Label for="bitbucket-redirect">Redirect URL</Label>
				<Input
					id="bitbucket-redirect"
					bind:value={bitbucket.redirect_url}
					placeholder="https://example.com/auth/bitbucket/callback"
				/>
			</div>
		</Card.Content>
		<Card.Footer class="flex justify-end">
			<Button onclick={() => handleSaveProvider('bitbucket', bitbucket)}>Save Bitbucket</Button>
		</Card.Footer>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<div class="flex items-center gap-2">
				<Chrome class="h-5 w-5" />
				<Card.Title>Google OAuth</Card.Title>
			</div>
			<Card.Description>Configure Google OAuth authentication.</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			<div class="space-y-2">
				<Label for="google-client-id">Client ID</Label>
				<Input
					id="google-client-id"
					bind:value={google.client_id}
					placeholder="123456789-abc.apps.googleusercontent.com"
				/>
			</div>

			<div class="space-y-2">
				<Label for="google-secret">Client Secret</Label>
				<Input
					id="google-secret"
					type="password"
					bind:value={google.secret}
					placeholder="GOCSPX-..."
				/>
			</div>

			<div class="space-y-2">
				<Label for="google-redirect">Redirect URL</Label>
				<Input
					id="google-redirect"
					bind:value={google.redirect_url}
					placeholder="https://example.com/auth/google/callback"
				/>
			</div>

			<div class="space-y-2">
				<Label for="google-tenant">Tenant (optional)</Label>
				<Input id="google-tenant" bind:value={google.tenant} placeholder="example.com" />
				<p class="text-xs text-muted-foreground">
					Hosted domain for Google Workspace. Leave empty for personal accounts.
				</p>
			</div>
		</Card.Content>
		<Card.Footer class="flex justify-end">
			<Button onclick={() => handleSaveProvider('google', google)}>Save Google</Button>
		</Card.Footer>
	</Card.Root>

	<div
		class="bg-info-subtle border border-info rounded-lg p-4"
	>
		<p class="text-sm font-medium text-info-foreground">ℹ️ Note</p>
		<p class="text-xs text-info-foreground mt-1">
			After configuring OAuth providers, the corresponding login buttons will be enabled on the
			login and signup pages. Make sure to save your changes and restart the application if needed.
		</p>
	</div>
</div>

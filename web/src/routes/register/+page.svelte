<script>
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Alert, AlertDescription } from '$lib/components/ui/alert';
	import { Eye, EyeOff, Mail, Lock, User, AlertCircle } from 'lucide-svelte';

	let name = $state('');
	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let showPassword = $state(false);
	let showConfirmPassword = $state(false);
	let isLoading = $state(false);
	let error = $state('');

	const isEmailValid = $derived(email.includes('@') && email.includes('.'));
	const isPasswordValid = $derived(password.length >= 6);
	const doPasswordsMatch = $derived(password === confirmPassword && confirmPassword !== '');
	const isFormValid = $derived(
		name.trim() !== '' && 
		isEmailValid && 
		isPasswordValid && 
		doPasswordsMatch
	);

	async function handleSubmit(event) {
		event.preventDefault();
		
		if (!isFormValid) {
			error = 'Please fill in all fields correctly';
			return;
		}

		isLoading = true;
		error = '';

		try {
			await new Promise(resolve => setTimeout(resolve, 1500));
			// Simulate successful registration
			goto('/login');
		} catch (err) {
			error = 'Registration failed. Please try again.';
		} finally {
			isLoading = false;
		}
	}

	$effect(() => {
		if (name || email || password || confirmPassword) {
			error = '';
		}
	});
</script>

<svelte:head>
	<title>Sign Up - Your App</title>
	<meta name="description" content="Create your account" />
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full space-y-8">
		<div class="text-center">
			<h1 class="text-3xl font-bold text-gray-900">Create your account</h1>
			<p class="mt-2 text-sm text-gray-600">
				Already have an account? 
				<a href="/login" class="font-medium text-primary hover:text-primary/80">
					Sign in
				</a>
			</p>
		</div>

		<Card>
			<CardHeader>
				<CardTitle>Sign up</CardTitle>
				<CardDescription>
					Enter your information to create your account
				</CardDescription>
			</CardHeader>
			<CardContent>
				<form onsubmit={handleSubmit} class="space-y-6">
					{#if error}
						<Alert variant="destructive">
							<AlertCircle class="h-4 w-4" />
							<AlertDescription>{error}</AlertDescription>
						</Alert>
					{/if}

					<div class="space-y-2">
						<Label for="name">Full name</Label>
						<div class="relative">
							<User class="absolute left-3 top-3 h-4 w-4 text-gray-400" />
							<Input
								id="name"
								type="text"
								bind:value={name}
								placeholder="Enter your full name"
								class="pl-10"
								required
								autocomplete="name"
							/>
						</div>
					</div>

					<div class="space-y-2">
						<Label for="email">Email address</Label>
						<div class="relative">
							<Mail class="absolute left-3 top-3 h-4 w-4 text-gray-400" />
							<Input
								id="email"
								type="email"
								bind:value={email}
								placeholder="Enter your email"
								class="pl-10"
								required
								autocomplete="email"
							/>
						</div>
						{#if email && !isEmailValid}
							<p class="text-sm text-red-600">Please enter a valid email address</p>
						{/if}
					</div>

					<div class="space-y-2">
						<Label for="password">Password</Label>
						<div class="relative">
							<Lock class="absolute left-3 top-3 h-4 w-4 text-gray-400" />
							<Input
								id="password"
								type={showPassword ? 'text' : 'password'}
								bind:value={password}
								placeholder="Create a password"
								class="pl-10 pr-10"
								required
								autocomplete="new-password"
							/>
							<button
								type="button"
								onclick={() => showPassword = !showPassword}
								class="absolute right-3 top-3 text-gray-400 hover:text-gray-600"
								aria-label={showPassword ? 'Hide password' : 'Show password'}
							>
								{#if showPassword}
									<EyeOff class="h-4 w-4" />
								{:else}
									<Eye class="h-4 w-4" />
								{/if}
							</button>
						</div>
						{#if password && !isPasswordValid}
							<p class="text-sm text-red-600">Password must be at least 6 characters</p>
						{/if}
					</div>

					<div class="space-y-2">
						<Label for="confirmPassword">Confirm password</Label>
						<div class="relative">
							<Lock class="absolute left-3 top-3 h-4 w-4 text-gray-400" />
							<Input
								id="confirmPassword"
								type={showConfirmPassword ? 'text' : 'password'}
								bind:value={confirmPassword}
								placeholder="Confirm your password"
								class="pl-10 pr-10"
								required
								autocomplete="new-password"
							/>
							<button
								type="button"
								onclick={() => showConfirmPassword = !showConfirmPassword}
								class="absolute right-3 top-3 text-gray-400 hover:text-gray-600"
								aria-label={showConfirmPassword ? 'Hide password' : 'Show password'}
							>
								{#if showConfirmPassword}
									<EyeOff class="h-4 w-4" />
								{:else}
									<Eye class="h-4 w-4" />
								{/if}
							</button>
						</div>
						{#if confirmPassword && !doPasswordsMatch}
							<p class="text-sm text-red-600">Passwords do not match</p>
						{/if}
					</div>

					<Button
						type="submit"
						class="w-full"
						disabled={!isFormValid || isLoading}
					>
						{#if isLoading}
							<div class="flex items-center">
								<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
								Creating account...
							</div>
						{:else}
							Create account
						{/if}
					</Button>
				</form>
			</CardContent>
		</Card>
	</div>
</div>

<script>
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input/index';
	import { Label } from '$lib/components/ui/label';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card/index';
	import { Alert, AlertDescription } from '$lib/components/ui/alert';
	import { Eye, EyeOff, Mail, Lock, CircleAlert } from 'lucide-svelte';

	// Reactive state using Svelte 5 runes
	let email = $state('');
	let password = $state('');
	let showPassword = $state(false);
	let isLoading = $state(false);
	let error = $state('');

	// Form validation using derived state
	const isEmailValid = $derived(email.includes('@') && email.includes('.'));
	const isPasswordValid = $derived(password.length >= 6);
	const isFormValid = $derived(isEmailValid && isPasswordValid && email.trim() !== '');

	// Handle form submission
	async function handleSubmit(event) {
		event.preventDefault();

		if (!isFormValid) {
			error = 'Please fill in all fields correctly';
			return;
		}

		isLoading = true;
		error = '';

		try {
			// Simulate API call
			await new Promise((resolve) => setTimeout(resolve, 1500));

			// Mock authentication logic
			if (email === 'demo@example.com' && password === 'password') {
				// Successful login - redirect to dashboard
				goto('/dashboard');
			} else {
				error = 'Invalid email or password';
			}
		} catch (err) {
			error = 'An error occurred. Please try again.';
		} finally {
			isLoading = false;
		}
	}

	// Toggle password visibility
	function togglePasswordVisibility() {
		showPassword = !showPassword;
	}

	// Clear error when user starts typing
	$effect(() => {
		if (email || password) {
			error = '';
		}
	});
</script>

<svelte:head>
	<title>Login - Your App</title>
	<meta name="description" content="Sign in to your account" />
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full space-y-8">
		<div class="text-center">
			<h1 class="text-3xl font-bold text-gray-900">Welcome back</h1>
			<p class="mt-2 text-sm text-gray-600">
				Don't have an account?
				<a href="/register" class="font-medium text-primary hover:text-primary/80"> Sign up </a>
			</p>
		</div>

		<Card>
			<CardHeader>
				<CardTitle>Sign in to your account</CardTitle>
				<CardDescription>Enter your email and password to access your account</CardDescription>
			</CardHeader>
			<CardContent>
				<form onsubmit={handleSubmit} class="space-y-6">
					{#if error}
						<Alert variant="destructive">
							<CircleAlert class="h-4 w-4" />
							<AlertDescription>{error}</AlertDescription>
						</Alert>
					{/if}

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
								placeholder="Enter your password"
								class="pl-10 pr-10"
								required
								autocomplete="current-password"
							/>
							<button
								type="button"
								onclick={togglePasswordVisibility}
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

					<div class="flex items-center justify-between">
						<div class="flex items-center">
							<input
								id="remember-me"
								name="remember-me"
								type="checkbox"
								class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded"
							/>
							<Label for="remember-me" class="ml-2 text-sm">Remember me</Label>
						</div>

						<div class="text-sm">
							<a href="/forgot-password" class="font-medium text-primary hover:text-primary/80">
								Forgot your password?
							</a>
						</div>
					</div>

					<Button type="submit" class="w-full" disabled={!isFormValid || isLoading}>
						{#if isLoading}
							<div class="flex items-center">
								<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
								Signing in...
							</div>
						{:else}
							Sign in
						{/if}
					</Button>
				</form>

				<div class="mt-6">
					<div class="relative">
						<div class="absolute inset-0 flex items-center">
							<div class="w-full border-t border-gray-300"></div>
						</div>
						<div class="relative flex justify-center text-sm">
							<span class="px-2 bg-white text-gray-500">Or continue with</span>
						</div>
					</div>

					<div class="mt-6 grid grid-cols-2 gap-3">
						<Button variant="outline" class="w-full">
							<svg class="w-4 h-4 mr-2" viewBox="0 0 24 24">
								<path
									fill="currentColor"
									d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
								/>
								<path
									fill="currentColor"
									d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
								/>
								<path
									fill="currentColor"
									d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
								/>
								<path
									fill="currentColor"
									d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
								/>
							</svg>
							Google
						</Button>
						<Button variant="outline" class="w-full">
							<svg class="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 24 24">
								<path
									d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"
								/>
							</svg>
							Facebook
						</Button>
					</div>
				</div>

				<div class="mt-6 text-center text-sm text-gray-600">
					<p>Demo credentials:</p>
					<p class="font-mono text-xs mt-1">
						Email: demo@example.com<br />
						Password: password
					</p>
				</div>
			</CardContent>
		</Card>
	</div>
</div>

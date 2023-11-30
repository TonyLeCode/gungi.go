<script lang="ts">
	import { superForm } from 'sveltekit-superforms/client';
	import type { PageData } from './$types.js';
	import { default as confetti } from 'canvas-confetti';
	import { browser } from '$app/environment';
	import { onDestroy } from 'svelte';

	export let data: PageData;
	const { form, errors, constraints, enhance, message } = superForm(data.form, { taintedMessage: null });
	const unsub = message.subscribe((msg) => {
		if (msg === 'success' && browser) {
			const count = 110;
			const defaults = {
				origin: { y: 0.7 },
			};
			function fire(particleRatio: number, opts: object) {
				confetti({
					...defaults,
					...opts,
					particleCount: Math.floor(count * particleRatio),
				});
			}
			fire(0.2, {
				spread: 26,
				startVelocity: 55,
			});
			fire(0.2, {
				spread: 60,
			});
			fire(0.2, {
				spread: 150,
				decay: 0.91,
				scalar: 0.8,
			});
			fire(0.2, {
				spread: 180,
				startVelocity: 25,
				decay: 0.92,
				scalar: 1.2,
			});
			fire(0.2, {
				spread: 240,
				startVelocity: 45,
			});
		}
	});
	onDestroy(unsub);
</script>

<main>
	{#if $message === 'success'}
		<h2 class="registered">🎉 Registration Successful! 🎉</h2>
		<h3>Please check your email for confirmation</h3>
	{/if}
	<form class="register" method="POST" use:enhance>
		<h2>Register</h2>
		<fieldset>
			<div class="input-group">
				<label for="email">Email:</label>
				<input
					id="email"
					name="email"
					type="email"
					bind:value={$form.email}
					aria-invalid={$errors.email ? 'true' : undefined}
					{...$constraints.email}
				/>
				{#if $errors.email}<div class="invalid">{$errors.email}</div>{/if}
			</div>
			<div class="input-group">
				<label for="password">Password:</label>
				<input
					id="password"
					name="password"
					type="password"
					bind:value={$form.password}
					aria-invalid={$errors.password ? 'true' : undefined}
					{...$constraints.password}
				/>
				{#if $errors.password}<div class="invalid">{$errors.password}</div>{/if}
			</div>
		</fieldset>
		<button class="button-primary" type="submit">Register</button>
		<a href="/login" class="button-ghost">Login</a>
	</form>
</main>

<style lang="scss">
	main {
		min-height: calc(100vh - 3rem);
		display: flex;
		justify-content: center;
		align-items: center;
		flex-direction: column;
	}
	input {
		display: block;
		width: 100%;
		border-radius: 4px;
		padding: 0.25rem 0.75rem;
		background-color: rgb(var(--bg-3));
		&:focus {
			outline: 2px solid rgb(var(--primary));
		}
	}
	.input-group {
		position: relative;
		margin-bottom: 1.25rem;
	}
	.invalid {
		font-weight: 300;
		margin-top: 4px;
	}
	.registered {
		margin-top: 4rem;
	}
	h2 {
		text-align: center;
		font-size: 1.25rem;
		font-weight: 600;
	}
	fieldset {
		padding: 1rem 0;
	}
	.button-ghost {
		text-align: center;
	}
	.register {
		gap: 1rem;
		display: flex;
		flex-direction: column;
		max-width: 20rem;
		margin: auto;
		padding: 6rem 3rem;
		width: 100%;
		background-color: rgb(var(--bg-2));
		box-shadow:
			0px 5px 15px rgba(0, 0, 0, 0.07),
			0px 2px 5px rgba(0, 0, 0, 0.05);
		box-sizing: content-box;
		min-height: 25rem;
		//TODO handle margin-bottom/layout with form messages
		margin-bottom: 10rem;
	}
	a {
		margin-bottom: 1rem;
	}
</style>

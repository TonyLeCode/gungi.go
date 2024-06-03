<script lang="ts">
	import { superForm } from 'sveltekit-superforms/client';
	import { page } from '$app/stores';
	import { zodClient } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';

	const schema = z.object({
		username: z.string().min(3).max(28),
	});

	let { data } = $props();
	let onboard = $page.url.searchParams.get('onboard');
	let onboardBool = $derived(onboard === 'true');
	let username = $derived(data.session?.user.user_metadata.username);
	const { form, errors, constraints, enhance } = superForm(data.form, {
		taintedMessage: null,
		validators: zodClient(schema),
	});
</script>

<main>
	{#if onboardBool}
		<h2 class="onboard">You have been given a randomly generated username, you can change it here now or later.</h2>
	{/if}
	<form class="register" method="POST" use:enhance>
		<h2>Change Username</h2>
		<fieldset>
			<div class="input-group">
				<div class="current-username-label">Current Username:</div>
				<div class="current-username">{username}</div>
			</div>
			<div class="input-group">
				<label for="username">Change Username To:</label>
				<input
					id="username"
					name="username"
					type="input"
					bind:value={$form.username}
					aria-invalid={$errors.username ? 'true' : undefined}
					{...$constraints.username}
				/>
				{#if $errors.username}<div class="invalid">{$errors.username}</div>{/if}
			</div>
		</fieldset>
		<button class="button-primary" type="submit">Change Username</button>
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
	.current-username {
		background-color: rgb(var(--bg-3));
		color: rgba(var(--font), 0.6);
		border-radius: 4px;
		padding: 0.25rem 0.75rem;
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
	.onboard {
		margin-top: 2rem;
		max-width: 50ch;
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
		margin-bottom: 10rem;
	}
</style>

<script lang="ts">
	import { superForm } from 'sveltekit-superforms/client';
	import { zodClient } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';

	const schema = z.object({
		username: z
			.string()
			.min(3, { message: 'Username must contain at least 3 character(s)' })
			.max(28, { message: 'Username must contain less than 28 character(s)' })
			.trim()
			.regex(/^[a-zA-Z0-9 _-]+$/, {
				message: 'Username may only contain letters, numbers, spaces, hyphens, and underscores',
			}),
	});

	let { data } = $props();
	const { form, errors, constraints, enhance, message } = superForm(data.form, {
		taintedMessage: null,
		validators: zodClient(schema),
	});
</script>

<main>
	<h2 class="onboard">Welcome, please create a username. It will be displayed to others.</h2>
	<form class="form-container" method="POST" use:enhance>
		<fieldset>
			<div class="input-group">
				<label for="username">Username:</label>
				<input
					id="username"
					name="username"
					type="input"
					bind:value={$form.username}
					aria-invalid={$errors.username ? 'true' : undefined}
					{...$constraints.username}
				/>
				{#if $errors.username}
					{#each $errors.username as error}
						<div class="invalid">{error}</div>
					{/each}
				{/if}
			</div>
		</fieldset>
		<button class="button-primary" type="submit">Submit Username</button>
	</form>
</main>

<style lang="scss">
	main {
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
		margin-bottom: 1rem;
	}
	.invalid {
		font-weight: 300;
		margin-top: 4px;
	}
	.onboard {
		margin-top: 2rem;
		max-width: 45ch;
		margin-bottom: 2rem;
	}
	h2 {
		text-align: center;
		font-size: 1.25rem;
		font-weight: 600;

	}
	.button-ghost {
		text-align: center;
	}
	.form-container {
		gap: 1rem;
		display: flex;
		flex-direction: column;
		max-width: 26rem;
		margin: auto;
		padding: 3rem 3rem;
		width: 100%;
		background-color: rgb(var(--bg-2));
		box-shadow:
			0px 5px 15px rgba(0, 0, 0, 0.07),
			0px 2px 5px rgba(0, 0, 0, 0.05);
	}
</style>

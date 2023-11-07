<script lang="ts">
	import { superForm } from 'sveltekit-superforms/client';
	import type { PageData } from './$types.js';

	export let data: PageData;
	const { form, errors, constraints, enhance, message } = superForm(data.form, { taintedMessage: null });
</script>

<main>
	<form class="login" method="POST" use:enhance>
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
			{#if $message}<div class="invalid">{$message}</div>{/if}
		</fieldset>
		<button class="button-primary" type="submit">Login</button>
		<a href="/register" class="button-ghost">Register</a>
	</form>
</main>

<style lang="scss">
	main {
		min-height: calc(100vh - 3rem);
		display: flex;
		justify-content: center;
		align-items: center;
	}
	input {
		/* border: 1.5px solid rgba(var(--primary), 0.25); */
		display: block;
		width: 100%;
		border-radius: 4px;
		padding: 0.25rem 0.75rem;
		/* background-color: rgb(var(--bg)); */
		background-color: rgb(var(--bg-3));
		/* background-color: rgb(236, 236, 236); */
		/* box-shadow: inset 0px 0px 15px 2px rgba(0, 0, 0, 0.045); */
		&:focus {
			/* outline: 5px solid rgb(var(--primary)); */
			outline: 2px solid rgb(var(--primary));
			// outline-offset: 2px;
		}
	}
	.input-group {
		position: relative;
		margin-bottom: 1.25rem;
	}
	.invalid {
		// position: absolute;
		font-weight: 300;
		// color: rgba(var(--font), .5);
		margin-top: 4px;
		// bottom: -1.75rem;
	}
	fieldset {
		padding: 1rem 0;
	}
	.button-ghost {
		text-align: center;
	}
	.login {
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
	}
	a {
		margin-bottom: 1rem;
	}
</style>

<script lang="ts">
	import { superForm } from 'sveltekit-superforms/client';

	let { data } = $props();
	const { form, errors, constraints, enhance, message } = superForm(data.form, { taintedMessage: null });
</script>

<svelte:head>
	<title>Login | White Monarch Server</title>
</svelte:head>

<main>
	<form class="login" method="POST" use:enhance>
		<h2>Login</h2>
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
		display: block;
		width: 100%;
		border-radius: 4px;
		padding: 0.25rem 0.75rem;
		background-color: rgb(var(--bg-3));
		&:focus {
			outline: 2px solid rgb(var(--primary));
		}
	}
	h2 {
		text-align: center;
		font-size: 1.25rem;
		font-weight: 600;
	}
	.input-group {
		position: relative;
		margin-bottom: 1.25rem;
	}
	.invalid {
		font-weight: 300;
		margin-top: 4px;
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
		margin-bottom: 10rem;
	}
	a {
		margin-bottom: 1rem;
	}
</style>

<script lang="ts">
	import Modal from './Modal.svelte';

	export let showModal: boolean;
</script>

<Modal bind:showModal>
	<form class="login" method="POST" action="/login">
		<fieldset>
			<label for="email">Email:</label>
			<input name="email" type="email" />
			<label for="password">Password:</label>
			<input name="password" type="password" />
		</fieldset>
		<button class="button-primary" type="submit">Login</button>
		<a href="/register" class="button-ghost">Register</a>
	</form>
</Modal>

<style>
	input {
		border: 1.5px solid rgba(var(--primary), 0.25);
		padding: 0.25rem 0.75rem;
		background-color: rgb(var(--bg));
	}

	fieldset {
		margin: 1rem 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.button-ghost {
		text-align: center;
	}

	.login {
		gap: 1rem;
		display: flex;
		flex-direction: column;
		min-width: 20rem;
	}

	a {
		margin-bottom: 1rem;
	}
</style>

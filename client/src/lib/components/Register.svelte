<script lang="ts">
	import { supabase } from '$lib/supabaseClient';
	import Modal from './Modal.svelte';

	export let showModal: boolean;
	let email: string;
	let password: string;
	let username: string;

	const handleSignUp = async () => {
		const { data, error } = await supabase.auth.signUp({
			email: email,
			password: password,
			options: {
				data: { username: username },
			},
		});

		if (error) {
			console.log(error);
		} else {
			console.log(data);
		}
	};
</script>

<svelte:head>
	<title>Gungi.go</title>
</svelte:head>

<Modal bind:showModal>
	<form class="login" action="/register" on:submit|preventDefault={handleSignUp}>
		<fieldset>
			<label for="username">Username:</label>
			<input bind:value={username} type="text" />
			<label for="email">Email:</label>
			<input bind:value={email} type="email" />
			<label for="password">Password:</label>
			<input bind:value={password} type="password" />
		</fieldset>
		<button class="button-primary" type="submit">Register</button>
		<a href="/login" class="button-ghost">Login</a>
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

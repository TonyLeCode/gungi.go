<script lang="ts">
	import Navbar from '$lib/components/Navbar.svelte';
	import { supabase } from '$lib/supabaseClient';
	// TODO verify unique username
	//TODO verification

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

<form class="login" on:submit|preventDefault={handleSignUp}>
	<fieldset>
		<label for="username">Username:</label>
		<input id="username" bind:value={username} type="text" />
		<label for="email">Email:</label>
		<input id="email" bind:value={email} type="email" />
		<label for="password">Password:</label>
		<input id="password" bind:value={password} type="password" />
	</fieldset>
	<button class="register-button" type="submit">Register</button>
	<button class="register-button">Login</button>
</form>

<style>
	input {
		border: 1px solid black;
		padding: 0.25rem 0.75rem;
	}

	fieldset {
		margin: 1rem 0;
		display: flex;
		flex-direction: column;
	}

	.login {
		display: flex;
		flex-direction: column;
		max-width: 20rem;
		margin: auto;
		border: 1px solid black;
		padding: 2rem;
	}

	.register-button {
		border: 1px solid black;
	}
</style>

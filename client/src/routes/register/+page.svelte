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

<main>
	<form class="login" on:submit|preventDefault={handleSignUp}>
		<fieldset>
			<label for="username">Username:</label>
			<input id="username" bind:value={username} type="text" />
			<label for="email">Email:</label>
			<input id="email" bind:value={email} type="email" />
			<label for="password">Password:</label>
			<input id="password" bind:value={password} type="password" />
		</fieldset>
		<button class="button-primary" type="submit">Register</button>
		<a href='/login' class="button-ghost">Login</a>
	</form>
</main>

<style>
	main{
		height: calc(100vh - 3rem);
		display:flex;
		justify-content: center;
	}
	input {
		border: 1.5px solid rgba(var(--primary), 0.25);
		padding: 0.25rem 0.75rem;
		background-color: rgb(var(--white));
	}
	form{
	}
	
	fieldset {
		margin: 1rem 0;
		display: flex;
		flex-direction: column;
		gap: .5rem;
	}

	.button-ghost{
		text-align: center;
	}
	
	.login {
		gap: 1rem;
		display: flex;
		flex-direction: column;
		max-width: 20rem;
		margin: auto;
		padding: 8rem 6rem;
		width: 100%;
		margin-bottom: 12rem;
		background-color: rgb(var(--white-2));
		box-shadow: 0px 5px 15px rgba(0, 0, 0, 0.07);
		box-sizing: content-box;
	}

	a{
		margin-bottom: 1rem;
	}
</style>

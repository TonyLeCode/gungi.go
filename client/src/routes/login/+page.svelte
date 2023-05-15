<script lang="ts">
	import { supabase } from '$lib/supabaseClient';
	import { fail, redirect } from '@sveltejs/kit';

	// const handleSignIn = async () => {
	//   const { data, error } = await supabase.auth.signInWithPassword({
	//     email: email,
	//     password: password,
	//   })

	//   if (error) {
	//     console.log(error);
	//     return fail(500, {
	//       error: 'Server error. Try again later.',
	//       values: {
	//         email
	//       }
	//     });
	//   } else {
	//     console.log("logged in", data);
	//   }
	//   throw redirect(303, '/');
	// }
</script>

<svelte:head>
	<title>Gungi.go</title>
</svelte:head>

<main>
	<form class="login" method="POST">
		<fieldset>
			<label for="email">Email:</label>
			<input name="email" id="email" type="email" />
			<label for="password">Password:</label>
			<input name="password" id="password" type="password" />
		</fieldset>
		<button class="button-primary">Register</button>
		<button class="button-primary" type="submit">Login</button>
	</form>
</main>

<style>
	main {
		height: calc(100vh - 3rem);
		display: flex;
		justify-content: center;
	}
	input {
		border: 1.5px solid rgba(var(--primary), 0.25);
		padding: 0.25rem 0.75rem;
		background-color: rgb(var(--white));
	}
	form {
	}

	fieldset {
		margin: 1rem 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
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

	button:nth-last-of-type(1) {
		margin-bottom: 1rem;
	}
</style>

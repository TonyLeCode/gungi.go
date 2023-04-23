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

<form class="login" method="POST">
	<fieldset>
		<label
			>Email:
			<input name="email" class="email" type="email" />
		</label>
		<label
			>Password:
			<input name="password" class="password" type="password" />
		</label>
	</fieldset>
	<button class="login-button">Register</button>
	<button class="login-button" type="submit">Login</button>
</form>

<style>
	input {
		border: 1px solid black;
		padding: 0.25rem 0.75rem;
	}

	.login {
		display: flex;
		flex-direction: column;
		max-width: 20rem;
		margin: auto;
		border: 1px solid black;
		padding: 2rem;
	}

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

	.login-button {
		border: 1px solid black;
	}
</style>

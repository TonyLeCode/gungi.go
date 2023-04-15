<script lang="ts">
	import Navbar from "$lib/components/Navbar.svelte";
  import { supabase } from '$lib/supabaseClient';

  let email: string;
  let password: string;

  const handleSignIn = async () => {
    const { data, error } = await supabase.auth.signInWithPassword({
      email: email,
      password: password,
    })

    if (error) {
      console.log(error);
    } else {
      console.log(data);
    }
  }
</script>

<svelte:head>
	<title>Gungi.go</title>
</svelte:head>

<Navbar />
<form class="login" on:submit|preventDefault="{handleSignIn}">
	<fieldset>
    <label for="email">Email:</label>
    <input id="email" bind:value="{email}" class="email" type="email" />
    <label for="password">Password:</label>
    <input id="password" bind:value="{password}" class="password" type="password" />
  </fieldset>
  <button class='login-button'>Register</button>
  <button class='login-button' type='submit'>Login</button>
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

  input{
    border:1px solid black;
    padding: .25rem .75rem;
  }

  fieldset{
    margin: 1rem 0;
    display: flex;
    flex-direction: column;
  }

  .login{
    display: flex;
    flex-direction: column;
    max-width: 20rem;
    margin:auto;
    border: 1px solid black;
    padding: 2rem;
  }

  .login-button{
    border: 1px solid black;
  }
</style>

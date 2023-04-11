<script>
	import Navbar from '../../components/Navbar.svelte';
</script>

<svelte:head>
	<title>Gungi.go</title>
</svelte:head>

<Navbar />
<div class="login">
	<fieldset>
    <label for="email">Email:</label>
    <input id="email" class="email" type="email" />
    <label for="password">Password:</label>
    <input id="password" class="password" type="password" />
  </fieldset>
  <button class='login-button'>Register</button>
  <button class='login-button'>Login</button>
</div>

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

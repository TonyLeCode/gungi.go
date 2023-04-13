<script>
	import Navbar from "$lib/components/Navbar.svelte";

</script>
<svelte:head>
  <title>Gungi.go</title>
</svelte:head>

<Navbar />
<form class='login'>
  <fieldset>
    <label for="username">Username:</label>
    <input id='username' type="text">
    <label for="email">Email:</label>
    <input id='email' type="email">
    <label for="password">Password:</label>
    <input id='password' type="password">
    <label for="password">Confirm Password:</label>
    <input id='confirmpassword' type="password">
  </fieldset>
  <button class='register-button'>Register</button>
  <button class='register-button'>Login</button>
</form>

<style>
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

  .register-button{
    border: 1px solid black;
  }
</style>
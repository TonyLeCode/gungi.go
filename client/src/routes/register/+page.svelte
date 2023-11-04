<script lang='ts'>
	export let form;

	if (form){
		console.log("form: ", form)
		//TODO, make notification
	}
</script>

<main>
  <form class="register" method="POST">
		<fieldset>
			<label for="username">Username:</label>
			<input id="username" name="username" type="input" required />
			<label for="email">Email:</label>
			<input id="email" name="email" type="email" required />
			<label for="password">Password:</label>
			<input id="password" name="password" type="password" required />
		</fieldset>
		<button class="button-primary" type="submit">Register</button>
		<a href="/login" class="button-ghost">Login</a>
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
		/* border: 1.5px solid rgba(var(--primary), 0.25); */
    border-radius: 4px;
		padding: 0.25rem 0.75rem;
		/* background-color: rgb(var(--bg)); */
		background-color: rgb(var(--bg-3));
		/* background-color: rgb(236, 236, 236); */
    /* box-shadow: inset 0px 0px 15px 2px rgba(0, 0, 0, 0.045); */
    &:focus {
		/* outline: 5px solid rgb(var(--primary)); */
		outline: 2px solid rgb(var(--primary));
		// outline-offset: 2px;
	}
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
	.register {
		gap: 1rem;
		display: flex;
		flex-direction: column;
		max-width: 20rem;
		margin: auto;
		padding: 6rem 3rem;
		width: 100%;
		background-color: rgb(var(--bg-2));
		box-shadow: 0px 5px 15px rgba(0, 0, 0, 0.07), 0px 2px 5px rgba(0, 0, 0, 0.05);
		box-sizing: content-box;
    min-height: 25rem;
	}
	a {
		margin-bottom: 1rem;
	}
</style>

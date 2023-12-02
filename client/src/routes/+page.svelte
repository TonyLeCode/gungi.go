<script>
	export let data;
</script>

<svelte:head>
	<title>White Monarch Server</title>
</svelte:head>

<main>
	<h1>White Monarch Server</h1>
	<p>Welcome! This is a free web app to play Gungi online. You can spectate or play live and correspondence games.</p>
	{#if data.session}
		<a class="button-primary" href="/play">Play</a>
	{:else}
		<div>
			<a class="button-primary" href="/login">Login</a>
			<a class="button-primary" href="/register">Register</a>
		</div>
	{/if}
</main>

<style>
	main {
		max-width: 75ch;
		margin: auto;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		min-height: 50vh;
		padding: 1rem;
	}

	h1 {
		font-size: 2.5rem;
		text-align: center;
		font-weight: 600;
		/* letter-spacing: .1rem; */
	}
	p {
		margin-top: 0.5rem;
		width: 100%;
	}
	div {
		margin-top: 0.5rem;
		display: flex;
		gap: 1.5rem;
	}
	a {
		/* font-size: 1.2rem; */
		padding: 0.5rem 2rem;
	}

	@media only screen and (min-width: 600px){
		h1 {
			font-size: 3rem;
		}
	}
</style>

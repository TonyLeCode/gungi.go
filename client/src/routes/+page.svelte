<script>
	export let data;
</script>

<svelte:head>
	<title>White Monarch Server</title>
</svelte:head>

<main>
	<h1>White Monarch Server</h1>
	<img src="/home-splash.png" alt="home-splash">
	<!-- <p>Welcome! This is a free web app to play Gungi online. You can spectate or play live and correspondence games.</p> -->
	<p>"You should enjoy the little detours to the fullest. Because that's where you'll find things more important than what you want."</p>
	<p class="author">- Ging Freecss</p>
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
		/* max-width: 75ch; */
		max-width: 65rem;
		margin: auto;
		display: flex;
		flex-direction: column;
		align-items: center;
		min-height: 50vh;
		padding: clamp(0.75rem, 2vw, 2rem);
	}

	h1 {
		font-size: clamp(1.5rem, 5vw, 2.5rem);
		text-align: center;
		font-weight: 600;
		margin-top: 1rem;
		/* letter-spacing: .1rem; */
		@media (min-width: 767px) {
			font-size: 2.5rem;
		}
	}
	p {
		margin-top: 1rem;
		padding: 0 1.5rem;
		width: 100%;
		max-width: 70ch;
		font-size: 1.25rem;
		font-size: clamp(0.875rem, 2vw, 1.25rem);
	}
	.author {
		text-align: right;
		margin-top: 0;
		margin-bottom: 1rem;
		font-weight: 600;
	}
	div {
		margin-top: 0.5rem;
		display: flex;
		gap: 1.5rem;
	}
	a {
		/* font-size: 1.2rem; */
		font-size: clamp(0.875rem, 2vw, 1.25rem);
		padding: 0.375rem 1.5rem;
		@media (min-width: 767px) {
			padding: 0.5rem 2rem;
		}
	}

	img{
		margin: clamp(1rem, 2vw, 2rem);
	}
</style>

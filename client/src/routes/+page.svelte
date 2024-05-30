<script>
	export let data;
</script>

<svelte:head>
	<title>White Monarch Server</title>
	<meta
		name="description"
		content="White Monarch Server is a server where you can play Gungi online against other players. Gungi is a fictional strategy game in Hunter x Hunter, made by Yoshihiro Togashi."
	/>
</svelte:head>

<main>
	<section class="intro">
		<h1>White Monarch Server</h1>
		<img src="/home-splash.png" alt="home-splash" />
		<!-- <p>Welcome! This is a free web app to play Gungi online. You can spectate or play live and correspondence games.</p> -->
		<p>
			"You should enjoy the little detours to the fullest. Because that's where you'll find things more important than
			what you want."
		</p>
		<p class="author">- Ging Freecss</p>
		{#if data.session}
			<a class="button-primary" href="/play">Play</a>
		{:else}
			<div>
				<a class="button-primary" href="/login">Login</a>
				<a class="button-primary" href="/register">Register</a>
			</div>
		{/if}
	</section>
	<section class="description">
		<h2>Gungi</h2>
		<p>Gungi is a fictional strategy game in Hunter x Hunter, made by Yoshihiro Togashi. It is similar to Chess or Shogi except players can create their own formations at the start of the game, and pieces may be stacked onto each other.</p>
		<p>Since the game is unofficial, there are various rulesets, currently we only have the revised community ruleset.</p>
		<p>I've created this site to provide a platform for players to play correspondence games against each other. I hope you enjoy your stay.</p>
	</section>
</main>

<style>
	main {
		/* max-width: 75ch; */
		max-width: 65rem;
		margin: auto;
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: clamp(0.75rem, 2vw, 2rem);
	}

	.intro {
		margin: auto;
		display: flex;
		flex-direction: column;
		align-items: center;
		margin-bottom: clamp(4rem, 15vw, 14rem);
	}

	.description {
		margin-bottom: 20rem;
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
	h2 {
		font-size: clamp(1.25rem, 5vw, 2rem);
		text-align: center;
		font-weight: 600;
		margin-top: 1rem;
		/* letter-spacing: .1rem; */
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

	img {
		margin: clamp(1rem, 2vw, 2rem);
	}
</style>

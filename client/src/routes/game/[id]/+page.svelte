<script lang="ts">
	import Board from '$lib/components/Board.svelte';

	export let data;
	// console.log(data?.params.id);
	console.log(data)
</script>

<svelte:head>
	<title>Gungi.go</title>
</svelte:head>

<main>
	<section>
		<h3>Player1</h3>
		<Board />
		<h3>Player2</h3>
	</section>
	<aside class="side-menu">
		<div class="hands">
			<div>
				hello
			</div>
			<div>
				hello
			</div>
		</div>
	</aside>
</main>

<style>
	main {
		display: flex;
		max-width: 90rem;
		margin: 0 auto;
	}

	.side-menu {
		background-color: gray;
		margin-left: auto;
		max-width: 35rem;
		width: 100%;
	}

	.hands {
		display: grid;
		gap: 1rem;
	}
	
	.hands div {
		background-color:yellow;
	}
	/* section {
		max-width: 70rem;
		margin: 0 auto;
		margin-top: 2rem;
		padding: 0 2rem;
		text-align: center;
	}

	.gameList {
		display: grid;
		gap: 1rem;
		grid-template-columns: repeat(auto-fit, 20rem);
		padding: 1rem;
		justify-content: center;
	} */
</style>

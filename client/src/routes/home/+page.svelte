<script lang="ts">
	import BoardSimple from '$lib/components/BoardSimple.svelte';

	export let data;
	$: console.log(data.data);
</script>

<svelte:head>
	<title>Gungi.go</title>
</svelte:head>

<main>
	<section>
		<h2>Current Games</h2>
		<ul class="gameList">
			{#each data.data as game}
				<li><a href={`/game/${game.id}`}><BoardSimple gameData={game} /></a></li>
			{/each}
			<BoardSimple gameData={{current_state:'3,k,5/3,psc,m,4/5,pyt,3/9/9/9/5,PT,P,2/6,F,PS,1/7,M,1 6446122122210/6446212121210 w'}} />
			<BoardSimple gameData={{current_state:'3,k,5/3,psc,m,4/5,pyt,3/9/9/9/5,PT,P,2/6,F,PS,1/7,M,1 6446122122210/6446212121210 w'}} />
			<BoardSimple gameData={{current_state:'3,k,5/3,psc,m,4/5,pyt,3/9/9/9/5,PT,P,2/6,F,PS,1/7,M,1 6446122122210/6446212121210 w'}} />
			<BoardSimple gameData={{current_state:'3,k,5/3,psc,m,4/5,pyt,3/9/9/9/5,PT,P,2/6,F,PS,1/7,M,1 6446122122210/6446212121210 w'}} />
		</ul>
	</section>
	<section>
		<h2>Game History</h2>
	</section>
</main>

<style>
	section {
		max-width: 70rem;
		margin: 0 auto;
		margin-top: 2rem;
		padding: 0 2rem;
		text-align: center;
	}

	h2 {
		font-size: 1.25rem;
	}

	li:hover {
		outline: rgba(158, 250, 193, 0.753) solid .5rem;
		/* box-sizing: content-box; */
	}

	.gameList {
		display: grid;
		gap: 1rem;
		grid-template-columns: repeat(auto-fit, 20rem);
		padding: 1rem;
		justify-content: center;
	}
</style>

<script lang="ts">
	import BoardSimple from '$lib/components/BoardSimple.svelte';

	interface Game {
		completed: boolean;
		current_state: string;
		date_started: Date;
		fen: {
			String: string;
			Valid: boolean;
		}
		id: string;
		username1: string;
		username2: string;
	}

	function isTurn(game: Game){
		const fields = game.current_state.split(' ')
		const turnColor = fields[2]
		const turnPlayer = () => {
			return turnColor === 'w' ? game.username1 : game.username2
		}
		return turnPlayer() === data.session?.user.user_metadata.username
		// if(data.session){
		// 	return data.session.user.user_metadata.username === user
		// } else {
		// 	return false
		// }
	}

	export let data;
	$: console.log(data);
	// $: console.log(data.data);
</script>

<svelte:head>
	<title>Gungi.go</title>
</svelte:head>

<main>
	<section>
		<h2>Current Games</h2>
		<ul class="gameList">
			{#each data.data as game}
				<li class:your-turn={isTurn(game)}><a href={`/game/${game.id}`}><BoardSimple gameData={game} /></a></li>
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

	li {
		border-radius: 4px;
	}
	
	li:hover {
		outline: rgb(240, 80, 17) solid .6rem;
		/* box-sizing: content-box; */
	}

	.your-turn {
		outline: rgba(255, 136, 81, 0.829) solid .5rem;
	}

	.gameList {
		display: grid;
		gap: 2rem;
		grid-template-columns: repeat(auto-fit, 20rem);
		padding: 1rem;
		justify-content: center;
	}
</style>

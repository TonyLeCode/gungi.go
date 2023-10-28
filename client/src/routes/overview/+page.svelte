<script lang="ts">
	import BoardSimple from '$lib/components/BoardSimple.svelte';
	import type {Game} from './+page.server'

	export let data;
	const username = data.session?.user.user_metadata.username ?? ''
	$: sortedGames = [...data.data].sort((a, b) => {
		if (turnPlayer(a) === turnPlayer(b)){
			return 0
		}
		if (isUserTurn(a)){
			return -1
		} else if (isUserTurn(b)){
			return 1
		}
		return 0
	})

	function isUser(playername: string){
		return username === playername
	}

	function getUserColor(player1: string, player2: string){
		if (username === player1) {
			return 'w'
		} else if (username === player2){
			return 'b'
		}
		return 'spectator'
	}

	function turnPlayer(game: Game){
		const fields = game.current_state.split(' ');
		const turnColor = fields[2];
		return turnColor === 'w' ? game.username1 : game.username2;
	}

	function isUserTurn(game: Game) {
		return turnPlayer(game) === username
		// if(data.session){
		// 	return data.session.user.user_metadata.username === user
		// } else {
		// 	return false
		// }
	}

	
	$: console.log(data.data);
</script>

<svelte:head>
	<title>Gungi.go</title>
</svelte:head>

<main>
	<section>
		<h2>Current Games</h2>
		<ul class="gameList">
			{#each sortedGames as game}
				<li class:your-turn={isUserTurn(game)}>
					<div class="name name-1" class:is-user={isUser(game.username1)}>{game.username1}</div>
					<a href={`/game/${game.id}`}><BoardSimple gameData={game} userColor={getUserColor(game.username1, game.username2)} /></a>
					<div class="name name-2" class:is-user={isUser(game.username2)}>{game.username2}</div>
				</li>
			{/each}
			<!-- <BoardSimple gameData={{current_state:'3,k,5/3,psc,m,4/5,pyt,3/9/9/9/5,PT,P,2/6,F,PS,1/7,M,1 6446122122210/6446212121210 w'}} />
			<BoardSimple gameData={{current_state:'3,k,5/3,psc,m,4/5,pyt,3/9/9/9/5,PT,P,2/6,F,PS,1/7,M,1 6446122122210/6446212121210 w'}} />
			<BoardSimple gameData={{current_state:'3,k,5/3,psc,m,4/5,pyt,3/9/9/9/5,PT,P,2/6,F,PS,1/7,M,1 6446122122210/6446212121210 w'}} />
			<BoardSimple gameData={{current_state:'3,k,5/3,psc,m,4/5,pyt,3/9/9/9/5,PT,P,2/6,F,PS,1/7,M,1 6446122122210/6446212121210 w'}} /> -->
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
		position: relative;
	}

	li:hover {
		outline: rgb(240, 80, 17) solid 6px;
		/* box-sizing: content-box; */
	}

	.your-turn {
		/* outline: rgba(255, 136, 81, 0.829) solid 2px; */
		outline: rgb(240, 80, 17) solid 4px;
	}

	.gameList {
		display: grid;
		gap: 4rem 2rem;
		grid-template-columns: repeat(auto-fit, 20rem);
		padding: 1rem;
		justify-content: center;
		margin-top: 1rem;
	}

	.gameList a {
		padding: 0;
	}

	.name {
		position: absolute;
	}
	.name-1 {
		top: -1.5rem;
		left: 0;
		/* background-color: red; */
		z-index: 10;
	}
	.name-2 {
		right: 0;
		top: -1.5rem;
	}
	.is-user {
		color: rgb(var(--primary));
		/* font-weight: 600; */
	}
</style>

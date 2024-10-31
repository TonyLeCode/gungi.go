<script lang="ts">
	import { dev } from '$app/environment';
	import BoardSimple from '$lib/components/BoardSimple.svelte';
	import { createPaginationStore } from '$lib/store/paginationStore.svelte';
	import type { Game } from './+page.server';
	//TODO make responsive

	let { data } = $props();

	let username = $derived(data.username as string);
	let ongoingGames = $derived(data.data.ongoingGames);
	let sortedOngoingGames = $derived(
		ongoingGames.sort((game1, game2) => {
			if (turnPlayer(game1) === turnPlayer(game2)) return 0;
			if (isUserTurn(game1)) {
				return -1;
			} else if (isUserTurn(game2)) {
				return 1;
			}
			return 0;
		})
	);

	let paginationStore = createPaginationStore(Math.max(Math.ceil(data.data.gameHistoryCount / 10), 1));
	interface completedGames {
		[index: number]: Game[];
	}
	let completedGames = $state<completedGames>({ 1: data.data.completedGames });

	$effect(() => {
		if (!completedGames[paginationStore.currentPage]) {
			const url = `/getongoinggames?offset=${10 * (paginationStore.currentPage - 1)}`;
			const res = fetch(url);
			res
				.then((res) => {
					if (res.ok) {
						return res.json();
					}
				})
				.then((data) => {
					completedGames[paginationStore.currentPage] = data;
				});
		}
	});

	function isUser(playername: string) {
		return username === playername;
	}

	function getUserColor(username1: string, username2: string) {
		if (username === username1) {
			return 'w';
		} else if (username === username2) {
			return 'b';
		}
		return 'spectator';
	}

	function turnPlayer(game: Game) {
		const fields = game.current_state.split(' ');
		const turnColor = fields[2];
		return turnColor === 'w' ? game.username1 : game.username2;
	}

	function isUserTurn(game: Game) {
		return turnPlayer(game) === username;
	}
</script>

<svelte:head>
	<title>Overview | White Monarch Server</title>
</svelte:head>

<main>
	<section>
		<h2>Current Games</h2>
		<ul class="gameList">
			{#if sortedOngoingGames.length === 0}
				<p class="empty-list-message">You currently have no ongoing games</p>
			{:else}
				{#each sortedOngoingGames as game}
					<li class:your-turn={isUserTurn(game)}>
						<div class="name name-1" class:is-user={isUser(game.username1)}>{game.username1}</div>
						<a href={`/game/${game.public_id}`}>
							<BoardSimple gameData={game} userColor={getUserColor(game.username1, game.username2)} />
						</a>
						<div class="name name-2" class:is-user={isUser(game.username2)}>{game.username2}</div>
					</li>
				{/each}
			{/if}
		</ul>
	</section>
	<section>
		<h2>Game History</h2>
		<ul class="gameHistoryList">
			{#if completedGames[1].length === 0}
				<p class="empty-list-message">You currently have no completed games</p>
			{:else}
				{#each completedGames[paginationStore.currentPage] as game}
					<li class="historyItem">
						<a href={`/game/${game.public_id}`}>
							<div>{game.date_started?.toString().slice(0, 10)}</div>
							<div>{game.username1 !== username ? game.username1 : game.username2}</div>
							<div>{game.type}</div>
							<div>{game.ruleset}</div>
							{#if game.result === 'w/r' || game.result === 'w'}
								<div>W</div>
							{:else if game.result === 'b/r' || game.result === 'b'}
								<div>B</div>
							{:else}
								<div>Draw</div>
							{/if}
						</a>
					</li>
				{/each}
			{/if}
		</ul>
		{#if completedGames[1].length > 0}
			<div class="pagination-controls">
				<button class="button-primary" onclick={() => paginationStore.prev()} disabled={!paginationStore.hasPrev}
					>&lt;</button
				>
				<input
					class="page-input"
					type="number"
					name="page"
					min="1"
					bind:value={paginationStore.currentPage}
					max={paginationStore.totalPages}
				/>
				/
				{paginationStore.totalPages}
				<button class="button-primary" onclick={() => paginationStore.next()} disabled={!paginationStore.hasNext}
					>&gt;</button
				>
			</div>
		{/if}
	</section>
</main>

<style lang="scss">
	section {
		max-width: 70rem;
		margin: 0 auto;
		margin-top: 2rem;
		padding: 0 0.5rem;
		text-align: center;
		overflow: auto;
		@media (min-width: 767px) {
			padding: 0 2rem;
		}
	}

	h2 {
		font-size: 1.25rem;
	}

	.empty-list-message {
		font-weight: 300;
		color: rgba(var(--font), 0.8);
	}

	.gameList li {
		border-radius: 4px;
		position: relative;
	}

	.gameList li:hover {
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
		z-index: 3;
	}
	.name-2 {
		right: 0;
		top: -1.5rem;
	}
	.is-user {
		color: rgb(var(--primary));
		/* font-weight: 600; */
	}
	.gameHistoryList {
		display: grid;
		min-height: 4.5rem;
		gap: 0.5rem;
		margin: auto;
		margin-top: 1rem;
		margin-bottom: 1.5rem;
		max-width: 50rem;
	}
	.historyItem {
	}
	.historyItem a {
		cursor: pointer;
		text-align: left;
		display: grid;
		grid-template-columns: 1fr 2fr 1fr 1fr 1fr;
		width: 100%;
		min-height: 3.5rem;
		gap: 1rem;
		align-items: center;
		padding: 0.5rem 1.25rem;
		border-radius: 4px;
		background-color: rgb(var(--bg-2));
		box-shadow:
			0px 5px 25px rgba(0, 0, 0, 0.05),
			0px 2px 5px rgba(0, 0, 0, 0.05);
		transition-duration: 150ms;
		transition-property: background-color;
		&:hover:not([disabled]) {
			background-color: rgb(var(--primary));
			color: rgb(242, 242, 242);
			// color: rgb(var(--bg-2));
		}
		&:active:not([disabled]) {
			background-color: rgb(var(--primary-3));
			color: rgb(242, 242, 242);
			// color: rgb(var(--bg-2));
		}
		&:focus {
			outline: 2px solid rgba(var(--primary), 0.5);
			outline-offset: 2px;
		}
		&:disabled {
			background-color: rgb(225, 225, 225);
			font-weight: 300;
			box-shadow: none;
		}
	}
	.pagination-controls {
		margin-bottom: 2rem;
	}
	.page-input {
		background-color: rgb(var(--bg-2));
		padding: 0.25rem 0.75rem;
	}
	.page-input::-webkit-outer-spin-button,
	.page-input::-webkit-inner-spin-button {
		-webkit-appearance: none;
	}
	.page-input[type='number'] {
		-moz-appearance: textfield;
		appearance: textfield;
	}
</style>

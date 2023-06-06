<script lang="ts">
	import PieceHand from './PieceHand.svelte';
	import type { dragAndDropFunction } from '$lib/utils/dragAndDrop';

	export let playerColor: string;
	export let player1: string;
	export let player2: string;
	export let hands: number[][];
	export let onBoardBlack: number | undefined;
	export let onBoardWhite: number | undefined;
	export let dragAndDrop: dragAndDropFunction;
	export let reversed: boolean;
	export let isPlayerTurn: boolean;
</script>

<div class="hands">
	<div class="hand-container">
		<div class="hand-info">
			<h3 class="name">{playerColor === 'w' ? player2 : player1}</h3>
			<div>On Board: <span>{playerColor === 'w' ? onBoardBlack : onBoardWhite}</span></div>
			<div>
				In Hand: <span
					>{hands[playerColor === 'w' ? 1 : 0].reduce((a, b) => {
						return a + b;
					})}</span
				>
			</div>
		</div>
		<div class="hand">
			{#each hands[playerColor === 'w' ? 1 : 0] as amount, i}
				{#if amount != 0}
					<PieceHand on:drop {dragAndDrop} {reversed} {playerColor} {isPlayerTurn} color={playerColor === 'w' ? 'b' : 'w'} piece={i} {amount} />
				{/if}
			{/each}
		</div>
	</div>
	<div class="hand-container">
		<div class="hand-info">
			<h3 class="name">{playerColor === 'w' ? player1 : player2}</h3>
			<span>On Board: {playerColor === 'w' ? onBoardWhite : onBoardBlack}</span>
			<span
				>In Hand: {hands[playerColor === 'w' ? 0 : 1].reduce((a, b) => {
					return a + b;
				})}</span
			>
		</div>
		<div class="hand">
			{#each hands[playerColor === 'w' ? 0 : 1] as amount, i}
				{#if amount != 0}
					<PieceHand on:drop {dragAndDrop} {reversed} {playerColor} {isPlayerTurn} color={playerColor === 'w' ? 'w' : 'b'} piece={i} {amount} />
				{/if}
			{/each}
		</div>
	</div>
	<div class="buttons">
		<button class="button-primary">resign</button>
		<button class="button-primary">request undo</button>
		<button class="button-primary" disabled>confirm move</button>
	</div>
</div>

<style lang="scss">
	.name {
		/* color: rgb(var(--primary)); */
		margin-right: auto;
		margin-left: 0.5rem;
		position: relative;
		&::before {
			content: '';
			width: 15px;
			height: 15px;
			border-radius: 50%;
			background-color: rgb(var(--primary));
			display: block;
			position: absolute;
			left: -1.25rem;
			margin: auto;
			top: 0;
			bottom: 0;
		}
	}

	.hand-container {
		display: flex;
		flex-direction: column;
		/* border: 2px solid rgba(255, 77, 7, 0.7); */
		background-color: rgb(var(--bg-2));
		/* background-color: rgb(255, 255, 255); */
		box-shadow: 0px 5px 15px rgba(0, 0, 0, 0.07);
		border-radius: 8px;
		padding: 1.5rem 2rem;
	}

	.hands {
		display: grid;
		gap: 1rem;
	}
	.hand {
		/* background-color: rgb(199, 199, 199); */
		display: grid;
		grid-template-columns: repeat(auto-fit, 3rem);
		gap: 1rem;
		padding: 0.5rem 0.5rem;
	}

	.hand-info {
		display: flex;
		gap: 1rem;
		padding: 0.5rem 1rem;
		padding-bottom: 1rem;
		/* justify-content: space-between; */
	}

	.buttons {
		/* background-color: red; */
		display: flex;
		gap: 1rem;
		justify-content: center;
		padding: 1rem 0;
	}
</style>

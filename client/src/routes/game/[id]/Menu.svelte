<script lang="ts">
	import type { Droppable } from '$lib/store/dragAndDrop.svelte';
	import MenuHand from './MenuHand.svelte';
	import MenuMoveHistory from './MenuMoveHistory.svelte';

	type DropItem = {
		destinationIndex: number;
		destinationStack: number[];
	};

	let {
		selectedStack,
		ready,
		placeHandMove,
		droppable,
		selectHandPiece,
		resign,
		undo
	}: {
		selectedStack: number[];
		ready: () => void;
		placeHandMove: (fromPiece: number, toCoord: number) => void;
		droppable: Droppable<DropItem>;
		selectHandPiece: (piece: number) => void;
		resign: () => void;
		undo: () => void;
	} = $props();

	let menuState = $state('hand');
</script>

<aside class="side-menu">
	<div class="tabs">
		<button class="tab" class:active-tab={menuState === 'hand'} onclick={() => (menuState = 'hand')}>Hand</button>
		<button class="tab" class:active-tab={menuState === 'move history'} onclick={() => (menuState = 'move history')}
			>Move History</button
		>
		<!-- <button class="tab" class:active-tab={menuState === "chat"} onclick={() => menuState = "chat"}>Chat</button> -->
	</div>
	{#if menuState === 'hand'}
		<MenuHand {resign} {undo} {selectHandPiece} {selectedStack} {ready} {placeHandMove} {droppable} />
	{:else if menuState === 'move history'}
		<MenuMoveHistory />
	{/if}
</aside>

<style lang="scss">
	.side-menu {
		display: flex;
		flex-direction: column;
		margin: 0 auto;
		margin-top: 0.5rem;
		max-width: 44rem;
		width: 100%;
		padding: 0 2rem;
		@media (min-width: 1200px) {
			gap: 1rem;
			// margin: auto 0;
			margin-top: 5rem;
			margin-left: auto;
			max-width: 36rem;
			padding: 0;
		}
	}

	.tabs {
		display: grid;
		// grid-template-columns: repeat(3, 1fr);
		grid-template-columns: repeat(2, 1fr);
		justify-content: center;
		margin: 1.5rem 0;
		border-radius: 4px;
		order: 1;
		@media (min-width: 1200px) {
			order: 0;
			margin: 0.5rem 0;
		}
	}

	.tab {
		background-color: rgb(var(--bg-2));
		padding: 0.25rem 0.5rem;
		&:hover {
			background-color: rgb(var(--primary));
			color: white;
		}
		@media (min-width: 1200px) {
			padding: 0.5rem 1rem;
		}
	}

	.active-tab {
		background-color: rgb(var(--primary));
		color: white;
	}
</style>

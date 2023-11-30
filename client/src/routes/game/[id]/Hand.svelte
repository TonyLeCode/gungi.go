<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import PieceHand from './PieceHand.svelte';
	import type { dragAndDropFunction } from '$lib/utils/dragAndDrop';

	import {
		manualFlipContext,
		player1HandListContext,
		player2HandListContext,
		userColorContext,
		player1NameContext,
		player2NameContext,
		isPlayer1ReadyContext,
		isPlayer2ReadyContext,
		player1ArmyCountContext,
		player2ArmyCountContext,
		player1HandCountContext,
		player2HandCountContext,
		isUserTurnContext,
	} from './+page.svelte';
	import Stack from './Stack.svelte';

	const manualFlip = manualFlipContext.get();
	const player1HandList = player1HandListContext.get();
	const player2HandList = player2HandListContext.get();
	const isPlayer1Ready = isPlayer1ReadyContext.get();
	const isPlayer2Ready = isPlayer2ReadyContext.get();
	const userColor = userColorContext.get();
	const player1Name = player1NameContext.get();
	const player2Name = player2NameContext.get();
	const player1ArmyCount = player1ArmyCountContext.get();
	const player2ArmyCount = player2ArmyCountContext.get();
	const player1HandCount = player1HandCountContext.get();
	const player2HandCount = player2HandCountContext.get();
	const isUserTurn = isUserTurnContext.get();
	export let dragAndDrop: dragAndDropFunction;
	export let stack: number[];

	const dispatch = createEventDispatcher();

	function readyDisplay(num: number, playerColor: string): string {
		if (num === 0) {
			if (playerColor === 'w' && $isPlayer2Ready) {
				return '- ready';
			}
			if (playerColor === 'b' && $isPlayer1Ready) {
				return '- ready';
			}
		} else if (num === 1) {
			if (playerColor === 'w' && $isPlayer1Ready) {
				return '- ready';
			}
			if (playerColor === 'b' && $isPlayer2Ready) {
				return '- ready';
			}
		}
		return '';
	}

	function handleFlipBoardButton() {
		manualFlip.update((val) => {
			return !val;
		});
	}

	function handleResignButton() {
		dispatch('resign');
	}
	function handleUndoButton() {
		dispatch('undo');
	}
	function handleReadyButton() {
		dispatch('ready');
	}
</script>

<div class="hands">
	<Stack {stack} />
	<div class="hand-container">
		<div class="hand-info">
			<h3 class="name">{$userColor === 'w' ? $player2Name : $player1Name} {readyDisplay(0, $userColor)}</h3>
			<div>On Board: <span>{$userColor === 'w' ? $player2ArmyCount : $player1ArmyCount}</span></div>
			<div>
				In Hand: <span>{$userColor === 'b' ? $player1HandCount : $player2HandCount}</span>
			</div>
		</div>
		<div class="hand">
			{#each $userColor === 'b' ? $player1HandList : $player2HandList as amount, i}
				{#if amount != 0}
					<PieceHand
						on:drop
						{dragAndDrop}
						playerColor={$userColor}
						isPlayerTurn={$isUserTurn}
						color={$userColor === 'w' ? 'b' : 'w'}
						piece={i}
						{amount}
					/>
				{/if}
			{/each}
		</div>
	</div>
	<div class="hand-container">
		<div class="hand-info">
			<h3 class="name">{$userColor === 'w' ? $player1Name : $player2Name} {readyDisplay(1, $userColor)}</h3>
			<div>On Board: <span>{$userColor === 'w' ? $player1ArmyCount : $player2ArmyCount}</span></div>
			<div>In Hand:<span>{$userColor === 'w' ? $player1HandCount : $player2HandCount}</span></div>
		</div>
		<div class="hand">
			{#each $userColor === 'w' ? $player1HandList : $player2HandList as amount, i}
				{#if amount != 0}
					<PieceHand
						on:drop
						{dragAndDrop}
						playerColor={$userColor}
						isPlayerTurn={$isUserTurn}
						color={$userColor === 'w' ? 'w' : 'b'}
						piece={i}
						{amount}
					/>
				{/if}
			{/each}
		</div>
	</div>
	<div class="buttons">
		<button class="button-primary" on:click={handleResignButton}>resign</button>
		<button class="button-primary" on:click={handleUndoButton}>request undo</button>
		<button class="button-primary" disabled>confirm move</button>
		{#if !$isPlayer1Ready || !$isPlayer2Ready}
			<button class="button-primary" on:click={handleReadyButton}>ready</button>
		{/if}
		<button class="button-primary" on:click={handleFlipBoardButton}>flip board</button>
	</div>
</div>

<style lang="scss">
	.name {
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
		background-color: rgb(var(--bg-2));
		box-shadow:
			0px 2px 55px rgba(0, 0, 0, 0.07),
			0px 4px 15px rgba(0, 0, 0, 0.05);
		border-radius: 8px;
		padding: 1.5rem 2rem;
	}

	.hands {
		display: grid;
		gap: 1rem;
	}
	.hand {
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
	}

	.buttons {
		display: flex;
		gap: 1rem;
		justify-content: center;
		padding: 1rem 0;
	}
</style>

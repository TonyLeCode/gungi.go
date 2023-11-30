<script lang="ts">
	import { onDestroy } from 'svelte';
	import { moveHistoryContext, gameStateContext } from './+page.svelte';
	import { get } from 'svelte/store';

	const moveHistory = moveHistoryContext.get();
	const gameState = gameStateContext.get();
	let currentIndex = get(moveHistory).length - 1;
	moveHistory.subscribe((val) => {
		currentIndex = val.length - 1;
	});
	let containerRef: HTMLOListElement;
	let autoplay = false;
	let autoplayInterval = 1000;
	let intervalId: NodeJS.Timer | null;

	$: if (containerRef) {
		containerRef.scrollTop = currentIndex * 36;
	}

	function handleCopy() {
		const state = get(gameState);
		console.log(state);
		let headings = `[game_id: ${state.id}]\n[ruleset: ${state.ruleset}]\n[type: ${state.type}]\n[date_started: ${state.date_started}]\n[white: ${state.player1}]\n[black: ${state.player2}]\n\n`;
		const history = get(moveHistory).join(' ');
		navigator.clipboard.writeText(headings + history);
	}

	function startAutoplay() {
		intervalId = setInterval(() => {
			if (currentIndex === get(moveHistory).length - 1) {
				if (intervalId) {
					clearInterval(intervalId);
					autoplay = false;
					return;
				}
			}
			currentIndex++;
		}, autoplayInterval);

		autoplay = true;
	}
	function stopAutoplay() {
		if (intervalId) {
			clearInterval(intervalId);
			intervalId = null;
		}

		autoplay = false;
	}
	function intervalChange() {
		if (intervalId) {
			clearInterval(intervalId);
			intervalId = setInterval(() => {
				if (currentIndex === get(moveHistory).length - 1) {
					if (intervalId) {
						clearInterval(intervalId);
						autoplay = false;
						return;
					}
				}
				currentIndex++;
			}, autoplayInterval);
		}
	}

	function next() {
		if (currentIndex === get(moveHistory).length - 1) {
			return;
		}
		currentIndex++;
	}
	function prev() {
		if (currentIndex === 0) {
			return;
		}
		currentIndex--;
	}

	onDestroy(() => {
		if (intervalId) {
			clearInterval(intervalId);
		}
	});
</script>

<ol class="history-list" bind:this={containerRef}>
	{#each $moveHistory as move, index (index + JSON.stringify($moveHistory))}
		<li>
			<button
				class={`${index === currentIndex ? 'current-index' : ''}`}
				on:click={() => {
					currentIndex = index;
				}}
			>
				<span class="index">{index + 1}:</span>
				<span>{move}</span>
			</button>
		</li>
	{/each}
</ol>
<input
	class="replay-bar"
	bind:value={currentIndex}
	type="range"
	step="1"
	name="replay"
	min="0"
	max={$moveHistory.length - 1}
/>
<div class="controls">
	<button class="button-primary" on:click={prev} disabled={currentIndex === 0}>&lt;</button>
	<div class="counter">
		{currentIndex + 1}/{$moveHistory.length}
	</div>
	<button class="button-primary" on:click={next} disabled={currentIndex === $moveHistory.length - 1}> &gt;</button>
	<div class="autoplay-controls">
		<button class="button-primary" disabled={!autoplay} on:click={stopAutoplay}>stop</button>
		<button
			class="button-primary"
			disabled={autoplay || currentIndex === $moveHistory.length - 1}
			on:click={startAutoplay}>autoplay</button
		>
		<input type="number" name="autoplay" bind:value={autoplayInterval} on:change={intervalChange} />
		<button class="button-primary" on:click={handleCopy}>Copy</button>
	</div>
</div>

<style lang="scss">
	.counter {
		width: 3.5rem;
		margin: 0 0.25rem;
		display: flex;
		justify-content: center;
		align-items: center;
	}
	.autoplay-controls {
		// margin-left: auto;
		margin-left: 2rem;
		display: flex;
		gap: 0.5rem;
	}
	.autoplay-controls input {
		width: 5rem;
		padding: 0 0.5rem;
	}
	.replay-bar {
		width: 100%;
		margin-bottom: 1rem;
		margin-top: 1rem;
		background-color: rgb(var(--bg-2));
		// height: 4rem;
	}
	.controls {
		display: flex;
		background-color: rgb(var(--bg-2));
		box-shadow:
			0px 2px 55px rgba(0, 0, 0, 0.07),
			0px 4px 15px rgba(0, 0, 0, 0.05);
		margin-bottom: 1rem;
		padding: 1rem;
		justify-content: center;
	}
	.history-list {
		background-color: rgb(var(--bg-2));
		box-shadow:
			0px 2px 55px rgba(0, 0, 0, 0.07),
			0px 4px 15px rgba(0, 0, 0, 0.05);
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		border-radius: 4px;
		height: 30rem;
		overflow: auto;
	}
	li {
		// background-color: rgb(230, 230, 230);
	}
	li button {
		width: 100%;
		text-align: left;
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		&:hover {
			background-color: rgb(var(--primary));
			color: white;
		}
	}

	.index {
		width: 2rem;
		display: inline-block;
	}
	.current-index {
		background-color: rgb(var(--primary));
		color: white;
	}
</style>

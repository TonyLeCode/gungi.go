<script lang="ts">
	import { getGameStore } from '$lib/store/gameState.svelte';
	import { onMount } from 'svelte';

	const boardStore = getGameStore();

	let currentMoveHistoryIndex = $state(boardStore.moveHistory.length - 1);
	let autoplayInterval = $state(500);
	let autoplay = $state(false);
	let intervalId = $state<number | null>(null);
	let containerRef: HTMLOListElement;

	$effect(() => {
		currentMoveHistoryIndex = boardStore.moveHistory.length - 1;
	})
	$effect(() => {
		containerRef.scrollTop = Math.floor(currentMoveHistoryIndex / 2) * 36;
	});

	let url = $state('');
	let fileText = $state('');
	let fileName = $state('');

	let pliedHistory = $derived.by(() => {
		const newList = [];
		let history = boardStore.moveHistory;
		for (let i = 0; i < history.length; i += 2) {
			let tuple = [history[i], history[i + 1]];
			newList.push(tuple);
		}
		return newList;
	});

	function prev() {
		if (currentMoveHistoryIndex > 0) {
			currentMoveHistoryIndex--;
		}
	}

	function next() {
		if (currentMoveHistoryIndex < boardStore.moveHistory.length - 1) {
			currentMoveHistoryIndex++;
		}
	}

	function startAutoplay() {
		intervalId = window.setInterval(() => {
			console.log('iteration');
			if (currentMoveHistoryIndex === boardStore.moveHistory.length - 1) {
				if (intervalId) {
					clearInterval(intervalId);
					autoplay = false;
					return;
				}
			}
			currentMoveHistoryIndex++;
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
			intervalId = window.setInterval(() => {
				if (currentMoveHistoryIndex === boardStore.moveHistory.length - 1) {
					if (intervalId) {
						clearInterval(intervalId);
						autoplay = false;
						return;
					}
				}
				currentMoveHistoryIndex++;
			}, autoplayInterval);
		}
	}

	onMount(() => {
		$effect(() => {
			const headings = `[game_id : ${boardStore.id}]\n[ruleset : ${boardStore.ruleset}]\n[type : ${boardStore.type}]\n[date_started : ${boardStore.date_started}]\n[white : ${boardStore.player1}]\n[black : ${boardStore.player2}]\n\n`;
			fileText = headings + boardStore.moveHistory.join(' ');
			url = URL.createObjectURL(new Blob([fileText], { type: 'text/plain' }));
			const date = new Date(boardStore.date_started);
			const year = date.getFullYear();
			const month = String(date.getMonth() + 1).padStart(2, '0');
			const day = String(date.getDate()).padStart(2, '0');
			const formattedDate = `${year}-${month}-${day}`;
			fileName = `${boardStore.player1}_vs_${boardStore.player2}_${formattedDate}`;
		});
	});

	function handleCopy() {
		navigator.clipboard.writeText(fileText);
	}
</script>

<h3>Move History: </h3>
<ol class="history-list" bind:this={containerRef}>
	{#each pliedHistory as move, index}
		<li>
			<span class="move-number">{index + 1}</span>
			<button
				class:current-index={Math.floor(currentMoveHistoryIndex / 2) === index && currentMoveHistoryIndex % 2 === 0}
				onclick={() => (currentMoveHistoryIndex = index * 2)}
			>
				{move[0]}
			</button>
			{#if move[1] !== undefined}
				<button
					class:current-index={Math.floor(currentMoveHistoryIndex / 2) === index && currentMoveHistoryIndex % 2 === 1}
					onclick={() => (currentMoveHistoryIndex = index * 2 + 1)}
				>
					{move[1]}
				</button>
			{/if}
		</li>
	{/each}
</ol>
<div class="controls">
	<div class="step-controls">
		<button class="button-primary" onclick={() => (currentMoveHistoryIndex = 0)}>&#171;</button>
		<button class="button-primary" onclick={prev}>&lt;</button>
		<button class="button-primary" onclick={next}>&gt;</button>
		<button class="button-primary" onclick={() => (currentMoveHistoryIndex = boardStore.moveHistory.length - 1)}
			>&#187;</button
		>
	</div>
	<div class="autoplay-controls">
		<button class="button-primary" disabled={!autoplay} onclick={stopAutoplay}>stop</button>
		<button class="button-primary" disabled={autoplay} onclick={startAutoplay}>autoplay</button>
		<input
			type="number"
			name="autoplay"
			min="10"
			max="10000"
			step="10"
			bind:value={autoplayInterval}
			onchange={intervalChange}
		/>
	</div>
</div>
<div class="options">
	<button class="button-primary" onclick={handleCopy}>Copy</button>
	<a class="button-primary" href={url} download={fileName}>Download</a>
</div>

<style lang="scss">
  h3 {
    margin: 0.5rem 0;
    order: -1;
    @media (min-width: 1200px) {
      order: 0;
    }
  }
	.history-list {
		background-color: rgb(var(--bg-2));
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		border-radius: 4px;
		height: 25rem;
		overflow: auto;
    margin: 0.5rem 0;
		border-radius: 4px;
		@media (min-width: 767px) {
			height: 25rem;
		}
	}

	.move-number {
		text-align: right;
		padding: 0.25rem 0.5rem;
	}

	li {
		display: grid;
		grid-template-columns: 3rem 2fr 2fr;
		padding-right: 4rem;
		@media (min-width: 1200px) {
			padding-right: 14rem;
		}
	}

	li button {
		text-align: left;
		padding: 0.25rem 0.5rem;
		&:hover {
			background-color: rgb(var(--primary));
			color: white;
		}
	}

	.current-index {
		background-color: rgb(var(--primary));
		color: white;
	}

	.controls {
    display: flex;
		gap: 0.5rem;
		justify-content: center;
    // margin: 1rem 0;
    order: -1;
    @media (min-width: 1200px) {
      order: 0;
    }
	}

  .step-controls {
    display: flex;
    gap: 0.5rem;
    justify-content: space-evenly;
  }

	.autoplay-controls {
		display: flex;
		gap: 0.5rem;
	}

	.autoplay-controls input {
		width: 3.5rem;
		padding: 0.25rem 0.5rem;
    @media (min-width: 767px) {
      width: 5rem;
    }
	}

	span {
		padding: 0.25rem 0.5rem;
	}

	.options {
		display: flex;
		gap: 1rem;
		justify-content: center;
	}

	input {
		background-color: rgb(var(--bg-3));
	}
</style>

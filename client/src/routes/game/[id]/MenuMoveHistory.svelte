<script lang="ts">
	import { getGameStore } from '$lib/store/gameState.svelte';
	import { getReplayStore } from '$lib/store/replayStore.svelte';
	import { onMount } from 'svelte';

	const boardStore = getGameStore();

	const replayStore = getReplayStore();
	let autoplayInterval = $state(500);
	let autoplay = $state(false);
	let intervalId = $state<number | null>(null);
	let containerRef: HTMLOListElement;

	$effect(() => {
		const childElement = containerRef.children[replayStore.pagination.currentPage - 1];
		if (!childElement || !(childElement instanceof HTMLLIElement)) {
			console.error('no child element');
			return;
		}

		const height = childElement.clientHeight + 4;
		containerRef.scrollTop = (replayStore.pagination.currentPage - 1) * height;
	});

	let url = $state('');
	let fileText = $state('');
	let fileName = $state('');

	function startAutoplay() {
		intervalId = window.setInterval(() => {
			console.log('iteration');
			if (!replayStore.pagination.hasNext) {
				if (intervalId) {
					clearInterval(intervalId);
					autoplay = false;
					return;
				}
			}
			replayStore.next();
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
				if (!replayStore.pagination.hasNext) {
					if (intervalId) {
						clearInterval(intervalId);
						autoplay = false;
						return;
					}
				}
				replayStore.next();
			}, autoplayInterval);
		}
	}

	onMount(() => {
		$effect(() => {
			const headings = `[game_id : ${boardStore.id}]\n[ruleset : ${boardStore.ruleset}]\n[type : ${boardStore.type}]\n[date_started : ${boardStore.date_started}]\n[white : ${boardStore.player1}]\n[black : ${boardStore.player2}]\n\n`;
			fileText = headings + boardStore.moveHistory.join(' ');
			url = URL.createObjectURL(new Blob([fileText], { type: 'text/plain' }));
			const date = new Date(boardStore.date_started as Date);
			const year = date.getFullYear();
			const month = String(date.getMonth() + 1).padStart(2, '0');
			const day = String(date.getDate()).padStart(2, '0');
			const formattedDate = `${year}-${month}-${day}`;
			fileName = `${boardStore.player1}_vs_${boardStore.player2}_${formattedDate}`;
		});
	});

	onMount(() => {
		return () => {
			replayStore.setPage(boardStore.moveHistory.length);
		}
	})
	
	function handleCopy() {
		navigator.clipboard.writeText(fileText);
	}
</script>

<h3>Move History:</h3>
<ol class="history-list" bind:this={containerRef}>
	{#each boardStore.moveHistory as move, index}
		<li class="move-history-item">
			<span class="move-history-number">{index + 1}</span>
			<button
				class="move-history-button"
				class:current-index={replayStore.pagination.currentPage === index + 1}
				onclick={() => replayStore.setPage(index + 1)}
			>
				{move}
			</button>
		</li>
	{/each}
</ol>
<div class="controls">
	<div class="step-controls">
		<button class="button-primary" disabled={!replayStore.pagination.hasPrev} onclick={() => replayStore.setPage(1)}>&#171;</button>
		<button class="button-primary" disabled={!replayStore.pagination.hasPrev} onclick={() => replayStore.prev()}>&lt;</button>
		<button class="button-primary" disabled={!replayStore.pagination.hasNext} onclick={() => replayStore.next()}>&gt;</button>
		<button class="button-primary" disabled={!replayStore.pagination.hasNext} onclick={() => replayStore.setPage(replayStore.pagination.totalPages)}>&#187;</button
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
		// @media (min-width: 767px) {
		// 	height: 25rem;
		// }
	}

	.move-history-number {
		text-align: right;
		padding: 0.25rem 0.5rem;
	}

	.move-history-item {
		display: flex;
	}

	.move-history-button {
		text-align: left;
		padding: 0.25rem 0.5rem;
		flex-grow: 1;
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
		margin: 0.5rem 0;
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
		margin-top: 0.5rem;
		@media (min-width: 1200px) {
			margin-top: 0;
		}
	}

	input {
		background-color: rgb(var(--bg-3));
	}
</style>

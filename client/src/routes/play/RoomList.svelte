<script lang="ts">
	import Modal from '$lib/components/Modal.svelte';
	import { createPaginationStore } from '$lib/store/paginationStore.svelte';
	import { getWebsocketStore } from '$lib/store/websocketStore.svelte';

	type Info = {
		id: string;
		host: string;
		description: string;
		type: string;
		color: string;
		rules: string;
	};

	let websocketStore = getWebsocketStore();
	let {
		username,
		roomList,
		heading,
		roomDialog,
		roomDialogueInfo = $bindable(),
		accept,
	}: {
		username: string;
		roomList: Info[];
		heading: string;
		roomDialog: ReturnType<typeof Modal> | undefined;
		roomDialogueInfo: Info;
		accept: (roomid: string) => void;
	} = $props();

	let totalPages = $derived.by(() => {
		const pages = Math.ceil(roomList.length / 10)
		if (pages < 1) return 1
		return pages
	});
	const paginationStore = createPaginationStore(totalPages);

	$effect(() => {
		let pages = totalPages
		if (pages < 1) pages = 1
		paginationStore.setTotalPages(pages);
	});

	let spectator = $derived(username == null);

	function handleCancel(roomid: string) {
		const payload = {
			type: 'cancelPlayRoom',
			payload: roomid,
		};
		websocketStore.send(payload);
	}
</script>

<h2 class="fly-up-fade">{heading}</h2>
{#if roomList.length != 0}
	<ul class="room-list">
		{#each roomList.slice((paginationStore.currentPage - 1) * 10, paginationStore.currentPage * 10) as room, index (room.id)}
			<li class="room-item fly-up-fade" style={`animation-delay:${String((index + 1) * 25)}ms;`}>
				<button
					disabled={room.host === username || spectator ? true : false}
					onclick={() => {
						roomDialogueInfo = room;
						roomDialog?.open();
					}}
					class="room"
				>
					<div>{room.host}</div>
					<div class="ruleset">{room.rules}</div>
					<div>{room.description}</div>
				</button>
				{#if room.host === username}
					<button
						disabled={spectator}
						onclick={() => {
							handleCancel(room.id);
						}}
						class="cancel button-ghost">Cancel</button
					>
				{:else}
					<button
						class="accept button-primary"
						disabled={spectator}
						onclick={() => {
							accept(room.id);
						}}>Accept</button
					>
				{/if}
			</li>
		{/each}
	</ul>
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
{:else}
	<p class="empty fly-up-fade">Looks like there are no {heading.toLowerCase()} available</p>
{/if}

<style lang="scss">
	h2 {
		margin-bottom: 0.5rem;
	}
	.room-item {
		display: flex;
	}
	.empty {
		margin-bottom: 2rem;
		font-weight: 300;
		animation-delay: 50ms;
	}
	.room-list {
		display: grid;
		min-height: 4.5rem;
		word-break: break-all;
		gap: 0.5rem;
		margin-bottom: 1.5rem;
		@media (min-width: 767px) {
		}
	}
	.room {
		text-align: left;
		display: grid;
		grid-template-columns: 12rem 8rem 1fr;
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
			color: white;
		}
		&:active:not([disabled]) {
			background-color: rgb(var(--primary-3));
			color: white;
		}
		&:focus {
			outline: 2px solid rgba(var(--primary), 0.5);
			outline-offset: 2px;
		}
		&:disabled {
			// background-color: rgba(36, 36, 36, 0.5);
			// background-color: rgb(225, 225, 225);
			background-color: rgb(var(--bg-4));
			color: rgba(var(--font), 0.8);
			font-weight: 300;
			box-shadow: none;
		}
	}
	.ruleset {
		font-weight: 300;
		text-transform: capitalize;
	}
	.cancel,
	.accept {
		word-break: normal;
		margin: auto 1rem;
		margin-right: 0;
		height: min-content;
		width: 5.5rem;
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

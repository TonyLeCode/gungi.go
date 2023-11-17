<script lang="ts">
	import { ws } from '$lib/store/websocket';
	type Info = {
		id: string;
		host: string;
		description: string;
		type: string;
		color: string;
		rules: string;
	};
	export let username: string;

	export let roomList: Info[];
	export let heading: string;

	export let showRoomDialogue: boolean;
	export let roomDialogueInfo: Info;
	export let accept: (roomid: string) => void;

	function handleCancel(roomid: string) {
		const payload = {
			type: 'cancelPlayRoom',
			payload: roomid,
		};
		ws?.send(payload);
	}
</script>

<h2 class="fly-up-fade">{heading}</h2>
{#if roomList.length != 0}
	<ul class="room-list">
		{#each roomList ?? [] as room, index (room.id)}
			<li class="room-item fly-up-fade" style={`animation-delay:${String((index + 1) * 25)}ms;`}>
				<button
					disabled={room.host === username ? true : false}
					on:click={() => {
						roomDialogueInfo = room;
						showRoomDialogue = true;
					}}
					class="room"
				>
					<div>{room.host}</div>
					<div class="ruleset">{room.rules}</div>
					<div>{room.description}</div>
				</button>
				{#if room.host === username}
					<button
						on:click={() => {
							handleCancel(room.id);
						}}
						class="cancel button-ghost">Cancel</button
					>
				{:else}
					<button
						class="accept button-primary"
						on:click={() => {
							accept(room.id);
						}}>Accept</button
					>
				{/if}
			</li>
		{/each}
	</ul>
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
		// box-shadow: 0px 5px 15px rgba(0, 0, 0, 0.05);
		box-shadow:
			0px 5px 25px rgba(0, 0, 0, 0.05),
			0px 2px 5px rgba(0, 0, 0, 0.05);
		transition-duration: 150ms;
		transition-property: background-color;
		&:hover:not([disabled]) {
			background-color: rgb(var(--primary));
			color: rgb(var(--bg-2));
		}
		&:active:not([disabled]) {
			background-color: rgb(var(--primary-3));
			color: rgb(var(--bg-2));
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
</style>

<script lang="ts">
	import { fly } from 'svelte/transition';

	type Info = {
		roomid: string;
		host: string;
		description: string;
		type: string;
		color: string;
		rules: string;
	};
	export let roomList: Info[];

	export let showRoomDialogue: boolean;
	export let roomDialogueInfo: Info;
</script>

<ul class="room-list">
	{#each roomList ?? [] as room, index (room.roomid)}
		<li class='fly-up' style={`animation-delay:${String((index + 1) * 25)}ms;`}>
			<button
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
		</li>
	{/each}
</ul>

<style lang="scss">
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
		min-height: 4rem;
		gap: 1rem;
		align-items: center;
		padding: 0.5rem 2rem;
		border-radius: 4px;
		background-color: rgb(var(--bg-2));
		box-shadow: 0px 5px 15px rgba(0, 0, 0, 0.05);
		transition-duration: 150ms;
		transition-property: background-color;
		&:hover {
			background-color: rgb(var(--primary));
			color: rgb(var(--bg-2));
		}
		&:active {
			background-color: rgb(var(--primary-3));
			color: rgb(var(--bg-2));
		}
		&:focus {
			outline: 2px solid rgba(var(--primary), 0.5);
			outline-offset: 2px;
		}
	}
	.ruleset {
		font-weight: 300;
		text-transform: capitalize;
	}
</style>

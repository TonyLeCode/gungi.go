<script lang='ts'>
  type Info = {
		name: string;
		description: string;
		type: string;
		color: string;
		rules: string;
	};
  export let roomList;

  export let showRoomDialogue:boolean;
  export let roomDialogueInfo: Info;
</script>

<ul class="room-list">
  <li>
    <button
      on:click={() => {
        showRoomDialogue = true;
      }}
      class="room"
    >
      <div>PlayerName</div>
      <div class="ruleset">Ruleset</div>
      <div>DescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescription</div>
    </button>
  </li>
  <li>
    <button
      on:click={() => {
        showRoomDialogue = true;
      }}
      class="room"
    >
      <div>PlayerName</div>
      <div class="ruleset">Ruleset</div>
      <div>DescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescription</div>
    </button>
  </li>
</ul>

<style lang='scss'>
  .room-list {
		display: grid;
		min-height: 4.5rem;
		word-break: break-all;
		gap: 1rem;
	}
	.room {
		text-align: left;
		display: grid;
		grid-template-columns: 12rem 8rem 1fr;
		gap: 1rem;
		align-items: center;
		padding: 0.5rem 2rem;
		border-radius: 4px;
		background-color: rgb(var(--bg-2));
		box-shadow: 0px 5px 15px rgba(0, 0, 0, 0.07);
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
	}
</style>
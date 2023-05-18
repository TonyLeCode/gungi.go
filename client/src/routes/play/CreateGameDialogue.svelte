<script lang="ts">
	import Modal from '$lib/components/Modal.svelte';

	export let showModal: boolean;
</script>


<Modal bind:showModal>
	<div class="options">
    <h3>Create Game</h3>
		<fieldset class="type">
			<legend>Type:</legend>
			<label>
				<input type="radio" name="type" value="live" />
				Live
			</label>
			<label>
				<input type="radio" name="type" value="correspondence" checked />
				Correspondence
			</label>
		</fieldset>
		<fieldset class="color">
			<legend>Color:</legend>
			<label>
				<input type="radio" name="color" value="white" />
				White
			</label>
			<label>
				<input type="radio" name="color" value="black" />
				Black
			</label>
			<label>
				<input type="radio" name="color" value="random" checked />
				Random
			</label>
		</fieldset>
		<fieldset class="select">
			<!-- <label class="color">
				Your Color:
				<select name="color">
					<option value="white">White</option>
					<option value="black">Black</option>
					<option value="random">Random</option>
				</select>
			</label> -->
			<label class="ruleset">
				Ruleset:
				<select name="ruleset">
					<option value="default">Default</option>
					<option value="universal-music">Universal Music</option>
					<option value="revised">Revised</option>
				</select>
			</label>
		</fieldset>
		<fieldset>
			<label class="description">
				Description:
				<textarea name="description" cols="30" rows="2" maxlength="50" />
			</label>
		</fieldset>
	</div>
	<button class="button-primary">Create Challenge</button>
</Modal>

<style lang="scss">
  h3{
    font-size: 1.25rem;
    font-weight: 600;
		text-align: center;
  }
	textarea {
		resize: none;
	}
	textarea,
	select {
		border: 1px black solid;
		border-radius: 4px;
		padding: 0.25rem 0.5rem;
		border: 1.5px solid rgba(var(--primary), 0.25);
		background-color: rgb(var(--bg));
	}
	.description,
	.select,
	.color,
	.ruleset,
  .type {
		display: flex;
		flex-direction: column;
	}
	.description {
		margin-bottom: 1rem;
	}
	.options {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}
	button{
		width: 100%;
	}
</style>

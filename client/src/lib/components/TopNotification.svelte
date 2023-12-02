<script lang="ts">
	export let text: string;
</script>

<div class="fly-up-fade">{text}</div>

<style lang="scss">
	div {
		position: fixed;
		top: 1rem;
		text-align: center;
		left: 5px;
		right: 5px;
		background-color: rgb(var(--primary));
		color: rgb(var(--bg));
		padding: 1rem;
		margin: 0 auto;
		max-width: 60rem;
		border-radius: 4px;
		z-index: 4;
		box-shadow:
			0px 5px 15px rgba(0, 0, 0, 0.07),
			0px 2px 5px rgba(0, 0, 0, 0.05);
		animation-duration: 400ms;
	}
</style>

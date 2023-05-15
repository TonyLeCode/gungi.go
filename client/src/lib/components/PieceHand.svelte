<script lang="ts">
	import { DecodePiece } from '$lib/utils/utils';

	export let color: string;
	export let piece: number;
	export let amount: number;
	const decodedPiece = DecodePiece(piece).toLowerCase();
</script>

<div class={`hand ${color === 'b' ? 'dark' : ''}`}>
	<img src={`/pieces/${color}1${decodedPiece}.svg`} alt="" />
	<div class="amount">
		{amount}
	</div>
</div>

<style>
	img {
		width: 3rem;
		border-radius: 50%;
		box-shadow: 0px 3px 10px rgba(0, 0, 0, 0.17);
	}
	img:hover {
		outline: var(--primary) 5px solid;
	}
	.hand {
		position: relative;
		--bg-color: var(--primary);
		--fill-color: var(--white);
	}
	.dark {
		--bg-color: var(--white);
		--fill-color: var(--primary);
	}
	.amount {
		text-align: center;
		font-size: 0.7rem;
		background-color: var(--bg-color);
		color: var(--fill-color);
		border: 1px solid var(--fill-color);
		border-radius: 50%;
		width: 22px;
		height: 22px;
		display: flex;
		align-items: center;
		justify-content: center;
		position: absolute;
		right: -4px;
		top: -4px;
	}
</style>

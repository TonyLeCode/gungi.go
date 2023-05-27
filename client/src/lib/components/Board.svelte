<script lang="ts">
	import { reverseList } from '$lib/helpers';
	import type { dragAndDropFunction, dropFunction } from '$lib/utils/dragAndDrop';
	import {
		DecodePiece,
		DecodePieceFull,
		FenToBoard,
		FenToHand,
		GetImage,
		GetPieceColor,
		IndexToCoords,
	} from '$lib/utils/utils';

	// export const boardState = new Array(81).fill(['']);
	export let gameData;
	export let reversed: boolean;
	export let dragAndDrop: dragAndDropFunction;
	export let drop: dropFunction;

	function temp(a: any, b: any) {
		return function temp2(c: any) {
			// console.log(a)
			// // console.log(DecodePieceFull(piece))
			// console.log(b)
			let correctedIndexB = b;
			let correctedIndexC = c;
			if (reversed) {
				correctedIndexB = 80 - b;
				correctedIndexC = 80 - c;
			}
			const [file, rank] = IndexToCoords(correctedIndexB);
			const [file2, rank2] = IndexToCoords(correctedIndexC);
			alert(
				`From: ${file.toUpperCase()}${rank} \nDestination: ${file2.toUpperCase()}${rank2} \nPiece: ${DecodePieceFull(
					a
				)}`
			);
		};
	}

	function GetImage2(tier: number, piece: number): string {
		const encodedPiece = DecodePiece(piece).toLowerCase();
		const color = GetPieceColor(piece);
		return `/pieces/${color}${tier}${encodedPiece}.svg`;
	}

	function reverseIfBlack<T>(arr: T[]): T[] {
		if (reversed) {
			return reverseList(arr);
		} else return arr;
	}
	let boardState = reverseIfBlack(FenToBoard(gameData.current_state));
	let fileCoords = reverseIfBlack([1, 2, 3, 4, 5, 6, 7, 8, 9]);
	let rankCoords = reverseIfBlack(['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i']);

	function dropEvent(index: number) {
		let correctedIndex = index
		if(reversed){
			correctedIndex = 80 - index
		}
		//TODO temporary
		function a() {
			console.log(correctedIndex);
			const [file, rank] = IndexToCoords(correctedIndex);
			console.log("File: ", file, " Rank: ", rank)
		}
		return {
			mouseEnterItem: correctedIndex,
			mouseEnterEvent: a,
		};
	}
</script>

<div class="board">
	{#each boardState as square, index}
		<div
			class="square"
			use:drop={dropEvent(index)}
			on:focus={() => {
				console.log('');
			}}
		>
			{#if square.length > 0}
				<img draggable="false" use:dragAndDrop class="piece" src={GetImage(square)} alt="" />
				{#if square.length > 1}
					<img
						draggable="false"
						class="piece-under"
						src={GetImage2(square.length - 1, square[square.length - 2])}
						alt=""
					/>
				{/if}
			{/if}
		</div>
	{/each}

	<div class="file">
		{#each reverseList(fileCoords) as file}
			<div class="">{file}</div>
		{/each}
	</div>

	<div class="rank">
		{#each rankCoords as rank}
			<div class="">{rank}</div>
		{/each}
	</div>
</div>

<style>
	img {
		/* box-shadow: 0px 7px 15px rgba(230, 106, 5, 0.527);
		border-radius: 50%; */
	}

	.board {
		box-shadow: 0px 7px 50px 5px rgba(230, 106, 5, 0.25);
		display: grid;
		grid-template-columns: repeat(9, minmax(20px, 1fr));
		grid-template-rows: repeat(9, minmax(20px, 1fr));
		gap: 2px;
		max-width: 45rem;
		aspect-ratio: 1/1;
		background-color: rgb(226, 147, 49);
		padding: 2px;
		margin: 2rem;
		position: relative;
	}

	.square {
		background-color: rgb(254 215 170);
		position: relative;
	}

	.file {
		display: grid;
		position: absolute;
		height: 100%;
		margin-left: -1rem;
		align-items: center;
	}

	.rank {
		display: grid;
		position: absolute;
		grid-auto-flow: column;
		width: 100%;
		bottom: -1.5rem;
		text-align: center;
	}

	.piece {
		padding: 0.375rem;
		position: relative;
		z-index: 2;
		user-select: none;
	}
	.piece-under {
		padding: 0.375rem;
		position: absolute;
		left: 0;
		top: 0;
		right: 0;
		z-index: 1;
		user-select: none;
	}
	.piece:hover {
		background-color: rgba(255, 77, 7, 0.479);
	}
</style>

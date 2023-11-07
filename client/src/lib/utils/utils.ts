export function FenToBoard(fen: string): number[][] {
	const newBoard: number[][] = new Array(81).fill([]);
	const fields = fen.split(' ');
	const split = fields[0].split('/');
	if (split.length != 9) {
		//err
	}

	split.forEach((row, index) => {
		const split2 = row.split(',');
		let fileIndex = 0;
		split2.forEach((column) => {
			const newStack: number[] = [];
			const decodeStack = column.split('');
			decodeStack.forEach((piece) => {
				const charCode = piece.charCodeAt(0);
				if (charCode >= '1'.charCodeAt(0) && charCode <= '9'.charCodeAt(0)) {
					const skipNumber = Number(piece);
					fileIndex += skipNumber - 1;
				} else {
					newStack.push(EncodePiece(piece));
				}
			});
			if (newStack.length > 0) {
				newBoard[CoordsToIndex(fileIndex, index)] = newStack;
			}
			fileIndex += 1;
		});
	});
	return newBoard;
}

export function FenToHand(fen: string): number[][] {
	const fields = fen.split(' ')[1];
	const split = fields.split('/');

	const hands: number[][] = [[], []];

	for (const x of split[0]) {
		hands[0].push(Number(x));
	}
	for (const x of split[1]) {
		hands[1].push(Number(x));
	}
	return hands;
}

export function GetImage(stack: number[]): string {
	const topPiece = GetTopStack(stack);
	const encodedPiece = DecodePiece(topPiece).toLowerCase();
	const color = GetPieceColor(topPiece);
	return `/pieces/${color}${stack.length}${encodedPiece}.svg`;
}

export function GetPieceColor(piece: number): string {
	if (piece < 13) {
		return 'w';
	} else {
		return 'b';
	}
}

export function PieceIsPlayerColor(piece: number, playerColor: string): boolean {
	return GetPieceColor(piece) === playerColor;
}

export function GetTopStack(stack: number[]): number {
	return stack[stack.length - 1];
}

export function CoordsToIndex(file: number, rank: number): number {
	return file + rank * 9;
}

export function IndexToCoords(index: number): string[] {
	const file = (index % 9) + 1;
	const rank = Math.floor(index / 9) + 1;
	return [FileToLetter(file), String(RankInvert(rank))];
}

export function RankInvert(num: number): number {
	return 10 - num;
}

export function FileToLetter(num: number): string {
	const pieceEnums = {
		1: 'a',
		2: 'b',
		3: 'c',
		4: 'd',
		5: 'e',
		6: 'f',
		7: 'g',
		8: 'h',
		9: 'i',
	};
	return pieceEnums[num];
}

type EncodePieceEnums = {
	[key: string]: number;
};

export function EncodePiece(decodedPiece: string): number {
	const pieceEnums: EncodePieceEnums = {
		P: 0,
		L: 1,
		S: 2,
		G: 3,
		F: 4,
		K: 5,
		Y: 6,
		B: 7,
		W: 8,
		C: 9,
		N: 10,
		T: 11,
		M: 12,
		p: 13,
		l: 14,
		s: 15,
		g: 16,
		f: 17,
		k: 18,
		y: 19,
		b: 20,
		w: 21,
		c: 22,
		n: 23,
		t: 24,
		m: 25,
	};

	return pieceEnums[decodedPiece];
}

type DecodePieceEnums = {
	[key: string]: string;
};

export function DecodePiece(encodedPiece: number): string {
	if (encodedPiece > 25 || encodedPiece < 0) {
		return '';
	}
	const pieceEnums: DecodePieceEnums = {
		0: 'P',
		1: 'L',
		2: 'S',
		3: 'G',
		4: 'F',
		5: 'K',
		6: 'Y',
		7: 'B',
		8: 'W',
		9: 'C',
		10: 'N',
		11: 'T',
		12: 'M',
		13: 'p',
		14: 'l',
		15: 's',
		16: 'g',
		17: 'f',
		18: 'k',
		19: 'y',
		20: 'b',
		21: 'w',
		22: 'c',
		23: 'n',
		24: 't',
		25: 'm',
	};

	return pieceEnums[encodedPiece];
}

export function DecodePieceFull(encodedPiece: number | string): string {
	const piece = Number(encodedPiece) % 13;
	const pieceEnums: DecodePieceEnums = {
		0: 'Pawn',
		1: 'Lieutenant General',
		2: 'Major General',
		3: 'General',
		4: 'Fortress',
		5: 'Knight',
		6: 'Archer',
		7: 'Musketeer',
		8: 'Samurai',
		9: 'Cannon',
		10: 'Spy',
		11: 'Tactician',
		12: 'Marshal',
	};

	return pieceEnums[piece];
}

import { DecodePiece } from './utils';

// example fen:

// 9/m,1,c,l,n,4/3,pt,g,3,P/fk,pl,spn,s,p,S,p,pG,1/5,l,3/P,P,P,3,S,S,1/3,P,L,PL,P,L,1/p,2,L,1,S,M,2/9 0005222122100/1025012121000 b 11
export function parseFen(fen: string) {
	const fields = fen.split(' ');
	if (fields.length !== 4) {
		throw new Error('Invalid FEN');
	}

	const pieces = fields[0];

	const hands = fields[1].split('/');
	const player1HandList = [];
	let tempHand = hands[0];
	for (let i = 0; i < tempHand.length; i++) {
		player1HandList.push(Number(tempHand[i]));
	}
	const player2HandList = [];
	tempHand = hands[1];
	for (let i = 0; i < tempHand.length; i++) {
		player2HandList.push(Number(tempHand[i]));
	}

	const turnColor = fields[2];
	const isPlayer1Ready = fields[3][0] === '1';
	const isPlayer2Ready = fields[3][1] === '1';

	return { pieces, player1HandList, player2HandList, turnColor, isPlayer1Ready, isPlayer2Ready };
}

export function serializeFen(
	boardArray: number[][],
	player1HandList: number[],
	player2HandList: number[],
	turnColor: string,
	isPlayer1Ready: boolean,
	isPlayer2Ready: boolean
) {
  const player1Ready = isPlayer1Ready ? '1' : '0';
  const player2Ready = isPlayer2Ready ? '1' : '0';
	return `${boardToString(boardArray)} ${player1HandList.join('')}/${player2HandList.join('')} ${turnColor} ${player1Ready}${player2Ready}`;
}

function boardToString(boardArray: number[][]) {
	let boardString = '';
	let skipIndex = 0;
	for (let i = 0; i < boardArray.length; i++) {
		const square = boardArray[i];

		if (square.length === 0) {
			skipIndex++;
		} else {
			let stackStr = '';

			for (let j = 0; j < square.length; j++) {
				stackStr += DecodePiece(square[j]);
			}
			if (skipIndex != i % 9 && skipIndex != 0) {
				boardString += ',';
			}
			if (skipIndex != 0) {
				boardString += skipIndex;
				skipIndex = 0;
			}
			if (i % 9 != 0 || skipIndex != 0) {
				boardString += ',';
			}

			boardString += stackStr;
		}

		if (i % 9 == 8) {
			if (skipIndex != 0) {
				if (skipIndex != 9) {
					boardString += ',';
				}
				boardString += skipIndex;
				skipIndex = 0;
			}
			if (i != 80) {
				boardString += '/';
			}
		}
	}

	return boardString;
}

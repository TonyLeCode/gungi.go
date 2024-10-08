import { CoordsToSquare, EncodePiece, LetterToFile } from './utils';

export function getSquareCoords(str: string): number[] {
	if (!str) return [];
	if (str === 'w-r' || str === 'b-r') return [];

	if (str.includes('x') && str.length === 7) {
		// attack
		const fromCoord = CoordsToSquare(LetterToFile(str[2]), Number(str[1]));
		const toCoord = CoordsToSquare(LetterToFile(str[6]), Number(str[5]));
		return [fromCoord, toCoord];
	} else if (str.includes('-') && str.length === 6) {
		// stack/move
		const fromCoord = CoordsToSquare(LetterToFile(str[2]), Number(str[1]));
		const toCoord = CoordsToSquare(LetterToFile(str[5]), Number(str[4]));
		return [fromCoord, toCoord];
	} else if (str.length === 3) {
		// place
		const toCoord = CoordsToSquare(LetterToFile(str[2]), Number(str[1]));
		return [toCoord];
	}
	return [];
}

// Example Moves
// m8E-7D Move and Stack
// L3F Place
// w-r Ready
// L2Exf3E Attack

// Move is also stack in move history
export function parseMove(
	str: string
):
	| { moveType: 'move'; fromCoord: number; toCoord: number; piece: number }
	| { moveType: 'attack'; fromCoord: number; toCoord: number; fromPiece: number; toPiece: number }
	| { moveType: 'place'; piece: number; toCoord: number }
	| { moveType: 'ready'; playerColor: 'w' | 'b' } {
	const attackRegex = /[PLSGFKYBWCNTMplsgfkybwcntm][1-9][A-I]x[PLSGFKYBWCNTMplsgfkybwcntm][1-9][A-I]/;
	const moveStackRegex = /[PLSGFKYBWCNTMplsgfkybwcntm][1-9][A-I]-[1-9][A-I]/;
	const placeRegex = /^[PLSGFKYBWCNTMplsgfkybwcntm][1-9][A-I]$/;
	const readyRegex = /[wb]-r/;

	const attackMatch = str.match(attackRegex);
	if (attackMatch) {
		const str = attackMatch[0];
		const moveType = 'attack';
		const fromPiece = EncodePiece(str[0]);
		const toPiece = EncodePiece(str[4]);

		const fromRank = Number(str[1]);
		const fromFile = LetterToFile(str[2]);
		const fromCoord = CoordsToSquare(fromFile, fromRank);

		const toRank = Number(str[5]);
		const toFile = LetterToFile(str[6]);
		const toCoord = CoordsToSquare(toFile, toRank);
		return {
			moveType,
			fromCoord,
			toCoord,
			fromPiece,
			toPiece,
		};
	}

	const moveStackMatch = str.match(moveStackRegex);
	if (moveStackMatch) {
		const str = moveStackMatch[0];
		const moveType = 'move';
		const piece = EncodePiece(str[0]);

		const fromRank = Number(str[1]);
		const fromFile = LetterToFile(str[2]);
		const fromCoord = CoordsToSquare(fromFile, fromRank);

		const toRank = Number(str[4]);
		const toFile = LetterToFile(str[5]);
		const toCoord = CoordsToSquare(toFile, toRank);
		return {
			moveType,
			fromCoord,
			toCoord,
			piece,
		};
	}

	const placeMatch = str.match(placeRegex);
	if (placeMatch) {
		const str = placeMatch[0];
		const moveType = 'place';
		const piece = EncodePiece(str[0]);

		const toRank = Number(str[1]);
		const toFile = LetterToFile(str[2]);
		const toCoord = CoordsToSquare(toFile, toRank);
		return {
			moveType,
			toCoord,
			piece,
		};
	}

	const readyMatch = str.match(readyRegex);
	if (readyMatch) {
		const str = readyMatch[0];
		const moveType = 'ready';

		if (str[0] !== 'w' && str[0] !== 'b') throw new Error('Invalid Move');

		const playerColor = str[0];
		return {
			moveType,
			playerColor,
		};
	}

	throw new Error('Invalid Move');
}

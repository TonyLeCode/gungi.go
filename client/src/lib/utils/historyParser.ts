import { CoordsToSquare, LetterToFile } from "./utils";

export function getSquareCoords(str: string): number[]{
  if (!str) return [];
  if (str === "w-r" || str === "b-r") return [];
  
  if (str.includes("x") && str.length === 7) {
    // attack
    const fromCoord = CoordsToSquare(LetterToFile(str[2]), Number(str[1]))
    const toCoord = CoordsToSquare(LetterToFile(str[6]), Number(str[5]))
    return [fromCoord, toCoord]
  } else if (str.includes("-") && str.length === 6) {
    // stack/move
    const fromCoord = CoordsToSquare(LetterToFile(str[2]), Number(str[1]))
    const toCoord = CoordsToSquare(LetterToFile(str[5]), Number(str[4]))
    return [fromCoord, toCoord]
  } else if (str.length === 3) {
    // place
    const toCoord = CoordsToSquare(LetterToFile(str[2]), Number(str[1]))
    return [toCoord]
  }
  return []
}

export function parseMove(str: string) {
  if (!str) return {};
  if (str === "w-r" || str === "b-r") return {};
  
  if (str.includes("x") && str.length === 7) {
    // attack
    const fromCoord = CoordsToSquare(LetterToFile(str[2]), Number(str[1]))
    const toCoord = CoordsToSquare(LetterToFile(str[6]), Number(str[5]))
    return [fromCoord, toCoord]
  } else if (str.includes("-") && str.length === 6) {
    // stack/move
    const fromCoord = CoordsToSquare(LetterToFile(str[2]), Number(str[1]))
    const toCoord = CoordsToSquare(LetterToFile(str[5]), Number(str[4]))
    return [fromCoord, toCoord]
  } else if (str.length === 3) {
    // place
    const toCoord = CoordsToSquare(LetterToFile(str[2]), Number(str[1]))
    return [toCoord]
  }
}
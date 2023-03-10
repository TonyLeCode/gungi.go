function reverseList<T>(arr: T[]): T[] {
	const reversed: T[] = [];
	for (let i = arr.length - 1; i >= 0; i--) {
		reversed.push(arr[i]);
	}
	return reversed;
}

export {reverseList}
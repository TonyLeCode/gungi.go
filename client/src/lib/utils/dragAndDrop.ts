export interface DragAndDropType {
	initialMouseX: number;
	initialMouseY: number;
	offsetX: number;
	offsetY: number;
	dragElement: HTMLElement | null;
	callback: DragAndDropCallback;
	hoverItem: unknown;
	mouseLeave: () => void;
	mouseOver: (item: unknown) => void;
	dragMouse: (e: MouseEvent) => void;
	releaseElement: () => void;
	startDragMouse: (e: MouseEvent) => void;
	dragAndDropAction: (node: HTMLElement) => void;
}
type DragAndDropCallback = ((item: unknown) => void) | null;

export function dragAndDrop() {
	const dragAndDropObj: DragAndDropType = {
		initialMouseX: 0,
		initialMouseY: 0,
		offsetX: 0,
		offsetY: 0,
		dragElement: null,
		callback: null,
		hoverItem: null,
		mouseLeave: function () {
			if (dragAndDropObj.dragElement) {
				dragAndDropObj.hoverItem = null;
			}
		},
		mouseOver: function (item: unknown) {
			if (dragAndDropObj.dragElement) {
				dragAndDropObj.hoverItem = item;
			}
		},

		dragMouse: function (e: MouseEvent) {
			// console.log(dragAndDropObj.hoverItem)
			if (dragAndDropObj.dragElement) {
				const dx = e.clientX + dragAndDropObj.offsetX - dragAndDropObj.initialMouseX;
				const dy = e.clientY + dragAndDropObj.offsetY - dragAndDropObj.initialMouseY;
				dragAndDropObj.dragElement.style.left = `${dx}px`;
				dragAndDropObj.dragElement.style.top = `${dy}px`;
			}
		},

		releaseElement: function () {
			if (dragAndDropObj.dragElement) {
				dragAndDropObj.dragElement.style.left = '0';
				dragAndDropObj.dragElement.style.top = '0';
				document.removeEventListener('mousemove', this.dragMouse);
				document.removeEventListener('mouseup', this.releaseElement);
				dragAndDropObj.dragElement.style.pointerEvents = 'auto';
				// if (hoveringItemIndex != null) {
				// 	data[hoveringItemIndex].push(data[draggingIndex].pop());
				// 	data = data;
				// }
				if (dragAndDropObj.callback && dragAndDropObj.hoverItem !== null) {
					dragAndDropObj.callback(dragAndDropObj.hoverItem);
					// dragAndDropObj.callback();
				}
				dragAndDropObj.dragElement.style.zIndex = '2';
				// console.log(data)
				dragAndDropObj.dragElement = null;
				dragAndDropObj.hoverItem = null;
			}
		},

		startDragMouse: function (e: MouseEvent) {
			dragAndDropObj.initialMouseX = e.clientX;
			dragAndDropObj.initialMouseY = e.clientY;
			const target = e.target as HTMLElement;
			dragAndDropObj.offsetX = e.offsetX - target.offsetWidth / 2;
			dragAndDropObj.offsetY = e.offsetY - target.offsetHeight / 2;
			// draggingIndex = itemIndex;
			dragAndDropObj.dragElement = target;
			dragAndDropObj.dragElement.style.pointerEvents = 'none';
			dragAndDropObj.dragElement.style.zIndex = '3';
			document.addEventListener('mousemove', this.dragMouse);
			document.addEventListener('mouseup', this.releaseElement);
		},

		dragAndDropAction: function (node: HTMLElement) {
			node.addEventListener('mousedown', this.startDragMouse);
		},
	};
	dragAndDropObj.mouseLeave = dragAndDropObj.mouseLeave.bind(dragAndDropObj);
	dragAndDropObj.mouseOver = dragAndDropObj.mouseOver.bind(dragAndDropObj);
	dragAndDropObj.dragMouse = dragAndDropObj.dragMouse.bind(dragAndDropObj);
	dragAndDropObj.releaseElement = dragAndDropObj.releaseElement.bind(dragAndDropObj);
	dragAndDropObj.startDragMouse = dragAndDropObj.startDragMouse.bind(dragAndDropObj);
	dragAndDropObj.dragAndDropAction = dragAndDropObj.dragAndDropAction.bind(dragAndDropObj);

	return dragAndDropObj;
}

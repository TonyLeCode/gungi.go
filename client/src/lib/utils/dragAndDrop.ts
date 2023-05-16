interface DragAndDropType {
	initialMouseX: number;
	initialMouseY: number;
	offsetX: number;
	offsetY: number;
	dragObj: unknown;
	dragElement: HTMLElement | null;
  callback: DragAndDropCallback;
  hoverItem: unknown;
}
type DragAndDropCallback = ((item: unknown) => void) | null

export function dragAndDrop() {
	const dragAndDropObj: DragAndDropType = {
		initialMouseX: 0,
		initialMouseY: 0,
    offsetX: 0,
    offsetY: 0,
		dragObj: null,
		dragElement: null,
    callback: null,
    hoverItem: null
	}

  function mouseLeave(){
    if (dragAndDropObj.dragElement) {
      dragAndDropObj.hoverItem = null
    }
  }
  function mouseOver(item: unknown){
    if (dragAndDropObj.dragElement) {
      dragAndDropObj.hoverItem = item
    }
  }

  function dragMouse(e: MouseEvent) {
    // console.log(dragAndDropObj.hoverItem)
    if (dragAndDropObj.dragElement) {
      const dx = e.clientX + dragAndDropObj.offsetX - dragAndDropObj.initialMouseX ;
      const dy = e.clientY + dragAndDropObj.offsetY - dragAndDropObj.initialMouseY ;
      dragAndDropObj.dragElement.style.left = `${dx}px`;
      dragAndDropObj.dragElement.style.top = `${dy}px`;
    }
  }
  
  function releaseElement() {
    if (dragAndDropObj.dragElement) {
      dragAndDropObj.dragElement.style.left = '0';
      dragAndDropObj.dragElement.style.top = '0';
      document.removeEventListener('mousemove', dragMouse);
      document.removeEventListener('mouseup', releaseElement);
      dragAndDropObj.dragElement.style.pointerEvents = 'auto';
      // if (hoveringItemIndex != null) {
      // 	data[hoveringItemIndex].push(data[draggingIndex].pop());
      // 	data = data;
      // }
      if(dragAndDropObj.callback && dragAndDropObj.hoverItem !== null){
        dragAndDropObj.callback(dragAndDropObj.hoverItem);
        // dragAndDropObj.callback();
      }
      dragAndDropObj.dragElement.style.zIndex = '2';
      // console.log(data)
      dragAndDropObj.dragElement = null;
    }
  }
  
  function startDragMouse(e: MouseEvent) {
    dragAndDropObj.initialMouseX = e.clientX;
    dragAndDropObj.initialMouseY = e.clientY;
    const target = e.target as HTMLElement;
    dragAndDropObj.offsetX = e.offsetX - target.offsetWidth / 2
    dragAndDropObj.offsetY = e.offsetY - target.offsetHeight / 2
    // draggingIndex = itemIndex;
    dragAndDropObj.dragElement = target;
    dragAndDropObj.dragElement.style.pointerEvents = 'none';
    dragAndDropObj.dragElement.style.zIndex = '10';
    document.addEventListener('mousemove', dragMouse);
    document.addEventListener('mouseup', releaseElement);
  }
  
  function dragAndDropAction(node: HTMLElement){
    node.addEventListener('mousedown', startDragMouse)
  }
  return [dragAndDropObj, dragAndDropAction, mouseOver, mouseLeave]
}

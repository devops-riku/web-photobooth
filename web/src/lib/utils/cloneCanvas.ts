export function cloneCanvas(source: HTMLCanvasElement) {
  const c = document.createElement('canvas');
  c.width = source.width;
  c.height = source.height;
  c.getContext('2d')!.drawImage(source, 0, 0);
  return c;
}

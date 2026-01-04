export async function composeStrip(
  images: string[],
  width: number,
  count: number
): Promise<HTMLCanvasElement> {
  const photoHeight = Math.floor((width * 3) / 2);
  const canvas = document.createElement('canvas');

  canvas.width = width;
  canvas.height = photoHeight * count;

  const ctx = canvas.getContext('2d')!;

  const loaded = await Promise.all(
    images.map(src =>
      new Promise<HTMLImageElement>((resolve) => {
        const img = new Image();
        img.src = src;
        img.onload = () => resolve(img);
      })
    )
  );

  loaded.forEach((img, i) => {
    ctx.drawImage(img, 0, i * photoHeight, width, photoHeight);
  });

  return canvas;
}

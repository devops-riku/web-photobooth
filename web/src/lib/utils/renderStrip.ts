import { computeStripLayout } from './stripLayout';
import { applyGLFXFilter } from './glfxFilters';

export async function renderStripCanvas(
    images: string[],
    photoCount: 2 | 3 | 4,
    dpi = 300,
    filter = 'none',
    topIn = 0.6,
    bottomIn = 0.8,
    bgColor = '#ffffff'
): Promise<HTMLCanvasElement> {
    const layout = computeStripLayout(photoCount, dpi, topIn, bottomIn);

    const canvas = document.createElement('canvas');
    canvas.width = layout.canvasWidthPx;
    canvas.height = layout.canvasHeightPx;

    const ctx = canvas.getContext('2d')!;
    ctx.fillStyle = bgColor;
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    const loadedImages = await Promise.all(
        images.slice(0, photoCount).map(
            (src) =>
                new Promise<HTMLImageElement>((resolve, reject) => {
                    const img = new Image();
                    if (!src.startsWith('data:')) {
                        img.crossOrigin = 'anonymous';
                    }
                    img.onload = () => resolve(img);
                    img.onerror = () => reject(new Error(`Failed to load image: ${src.substring(0, 20)}...`));
                    img.src = src;
                })
        )
    );

    let y = layout.topCanvasPx;

    for (let i = 0; i < loadedImages.length; i++) {
        const img = loadedImages[i];

        // 1. Create a temporary canvas for the square photo
        const photoCanvas = document.createElement('canvas');
        photoCanvas.width = layout.contentWidthPx;
        photoCanvas.height = layout.photoHeightPx;
        const pCtx = photoCanvas.getContext('2d')!;

        // 2. Center-crop to square
        const minDim = Math.min(img.width, img.height);
        const sx = (img.width - minDim) / 2;
        const sy = (img.height - minDim) / 2;

        pCtx.drawImage(
            img,
            sx, sy, minDim, minDim, // source
            0, 0, layout.contentWidthPx, layout.photoHeightPx // destination
        );

        // 3. Apply filter ONLY to the photo
        let filteredPhoto: HTMLCanvasElement | HTMLImageElement = photoCanvas;
        if (filter !== 'none') {
            filteredPhoto = await applyGLFXFilter(photoCanvas, filter);
        }

        // 4. Draw onto main strip
        ctx.drawImage(
            filteredPhoto,
            layout.contentX,
            y,
            layout.contentWidthPx,
            layout.photoHeightPx
        );

        y += layout.photoHeightPx + layout.gapPx;
    }

    return canvas;
}



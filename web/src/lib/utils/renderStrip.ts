import { computeStripLayout } from './stripLayout';
import { applyGLFXFilter } from './glfxFilters';

export async function renderStripCanvas(
    images: string[],
    photoCount: 2 | 3 | 4,
    dpi = 300,
    filter = 'none',
    topIn = 0.6,
    bottomIn = 0.8,
    bgColor = '#ffffff',
    roundedCorners = true
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

        // 4. Draw onto main strip with optional rounded corners
        ctx.save();

        const px = layout.contentX;
        const py = y;
        const pw = layout.contentWidthPx;
        const ph = layout.photoHeightPx;

        if (roundedCorners) {
            const radius = Math.floor(dpi * 0.05); // Subtle rounding
            ctx.beginPath();
            ctx.moveTo(px + radius, py);
            ctx.lineTo(px + pw - radius, py);
            ctx.quadraticCurveTo(px + pw, py, px + pw, py + radius);
            ctx.lineTo(px + pw, py + ph - radius);
            ctx.quadraticCurveTo(px + pw, py + ph, px + pw - radius, py + ph);
            ctx.lineTo(px + radius, py + ph);
            ctx.quadraticCurveTo(px, py + ph, px, py + ph - radius);
            ctx.lineTo(px, py + radius);
            ctx.quadraticCurveTo(px, py, px + radius, py);
            ctx.closePath();
            ctx.clip();
        }

        ctx.drawImage(
            filteredPhoto,
            px,
            py,
            pw,
            ph
        );
        ctx.restore();

        y += layout.photoHeightPx + layout.gapPx;
    }

    return canvas;
}



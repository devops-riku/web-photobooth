const STRIP_WIDTH_IN = 2;
const GAP_IN = 0.1;
const SIDE_MARGIN_IN = 0.1;
const CONTENT_WIDTH_IN = STRIP_WIDTH_IN - SIDE_MARGIN_IN * 2;

function inchToPx(inch: number, dpi: number) {
    return Math.round(inch * dpi);
}

export type LayoutResult = {
    canvasWidthPx: number;
    canvasHeightPx: number;
    contentX: number;
    contentWidthPx: number;
    photoHeightPx: number;
    topCanvasPx: number;
    bottomCanvasPx: number;
    gapPx: number;
};

export function computeStripLayout(
    photoCount: 2 | 3 | 4,
    dpi = 300,
    topIn = 0.6,
    bottomIn = 0.8
): LayoutResult {
    // Each photo is square: height = contentWidth
    const photoHeightIn = CONTENT_WIDTH_IN;

    // Total height = top margin + (photos + gaps) + bottom margin
    const totalHeightIn = topIn + (photoHeightIn * photoCount) + (GAP_IN * (photoCount - 1)) + bottomIn;

    return {
        canvasWidthPx: inchToPx(STRIP_WIDTH_IN, dpi),
        canvasHeightPx: inchToPx(totalHeightIn, dpi),

        contentX: inchToPx(SIDE_MARGIN_IN, dpi),
        contentWidthPx: inchToPx(CONTENT_WIDTH_IN, dpi),

        photoHeightPx: inchToPx(photoHeightIn, dpi),
        topCanvasPx: inchToPx(topIn, dpi),
        bottomCanvasPx: inchToPx(bottomIn, dpi),
        gapPx: inchToPx(GAP_IN, dpi)
    };
}

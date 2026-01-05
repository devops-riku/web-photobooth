let fx: any = null;
let glCanvas: any = null;

export async function applyGLFXFilter(
    sourceCanvas: HTMLCanvasElement,
    filter: string
): Promise<HTMLCanvasElement> {
    if (!fx) {
        const module = await import('glfx');
        fx = module.default || module;
    }

    if (!glCanvas) {
        glCanvas = fx.canvas();
    }

    const texture = glCanvas.texture(sourceCanvas);

    glCanvas.draw(texture);

    switch (filter) {
        case 'cinematic':
            glCanvas
                .brightnessContrast(0.0, 0.15)
                .vignette(0.4, 0.6)
                .hueSaturation(0.0, -0.15);
            break;

        case 'film':
            glCanvas
                .sepia(0.3)
                .noise(0.08)
                .brightnessContrast(-0.05, 0.12);
            break;

        case 'warm':
            glCanvas
                .hueSaturation(0.05, 0.1)
                .brightnessContrast(0.05, 0.1);
            break;

        case 'bw':
            glCanvas
                .hueSaturation(0, -1)
                .brightnessContrast(0.0, 0.2);
            break;

        default:
            // no filter
            break;
    }

    glCanvas.update();

    const output = document.createElement('canvas');
    output.width = sourceCanvas.width;
    output.height = sourceCanvas.height;
    output.getContext('2d')!.drawImage(glCanvas, 0, 0);

    texture.destroy();

    return output;
}

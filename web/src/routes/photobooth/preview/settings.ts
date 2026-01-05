export const PREVIEW_SETTINGS = {
    DEFAULT_FILTER: 'none',
    DEFAULT_CAPTION: '',
    DEFAULT_CAPTION_SIZE: 60,
    DEFAULT_FONT: 'Lobster',
    STRIP_THUMBNAIL_WIDTH: 300,
    AVAILABLE_FILTERS: ['none', 'cinematic', 'film', 'warm', 'bw'],
    FONTS: [
        { id: 'Lobster', label: 'Lobster' },
        { id: 'Pacifico', label: 'Pacifico' },
        { id: 'Caveat', label: 'Caveat' },
        { id: 'Dancing Script', label: 'Dancing' },
        { id: 'Bebas Neue', label: 'Bebas' },
        { id: 'Righteous', label: 'Righteous' },
        { id: 'Abril Fatface', label: 'Abril' },
        { id: 'Cormorant Garamond', label: 'Classic' },
        { id: 'Permanent Marker', label: 'Marker' },
        { id: 'Special Elite', label: 'Typewriter' },
        { id: 'Monoton', label: 'Retro' },
        { id: 'Montserrat', label: 'Clean' }
    ],
    // 📏 Layout & Spacing (Pixel-perfect control)
    BRAND_GAP_PX: 50,             // Space between logo and qr
    BRAND_SIDE_PADDING_PX: 0,     // Horizontal padding for branding (makes it narrower than photos)
    BRAND_TOP_PX: 10,             // Padding from top edge to branding
    BRAND_BOT_PX: 40,             // Padding from branding to the first photo

    TIMESTAMP_TOP_PX: 70,         // Padding from the last photo to the timestamp
    TIMESTAMP_BOT_PX: 30,         // Padding from the timestamp to the caption

    CAPTION_BOT_PX: 70,           // Padding from the caption to the bottom edge

    CAPTION_COLOR: '#581c87',     // Dark purple for readability
    TIMESTAMP_COLOR: '#a855f7'    // Lighter purple for date/time
};

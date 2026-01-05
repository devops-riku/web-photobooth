/**
 * Centralized configuration for the web frontend.
 * Backend is always reached via nginx or dev proxy.
 */
export const API_CONFIG = {
    // IMPORTANT: relative base URL
    BASE_URL: '',

    ENDPOINTS: {
        LOGIN: '/api/auth/login',
        SIGNUP: '/api/auth/signup',
        SAVE_STRIP: '/api/strips/save',
        GUEST_SAVE: '/api/strips/guest-save',
        PUBLIC_STRIP: '/api/strips/public/', // + id
        GET_STRIPS: '/api/strips/my-strips',
        STRIP_DETAIL: '/api/strips/', // + id
        SYNC_USER: '/api/auth/sync'
    },

    // App URL (optional, used only for redirects)
    APP_URL: typeof window !== 'undefined' ? window.location.origin : '',

    // CDN URL (this one CAN stay env-based)
    DO_SPACES_URL:
        import.meta.env.VITE_DO_SPACES_URL ||
        'https://innotekinc.sgp1.cdn.digitaloceanspaces.com'
};

/**
 * Global brand configuration to avoid hardcoded strings/paths.
 */
export const BRAND_CONFIG = {
    NAME: 'Wuby',
    LOGO: '/wuby_logo.png',
    MASCOT: '/booth_mascot_pastel.png',
    TAGLINE: 'Digital Memories'
};

/**
 * Helper to get the full URL for a specific endpoint.
 */
export function getApiUrl(
    endpoint: keyof typeof API_CONFIG.ENDPOINTS,
    id?: string | number
): string {
    const path = API_CONFIG.ENDPOINTS[endpoint];
    return id ? `${path}${id}` : path;
}

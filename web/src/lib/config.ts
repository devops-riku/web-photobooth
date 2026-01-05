/**
 * Centralized configuration for the web frontend to interact with the Go backend.
 * Update these values when moving to production.
 */
export const API_CONFIG = {
    // Base URL for the Go monolithic backend
    BASE_URL: 'http://localhost:3000',

    // API Endpoints
    ENDPOINTS: {
        LOGIN: '/api/auth/login',
        SIGNUP: '/api/auth/signup',
        SAVE_STRIP: '/api/strips/save',
        GUEST_SAVE: '/api/strips/guest-save',
        PUBLIC_STRIP: '/api/strips/public/', // Base for /:id
        GET_STRIPS: '/api/strips/my-strips',
        STRIP_DETAIL: '/api/strips/', // Base for /:id
        SYNC_USER: '/api/auth/sync'
    },
    APP_URL: 'http://localhost:5173',
    // DO Spaces CDN URL
    DO_SPACES_URL: 'https://innotekinc.sgp1.cdn.digitaloceanspaces.com'
};

/**
 * Helper to get the full URL for a specific endpoint.
 */
export function getApiUrl(endpoint: keyof typeof API_CONFIG.ENDPOINTS, id?: string | number): string {
    const url = `${API_CONFIG.BASE_URL}${API_CONFIG.ENDPOINTS[endpoint]}`;
    return id ? `${url}${id}` : url;
}

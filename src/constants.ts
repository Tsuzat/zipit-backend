export const DB_URL = process.env.DB_URL;
export const PORT = process.env.PORT || 8080;
export const CROSS_ORIGIN = process.env.CROSS_ORIGIN || "*";

// TODO: More Secure way is to generate JWT Secret at runtime
export const ACCESS_TOKEN_SECRET = process.env.ACCESS_TOKEN_SECRET;
export const ACCESS_TOKEN_EXPIRY = process.env.ACCESS_TOKEN_EXPIRY;
export const REFRESH_TOKEN_EXPIRY = process.env.REFRESH_TOKEN_EXPIRY;
export const REFRESH_TOKEN_SECRET = process.env.REFRESH_TOKEN_SECRET;

// Email Related
export const EMAIL_HOST = process.env.EMAIL_HOST;
export const EMAIL_PORT = process.env.EMAIL_PORT;
export const EMAIL_USERNAME = process.env.EMAIL_USERNAME;
export const EMAIL_PASSWORD = process.env.EMAIL_PASSWORD;
export const EMAIL_FROM = process.env.EMAIL_FROM;
export const EMAIL_PROTOCOL = process.env.EMAIL_PROTOCOL;

// Backend URL
export const BACKEND_URL = process.env.BACKEND_URL;
export const FRONTEND_URL = process.env.FRONTEND_URL;

// Verification Token Expiry
export const VERIFICATION_TOKEN_EXPIRY = parseInt(
  process.env.VERIFICATION_TOKEN_EXPIRY || "30"
);

// REDIS
export const REDIS_URL = process.env.REDIS_URL;
export const REDIS_KEY_EXPIRY = parseInt(process.env.REDIS_KEY_EXPIRY || "300");
// Cookie Options
export const COOKIE_OPTIONS = { httpOnly: true, secure: false };

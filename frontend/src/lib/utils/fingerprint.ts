/**
 * Device Fingerprinting for JEDUG
 * Generates a unique-ish hash based on device characteristics
 * Used for no-login-first user identification
 */

import { browser } from "$app/environment";

const FINGERPRINT_KEY = "jedug-device-fp";

/**
 * Generate a simple but robust device fingerprint
 */
async function generateFingerprint(): Promise<string> {
  const components: string[] = [];

  // Screen info
  components.push(`${screen.width}x${screen.height}`);
  components.push(`${screen.colorDepth}`);
  components.push(`${screen.pixelDepth}`);

  // Platform & language
  components.push(navigator.language);
  components.push(navigator.platform);
  components.push(String(navigator.hardwareConcurrency || "unknown"));

  // Timezone
  components.push(Intl.DateTimeFormat().resolvedOptions().timeZone);

  // Canvas fingerprint (subtle difference per device)
  try {
    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d");
    if (ctx) {
      canvas.width = 200;
      canvas.height = 50;
      ctx.textBaseline = "top";
      ctx.font = "14px Arial";
      ctx.fillStyle = "#f60";
      ctx.fillRect(10, 5, 80, 20);
      ctx.fillStyle = "#069";
      ctx.fillText("JEDUG-FP", 2, 15);
      ctx.fillStyle = "rgba(102,204,0,0.7)";
      ctx.fillText("JEDUG-FP", 4, 17);
      components.push(canvas.toDataURL().slice(-50));
    }
  } catch {
    components.push("no-canvas");
  }

  // Hash using SubtleCrypto
  const data = components.join("|");
  const encoder = new TextEncoder();
  const hashBuffer = await crypto.subtle.digest(
    "SHA-256",
    encoder.encode(data),
  );
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.map((b) => b.toString(16).padStart(2, "0")).join("");
}

/**
 * Get or generate device fingerprint
 * Cached in localStorage for consistency
 */
export async function getDeviceFingerprint(): Promise<string> {
  if (!browser) return "server";

  // Check cache
  const cached = localStorage.getItem(FINGERPRINT_KEY);
  if (cached) return cached;

  // Generate new
  const fp = await generateFingerprint();
  localStorage.setItem(FINGERPRINT_KEY, fp);
  return fp;
}

/**
 * Get fingerprint synchronously (returns cached or empty string)
 */
export function getDeviceFingerprintSync(): string {
  if (!browser) return "server";
  return localStorage.getItem(FINGERPRINT_KEY) || "";
}

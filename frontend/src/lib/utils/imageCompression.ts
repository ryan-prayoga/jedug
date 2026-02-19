/**
 * Client-Side Image Compression for JEDUG
 * - Resize to max-width 1080px
 * - Compress to WebP (Quality 80%)
 * - Strip metadata
 * Result: 10MB -> ~150KB
 */
export async function compressImage(
  file: File,
  maxWidth: number = 1080,
  quality: number = 0.8,
): Promise<Blob> {
  return new Promise((resolve, reject) => {
    const img = new Image();
    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d");

    if (!ctx) {
      reject(new Error("Canvas context not available"));
      return;
    }

    img.onload = () => {
      // Calculate new dimensions
      let width = img.width;
      let height = img.height;

      if (width > maxWidth) {
        height = Math.round((height * maxWidth) / width);
        width = maxWidth;
      }

      canvas.width = width;
      canvas.height = height;

      // Draw image (this strips EXIF metadata)
      ctx.drawImage(img, 0, 0, width, height);

      // Convert to WebP
      canvas.toBlob(
        (blob) => {
          if (blob) {
            resolve(blob);
          } else {
            // Fallback to JPEG if WebP not supported
            canvas.toBlob(
              (jpegBlob) => {
                if (jpegBlob) {
                  resolve(jpegBlob);
                } else {
                  reject(new Error("Failed to compress image"));
                }
              },
              "image/jpeg",
              quality,
            );
          }
        },
        "image/webp",
        quality,
      );
    };

    img.onerror = () => reject(new Error("Failed to load image"));

    // Read file as data URL
    const reader = new FileReader();
    reader.onload = (e) => {
      if (e.target?.result) {
        img.src = e.target.result as string;
      }
    };
    reader.onerror = () => reject(new Error("Failed to read file"));
    reader.readAsDataURL(file);
  });
}

/**
 * Compress multiple images
 */
export async function compressImages(files: File[]): Promise<Blob[]> {
  return Promise.all(files.map((f) => compressImage(f)));
}

/**
 * Convert a Blob to a data URL for preview
 */
export function blobToDataURL(blob: Blob): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => resolve(e.target?.result as string);
    reader.onerror = reject;
    reader.readAsDataURL(blob);
  });
}

/**
 * Format file size to human readable
 */
export function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

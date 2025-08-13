/**
 * Converts a hex color string to RGB values
 * @param hex - Hex color string (with or without #, e.g., "#ff0000" or "ff0000")
 * @returns Object with r, g, b values (0-255) or null if invalid
 */
export function hexToRgb(hex: string): { r: number; g: number; b: number, a: number } {
  const cleanHex = hex.replace('#', '');

  if (!/^[0-9a-fA-F]{3}$|^[0-9a-fA-F]{6}$/.test(cleanHex)) {
    return { r: 255, g: 255, b: 255, a: 1 };
  }

  let r: number, g: number, b: number;

  if (cleanHex.length === 3) {
    r = parseInt(cleanHex[0] + cleanHex[0], 16);
    g = parseInt(cleanHex[1] + cleanHex[1], 16);
    b = parseInt(cleanHex[2] + cleanHex[2], 16);
  } else {
    r = parseInt(cleanHex.slice(0, 2), 16);
    g = parseInt(cleanHex.slice(2, 4), 16);
    b = parseInt(cleanHex.slice(4, 6), 16);
  }

  return { r, g, b, a: 1 };
}

/**
 * Converts RGB values to hex color string
 * @param r - Red value (0-255)
 * @param g - Green value (0-255)
 * @param b - Blue value (0-255)
 * @param includeHash - Whether to include # prefix (default: true)
 * @returns Hex color string or null if invalid RGB values
 */
export function rgbToHex(
  r: number,
  g: number,
  b: number,
  includeHash: boolean = true
): string {
  if (
    !Number.isInteger(r) ||
    r < 0 ||
    r > 255 ||
    !Number.isInteger(g) ||
    g < 0 ||
    g > 255 ||
    !Number.isInteger(b) ||
    b < 0 ||
    b > 255
  ) {
    return includeHash ? `#ffffff` : "ffffff";
  }

  const hexR = r.toString(16).padStart(2, '0');
  const hexG = g.toString(16).padStart(2, '0');
  const hexB = b.toString(16).padStart(2, '0');

  const hex = hexR + hexG + hexB;

  return includeHash ? `#${hex}` : hex;
}

/**
 * Converts RGB object to hex color string
 * @param rgb - Object with r, g, b properties
 * @param includeHash - Whether to include # prefix (default: true)
 * @returns Hex color string or null if invalid RGB values
 */
export function rgbObjectToHex(
  rgb: { r: number; g: number; b: number },
  includeHash: boolean = true
): string | null {
  return rgbToHex(rgb.r, rgb.g, rgb.b, includeHash);
}

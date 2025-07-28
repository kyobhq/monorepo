export function generateRandomId(n: number = 8) {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';

  for (let i = 0; i < n; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }

  return result;
}

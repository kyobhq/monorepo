export function print(...args: any) {
  if (import.meta.env.DEV) {
    console.log(...args);
  }
}

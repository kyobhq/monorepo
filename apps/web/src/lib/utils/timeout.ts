import { core } from "stores/core.svelte";

export function setControllableTimeout(callback: () => void, delay: number) {
  const timeoutId = setTimeout(callback, delay);

  return {
    executeNow() {
      clearTimeout(timeoutId)
      callback()
    },
    clear() {
      clearTimeout(timeoutId)
      core.callTimeout = undefined
    }
  }
}

import type { APIError } from "$lib/types/errors";

export function print(...args: any) {
  if (import.meta.env.DEV) {
    console.log(...args);
  }
}

export function logErr(error: APIError) {
  if (import.meta.env.DEV) {
    console.log(`
[API Error]
Status: ${error.status}
Code: ${error.code}
Message: ${error.message}
Cause: ${error.cause}
    `)
  }
}

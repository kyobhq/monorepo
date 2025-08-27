/**
 * Converts an array to an object using a property name as the key
 * Use this instead of a reduce please
 */
export function keyByProperty<T, K extends keyof T>(array: T[], property: K): Record<string, T> {
	const result = {} as Record<string, T>;
	for (const item of array) {
		const key = String(item[property]);
		result[key] = item;
	}
	return result;
}

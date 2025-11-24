import timezones from '$lib/assets/timezones.json';

export interface Timezone {
  label: string;
  tzCode: string;
  name: string;
  utc: string;
}

/**
 * Get all available timezones
 */
export function getTimezones(): Timezone[] {
  return timezones as Timezone[];
}

/**
 * Get a specific timezone by its timezone code
 * @param tzCode - The timezone code (e.g., "America/New_York")
 * @returns The timezone object or undefined if not found
 */
export function getTimezoneByCode(tzCode: string): Timezone | undefined {
  return timezones.find((tz) => tz.tzCode === tzCode);
}

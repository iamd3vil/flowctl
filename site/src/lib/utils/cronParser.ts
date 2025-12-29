import { CronExpressionParser } from 'cron-parser';

/**
 * Get the next run time for a cron expression with a specific timezone
 * @param cron - The cron expression (e.g., "0 0 * * *")
 * @param timezone - The timezone code (e.g., "America/New_York")
 * @returns The next run Date or null if parsing fails
 */
export function getNextCronRun(cron: string, timezone: string): Date | null {
  try {
    const interval = CronExpressionParser.parse(cron, {
      tz: timezone,
      currentDate: new Date()
    });
    return interval.next().toDate();
  } catch (error) {
    console.error('Failed to parse cron expression:', cron, error);
    return null;
  }
}

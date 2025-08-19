/**
 * Creates a slug from a given name by converting to lowercase,
 * removing special characters, and replacing spaces with underscores
 */
export function createSlug(name: string): string {
  return name
    .toLowerCase()
    .replace(/[^a-z0-9\s]/g, '') // Remove special characters except spaces
    .replace(/\s+/g, '_') // Replace spaces with underscores
    .replace(/_{2,}/g, '_') // Replace multiple underscores with single
    .replace(/^_|_$/g, ''); // Remove leading/trailing underscores
}

/**
 * Validates a cron expression format
 * Returns true if valid, false otherwise
 */
export function isValidCronExpression(cronExpr: string): boolean {
  if (!cronExpr || cronExpr.trim() === '') {
    return true; // Empty is valid (means no schedule)
  }

  const parts = cronExpr.trim().split(/\s+/);
  
  // Cron should have 5 parts: minute hour day month weekday
  if (parts.length !== 5) {
    return false;
  }

  // Basic validation for each part
  const validators = [
    (val: string) => isValidCronPart(val, 0, 59),   // minute (0-59)
    (val: string) => isValidCronPart(val, 0, 23),   // hour (0-23)
    (val: string) => isValidCronPart(val, 1, 31),   // day (1-31)
    (val: string) => isValidCronPart(val, 1, 12),   // month (1-12)
    (val: string) => isValidCronPart(val, 0, 7),    // weekday (0-7, both 0 and 7 = Sunday)
  ];

  return parts.every((part, index) => validators[index](part));
}

/**
 * Validates a single cron part (minute, hour, day, etc.)
 */
function isValidCronPart(part: string, min: number, max: number): boolean {
  // Allow wildcards
  if (part === '*') return true;
  
  // Allow step values like */5
  if (part.includes('*/')) {
    const [base, step] = part.split('/');
    if (base !== '*') return false;
    const stepNum = parseInt(step);
    return !isNaN(stepNum) && stepNum > 0 && stepNum <= max;
  }
  
  // Allow ranges like 1-5
  if (part.includes('-')) {
    const [start, end] = part.split('-');
    const startNum = parseInt(start);
    const endNum = parseInt(end);
    return !isNaN(startNum) && !isNaN(endNum) && 
           startNum >= min && endNum <= max && startNum <= endNum;
  }
  
  // Allow lists like 1,3,5
  if (part.includes(',')) {
    const values = part.split(',');
    return values.every(val => {
      const num = parseInt(val.trim());
      return !isNaN(num) && num >= min && num <= max;
    });
  }
  
  // Single number
  const num = parseInt(part);
  return !isNaN(num) && num >= min && num <= max;
}
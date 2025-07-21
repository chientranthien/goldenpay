export const config = {
  host: process.env.REACT_APP_API_URL as string
};

export const ToJSON = (obj: any): string => JSON.stringify(toSnakeCase(obj, true))
export const FromJSON = (j: string): any => toCamelCase(JSON.parse(j), true)

function toSnakeCase(data: any, deep: boolean = true): any {
  if (!data || typeof data !== "object") return data;

  if (Array.isArray(data)) {
    return data.map((item) => toSnakeCase(item, deep));
  }

  return Object.keys(data).reduce((acc: any, key: string) => {
    const newKey = key.replace(/([A-Z])/g, (match, p1) => `_${p1.toLowerCase()}`);
    const value = data[key];

    // Deep conversion for nested objects and arrays
    acc[newKey] = deep && typeof value === "object" ? toSnakeCase(value, true) : value;
    return acc;
  }, {});
}

function toCamelCase(data: any, deep: boolean = true): any {
  if (!data || typeof data !== "object") return data;

  if (Array.isArray(data)) {
    return data.map((item) => toCamelCase(item, deep));
  }

  return Object.keys(data).reduce((acc: any, key: string) => {
    const newKey = key.replace(/_([a-z])/g, (match, p1) => p1.toUpperCase());
    const value = data[key];

    // Deep conversion for nested objects and arrays
    acc[newKey] = deep && typeof value === "object" ? toCamelCase(value, true) : value;
    return acc;
  }, {});
}
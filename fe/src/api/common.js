export const config = {
  host: "http://localhost:5000/api/v1/"
};

export const ToJSON = obj => JSON.stringify(toSnakeCase(obj, true))
export const FromJSON = j => toCamelCase(JSON.parse(j), true)

function toSnakeCase(data, deep = true) {
  if (!data || typeof data !== "object") return data;

  if (Array.isArray(data)) {
    return data.map((item) => toSnakeCase(item, deep));
  }

  return Object.keys(data).reduce((acc, key) => {
    const newKey = key.replace(/([A-Z])/g, (match, p1) => `_${p1.toLowerCase()}`);
    const value = data[key];

    // Deep conversion for nested objects and arrays
    acc[newKey] = deep && typeof value === "object" ? toSnakeCase(value, true) : value;
    return acc;
  }, {});
}

function toCamelCase(data, deep = true) {
  if (!data || typeof data !== "object") return data;

  if (Array.isArray(data)) {
    return data.map((item) => toCamelCase(item, deep));
  }

  return Object.keys(data).reduce((acc, key) => {
    const newKey = key.replace(/_([a-z])/g, (match, p1) => p1.toUpperCase());
    const value = data[key];

    // Deep conversion for nested objects and arrays
    acc[newKey] = deep && typeof value === "object" ? toCamelCase(value, true) : value;
    return acc;
  }, {});
}
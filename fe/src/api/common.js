export const config = {
  host: "http://localhost:5000/api/v1/"
};

const CamelToSnakeCase = str => str.replace(/[A-Z]/g, letter => `_${letter.toLowerCase()}`);
const SnakeToCamel = str =>
  str.toLowerCase().replace(/([-_][a-z])/g, group =>
    group
      .toUpperCase()
      .replace('-', '')
      .replace('_', '')
  );

export const ToJSON = obj => CamelToSnakeCase(JSON.stringify(obj))
export const FromJSON = j => JSON.parse(SnakeToCamel(j))
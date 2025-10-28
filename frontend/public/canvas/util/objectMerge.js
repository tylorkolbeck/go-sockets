export function mergeObject(defaultObj, newObj) {
  const rtnObj = {};

  for (const key in defaultObj) {
    if (newObj[key]) {
      rtnObj[key] = newObj[key];
    } else {
      rtnObj[key] = defaultObj[key];
    }
  }

  return structuredClone(rtnObj);
}

type Environment = {
  MOCKED: boolean;
}

const ENV: Environment = {
  MOCKED: false,
};

// eslint-disable-next-line
const covertToType = (value: any) => {
  if (value === "true") {
    return true;
  }

  if (value === "false") {
    return false;
  }

  if (value === "null") {
    return null;
  }

  if (value === "undefined") {
    return undefined;
  }

  if (!isNaN(value)) {
    return Number(value);
  }

  return value;
};

/* eslint-disable */
for (const key in ENV) {
  if (import.meta.env?.hasOwnProperty(`VITE_${key}`)) {
    (ENV as any)[key] = covertToType(import.meta.env[`VITE_${key}`]);
  }

  if ((window as any).SERVER_DATA?.hasOwnProperty(key)) {
    (ENV as any)[key] = covertToType((window as any).SERVER_DATA[key]);
  }
}

(window as any).environment = ENV;
/* eslint-enable */

console.table(ENV);

export { ENV };

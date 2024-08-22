import axios from "axios";
import AxiosMockAdapter from "axios-mock-adapter";
// import { ENV } from "../../constants/env";
import registerLoggerInterceptor from "../interceptors/logger";

const isMocked = import.meta.env.VITE_MOCKED === "true";

// Mock axios instance
const mockInstance = axios.create();
registerLoggerInterceptor(mockInstance);



export const mock = new AxiosMockAdapter(mockInstance, {
  delayResponse: 200,
});

// Real axios instance
const realInstance = axios.create();

if (isMocked) {
  console.debug("Using mocked axios client");
}

export default isMocked ? mockInstance : realInstance;

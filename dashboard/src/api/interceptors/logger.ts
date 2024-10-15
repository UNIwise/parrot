import { AxiosInstance } from "axios";
import * as AxiosLogger from "axios-logger";

function registerLoggerInterceptor(inst: AxiosInstance) {
  inst.interceptors.request.use((request) => {
    return AxiosLogger.requestLogger(request, {
      prefixText: "Mocked",
      dateFormat: "HH:MM:ss",
      headers: false,
      data: false,
    });
  });
}

export default registerLoggerInterceptor;

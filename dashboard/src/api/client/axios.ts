import axios from 'axios';
import AxiosMockAdapter from 'axios-mock-adapter';
import { ENV } from '../../constants/env';
import registerLoggerInterceptor from '../interceptors/logger';

// Mock axios instance
const mockInstance = axios.create();

registerLoggerInterceptor(mockInstance);

export const mock = new AxiosMockAdapter(mockInstance, {
  delayResponse: 200,
});

// Real axios instance
const realInstance = axios.create();

if (ENV.MOCKED) {
  console.debug('Using mocked axios client');
}

export default ENV.MOCKED ? mockInstance : realInstance;

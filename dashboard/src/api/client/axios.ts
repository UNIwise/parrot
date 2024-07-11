import axios from 'axios';


import registerAuthorizationInterceptor from '../interceptors/authorization';
import registerLoggerInterceptor from '../interceptors/logger';
import { ENV } from '../../constants/env';
import AxiosMockAdapter from 'axios-mock-adapter';

// Mock axios instance
const mockInstance = axios.create();
registerAuthorizationInterceptor(mockInstance);
registerLoggerInterceptor(mockInstance);

export const mock = new AxiosMockAdapter(mockInstance, {
  delayResponse: 200,
});

// Real axios instance
const realInstance = axios.create();
registerAuthorizationInterceptor(realInstance);

if (ENV.MOCKED) {
  console.debug('Using mocked axios client');
}

export default ENV.MOCKED ? mockInstance : realInstance;

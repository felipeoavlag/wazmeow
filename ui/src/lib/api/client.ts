import axios, { AxiosInstance, AxiosResponse } from 'axios';
import { ApiResponse } from '@/lib/types/api';

class ApiClient {
  private client: AxiosInstance;

  constructor(baseURL: string = 'http://localhost:8080') {
    this.client = axios.create({
      baseURL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Interceptor para requests
    this.client.interceptors.request.use(
      (config) => {
        // Adicionar token de autenticação se necessário
        // const token = localStorage.getItem('auth_token');
        // if (token) {
        //   config.headers.Authorization = `Bearer ${token}`;
        // }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Interceptor para responses
    this.client.interceptors.response.use(
      (response: AxiosResponse<ApiResponse>) => {
        return response;
      },
      (error) => {
        // Tratamento global de erros
        if (error.response?.status === 401) {
          // Redirect para login se necessário
          // window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    );
  }

  // Métodos HTTP genéricos
  async get<T = any>(url: string, params?: any): Promise<ApiResponse<T>> {
    const response = await this.client.get<ApiResponse<T>>(url, { params });
    return response.data;
  }

  async post<T = any>(url: string, data?: any): Promise<ApiResponse<T>> {
    const response = await this.client.post<ApiResponse<T>>(url, data);
    return response.data;
  }

  async put<T = any>(url: string, data?: any): Promise<ApiResponse<T>> {
    const response = await this.client.put<ApiResponse<T>>(url, data);
    return response.data;
  }

  async delete<T = any>(url: string): Promise<ApiResponse<T>> {
    const response = await this.client.delete<ApiResponse<T>>(url);
    return response.data;
  }

  // Método para configurar a URL base
  setBaseURL(baseURL: string) {
    this.client.defaults.baseURL = baseURL;
  }

  // Método para configurar timeout
  setTimeout(timeout: number) {
    this.client.defaults.timeout = timeout;
  }
}

// Instância singleton do cliente
export const apiClient = new ApiClient();

// Hook para configurar a API
export const useApiConfig = () => {
  const setBaseURL = (url: string) => {
    apiClient.setBaseURL(url);
  };

  const setTimeout = (timeout: number) => {
    apiClient.setTimeout(timeout);
  };

  return {
    setBaseURL,
    setTimeout,
  };
};
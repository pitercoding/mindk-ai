import { getToken } from "@clerk/react";

const API_URL =
  import.meta.env.VITE_API_URL ??
  "http://localhost:8080";

export async function apiClient<T>(
  endpoint: string,
  options?: RequestInit,
): Promise<T> {

  const token = await getToken();

  const response = await fetch(
    `${API_URL}${endpoint}`,
    {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...options?.headers,
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
    },
  );

  if (!response.ok) {
    throw new Error(
      `API Error: ${response.status}`,
    );
  }

  const text = await response.text();

  if (!text) {
    return undefined as T;
  }

  try {
    return JSON.parse(text);
  } catch {
    throw new Error(
      "Invalid JSON response from API",
    );
  }
}

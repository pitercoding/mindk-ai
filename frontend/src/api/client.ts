const API_URL =
  import.meta.env.VITE_API_URL ??
  "http://localhost:8080";

export async function apiClient<T>(
  endpoint: string,
  options?: RequestInit,
): Promise<T> {

  const response = await fetch(
    `${API_URL}${endpoint}`,
    {
      headers: {
        "Content-Type": "application/json",
      },
      ...options,
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

  return JSON.parse(text);
}

import { beforeEach, describe, expect, it, vi } from "vitest";

const getToken = vi.fn();

vi.mock("@clerk/react", () => ({
    getToken: () => getToken(),
}));

import { API_URL, apiClient } from "@/api/client";

function mockFetch(response: Partial<Response> & { text?: () => Promise<string> }) {
    const fetchMock = vi.fn().mockResolvedValue(response);
    vi.stubGlobal("fetch", fetchMock);
    return fetchMock;
}

describe("apiClient", () => {

    beforeEach(() => {
        getToken.mockReset();
        getToken.mockResolvedValue("token-123");
        vi.unstubAllGlobals();
    });

    it("calls the API with the bearer token and JSON content type", async () => {
        const fetchMock = mockFetch({
            ok: true,
            text: async () => JSON.stringify([{ id: 1 }]),
        });

        const result = await apiClient<{ id: number }[]>("/notes");

        expect(result).toEqual([{ id: 1 }]);
        expect(fetchMock).toHaveBeenCalledTimes(1);

        const [url, options] = fetchMock.mock.calls[0];
        expect(url).toBe(`${API_URL}/notes`);
        expect(options.headers).toMatchObject({
            "Content-Type": "application/json",
            Authorization: "Bearer token-123",
        });
    });

    it("omits the Authorization header when there is no token", async () => {
        getToken.mockResolvedValue(null);
        const fetchMock = mockFetch({ ok: true, text: async () => "{}" });

        await apiClient("/health");

        const [, options] = fetchMock.mock.calls[0];
        expect(options.headers).not.toHaveProperty("Authorization");
    });

    it("lets callers override headers and pass request options", async () => {
        const fetchMock = mockFetch({ ok: true, text: async () => "{}" });

        await apiClient("/notes", {
            method: "POST",
            body: "{}",
            headers: { "X-Test": "1" },
        });

        const [, options] = fetchMock.mock.calls[0];
        expect(options.method).toBe("POST");
        expect(options.headers).toMatchObject({ "X-Test": "1" });
    });

    it("returns undefined for an empty response body", async () => {
        mockFetch({ ok: true, text: async () => "" });

        await expect(apiClient("/notes/1", { method: "DELETE" }))
            .resolves.toBeUndefined();
    });

    it("throws with the status code when the response is not ok", async () => {
        mockFetch({ ok: false, status: 404, text: async () => "" });

        await expect(apiClient("/missing")).rejects.toThrow("API Error: 404");
    });

    it("throws when the body is not valid JSON", async () => {
        mockFetch({ ok: true, text: async () => "<html>oops</html>" });

        await expect(apiClient("/notes")).rejects.toThrow("Invalid JSON response from API");
    });
});

import { apiClient } from "@/api/client";
import type { Document } from "@/types/document";

export async function getDocuments(): Promise<Document[]> {

    const response =
        await apiClient<Document[]>("/documents");

    return response ?? [];
}


export async function deleteDocument(
    id: number
): Promise<void> {

    await apiClient(
        `/documents/${id}`,
        {
            method: "DELETE",
        },
    );
}

export async function searchDocuments(
    query: string
): Promise<Document[]> {

    const response =
        await apiClient<Document[]>(
            `/documents/search?q=${encodeURIComponent(query)}`
        );

    return response ?? [];
}

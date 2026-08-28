import { apiClient } from "@/api/client";

import type { ChatSession } from "@/types/chat-session";


export async function getSessions(): Promise<ChatSession[]> {

    return apiClient<ChatSession[]>("/chat/sessions");
}


export async function updateSession(
    id: number,
    session: Pick<ChatSession, "title" | "mode">,
): Promise<ChatSession> {

    return apiClient<ChatSession>(
        `/chat/sessions/${id}`,
        {
            method: "PUT",
            body: JSON.stringify(session),
        },
    );
}


export async function deleteSession(
    id: number,
): Promise<void> {

    return apiClient<void>(
        `/chat/sessions/${id}`,
        { method: "DELETE" },
    );
}

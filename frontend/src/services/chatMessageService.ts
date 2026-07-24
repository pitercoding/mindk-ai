import { apiClient } from "@/api/client";

import type { ChatMessage } from "@/types/chat";


export async function getMessagesByNote(
    noteId: number,
): Promise<ChatMessage[]> {

    return apiClient<ChatMessage[]>(
        `/chat/messages/${noteId}`,
    );
}


export async function saveMessage(
    message: Omit<ChatMessage, "id" | "created_at">,
): Promise<ChatMessage> {

    return apiClient<ChatMessage>(
        "/chat/messages",
        {
            method: "POST",
            body: JSON.stringify(message),
        },
    );
}
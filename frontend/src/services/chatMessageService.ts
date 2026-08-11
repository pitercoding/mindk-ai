import { apiClient } from "@/api/client";

import type { ChatMessage } from "@/types/chat";


export async function getMessagesBySession(
    sessionId: number,
): Promise<ChatMessage[]> {

    return apiClient<ChatMessage[]>(
        `/chat/messages/${sessionId}`,
    );
}

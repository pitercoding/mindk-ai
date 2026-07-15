import { apiClient } from "@/api/client";
import type { ChatHistory } from "@/types/chat-history";

interface ChatHistoryResponse {
    data: ChatHistory[];
    page: number;
    limit: number;
    total: number;
}

export async function getChatHistory(
    page = 1,
    limit = 10,
): Promise<ChatHistoryResponse> {
    return apiClient<ChatHistoryResponse>(
        `/chat/history?page=${page}&limit=${limit}`,
    );
}
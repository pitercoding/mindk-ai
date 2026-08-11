import { apiClient } from "@/api/client";
import type { ChatRequest, ChatResponse } from "@/types/chat";

export async function ask(request: ChatRequest): Promise<ChatResponse> {

    return apiClient<ChatResponse>("/chat", { method: "POST", body: JSON.stringify(request) });
};

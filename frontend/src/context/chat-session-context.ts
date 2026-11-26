import { createContext, useContext } from "react";

import type { ChatMode, ChatSession } from "@/types/chat-session";
import type { Message } from "@/types/message";

export interface NewSessionOptions {
    mode: ChatMode;
    noteId?: number;
    title?: string;
}

export interface ChatSessionContextType {
    sessions: ChatSession[];
    isLoadingSessions: boolean;
    sessionsError: string | null;
    currentSessionId: number | null;
    currentSession: ChatSession | null;
    messages: Message[];
    isSending: boolean;
    isLoadingMessages: boolean;
    selectSession: (id: number) => Promise<void>;
    newChat: () => void;
    sendMessage: (text: string, options?: NewSessionOptions) => Promise<void>;
    renameSession: (id: number, title: string) => Promise<void>;
    deleteSession: (id: number) => Promise<void>;
}

export const ChatSessionContext = createContext<
    ChatSessionContextType | undefined
>(undefined);

export function useChatSession() {

    const context = useContext(
        ChatSessionContext,
    );

    if (!context) {
        throw new Error(
            "useChatSession must be used inside ChatSessionProvider"
        );
    }

    return context;
}

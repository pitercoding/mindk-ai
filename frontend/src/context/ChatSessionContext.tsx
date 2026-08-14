import {
    createContext,
    useCallback,
    useContext,
    useEffect,
    useState,
    type ReactNode,
} from "react";

import { ask } from "@/services/chatService";
import { getMessagesBySession } from "@/services/chatMessageService";
import {
    deleteSession as deleteSessionRequest,
    getSessions,
    updateSession as updateSessionRequest,
} from "@/services/chatSessionService";

import type { ChatMode, ChatSession } from "@/types/chat-session";
import type { Message } from "@/types/message";

interface NewSessionOptions {
    mode: ChatMode;
    noteId?: number;
    title?: string;
}

interface ChatSessionContextType {
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

const ChatSessionContext = createContext<
    ChatSessionContextType | undefined
>(undefined);

interface ChatSessionProviderProps {
    children: ReactNode;
}

export function ChatSessionProvider({
    children,
}: ChatSessionProviderProps) {

    const [sessions, setSessions] = useState<ChatSession[]>([]);
    const [isLoadingSessions, setIsLoadingSessions] = useState(true);
    const [sessionsError, setSessionsError] = useState<string | null>(null);
    const [currentSessionId, setCurrentSessionId] = useState<number | null>(null);
    const [messages, setMessages] = useState<Message[]>([]);
    const [isSending, setIsSending] = useState(false);
    const [isLoadingMessages, setIsLoadingMessages] = useState(false);

    const refreshSessions = useCallback(async () => {

        setIsLoadingSessions(true);
        setSessionsError(null);

        try {
            const result = await getSessions();
            setSessions(result ?? []);
        } catch (error) {
            console.error("Failed to load chat sessions:", error);
            setSessionsError("Unable to load conversations.");
        } finally {
            setIsLoadingSessions(false);
        }

    }, []);

    useEffect(() => {
        refreshSessions();
    }, [refreshSessions]);

    const currentSession =
        sessions.find(session => session.id === currentSessionId) ?? null;

    async function selectSession(id: number) {

        setCurrentSessionId(id);
        setMessages([]);
        setIsLoadingMessages(true);

        try {

            const history = (await getMessagesBySession(id)) ?? [];

            setMessages(
                history.map(message => ({
                    id: message.id,
                    role: message.role,
                    content: message.content,
                })),
            );

        } catch (error) {

            console.error("Failed to load session messages:", error);
            setMessages([]);

        } finally {
            setIsLoadingMessages(false);
        }
    }

    function newChat() {
        setCurrentSessionId(null);
        setMessages([]);
    }

    async function sendMessage(
        text: string,
        options?: NewSessionOptions,
    ) {

        const userMessage: Message = {
            id: crypto.randomUUID(),
            role: "user",
            content: text,
        };

        setMessages(previous => [...previous, userMessage]);
        setIsSending(true);

        try {

            const response = await ask({
                session_id: currentSessionId ?? 0,
                message: text,
                mode: options?.mode,
                note_id: options?.noteId,
                title: options?.title,
            });

            setMessages(previous => [
                ...previous,
                {
                    id: crypto.randomUUID(),
                    role: "assistant",
                    content: response.answer,
                    sources: response.sources,
                },
            ]);

            if (currentSessionId === null) {
                setCurrentSessionId(response.session_id);
                await refreshSessions();
            }

        } catch (error) {

            console.error("Failed to send chat message:", error);

            setMessages(previous => [
                ...previous,
                {
                    id: crypto.randomUUID(),
                    role: "assistant",
                    content: "Unable to process your request. Please try again.",
                    status: "error",
                },
            ]);

        } finally {
            setIsSending(false);
        }
    }

    async function renameSession(id: number, title: string) {

        const session = sessions.find(item => item.id === id);

        if (!session) {
            return;
        }

        const updated = await updateSessionRequest(id, {
            title,
            mode: session.mode,
        });

        setSessions(previous =>
            previous.map(item => (item.id === id ? updated : item)),
        );
    }

    async function deleteSession(id: number) {

        await deleteSessionRequest(id);

        setSessions(previous =>
            previous.filter(item => item.id !== id),
        );

        if (id === currentSessionId) {
            newChat();
        }
    }

    return (
        <ChatSessionContext.Provider
            value={{
                sessions,
                isLoadingSessions,
                sessionsError,
                currentSessionId,
                currentSession,
                messages,
                isSending,
                isLoadingMessages,
                selectSession,
                newChat,
                sendMessage,
                renameSession,
                deleteSession,
            }}
        >
            {children}
        </ChatSessionContext.Provider>
    );
}

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

import {
    useCallback,
    useEffect,
    useRef,
    useState,
    type ReactNode,
} from "react";
import { useAuth } from "@clerk/react";

import { ask } from "@/services/chatService";
import { getMessagesBySession } from "@/services/chatMessageService";
import {
    deleteSession as deleteSessionRequest,
    getSessions,
    updateSession as updateSessionRequest,
} from "@/services/chatSessionService";

import {
    ChatSessionContext,
    type NewSessionOptions,
} from "@/context/chat-session-context";

import type { ChatSession } from "@/types/chat-session";
import type { Message } from "@/types/message";

interface ChatSessionProviderProps {
    children: ReactNode;
}

export function ChatSessionProvider({
    children,
}: ChatSessionProviderProps) {

    const { isLoaded, isSignedIn } = useAuth();

    const [sessions, setSessions] = useState<ChatSession[]>([]);
    const [isLoadingSessions, setIsLoadingSessions] = useState(true);
    const [sessionsError, setSessionsError] = useState<string | null>(null);
    const [currentSessionId, setCurrentSessionId] = useState<number | null>(null);
    const [messages, setMessages] = useState<Message[]>([]);
    const [isSending, setIsSending] = useState(false);
    const [isLoadingMessages, setIsLoadingMessages] = useState(false);

    // Bumped whenever the active conversation changes (session switch, new chat).
    // In-flight requests started under a previous value are ignored on completion,
    // so a slow response can't overwrite the messages of a session the user has since left.
    const activeContextRef = useRef(0);

    // Fetches the conversations without touching state before the first await,
    // so it is safe to call from an effect (isLoadingSessions starts as true).
    const fetchSessions = useCallback(async () => {

        try {
            const result = await getSessions();
            setSessions(result ?? []);
            setSessionsError(null);
        } catch (error) {
            console.error("Failed to load chat sessions:", error);
            setSessionsError("Unable to load conversations.");
        } finally {
            setIsLoadingSessions(false);
        }

    }, []);

    const refreshSessions = useCallback(async () => {

        setIsLoadingSessions(true);
        setSessionsError(null);

        await fetchSessions();

    }, [fetchSessions]);

    useEffect(() => {
        if (!isLoaded || !isSignedIn) {
            return;
        }
        // Deferred to a microtask: fetchSessions sets state once its request
        // resolves, and react-hooks/set-state-in-effect disallows a direct call.
        Promise.resolve().then(fetchSessions);
    }, [isLoaded, isSignedIn, fetchSessions]);

    const currentSession =
        sessions.find(session => session.id === currentSessionId) ?? null;

    async function selectSession(id: number) {

        const contextId = ++activeContextRef.current;

        setIsSending(false);
        setCurrentSessionId(id);
        setMessages([]);
        setIsLoadingMessages(true);

        try {

            const history = (await getMessagesBySession(id)) ?? [];

            if (activeContextRef.current !== contextId) {
                return;
            }

            setMessages(
                history.map(message => ({
                    id: message.id,
                    role: message.role,
                    content: message.content,
                    createdAt: message.created_at,
                })),
            );

        } catch (error) {

            if (activeContextRef.current !== contextId) {
                return;
            }

            console.error("Failed to load session messages:", error);
            setMessages([]);

        } finally {
            if (activeContextRef.current === contextId) {
                setIsLoadingMessages(false);
            }
        }
    }

    function newChat() {
        ++activeContextRef.current;
        setIsSending(false);
        // A history request still in flight is now stale and will be ignored,
        // so it can no longer clear the loading flag on its own.
        setIsLoadingMessages(false);
        setCurrentSessionId(null);
        setMessages([]);
    }

    async function sendMessage(
        text: string,
        options?: NewSessionOptions,
    ) {

        const contextId = activeContextRef.current;
        const sessionIdAtSend = currentSessionId;

        const userMessage: Message = {
            id: crypto.randomUUID(),
            role: "user",
            content: text,
            createdAt: new Date().toISOString(),
        };

        setMessages(previous => [...previous, userMessage]);
        setIsSending(true);

        try {

            const response = await ask({
                session_id: sessionIdAtSend ?? 0,
                message: text,
                mode: options?.mode,
                note_id: options?.noteId,
                title: options?.title,
            });

            if (activeContextRef.current !== contextId) {
                return;
            }

            setMessages(previous => [
                ...previous,
                {
                    id: crypto.randomUUID(),
                    role: "assistant",
                    content: response.answer,
                    sources: response.sources,
                    createdAt: new Date().toISOString(),
                },
            ]);

            if (sessionIdAtSend === null) {
                setCurrentSessionId(response.session_id);
                await refreshSessions();
            }

        } catch (error) {

            if (activeContextRef.current !== contextId) {
                return;
            }

            console.error("Failed to send chat message:", error);

            setMessages(previous => [
                ...previous,
                {
                    id: crypto.randomUUID(),
                    role: "assistant",
                    content: "Unable to process your request. Please try again.",
                    status: "error",
                    createdAt: new Date().toISOString(),
                },
            ]);

        } finally {
            if (activeContextRef.current === contextId) {
                setIsSending(false);
            }
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

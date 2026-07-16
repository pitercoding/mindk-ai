import { useEffect, useState } from "react";

import ChatInput from "@/components/chat/ChatInput";
import ChatMessageList from "@/components/chat/ChatMessageList";
import ChatHistoryList from "@/components/chat/ChatHistoryList";

import { sendMessage } from "@/services/chatService";
import { getChatHistory } from "@/services/chatHistoryService";

import type { Message } from "@/types/message";
import type { ChatHistory } from "@/types/chat-history";

export default function ChatPage() {

    const initialMessage: Message = {
        id: "1",
        role: "assistant",
        content: "Hello! How can I help you today?",
    };

    const [messages, setMessages] = useState<Message[]>([
        initialMessage,
    ]);

    const [isLoading, setIsLoading] = useState(false);

    const [history, setHistory] = useState<ChatHistory[]>([]);

    async function loadHistory() {
        try {
            const response = await getChatHistory();

            setHistory(response.data);

            if (response.data.length > 0) {
                const latest = response.data[0];

                setMessages([
                    {
                        id: crypto.randomUUID(),
                        role: "user",
                        content: latest.question,
                    },
                    {
                        id: crypto.randomUUID(),
                        role: "assistant",
                        content: latest.answer,
                    },
                ]);
            }

        } catch (error) {
            console.error(
                "Failed to load chat history:",
                error
            );
        }
    }

    useEffect(() => {
        loadHistory();
    }, []);

    async function handleSend(message: string) {
        const userMessage: Message = {
            id: crypto.randomUUID(),
            role: "user",
            content: message,
        };

        setMessages((previousMessages) => [
            ...previousMessages,
            userMessage,
        ]);

        setIsLoading(true);

        try {

            const loadingId = crypto.randomUUID();

            const loadingMessage: Message = {
                id: loadingId,
                role: "assistant",
                content: "Thinking...",
            };

            setMessages((previousMessages) => [
                ...previousMessages,
                loadingMessage,
            ]);

            const response = await sendMessage({
                message,
            });

            setMessages((previousMessages) =>
                previousMessages.filter(
                    (message) => message.id !== loadingId
                )
            );

            const assistantMessage: Message = {
                id: crypto.randomUUID(),
                role: "assistant",
                content: response.answer,
            };

            setMessages((previousMessages) => [
                ...previousMessages,
                assistantMessage,
            ]);

            await loadHistory();

        } catch {
            const errorMessage: Message = {
                id: crypto.randomUUID(),
                role: "assistant",
                content: "Sorry, something went wrong.",
            };

            setMessages((previousMessages) => [
                ...previousMessages,
                errorMessage,
            ]);
        } finally {
            setIsLoading(false);
        }
    }

    function handleHistorySelect(item: ChatHistory) {

        setMessages([
            {
                id: crypto.randomUUID(),
                role: "user",
                content: item.question,
            },
            {
                id: crypto.randomUUID(),
                role: "assistant",
                content: item.answer,
            },
        ]);
    }

    function handleNewChat() {
        setMessages([
            {
                ...initialMessage,
                id: crypto.randomUUID(),
            },
        ]);
    }

    return (
        <main className="chat-page">

            <aside className="history-sidebar">
                <ChatHistoryList
                    history={history}
                    onSelect={handleHistorySelect}
                    onNewChat={handleNewChat}
                />
            </aside>

            <section className="chat-container">

                <header className="chat-header">
                    <h1>MindK AI</h1>
                </header>

                <section className="chat-content">
                    <ChatMessageList messages={messages} />
                </section>

                <footer className="chat-footer">
                    <ChatInput
                        onSend={handleSend}
                        disabled={isLoading}
                    />
                </footer>

            </section>

        </main>
    );
}
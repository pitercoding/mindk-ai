import { useState } from "react";

import ChatInput from "@/components/chat/ChatInput";
import ChatMessageList from "@/components/chat/ChatMessageList";

import { sendMessage } from "@/services/chatService";

import type { Message } from "@/types/message";

export default function ChatPage() {
    const [messages, setMessages] = useState<Message[]>([
        {
            id: "1",
            role: "assistant",
            content: "Hello! How can I help you today?",
        },
    ]);

    const [isLoading, setIsLoading] = useState(false);

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
                loadingMessage
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

    return (
        <main className="chat-page">

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

        </main>
    );
}
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

        try {
            const response = await sendMessage({
                message,
            });

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
        }
    }

    return (
        <main>
            <h1>MindK AI</h1>

            <ChatMessageList messages={messages} />

            <ChatInput onSend={handleSend} />
        </main>
    );
}
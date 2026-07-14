import { useState } from "react";

import ChatInput from "@/components/chat/ChatInput";
import ChatMessageList from "@/components/chat/ChatMessageList";

import type { Message } from "@/types/message";

export default function ChatPage() {
    const [messages, setMessages] = useState<Message[]>([
        {
            id: "1",
            role: "assistant",
            content: "Hello! How can I help you today?",
        },
    ]);

    function handleSend(message: string) {
        const newMessage: Message = {
            id: crypto.randomUUID(),
            role: "user",
            content: message,              
        };

        setMessages((previousMessages) => [
            ...previousMessages, newMessage,
        ]);
    }

    return (
        <main>
            <h1>MindK AI</h1>

            <ChatMessageList
            messages={messages}
            />

            <ChatInput onSend={handleSend} />
        </main>
    );
}
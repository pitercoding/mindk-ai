import { useEffect, useRef } from "react";

import ChatMessage from "./ChatMessage";

import type { Message } from "@/types/message";

interface ChatMessageListProps {
    messages: Message[];
}

export default function ChatMessageList({
    messages,
}: ChatMessageListProps) {

    const messagesEndRef = useRef<HTMLDivElement | null>(null);

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({
            behavior: "smooth",
        });
    }, [messages]);

    return (
        <div className="message-list">

            {messages.map((message) => (
                <ChatMessage
                    key={message.id}
                    message={message}
                />
            ))}

            <div ref={messagesEndRef} />
            
        </div>
    );
}

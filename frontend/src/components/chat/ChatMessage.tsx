import type { Message } from "@/types/message";

interface ChatMessageProps {
    message: Message;
}

export default function ChatMessage({
        message
    }: ChatMessageProps) {
    return (
        <div>
            <strong>{message.role}</strong>

            <p>{message.content}</p>
        </div>
    );
}
import { useState, type FormEvent } from "react";

interface ChatInputProps {
    onSend: (message: string) => void;
}

export default function ChatInput({
    onSend
}: ChatInputProps) {
    const [message, setMessage] = useState("");

    function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault();

        const trimmedMessage = message.trim();

        if (!trimmedMessage) {
            return;
        }

        onSend(trimmedMessage);

        setMessage("");
    }

    return (
        <form onSubmit={handleSubmit}>
            <input 
                type="text"
                placeholder="Ask something..."
                value={message}
                onChange={(event) => setMessage(event.target.value)}
            />

            <button type="submit">
                Send
            </button>
        </form>
    );
}
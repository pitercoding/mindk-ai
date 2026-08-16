import { useEffect, useRef, useState, type FormEvent } from "react";

interface ChatInputProps {
    onSend: (message: string) => void;
    disabled?: boolean;
}

export default function ChatInput({
    onSend,
    disabled,
}: ChatInputProps) {
    const [message, setMessage] = useState("");
    const textareaRef = useRef<HTMLTextAreaElement | null>(null);
    const wasDisabledRef = useRef(disabled);

    // Auto-grow within the CSS min/max-height bounds as content changes.
    useEffect(() => {
        const el = textareaRef.current;

        if (!el) {
            return;
        }

        el.style.height = "auto";
        el.style.height = `${el.scrollHeight}px`;
    }, [message]);

    // Return focus once sending/loading finishes, so typing the next
    // message doesn't require re-clicking the textarea.
    useEffect(() => {
        if (wasDisabledRef.current && !disabled) {
            textareaRef.current?.focus();
        }

        wasDisabledRef.current = disabled;
    }, [disabled]);

    function trySend() {
        const trimmedMessage = message.trim();

        if (!trimmedMessage) {
            return;
        }

        onSend(trimmedMessage);

        setMessage("");
    }

    function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault();
        trySend();
    }

    function handleKeyDown(
        event: React.KeyboardEvent<HTMLTextAreaElement>
    ) {

        if (
            event.key === "Enter" &&
            !event.shiftKey
        ) {

            event.preventDefault();
            trySend();
        }
    }

    return (
        <form
            className="chat-input"
            onSubmit={handleSubmit}>
            <textarea
                ref={textareaRef}
                placeholder="Ask something..."
                value={message}
                onChange={(event) => setMessage(event.target.value)}
                onKeyDown={handleKeyDown}
                disabled={disabled}
                rows={1}
            />

            <button
                type="submit"
                disabled={disabled}
            >
                Send
            </button>
        </form>
    );
}

import { useEffect, useState } from "react";

import Chat from "@/components/chat/Chat";

import { useSelectedNote } from "@/context/SelectedNoteContext";

import { sendMessage } from "@/services/chatService";

import type { Message } from "@/types/message";


export default function ChatPanel() {

    const initialMessage: Message = {
        id: "1",
        role: "assistant",
        content: "Hello! How can I help you today?",
    };

    const [messages, setMessages] = useState<Message[]>([
        initialMessage,
    ]);

    const [isLoading, setIsLoading] = useState(false);

    const {
        selectedNote,
    } = useSelectedNote();

    useEffect(() => {

        if (!selectedNote) {
            return;
        }

        setMessages([
            {
                id: crypto.randomUUID(),
                role: "assistant",
                content:
                    `Olá! Estou analisando a nota "${selectedNote.title}". Como posso ajudar?`,
            },
        ]);

    }, [selectedNote]);

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
                context: selectedNote
                    ? {
                        title: selectedNote.title,
                        content: selectedNote.content,
                    }
                    : undefined,
            });

            setMessages((previousMessages) =>
                previousMessages.filter(
                    (item) =>
                        item.id !== loadingId
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

        } catch (error) {

            console.error(
                "Failed to send dashboard chat message:",
                error
            );

            setMessages((previousMessages) => [
                ...previousMessages,
                {
                    id: crypto.randomUUID(),
                    role: "assistant",
                    content: "Sorry, something went wrong.",
                },
            ]);

        } finally {
            setIsLoading(false);
        }
    }

    return (

        <section className="dashboard-panel chat-panel">

            <header className="panel-header">

                <div>
                    <h2>MINDK CHAT</h2>

                    <span>Ask MindK anything about your data</span>
                </div>

                <button>...</button>

            </header>

            {selectedNote && (

                <div className="chat-context">
                    <span>Context:</span>
                    <strong> {selectedNote.title}</strong>
                </div>
            )}

            <Chat
                messages={messages}
                isLoading={isLoading}
                onSend={handleSend}
            />

        </section>
    );
}

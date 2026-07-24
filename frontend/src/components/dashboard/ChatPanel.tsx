import { useEffect, useState } from "react";

import Chat from "@/components/chat/Chat";

import { useSelectedNote } from "@/context/SelectedNoteContext";

import { sendMessage } from "@/services/chatService";
import { getMessagesByNote } from "@/services/chatMessageService";

import type { Message } from "@/types/message";


export default function ChatPanel() {

    const initialMessage: Message = {
        id: "1",
        role: "assistant",
        content: "Hello! How can I help you today?",
    };

    const [chatSessions, setChatSessions] = useState<
        Record<number, Message[]>
    >({});

    const [isLoading, setIsLoading] = useState(false);

    const {
        selectedNote,
    } = useSelectedNote();

    const messages =
        selectedNote
            ? chatSessions[selectedNote.id] ?? []
            : [];

    useEffect(() => {

        if (!selectedNote) {
            return;
        }

        const note = selectedNote;

        async function loadMessages() {

            try {

                const history =
                    (await getMessagesByNote(
                        note.id,
                    )) ?? [];

                setChatSessions(
                    (previousSessions) => ({

                        ...previousSessions,

                        [note.id]:

                            history.length > 0

                                ? history.map((message) => ({
                                    id: message.id,
                                    role: message.role,
                                    content: message.content,
                                }))

                                : [
                                    {
                                        id: crypto.randomUUID(),
                                        role: "assistant",
                                        content:
                                            `Hello! I am analyzing the note "${note.title}". How can I help?`,
                                    },
                                ],
                    }),
                );

            } catch (error) {

                console.error(
                    "Failed to load chat messages:",
                    error,
                );
            }
        }
        loadMessages();
    }, [selectedNote]);

    function updateCurrentMessages(
        updater: (messages: Message[]) => Message[]
    ) {

        if (!selectedNote) {
            return;
        }

        setChatSessions((previousSessions) => {

            const currentMessages =
                previousSessions[selectedNote.id] ??
                [initialMessage];

            return {
                ...previousSessions,
                [selectedNote.id]: updater(currentMessages),
            };

        });
    }

    async function handleSend(message: string) {

        const userMessage: Message = {
            id: crypto.randomUUID(),
            role: "user",
            content: message,
        };

        updateCurrentMessages((previousMessages) => [
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

            updateCurrentMessages((previousMessages) => [
                ...previousMessages,
                loadingMessage,
            ]);

            const response = await sendMessage({
                message,
                context: selectedNote
                    ? {
                        note_id: selectedNote.id,
                        title: selectedNote.title,
                        content: selectedNote.content,
                    }
                    : undefined,
            });

            updateCurrentMessages((previousMessages) =>
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

            updateCurrentMessages((previousMessages) => [
                ...previousMessages,
                assistantMessage,
            ]);

        } catch (error) {

            console.error(
                "Failed to send dashboard chat message:",
                error
            );

            updateCurrentMessages((previousMessages) => [
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
